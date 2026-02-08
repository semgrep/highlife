package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/charmbracelet/log"
	"github.com/semgrep/highlife/internal/config"
	"github.com/semgrep/highlife/internal/flock"
	"github.com/semgrep/highlife/internal/gitops"
	"github.com/semgrep/highlife/internal/paths"
)

type TrackCmd struct {
	Paths []string `arg:"" optional:"" default:"Brewfile" help:"Paths to Brewfiles within the current git repo."`
}

func (c *TrackCmd) Run(g *Globals) error {
	// Resolve each path and group by repo URL.
	type resolved struct {
		url     string
		relPath string
	}
	var items []resolved

	for _, p := range c.Paths {
		absPath, err := filepath.Abs(p)
		if err != nil {
			return fmt.Errorf("resolve path %s: %w", p, err)
		}

		info, err := os.Stat(absPath)
		if err != nil {
			return fmt.Errorf("file not found: %s", p)
		}

		if info.IsDir() {
			absPath = filepath.Join(absPath, "Brewfile")
			if _, err := os.Stat(absPath); err != nil {
				return fmt.Errorf("no Brewfile found in directory: %s", p)
			}
		}

		dir := filepath.Dir(absPath)
		root, err := gitops.RepoRoot(dir)
		if err != nil {
			return err
		}

		url, err := gitops.RemoteURL(root)
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(root, absPath)
		if err != nil {
			return fmt.Errorf("compute relative path: %w", err)
		}

		items = append(items, resolved{url: url, relPath: rel})
	}

	lf, err := flock.Acquire(paths.LockFile())
	if err != nil {
		return err
	}
	defer func() { _ = flock.Release(lf) }()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	for _, item := range items {
		if cfg.HasSource(item.url, item.relPath) {
			return fmt.Errorf("source already exists: %s %s", item.url, item.relPath)
		}
	}

	// Group paths by URL for EnsureRepo calls.
	byURL := make(map[string][]string)
	for _, item := range items {
		byURL[item.url] = append(byURL[item.url], item.relPath)
	}

	for url, relPaths := range byURL {
		allPaths := slices.Concat(cfg.PathsForURL(url), relPaths)
		dir, err := gitops.EnsureRepo(url, allPaths)
		if err != nil {
			return fmt.Errorf("fetch repo: %w", err)
		}

		for _, rp := range relPaths {
			if _, err := os.Stat(filepath.Join(dir, rp)); err != nil {
				return fmt.Errorf("file not found in repo: %s", rp)
			}
		}
	}

	for _, item := range items {
		cfg.Sources = append(cfg.Sources, config.Source{
			URL:  item.url,
			Path: item.relPath,
		})
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	for _, item := range items {
		log.Info("added", "repo", item.url, "path", item.relPath)
	}
	return nil
}
