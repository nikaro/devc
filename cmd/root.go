package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootProjectName string
var rootDockerComposeFile string
var rootServiceName string
var rootCommand []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "devc",
	Short: "A CLI tool to manage your devcontainers using Docker-Compose",
	Long:  ``,
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
	rootCmd.PersistentFlags().StringVarP(&rootProjectName, "project-name", "p", "", "alternate project name")
	rootCmd.PersistentFlags().StringVarP(&rootDockerComposeFile, "file", "f", "", "alternate Compose file")

	if _, err := os.Stat(".devcontainer/"); err == nil {
		// load settings only if devcontainer configuration is found
		if rootProjectName == "" {
			rootProjectName = GetConfig("name")
		}
		if rootDockerComposeFile == "" {
			rootDockerComposeFile = GetConfig("dockerComposeFile")
		}
		rootServiceName = GetConfig("service")
		rootCommand = append(
			rootCommand,
			"docker-compose",
			"-p", rootProjectName,
			"-f", rootDockerComposeFile,
		)
	} else if len(os.Args) > 1 && os.Args[1] != "completion" && os.Args[1] != "-h" && os.Args[1] != "--help" {
		// allow usage of help and completion even when no devcontainer config is found
		fmt.Println("devcontainer directory not found")
		os.Exit(1)
	}
}
