package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/term"

	"github.com/tpetr/highlife/internal/brewbundle"
	"github.com/tpetr/highlife/internal/config"
	"github.com/tpetr/highlife/internal/gitops"
	"github.com/tpetr/highlife/internal/state"
)

type SyncCmd struct{}

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

		repoDir, err := gitops.EnsureRepo(src.URL)
		if err != nil {
			result.Success = false
			result.Error = state.TruncateError(err.Error(), 1024)
			results = append(results, result)
			log.Printf("error: %s", err)
			anyFailed = true
			continue
		}

		output, err := brewbundle.Run(repoDir, src.Path)
		if err != nil {
			result.Success = false
			result.Error = state.TruncateError(err.Error(), 1024)
			log.Printf("error: %s", err)
			anyFailed = true
		} else {
			result.Success = true
			log.Printf("ok: %s %s", src.URL, src.Path)
			_ = output
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
