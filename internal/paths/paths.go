package paths

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
)

const appName = "highlife"

func ConfigDir() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, appName)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", appName)
}

func StateDir() string {
	if dir := os.Getenv("XDG_STATE_HOME"); dir != "" {
		return filepath.Join(dir, appName)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "state", appName)
}

func ConfigFile() string {
	return filepath.Join(ConfigDir(), "config.json")
}

func StateFile() string {
	return filepath.Join(StateDir(), "state.json")
}

func LogFile() string {
	return filepath.Join(StateDir(), appName+".log")
}

func NewsyslogConf() string {
	return filepath.Join(StateDir(), "newsyslog.conf")
}

func ReposDir() string {
	return filepath.Join(StateDir(), "repos")
}

func LockFile() string {
	return filepath.Join(StateDir(), appName+".lock")
}

func RepoDir(gitURL string) string {
	h := sha256.Sum256([]byte(gitURL))
	return filepath.Join(ReposDir(), fmt.Sprintf("%x", h[:6]))
}
