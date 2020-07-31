package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
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

var httpClient *http.Client

const (
	MaxIdleConnections int = 20
	RequestTimeout     int = 5
)

func init() {
	httpClient = createHTTPClient()
}

// Config contains Jira configuration values.
type Config struct {
	Username  string
	AccountID string
	Token     string
	APIURL    string
	WebURL    string
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

// GetIssue returns the JSON representation of a Jira issue.
// It does this by making an HTTP request to the issue tracker API.
// Reference: https://stackoverflow.com/questions/12864302
func GetIssue(issueID string, c *Config) (Issue, error) {
	fmt.Printf("Retrieving info for %s Jira issue...\n", issueID)

	var i Issue
	u := joinURLPath(c.APIURL, APIIssuePath, issueID)
	res, err := makeRequest("GET", u, nil, c)
	if err != nil {
		return Issue{}, err
	}
	defer res.Body.Close()

	if !statusSuccess(res) {
		var e Error
		if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
			return Issue{}, err
		}
		return Issue{}, fmt.Errorf("%s: %s", res.Status, e.Messages)
	}

	err = json.NewDecoder(res.Body).Decode(&i)
	return i, err
}

// AssignUser assigns a user to the Jira issue.
func AssignUser(accountID, issueID string, c *Config) error {
	fmt.Printf("Assigning Jira issue %s to user with account ID %s...\n", issueID, accountID)

	reqBody, err := json.Marshal(map[string]string{
		"accountId": accountID,
	})
	if err != nil {
		return err
	}

	u := joinURLPath(c.APIURL, APIIssuePath, issueID, "assignee")
	res, err := makeRequest("PUT", u, reqBody, c)
	if err != nil {
		return err
	}

	_, err = readBody(res.Body)
	if err != nil {
		return err
	}

	if !statusSuccess(res) {
		return fmt.Errorf("%s: %s", res.Status, "unable to assign user")
	}

	return nil
}
