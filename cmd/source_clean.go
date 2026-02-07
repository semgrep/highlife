package cmd

import (
	"fmt"

	"github.com/semgrep/highlife/internal/gitops"
)

type SourceCleanCmd struct {
	Debug bool `help:"Print paths being deleted." default:"false"`
}

func (c *SourceCleanCmd) Run() error {
	removed, err := gitops.RemoveAllRepos()
	if err != nil {
		return fmt.Errorf("clean repos: %w", err)
	}

	if c.Debug {
		for _, p := range removed {
			fmt.Println(p)
		}
	}

	fmt.Printf("cleaned %d cached repo(s)\n", len(removed))
	return nil
}
