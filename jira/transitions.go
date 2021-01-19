package jira

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/greganswer/workflow/issues"
)

// Transition names.
const (
	inProgress = "In Progress"
	codeReview = "Code Review"
)

// Transitions is the data model for the transition API response.
type transitions struct {
	Transitions []transition `json:"transitions"`
}

// Transition is the data model for Jira ticket statuses.
type transition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// findByName searches a slice of transitions by name.
// Time: O(n) - Iterate over Transitions
// Space: O(1)
func (t *transitions) findByName(name string) (*transition, error) {
	for i := range t.Transitions {
		if t.Transitions[i].Name == name {
			return &t.Transitions[i], nil
		}
	}
	return nil, fmt.Errorf("transition not found. name: %s", name)
}

// TransitionToInProgress updates the status Jira issue to "In Progress".
func TransitionToInProgress(issue issues.Issue, c *Config) error {
	return transitionIssue(inProgress, issue, c)
}

// TransitionToCodeReview updates the status Jira issue to "Code Review".
func TransitionToCodeReview(issue issues.Issue, c *Config) error {
	return transitionIssue(codeReview, issue, c)
}

func transitionIssue(name string, issue issues.Issue, c *Config) error {
	if issue.Status == name {
		fmt.Printf("Jira issue %s status already set to '%s'\n", issue.ID, name)
		return nil
	}

	fmt.Printf("Transitioning Jira issue %s to '%s' status...\n", issue.ID, name)

	t, err := getTransitionByName(name, issue.ID, c)
	if err != nil {
		return errors.Wrap(err, "getTransitionByName failed")
	}

	reqBody, err := json.Marshal(map[string]map[string]interface{}{"transition": {"id": t.ID}})
	if err != nil {
		return errors.Wrap(err, "JSON marshal failed")
	}

	res, err := makeTransitionRequest(c, issue, reqBody)
	if err != nil {
		return errors.Wrap(err, "makeTransitionRequest failed")
	}

	resBody, err := readBody(res.Body)
	if err != nil {
		return errors.Wrap(err, "read failed")
	}

	if !statusSuccess(res) {
		// var e Error
		// if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
		// 	return err
		// }
		// return fmt.Errorf("%s: %s", res.Status, e.Messages)
		return fmt.Errorf("issue transition failed with %s status: %s", res.Status, resBody)
	}

	return nil
}

func getTransitions(issueID string, c *Config) (transitions, error) {
	fmt.Printf("Retrieving transitions for %s Jira issue...\n", issueID)

	var ts transitions
	URL := joinURLPath(c.APIURL, APIIssuePath, issueID, "transitions")
	res, err := makeRequest("GET", URL, nil, c)
	if err != nil {
		return ts, errors.Wrap(err, "makeRequest failed")
	}
	// TODO: Try to use readBody() instead
	defer res.Body.Close()

	if !statusSuccess(res) {
		var e errorResponse
		if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
			return ts, errors.Wrap(err, "decode failed")
		}
		err = fmt.Errorf("%s: %s", res.Status, e.Messages)
		return ts, errors.Wrap(err, "HTTP status not success")
	}

	err = json.NewDecoder(res.Body).Decode(&ts)

	return ts, err
}

func getTransitionByName(name string, issueID string, c *Config) (*transition, error) {
	ts, err := getTransitions(issueID, c)
	if err != nil {
		return nil, errors.Wrap(err, "getTransitions failed")
	}

	return ts.findByName(name)
}

func makeTransitionRequest(c *Config, issue issues.Issue, reqBody []byte) (*http.Response, error) {
	URL := joinURLPath(c.APIURL, APIIssuePath, issue.ID, "transitions")
	return makeRequest("POST", URL, reqBody, c)
}
