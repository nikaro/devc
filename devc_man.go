package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var manOutDir string

func init() {
	manCmd.PersistentFlags().StringVarP(&manOutDir, "output", "o", "man", "output directory")
	rootCmd.AddCommand(manCmd)
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}

func (d *DevContainer) Man() {
	header := &doc.GenManHeader{}
	err := doc.GenManTree(rootCmd, header, manOutDir)
	if err != nil {
		log.Error().Err(err).Send()
		os.Exit(1)
	}
}

var manCmd = &cobra.Command{
	Use:    "man",
	Short:  "Generate manpage",
	Hidden: true,
	Run: func(_ *cobra.Command, _ []string) {
		devc.Man()
	},
}
