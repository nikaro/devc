package utils

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// GetConfig parse the devcontainer.json file and return a viper config object
func GetConfig() (*viper.Viper, error) {
	config := viper.New()
	config.AddConfigPath(".devcontainer/")
	config.SetConfigName("devcontainer")
	config.SetConfigType("json")

	// set aliases
	config.RegisterAlias("dockerfile", "build.dockerfile")
	config.RegisterAlias("context", "build.context")

	// set defaults
	path, _ := os.Getwd()
	dirName := filepath.Base(path)
	config.SetDefault("name", dirName)
	config.SetDefault("build.context", ".")
	config.SetDefault("updateRemoteUserUID", true)
	config.SetDefault("overrideCommand", true)

	err := config.ReadInConfig()

	return config, err
}

// CheckMutuallyExclusiveSettings does what its name says
func CheckMutuallyExclusiveSettings(config *viper.Viper) error {
	switch {
	case config.Get("image") != nil && config.Get("dockerFile") != nil:
		return errors.New("you cannot use both 'image' and 'dockerFile' settings")
	case config.Get("image") != nil && config.Get("build.dockerfile") != nil:
		return errors.New("you cannot use both 'image' and 'build.dockerfile' settings")
	case config.Get("image") != nil && config.Get("dockerComposeFile") != nil:
		return errors.New("you cannot use both 'image' and 'dockerComposeFile' settings")
	case config.Get("dockerFile") != nil && config.Get("dockerComposeFile") != nil:
		return errors.New("you cannot use both 'dockerFile' and 'dockerComposeFile' settings")
	case config.Get("build.dockerfile") != nil && config.Get("dockerComposeFile") != nil:
		return errors.New("you cannot use both 'build.dockerfile' and 'dockerComposeFile' settings")
	default:
		return nil
	}
}
