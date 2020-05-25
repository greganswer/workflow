package git

import (
	"fmt"
	"os/exec"
	"strings"
)

var RepoIsDirtyErr = fmt.Errorf("repository has unstaged changes")
var NotInitializedErr = fmt.Errorf("git repository has not been initialized")

// Checkout branch by name.
func Checkout(branch string) error {
	return executeAndStream("git", "checkout", branch)
}

// CreateBranch creates a new git branch.
func CreateBranch(name string) error {
	return executeAndStream("git", "checkout", "-b", name)
}

// CurrentBranch returns the current branch for this Git repo.
func CurrentBranch() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	return strings.Trim(string(out), "\n"), err
}

// DirIsClean returns false if there are changes in the repo.
func RepoIsClean() bool {
	return exec.Command("git", "diff", "--exit-code").Run() == nil
}

// RootDir is the root directory of the Git project.
// Reference: https://stackoverflow.com/a/957978
func RootDir() string {
	out, _ := exec.Command("git", "rev-parse", "--show-toplevel").CombinedOutput()
	return strings.TrimSuffix(string(out), "\n")
}

// Pull gets new changes from the remote repo.
func Pull() error {
	return executeAndStream("git", "pull")
}
