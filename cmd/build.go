package cmd

import (
	"github.com/spf13/cobra"
)

var buildCompress bool
var buildForceRemove bool
var buildNoCache bool
var buildPull bool

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build or rebuild devcontainer services",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		rootCommand = append(rootCommand, "build")
		if buildCompress {
			rootCommand = append(rootCommand, "--compress")
		}
		if buildForceRemove {
			rootCommand = append(rootCommand, "--force-rm")
		}
		if buildNoCache {
			rootCommand = append(rootCommand, "--no-cache")
		}
		if buildPull {
			rootCommand = append(rootCommand, "--pull")
		}
		Run(rootCommand, rootVerbose)
	},
}

func init() {
	buildCmd.PersistentFlags().BoolVarP(&buildCompress, "compress", "", false, "compress the build context using gzip")
	buildCmd.PersistentFlags().BoolVarP(&buildForceRemove, "force-rm", "", false, "always remove intermediate containers")
	buildCmd.PersistentFlags().BoolVarP(&buildNoCache, "no-cache", "", false, "do not use cache when building the image")
	buildCmd.PersistentFlags().BoolVarP(&buildPull, "pull", "", false, "always attempt to pull a newer version of the image")

	rootCmd.AddCommand(buildCmd)
}
