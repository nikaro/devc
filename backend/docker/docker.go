package docker

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/nikaro/devc/utils"
)

// Docker type
type Docker struct {
	verbose    bool
	path       string
	image      string
	dockerfile string
	context    string
	container  string
	command    []string
	args       []string
	run        func(command []string, verbose bool) error
	runOut     func(command []string) (string, error)
}

// New returns initialized Docker instance
func New() *Docker {
	d := new(Docker)
	d.run = utils.Run
	d.runOut = utils.RunOut
	d.SetPath("")
	d.SetImage("")
	d.SetContainer("")

	return d
}

// SetVerbose set the context path
func (d *Docker) SetVerbose(verbose bool) {
	d.verbose = verbose
}

// SetPath set the path
func (d *Docker) SetPath(path string) {
	if path != "" {
		d.path = path
	} else {
		d.path, _ = os.Getwd()
	}
}

// GetPath get the image name
func (d *Docker) GetPath() (path string) {
	return d.path
}

// SetImage set the image name
func (d *Docker) SetImage(image string) {
	if image != "" {
		d.image = image
	} else {
		d.image = "vsc-" + filepath.Base(d.path) + "-" + utils.Md5Sum(d.path)
	}
}

// GetImage get the image name
func (d *Docker) GetImage() (image string) {
	return d.image
}

// SetContainer set the latest container ID based on the path
func (d *Docker) SetContainer(container string) {
	if container != "" {
		d.container = container
	} else {
		cmd := []string{
			"docker",
			"ps",
			"--quiet",
			"--all",
			"--filter", "label=vsch.local.folder=" + d.path,
		}
		out, _ := d.runOut(cmd)
		if containers := strings.Split(out, "\n"); len(containers) > 0 {
			d.container = containers[0]
		}
	}
}

// GetContainer get the container ID
func (d *Docker) GetContainer() (container string) {
	return d.container
}

// SetDockerfile set the Dockerfile path
func (d *Docker) SetDockerfile(path string) {
	d.dockerfile = path
}

// SetContext set the context path
func (d *Docker) SetContext(path string) {
	d.context = path
}

// SetCommand set the command
func (d *Docker) SetCommand(command []string) {
	d.command = command
}

// SetArgs set the arguments list
func (d *Docker) SetArgs(args []string) {
	d.args = args
}

// Build the image for the given Dockerfile
func (d *Docker) Build() (err error) {
	cmdArgs := []string{"docker", "build", "--tag", d.image, "--file", d.dockerfile}
	cmdArgs = append(cmdArgs, d.args...)
	cmdArgs = append(cmdArgs, d.context)

	return d.run(cmdArgs, d.verbose)
}

// Create the container with the given image
func (d *Docker) Create() (err error) {
	cmdArgs := []string{"docker", "create"}
	cmdArgs = append(cmdArgs, d.args...)
	cmdArgs = append(cmdArgs, d.image)
	if len(d.command) > 0 {
		cmdArgs = append(cmdArgs, d.command...)
	}

	return d.run(cmdArgs, d.verbose)
}

// Remove the container
func (d *Docker) Remove() (err error) {
	cmdArgs := []string{"docker", "rm"}
	cmdArgs = append(cmdArgs, d.args...)
	cmdArgs = append(cmdArgs, d.container)

	return d.run(cmdArgs, d.verbose)
}

// Start the given container
func (d *Docker) Start() (err error) {
	cmdArgs := []string{"docker", "start"}
	cmdArgs = append(cmdArgs, d.args...)
	cmdArgs = append(cmdArgs, d.container)

	return d.run(cmdArgs, d.verbose)
}

// Stop the given container
func (d *Docker) Stop() (err error) {
	cmdArgs := []string{"docker", "stop"}
	cmdArgs = append(cmdArgs, d.args...)
	cmdArgs = append(cmdArgs, d.container)

	return d.run(cmdArgs, d.verbose)
}

// List return the list of containers based on the given path
func (d *Docker) List() (err error) {
	cmdArgs := []string{"docker", "ps", "--filter", "label=vsch.local.folder=" + d.path}
	cmdArgs = append(cmdArgs, d.args...)

	return d.run(cmdArgs, d.verbose)
}

// Exec execute the given command into the given container
func (d *Docker) Exec() (err error) {
	cmdArgs := []string{"docker", "exec"}
	cmdArgs = append(cmdArgs, d.args...)
	cmdArgs = append(cmdArgs, d.container)
	cmdArgs = append(cmdArgs, d.command...)

	return d.run(cmdArgs, d.verbose)
}
