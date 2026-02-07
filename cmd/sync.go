package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/term"

	"github.com/semgrep/highlife/internal/brewbundle"
	"github.com/semgrep/highlife/internal/config"
	"github.com/semgrep/highlife/internal/gitops"
	"github.com/semgrep/highlife/internal/state"
)

type SyncCmd struct {
	Debug  bool `help:"Print commands before running them." default:"false"`
	DryRun bool `help:"Print commands without running them." name:"dry-run" default:"false"`
}

func (c *SyncCmd) Run() error {
	if term.IsTerminal(int(os.Stderr.Fd())) {
		log.SetOutput(os.Stderr)
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if len(cfg.Sources) == 0 {
		log.Println("no sources configured")
		return nil
	}

	var results []state.SourceResult
	var anyFailed bool

	for _, src := range cfg.Sources {
		log.Printf("syncing %s %s", src.URL, src.Path)

		result := state.SourceResult{
			URL:    src.URL,
			Path:   src.Path,
			SyncAt: time.Now(),
		}

		repoDir, err := gitops.EnsureRepo(src.URL, src.Path, c.Debug || c.DryRun)
		if err != nil {
			result.Success = false
			result.Error = state.TruncateError(err.Error(), 1024)
			results = append(results, result)
			log.Printf("error: %s", err)
			anyFailed = true
			continue
		}

		if err := brewbundle.Run(repoDir, src.Path, brewbundle.Options{
			Debug:  c.Debug,
			DryRun: c.DryRun,
		}); err != nil {
			result.Success = false
			result.Error = state.TruncateError(err.Error(), 1024)
			log.Printf("error: %s", err)
			anyFailed = true
		} else {
			result.Success = true
			log.Printf("ok: %s %s", src.URL, src.Path)
		}

		results = append(results, result)
	}

	st := &state.State{
		LastSync: time.Now(),
		Results:  results,
	}
	if err := state.Save(st); err != nil {
		return fmt.Errorf("save state: %w", err)
	}

	if anyFailed {
		return fmt.Errorf("one or more sources failed to sync")
	}
	return nil
}
