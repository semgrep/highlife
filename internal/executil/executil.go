package executil

import (
	"bytes"
	"context"
	"io"
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

// RunContextCapture executes a command with the given context, capturing
// combined stdout/stderr into a buffer while also streaming to os.Stderr.
func RunContextCapture(ctx context.Context, name string, args ...string) ([]byte, error) {
	// nosemgrep: go.lang.security.audit.dangerous-exec-command.dangerous-exec-command
	cmd := exec.CommandContext(ctx, name, args...)
	log.Debug("exec", "cmd", cmd.String())
	var buf bytes.Buffer
	w := io.MultiWriter(os.Stderr, &buf)
	cmd.Stdout = w
	cmd.Stderr = w
	err := cmd.Run()
	return buf.Bytes(), err
}

// RunQuiet executes a command, logging the command string at debug level
// but discarding all output.
func RunQuiet(name string, args ...string) error {
	// nosemgrep: go.lang.security.audit.dangerous-exec-command.dangerous-exec-command
	cmd := exec.Command(name, args...)
	log.Debug("exec", "cmd", cmd.String())
	return cmd.Run()
}
