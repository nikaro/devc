package cmd

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List containers",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		command = append(command, "ps")
		Run(command)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
