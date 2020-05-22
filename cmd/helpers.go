package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pkg/browser"

	"github.com/fatih/color"

	"github.com/manifoldco/promptui"
)

func title(s string) {
	c := color.New(color.FgMagenta, color.Bold)
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

var myClient = &http.Client{Timeout: 10 * time.Second}

// Reference: https://stackoverflow.com/a/31129967
func getJSON(url string, target interface{}) error {
	r, err := myClient.Get(url)
	failIfError(err)
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
