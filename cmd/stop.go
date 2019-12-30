package cmd

import (
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop devcontainer services",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		rootCommand = append(rootCommand, "stop")
		Run(rootCommand, rootVerbose)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
