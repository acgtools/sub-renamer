package cmd

import "github.com/spf13/viper"

type Config struct {
	VidExt []string // video extension
	SubExt []string // subtitle extension
	Log    *LogConfig
}

type LogConfig struct {
	Level string
}

func NewConfig() (*Config, error) {
	var config *Config

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}
