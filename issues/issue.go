package issues

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/greganswer/workflow/jira"
)

const TitleMaxLengthForBranchName = 34

var httpClient = &http.Client{Timeout: 10 * time.Second}

// Issue contains the ticket information.
type Issue struct {
	ID     string
	Title  string
	Type   string
	APIURL string
	WebURL string
}

// NewFromJira creates an issue by making an HTTP request to the issue tracker API.
// Reference: https://stackoverflow.com/questions/12864302
func NewFromJira(issueID string, j jira.Config) (Issue, error) {
	u := joinURLPath(j.WebURL, jira.APIIssuePath, issueID)
	request, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return Issue{}, err
	}

	request.Header.Set("Content-type", "application/json")
	request.SetBasicAuth(j.Username, j.Token)
	res, err := httpClient.Do(request)
	if err != nil {
		return Issue{}, err
	}
	defer res.Body.Close()

	var jiraIssue jira.Issue
	json.NewDecoder(res.Body).Decode(&jiraIssue)

	return Issue{
		ID:     jiraIssue.Key,
		Title:  jiraIssue.Fields.Summary,
		Type:   jiraIssue.Fields.IssueType.Name,
		APIURL: jiraIssue.Self,
		WebURL: joinURLPath(j.WebURL, jira.WebIssuePath, issueID),
	}, nil
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

	return fmt.Sprintf("%.60s", strings.ToLower(branch))
}

// BranchPrefix returns the Git flow branch prefixes based on the Issue type.
func (i Issue) BranchPrefix() string {
	switch i.Type {
	case "Story":
		return "feature/"
	case "Bug":
		return "bug/"
	default:
		return "task/"
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
