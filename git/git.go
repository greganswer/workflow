package git

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

var RepoIsDirtyErr = fmt.Errorf("repository has unstaged changes")

// TODO: REMOVE ME
func todo(message string) {
	fmt.Println(color.YellowString("TODO:"), fmt.Sprintf("Implement git.%s", message))
}

// RootDir is the root directory of the Git project.
// Reference: https://stackoverflow.com/a/957978
func RootDir() string {
	out, _ := exec.Command("git", "rev-parse", "--show-toplevel").CombinedOutput()
	return strings.TrimSuffix(string(out), "\n")
}

// Checkout branch by name.
func Checkout(branch string) error {
	todo("Checkout")
	return nil
	// return exec.Command("git", "checkout", branch).Run()
}

// CreateBranch creates a new git branch.
func CreateBranch(name string) error {
	todo("CreateBranch")
	return nil
	// return exec.Command("git", "checkout", "-b", name).Run()
}

// DirIsClean returns false if there are changes in the repo.
func RepoIsClean() bool {
	return exec.Command("git", "diff", "--exit-code").Run() == nil
}

// Pull gets new changes from the remote repo.
func Pull() error {
	todo("Pull")
	return nil
	// return exec.Command("git", "pull").Run()
}
