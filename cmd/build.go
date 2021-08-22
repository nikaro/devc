package cmd

import (
	"github.com/nikaro/devc/backend/docker"
	"github.com/nikaro/devc/backend/dockercompose"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build or rebuild devcontainer services",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		switch rootBackend {
		case "dockerCompose":
			projectName := rootConfig.GetString("name") + "_devcontainer"
			dockerComposeFile := ".devcontainer/" + rootConfig.GetString("dockercomposefile")
			dockercompose.Build(rootVerbose, projectName, dockerComposeFile)
		case "docker":
			buildDocker := docker.New()
			buildDocker.SetVerbose(rootVerbose)
			buildDocker.SetDockerfile(".devcontainer/" + rootConfig.GetString("build.dockerfile"))
			buildDocker.SetContext(rootConfig.GetString("build.context"))
			buildDocker.SetArgs(buildGetDockerArgs(rootConfig))
			buildDocker.Build()
		}
	},
}

func buildGetDockerArgs(config *viper.Viper) (args []string) {
	// append build args
	for k, v := range config.GetStringMapString("build.args") {
		args = append(args, "--build-arg", k+"="+v)
	}
	// append target
	if target := config.GetString("build.target"); target != "" {
		args = append(args, "--target", target)
	}

	return args
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
