package cmd

import (
	"fmt"
	"path"

	"github.com/spf13/viper"

	"github.com/greganswer/workflow/file"
	"github.com/greganswer/workflow/git"
	"github.com/greganswer/workflow/github"
	"github.com/greganswer/workflow/jira"
)

type Config struct {
	Global *viper.Viper
	Local  *viper.Viper
	Jira   *jira.Config
}

type Setting struct {
	Parent         *viper.Viper
	Key            string
	Description    string
	InstructionURL string
	Label          string
}

// init initialize the configs and validates required values.
func (c *Config) init() {
	const filename = ".workflow.yml"

	c.Global = viper.New()
	c.Local = viper.New()

	c.Local.SetConfigFile(
		path.Join(git.RootDir(), filename),
	)

	c.Global.SetConfigFile(
		path.Join(currentUser.HomeDir, filename),
	)

	// TODO: configs := []*viper.Viper{c.Local, c.Global}
	configs := []*viper.Viper{c.Global}
	for _, v := range configs {
		file.Touch(v.ConfigFileUsed())
		failIfError(v.ReadInConfig())
	}

	failIfError(c.validate())
	failIfError(c.update())
	c.initJira()
}

// validate each required setting in the configs.
func (c *Config) validate() error {
	for _, s := range c.settings() {
		value := s.Parent.GetString(s.Key)
		if value == "" {
			fmt.Println(s.Description)

			if s.InstructionURL != "" && confirm("Open URL with instructions") {
				openURL(jira.APIInstructionsURL)
			}

			value, err := promptString(s.Label)
			if err != nil {
				return err
			}

			s.Parent.Set(s.Key, value)
		}
	}
	return nil
}

// settings contains the list of Setting data.
func (c *Config) settings() []Setting {
	return []Setting{
		{
			Parent:         c.Global,
			Key:            jira.UsernameConfigKey,
			Description:    "Your Jira username is required to access issue info from Jira's API.",
			InstructionURL: jira.APIInstructionsURL,
			Label:          "Jira username",
		},
		{
			Parent:      c.Global,
			Key:         jira.TokenConfigKey,
			Description: "A Jira token is required to access issue info from Jira's API.",
			Label:       "Jira token",
		},
		{
			Parent:      c.Global,
			Key:         github.UsernameConfigKey,
			Description: "Your GitHub username is required to assign pull requests.",
			Label:       "GitHub username",
		},
		// {
		// 	Parent:      c.Local,
		// 	Key:         jira.APIConfigKey,
		// 	Description: "The project Jira API URL is required to access issue info.",
		// 	Label:       "Jira API URL",
		// },
	}
}

// update the config files.
func (c *Config) update() error {
	if err := c.Global.WriteConfig(); err != nil {
		return err
	}
	// return c.Local.WriteConfig() // TODO: Get local working.
	return nil
}

// initJira from global and local configs.
func (c *Config) initJira() {
	c.Jira = &jira.Config{
		Username: c.Global.GetString(jira.UsernameConfigKey),
		Token:    c.Global.GetString(jira.TokenConfigKey),
		APIURL:   c.Local.GetString(jira.APIConfigKey),
		WebURL:   c.Local.GetString(jira.WebConfigKey),
	}
}
