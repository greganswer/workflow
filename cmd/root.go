package cmd

import (
	"os/user"

	"github.com/spf13/cobra"

	"github.com/greganswer/workflow/git"
)

var currentUser *user.User

var config = &Config{}

// rootCmd represents the base command when called without any sub commands.
var rootCmd = &cobra.Command{
	Use:              "workflow",
	Version:          "0.1.0",
	Short:            "Automate software development workflows using the command line",
	PersistentPreRun: persistentPreRun,
}

// persistentPreRun runs settings before each command
func persistentPreRun(*cobra.Command, []string) {
	if git.RootDir() == "" {
		failIfError(git.NotInitializedErr)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	failIfError(rootCmd.Execute())
}

func init() {
	var err error
	currentUser, err = user.Current()
	failIfError(err)

	cobra.OnInitialize(config.init)
	rootCmd.PersistentFlags().StringP("base-branch", "b", "develop", "base branch to perform command on")
}
