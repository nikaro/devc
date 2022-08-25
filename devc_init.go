package main

import (
	"encoding/json"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var initJSON []byte

func init() {
	rootCmd.AddCommand(initCmd)
}

func (d *DevContainer) Init() {
	fs := afero.NewOsFs()
	if exists, _ := afero.DirExists(fs, ".devcontainer"); !exists {
		if err := fs.Mkdir(".devcontainer", 0755); err != nil {
			log.Fatal().Err(err).Msg("cannot create .devcontainer directory")
		}
		log.Info().Msg(".devcontainer directory created")
	}
	if exists, _ := afero.Exists(fs, ".devcontainer/devcontainer.json"); !exists {
		d.JSON.Image = "alpine:latest"
		initJSON, _ = json.MarshalIndent(d.JSON, "", "  ")
		if err := afero.WriteFile(fs, ".devcontainer/devcontainer.json", initJSON, 0644); err != nil {
			log.Fatal().Err(err).Msg("cannot write devcontainer.json file")
		}
		log.Info().Msg("devcontainer.json file created")
	}
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize devcontainer configuration",
	Run: func(_ *cobra.Command, _ []string) {
		devc.Init()
	},
}
