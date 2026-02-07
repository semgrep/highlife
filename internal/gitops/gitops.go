package gitops

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tpetr/highlife/internal/paths"
)

// EnsureRepo clones the repo if it doesn't exist, or pulls latest if it does.
// Only the given path is checked out via sparse checkout.
// Returns the local directory path.
func EnsureRepo(url, path string, debug bool) (string, error) {
	dir := paths.RepoDir(url)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := clone(url, dir, path, debug); err != nil {
			return "", fmt.Errorf("git clone %s: %w", url, err)
		}
		return dir, nil
	}

	if err := pull(dir, path, debug); err != nil {
		return "", fmt.Errorf("git pull in %s: %w", dir, err)
	}
	return dir, nil
}

func runCmd(cmd *exec.Cmd, debug bool) error {
	if debug {
		fmt.Fprintln(os.Stderr, cmd.String())
	}
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func clone(url, dir, path string, debug bool) error {
	if err := runCmd(exec.Command("git", "clone", "--depth=1", "--filter=blob:none", "--no-checkout", "--no-recurse-submodules", url, dir), debug); err != nil {
		return err
	}

	if err := runCmd(exec.Command("git", "-C", dir, "sparse-checkout", "set", "--no-cone", "/"+path), debug); err != nil {
		return err
	}

	return runCmd(exec.Command("git", "-C", dir, "checkout", "--no-recurse-submodules"), debug)
}

func pull(dir, path string, debug bool) error {
	if err := runCmd(exec.Command("git", "-C", dir, "sparse-checkout", "set", "--no-cone", "/"+path), debug); err != nil {
		return err
	}

	return runCmd(exec.Command("git", "-C", dir, "pull", "--no-recurse-submodules", "--ff-only"), debug)
}

// RemoveAllRepos deletes all cached repo directories and returns the paths removed.
func RemoveAllRepos() ([]string, error) {
	dir := paths.ReposDir()
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
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
