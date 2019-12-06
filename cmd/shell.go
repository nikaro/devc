package cmd

import (
	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Execute a bash shell inside the running devcontainer",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		command = append(command, "exec", serviceName, "bash")
		Run(command)
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}
