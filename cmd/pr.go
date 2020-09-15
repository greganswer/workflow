package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/greganswer/workflow/git"
	"github.com/greganswer/workflow/github"
	"github.com/greganswer/workflow/issues"
	"github.com/greganswer/workflow/jira"
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

// TODO: Handle uncommitted changes
//		1. Prompt user to use Issue Title or a custom one
//		2. Create a commit with the Issue title
func runPrCmd(cmd *cobra.Command, _ []string) {
	branch, err := git.CurrentBranch()
	failIfError(err)

	ID := issues.ParseIDFromBranch(branch)
	issue, err := jira.GetIssue(ID, config.Jira)
	failIfError(err)

	baseBranch, _ := cmd.Flags().GetString("base")
	reviewers := os.Getenv("WORKFLOW_PR_REVIEWERS")
	pr, err := github.NewPr(issue, baseBranch, reviewers)
	warnIfError(err)

	displayIssueAndPRInfo(issue, pr)

	if !confirm("Create this pull request") {
		os.Exit(1)
	}

	failIfError(pr.Create())
	failIfError(github.OpenPR(branch))
	failIfError(jira.TransitionToCodeReview(issue, config.Jira))
}

// displayIssueAndPRInfo in a nicely formatted way.
func displayIssueAndPRInfo(i issues.Issue, pr github.PullRequest) {
	cyan := color.New(color.FgHiCyan).SprintFunc()
	fmt.Println()
	displayIssueInfo(i)

	title("  Pull request:")
	fmt.Println(cyan("    Title:"), pr.Title)
	fmt.Println(cyan("    Base:"), pr.Base)
	fmt.Println(cyan("    Reviewers:"), pr.Reviewers)
	fmt.Println(cyan("    Template:"), pr.Template)

	fmt.Println()
}
