package cmd

import (
	"os"

	"git.sr.ht/nka/devc/backend/docker"
	"git.sr.ht/nka/devc/backend/dockercompose"
	"github.com/spf13/cobra"
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
			path, _ := os.Getwd()
			container, _ := docker.GetContainer(path)
			// create if it does not exist
			if container == "" {
				image := rootConfig.GetString("image")
				if image != "" {
					image = docker.GetImageName(path)
					buildCmd.Run(cmd, args) // ensure image is built
				}
				var command []string
				if rootConfig.GetBool("overrideCommand") {
					command = append(command, "/bin/sh", "-c", "while sleep 1000; do :; done")
				}
				createArgs := []string{"--label", "vsch.local.folder=" + path}
				if mount := rootConfig.GetString("workspaceMount"); mount != "" {
					createArgs = append(createArgs, "--mount", mount)
				}
				for _, mount := range rootConfig.GetStringSlice("mounts") {
					createArgs = append(createArgs, "--mount", mount)
				}
				for _, port := range rootConfig.GetStringSlice("appPort") {
					createArgs = append(createArgs, "--publish", port)
				}
				for _, port := range rootConfig.GetStringSlice("forwardPorts") {
					createArgs = append(createArgs, "--publish", port)
				}
				for containerEnvKey, containerEnvValue := range rootConfig.GetStringMapString("containerEnv") {
					startArgs = append(startArgs, "--env", containerEnvKey+"="+containerEnvValue)
				}
				if containerUser := rootConfig.GetString("containerUser"); containerUser != "" {
					createArgs = append(createArgs, "--user", containerUser)
				}
				for _, runArg := range rootConfig.GetStringSlice("runArgs") {
					createArgs = append(createArgs, runArg)
				}
				docker.Create(rootVerbose, image, command, createArgs...)
				container, _ = docker.GetContainer(path)
			}
			docker.Start(rootVerbose, container, startArgs...)
		}
	},
}

func init() {
	startCmd.PersistentFlags().BoolVarP(&startUp, "up", "u", false, "when docker-compose is used, create and start containers")

	rootCmd.AddCommand(startCmd)
}
