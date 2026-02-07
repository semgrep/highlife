package gitops

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/semgrep/highlife/internal/executil"
	"github.com/semgrep/highlife/internal/paths"
)

// EnsureRepo clones the repo if it doesn't exist, or pulls latest if it does.
// All given paths are checked out via sparse checkout.
// Returns the local directory path.
func EnsureRepo(url string, filePaths []string) (string, error) {
	dir := paths.RepoDir(url)

	if _, err := os.Stat(dir); errors.Is(err, fs.ErrNotExist) {
		if err := clone(url, dir, filePaths); err != nil {
			return "", fmt.Errorf("git clone %s: %w", url, err)
		}
		return dir, nil
	}

	if err := pull(dir, filePaths); err != nil {
		return "", fmt.Errorf("git pull in %s: %w", dir, err)
	}
	return dir, nil
}

func sparseCheckoutArgs(dir string, filePaths []string) []string {
	args := []string{"-C", dir, "sparse-checkout", "set", "--no-cone"}
	for _, p := range filePaths {
		args = append(args, "/"+p)
	}
	return args
}

func clone(url, dir string, filePaths []string) error {
	if err := executil.Run("git", "clone", "--depth=1", "--filter=blob:none", "--no-checkout", "--no-recurse-submodules", url, dir); err != nil {
		return err
	}

	if err := executil.Run("git", sparseCheckoutArgs(dir, filePaths)...); err != nil {
		return err
	}

	return executil.Run("git", "-C", dir, "checkout", "--no-recurse-submodules")
}

func pull(dir string, filePaths []string) error {
	if err := executil.Run("git", sparseCheckoutArgs(dir, filePaths)...); err != nil {
		return err
	}

	return executil.Run("git", "-C", dir, "pull", "--no-recurse-submodules", "--ff-only")
}

// RemoveAllRepos deletes all cached repo directories and returns the paths removed.
func RemoveAllRepos() ([]string, error) {
	dir := paths.ReposDir()
	entries, err := os.ReadDir(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var removed []string
	for _, e := range entries {
		p := filepath.Join(dir, e.Name())
		if err := os.RemoveAll(p); err != nil {
			return removed, err
		}
		removed = append(removed, p)
	}
	return removed, nil
}

// RemoveRepo deletes the cached repo directory for the given URL.
func RemoveRepo(url string) error {
	dir := paths.RepoDir(url)
	return os.RemoveAll(dir)
}
