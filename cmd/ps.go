package cmd

import (
	"github.com/spf13/cobra"
)

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List containers",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		rootCommand = append(rootCommand, "ps")
		Run(rootCommand, rootVerbose)
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
}
