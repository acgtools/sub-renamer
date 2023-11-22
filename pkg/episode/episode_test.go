package episode

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_getEpPosInName1(t *testing.T) {
	for _, tc := range []struct {
		name          string
		fileName1     string
		fileName2     string
		epStartIndex  int
		hasErr        bool
		diffPattenErr bool
		notFoundErr   bool
	}{
		{
			name:          "Both empty string",
			fileName1:     "",
			fileName2:     "",
			epStartIndex:  -1,
			hasErr:        true,
			diffPattenErr: false,
			notFoundErr:   true,
		},
		{
			name:          "FileName1 is empty string",
			fileName1:     "",
			fileName2:     "[VCB-Studio] FAIRY TAIL [12][Ma10p_1080p][x265_flac].mkv",
			epStartIndex:  -1,
			hasErr:        true,
			diffPattenErr: true,
			notFoundErr:   false,
		},
		{
			name:          "FileName2 is empty string",
			fileName1:     "[VCB-Studio] FAIRY TAIL [12][Ma10p_1080p][x265_flac].mkv",
			fileName2:     "",
			epStartIndex:  -1,
			hasErr:        true,
			diffPattenErr: true,
			notFoundErr:   false,
		},
		{
			name:          "FileNames in different pattern",
			fileName1:     "[VCB-Studio] FAIRY TAIL [12][Ma10p_1080p].mkv",
			fileName2:     "[VCB-Studio] FAIRY TAIL [12][Ma10p_1080p][x265_flac].mkv",
			epStartIndex:  -1,
			hasErr:        true,
			diffPattenErr: true,
			notFoundErr:   false,
		},
		{
			name:         "FileNames in same pattern",
			fileName1:    "[VCB-Studio] FAIRY TAIL [12][Ma10p_1080p][x265_flac].mkv",
			fileName2:    "[VCB-Studio] FAIRY TAIL [13][Ma10p_1080p][x265_flac].mkv",
			epStartIndex: 25,
			hasErr:       false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			idx, err := getEpPosInName(tc.fileName1, tc.fileName2)

			assert.Equal(t, tc.epStartIndex, idx)

			if tc.hasErr {
				assert.NotNil(t, err)
				switch {
				case tc.diffPattenErr:
					assert.ErrorContains(t, err, "file names are not in same pattern")
				case tc.notFoundErr:
					assert.ErrorContains(t, err, "episode number not found")
				}
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestParseEpisodes(t *testing.T) {
	for _, tc := range []struct {
		name      string
		dir       string
		episodes  map[int]string
		wantErr   bool
		customErr bool
		err       error
	}{
		{
			name:     "Empty string",
			dir:      "",
			episodes: nil,
			wantErr:  true,
			err:      os.ErrNotExist,
		},
		{
			name:     "Directory not exists",
			dir:      "not exists path",
			episodes: nil,
			wantErr:  true,
			err:      os.ErrNotExist,
		},
		{
			name:      "Number of files less than 2",
			dir:       "../../test-data/parse-episodes/case01",
			episodes:  nil,
			wantErr:   true,
			customErr: true,
			err:       errors.New("number of file must be greater than 2"),
		},
		{
			name:      "Files have different pattern",
			dir:       "../../test-data/parse-episodes/case02",
			episodes:  nil,
			wantErr:   true,
			customErr: true,
			err:       errors.New("file names are not in same pattern"),
		},
		{
			name: "Parse episodes",
			dir:  "../../test-data/parse-episodes/case03",
			episodes: map[int]string{
				12:  "[VCB-Studio] FAIRY TAIL [12][Ma10p_1080p][x265_flac].mkv",
				13:  "[VCB-Studio] FAIRY TAIL [13][Ma10p_1080p][x265_flac].mkv",
				169: "[VCB-Studio] FAIRY TAIL [169][Ma10p_1080p][x265_flac].mkv",
			},
			wantErr: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			eps, err := ParseEpisodes(tc.dir)

			assert.Equal(t, tc.episodes, eps)

			if tc.wantErr {
				if !tc.customErr {
					assert.ErrorIs(t, err, tc.err)
				} else {
					assert.ErrorContains(t, err, tc.err.Error())
				}
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
