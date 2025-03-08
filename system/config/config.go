package config

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"

	"github.com/spf13/viper"
)

var HomeDirPath, _ = homedir.Dir()

const configDirName = ".alpaca"

const configName = "config.json"

func GetConfigPath() string {
	if viper.GetString("metadata") != "" {
		MetadataPath := viper.GetString("metadata")
		return filepath.Join(MetadataPath, configName)
	}
	return filepath.Join(HomeDirPath, configDirName, configName)
}
