package launchd

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/semgrep/highlife/internal/executil"
	"github.com/semgrep/highlife/internal/paths"
)

const label = "com.semgre.highlife.sync"

var plistTemplate = template.Must(template.New("plist").Parse(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>{{ .Label }}</string>
    <key>ProgramArguments</key>
    <array>
        <string>{{ .Executable }}</string>
        <string>sync</string>
        <string>--brew-min-interval</string>
        <string>{{ .BrewMinInterval }}</string>
    </array>
    <key>StartInterval</key>
    <integer>{{ .Interval }}</integer>
    <key>RunAtLoad</key>
    <{{ .RunAtLoad }}/>
    <key>StandardOutPath</key>
    <string>{{ .LogFile }}</string>
    <key>StandardErrorPath</key>
    <string>{{ .LogFile }}</string>
</dict>
</plist>
`))

func plistPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "LaunchAgents", label+".plist")
}

type plistData struct {
	Label           string
	Executable      string
	LogFile         string
	Interval        int
	BrewMinInterval int
	RunAtLoad       string
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

	runAtLoad := "false"
	if opts.RunAtLoad {
		runAtLoad = "true"
	}

	data := plistData{
		Label:      label,
		Executable: exe,
		LogFile:    paths.LogFile(),
		Interval:   opts.IntervalMinutes * 60,
		BrewMinInterval: opts.BrewMinIntervalMins,
		RunAtLoad:  runAtLoad,
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

	if err := plistTemplate.Execute(f, data); err != nil {
		return err
	}

	if err := executil.Run("launchctl", "load", plistPath()); err != nil {
		return fmt.Errorf("launchctl load: %w", err)
	}

	log.Info("installed", "path", plistPath())
	return nil
}

func Uninstall() error {
	path := plistPath()

	_ = executil.Run("launchctl", "unload", path) // ignore error if not loaded

	if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("remove plist: %w", err)
	}

	log.Info("uninstalled", "path", path)
	return nil
}
