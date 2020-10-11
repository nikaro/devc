package cmd

import (
	"os"

	"git.sr.ht/nka/devc/backend/docker"
	"git.sr.ht/nka/devc/backend/dockercompose"
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
			path, _ := os.Getwd()
			docker.List(rootVerbose, path, psArgs...)
		}
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
}
