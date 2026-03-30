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
// It returns the combined stdout/stderr output and any error.
func Run(repoDir, brewfilePath string, opts Options) (string, error) {
	fullPath := filepath.Join(repoDir, brewfilePath)
	args := []string{"bundle", "--file=" + fullPath}

	if opts.DryRun {
		log.Debug("dry-run, skipping", "cmd", "brew bundle", "file", fullPath)
		return "", nil
	}

	ctx := context.Background()
	var cancel context.CancelFunc
	if opts.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	output, err := executil.RunContextCapture(ctx, "brew", args...)
	if err != nil {
		return string(output), fmt.Errorf("brew bundle: %w", err)
	}
	return string(output), nil
}
