package episode

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/dreamjz/sub-renamer/pkg/util"
)

const (
	minFileNum    = 2
	minFileExtLen = 2
)

func AutoRename(vidDir, subDir string, vidExt, subExt []string) error {
	var err error

	if !filepath.IsAbs(vidDir) {
		vidDir, err = filepath.Abs(vidDir)
		if err != nil {
			return fmt.Errorf("failed to convert video path %q to absolute path: %w", vidDir, err)
		}
	}
	if !filepath.IsAbs(subDir) {
		subDir, err = filepath.Abs(subDir)
		if err != nil {
			return fmt.Errorf("failed to convert subtitle path %q to absolute path: %w", subDir, err)
		}
	}

	vidExtSet := util.SliceToSet(vidExt)
	vidMap, err := parseEpisodes(vidDir, vidExtSet)
	if err != nil {
		return fmt.Errorf("failed to parse video episode: %w", err)
	}

	subExtSet := util.SliceToSet(subExt)
	subMap, err := parseEpisodes(subDir, subExtSet)
	if err != nil {
		return fmt.Errorf("failed to parse subtitle episode: %w", err)
	}

	for ep, vidName := range vidMap {
		subName, ok := subMap[ep]
		if !ok {
			continue
		}

		subExt := filepath.Ext(subName)

		newSubName := strings.TrimSuffix(vidName, filepath.Ext(vidName))

		err = os.Rename(filepath.Join(subDir, subName), filepath.Join(subDir, newSubName+subExt))
		if err != nil {
			return fmt.Errorf("failed to rename subtitle file: %w", err)
		}
	}

	return nil
}

func parseEpisodes(dir string, supportedExt map[string]struct{}) (map[int]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %q, %w", dir, err)
	}

	if len(entries) < minFileNum {
		return nil, errors.New("number of file must be greater than 2")
	}

	filteredEntries, err := filterFiles(entries, supportedExt)
	if err != nil {
		return nil, err
	}
	if len(filteredEntries) < minFileNum {
		return nil, errors.New("number of file must be greater than 2")
	}

	epStartIndex, err := getEpPosInName(filteredEntries[0].Name(), filteredEntries[1].Name())
	if err != nil {
		return nil, fmt.Errorf("failed to get episode position in file name: %q, %w", filteredEntries[0].Name(), err)
	}

	nameEpMap := make(map[int]string, len(filteredEntries))
	for _, entry := range filteredEntries {
		fileName := entry.Name()
		if epStartIndex > len(fileName)-1 {
			slog.Warn(fmt.Sprintf("cannot get episode number, skip: %q", fileName))
			continue
		}
		epNum := getEpisodeNum(fileName, epStartIndex)
		nameEpMap[epNum] = fileName
	}

	return nameEpMap, nil
}

func getEpPosInName(fileName1, fileName2 string) (int, error) {
	r := regexp.MustCompile(`\d+`)
	numMatchIndex1 := r.FindAllStringIndex(fileName1, -1)
	numMatchIndex2 := r.FindAllStringIndex(fileName2, -1)

	if len(numMatchIndex1) != len(numMatchIndex2) {
		return -1, errors.New("file names are not in same pattern")
	}

	for _, subMatch := range numMatchIndex1 {
		num1 := fileName1[subMatch[0]:subMatch[1]]
		num2 := fileName2[subMatch[0]:subMatch[1]]

		if num1 != num2 {
			return subMatch[0], nil
		}
	}

	return -1, errors.New("episode number not found")
}

func getEpisodeNum(fileName string, start int) int {
	end := start + 1
	for end <= len(fileName)-1 && isDigit(fileName[end]) {
		end++
	}

	ep, _ := strconv.Atoi(fileName[start:end])
	return ep
}

func isDigit(b byte) bool {
	if b >= '0' && b <= '9' {
		return true
	}
	return false
}

func filterFiles(entries []os.DirEntry, supportedExt map[string]struct{}) ([]os.DirEntry, error) {
	ext, err := determineFileExt(entries, supportedExt)
	if err != nil {
		return nil, err
	}

	filteredEntries := util.SliceFilter(entries, func(e os.DirEntry) bool {
		if e.IsDir() {
			return false
		}
		return fileExtTrimDot(e.Name()) == ext
	})

	return filteredEntries, nil
}

func determineFileExt(entries []os.DirEntry, supportedExt map[string]struct{}) (string, error) {
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		fileExt := fileExtTrimDot(e.Name())

		if _, ok := supportedExt[fileExt]; ok {
			return fileExt, nil
		}
	}

	return "", errors.New("no supported files found")
}

func fileExtTrimDot(name string) string {
	ext := filepath.Ext(name)
	if len(ext) < minFileExtLen {
		return ""
	}

	return ext[1:]
}
