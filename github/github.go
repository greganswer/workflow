package github

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"

	"github.com/greganswer/workflow/file"
	"github.com/greganswer/workflow/git"
	"github.com/greganswer/workflow/issues"
)

const CLIInstallationInstructions = "https://cli.github.com"
const UsernameConfigKey = "github.username"

var PRBodyTemplatePath = path.Join(git.RootDir(), ".github", "PULL_REQUEST_TEMPLATE.md")

// PullRequest contains GitHub Pull Request data.
type PullRequest struct {
	Assignee  string
	Reviewers string
	Base      string
	Body      string
	Title     string
	Template  string
}

// NewPr create the Pull Request data structure ready to create a PR on GitHub.
func NewPr(issue issues.Issue, baseBranch, assignee, reviewers string) (PullRequest, error) {
	template := "None"
	body := fmt.Sprintf("## [Issue #%s](%s)\n\n", issue.ID, issue.WebURL)

	exists, err := file.Exists(PRBodyTemplatePath)
	if exists {
		b, err := ioutil.ReadFile(PRBodyTemplatePath)
		if err == nil {
			template = PRBodyTemplatePath
			body += string(b)
		}
	}

	return PullRequest{
		Title:     issue.String(),
		Base:      baseBranch,
		Template:  template,
		Body:      body,
		Assignee:  assignee,
		Reviewers: reviewers,
	}, err
}

// Create a Pull Request on GitHub.
// Reference: https://cli.github.com/manual/gh_pr_create
func (p *PullRequest) Create() error {
	fmt.Println("Creating Pull Request on GitHub...")

	return executeAndStream("gh", "pr", "create",
		"--base", p.Base,
		"--title", p.Title,
		"--body", p.Body,
		"--assignee", p.Assignee,
		"--reviewer", p.Reviewers,
	)
}

// CLIExists returns true if the "gh" app exists.
func CLIExists() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

// Open the PR for the given branch.
// Reference: https://cli.github.com/manual/gh_pr_view
func OpenPR(branch string) error {
	return executeAndStream("gh", "pr", "view", branch, "--web")
}
