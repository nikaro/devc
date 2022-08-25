package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func init() {
	rootCmd.AddCommand(man)
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}

func (d *DevContainer) Man() {
	header := &doc.GenManHeader{}
	err := doc.GenMan(rootCmd, header, os.Stdout)
	if err != nil {
		log.Error().Err(err).Send()
		os.Exit(1)
	}
}

var man = &cobra.Command{
	Use:    "man",
	Short:  "Generate manpage",
	Hidden: true,
	Run: func(_ *cobra.Command, _ []string) {
		devc.Man()
	},
}
