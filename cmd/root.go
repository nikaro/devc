package cmd

import (
	"fmt"
	"os"

	"git.sr.ht/~nka/devc/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootBackend string
var rootCommand []string
var rootConfig *viper.Viper
var rootVerbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:              "devc",
	Short:            "A CLI tool to manage your devcontainers using Docker-Compose",
	Long:             ``,
	PersistentPreRun: preRun,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&rootVerbose, "verbose", "v", false, "show commands used")
}

func preRun(command *cobra.Command, args []string) {
	if _, err := os.Stat(".devcontainer/devcontainer.json"); err == nil {
		rootConfig, err = utils.GetConfig()
		if err != nil {
			panic(fmt.Errorf("fatal error config file: %s", err))
		}

		if err := utils.CheckMutuallyExclusiveSettings(rootConfig); err != nil {
			panic(err)
		}

		// determine the command to use
		if rootConfig.Get("dockerComposeFile") != nil {
			rootBackend = "dockerCompose"
		} else if rootConfig.Get("build.dockerfile") != nil {
			rootBackend = "docker"
		} else {
			panic("cannot determine which command to use")
		}
	} else if len(os.Args) > 1 && os.Args[1] != "completion" && os.Args[1] != "-h" && os.Args[1] != "--help" && os.Args[1] != "init" {
		// allow usage of help and completion even when no devcontainer config is found
		fmt.Println("devcontainer settings not found")
		os.Exit(1)
	}
}
