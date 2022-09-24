package main

import (
	"time"

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

func (d *DevContainer) shellPostAttach() {
	if len(devc.JSON.PostAttachCommand) > 0 {
		// wait a bit to ensure shell is started
		time.Sleep(1 * time.Second)
		if _, err := devc.Engine.Exec(devc.JSON.PostAttachCommand, true, false); err != nil {
			log.Fatal().Err(err).Msg("cannot execute postAttachCommand")
		}
	}
}

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Execute a shell inside devcontainer",
	Run: func(_ *cobra.Command, _ []string) {
		// ensure it is started before attaching starting a shell
		devc.Start()
		// run asynchronously to avoid blocking shell attach
		go devc.shellPostAttach()
		devc.Shell()
	},
}
