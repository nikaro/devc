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
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// global log var
var log zerolog.Logger

// devcontainer.json structure
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

// interface that each engine (docker, compose, podman, k8s, etc.) must implement
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

// devcontainer meta-structure
type DevContainer struct {
	Engine               Engine
	JSON                 DevContainerJSON
	WorkingDirectoryPath string
	WorkingDirectoryName string
}

var version string

// global devcontainer var
var devc DevContainer

// DEVC COMMANDS

// cli args
var rootVerbose int
var initJSON []byte
var manOutDir string
var shellBin string
var stopRemove bool

func init() {
	// devc command
	rootCmd.PersistentFlags().CountVarP(&rootVerbose, "verbose", "v", "enable verbose output")
	// build sub-command
	rootCmd.AddCommand(buildCmd)
	// init sub-command
	rootCmd.AddCommand(initCmd)
	// list sub-command
	rootCmd.AddCommand(listCmd)
	// man sub-command
	manCmd.PersistentFlags().StringVarP(&manOutDir, "output", "o", "man", "output directory")
	rootCmd.AddCommand(manCmd)
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	// shell sub-command
	shellCmd.PersistentFlags().StringVarP(&shellBin, "shell", "s", "sh", "override shell")
	rootCmd.AddCommand(shellCmd)
	// start sub-command
	rootCmd.AddCommand(startCmd)
	// stop sub-command
	stopCmd.PersistentFlags().BoolVarP(&stopRemove, "remove", "r", false, "remove containers and networks")
	rootCmd.AddCommand(stopCmd)
}

var rootCmd = &cobra.Command{
	Use:              "devc",
	Version:          version,
	Short:            "devc is a devcontainer managment tool",
	Long:             ``,
	PersistentPreRun: devc.PreRun,
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build devcontainer",
	Run:   devc.Build,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize devcontainer configuration",
	Run:   devc.Init,
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "ps"},
	Short:   "List devcontainers",
	Run:     devc.List,
}

var manCmd = &cobra.Command{
	Use:    "man",
	Short:  "Generate manpage",
	Hidden: true,
	Run:    devc.Man,
}

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Execute a shell inside devcontainer",
	Run:   devc.Shell,
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start devcontainer",
	Run:   devc.Start,
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop devcontainer",
	Run:   devc.Stop,
}

func (d *DevContainer) PreRun(_ *cobra.Command, _ []string) {
	d.SetLogLevel()
	if lo.None([]string{"completion", "init", "man", "-h", "--help", "help"}, os.Args) {
		d.ParseConfig()
		d.SetDefaults()
		d.CheckConfig()
		d.InitializeCommand()
		d.ResolveVars()
		d.SetEngine()
	}
}

func (d *DevContainer) Build(_ *cobra.Command, _ []string) {
	if _, err := d.Engine.Build(); err != nil {
		log.Fatal().Err(err).Msg("cannot build")
	}
	if _, err := d.Engine.Create(); err != nil {
		log.Fatal().Err(err).Msg("cannot create")
	}
	d.OnCreateCommand()
	// TODO: figure out what updateContentCommand does and add it
	// go d.UpdateContent()
	d.PostCreateCommand()
}

func (d *DevContainer) Init(_ *cobra.Command, _ []string) {
	fs := afero.NewOsFs()
	if exists, _ := afero.DirExists(fs, ".devcontainer"); !exists {
		if err := fs.Mkdir(".devcontainer", 0755); err != nil {
			log.Fatal().Err(err).Msg("cannot create .devcontainer directory")
		}
		log.Info().Msg(".devcontainer directory created")
	}
	if exists, _ := afero.Exists(fs, ".devcontainer/devcontainer.json"); !exists {
		d.JSON.Image = "alpine:latest"
		initJSON, _ = json.MarshalIndent(d.JSON, "", "  ")
		if err := afero.WriteFile(fs, ".devcontainer/devcontainer.json", initJSON, 0644); err != nil {
			log.Fatal().Err(err).Msg("cannot write devcontainer.json file")
		}
		log.Info().Msg("devcontainer.json file created")
	}
}

func (d *DevContainer) List(_ *cobra.Command, _ []string) {
	if _, err := d.Engine.List(); err != nil {
		log.Fatal().Err(err).Msg("cannot list")
	}
}

func (d *DevContainer) Man(_ *cobra.Command, _ []string) {
	header := &doc.GenManHeader{}
	err := doc.GenManTree(rootCmd, header, manOutDir)
	if err != nil {
		log.Error().Err(err).Send()
		os.Exit(1)
	}
}

func (d *DevContainer) Shell(cmd *cobra.Command, args []string) {
	// ensure it is started before starting a shell
	d.Start(cmd, args)
	// run asynchronously to avoid blocking shell start
	go d.PostAttachCommand()
	if _, err := d.Engine.Exec([]string{shellBin}, true, false); err != nil {
		log.Fatal().Err(err).Msg("cannot execute a shell")
	}
}

func (d *DevContainer) Start(_ *cobra.Command, _ []string) {
	if _, err := d.Engine.Start(); err != nil {
		log.Fatal().Err(err).Msg("cannot start")
	}
	d.PostStartCommand()
}

func (d *DevContainer) Stop(_ *cobra.Command, _ []string) {
	if _, err := d.Engine.Stop(); err != nil {
		log.Fatal().Err(err).Msg("cannot stop")
	}
	if stopRemove {
		if _, err := d.Engine.Remove(); err != nil {
			log.Fatal().Err(err).Msg("cannot remove")
		}
	}
}

// INIT/POST/ON STEPS

func (d *DevContainer) InitializeCommand() {
	if len(devc.JSON.InitializeCommand) > 0 {
		if _, err := execCmd(devc.JSON.InitializeCommand, false); err != nil {
			log.Fatal().Err(err).Msg("cannot run initializeCommand")
		}
	}
}

func (d *DevContainer) OnCreateCommand() {
	if len(d.JSON.OnCreateCommand) > 0 {
		if _, err := d.Engine.Exec(d.JSON.OnCreateCommand, true, false); err != nil {
			log.Fatal().Err(err).Msg("cannot execute onCreateCommand")
		}
	}
}

func (d *DevContainer) PostCreateCommand() {
	if len(d.JSON.PostCreateCommand) > 0 {
		if _, err := d.Engine.Exec(d.JSON.PostCreateCommand, true, false); err != nil {
			log.Fatal().Err(err).Msg("cannot execute postCreateCommand")
		}
	}
}

func (d *DevContainer) PostAttachCommand() {
	if len(devc.JSON.PostAttachCommand) > 0 {
		// wait a bit to ensure shell is started
		time.Sleep(1 * time.Second)
		if _, err := devc.Engine.Exec(devc.JSON.PostAttachCommand, true, false); err != nil {
			log.Fatal().Err(err).Msg("cannot execute postAttachCommand")
		}
	}
}

func (d *DevContainer) PostStartCommand() {
	if len(devc.JSON.PostStartCommand) > 0 {
		if _, err := devc.Engine.Exec(devc.JSON.PostStartCommand, true, false); err != nil {
			log.Fatal().Err(err).Msg("cannot execute postStartCommand")
		}
	}
}

// MAIN

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Send()
	}
}
