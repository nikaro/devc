package cmd

import (
	"os"

	"git.sr.ht/~nka/devc/backend/docker"
	"git.sr.ht/~nka/devc/backend/dockercompose"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startUp bool
var startArgs []string

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start devcontainer services",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		switch rootBackend {
		case "dockerCompose":
			projectName := rootConfig.GetString("name") + "_devcontainer"
			dockerComposeFile := ".devcontainer/" + rootConfig.GetString("dockercomposefile")
			// append additonal services to start
			for _, service := range rootConfig.GetStringSlice("runServices") {
				startArgs = append(startArgs, service)
			}
			// use "up" if it does not exist or if flag is used
			if containers, _ := dockercompose.GetContainers(projectName, dockerComposeFile); len(containers) == 0 || startUp {
				startArgs = append(startArgs, "--detach")
				dockercompose.Up(rootVerbose, projectName, dockerComposeFile, startArgs...)
			} else {
				dockercompose.Start(rootVerbose, projectName, dockerComposeFile, startArgs...)
			}
		case "docker":
			startDocker := docker.New()
			startDocker.SetVerbose(rootVerbose)
			// create if it does not exist
			if startDocker.GetContainer() == "" {
				if image := rootConfig.GetString("image"); image != "" {
					startDocker.SetImage(image) // use image from devcontainer.json
				} else {
					buildCmd.Run(cmd, args) // ensure image is built
				}
				if rootConfig.GetBool("overrideCommand") {
					startDocker.SetCommand([]string{"/bin/sh", "-c", "while sleep 1000; do :; done"})
				}
				startDocker.SetArgs(createGetDockerArgs(rootConfig)) // update args for container creation
				startDocker.Create()                                 // ensure container is created
				startDocker.SetArgs([]string{})                      // empty args previously used for creation
				startDocker.SetContainer("")                         // update container reference with the created one
			}
			startDocker.Start()
		}
	},
}

func createGetDockerArgs(config *viper.Viper) (args []string) {
	path, _ := os.Getwd()
	args = append(args, "--label", "vsch.local.folder="+path)
	if mount := config.GetString("workspaceMount"); mount != "" {
		args = append(args, "--mount", mount)
	}
	for _, mount := range config.GetStringSlice("mounts") {
		args = append(args, "--mount", mount)
	}
	for _, port := range config.GetStringSlice("appPort") {
		args = append(args, "--publish", port)
	}
	for _, port := range config.GetStringSlice("forwardPorts") {
		args = append(args, "--publish", port)
	}
	for containerEnvKey, containerEnvValue := range config.GetStringMapString("containerEnv") {
		startArgs = append(startArgs, "--env", containerEnvKey+"="+containerEnvValue)
	}
	if containerUser := config.GetString("containerUser"); containerUser != "" {
		args = append(args, "--user", containerUser)
	}
	for _, runArg := range config.GetStringSlice("runArgs") {
		args = append(args, runArg)
	}

	return args
}

func init() {
	startCmd.PersistentFlags().BoolVarP(&startUp, "up", "u", false, "when docker-compose is used, create and start containers")

	rootCmd.AddCommand(startCmd)
}
