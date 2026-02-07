package config

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/semgrep/highlife/internal/paths"
)

type Source struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}

type Config struct {
	Sources []Source `json:"sources"`
}

func Load() (*Config, error) {
	data, err := os.ReadFile(paths.ConfigFile())
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return &Config{}, nil
		}
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	dir := paths.ConfigDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp := paths.ConfigFile() + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, paths.ConfigFile())
}

func (c *Config) HasSource(url, path string) bool {
	for _, s := range c.Sources {
		if s.URL == url && s.Path == path {
			return true
		}
	}
	return false
}

func (c *Config) RemoveByURL(url string) []Source {
	var removed []Source
	var kept []Source
	for _, s := range c.Sources {
		if s.URL == url {
			removed = append(removed, s)
		} else {
			kept = append(kept, s)
		}
	}
	c.Sources = kept
	return removed
}

func (c *Config) RemoveByURLAndPath(url, path string) []Source {
	var removed []Source
	var kept []Source
	for _, s := range c.Sources {
		if s.URL == url && s.Path == path {
			removed = append(removed, s)
		} else {
			kept = append(kept, s)
		}
	}
	c.Sources = kept
	return removed
}

func (c *Config) SourceURLs() []string {
	seen := map[string]bool{}
	var urls []string
	for _, s := range c.Sources {
		if !seen[s.URL] {
			seen[s.URL] = true
			urls = append(urls, s.URL)
		}
	}
	return urls
}

// PathsForURL returns all paths configured for the given URL.
func (c *Config) PathsForURL(url string) []string {
	var paths []string
	for _, s := range c.Sources {
		if s.URL == url {
			paths = append(paths, s.Path)
		}
	}
	return paths
}

// URLStillReferenced returns true if any remaining source uses this URL.
func (c *Config) URLStillReferenced(url string) bool {
	for _, s := range c.Sources {
		if s.URL == url {
			return true
		}
	}
	return false
}

func DefaultPath() string {
	return filepath.Clean("Brewfile")
}
