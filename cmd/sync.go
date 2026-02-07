package cmd

import (
	"fmt"
	"time"

	"github.com/charmbracelet/log"

	"github.com/semgrep/highlife/internal/brewbundle"
	"github.com/semgrep/highlife/internal/config"
	"github.com/semgrep/highlife/internal/flock"
	"github.com/semgrep/highlife/internal/gitops"
	"github.com/semgrep/highlife/internal/paths"
	"github.com/semgrep/highlife/internal/state"
)

type SyncCmd struct {
	DryRun bool `help:"Print commands without running them." name:"dry-run" default:"false"`
	Delay  int  `help:"Skip sync if last successful sync was less than this many minutes ago." default:"0"`
}

func (c *SyncCmd) Run(g *Globals) error {
	lf, err := flock.Acquire(paths.LockFile())
	if err != nil {
		return err
	}
	defer func() { _ = flock.Release(lf) }()

	prev, err := state.Load()
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	if c.Delay > 0 && !prev.LastSuccessfulSync.IsZero() {
		elapsed := time.Since(prev.LastSuccessfulSync)
		if elapsed < time.Duration(c.Delay)*time.Minute {
			log.Info("skipping sync", "last_successful_sync", elapsed.Round(time.Second), "delay", fmt.Sprintf("%dm", c.Delay))
			return nil
		}
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if len(cfg.Sources) == 0 {
		log.Warn("no sources configured")
		return nil
	}

	// Clone/pull each unique URL once with all its paths for sparse checkout.
	repoDirs := map[string]string{}
	repoErrors := map[string]error{}
	for _, url := range cfg.SourceURLs() {
		paths := cfg.PathsForURL(url)
		log.Info("fetching repo", "url", url, "paths", len(paths))
		dir, err := gitops.EnsureRepo(url, paths)
		if err != nil {
			repoErrors[url] = err
			log.Error("fetch failed", "url", url, "err", err)
		} else {
			repoDirs[url] = dir
		}
	}

	var results []state.SourceResult
	var anyFailed bool

	for _, src := range cfg.Sources {
		start := time.Now()
		result := state.SourceResult{
			URL:    src.URL,
			Path:   src.Path,
			SyncAt: start,
		}

		if err, ok := repoErrors[src.URL]; ok {
			result.Success = false
			result.Error = state.TruncateError(err.Error(), 1024)
			result.Duration = time.Since(start)
			results = append(results, result)
			anyFailed = true
			continue
		}

		log.Info("running brew bundle", "url", src.URL, "path", src.Path)
		if err := brewbundle.Run(repoDirs[src.URL], src.Path, brewbundle.Options{
			DryRun:  c.DryRun,
			Timeout: g.BrewTimeout,
		}); err != nil {
			result.Success = false
			result.Error = state.TruncateError(err.Error(), 1024)
			log.Error("brew bundle failed", "url", src.URL, "path", src.Path, "err", err)
			anyFailed = true
		} else {
			result.Success = true
			log.Info("brew bundle ok", "url", src.URL, "path", src.Path)
		}

		result.Duration = time.Since(start)
		results = append(results, result)
	}

	now := time.Now()
	st := &state.State{
		LastSync:           now,
		LastSuccessfulSync: prev.LastSuccessfulSync,
		Results:            results,
	}
	if !anyFailed {
		st.LastSuccessfulSync = now
	}

	if err := state.Save(st); err != nil {
		return fmt.Errorf("save state: %w", err)
	}

	if anyFailed {
		return fmt.Errorf("one or more sources failed to sync")
	}
	return nil
}
