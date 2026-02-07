package cmd

import (
	"fmt"
	"time"

	"github.com/charmbracelet/log"

	"github.com/semgrep/highlife/internal/brewbundle"
	"github.com/semgrep/highlife/internal/config"
	"github.com/semgrep/highlife/internal/gitops"
	"github.com/semgrep/highlife/internal/state"
)

type SyncCmd struct {
	DryRun bool `help:"Print commands without running them." name:"dry-run" default:"false"`
	Delay  int  `help:"Skip sync if last successful sync was less than this many minutes ago." default:"0"`
}

func (c *SyncCmd) Run(g *Globals) error {
	if c.Delay > 0 {
		st, err := state.Load()
		if err != nil {
			return fmt.Errorf("load state: %w", err)
		}
		if !st.LastSync.IsZero() && allSucceeded(st.Results) {
			elapsed := time.Since(st.LastSync)
			if elapsed < time.Duration(c.Delay)*time.Minute {
				log.Info("skipping sync", "last_sync", elapsed.Round(time.Second), "delay", fmt.Sprintf("%dm", c.Delay))
				return nil
			}
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
		result := state.SourceResult{
			URL:    src.URL,
			Path:   src.Path,
			SyncAt: time.Now(),
		}

		if err, ok := repoErrors[src.URL]; ok {
			result.Success = false
			result.Error = state.TruncateError(err.Error(), 1024)
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

func allSucceeded(results []state.SourceResult) bool {
	for _, r := range results {
		if !r.Success {
			return false
		}
	}
	return len(results) > 0
}
