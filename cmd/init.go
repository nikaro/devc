package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initName string
var initDockerfile string
var initWorkspace string
var initMount string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create an initial devcontainer configuration",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		os.Mkdir(".devcontainer", 0755)
		mkDevcontainer()
		mkDockerfile()
	},
}

func init() {
	path, _ := os.Getwd()
	initCmd.PersistentFlags().StringVarP(&initName, "name", "n", "myprojectname", "project name")
	initCmd.PersistentFlags().StringVarP(&initDockerfile, "dockerfile", "d", "Dockerfile", "path to Dockerfile")
	initCmd.PersistentFlags().StringVarP(&initWorkspace, "workspace", "w", "/workspace", "working directory inside the container")
	initCmd.PersistentFlags().StringVarP(&initMount, "mount", "m", path, "path to mount into container workspace")

	rootCmd.AddCommand(initCmd)
}

func mkDevcontainer() {
	devcontainer := viper.New()
	devcontainer.AddConfigPath(".devcontainer/")
	devcontainer.SetConfigName("devcontainer")
	devcontainer.SetConfigType("json")
	devcontainer.SetDefault("name", initName)
	devcontainer.SetDefault("build.dockerfile", initDockerfile)
	devcontainer.SetDefault("workspaceFolder", initWorkspace)
	devcontainer.SetDefault("workspaceMount", "type=bind,src="+initMount+",dst="+initWorkspace)
	devcontainer.SetDefault("extensions", []string{})
	devcontainer.SetDefault("settings", map[string]string{})
	devcontainer.SafeWriteConfigAs(".devcontainer/devcontainer.json")
}

func mkDockerfile() {
	dockerfile, _ := os.Create(".devcontainer/Dockerfile")
	defer dockerfile.Close()
	dockerfile.WriteString("FROM alpine:edge\n\nRUN apk add --no-cache bash bash-completion\n")
}
