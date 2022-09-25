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

func (d *DevContainer) OnCreate() {
	if len(d.JSON.OnCreateCommand) > 0 {
		if _, err := d.Engine.Exec(d.JSON.OnCreateCommand, true, false); err != nil {
			log.Fatal().Err(err).Msg("cannot execute onCreateCommand")
		}
	}
}

func (d *DevContainer) PostCreate() {
	if len(d.JSON.PostCreateCommand) > 0 {
		if _, err := d.Engine.Exec(d.JSON.PostCreateCommand, true, false); err != nil {
			log.Fatal().Err(err).Msg("cannot execute postCreateCommand")
		}
	}
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build devcontainer",
	Run: func(_ *cobra.Command, _ []string) {
		devc.Build()
		devc.OnCreate()
		// TODO: figure out what updateContentCommand does and add it
		// go devc.UpdateContent()
		devc.PostCreate()
	},
}
