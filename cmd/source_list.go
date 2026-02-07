package cmd

import (
	"fmt"

	"github.com/semgrep/highlife/internal/config"
)

type SourceListCmd struct{}

func (c *SourceListCmd) Run(g *Globals) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	for _, s := range cfg.Sources {
		fmt.Printf("%s\t%s\n", s.URL, s.Path)
	}
	return nil
}
