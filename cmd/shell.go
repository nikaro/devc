package cmd

import (
	"github.com/spf13/cobra"
)

var shellShell string

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Execute a shell inside the running devcontainer",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		rootCommand = append(rootCommand, "exec", rootServiceName, shellShell)
		Run(rootCommand, rootVerbose)
	},
}

func init() {
	shellCmd.PersistentFlags().StringVarP(&shellShell, "shell", "", "bash", "override shell")

	rootCmd.AddCommand(shellCmd)
}
