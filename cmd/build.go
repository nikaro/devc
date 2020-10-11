package cmd

import (
	"os"

	"git.sr.ht/nka/devc/backend/docker"
	"git.sr.ht/nka/devc/backend/dockercompose"
	"github.com/spf13/cobra"
)

var buildArgs []string

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build or rebuild devcontainer services",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		switch rootBackend {
		case "dockerCompose":
			projectName := rootConfig.GetString("name") + "_devcontainer"
			dockerComposeFile := ".devcontainer/" + rootConfig.GetString("dockercomposefile")
			dockercompose.Build(rootVerbose, projectName, dockerComposeFile, buildArgs...)
		case "docker":
			path, _ := os.Getwd()
			image := docker.GetImageName(path)
			dockerFile := ".devcontainer/" + rootConfig.GetString("build.dockerfile")
			context := rootConfig.GetString("build.context")
			// append build args
			buildArgsConfig := rootConfig.GetStringMapString("build.args")
			for _, buildArg := range buildArgsConfig {
				buildArgs = append(buildArgs, "--build-arg", buildArg+"="+buildArgsConfig[buildArg])
			}
			// append target
			if target := rootConfig.GetString("build.target"); target != "" {
				buildArgs = append(buildArgs, "--target", target)
			}
			docker.Build(rootVerbose, image, dockerFile, context, buildArgs...)
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
