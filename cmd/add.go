package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/charmbracelet/log"
	"github.com/semgrep/highlife/internal/config"
	"github.com/semgrep/highlife/internal/gitops"
)

type AddCmd struct {
	URL   string   `arg:"" help:"Git repository URL."`
	Paths []string `arg:"" optional:"" default:"Brewfile" help:"Paths to Brewfiles within the repo."`
}

func (c *AddCmd) Run(g *Globals) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	for _, p := range c.Paths {
		if cfg.HasSource(c.URL, p) {
			return fmt.Errorf("source already exists: %s %s", c.URL, p)
		}
	}

	// Eagerly clone/pull to verify the files exist before saving.
	// Include any existing paths for this URL so sparse checkout covers everything.
	allPaths := slices.Concat(cfg.PathsForURL(c.URL), c.Paths)
	dir, err := gitops.EnsureRepo(c.URL, allPaths)
	if err != nil {
		return fmt.Errorf("fetch repo: %w", err)
	}

	for _, p := range c.Paths {
		if _, err := os.Stat(filepath.Join(dir, p)); err != nil {
			return fmt.Errorf("file not found in repo: %s", p)
		}
	}

	for _, p := range c.Paths {
		cfg.Sources = append(cfg.Sources, config.Source{
			URL:  c.URL,
			Path: p,
		})
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	for _, p := range c.Paths {
		log.Info("added", "repo", c.URL, "path", p)
	}
	return nil
}
