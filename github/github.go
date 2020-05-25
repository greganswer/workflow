package github

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/greganswer/workflow/file"
	"github.com/greganswer/workflow/git"
	"github.com/greganswer/workflow/issues"
)

var PRBodyTemplatePath = path.Join(git.RootDir(), ".github", "PULL_REQUEST_TEMPLATE.md")

// PullRequest contains GitHub Pull Request data.
type PullRequest struct {
	Assignee string
	Base     string
	Body     string
	Title    string
	Template string
}

// NewPr create the Pull Request data structure ready to create a PR on GitHub.
func NewPr(issue issues.Issue, baseBranch string) (PullRequest, error) {
	template := "None"
	body := fmt.Sprintf("# [Issue #%s](%s)\n", issue.ID, issue.WebURL)

	exists, err := file.Exists(PRBodyTemplatePath)
	if exists {
		b, err := ioutil.ReadFile(PRBodyTemplatePath)
		if err == nil {
			template = PRBodyTemplatePath
			body += string(b)
		}
	}

	return PullRequest{
		Title:    issue.String(),
		Base:     baseBranch,
		Template: template,
		Body:     body,
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
		"--web",
	)
}
