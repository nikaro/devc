package cmd

import (
	"github.com/spf13/cobra"
)

var downRemoveImages bool
var downRemoveVolumes bool
var downRemoveOrphans bool

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop devcontainer services",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		rootCommand = append(rootCommand, "down")
		if downRemoveImages {
			rootCommand = append(rootCommand, "--rmi", "all")
		}
		if downRemoveVolumes {
			rootCommand = append(rootCommand, "--volumes")
		}
		if downRemoveOrphans {
			rootCommand = append(rootCommand, "--remove-orphans")
		}
		Run(rootCommand)
	},
}

func init() {
	downCmd.PersistentFlags().BoolVarP(&downRemoveImages, "remove-images", "", false, "remove all images used by any service")
	downCmd.PersistentFlags().BoolVarP(&downRemoveVolumes, "remove-volumes", "", false, "remove named volumes declared in the volumes section of the Compose file and anonymous volumes attached to containers")
	downCmd.PersistentFlags().BoolVarP(&downRemoveOrphans, "remove-orphans", "", false, "remove containers for services not defined in the Compose file")

	rootCmd.AddCommand(downCmd)
}
