package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/greganswer/workflow/git"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show the logs for this branch",
	Long:  "Show the Datadog logs for the preview environment based on the branch name",
	Run:   runLogsCmd,
}

func runLogsCmd(_ *cobra.Command, _ []string) {
	branch, err := git.CurrentBranch()
	branchInLowercase := strings.ToLower(branch)
	failIfError(err)

	words := strings.Split(branchInLowercase, "-")
	backslashes := strings.Join(words, "\\-")

	base := "https://app.datadoghq.com/logs"
	u, err := url.Parse(base)
	failIfError(err)

	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("cols", "service,container_name")
	q.Add("query", fmt.Sprintf("container_name:*%s*", backslashes))
	u.RawQuery = q.Encode()

	failIfError(browser.OpenURL(u.String()))
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
