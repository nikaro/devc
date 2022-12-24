// devc is a devcontainer managment tool

// Do not want to mess your minimal, clean and fine tuned OS installation with a
// bunch of compilers, linters, debuggers, interpreters, minifiers, unminifiers,
// beautifiers, etc...? Visual Studio Code nailed it with DevContainers. But do
// not want to use VSCode neither? That's where devc comes in, it a simple CLI
// that wrap docker/docker-compose and run the (almost) same commands that
// VSCode runs behind the scenes.
package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

// global log var
var log zerolog.Logger

// interface that each engine (docker, compose, podman, k8s, etc.) must implement
type Engine interface {
	Init(config *DevContainer) error
	IsBuilt() (bool, error)
	IsCreated() (bool, error)
	IsRunning() (bool, error)
	Build() (string, error)
	Create() (string, error)
	Remove() (string, error)
	Start() (string, error)
	Stop() (string, error)
	List() (string, error)
	Run(command []string) (string, error)
	Exec(command []string) (string, error)
	ResolveEnv(env string) string
}

// devcontainer meta-structure
type DevContainer struct {
	ConfigDir            string
	Config               *viper.Viper
	Engine               Engine
	WorkingDirectoryPath string
	WorkingDirectoryName string
}

var version string

// global devcontainer var
var devc DevContainer

// DEVC COMMANDS

// cli args
var rootConfigDir string
var rootVerbose int
var manOutDir string
var shellBin string
var stopRemove bool

func init() {
	// devc command
	rootCmd.PersistentFlags().StringVarP(&rootConfigDir, "config-dir", "c", ".devcontainer", "custom devcontainer directory")
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
		d.SetAliases()
		d.NormalizeTypes()
		d.SetDefaults()
		d.CheckConfig()
		d.InitializeCommand()
		d.ResolveVars()
		d.SetEngine()
	}
}

func (d *DevContainer) Build(cmd *cobra.Command, args []string) {
	if built, _ := d.Engine.IsBuilt(); !built {
		if _, err := d.Engine.Build(); err != nil {
			log.Fatal().Err(err).Msg("cannot build")
		}
	}
}

func (d *DevContainer) Init(_ *cobra.Command, _ []string) {
	if err := os.Mkdir(d.ConfigDir, 0755); err != nil {
		log.Fatal().Err(err).Msg("cannot create directory")
	}
	log.Info().Msg("directory created")
	d.Config.Set("image", "alpine:latest")
	if err := d.Config.SafeWriteConfigAs(filepath.Join(d.ConfigDir, "devcontainer.json")); err != nil {
		log.Fatal().Err(err).Msg("cannot write file")
	}
	log.Info().Msg("file created")
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
	// ensure container is started before starting a shell
	d.Start(cmd, args)
	// run post command asynchronously to avoid blocking shell start
	go d.PostAttachCommand()
	if _, err := d.Engine.Exec([]string{shellBin}); err != nil {
		log.Fatal().Err(err).Msg("cannot execute a shell")
	}
}

func (d *DevContainer) Start(cmd *cobra.Command, args []string) {
	created, _ := d.Engine.IsCreated()
	if !created {
		if _, err := d.Engine.Create(); err != nil {
			log.Fatal().Err(err).Msg("cannot create")
		}
	}
	if running, _ := d.Engine.IsRunning(); !running {
		if _, err := d.Engine.Start(); err != nil {
			log.Fatal().Err(err).Msg("cannot start")
		}
		if !created {
			d.OnCreateCommand()
			// TODO: figure out what updateContentCommand does and add it
			// go d.UpdateContent()
			d.PostCreateCommand()
		}
		d.PostStartCommand()
	}
}

func (d *DevContainer) Stop(_ *cobra.Command, _ []string) {
	if created, _ := d.Engine.IsCreated(); created {
		if running, _ := d.Engine.IsRunning(); running {
			if _, err := d.Engine.Stop(); err != nil {
				log.Fatal().Err(err).Msg("cannot stop")
			}
		}
		if stopRemove {
			if _, err := d.Engine.Remove(); err != nil {
				log.Fatal().Err(err).Msg("cannot remove")
			}
		}
	}
}

// INIT/POST/ON STEPS

func (d *DevContainer) InitializeCommand() {
	var cmd []string
	switch d.Config.Get("initializeCommand").(type) {
	case string:
		cmd = []string{"sh", "-c", d.Config.GetString("initializeCommand")}
	case []interface{}:
		cmd = d.Config.GetStringSlice("initializeCommand")
	}
	if len(cmd) > 0 {
		// execute on the host
		if _, err := execCmd(cmd, false); err != nil {
			log.Fatal().Err(err).Msgf("cannot run %s", "initializeCommand")
		}
	}
}

func (d *DevContainer) cmd(step string, wait bool) {
	var cmd []string
	switch d.Config.Get(step).(type) {
	case string:
		cmd = []string{"sh", "-c", d.Config.GetString(step)}
	case []interface{}:
		cmd = d.Config.GetStringSlice(step)
	}
	if len(cmd) > 0 {
		// wait a bit to ensure shell is started
		if wait {
			time.Sleep(1 * time.Second)
		}
		// execute inside the container
		if _, err := d.Engine.Exec(cmd); err != nil {
			log.Fatal().Err(err).Msgf("cannot run %s", step)
		}
	}
}

func (d *DevContainer) OnCreateCommand() {
	d.cmd("onCreateCommand", false)
}

func (d *DevContainer) PostCreateCommand() {
	d.cmd("postCreateCommand", false)
}

func (d *DevContainer) PostAttachCommand() {
	d.cmd("postAttachCommand", true)
}

func (d *DevContainer) PostStartCommand() {
	d.cmd("postStartCommand", false)
}

// MAIN

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Send()
	}
}
