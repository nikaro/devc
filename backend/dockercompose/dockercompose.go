package dockercompose

import (
	"os/exec"
	"strings"

	"git.sr.ht/nka/devc/utils"
)

func setCommand(projectName string, dockerComposeFile string, action string, args ...string) (command []string) {
	command = []string{
		"docker-compose",
		"--project-name", projectName,
		"--file", dockerComposeFile,
		action,
	}
	command = append(command, args...)

	return command
}

// GetContainers return the list of containers for the given project and docker-compose file
func GetContainers(projectName string, dockerComposeFile string) (containers []string, err error) {
	cmdArgs := setCommand(projectName, dockerComposeFile, "ps", "--quiet", "--all")
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	out, err := cmd.CombinedOutput()
	split := strings.Split(string(out), "\n")
	containers = utils.RemoveFromSlice(split, "")

	return containers, err
}

// Build run docker-compose build with the given arguments
func Build(verbose bool, projectName string, dockerComposeFile string, args ...string) (err error) {
	cmd := setCommand(projectName, dockerComposeFile, "build", args...)

	return utils.Run(cmd, verbose)
}

// Up run docker-compose up with the given arguments
func Up(verbose bool, projectName string, dockerComposeFile string, args ...string) (err error) {
	cmd := setCommand(projectName, dockerComposeFile, "up", args...)

	return utils.Run(cmd, verbose)
}

// Start run docker-compose start with the given arguments
func Start(verbose bool, projectName string, dockerComposeFile string, args ...string) (err error) {
	cmd := setCommand(projectName, dockerComposeFile, "start", args...)

	return utils.Run(cmd, verbose)
}

// Down run docker-compose down with the given arguments
func Down(verbose bool, projectName string, dockerComposeFile string, args ...string) (err error) {
	cmd := setCommand(projectName, dockerComposeFile, "down", args...)

	return utils.Run(cmd, verbose)
}

// Stop run docker-compose stop with the given arguments
func Stop(verbose bool, projectName string, dockerComposeFile string, args ...string) (err error) {
	cmd := setCommand(projectName, dockerComposeFile, "stop", args...)

	return utils.Run(cmd, verbose)
}

// List run docker-compose ps with the given arguments
func List(verbose bool, projectName string, dockerComposeFile string, args ...string) (err error) {
	cmd := setCommand(projectName, dockerComposeFile, "ps", args...)

	return utils.Run(cmd, verbose)
}

// Exec run docker-compose exec with the given arguments
func Exec(verbose bool, projectName string, dockerComposeFile string, args ...string) (err error) {
	cmd := setCommand(projectName, dockerComposeFile, "exec", args...)

	return utils.Run(cmd, verbose)
}
