package main

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

func (d *DevContainer) Start() {
	if _, err := d.Engine.Start(); err != nil {
		log.Fatal().Err(err).Msg("cannot start")
	}
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start devcontainer",
	Run: func(_ *cobra.Command, _ []string) {
		devc.Start()
		if len(devc.JSON.PostStartCommand) > 0 {
			if _, err := devc.Engine.Exec(devc.JSON.PostStartCommand, false, false); err != nil {
				log.Fatal().Err(err).Msg("cannot execute postStartCommand")
			}
		}
	},
}
