package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/greganswer/workflow/git"
	"github.com/greganswer/workflow/github"
	"github.com/greganswer/workflow/issues"
)

// prCmd represents the pr command
var prCmd = &cobra.Command{
	Use:    "pr",
	Short:  "Create a GitHub Pull Request for the specified branch.",
	PreRun: preRunPrCmd,
	Run:    runPrCmd,
}

func init() {
	rootCmd.AddCommand(prCmd)
}

func preRunPrCmd(cmd *cobra.Command, _ []string) {
	force, _ := cmd.Flags().GetBool("force")
	if !force && !git.RepoIsClean() {
		failIfError(git.RepoIsDirtyErr)
	}
	if !github.CLIExists() {
		fmt.Println("The 'gh' CLI app is required to execute this command.")
		if confirm("Open URL with instructions") {
			openURL(github.CLIInstallationInstructions)
		}
		os.Exit(1)
	}
}

func runPrCmd(cmd *cobra.Command, _ []string) {
	branch, err := git.CurrentBranch()
	failIfError(err)

	issue, err := issues.NewFromBranch(branch, config.Jira)
	failIfError(err)

	baseBranch, _ := cmd.Flags().GetString("base")
	pr, err := github.NewPr(issue, baseBranch)
	warnIfError(err)

	displayIssueAndPRInfo(issue, pr)

	if !confirm("Create this pull request") {
		os.Exit(1)
	}

	failIfError(pr.Create())
}

// displayIssueAndPRInfo in a nicely formatted way.
func displayIssueAndPRInfo(i issues.Issue, pr github.PullRequest) {
	cyan := color.New(color.FgHiCyan).SprintFunc()
	fmt.Println()
	displayIssueInfo(i)

	title("  Pull request:")
	fmt.Println(cyan("    Title:"), pr.Title)
	fmt.Println(cyan("    Base:"), pr.Base)
	fmt.Println(cyan("    Assignee:"), pr.Assignee)
	fmt.Println(cyan("    Template:"), pr.Template)

	fmt.Println()
}
