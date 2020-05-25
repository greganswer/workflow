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
	Draft    bool
	Title    string
	Template string
}

// NewPr create the Pull Request data structure ready to create a PR on GitHub.
func NewPr(issue issues.Issue, baseBranch string, draft bool) (PullRequest, error) {
	template := "None"
	body := ""

	exists, err := file.Exists(PRBodyTemplatePath)
	if exists {
		b, err := ioutil.ReadFile(PRBodyTemplatePath)
		if err == nil {
			template = PRBodyTemplatePath
			body = string(b)
		}
	}

	return PullRequest{
		Title:    issue.String(),
		Base:     baseBranch,
		Draft:    draft,
		Template: template,
		Body:     body,
	}, err
}

// Create a Pull Request on GitHub.
// Reference: https://cli.github.com/manual/gh_pr_create
func (p *PullRequest) Create() error {
	fmt.Println("Creating Pull Request on GitHub...")

	// -a, --assignee login   Assign a person by their login
	// -B, --base string      The branch into which you want your code merged
	// -b, --body string      Supply a body. Will prompt for one otherwise.
	// -d, --draft            Mark pull request as a draft
	// -r, --reviewer login   Request a review from someone by their login
	// -t, --title string     Supply a title. Will prompt for one otherwise.
	// -w, --web              Open the web browser to create a pull request

	return nil
}
