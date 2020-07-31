package jira

import (
	"encoding/json"
	"fmt"
)

// Transitions is the data model for the transition API response.
type Transitions struct {
	Transitions []Transition `json:"transitions"`
}

// Transition is the data model for Jira ticket statuses.
type Transition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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

// TransitionToInProgress updates the status Jira issue to "In Progress".
func TransitionToInProgress(issueID string, c *Config) error {
	return transitionIssue("In Progress", issueID, c)
}

// TransitionToCodeReview updates the status Jira issue to "Code Review".
func TransitionToCodeReview(issueID string, c *Config) error {
	return transitionIssue("Code Review", issueID, c)
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
