package issues

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/greganswer/workflow/jira"
)

const branchNameMaxLength = "%.50s"

// Issue contains the issue information.
type Issue struct {
	ID     string
	Title  string
	Type   string
	Status string
	APIURL string
	WebURL string
}

// NewFromJira converts a Jira Issue to an Issue entity.
// It does this by making an HTTP request to the issue tracker API.
func NewFromJira(issueID string, c *jira.Config) (Issue, error) {
	i, err := jira.GetIssue(issueID, c)
	if err != nil {
		return Issue{}, err
	}

	return Issue{
		ID:     i.Key,
		Title:  i.Fields.Summary,
		Type:   i.Fields.IssueType.Name,
		Status: i.Fields.Status.Name,
		APIURL: i.Self,
		WebURL: joinURLPath(c.WebURL, jira.WebIssuePath, issueID),
	}, nil
}

func NewFromBranch(branch string, c *jira.Config) (Issue, error) {
	issueID := parseIssueIdFromBranch(branch)
	return NewFromJira(issueID, c)
}

// String representation of an issue.
func (i Issue) String() string {
	if i.ID != "" && i.Title != "" {
		return fmt.Sprintf("%s: %s", i.ID, i.Title)
	}
	return i.ID
}

// BranchName from issue ID and title.
// Ref: https://github.com/lakshmichandrakala/go-parameterize
func (i Issue) BranchName() string {
	reAlphaNum := regexp.MustCompile("[^A-Za-z0-9]+")
	reTrim := regexp.MustCompile("^-|-$")

	title := reAlphaNum.ReplaceAllString(i.Title, "-")
	title = reTrim.ReplaceAllString(title, "")

	id := reAlphaNum.ReplaceAllString(i.ID, "-")
	id = reTrim.ReplaceAllString(id, "")

	branch := strings.Join([]string{i.BranchPrefix() + id, title}, "-")

	return fmt.Sprintf(branchNameMaxLength, strings.ToLower(branch))
}

// BranchPrefix returns the Git flow branch prefixes based on the Issue type.
func (i Issue) BranchPrefix() string {
	switch i.Type {
	case "Story":
		return "feature-"
	case "Bug":
		return "bug-"
	default:
		return "task-"
	}
}

func joinURLPath(base string, elem ...string) string {
	u, err := url.Parse(base)
	if err != nil {
		return ""
	}
	u.Path = path.Join(append([]string{u.Path}, elem...)...)
	return u.String()
}

func parseIssueIdFromBranch(branch string) string {
	parts := strings.Split(branch, "-")
	if len(parts) < 4 {
		return ""
	}
	for i, part := range parts {
		if _, err := strconv.Atoi(part); err == nil {
			return strings.Join(parts[1:i+1], "-")
		}
	}
	return ""
}
