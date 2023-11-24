package cmd

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Copy bool
	Log  *LogConfig
}

type LogConfig struct {
	Level string
}

func NewConfig() (*Config, error) {
	var config *Config

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}
