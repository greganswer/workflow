package cmd

import (
	"fmt"
	"os"

	"github.com/greganswer/workflow/git"

	"github.com/fatih/color"

	"github.com/greganswer/workflow/issues"

	"github.com/spf13/viper"

	"github.com/greganswer/workflow/jira"

	"github.com/spf13/cobra"
)

// startCmd represents the start command.
var startCmd = &cobra.Command{
	Use:    "start",
	Short:  "Start your workflow with the ID of a Jira issue",
	PreRun: preRunStartCmd,
	Run:    runStartCmd,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func preRunStartCmd(cmd *cobra.Command, args []string) {
	// TODO: Enable this.
	//if !git.RepoIsClean() {
	//	failIfError(git.RepoIsDirtyErr)
	//}
}

func runStartCmd(cmd *cobra.Command, args []string) {
	c := newJiraConfig(globalConfig, localConfig)
	id := args[0]
	fmt.Printf("Retrieving info for %s...\n", id)
	issue, err := issues.NewFromJira(id, c)
	failIfError(err)

	// TODO: Use base-branch persistent flag
	baseBranch := "develop"
	displayIssueAndBranchInfo(issue, baseBranch)
	if !confirm("Create this branch") {
		os.Exit(0)
	}

	failIfError(git.Checkout(baseBranch))
	failIfError(git.Pull())
	failIfError(git.CreateBranch(issue.BranchName()))
	failIfError(jira.TransitionIssueToInProgress(issue.ID, c))
}

// newJiraConfig from global and local configs.
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

func displayIssueAndBranchInfo(i issues.Issue, parent string) {
	cyan := color.New(color.FgHiCyan).SprintFunc()
	fmt.Println()

	title("  Issue:")
	fmt.Println(cyan("    ID:"), i.ID)
	fmt.Println(cyan("    Title:"), i.Title)
	fmt.Println(cyan("    Type:"), i.Type)
	fmt.Println()

	title("  Branch:")
	fmt.Println(cyan("    Name:"), i.BranchName())
	fmt.Println(cyan("    Parent:"), parent)

	fmt.Println()
}
