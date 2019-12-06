package cmd

import (
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start devcontainer services",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		command = append(command, "up", "-d")
		Run(command)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
