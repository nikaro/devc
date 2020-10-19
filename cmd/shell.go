package cmd

import (
	"os"

	"git.sr.ht/nka/devc/backend/docker"
	"git.sr.ht/nka/devc/backend/dockercompose"
	"github.com/spf13/cobra"
)

var shellBin string
var shellArgs []string

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Execute a shell inside the running devcontainer",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// ensure it is started
		startCmd.Run(cmd, args)
		// append environment variables
		for remoteEnvKey, remoteEnvValue := range rootConfig.GetStringMapString("remoteEnv") {
			shellArgs = append(shellArgs, "--env", remoteEnvKey+"="+remoteEnvValue)
		}
		// append user
		if remoteUser := rootConfig.GetString("remoteUser"); remoteUser != "" {
			shellArgs = append(shellArgs, "--user", remoteUser)
		} else if containerUser := rootConfig.GetString("containerUser"); containerUser != "" {
			// fallback to containerUser if set
			shellArgs = append(shellArgs, "--user", containerUser)
		}
		// append workspace
		if workspaceFolder := rootConfig.GetString("workspaceFolder"); workspaceFolder != "" {
			shellArgs = append(shellArgs, "--workdir", workspaceFolder)
		}
		switch rootBackend {
		case "dockerCompose":
			projectName := rootConfig.GetString("name") + "_devcontainer"
			dockerComposeFile := ".devcontainer/" + rootConfig.GetString("dockercomposefile")
			// append service to start
			serviceName := rootConfig.GetString("service")
			shellArgs = append(shellArgs, serviceName)
			// append shell
			shellArgs = append(shellArgs, shellBin)
			// call command with args
			dockercompose.Exec(rootVerbose, projectName, dockerComposeFile, shellArgs...)
		case "docker":
			// get container name
			path, _ := os.Getwd()
			container, _ := docker.GetContainer(path)
			// append args
			shellArgs = append(shellArgs, "--interactive", "--tty")
			// call command with args
			docker.Exec(rootVerbose, container, shellBin, shellArgs...)
		}
	},
}

func init() {
	shellCmd.PersistentFlags().StringVarP(&shellBin, "shell", "", "sh", "override shell")

	rootCmd.AddCommand(shellCmd)
}
