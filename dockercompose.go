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
	File        string
	ProjectName string
	Running     bool
	RunServices []string
	Service     string
	User        string
	WorkDir     string
}

func (d *DockerCompose) cmd(args ...string) []string {
	cmd := []string{"docker", "compose", "--file", d.File, "--project-name", d.ProjectName}
	cmd = append(cmd, args...)

	return cmd
}

// Init initialize compose settings
func (d *DockerCompose) Init(config *DevContainer) error {
	d._ExecCmd = lo.Ternary(d._ExecCmd != nil, d._ExecCmd, execCmd)
	d.Envs = lo.MapToSlice(config.JSON.RemoteEnv, func(k string, v string) string { return k + "=" + v })
	d.File = filepath.Join(config.ConfigDir, config.JSON.DockerComposeFile)
	d.ProjectName = config.JSON.Name + "_devcontainer"
	d.RunServices = config.JSON.RunServices
	d.Service = config.JSON.Service
	d.User = config.JSON.RemoteUser
	d.WorkDir = config.JSON.WorkspaceFolder

	// check if already started
	if running, err := d.IsRunning(); err != nil {
		return err
	} else if running {
		d.Running = true
	}

	return nil
}

// IsRunning return the container status
func (d *DockerCompose) IsRunning() (bool, error) {
	cmdArgs := d.cmd("ps")
	cmdArgs = append(cmdArgs, "--quiet")
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
func (d *DockerCompose) Exec(command []string, withEnv bool, capture bool) (string, error) {
	// start service if not running
	if !d.Running {
		if _, err := d.Start(); err != nil {
			return "", err
		}
	}

	cmdArgs := d.cmd("exec")
	cmdArgs = append(cmdArgs, "--workdir", d.WorkDir)
	if d.User != "" {
		cmdArgs = append(cmdArgs, "--user", d.User)
	}
	// resolve containerEnv variables
	if withEnv {
		d.Envs = lo.Map(d.Envs, func(v string, _ int) string { return resolveContainerEnv(d, v) })
		for _, env := range d.Envs {
			cmdArgs = append(cmdArgs, "--env", env)
		}
	}
	cmdArgs = append(cmdArgs, d.Service)
	cmdArgs = append(cmdArgs, command...)

	return d._ExecCmd(cmdArgs, capture)
}

// ResolveEnv resolve environment variable from inside the container
func (d *DockerCompose) ResolveEnv(env string) string {
	cmd := []string{"sh", "-c", "env | awk -F'=' '/" + env + "=/ {print $2}'"}
	resolved, err := d.Exec(cmd, false, true)
	if err != nil {
		panic(err)
	}

	return resolved
}
