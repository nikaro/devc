package cmd

import (
	"github.com/spf13/cobra"
)

var exeCommand []string

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a command inside a running container",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		command = append(command, "exec", serviceName)
		command = append(command, exeCommand...)
		Run(command)
	},
}

func init() {
	execCmd.PersistentFlags().StringVarP(&serviceName, "service", "s", "", "override service name")
	execCmd.PersistentFlags().StringArrayVarP(&exeCommand, "command", "c", []string{}, "command and its arguments (required)")
	execCmd.MarkPersistentFlagRequired("command")

	rootCmd.AddCommand(execCmd)
}
