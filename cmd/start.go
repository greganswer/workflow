package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/greganswer/workflow/git"
	"github.com/greganswer/workflow/issues"
	"github.com/greganswer/workflow/jira"
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
	id := args[0]
	fmt.Printf("Retrieving info for %s...\n", id)
	issue, err := issues.NewFromJira(id, config.Jira)
	failIfError(err)

	baseBranch, _ := cmd.Flags().GetString("base-branch")
	displayIssueAndBranchInfo(issue, baseBranch)
	if !confirm("Create this branch") {
		os.Exit(0)
	}

	failIfError(git.Checkout(baseBranch))
	failIfError(git.Pull())
	failIfError(git.CreateBranch(issue.BranchName()))
	failIfError(jira.TransitionIssueToInProgress(issue.ID, config.Jira))
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

// displayIssueAndBranchInfo in a nicely formatted way.
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
