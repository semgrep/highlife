package launchd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"howett.net/plist"

	"github.com/semgrep/highlife/internal/executil"
	"github.com/semgrep/highlife/internal/paths"
)

const label = "com.semgrep.highlife.sync"
const filename = label + ".plist"

type launchdPlist struct {
	Label             string   `plist:"Label"`
	ProgramArguments  []string `plist:"ProgramArguments"`
	StartInterval     int      `plist:"StartInterval"`
	RunAtLoad         bool     `plist:"RunAtLoad"`
	StandardOutPath   string   `plist:"StandardOutPath"`
	StandardErrorPath string   `plist:"StandardErrorPath"`
}

// domainTarget returns the launchd domain target for the current user (e.g. "gui/501").
func domainTarget() string {
	return fmt.Sprintf("gui/%d", os.Getuid())
}

// domain returns the fully qualified launchd service domain (e.g. "gui/501/com.semgrep.highlife.sync").
func domain() string {
	return domainTarget() + "/" + label
}

func plistPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "LaunchAgents", filename)
}

type InstallOptions struct {
	IntervalMinutes     int
	BrewMinIntervalMins int
	RunAtLoad           bool
}

func Install(opts InstallOptions) error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve executable path: %w", err)
	}

	data := launchdPlist{
		Label: label,
		ProgramArguments: []string{
			exe, "--debug", "sync",
			"--connectivity-check",
			"--brew-min-interval", fmt.Sprintf("%d", opts.BrewMinIntervalMins),
		},
		StartInterval:     opts.IntervalMinutes * 60,
		RunAtLoad:         opts.RunAtLoad,
		StandardOutPath:   paths.LogFile(),
		StandardErrorPath: paths.LogFile(),
	}

	if err := writeNewsyslogConf(); err != nil {
		return fmt.Errorf("write newsyslog conf: %w", err)
	}

	dir := filepath.Dir(plistPath())
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	f, err := os.Create(plistPath())
	if err != nil {
		return err
	}
	defer f.Close()

	if err := plist.NewEncoder(f).Encode(data); err != nil {
		return err
	}

	// bootout any previously loaded version (ignore errors if not loaded).
	_ = executil.RunQuiet("launchctl", "bootout", domain())

	if err := executil.Run("launchctl", "bootstrap", domainTarget(), plistPath()); err != nil {
		return fmt.Errorf("launchctl bootstrap: %w", err)
	}

	log.Info("installed", "label", label, "path", plistPath(), "interval", fmt.Sprintf("%dm", opts.IntervalMinutes), "brew_min_interval", fmt.Sprintf("%dm", opts.BrewMinIntervalMins))
	return nil
}

func writeNewsyslogConf() error {
	conf := fmt.Sprintf("# logfilename\towner:group\tmode\tcount\tsize\twhen\tflags\n%s\t:\t644\t5\t1024\t*\tNZ\n", paths.LogFile())

	dir := filepath.Dir(paths.NewsyslogConf())
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	return os.WriteFile(paths.NewsyslogConf(), []byte(conf), 0o644)
}

// RotateLogs runs newsyslog to rotate the log file if needed.
func RotateLogs() {
	conf := paths.NewsyslogConf()
	if _, err := os.Stat(conf); err != nil {
		return
	}
	_ = executil.RunQuiet("/usr/sbin/newsyslog", "-f", conf)
}

func Uninstall() error {
	path := plistPath()

	_ = executil.RunQuiet("launchctl", "bootout", domain()) // ignore error if not loaded

	if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("remove plist: %w", err)
	}

	log.Info("uninstalled", "label", label, "path", path)
	return nil
}
