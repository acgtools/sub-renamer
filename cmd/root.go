package cmd

import (
	"errors"
	"log/slog"
	"os"

	"github.com/dreamjz/sub-renamer/pkg/episode"
	"github.com/dreamjz/sub-renamer/pkg/log"
	"github.com/spf13/cobra"
)

const minArgNum = 2

var logLevel string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sub-renamer",
	Short: "Auto rename subtitle files to match video files",
	Long:  "sub-renamer <video dir> <sub dir>",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < minArgNum {
			return errors.New("not enough args")
		}

		logLevel, err := log.ParseLevel(logLevel)
		if err != nil {
			return err //nolint:wrapcheck
		}
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
		slog.SetDefault(logger)

		return episode.AutoRename(args[0], args[1]) //nolint:wrapcheck
	},
}

func init() { //nolint:gochecknoinits
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
