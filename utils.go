package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"muzzammil.xyz/jsonc"
)

// runs the given command while attaching stdin, stdout and stderr
func execCmd(command []string, capture bool) (string, error) {
	var stdout []byte
	var err error

	cwd, _ := os.Getwd()
	cmd := exec.Command(command[0], command[1:]...)
	log.Info().Str("workdir", cwd).Str("command", cmd.String()).Send()
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if capture {
		stdout, err = cmd.Output()
	} else {
		cmd.Stdout = os.Stdout
		err = cmd.Run()
	}

	return strings.TrimSpace(string(stdout)), err
}

// return the md5 hash for a string
func md5sum(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	hash := hex.EncodeToString(hasher.Sum(nil))

	return hash
}

type resolve func(string) string

type matchStore struct {
	Old string
	New string
}

// replace all matches of the regexp in the string with the resolved value
func replace(s string, r *regexp.Regexp, f resolve) string {
	if r.MatchString(s) {
		matches := r.FindAllStringSubmatchIndex(s, -1)
		matchesStore := []matchStore{}
		for _, match := range matches {
			matchesStore = append(matchesStore, matchStore{s[match[0]:match[1]], f(s[match[2]:match[3]])})
		}
		for _, store := range matchesStore {
			s = strings.Replace(s, store.Old, store.New, 1)
		}
	}

	return s
}

// resolve ${localEnv:VARIABLE_NAME}
func resolveLocalEnv(s string) string {
	regexpLocalEnv := regexp.MustCompile(`\${localEnv:([[:word:]]+)}`)

	return replace(s, regexpLocalEnv, os.Getenv)
}

// resolve ${containerEnv:VARIABLE_NAME}
func resolveContainerEnv(e Engine, s string) string {
	regexpLocalEnv := regexp.MustCompile(`\${containerEnv:([[:word:]]+)}`)

	return replace(s, regexpLocalEnv, e.ResolveEnv)
}

// PRERUN UTILS

func (d *DevContainer) SetLogLevel() {
	log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

	switch rootVerbose {
	case 1:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case 2:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}
}

func (d *DevContainer) ParseConfig() {
	// return JSONC as JSON
	_, j, err := jsonc.ReadFromFile(".devcontainer/devcontainer.json")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot read devcontainer settings")
	}

	if err := json.Unmarshal(j, &d.JSON); err != nil {
		log.Fatal().Err(err).Msg("cannot parse json")
	}
	log.Debug().Str("devcontainer", fmt.Sprintf("%+v", d.JSON)).Send()
}

func (d *DevContainer) SetDefaults() {
	// set defaults values for mydevcontainer
	d.WorkingDirectoryPath, _ = os.Getwd()
	d.WorkingDirectoryName = filepath.Base(d.WorkingDirectoryPath)
	d.JSON.Build.Context = lo.Ternary(d.JSON.Build.Context != "", d.JSON.Build.Context, ".")
	d.JSON.Name = lo.Ternary(d.JSON.Name != "", d.JSON.Name, d.WorkingDirectoryName)
	d.JSON.OverrideCommand = lo.Ternary(d.JSON.OverrideCommand, d.JSON.OverrideCommand, true)
	d.JSON.UpdateRemoteUserUID = lo.Ternary(d.JSON.UpdateRemoteUserUID, d.JSON.UpdateRemoteUserUID, true)
	d.JSON.WorkspaceFolder = lo.Ternary(d.JSON.WorkspaceFolder != "", d.JSON.WorkspaceFolder, "/workspace")
	d.JSON.WorkspaceMount = lo.Ternary(d.JSON.WorkspaceMount != "", d.JSON.WorkspaceMount, "type=bind,source="+d.WorkingDirectoryPath+",target="+d.JSON.WorkspaceFolder+",consistency=cached")
}

func (d *DevContainer) CheckConfig() {
	// check required and conflicting settings
	if d.JSON.Image == "" && d.JSON.Build.Dockerfile == "" && d.JSON.DockerComposeFile == "" {
		log.Fatal().Msg("one of these settings is required: 'image', 'build.dockerfile', 'dockerComposeFile'")
	}
	if (d.JSON.Image != "" && d.JSON.Build.Dockerfile != "") ||
		(d.JSON.Image != "" && d.JSON.DockerComposeFile != "") ||
		(d.JSON.Build.Dockerfile != "" && d.JSON.DockerComposeFile != "") {
		log.Fatal().Msg("one of these settings conflicts with another one: 'image', 'build.dockerfile', 'dockerComposeFile'")
	}
	if d.JSON.DockerComposeFile != "" && d.JSON.Service == "" {
		log.Fatal().Msg("'service' setting is required when using 'dockerComposeFile'")
	}
}

func (d *DevContainer) SetEngine() {
	// determine container engine
	switch {
	case d.JSON.Image != "" || d.JSON.Build.Dockerfile != "":
		d.Engine = &Docker{}
	case d.JSON.DockerComposeFile != "":
		d.Engine = &DockerCompose{}
	default:
		log.Fatal().Msg("cannot determine devcontainer engine")
	}

	// initialize engine
	if err := d.Engine.Init(d); err != nil {
		log.Fatal().Err(err).Msg("cannot initialize")
	}
	log.Debug().Str("engine", fmt.Sprintf("%+v", d.Engine)).Send()
}

// https://containers.dev/implementors/json_reference/#variables-in-devcontainerjson
func (d *DevContainer) ResolveVars() {
	// resolve ${localEnv:VARIABLE_NAME}
	d.JSON.Build.Args = lo.MapValues(d.JSON.Build.Args, func(v string, _ string) string { return resolveLocalEnv(v) })
	d.JSON.Build.CacheFrom = resolveLocalEnv(d.JSON.Build.CacheFrom)
	d.JSON.Build.Context = resolveLocalEnv(d.JSON.Build.Context)
	d.JSON.Build.Dockerfile = resolveLocalEnv(d.JSON.Build.Dockerfile)
	d.JSON.Build.Target = resolveLocalEnv(d.JSON.Build.Target)
	d.JSON.ContainerEnv = lo.MapValues(d.JSON.ContainerEnv, func(v string, _ string) string { return resolveLocalEnv(v) })
	d.JSON.ContainerUser = resolveLocalEnv(d.JSON.ContainerUser)
	d.JSON.DockerComposeFile = resolveLocalEnv(d.JSON.DockerComposeFile)
	d.JSON.ForwardPorts = lo.Map(d.JSON.ForwardPorts, func(v string, _ int) string { return resolveLocalEnv(v) })
	d.JSON.Image = resolveLocalEnv(d.JSON.Image)
	d.JSON.Mounts = lo.Map(d.JSON.Mounts, func(v string, _ int) string { return resolveLocalEnv(v) })
	d.JSON.Name = resolveLocalEnv(d.JSON.Name)
	d.JSON.RemoteEnv = lo.MapValues(d.JSON.RemoteEnv, func(v string, _ string) string { return resolveLocalEnv(v) })
	d.JSON.RemoteUser = resolveLocalEnv(d.JSON.RemoteUser)
	d.JSON.RunArgs = lo.Map(d.JSON.RunArgs, func(v string, _ int) string { return resolveLocalEnv(v) })
	d.JSON.RunServices = lo.Map(d.JSON.RunServices, func(v string, _ int) string { return resolveLocalEnv(v) })
	d.JSON.Service = resolveLocalEnv(d.JSON.Service)
	d.JSON.WorkspaceFolder = resolveLocalEnv(d.JSON.WorkspaceFolder)
	d.JSON.WorkspaceMount = resolveLocalEnv(d.JSON.WorkspaceMount)

	// resolve ${localWorkspaceFolder}
	d.JSON.Build.Args = lo.MapValues(d.JSON.Build.Args, func(v string, _ string) string {
		return strings.ReplaceAll(v, "${localWorkspaceFolder}", d.WorkingDirectoryPath)
	})
	d.JSON.Build.CacheFrom = strings.ReplaceAll(d.JSON.Build.CacheFrom, "${localWorkspaceFolder}", d.WorkingDirectoryPath)
	d.JSON.Build.Context = strings.ReplaceAll(d.JSON.Build.Context, "${localWorkspaceFolder}", d.WorkingDirectoryPath)
	d.JSON.Build.Dockerfile = strings.ReplaceAll(d.JSON.Build.Dockerfile, "${localWorkspaceFolder}", d.WorkingDirectoryPath)
	d.JSON.Build.Target = strings.ReplaceAll(d.JSON.Build.Target, "${localWorkspaceFolder}", d.WorkingDirectoryPath)
	d.JSON.ContainerEnv = lo.MapValues(d.JSON.ContainerEnv, func(v string, _ string) string {
		return strings.ReplaceAll(v, "${localWorkspaceFolder}", d.WorkingDirectoryPath)
	})
	d.JSON.ContainerUser = strings.ReplaceAll(d.JSON.ContainerUser, "${localWorkspaceFolder}", d.WorkingDirectoryPath)
	d.JSON.DockerComposeFile = strings.ReplaceAll(d.JSON.DockerComposeFile, "${localWorkspaceFolder}", d.WorkingDirectoryPath)
	d.JSON.ForwardPorts = lo.Map(d.JSON.ForwardPorts, func(v string, _ int) string {
		return strings.ReplaceAll(v, "${localWorkspaceFolder}", d.WorkingDirectoryPath)
	})
	d.JSON.Image = strings.ReplaceAll(d.JSON.Image, "${localWorkspaceFolder}", d.WorkingDirectoryPath)
	d.JSON.Mounts = lo.Map(d.JSON.Mounts, func(v string, _ int) string {
		return strings.ReplaceAll(v, "${localWorkspaceFolder}", d.WorkingDirectoryPath)
	})
	d.JSON.Name = strings.ReplaceAll(d.JSON.Name, "${localWorkspaceFolder}", d.WorkingDirectoryPath)
	d.JSON.RemoteEnv = lo.MapValues(d.JSON.RemoteEnv, func(v string, _ string) string {
		return strings.ReplaceAll(v, "${localWorkspaceFolder}", d.WorkingDirectoryPath)
	})
	d.JSON.RemoteUser = strings.ReplaceAll(d.JSON.RemoteUser, "${localWorkspaceFolder}", d.WorkingDirectoryPath)
	d.JSON.RunArgs = lo.Map(d.JSON.RunArgs, func(v string, _ int) string {
		return strings.ReplaceAll(v, "${localWorkspaceFolder}", d.WorkingDirectoryPath)
	})
	d.JSON.RunServices = lo.Map(d.JSON.RunServices, func(v string, _ int) string {
		return strings.ReplaceAll(v, "${localWorkspaceFolder}", d.WorkingDirectoryPath)
	})
	d.JSON.Service = strings.ReplaceAll(d.JSON.Service, "${localWorkspaceFolder}", d.WorkingDirectoryPath)
	d.JSON.WorkspaceFolder = strings.ReplaceAll(d.JSON.WorkspaceFolder, "${localWorkspaceFolder}", d.WorkingDirectoryPath)
	d.JSON.WorkspaceMount = strings.ReplaceAll(d.JSON.WorkspaceMount, "${localWorkspaceFolder}", d.WorkingDirectoryPath)

	// resolve ${containerWorkspaceFolder}
	d.JSON.Build.Args = lo.MapValues(d.JSON.Build.Args, func(v string, _ string) string {
		return strings.ReplaceAll(v, "${containerWorkspaceFolder}", d.WorkingDirectoryPath)
	})
	d.JSON.Build.CacheFrom = strings.ReplaceAll(d.JSON.Build.CacheFrom, "${containerWorkspaceFolder}", d.JSON.WorkspaceFolder)
	d.JSON.Build.Context = strings.ReplaceAll(d.JSON.Build.Context, "${containerWorkspaceFolder}", d.JSON.WorkspaceFolder)
	d.JSON.Build.Dockerfile = strings.ReplaceAll(d.JSON.Build.Dockerfile, "${containerWorkspaceFolder}", d.JSON.WorkspaceFolder)
	d.JSON.Build.Target = strings.ReplaceAll(d.JSON.Build.Target, "${containerWorkspaceFolder}", d.JSON.WorkspaceFolder)
	d.JSON.ContainerEnv = lo.MapValues(d.JSON.ContainerEnv, func(v string, _ string) string {
		return strings.ReplaceAll(v, "${containerWorkspaceFolder}", d.JSON.WorkspaceFolder)
	})
	d.JSON.ContainerUser = strings.ReplaceAll(d.JSON.ContainerUser, "${containerWorkspaceFolder}", d.JSON.WorkspaceFolder)
	d.JSON.DockerComposeFile = strings.ReplaceAll(d.JSON.DockerComposeFile, "${containerWorkspaceFolder}", d.JSON.WorkspaceFolder)
	d.JSON.ForwardPorts = lo.Map(d.JSON.ForwardPorts, func(v string, _ int) string {
		return strings.ReplaceAll(v, "${containerWorkspaceFolder}", d.JSON.WorkspaceFolder)
	})
	d.JSON.Image = strings.ReplaceAll(d.JSON.Image, "${containerWorkspaceFolder}", d.JSON.WorkspaceFolder)
	d.JSON.Mounts = lo.Map(d.JSON.Mounts, func(v string, _ int) string {
		return strings.ReplaceAll(v, "${containerWorkspaceFolder}", d.JSON.WorkspaceFolder)
	})
	d.JSON.Name = strings.ReplaceAll(d.JSON.Name, "${containerWorkspaceFolder}", d.JSON.WorkspaceFolder)
	d.JSON.RemoteEnv = lo.MapValues(d.JSON.RemoteEnv, func(v string, _ string) string {
		return strings.ReplaceAll(v, "${containerWorkspaceFolder}", d.JSON.WorkspaceFolder)
	})
	d.JSON.RemoteUser = strings.ReplaceAll(d.JSON.RemoteUser, "${containerWorkspaceFolder}", d.JSON.WorkspaceFolder)
	d.JSON.RunArgs = lo.Map(d.JSON.RunArgs, func(v string, _ int) string {
		return strings.ReplaceAll(v, "${containerWorkspaceFolder}", d.JSON.WorkspaceFolder)
	})
	d.JSON.RunServices = lo.Map(d.JSON.RunServices, func(v string, _ int) string {
		return strings.ReplaceAll(v, "${containerWorkspaceFolder}", d.JSON.WorkspaceFolder)
	})
	d.JSON.Service = strings.ReplaceAll(d.JSON.Service, "${containerWorkspaceFolder}", d.JSON.WorkspaceFolder)
	d.JSON.WorkspaceMount = strings.ReplaceAll(d.JSON.WorkspaceMount, "${containerWorkspaceFolder}", d.JSON.WorkspaceFolder)

	// resolve ${localWorkspaceFolderBasename}
	d.JSON.Build.Args = lo.MapValues(d.JSON.Build.Args, func(v string, _ string) string {
		return strings.ReplaceAll(v, "${localWorkspaceFolderBasename}", d.WorkingDirectoryName)
	})
	d.JSON.Build.CacheFrom = strings.ReplaceAll(d.JSON.Build.CacheFrom, "${localWorkspaceFolderBasename}", d.WorkingDirectoryName)
	d.JSON.Build.Context = strings.ReplaceAll(d.JSON.Build.Context, "${localWorkspaceFolderBasename}", d.WorkingDirectoryName)
	d.JSON.Build.Dockerfile = strings.ReplaceAll(d.JSON.Build.Dockerfile, "${localWorkspaceFolderBasename}", d.WorkingDirectoryName)
	d.JSON.Build.Target = strings.ReplaceAll(d.JSON.Build.Target, "${localWorkspaceFolderBasename}", d.WorkingDirectoryName)
	d.JSON.ContainerEnv = lo.MapValues(d.JSON.ContainerEnv, func(v string, _ string) string {
		return strings.ReplaceAll(v, "${localWorkspaceFolderBasename}", d.WorkingDirectoryName)
	})
	d.JSON.ContainerUser = strings.ReplaceAll(d.JSON.ContainerUser, "${localWorkspaceFolderBasename}", d.WorkingDirectoryName)
	d.JSON.DockerComposeFile = strings.ReplaceAll(d.JSON.DockerComposeFile, "${localWorkspaceFolderBasename}", d.WorkingDirectoryName)
	d.JSON.ForwardPorts = lo.Map(d.JSON.ForwardPorts, func(v string, _ int) string {
		return strings.ReplaceAll(v, "${localWorkspaceFolderBasename}", d.WorkingDirectoryName)
	})
	d.JSON.Image = strings.ReplaceAll(d.JSON.Image, "${localWorkspaceFolderBasename}", d.WorkingDirectoryName)
	d.JSON.Mounts = lo.Map(d.JSON.Mounts, func(v string, _ int) string {
		return strings.ReplaceAll(v, "${localWorkspaceFolderBasename}", d.WorkingDirectoryName)
	})
	d.JSON.Name = strings.ReplaceAll(d.JSON.Name, "${localWorkspaceFolderBasename}", d.WorkingDirectoryName)
	d.JSON.RemoteEnv = lo.MapValues(d.JSON.RemoteEnv, func(v string, _ string) string {
		return strings.ReplaceAll(v, "${localWorkspaceFolderBasename}", d.WorkingDirectoryName)
	})
	d.JSON.RemoteUser = strings.ReplaceAll(d.JSON.RemoteUser, "${localWorkspaceFolderBasename}", d.WorkingDirectoryName)
	d.JSON.RunArgs = lo.Map(d.JSON.RunArgs, func(v string, _ int) string {
		return strings.ReplaceAll(v, "${localWorkspaceFolderBasename}", d.WorkingDirectoryName)
	})
	d.JSON.RunServices = lo.Map(d.JSON.RunServices, func(v string, _ int) string {
		return strings.ReplaceAll(v, "${localWorkspaceFolderBasename}", d.WorkingDirectoryName)
	})
	d.JSON.Service = strings.ReplaceAll(d.JSON.Service, "${localWorkspaceFolderBasename}", d.WorkingDirectoryName)
	d.JSON.WorkspaceFolder = strings.ReplaceAll(d.JSON.WorkspaceFolder, "${localWorkspaceFolderBasename}", d.WorkingDirectoryName)
	d.JSON.WorkspaceMount = strings.ReplaceAll(d.JSON.WorkspaceMount, "${localWorkspaceFolderBasename}", d.WorkingDirectoryName)

	// resolve ${containerWorkspaceFolderBasename}
	workspaceFolderName := filepath.Base(d.JSON.WorkspaceFolder)
	d.JSON.Build.Args = lo.MapValues(d.JSON.Build.Args, func(v string, _ string) string {
		return strings.ReplaceAll(v, "${containerWorkspaceFolderBasename}", workspaceFolderName)
	})
	d.JSON.Build.CacheFrom = strings.ReplaceAll(d.JSON.Build.CacheFrom, "${containerWorkspaceFolderBasename}", workspaceFolderName)
	d.JSON.Build.Context = strings.ReplaceAll(d.JSON.Build.Context, "${containerWorkspaceFolderBasename}", workspaceFolderName)
	d.JSON.Build.Dockerfile = strings.ReplaceAll(d.JSON.Build.Dockerfile, "${containerWorkspaceFolderBasename}", workspaceFolderName)
	d.JSON.Build.Target = strings.ReplaceAll(d.JSON.Build.Target, "${containerWorkspaceFolderBasename}", workspaceFolderName)
	d.JSON.ContainerEnv = lo.MapValues(d.JSON.ContainerEnv, func(v string, _ string) string {
		return strings.ReplaceAll(v, "${containerWorkspaceFolderBasename}", workspaceFolderName)
	})
	d.JSON.ContainerUser = strings.ReplaceAll(d.JSON.ContainerUser, "${containerWorkspaceFolderBasename}", workspaceFolderName)
	d.JSON.DockerComposeFile = strings.ReplaceAll(d.JSON.DockerComposeFile, "${containerWorkspaceFolderBasename}", workspaceFolderName)
	d.JSON.ForwardPorts = lo.Map(d.JSON.ForwardPorts, func(v string, _ int) string {
		return strings.ReplaceAll(v, "${containerWorkspaceFolderBasename}", workspaceFolderName)
	})
	d.JSON.Image = strings.ReplaceAll(d.JSON.Image, "${containerWorkspaceFolderBasename}", workspaceFolderName)
	d.JSON.Mounts = lo.Map(d.JSON.Mounts, func(v string, _ int) string {
		return strings.ReplaceAll(v, "${containerWorkspaceFolderBasename}", workspaceFolderName)
	})
	d.JSON.Name = strings.ReplaceAll(d.JSON.Name, "${containerWorkspaceFolderBasename}", workspaceFolderName)
	d.JSON.RemoteEnv = lo.MapValues(d.JSON.RemoteEnv, func(v string, _ string) string {
		return strings.ReplaceAll(v, "${containerWorkspaceFolderBasename}", workspaceFolderName)
	})
	d.JSON.RemoteUser = strings.ReplaceAll(d.JSON.RemoteUser, "${containerWorkspaceFolderBasename}", workspaceFolderName)
	d.JSON.RunArgs = lo.Map(d.JSON.RunArgs, func(v string, _ int) string {
		return strings.ReplaceAll(v, "${containerWorkspaceFolderBasename}", workspaceFolderName)
	})
	d.JSON.RunServices = lo.Map(d.JSON.RunServices, func(v string, _ int) string {
		return strings.ReplaceAll(v, "${containerWorkspaceFolderBasename}", workspaceFolderName)
	})
	d.JSON.Service = strings.ReplaceAll(d.JSON.Service, "${containerWorkspaceFolderBasename}", workspaceFolderName)
	d.JSON.WorkspaceMount = strings.ReplaceAll(d.JSON.WorkspaceMount, "${containerWorkspaceFolderBasename}", workspaceFolderName)
}
