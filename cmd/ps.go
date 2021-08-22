package cmd

import (
	"github.com/nikaro/devc/backend/docker"
	"github.com/nikaro/devc/backend/dockercompose"
	"github.com/spf13/cobra"
)

var psArgs []string

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List containers",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		switch rootBackend {
		case "dockerCompose":
			projectName := rootConfig.GetString("name") + "_devcontainer"
			dockerComposeFile := ".devcontainer/" + rootConfig.GetString("dockercomposefile")
			dockercompose.List(rootVerbose, projectName, dockerComposeFile, psArgs...)
		case "docker":
			psDocker := docker.New()
			psDocker.SetVerbose(rootVerbose)
			psDocker.List()
		}
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
}
