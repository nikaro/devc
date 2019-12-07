package cmd

import (
	"github.com/spf13/cobra"
)

var execCommand []string

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a command inside a running container",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		rootCommand = append(rootCommand, "exec", rootServiceName)
		rootCommand = append(rootCommand, execCommand...)
		Run(rootCommand)
	},
}

func init() {
	execCmd.PersistentFlags().StringVarP(&rootServiceName, "service", "", "", "override service name")
	execCmd.PersistentFlags().StringArrayVarP(&execCommand, "command", "", []string{}, "command and its arguments (required)")
	execCmd.MarkPersistentFlagRequired("command")

	rootCmd.AddCommand(execCmd)
}
