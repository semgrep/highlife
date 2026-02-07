package state

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"time"

	"github.com/semgrep/highlife/internal/paths"
)

type SourceResult struct {
	URL     string    `json:"url"`
	Path    string    `json:"path"`
	Success bool      `json:"success"`
	Error   string    `json:"error,omitempty"`
	SyncAt  time.Time `json:"sync_at"`
}

type State struct {
	LastSync time.Time      `json:"last_sync"`
	Results  []SourceResult `json:"results"`
}

func Load() (*State, error) {
	data, err := os.ReadFile(paths.StateFile())
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return &State{}, nil
		}
		return nil, err
	}
	var st State
	if err := json.Unmarshal(data, &st); err != nil {
		return nil, err
	}
	return &st, nil
}

func Save(st *State) error {
	dir := paths.StateDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp := paths.StateFile() + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, paths.StateFile())
}

// TruncateError truncates an error message to keep state.json small.
func TruncateError(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "... (truncated)"
}
