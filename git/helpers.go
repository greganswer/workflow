package git

import (
	"bufio"
	"fmt"
	"os/exec"
)

// executeAndStream executes a shell command and streams the output to the terminal.
// Reference: https://stackoverflow.com/a/45957859
func executeAndStream(name string, arg ...string) error {
	c := exec.Command(name, arg...)

	// Setup stdout and stderr.
	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := c.StderrPipe()
	if err != nil {
		return err
	}

	// Start the command.
	if err := c.Start(); err != nil {
		return err
	}

	// Stream stdout and stderr.
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	scanner = bufio.NewScanner(stderr)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return c.Wait()
}
