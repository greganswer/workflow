package git

import (
	"os/exec"
	"strings"
)

// RootDir is the root directory of the Git project.
// Reference: https://stackoverflow.com/a/957978
func RootDir() string {
	out, _ := exec.Command("git", "rev-parse", "--show-toplevel").CombinedOutput()
	return strings.TrimSuffix(string(out), "\n")
}
