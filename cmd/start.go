package cmd

import (
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start devcontainer services",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		rootCommand = append(rootCommand, "start")
		Run(rootCommand, rootVerbose)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
