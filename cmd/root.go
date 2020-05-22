package cmd

import (
	"os/user"
	"path"

	"github.com/greganswer/workflow/file"

	"github.com/greganswer/workflow/git"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var (
	currentUser    *user.User
	globalConfig   *viper.Viper
	localConfig    *viper.Viper
	configFilename = ".workflow.yml"
	configFileType = "yaml"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "workflow",
	Version: "0.1.0",
	Short:   "Automate software development workflows using the command line",
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

	cobra.OnInitialize(initGlobalConfig, initLocalConfig)
}

// Create the global config file.
// TODO: DRY these 2 functions up.
func initGlobalConfig() {
	globalConfig = viper.New()
	globalConfig.SetConfigName(configFilename)
	globalConfig.SetConfigType(configFileType)
	configFilePath := path.Join(currentUser.HomeDir, configFilename)
	globalConfig.SetConfigFile(configFilePath)
	file.Touch(configFilePath)
	failIfError(globalConfig.ReadInConfig())
}

// Create the local config file.
func initLocalConfig() {
	localConfig = viper.New()
	localConfig.SetConfigName(configFilename)
	globalConfig.SetConfigType(configFileType)
	configFilePath := path.Join(git.RootDir(), configFilename)
	localConfig.SetConfigFile(configFilePath)
	file.Touch(configFilePath)
	failIfError(localConfig.ReadInConfig())
}
