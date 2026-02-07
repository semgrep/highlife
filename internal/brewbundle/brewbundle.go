package brewbundle

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// Run executes "brew bundle" with the given Brewfile path inside repoDir.
// Returns combined stdout+stderr output and any error.
func Run(repoDir, brewfilePath string) (string, error) {
	fullPath := filepath.Join(repoDir, brewfilePath)
	cmd := exec.Command("brew", "bundle", "--file="+fullPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("brew bundle: %w\n%s", err, out)
	}
	return string(out), nil
}
