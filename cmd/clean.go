package cmd

import (
	"fmt"

	"github.com/charmbracelet/log"

	"github.com/semgrep/highlife/internal/gitops"
)

type CleanCmd struct{}

func (c *CleanCmd) Run(g *Globals) error {
	removed, err := gitops.RemoveAllRepos()
	if err != nil {
		return fmt.Errorf("clean repos: %w", err)
	}

	for _, p := range removed {
		log.Debug("removed repo", "path", p)
	}

	log.Debug("cleared repo cache", "count", len(removed))
	return nil
}
