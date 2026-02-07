package brewbundle

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
)

// Options controls how Run executes the brew bundle command.
type Options struct {
	DryRun  bool
	Timeout time.Duration
}

// Run executes "brew bundle" with the given Brewfile path inside repoDir.
func Run(repoDir, brewfilePath string, opts Options) error {
	fullPath := filepath.Join(repoDir, brewfilePath)

	var cmd *exec.Cmd
	if opts.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
		defer cancel()
		cmd = exec.CommandContext(ctx, "brew", "bundle", "--file="+fullPath)
	} else {
		cmd = exec.Command("brew", "bundle", "--file="+fullPath)
	}

	log.Debug("exec", "cmd", cmd.String())

	if opts.DryRun {
		return nil
	}

	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("brew bundle: %w", err)
	}
	return nil
}
