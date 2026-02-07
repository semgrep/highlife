package launchd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/semgrep/highlife/internal/paths"
)

const label = "com.highlife.sync"

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
        <string>--delay</string>
        <string>{{ .Delay }}</string>
    </array>
    <key>StartInterval</key>
    <integer>{{ .Interval }}</integer>
    <key>RunAtLoad</key>
    <true/>
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
	Label      string
	Executable string
	LogFile    string
	Interval   int
	Delay      int
}

type InstallOptions struct {
	IntervalMinutes int
	DelayMinutes    int
}

func Install(opts InstallOptions) error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve executable path: %w", err)
	}

	data := plistData{
		Label:      label,
		Executable: exe,
		LogFile:    paths.LogFile(),
		Interval:   opts.IntervalMinutes * 60,
		Delay:      opts.DelayMinutes,
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

	cmd := exec.Command("launchctl", "load", plistPath())
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("launchctl load: %w", err)
	}

	fmt.Fprintf(os.Stderr, "installed %s\n", plistPath())
	return nil
}

func Uninstall() error {
	path := plistPath()

	cmd := exec.Command("launchctl", "unload", path)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	_ = cmd.Run() // ignore error if not loaded

	if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("remove plist: %w", err)
	}

	fmt.Fprintf(os.Stderr, "uninstalled %s\n", path)
	return nil
}
