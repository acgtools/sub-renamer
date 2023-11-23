package episode_test

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/dreamjz/sub-renamer/cmd"
	"github.com/dreamjz/sub-renamer/pkg/episode"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

const (
	rootDir = "../../"

	testDir    = rootDir + "test/"
	exampleDir = testDir + "example/"
	originSub  = exampleDir + "origin-subDir/"
	vidDir     = exampleDir + "vid/"
	subDir     = exampleDir + "subDir/"

	cfgDir = rootDir
)

func TestMain(m *testing.M) {
	initConfig()
	clean()
	genVidSubFiles()

	code := m.Run()
	os.Exit(code)
}

func initConfig() {
	viper.SetConfigName("sub-renamer")
	viper.SetConfigType("yml")
	viper.AddConfigPath(cfgDir)

	if err := viper.ReadInConfig(); err != nil {
		var e viper.ConfigFileNotFoundError
		if !errors.As(err, &e) {
			log.Fatal("error reading config file", "error", err)
		}
	}
}

func clean() {
	_ = os.RemoveAll(exampleDir)
}

func genVidSubFiles() {
	const (
		vidNameFormat = "[VCB-Studio] FAIRY TAIL [%03d][Ma10p_1080p][x265_flac].mkv"
		subNameFormat = "[YYDM-11FANS][FAIRY_TAIL][%03d][BDRIP][720P][X264-10bit_AAC][B721D247].tc.ass"
	)

	_ = os.MkdirAll(vidDir, os.ModePerm)
	_ = os.MkdirAll(originSub, os.ModePerm)
	_ = os.MkdirAll(subDir, os.ModePerm)

	for i := 1; i <= 175; i++ {
		vid, err := os.Create(vidDir + fmt.Sprintf(vidNameFormat, i))
		if err != nil {
			log.Fatalf("failed to generate video file: %v", err)
		}
		_ = vid.Close()

		osub, err := os.Create(originSub + fmt.Sprintf(subNameFormat, i))
		if err != nil {
			log.Fatalf("failed to generate original subDir file: %v", err)
		}
		_ = osub.Close()

		sub, err := os.Create(subDir + fmt.Sprintf(subNameFormat, i))
		if err != nil {
			log.Fatalf("failed to generate subDir file: %v", err)
		}
		_ = sub.Close()
	}
}

func TestAutoRename(t *testing.T) {
	t.Parallel()

	cfg, err := cmd.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	for _, tc := range []struct {
		name    string
		vidDir  string
		subDir  string
		vidExt  []string
		subExt  []string
		wantErr bool
	}{
		{
			name:    "Case 01",
			vidDir:  vidDir,
			subDir:  subDir,
			vidExt:  cfg.VidExt,
			subExt:  cfg.SubExt,
			wantErr: false,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := episode.AutoRename(tc.vidDir, tc.subDir, tc.vidExt, tc.subExt)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
