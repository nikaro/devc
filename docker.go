package main

import (
	"os"
	"strings"

	"github.com/samber/lo"
)

type Docker struct {
	_ExecCmd        func([]string, bool) (string, error)
	Args            []string
	Capabilities    []string
	Command         []string
	Container       string
	ContainerUser   string
	EnableInit      bool
	EnablePrivilege bool
	Envs            []string
	Image           string
	ImageBuild      DockerImageBuild
	Mounts          []string
	Path            string
	Ports           []string
	RemoteEnvs      []string
	RemoteUser      string
	Running         bool
	SecurityOpts    []string
	WorkDir         string
}

type DockerImageBuild struct {
	Args       []string
	Dockerfile string
	CacheFrom  []string
	Context    string
	Tag        string
	Target     string
}

// Init initialize docker settings
func (d *Docker) Init(c *DevContainer) error {
	if err := os.Chdir(c.ConfigDir); err != nil {
		return err
	}

	d._ExecCmd = lo.Ternary(d._ExecCmd != nil, d._ExecCmd, execCmd)
	d.Args = c.Config.GetStringSlice("runArgs")
	d.Capabilities = c.Config.GetStringSlice("capAdd")
	d.Command = lo.Ternary(
		c.Config.GetBool("overrideCommand"),
		[]string{"/bin/sh", "-c", "while sleep 1000; do :; done"},
		nil,
	)
	d.ContainerUser = c.Config.GetString("containerUser")
	d.EnableInit = c.Config.GetBool("init")
	d.EnablePrivilege = c.Config.GetBool("privileged")
	d.Envs = lo.MapToSlice(
		c.Config.GetStringMapString("containerEnv"),
		func(k string, v string) string { return k + "=" + v },
	)
	d.Image = lo.Ternary(
		c.Config.IsSet("image"),
		c.Config.GetString("image"),
		"vsc-"+c.WorkingDirectoryName+"-"+md5sum(c.WorkingDirectoryPath),
	)
	d.ImageBuild.Args = lo.MapToSlice(
		c.Config.GetStringMapString("build.args"),
		func(k string, v string) string { return k + "=" + v },
	)
	d.ImageBuild.CacheFrom = c.Config.GetStringSlice("build.cacheFrom")
	d.ImageBuild.Context = c.Config.GetString("build.context")
	d.ImageBuild.Dockerfile = c.Config.GetString("build.dockerfile")
	d.ImageBuild.Target = c.Config.GetString("build.target")
	d.Mounts = c.Config.GetStringSlice("mounts")
	d.Mounts = append(d.Mounts, c.Config.GetString("workspaceMount"))
	d.Path = c.WorkingDirectoryPath
	d.Ports = c.Config.GetStringSlice("forwardPorts")
	d.RemoteEnvs = lo.MapToSlice(
		c.Config.GetStringMapString("remoteEnv"),
		func(k string, v string) string { return k + "=" + v },
	)
	d.RemoteUser = c.Config.GetString("remoteUser")
	d.WorkDir = c.Config.GetString("workspaceFolder")

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

	return d._ExecCmd(cmdArgs, true)
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
	out, err := d._ExecCmd(cmdArgs, true)
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
	for _, cache := range d.ImageBuild.CacheFrom {
		cmdArgs = append(cmdArgs, "--cache-from", cache)
	}
	for _, arg := range d.ImageBuild.Args {
		cmdArgs = append(cmdArgs, "--build-arg", arg)
	}
	cmdArgs = append(cmdArgs, d.ImageBuild.Context)

	return d._ExecCmd(cmdArgs, false)
}

// Create create the container with the given image
func (d *Docker) Create() (string, error) {
	// skip if container already exists
	if d.Container != "" {
		return "", nil
	}

	cmdArgs := []string{"docker", "container", "create"}
	cmdArgs = append(cmdArgs, "--label", "devcontainer.local_folder="+d.Path)
	if d.EnableInit {
		cmdArgs = append(cmdArgs, "--init")
	}
	for _, cap := range d.Capabilities {
		cmdArgs = append(cmdArgs, "--cap-add", cap)
	}
	for _, sec := range d.SecurityOpts {
		cmdArgs = append(cmdArgs, "--security-opt", sec)
	}
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

	return d._ExecCmd(cmdArgs, true)
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

	return d._ExecCmd(cmdArgs, true)
}

// Stop stop the given container
func (d *Docker) Stop() (string, error) {
	// skip if already stopped
	if !d.Running {
		return "", nil
	}

	cmdArgs := []string{"docker", "container", "stop"}
	cmdArgs = append(cmdArgs, d.Container)

	return d._ExecCmd(cmdArgs, true)
}

// Remove remove the container
func (d *Docker) Remove() (string, error) {
	// skip if container does not exist
	if d.Container == "" {
		return "", nil
	}
	cmdArgs := []string{"docker", "container", "rm"}
	cmdArgs = append(cmdArgs, d.Container)

	return d._ExecCmd(cmdArgs, true)
}

// List return the list of containers based on the given path
func (d *Docker) List() (string, error) {
	cmdArgs := []string{"docker", "container", "ls", "--filter", "label=devcontainer.local_folder=" + d.Path}

	return d._ExecCmd(cmdArgs, false)
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

	return d._ExecCmd(cmdArgs, capture)
}

// ResolveEnv resolve environment variable from inside the container
func (d *Docker) ResolveEnv(env string) string {
	cmd := []string{"echo", "$" + env}
	resolved, err := d.Exec(cmd, false, true)
	if err != nil {
		panic(err)
	}

	return resolved
}
