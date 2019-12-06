package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var projectName string
var dockerComposeFile string
var serviceName string
var command []string

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
	rootCmd.PersistentFlags().StringVarP(&projectName, "project-name", "p", "", "alternate project name")
	rootCmd.PersistentFlags().StringVarP(&dockerComposeFile, "file", "f", "", "alternate Compose file")

	if projectName == "" {
		projectName = GetConfig("name")
	}
	if dockerComposeFile == "" {
		dockerComposeFile = GetConfig("dockerComposeFile")
	}
	serviceName = GetConfig("service")
	command = append(
		command,
		"docker-compose",
		"-p", projectName,
		"-f", dockerComposeFile,
	)
}
