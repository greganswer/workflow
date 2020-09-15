package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/greganswer/workflow/git"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

// previewCmd represents the logs command
var previewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Open the preview link for this branch",
	Run:   runPreviewCmd,
}

func runPreviewCmd(_ *cobra.Command, _ []string) {
	branch, err := git.CurrentBranch()
	branchInLowercase := strings.ToLower(branch)
	failIfError(err)

	URL := fmt.Sprintf("https://%s.public-preview.%s", branchInLowercase, os.Getenv("WORKFLOW_MAIN_SITE_BASE_URL"))
	failIfError(browser.OpenURL(URL))
}

func init() {
	rootCmd.AddCommand(previewCmd)
}
