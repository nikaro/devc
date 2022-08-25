package main

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

func (d *DevContainer) Build() {
	if _, err := d.Engine.Build(); err != nil {
		log.Fatal().Err(err).Msg("cannot build")
	}
	if _, err := d.Engine.Create(); err != nil {
		log.Fatal().Err(err).Msg("cannot create")
	}
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build devcontainer",
	Run: func(_ *cobra.Command, _ []string) {
		devc.Build()
	},
}
