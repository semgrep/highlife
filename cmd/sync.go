package cmd

import (
	"fmt"
	"net"
	"time"

	"github.com/charmbracelet/log"

	"github.com/semgrep/highlife/internal/brewbundle"
	"github.com/semgrep/highlife/internal/config"
	"github.com/semgrep/highlife/internal/flock"
	"github.com/semgrep/highlife/internal/gitops"
	"github.com/semgrep/highlife/internal/paths"
	"github.com/semgrep/highlife/internal/state"
)

const (
	connectivityAddr    = "1.1.1.1:443"
	connectivityTimeout = 2 * time.Second
	maxErrorLength      = 1024
)

type SyncCmd struct {
	DryRun            bool `help:"Print commands without running them." name:"dry-run" default:"false"`
	BrewMinInterval   int  `help:"Skip brew bundle for unchanged Brewfiles if last successful sync was less than this many minutes ago." default:"0" name:"brew-min-interval"`
	ConnectivityCheck bool `help:"Check for internet connectivity before syncing." default:"false" name:"connectivity-check"`
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

	delayActive := c.BrewMinInterval > 0 && !prev.LastSuccessfulSync.IsZero() &&
		time.Since(prev.LastSuccessfulSync) < time.Duration(c.BrewMinInterval)*time.Minute

	if c.ConnectivityCheck {
		log.Debug("connectivity check", "addr", connectivityAddr, "timeout", connectivityTimeout)
		conn, err := net.DialTimeout("tcp", connectivityAddr, connectivityTimeout)
		if err != nil {
			log.Info("no internet connectivity, skipping sync")
			return nil
		}
		_ = conn.Close()
		log.Debug("connectivity check passed")
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
	repoHashes := make(map[string]string)
	for _, url := range cfg.SourceURLs() {
		filePaths := cfg.PathsForURL(url)
		log.Info("fetching repo", "url", url, "paths", len(filePaths))
		dir, err := gitops.EnsureRepo(url, filePaths)
		if err != nil {
			repoErrors[url] = err
			log.Error("fetch failed", "url", url, "err", err)
			continue
		}
		repoDirs[url] = dir

		hash, err := gitops.FileHash(dir, filePaths)
		if err != nil {
			log.Warn("file hash failed, will not skip", "url", url, "err", err)
		} else {
			repoHashes[url] = hash
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
			result.Error = state.TruncateError(err.Error(), maxErrorLength)
			result.Duration = time.Since(start)
			results = append(results, result)
			anyFailed = true
			continue
		}

		// Skip brew bundle if delay is active and the repo's files haven't changed.
		if hash, ok := repoHashes[src.URL]; ok && delayActive && hash == prev.FileHashes[src.URL] {
			log.Info("skipping brew bundle (unchanged)", "url", src.URL, "path", src.Path)
			result.Success = true
			result.Duration = time.Since(start)
			results = append(results, result)
			continue
		}

		log.Info("running brew bundle", "url", src.URL, "path", src.Path)
		if err := brewbundle.Run(repoDirs[src.URL], src.Path, brewbundle.Options{
			DryRun:  c.DryRun,
			Timeout: g.BrewTimeout,
		}); err != nil {
			result.Success = false
			result.Error = state.TruncateError(err.Error(), maxErrorLength)
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
		FileHashes:         repoHashes,
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
