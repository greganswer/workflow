package jira

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/greganswer/workflow/issues"
)

// User is a Jira user.
type user struct {
	ID   string `json:"accountId"`
	Name string `json:"displayName"`
}

// String representation of a Jira User.
func (a user) String() string {
	return fmt.Sprintf("%s (%s)", a.Name, a.ID)
}

// AssignUser assigns a user to the Jira issue.
func AssignUser(accountID string, issue issues.Issue, c *Config) error {
	u, err := findUserByID(accountID, c)
	if err != nil {
		return errors.Wrap(err, "findUserByID failed")
	}

	if issue.Assignee == u.Name {
		fmt.Printf("Jira issue %s is already assigned to %s\n", issue.ID, u)
		return nil
	}

	fmt.Printf("Assigning Jira issue %s to %s...\n", issue.ID, u)

	reqBody, err := json.Marshal(map[string]string{"accountId": u.ID})
	if err != nil {
		return errors.Wrap(err, "JSON marshal failed")
	}

	URL := joinURLPath(c.APIURL, APIIssuePath, issue.ID, "assignee")
	res, err := makeRequest("PUT", URL, reqBody, c)
	if err != nil {
		return errors.Wrap(err, "makeRequest failed")
	}

	_, err = readBody(res.Body)
	if err != nil {
		return errors.Wrap(err, "read failed")
	}

	if !statusSuccess(res) {
		return fmt.Errorf("%s: %s", res.Status, "unable to assign user")
	}

	return nil
}

func findUserByID(ID string, c *Config) (user, error) {
	fmt.Printf("Retrieving user by ID %s...\n", ID)

	var u user
	URL := joinURLPath(c.APIURL, APIUserPath, ID)
	res, err := makeRequest("GET", URL, nil, c)
	if err != nil {
		return u, errors.Wrap(err, "makeRequest failed")
	}
	defer res.Body.Close()

	if !statusSuccess(res) {
		var e errorResponse
		if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
			return u, errors.Wrap(err, "decode failed")
		}
		return u, fmt.Errorf("%s: %s", res.Status, e.Messages)
	}

	err = json.NewDecoder(res.Body).Decode(&u)
	return u, err
}
