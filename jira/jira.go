package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/fatih/color"
)

// Config keys.
const (
	UsernameConfigKey = "jira.username"
	TokenConfigKey    = "jira.token"
	APIConfigKey      = "jira.api_url"
	WebConfigKey      = "jira.api_url"
)

// URLs.
const (
	APIInstructionsURL = "https://confluence.atlassian.com/cloud/api-tokens-938839638.html"
	APIIssuePath       = "/rest/api/3/issue"
	WebIssuePath       = "/browse"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

// Config contains Jira configuration values.
type Config struct {
	Username string
	Token    string
	APIURL   string
	WebURL   string
}

// Issue is the data structure for an issue from Jira's JSON API response.
type Issue struct {
	// The API URL for the issue.
	Self string `json:"self"`
	// The ID of the issue.
	Key    string `json:"key"`
	Fields struct {
		// The title of the issue.
		Summary   string `json:"summary"`
		IssueType struct {
			// The issue type. Example: Story, Task, Sub-Task, etc.
			Name string `json:"name"`
		} `json:"issuetype"`
		Status struct {
			// The issue status. Example: Open, Closed, In Progress, etc.
			Name string `json:"name"`
		} `json:"status"`
		Priority struct {
			// The priority. Example: P0, P3, etc.
			Name string `json:"name"`
		} `json:"priority"`
		Assignee struct {
			Email string `json:"emailAddress"`
			Name  string `json:"displayName"`
		} `json:"assignee"`
	} `json:"fields"`
}

// Error is the data structure for an error response from Jira's JSON API.
type Error struct {
	Messages []string `json:"errorMessages"`
}

// TODO: REMOVE ME
func todo(message string) {
	fmt.Println(color.YellowString("TODO:"), fmt.Sprintf("Implement jira.%s", message))
}

// GetIssue returns the JSON representation of a Jira issue.
// It does this by making an HTTP request to the issue tracker API.
// Reference: https://stackoverflow.com/questions/12864302
func GetIssue(issueID string, c *Config) (Issue, error) {
	fmt.Printf("Retrieving info for %s Jira issue...\n", issueID)
	var i Issue
	u := joinURLPath(c.APIURL, APIIssuePath, issueID)
	request, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return Issue{}, err
	}

	request.Header.Set("Content-type", "application/json")
	request.SetBasicAuth(c.Username, c.Token)
	res, err := httpClient.Do(request)
	if err != nil {
		return Issue{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var e Error
		if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
			return Issue{}, err
		}
		return Issue{}, fmt.Errorf("%s: %s", res.Status, e.Messages)
	}

	err = json.NewDecoder(res.Body).Decode(&i)
	return i, err
}

// Reference: https://developer.atlassian.com/cloud/jira/platform/rest/v3/#api-rest-api-3-issue-issueIdOrKey-transitions-post
func TransitionIssueToInProgress(issueID string, c *Config) error {
	// TODO: Set status.name to "In Progress"
	// TODO: Set assignee to c.Username
	fmt.Printf("Transitioning Jira issue %s to 'In Progress'...\n", issueID)
	todo("TransitionIssueToInProgress")
	return nil
}

func joinURLPath(base string, elem ...string) string {
	u, err := url.Parse(base)
	if err != nil {
		return ""
	}
	u.Path = path.Join(append([]string{u.Path}, elem...)...)
	return u.String()
}
