package main

import (
	"github.com/spf13/cobra"
)

var stopRemove bool

func init() {
	stopCmd.PersistentFlags().BoolVarP(&stopRemove, "remove", "r", false, "remove containers and networks")
	rootCmd.AddCommand(stopCmd)
}

func (d *DevContainer) Stop() {
	if _, err := d.Engine.Stop(); err != nil {
		log.Fatal().Err(err).Msg("cannot stop")
	}
	if stopRemove {
		if _, err := d.Engine.Remove(); err != nil {
			log.Fatal().Err(err).Msg("cannot remove")
		}
	}
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop devcontainer",
	Run: func(_ *cobra.Command, _ []string) {
		devc.Stop()
	},
}
