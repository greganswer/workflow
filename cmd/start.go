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
	Run: func(cmd *cobra.Command, args []string) {
		jiraToken := globalConfig.GetString(jira.TokenConfigKey)
		if jiraToken == "" {
			setJiraTokenInGlobalConfig()
		}

		// Then the Jira ticket is opened
		// And the ticket type is copied
		// And the ticket ID is copied

		// TODO: Handle un-staged changes
		// Then the develop branch is checked out
		// And a new branch is created using the ticket type and ID

	},
}

// Inform the user that a Jira token is required
// Then open their browser to the page with instructions
// And prompt the user for a Jira token
// Then set the jira token in the global config.
func setJiraTokenInGlobalConfig() {
	fmt.Println("A Jira token is required.")

	if confirm("Open URL with instructions") {
		openURL(jira.APIInstructionsURL)
	}

	jiraToken, err := promptString("Jira token")
	failIfError(err)

	globalConfig.Set(jira.TokenConfigKey, jiraToken)
	failIfError(globalConfig.WriteConfig())
}

func init() {
	rootCmd.AddCommand(startCmd)
}
