package episode_test

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/acgtools/sub-renamer/pkg/episode"
	"github.com/stretchr/testify/require"
)

const (
	rootDir = "../../"

	testDir    = rootDir + "test/"
	exampleDir = testDir + "example/"

	originSub = "origin-subDir"
	vidDir    = "vid"
	subDir    = "sub"
)

func TestMain(m *testing.M) {
	initLogger()
	clean()

	code := m.Run()
	os.Exit(code)
}

func initLogger() {
	logLevel := slog.LevelDebug
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)
}

func clean() {
	_ = os.RemoveAll(exampleDir)
}

func genVidSubFiles(caseName string) {
	const (
		vidNameFormat = "[VCB-Studio] FAIRY TAIL [%03d][Ma10p_1080p][x265_flac].mkv"
		subNameFormat = "[YYDM-11FANS][FAIRY_TAIL][%03d][BDRIP][720P][X264-10bit_AAC][B721D247].tc.ass"
	)

	vidPath := filepath.Join(exampleDir, caseName, vidDir)
	originSubPath := filepath.Join(exampleDir, caseName, originSub)
	subPath := filepath.Join(exampleDir, caseName, subDir)

	_ = os.MkdirAll(vidPath, os.ModePerm)
	_ = os.MkdirAll(originSubPath, os.ModePerm)
	_ = os.MkdirAll(subPath, os.ModePerm)

	for i := 1; i <= 175; i++ {
		vid, err := os.Create(filepath.Join(vidPath, fmt.Sprintf(vidNameFormat, i)))
		if err != nil {
			log.Fatalf("failed to generate video file: %v", err)
		}
		_ = vid.Close()

		osub, err := os.Create(filepath.Join(originSubPath, fmt.Sprintf(subNameFormat, i)))
		if err != nil {
			log.Fatalf("failed to generate original subDir file: %v", err)
		}
		_ = osub.Close()

		sub, err := os.Create(filepath.Join(subPath, fmt.Sprintf(subNameFormat, i)))
		if err != nil {
			log.Fatalf("failed to generate subDir file: %v", err)
		}
		_ = sub.Close()
	}

	// some other dir or files
	dir1, dir2 := "tmp1/", "tmp2/"
	_ = os.Mkdir(filepath.Join(vidPath, dir1), os.ModePerm)
	_ = os.Mkdir(filepath.Join(vidPath, dir2), os.ModePerm)
	_ = os.Mkdir(filepath.Join(subPath, dir1), os.ModePerm)
	_ = os.Mkdir(filepath.Join(subPath, dir2), os.ModePerm)

	tf1, err := os.Create(filepath.Join(vidPath, dir1, "tmp.txt"))
	if err != nil {
		log.Fatalf("failed to generate other file: %v", err)
	}
	_ = tf1.Close()

	tf2, err := os.Create(filepath.Join(subPath, dir1, "tmp.txt"))
	if err != nil {
		log.Fatalf("failed to generate other file: %v", err)
	}
	_ = tf2.Close()

	file1, file2 := "file1.mkv", "file2.txt"
	f1, err := os.Create(filepath.Join(vidPath, file1))
	if err != nil {
		log.Fatalf("failed to generate other file: %v", err)
	}
	_ = f1.Close()

	f2, err := os.Create(filepath.Join(vidPath, file2))
	if err != nil {
		log.Fatalf("failed to generate other file: %v", err)
	}
	_ = f2.Close()

	file3, file4 := "file3.ass", "file4.txt"
	f3, err := os.Create(filepath.Join(subPath, file3))
	if err != nil {
		log.Fatalf("failed to generate other file: %v", err)
	}
	_ = f3.Close()

	f4, err := os.Create(filepath.Join(subPath, file4))
	if err != nil {
		log.Fatalf("failed to generate other file: %v", err)
	}
	_ = f4.Close()
}

func TestAutoRename(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name    string
		vidDir  string
		subDir  string
		cpy     bool
		wantErr bool
	}{
		{
			name:    "Not copy subs",
			cpy:     false,
			wantErr: false,
		},
		{
			name:    "Copy subs",
			cpy:     true,
			wantErr: false,
		},
	} {
		tc := tc
		tc.vidDir = filepath.Join(exampleDir, tc.name, vidDir)
		tc.subDir = filepath.Join(exampleDir, tc.name, subDir)

		genVidSubFiles(tc.name)

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := episode.AutoRename(tc.vidDir, tc.subDir, tc.cpy)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
