package cmd

import (
	"github.com/spf13/cobra"
)

var upDetach bool

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Create and start devcontainer services",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		rootCommand = append(rootCommand, "up")
		if upDetach {
			rootCommand = append(rootCommand, "--detach")
		}
		Run(rootCommand, rootVerbose)
	},
}

func init() {
	upCmd.PersistentFlags().BoolVarP(&upDetach, "detach", "", true, "run containers in the background")

	rootCmd.AddCommand(upCmd)
}
