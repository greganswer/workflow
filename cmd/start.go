package cmd

import (
	"fmt"

	"github.com/greganswer/workflow/jira"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start your workflow with the ID of a Jira ticket",
	Run:   runStartCmd,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

// TODO: On app initialize, validate all Config info
func runStartCmd(cmd *cobra.Command, args []string) {
	j, err := jira.NewConfig(globalConfig, localConfig)
	failIfError(err)
	// jira.Config{
	// 		Username string
	// 		Token string
	// 		APIURL string
	// }

	issue, err := issues.FromJira(args[0], j) // FromJira(ID string, j jira.Config)
	failIfError(err)

	// Issue URL is created like this:
	// 		fmt.Sprintf(jira.IssueURLFormat, j.APIURL, ID)
	// And stored in the Issue struct below:
	// Issue {
	// 		ID string
	// 		Title string
	// 		Type string
	//		APIURL string
	// 		WebURL string
	// }

	issue.BranchName()
	// TODO: Create some kind of validator.
	// if err, ask user to enter a shorter title or blank

	issue.URL

	jiraToken := getJiraToken()
	jiraAPIURL := getJiraAPIURL()
	jira.Username()
	jira.Token()
	issueURL := fmt.Sprintf("%s/rest/api/3/issue/%s", jira.API(), args[0])

	// Then the Jira ticket is opened
	getJSON(issueURL)
	// And the ticket type is copied
	// And the ticket ID is copied

	// TODO: Handle un-staged changes
	// Then the develop branch is checked out
	// And a new branch is created using the ticket type and ID

}

// Inform the user that a Jira token is required
// Then open their browser to the page with instructions
// And prompt the user for a Jira token
// Then set the jira token in the global config.
func getJiraToken() string {
	token := globalConfig.GetString(jira.TokenConfigKey)
	if token == "" {
		fmt.Println("A Jira token is required.")

		if confirm("Open URL with instructions") {
			openURL(jira.APIInstructionsURL)
		}

		token, err := promptString("Jira token")
		failIfError(err)

		globalConfig.Set(jira.TokenConfigKey, token)
		failIfError(globalConfig.WriteConfig())
	}
	return token
}

func getJiraAPIURL() string {
	URL := globalConfig.GetString(jira.APIConfigKey)
	if URL == "" {
		fmt.Println("The Jira API URL is required.")

		// TODO: Find instruction page.
		//if confirm("Open URL with instructions") {
		//	openURL(jira.APIInstructionsURL)
		//}

		URL, err := promptString("Jira API URL")
		failIfError(err)

		localConfig.Set(jira.APIConfigKey, URL)
		failIfError(globalConfig.WriteConfig())
	}
	return URL
}
