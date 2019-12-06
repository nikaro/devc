package cmd

import (
	"github.com/spf13/cobra"
)

var compress bool
var forceRemove bool
var noCache bool
var pull bool

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build or rebuild devcontainer services",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		command = append(command, "build")
		if compress {
			command = append(command, "--compress")
		}
		if forceRemove {
			command = append(command, "--force-rm")
		}
		if noCache {
			command = append(command, "--no-cache")
		}
		if pull {
			command = append(command, "--pull")
		}
		Run(command)
	},
}

func init() {
	buildCmd.PersistentFlags().BoolVarP(&compress, "compress", "z", false, "compress the build context using gzip")
	buildCmd.PersistentFlags().BoolVarP(&forceRemove, "force-rm", "r", false, "always remove intermediate containers")
	buildCmd.PersistentFlags().BoolVarP(&noCache, "no-cache", "n", false, "do not use cache when building the image")
	buildCmd.PersistentFlags().BoolVarP(&pull, "pull", "u", false, "always attempt to pull a newer version of the image")

	rootCmd.AddCommand(buildCmd)
}
