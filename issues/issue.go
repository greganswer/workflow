package issues

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const branchNameMaxLength = "%.40s"

// Issue contains the issue information.
type Issue struct {
	ID       string
	Title    string
	Type     string
	Status   string
	APIURL   string
	WebURL   string
	Assignee string
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

	branch := strings.Join([]string{i.branchPrefix() + id, title}, "-")
	shortName := fmt.Sprintf(branchNameMaxLength, strings.ToLower(branch))

	return strings.TrimSuffix(shortName, "-")
}

// BranchPrefix returns the Git flow branch prefixes based on the Issue type.
func (i Issue) branchPrefix() string {
	switch i.Type {
	case "Story":
		return "feature-"
	case "Bug":
		return "bug-"
	default:
		return "task-"
	}
}

// ParseIDFromBranch gets the Issue ID from the branch name.
func ParseIDFromBranch(branch string) string {
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
