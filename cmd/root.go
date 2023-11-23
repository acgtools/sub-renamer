package cmd

import (
	"errors"
	"log/slog"
	"os"

	"github.com/dreamjz/sub-renamer/pkg/episode"
	"github.com/dreamjz/sub-renamer/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const minArgNum = 2

var configFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sub-renamer",
	Short: "Auto rename subtitle files to match video files",
	Long:  "sub-renamer <video dir> <sub dir>",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < minArgNum {
			return errors.New("not enough args")
		}

		config, err := NewConfig()
		if err != nil {
			return err
		}

		logLevel, err := log.ParseLevel(config.Log.Level)
		if err != nil {
			return err //nolint:wrapcheck
		}
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
		slog.SetDefault(logger)

		return episode.AutoRename(args[0], args[1], config.VidExt, config.SubExt) //nolint:wrapcheck
	},
}

func init() { //nolint:gochecknoinits
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default to sub-renamer.yml)")
	rootCmd.PersistentFlags().String("log-level", "info", "log level")

	_ = viper.BindPFlag("log.level", rootCmd.Flags().Lookup("log-level"))
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("sub-renamer")
		viper.SetConfigType("yml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME")
	}

	if err := viper.ReadInConfig(); err != nil {
		var e viper.ConfigFileNotFoundError
		if !errors.As(err, &e) {
			slog.Error("error reading config file", "error", err)
			os.Exit(1)
		}
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
