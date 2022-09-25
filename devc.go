// devc is a devcontainer managment tool

// Do not want to mess your minimal, clean and fine tuned OS installation with a
// bunch of compilers, linters, debuggers, interpreters, minifiers, unminifiers,
// beautifiers, etc...? Visual Studio Code nailed it with DevContainers. But do
// not want to use VSCode neither? That's where devc comes in, it a simple CLI
// that wrap docker/docker-compose and run the (almost) same commands that
// VSCode runs behind the scenes.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"muzzammil.xyz/jsonc"
)

type DevContainer struct {
	Engine               Engine
	JSON                 DevContainerJSON
	WorkingDirectoryPath string
	WorkingDirectoryName string
}

type Engine interface {
	Init(config *DevContainer) error
	IsRunning() (bool, error)
	Build() (string, error)
	Create() (string, error)
	Remove() (string, error)
	Start() (string, error)
	Stop() (string, error)
	List() (string, error)
	Exec(command []string, withEnv bool, capture bool) (string, error)
	ResolveEnv(env string) string
}

type DevContainerJSON struct {
	Build               DevContainerJSONBuild `json:"build,omitempty"`
	ContainerEnv        map[string]string     `json:"containerEnv,omitempty"`
	ContainerUser       string                `json:"containerUser,omitempty"`
	Customizations      map[string]any        `json:"customizations,omitempty"`
	DockerComposeFile   string                `json:"dockerComposeFile,omitempty"`
	ForwardPorts        []string              `json:"forwardPorts,omitempty"`
	Image               string                `json:"image,omitempty"`
	InitializeCommand   []string              `json:"initializeCommand,omitempty"`
	Mounts              []string              `json:"mounts,omitempty"`
	Name                string                `json:"name,omitempty"`
	OnCreateCommand     []string              `json:"onCreateCommand,omitempty"`
	OverrideCommand     bool                  `json:"overrideCommand,omitempty"`
	PostAttachCommand   []string              `json:"postAttachCommand,omitempty"`
	PostCreateCommand   []string              `json:"postCreateCommand,omitempty"`
	PostStartCommand    []string              `json:"postStartCommand,omitempty"`
	RemoteEnv           map[string]string     `json:"remoteEnv,omitempty"`
	RemoteUser          string                `json:"remoteUser,omitempty"`
	RunArgs             []string              `json:"runArgs,omitempty"`
	RunServices         []string              `json:"runServices,omitempty"`
	Service             string                `json:"service,omitempty"`
	UpdateRemoteUserUID bool                  `json:"updateRemoteUserUID,omitempty"`
	WorkspaceFolder     string                `json:"workspaceFolder,omitempty"`
	WorkspaceMount      string                `json:"workspaceMount,omitempty"`
}

type DevContainerJSONBuild struct {
	Args       map[string]string `json:"args,omitempty"`
	CacheFrom  string            `json:"cacheFrom,omitempty"`
	Context    string            `json:"context,omitempty"`
	Dockerfile string            `json:"dockerfile,omitempty"`
	Target     string            `json:"target,omitempty"`
}

var devc DevContainer

var log zerolog.Logger

var rootVerbose int

var rootCmd = &cobra.Command{
	Use:   "devc",
	Short: "devc is a devcontainer managment tool",
	Long:  ``,
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		setLogLevel()
		if lo.None([]string{"completion", "init", "man", "-h", "--help", "help"}, os.Args) {
			parseConfig()
			setDefaults()
			checkConfig()
			initializeCommand()
			devc.resolveVars()
			setEngine()
		}
	},
}

func init() {
	rootCmd.PersistentFlags().CountVarP(&rootVerbose, "verbose", "v", "enable verbose output")
}

func setLogLevel() {
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

func parseConfig() {
	// return JSONC as JSON
	_, j, err := jsonc.ReadFromFile(".devcontainer/devcontainer.json")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot read devcontainer settings")
	}

	if err := json.Unmarshal(j, &devc.JSON); err != nil {
		log.Fatal().Err(err).Msg("cannot parse json")
	}
	log.Debug().Str("devcontainer", fmt.Sprintf("%+v", devc.JSON)).Send()
}

func setDefaults() {
	// set defaults values for mydevcontainer
	devc.WorkingDirectoryPath, _ = os.Getwd()
	devc.WorkingDirectoryName = filepath.Base(devc.WorkingDirectoryPath)
	devc.JSON.Build.Context = lo.Ternary(devc.JSON.Build.Context != "", devc.JSON.Build.Context, ".")
	devc.JSON.Name = lo.Ternary(devc.JSON.Name != "", devc.JSON.Name, devc.WorkingDirectoryName)
	devc.JSON.OverrideCommand = lo.Ternary(devc.JSON.OverrideCommand, devc.JSON.OverrideCommand, true)
	devc.JSON.UpdateRemoteUserUID = lo.Ternary(devc.JSON.UpdateRemoteUserUID, devc.JSON.UpdateRemoteUserUID, true)
	devc.JSON.WorkspaceFolder = lo.Ternary(devc.JSON.WorkspaceFolder != "", devc.JSON.WorkspaceFolder, "/workspace")
	devc.JSON.WorkspaceMount = lo.Ternary(devc.JSON.WorkspaceMount != "", devc.JSON.WorkspaceMount, "type=bind,source="+devc.WorkingDirectoryPath+",target="+devc.JSON.WorkspaceFolder+",consistency=cached")
}

func checkConfig() {
	// check required and conflicting settings
	if devc.JSON.Image == "" && devc.JSON.Build.Dockerfile == "" && devc.JSON.DockerComposeFile == "" {
		log.Fatal().Msg("one of these settings is required: 'image', 'build.dockerfile', 'dockerComposeFile'")
	}
	if (devc.JSON.Image != "" && devc.JSON.Build.Dockerfile != "") ||
		(devc.JSON.Image != "" && devc.JSON.DockerComposeFile != "") ||
		(devc.JSON.Build.Dockerfile != "" && devc.JSON.DockerComposeFile != "") {
		log.Fatal().Msg("one of these settings conflicts with another one: 'image', 'build.dockerfile', 'dockerComposeFile'")
	}
	if devc.JSON.DockerComposeFile != "" && devc.JSON.Service == "" {
		log.Fatal().Msg("'service' setting is required when using 'dockerComposeFile'")
	}
}

func setEngine() {
	// determine container engine
	switch {
	case devc.JSON.Image != "" || devc.JSON.Build.Dockerfile != "":
		devc.Engine = &Docker{}
	case devc.JSON.DockerComposeFile != "":
		devc.Engine = &DockerCompose{}
	default:
		log.Fatal().Msg("cannot determine devcontainer engine")
	}

	// initialize engine
	if err := devc.Engine.Init(&devc); err != nil {
		log.Fatal().Err(err).Msg("cannot initialize")
	}
	log.Debug().Str("engine", fmt.Sprintf("%+v", devc.Engine)).Send()
}

func initializeCommand() {
	if len(devc.JSON.InitializeCommand) > 0 {
		if _, err := execCmd(devc.JSON.InitializeCommand, false); err != nil {
			log.Fatal().Err(err).Msg("cannot run initializeCommand")
		}
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Send()
	}
}
