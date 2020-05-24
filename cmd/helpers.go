package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/pkg/browser"
)

func title(s string) {
	c := color.New(color.FgHiMagenta, color.Bold)
	c.Println(s)
}

func confirm(label string) bool {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}
	_, err := prompt.Run()
	return err == nil
}

func promptString(label string) (string, error) {
	validate := func(input string) error {
		if len(strings.TrimSpace(input)) < 1 {
			return fmt.Errorf("%s cannot be blank", label)
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}
	return prompt.Run()
}

// failIfError exits the program with a standardized error message if an error occurred.
func failIfError(err error) {
	if err != nil {
		red := color.New(color.FgRed, color.Bold).SprintFunc()
		os.Stderr.WriteString(fmt.Sprint(red("FAIL: "), err, "\n"))
		os.Exit(1)
	}
}

func openURL(URL string) {
	failIfError(browser.OpenURL(URL))
}
