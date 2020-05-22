package cmd

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/greganswer/workflow/issues"

	"github.com/spf13/viper"

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
	j := newJiraConfig(globalConfig, localConfig)
	issue, err := issues.NewFromJira(args[0], j)
	failIfError(err)

	ticketAndBranchInfo(issue, "develop")

	// TODO: Handle un-staged changes
	// Then the develop branch is checked out
	// And a new branch is created using the ticket type and ID

}

func newJiraConfig(global *viper.Viper, local *viper.Viper) jira.Config {
	return jira.Config{
		Username: global.GetString(jira.UsernameConfigKey),
		Token:    global.GetString(jira.TokenConfigKey),
		APIURL:   local.GetString(jira.APIConfigKey),
		WebURL:   local.GetString(jira.WebConfigKey),
	}
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

func ticketAndBranchInfo(i issues.Issue, parent string) {
	cyan := color.New(color.FgHiCyan).SprintFunc()
	fmt.Println()

	title("  Ticket:")
	fmt.Println(cyan("    ID:"), i.ID)
	fmt.Println(cyan("    Title:"), i.Title)
	fmt.Println(cyan("    Type:"), i.Type)

	title("  Branch:")
	fmt.Println(cyan("    Name:"), i.BranchName())
	fmt.Println(cyan("    Parent:"), parent)

	fmt.Println()
}
