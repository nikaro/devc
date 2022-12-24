package main

import (
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

// DockerCompose type
type DockerCompose struct {
	_ExecCmd    func([]string, bool) (string, error)
	Command     []string
	Containers  []string
	Envs        []string
	Files       []string
	ProjectName string
	Running     bool
	RunServices []string
	Service     string
	User        string
	WorkDir     string
}

func (d *DockerCompose) cmd(args ...string) []string {
	cmd := []string{"docker", "compose", "--project-name", d.ProjectName}
	for _, file := range d.Files {
		cmd = append(cmd, "--file", file)
	}
	cmd = append(cmd, args...)

	return cmd
}

// Init initialize compose settings
func (d *DockerCompose) Init(c *DevContainer) error {
	d._ExecCmd = lo.Ternary(d._ExecCmd != nil, d._ExecCmd, execCmd)
	d.Envs = lo.MapToSlice(
		c.Config.GetStringMapString("remoteEnv"),
		func(k string, v string) string { return k + "=" + v },
	)
	d.Files = lo.Map(
		c.Config.GetStringSlice("dockerComposeFile"),
		func(v string, _ int) string { return filepath.Join(c.ConfigDir, v) },
	)
	d.ProjectName = c.Config.GetString("name") + "_devcontainer"
	d.RunServices = c.Config.GetStringSlice("runServices")
	d.Service = c.Config.GetString("service")
	d.User = c.Config.GetString("remoteUser")
	d.WorkDir = c.Config.GetString("workspaceFolder")

	// check if already started
	if running, err := d.IsRunning(); err != nil {
		return err
	} else if running {
		d.Running = true
	}

	return nil
}

// IsBuilt return the image build status
func (d *DockerCompose) IsBuilt() (bool, error) {
	cmdArgs := d.cmd("images")
	cmdArgs = append(cmdArgs, "--quiet")
	out, err := d._ExecCmd(cmdArgs, true)
	images := lo.Filter(strings.Split(out, "\n"), func(x string, _ int) bool { return x != "" })
	built := len(images) > 0

	return built, err
}

// IsCreated return the container creation status
func (d *DockerCompose) IsCreated() (bool, error) {
	cmdArgs := d.cmd("ps")
	cmdArgs = append(cmdArgs, "--quiet")
	out, err := d._ExecCmd(cmdArgs, true)
	containers := lo.Filter(strings.Split(out, "\n"), func(x string, _ int) bool { return x != "" })
	created := len(containers) > 0

	return created, err
}

// IsRunning return the container running status
func (d *DockerCompose) IsRunning() (bool, error) {
	cmdArgs := d.cmd("ps")
	cmdArgs = append(cmdArgs, "--quiet")
	cmdArgs = append(cmdArgs, "--status", "running")
	out, err := d._ExecCmd(cmdArgs, true)
	containers := lo.Filter(strings.Split(out, "\n"), func(x string, _ int) bool { return x != "" })
	running := len(containers) > 0

	return running, err
}

// Build build the container with the given image
func (d *DockerCompose) Build() (string, error) {
	cmdArgs := d.cmd("build")
	if len(d.RunServices) > 0 {
		cmdArgs = append(cmdArgs, d.RunServices...)
	}

	return d._ExecCmd(cmdArgs, false)
}

// Create create the container with the given image
func (d *DockerCompose) Create() (string, error) {
	cmdArgs := d.cmd("create")
	if len(d.RunServices) > 0 {
		cmdArgs = append(cmdArgs, d.RunServices...)
	}

	return d._ExecCmd(cmdArgs, false)
}

// Start start the given container
func (d *DockerCompose) Start() (string, error) {
	cmdArgs := d.cmd("up", "--detach")
	if len(d.RunServices) > 0 {
		cmdArgs = append(cmdArgs, d.RunServices...)
	}

	return d._ExecCmd(cmdArgs, false)
}

// Stop stop the given container
func (d *DockerCompose) Stop() (string, error) {
	cmdArgs := d.cmd("stop")
	if len(d.RunServices) > 0 {
		cmdArgs = append(cmdArgs, d.RunServices...)
	}

	return d._ExecCmd(cmdArgs, false)
}

// Remove remove the given container
func (d *DockerCompose) Remove() (string, error) {
	cmdArgs := d.cmd("down", "--volumes")

	return d._ExecCmd(cmdArgs, false)
}

// List return the list of containers based on the given path
func (d *DockerCompose) List() (string, error) {
	cmdArgs := d.cmd("ls")

	return d._ExecCmd(cmdArgs, false)
}

// Exec execute the given command into the given container
func (d *DockerCompose) Run(command []string) (string, error) {
	cmdArgs := d.cmd("run")
	cmdArgs = append(cmdArgs, "--workdir", d.WorkDir)
	if d.User != "" {
		cmdArgs = append(cmdArgs, "--user", d.User)
	}
	cmdArgs = append(cmdArgs, d.Service)
	cmdArgs = append(cmdArgs, command...)

	return d._ExecCmd(cmdArgs, true)
}

// Exec execute the given command into the given container
func (d *DockerCompose) Exec(command []string) (string, error) {
	cmdArgs := d.cmd("exec")
	cmdArgs = append(cmdArgs, "--workdir", d.WorkDir)
	if d.User != "" {
		cmdArgs = append(cmdArgs, "--user", d.User)
	}
	// resolve containerEnv variables
	d.Envs = lo.Map(d.Envs, func(v string, _ int) string { return resolveContainerEnv(d, v) })
	for _, env := range d.Envs {
		cmdArgs = append(cmdArgs, "--env", env)
	}
	cmdArgs = append(cmdArgs, d.Service)
	cmdArgs = append(cmdArgs, command...)

	return d._ExecCmd(cmdArgs, false)
}

// ResolveEnv resolve environment variable from inside the container
func (d *DockerCompose) ResolveEnv(env string) string {
	cmd := []string{"echo", "$" + env}
	resolved, err := d.Run(cmd)
	if err != nil {
		panic(err)
	}

	return resolved
}
