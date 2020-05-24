package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/greganswer/workflow/git"
	"github.com/greganswer/workflow/issues"
	"github.com/greganswer/workflow/jira"
)

// startCmd represents the start command.
var startCmd = &cobra.Command{
	Use:    "start <issueID>",
	Short:  "Start your workflow with the ID of a Jira issue",
	PreRun: preRunStartCmd,
	Args:   validateStartCmdArgs,
	Run:    runStartCmd,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func validateStartCmdArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires the issueID argument")
	}
	return nil
}

func preRunStartCmd(*cobra.Command, []string) {
	// if !git.RepoIsClean() {
	// 	failIfError(git.RepoIsDirtyErr)
	// }
}

func runStartCmd(cmd *cobra.Command, args []string) {
	id := args[0]
	fmt.Printf("Retrieving info for %s...\n", id)
	issue, err := issues.NewFromJira(id, config.Jira)
	failIfError(err)

	baseBranch, _ := cmd.Flags().GetString("base")
	displayIssueAndBranchInfo(issue, baseBranch)
	if !confirm("Create this branch") {
		os.Exit(0)
	}

	failIfError(git.Checkout(baseBranch))
	failIfError(git.Pull())
	failIfError(git.CreateBranch(issue.BranchName()))
	failIfError(jira.TransitionIssueToInProgress(issue.ID, config.Jira))
}

// displayIssueAndBranchInfo in a nicely formatted way.
func displayIssueAndBranchInfo(i issues.Issue, base string) {
	cyan := color.New(color.FgHiCyan).SprintFunc()
	fmt.Println()

	title("  Issue:")
	fmt.Println(cyan("    ID:"), i.ID)
	fmt.Println(cyan("    Title:"), i.Title)
	fmt.Println(cyan("    Type:"), i.Type)
	fmt.Println()

	title("  Branch:")
	fmt.Println(cyan("    Name:"), i.BranchName())
	fmt.Println(cyan("    Base:"), base)

	fmt.Println()
}
