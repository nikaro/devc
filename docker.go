package main

import (
	"os"
	"strings"

	"github.com/samber/lo"
)

type Docker struct {
	Args          []string
	Command       []string
	Container     string
	ContainerUser string
	Envs          []string
	Image         string
	ImageBuild    DockerImageBuild
	Mounts        []string
	Path          string
	Ports         []string
	RemoteEnvs    []string
	RemoteUser    string
	Running       bool
	WorkDir       string
}

type DockerImageBuild struct {
	Args       []string
	Dockerfile string
	CacheFrom  string
	Context    string
	Tag        string
	Target     string
}

// Init initialize docker settings
func (d *Docker) Init(config *DevContainer) error {
	if err := os.Chdir(".devcontainer"); err != nil {
		return err
	}

	d.Args = config.JSON.RunArgs
	d.Command = lo.Ternary(config.JSON.OverrideCommand, []string{"/bin/sh", "-c", "while sleep 1000; do :; done"}, nil)
	d.ContainerUser = config.JSON.ContainerUser
	d.Envs = lo.MapToSlice(config.JSON.ContainerEnv, func(k string, v string) string { return k + "=" + v })
	d.Image = lo.Ternary(config.JSON.Image != "", config.JSON.Image, "vsc-"+config.WorkingDirectoryName+"-"+md5sum(config.WorkingDirectoryPath))
	d.Image = "vsc-" + config.WorkingDirectoryName + "-" + md5sum(config.WorkingDirectoryPath)
	d.ImageBuild.Args = lo.MapToSlice(config.JSON.Build.Args, func(k string, v string) string { return k + "=" + v })
	d.ImageBuild.CacheFrom = config.JSON.Build.CacheFrom
	d.ImageBuild.Context = config.JSON.Build.Context
	d.ImageBuild.Dockerfile = config.JSON.Build.Dockerfile
	d.ImageBuild.Target = config.JSON.Build.Target
	d.Mounts = config.JSON.Mounts
	d.Mounts = append(d.Mounts, config.JSON.WorkspaceMount)
	d.Path = config.WorkingDirectoryPath
	d.Ports = config.JSON.ForwardPorts
	d.RemoteEnvs = lo.MapToSlice(config.JSON.RemoteEnv, func(k string, v string) string { return k + "=" + v })
	d.RemoteUser = config.JSON.RemoteUser
	d.WorkDir = config.JSON.WorkspaceFolder

	// check if already created
	if container, err := d.GetContainer(); err != nil {
		return err
	} else if container != "" {
		d.Container = strings.TrimSpace(container)
	}

	// check if already started
	if running, err := d.IsRunning(); err != nil {
		return err
	} else if running {
		d.Running = true
	}

	return nil
}

// GetContainer return the container name
func (d *Docker) GetContainer() (string, error) {
	cmdArgs := []string{"docker", "container", "ls"}
	cmdArgs = append(cmdArgs, "--quiet")
	cmdArgs = append(cmdArgs, "--latest")
	cmdArgs = append(cmdArgs, "--filter", "label=devcontainer.local_folder="+d.Path)

	return execCmd(cmdArgs, true)
}

// IsRunning return the container status
func (d *Docker) IsRunning() (bool, error) {
	// skip if container does not exist
	if d.Container == "" {
		return false, nil
	}

	running := false
	cmdArgs := []string{"docker", "container", "ls"}
	cmdArgs = append(cmdArgs, "--quiet")
	cmdArgs = append(cmdArgs, "--filter", "label=devcontainer.local_folder="+d.Path)
	out, err := execCmd(cmdArgs, true)
	if d.Container != "" {
		running = lo.Contains(strings.Split(out, "\n"), d.Container)
	}

	return running, err
}

// Build build the image for the given Dockerfile
func (d *Docker) Build() (string, error) {
	// skip if there is no dockerfile to build
	if d.ImageBuild.Dockerfile == "" {
		return "", nil
	}

	cmdArgs := []string{"docker", "image", "build"}
	cmdArgs = append(cmdArgs, "--tag", d.Image)
	cmdArgs = append(cmdArgs, "--file", d.ImageBuild.Dockerfile)
	if d.ImageBuild.Target != "" {
		cmdArgs = append(cmdArgs, "--target", d.ImageBuild.Target)
	}
	if d.ImageBuild.CacheFrom != "" {
		cmdArgs = append(cmdArgs, "--cache-from", d.ImageBuild.CacheFrom)
	}
	for _, arg := range d.ImageBuild.Args {
		cmdArgs = append(cmdArgs, "--build-arg", arg)
	}
	cmdArgs = append(cmdArgs, d.ImageBuild.Context)

	return execCmd(cmdArgs, false)
}

// Create create the container with the given image
func (d *Docker) Create() (string, error) {
	// skip if container already exists
	if d.Container != "" {
		return "", nil
	}

	cmdArgs := []string{"docker", "container", "create"}
	cmdArgs = append(cmdArgs, "--label", "devcontainer.local_folder="+d.Path)
	for _, mount := range d.Mounts {
		cmdArgs = append(cmdArgs, "--mount", mount)
	}
	for _, port := range d.Ports {
		cmdArgs = append(cmdArgs, "--publish", port)
	}
	for _, env := range d.Envs {
		cmdArgs = append(cmdArgs, "--env", env)
	}
	if d.ContainerUser != "" {
		cmdArgs = append(cmdArgs, "--user", d.ContainerUser)
	}
	cmdArgs = append(cmdArgs, d.Image)
	if len(d.Command) > 0 {
		cmdArgs = append(cmdArgs, d.Command...)
	}

	return execCmd(cmdArgs, true)
}

// Start start the given container
func (d *Docker) Start() (string, error) {
	// build and create if not it does not exists
	if d.Container == "" {
		if _, err := d.Build(); err != nil {
			return "", err
		}
		if container, err := d.Create(); err != nil {
			return "", err
		} else {
			d.Container = container
		}
	}
	// skip if already started
	if d.Running {
		return "", nil
	}

	cmdArgs := []string{"docker", "container", "start"}
	cmdArgs = append(cmdArgs, d.Args...)
	cmdArgs = append(cmdArgs, d.Container)

	return execCmd(cmdArgs, true)
}

// Stop stop the given container
func (d *Docker) Stop() (string, error) {
	// skip if already stopped
	if !d.Running {
		return "", nil
	}

	cmdArgs := []string{"docker", "container", "stop"}
	cmdArgs = append(cmdArgs, d.Container)

	return execCmd(cmdArgs, true)
}

// Remove remove the container
func (d *Docker) Remove() (string, error) {
	// skip if container does not exist
	if d.Container == "" {
		return "", nil
	}
	cmdArgs := []string{"docker", "container", "rm"}
	cmdArgs = append(cmdArgs, d.Container)

	return execCmd(cmdArgs, true)
}

// List return the list of containers based on the given path
func (d *Docker) List() (string, error) {
	cmdArgs := []string{"docker", "container", "ls", "--filter", "label=devcontainer.local_folder=" + d.Path}

	return execCmd(cmdArgs, false)
}

// Exec execute the given command into the given container
func (d *Docker) Exec(command []string, withEnv bool, capture bool) (string, error) {
	// start container if not running
	if !d.Running {
		if _, err := d.Start(); err != nil {
			return "", err
		}
	}

	cmdArgs := []string{"docker", "container", "exec"}
	cmdArgs = append(cmdArgs, "--interactive", "--tty")
	cmdArgs = append(cmdArgs, "--workdir", d.WorkDir)
	if d.RemoteUser != "" {
		cmdArgs = append(cmdArgs, "--user", d.RemoteUser)
	}
	// resolve containerEnv variables
	if withEnv {
		d.RemoteEnvs = lo.Map(d.RemoteEnvs, func(v string, _ int) string { return resolveContainerEnv(d, v) })
		for _, env := range d.RemoteEnvs {
			cmdArgs = append(cmdArgs, "--env", env)
		}
	}
	cmdArgs = append(cmdArgs, d.Container)
	cmdArgs = append(cmdArgs, command...)

	return execCmd(cmdArgs, capture)
}

// ResolveEnv resolve environment variable from inside the container
func (d *Docker) ResolveEnv(env string) string {
	cmd := []string{"sh", "-c", "env | awk -F'=' '/" + env + "=/ {print $2}'"}
	resolved, err := d.Exec(cmd, false, true)
	if err != nil {
		panic(err)
	}

	return resolved
}
