package executil

import (
	"context"
	"os"
	"os/exec"

	"github.com/charmbracelet/log"
)

// RunContext executes a command with the given context, logging the command
// string at debug level and redirecting stdout/stderr to os.Stderr.
func RunContext(ctx context.Context, name string, args ...string) error {
	// nosemgrep: go.lang.security.audit.dangerous-exec-command.dangerous-exec-command
	cmd := exec.CommandContext(ctx, name, args...)
	log.Debug("exec", "cmd", cmd.String())
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Run executes a command, logging the command string at debug level and
// redirecting stdout/stderr to os.Stderr.
func Run(name string, args ...string) error {
	return RunContext(context.Background(), name, args...)
}

// RunQuiet executes a command, logging the command string at debug level
// but discarding all output.
func RunQuiet(name string, args ...string) error {
	// nosemgrep: go.lang.security.audit.dangerous-exec-command.dangerous-exec-command
	cmd := exec.Command(name, args...)
	log.Debug("exec", "cmd", cmd.String())
	return cmd.Run()
}
