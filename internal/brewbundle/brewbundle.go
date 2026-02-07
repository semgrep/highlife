package brewbundle

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"

	"github.com/semgrep/highlife/internal/executil"
)

// Options controls how Run executes the brew bundle command.
type Options struct {
	DryRun  bool
	Timeout time.Duration
}

// Run executes "brew bundle" with the given Brewfile path inside repoDir.
func Run(repoDir, brewfilePath string, opts Options) error {
	fullPath := filepath.Join(repoDir, brewfilePath)
	args := []string{"bundle", "--file=" + fullPath}

	if opts.DryRun {
		log.Debug("exec (dry-run)", "cmd", "brew "+fmt.Sprint(args))
		return nil
	}

	ctx := context.Background()
	var cancel context.CancelFunc
	if opts.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	if err := executil.RunContext(ctx, "brew", args...); err != nil {
		return fmt.Errorf("brew bundle: %w", err)
	}
	return nil
}
