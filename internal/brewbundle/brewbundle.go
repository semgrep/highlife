package brewbundle

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Options controls how Run executes the brew bundle command.
type Options struct {
	Debug  bool
	DryRun bool
}

// Run executes "brew bundle" with the given Brewfile path inside repoDir.
func Run(repoDir, brewfilePath string, opts Options) error {
	fullPath := filepath.Join(repoDir, brewfilePath)
	cmd := exec.Command("brew", "bundle", "--file="+fullPath)

	if opts.Debug || opts.DryRun {
		fmt.Fprintln(os.Stderr, cmd.String())
	}

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
