package cmd

import (
	"fmt"

	"github.com/charmbracelet/log"

	"github.com/semgrep/highlife/internal/flock"
	"github.com/semgrep/highlife/internal/gitops"
	"github.com/semgrep/highlife/internal/paths"
)

type CleanCmd struct{}

func (c *CleanCmd) Run(g *Globals) error {
	lf, err := flock.Acquire(paths.LockFile())
	if err != nil {
		return err
	}
	defer func() { _ = flock.Release(lf) }()

	removed, err := gitops.RemoveAllRepos()
	if err != nil {
		return fmt.Errorf("clean repos: %w", err)
	}

	for _, p := range removed {
		log.Debug("removed repo", "path", p)
	}

	log.Info("cleared repo cache", "count", len(removed))
	return nil
}
