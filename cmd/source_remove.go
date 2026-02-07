package cmd

import (
	"fmt"

	"github.com/semgrep/highlife/internal/config"
	"github.com/semgrep/highlife/internal/gitops"
)

type SourceRemoveCmd struct {
	URL  string `arg:"" help:"Git repository URL to remove."`
	Path string `arg:"" optional:"" help:"Path to Brewfile within the repo. If omitted, all sources for the URL are removed."`
}

func (c *SourceRemoveCmd) Run(g *Globals) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	var removed []config.Source
	if c.Path != "" {
		removed = cfg.RemoveByURLAndPath(c.URL, c.Path)
	} else {
		removed = cfg.RemoveByURL(c.URL)
	}
	if len(removed) == 0 {
		return fmt.Errorf("no sources found for %s", c.URL)
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	// Clean up cached repo if URL is no longer referenced.
	if !cfg.URLStillReferenced(c.URL) {
		_ = gitops.RemoveRepo(c.URL)
	}

	for _, s := range removed {
		fmt.Printf("removed %s %s\n", s.URL, s.Path)
	}
	return nil
}
