package cmd

import (
	"errors"
	"log/slog"
	"os"

	"github.com/acgtools/sub-renamer/pkg/episode"
	"github.com/acgtools/sub-renamer/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const minArgNum = 2

var version = "dev"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "sub-renamer",
	Short:   "Auto rename subtitle files to match video files",
	Long:    "sub-renamer <video dir> <sub dir>",
	Version: version,
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
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					return slog.Attr{}
				}
				return a
			},
		}))
		slog.SetDefault(logger)

		return episode.AutoRename(args[0], args[1], config.Copy) //nolint:wrapcheck
	},
}

func init() { //nolint:gochecknoinits
	rootCmd.PersistentFlags().String("log-level", "info", "log level, options: debug, info, warn, error")
	rootCmd.PersistentFlags().BoolP("copy", "c", false, "copy subtitles after renaming")

	_ = viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level"))
	_ = viper.BindPFlag("copy", rootCmd.PersistentFlags().Lookup("copy"))
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
