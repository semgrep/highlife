package gitops

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/tpetr/highlife/internal/paths"
)

// EnsureRepo clones the repo if it doesn't exist, or pulls latest if it does.
// Returns the local directory path.
func EnsureRepo(url string) (string, error) {
	dir := paths.RepoDir(url)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := clone(url, dir); err != nil {
			return "", fmt.Errorf("git clone %s: %w", url, err)
		}
		return dir, nil
	}

	if err := pull(dir); err != nil {
		return "", fmt.Errorf("git pull in %s: %w", dir, err)
	}
	return dir, nil
}

func clone(url, dir string) error {
	cmd := exec.Command("git", "clone", "--depth=1", url, dir)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func pull(dir string) error {
	cmd := exec.Command("git", "-C", dir, "pull", "--ff-only")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RemoveRepo deletes the cached repo directory for the given URL.
func RemoveRepo(url string) error {
	dir := paths.RepoDir(url)
	return os.RemoveAll(dir)
}
