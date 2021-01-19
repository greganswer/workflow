package jira

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/greganswer/workflow/issues"
)

// issueResponse is the data structure for an issue from Jira's JSON API response.
type issueResponse struct {
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
		Assignee user `json:"assignee"`
	} `json:"fields"`
}

// GetIssue returns the JSON representation of a Jira issue.
// It does this by making an HTTP request to the issue tracker API.
// Reference: https://stackoverflow.com/questions/12864302
func GetIssue(issueID string, c *Config) (issues.Issue, error) {
	fmt.Printf("Retrieving info for %s Jira issue...\n", issueID)

	var i issues.Issue
	URL := joinURLPath(c.APIURL, APIIssuePath, issueID)
	res, err := makeRequest("GET", URL, nil, c)
	if err != nil {
		return i, errors.Wrap(err, "makeRequest failed")
	}
	defer res.Body.Close()

	if !statusSuccess(res) {
		var e errorResponse
		if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
			return i, errors.Wrap(err, "decode failed")
		}
		return i, fmt.Errorf("get issue failed with %s HTTP status: %s", res.Status, e.Messages)
	}

	var data issueResponse
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return i, errors.Wrap(err, "decode failed")
	}

	return issues.Issue{
		ID:       data.Key,
		Title:    data.Fields.Summary,
		Type:     data.Fields.IssueType.Name,
		Status:   data.Fields.Status.Name,
		Assignee: data.Fields.Assignee.Name,
		APIURL:   data.Self,
		WebURL:   joinURLPath(c.WebURL, WebIssuePath, issueID),
	}, nil
}
