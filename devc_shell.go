package main

import (
	"github.com/spf13/cobra"
)

var shellBin string

func init() {
	shellCmd.PersistentFlags().StringVarP(&shellBin, "shell", "s", "sh", "override shell")
	rootCmd.AddCommand(shellCmd)
}

func (d *DevContainer) Shell() {
	if _, err := d.Engine.Exec([]string{shellBin}, true, false); err != nil {
		log.Fatal().Err(err).Msg("cannot execute a shell")
	}
}

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Execute a shell inside devcontainer",
	Run: func(_ *cobra.Command, _ []string) {
		devc.Shell()
	},
}
