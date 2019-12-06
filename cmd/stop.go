package cmd

import (
	"github.com/spf13/cobra"
)

var removeImages bool
var removeVolumes bool
var removeOrphans bool

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop devcontainer services",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		command = append(command, "down")
		if removeImages {
			command = append(command, "--rmi", "all")
		}
		if removeVolumes {
			command = append(command, "--volumes")
		}
		if removeOrphans {
			command = append(command, "--remove-orphans")
		}
		Run(command)
	},
}

func init() {
	stopCmd.PersistentFlags().BoolVarP(&removeImages, "remove-images", "i", false, "remove all images used by any service")
	stopCmd.PersistentFlags().BoolVarP(&removeVolumes, "remove-volumes", "v", false, "remove named volumes declared in the volumes section of the Compose file and anonymous volumes attached to containers")
	stopCmd.PersistentFlags().BoolVarP(&removeOrphans, "remove-orphans", "o", false, "remove containers for services not defined in the Compose file")

	rootCmd.AddCommand(stopCmd)
}
