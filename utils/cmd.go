package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Run run the given command while attaching stdin, stdout and stderr
func Run(command []string, verbose bool) error {
	if verbose {
		fmt.Println(strings.Join(command, " "))
	}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
