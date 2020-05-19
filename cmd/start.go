package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start your workflow with the ID of a Jira ticket",
	Run: func(cmd *cobra.Command, args []string) {
		// Given there is no config file
		// When this command is executed
		// Then their browser is opened to the page with instructions
		fmt.Println("A Jira token is required.")
		if confirm("Open URL with instructions") {
			openURL("https://confluence.atlassian.com/cloud/api-tokens-938839638.html")
		}

		// And the user is prompted for a Jira token
		_, err := promptString("Jira token")
		failIfError(err)

		// And the config file is created in the root of the project
		// And the jira token is set in the config file
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
