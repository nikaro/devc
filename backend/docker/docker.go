package docker

import (
	"crypto/md5"
	"encoding/hex"
	"os/exec"
	"path/filepath"
	"strings"

	"git.sr.ht/nka/devc/utils"
)

// GetImageName return the image for the given path
func GetImageName(path string) (name string) {
	hasher := md5.New()
	hasher.Write([]byte(path))
	md5 := hex.EncodeToString(hasher.Sum(nil))
	name = "vsc-" + filepath.Base(path) + "-" + md5

	return name
}

// GetContainer return the latest container ID based on the given path
func GetContainer(path string) (container string, err error) {
	cmdArgs := []string{"ps", "--quiet", "--all", "--filter", "label=vsch.local.folder=" + path}
	cmd := exec.Command("docker", cmdArgs...)
	out, err := cmd.CombinedOutput()
	if containers := strings.Split(string(out), "\n"); len(containers) > 0 {
		container = containers[0]
	}

	return container, err
}

// Build the image for the given Dockerfile
func Build(verbose bool, image string, dockerFile string, context string, args ...string) (err error) {
	cmdArgs := []string{"docker", "build", "--tag", image, "--file", dockerFile}
	cmdArgs = append(cmdArgs, args...)
	cmdArgs = append(cmdArgs, context)

	return utils.Run(cmdArgs, verbose)
}

// Create the container with the given image
func Create(verbose bool, image string, command []string, args ...string) (err error) {
	cmdArgs := []string{"docker", "create"}
	cmdArgs = append(cmdArgs, args...)
	cmdArgs = append(cmdArgs, image)
	if len(command) > 0 {
		cmdArgs = append(cmdArgs, command...)
	}

	return utils.Run(cmdArgs, verbose)
}

// Remove the container
func Remove(verbose bool, container string, args ...string) (err error) {
	cmdArgs := []string{"docker", "rm"}
	cmdArgs = append(cmdArgs, args...)
	cmdArgs = append(cmdArgs, container)

	return utils.Run(cmdArgs, verbose)
}

// Start the given container
func Start(verbose bool, container string, args ...string) (err error) {
	cmdArgs := []string{"docker", "start"}
	cmdArgs = append(cmdArgs, args...)
	cmdArgs = append(cmdArgs, container)

	return utils.Run(cmdArgs, verbose)
}

// Stop the given container
func Stop(verbose bool, container string, args ...string) (err error) {
	cmdArgs := []string{"docker", "stop"}
	cmdArgs = append(cmdArgs, args...)
	cmdArgs = append(cmdArgs, container)

	return utils.Run(cmdArgs, verbose)
}

// List return the list of containers based on the given path
func List(verbose bool, path string, args ...string) (err error) {
	cmdArgs := []string{"docker", "ps", "--filter", "label=vsch.local.folder=" + path}
	cmdArgs = append(cmdArgs, args...)

	return utils.Run(cmdArgs, verbose)
}

// Exec execute the given command into the given container
func Exec(verbose bool, container string, command string, args ...string) (err error) {
	cmdArgs := []string{"docker", "exec"}
	cmdArgs = append(cmdArgs, args...)
	cmdArgs = append(cmdArgs, container)
	cmdArgs = append(cmdArgs, command)

	return utils.Run(cmdArgs, verbose)
}
