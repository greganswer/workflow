package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"
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

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
		Timeout: time.Duration(RequestTimeout) * time.Second,
	}

	return client
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

// Transitions is the data model for the transition API response.
type Transitions struct {
	Transitions []Transition `json:"transitions"`
}

// Transition is the data model for Jira ticket statuses.
type Transition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Error is the data structure for an error response from Jira's JSON API.
type Error struct {
	Messages []string `json:"errorMessages"`
}

// findByName searches a slice of transitions by name.
// Time: O(n) - Iterate over Transitions
// Space: O(1)
func (t *Transitions) findByName(name string) (*Transition, error) {
	for i := range t.Transitions {
		if t.Transitions[i].Name == name {
			fmt.Printf("Found '%s' transition\n", name)
			return &t.Transitions[i], nil
		}
	}
	return nil, fmt.Errorf("transition not found. name: %s", name)
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

// TransitionToInProgress updates the status Jira issue to "In Progress".
func TransitionToInProgress(issueID string, c *Config) error {
	return transitionIssue("In Progress", issueID, c)
}

// TransitionToCodeReview updates the status Jira issue to "Code Review".
func TransitionToCodeReview(issueID string, c *Config) error {
	return transitionIssue("Code Review", issueID, c)
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

func transitionIssue(name, issueID string, c *Config) error {
	fmt.Printf("Transitioning Jira issue %s to '%s' status...\n", issueID, name)

	t, err := getTransitionByName(name, issueID, c)
	if err != nil {
		return err
	}

	reqBody, err := json.Marshal(map[string]map[string]interface{}{
		"transition": {
			"id": t.ID,
		},
	})
	if err != nil {
		return err
	}

	u := joinURLPath(c.APIURL, APIIssuePath, issueID, "transitions")
	res, err := makeRequest("POST", u, reqBody, c)
	if err != nil {
		return err
	}

	resBody, err := readBody(res.Body)
	if err != nil {
		return err
	}

	if !statusSuccess(res) {
		// var e Error
		// if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
		// 	return err
		// }
		// return fmt.Errorf("%s: %s", res.Status, e.Messages)
		return fmt.Errorf("%s: %s", res.Status, resBody)

	}

	return nil
}

func getTransitions(issueID string, c *Config) (Transitions, error) {
	fmt.Printf("Retrieving transitions for %s Jira issue...\n", issueID)

	var ts Transitions
	u := joinURLPath(c.APIURL, APIIssuePath, issueID, "transitions")
	res, err := makeRequest("GET", u, nil, c)
	if err != nil {
		return ts, err
	}
	// TODO: Try to use readBody() instead
	defer res.Body.Close()

	if !statusSuccess(res) {
		var e Error
		if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
			return ts, err
		}
		err = fmt.Errorf("%s: %s", res.Status, e.Messages)
		return ts, err
	}

	err = json.NewDecoder(res.Body).Decode(&ts)

	return ts, err
}

func getTransitionByName(name string, issueID string, c *Config) (*Transition, error) {
	ts, err := getTransitions(issueID, c)
	if err != nil {
		return nil, err
	}

	t, err := ts.findByName(name)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func makeRequest(method, u string, reqBody []byte, c *Config) (*http.Response, error) {
	req, err := http.NewRequest(method, u, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-type", "application/json")
	req.SetBasicAuth(c.Username, c.Token)
	return httpClient.Do(req)
}

func statusSuccess(res *http.Response) bool {
	return (res.StatusCode >= 200) || (res.StatusCode < 400)
}

// Read the body of the response to force the connection to close.
// Ref: https://stackoverflow.com/a/53589787
func readBody(readCloser io.ReadCloser) ([]byte, error) {
	defer readCloser.Close()
	body, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func joinURLPath(base string, elem ...string) string {
	u, err := url.Parse(base)
	if err != nil {
		return ""
	}
	u.Path = path.Join(append([]string{u.Path}, elem...)...)
	return u.String()
}
