package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Run runs the given command while attaching stdin, stdout and stderr
func Run(command []string, verbose bool) (err error) {
	if verbose {
		fmt.Println(strings.Join(command, " "))
	}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	return err
}

// RunOut runs the given command and return stdout as string
func RunOut(command []string) (out string, err error) {
	var stdout []byte
	cmd := exec.Command(command[0], command[1:]...)
	stdout, err = cmd.Output()

	return string(stdout), err
}
