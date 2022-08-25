package main

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

func (d *DevContainer) List() {
	if _, err := d.Engine.List(); err != nil {
		log.Fatal().Err(err).Msg("cannot list")
	}
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "ps"},
	Short:   "List devcontainers",
	Run: func(_ *cobra.Command, _ []string) {
		devc.List()
	},
}
