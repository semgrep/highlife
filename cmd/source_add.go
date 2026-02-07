package cmd

import (
	"fmt"

	"github.com/semgrep/highlife/internal/config"
)

type SourceAddCmd struct {
	URL  string `arg:"" help:"Git repository URL."`
	Path string `optional:"" default:"Brewfile" help:"Path to Brewfile within the repo."`
}

func (c *SourceAddCmd) Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if cfg.HasSource(c.URL, c.Path) {
		return fmt.Errorf("source already exists: %s %s", c.URL, c.Path)
	}

	cfg.Sources = append(cfg.Sources, config.Source{
		URL:  c.URL,
		Path: c.Path,
	})

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	fmt.Printf("added %s %s\n", c.URL, c.Path)
	return nil
}
