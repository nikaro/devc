package cmd

import (
	"github.com/nikaro/devc/backend/docker"
	"github.com/nikaro/devc/backend/dockercompose"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var shellBin string

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Execute a shell inside the running devcontainer",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// ensure it is started
		startCmd.Run(cmd, args)
		shellArgs := shellGetCommonArgs(rootConfig)
		if shellBin == "" {
			shellBin = rootConfig.GetString("shell")
		}
		switch rootBackend {
		case "dockerCompose":
			projectName := rootConfig.GetString("name") + "_devcontainer"
			dockerComposeFile := ".devcontainer/" + rootConfig.GetString("dockercomposefile")
			shellDockerComposeArgs := append(shellArgs, shellGetDockerComposeArgs(rootConfig)...)
			dockercompose.Exec(rootVerbose, projectName, dockerComposeFile, shellDockerComposeArgs...)
		case "docker":
			shellDockerArgs := append(shellArgs, shellGetDockerArgs(rootConfig)...)
			shellDocker := docker.New()
			shellDocker.SetVerbose(rootVerbose)
			shellDocker.SetArgs(shellDockerArgs)
			shellDocker.SetCommand([]string{shellBin})
			shellDocker.Exec()
		}
	},
}

func shellGetCommonArgs(config *viper.Viper) (args []string) {
	// append environment variables
	for remoteEnvKey, remoteEnvValue := range config.GetStringMapString("remoteEnv") {
		args = append(args, "--env", remoteEnvKey+"="+remoteEnvValue)
	}
	// append user
	if remoteUser := config.GetString("remoteUser"); remoteUser != "" {
		args = append(args, "--user", remoteUser)
	} else if containerUser := config.GetString("containerUser"); containerUser != "" {
		// fallback to containerUser if set
		args = append(args, "--user", containerUser)
	}
	// append workspace
	if workspaceFolder := config.GetString("workspaceFolder"); workspaceFolder != "" {
		args = append(args, "--workdir", workspaceFolder)
	}

	return args
}

func shellGetDockerArgs(config *viper.Viper) (args []string) {
	args = append(args, "--interactive", "--tty")

	return args
}

func shellGetDockerComposeArgs(config *viper.Viper) (args []string) {
	// append service to start
	args = append(args, config.GetString("service"))
	// append shell
	args = append(args, shellBin)

	return args
}

func init() {
	shellCmd.PersistentFlags().StringVarP(&shellBin, "shell", "s", "", "override shell")

	rootCmd.AddCommand(shellCmd)
}
