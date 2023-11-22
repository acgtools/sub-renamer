package cmd

import (
	"errors"
	"github.com/dreamjz/sub-renamer/pkg/episode"
	"github.com/dreamjz/sub-renamer/pkg/log"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

var logLevel string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sub-renamer",
	Short: "Auto-rename video and subtitle files",
	Long:  `TODO://`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("not enough args")
		}

		logLevel, err := log.ParseLogLevel(logLevel)
		if err != nil {
			return err
		}
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
		slog.SetDefault(logger)

		vidDir, subDir := args[0], args[1]
		if !filepath.IsAbs(vidDir) {
			vidDir, err = filepath.Abs(vidDir)
			if err != nil {
				return err
			}
		}
		if !filepath.IsAbs(subDir) {
			subDir, err = filepath.Abs(subDir)
			if err != nil {
				return err
			}
		}

		vidMap, err := episode.ParseEpisodes(vidDir)
		if err != nil {
			return err
		}
		subMap, err := episode.ParseEpisodes(subDir)
		if err != nil {
			return err
		}

		for ep, vidName := range vidMap {
			subName, ok := subMap[ep]
			if !ok {
				continue
			}

			subExt := filepath.Ext(subName)

			newSubName := strings.TrimSuffix(vidName, filepath.Ext(vidName))

			err := os.Rename(filepath.Join(subDir, subName), filepath.Join(subDir, newSubName+subExt))
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "debug", "log level")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
