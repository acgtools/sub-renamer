package episode

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/acgtools/sub-renamer/pkg/util"
)

const (
	minFileNum = 2
)

func AutoRename(vidDir, subDir string, cpy bool) error {
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

	slog.Info("Video path", "path", vidDir)
	slog.Info("Subtitle path", "path", subDir)

	slog.Info("Getting episode info...")
	vidMap, err := parseEpisodes(vidDir)
	if err != nil {
		return fmt.Errorf("failed to parse video episode: %w", err)
	}

	subMap, err := parseEpisodes(subDir)
	if err != nil {
		return fmt.Errorf("failed to parse subtitle episode: %w", err)
	}

	slog.Info("Renaming...")
	for ep, vidName := range vidMap {
		subName, ok := subMap[ep]
		if !ok {
			continue
		}

		subExt := filepath.Ext(subName)
		newSubName := strings.TrimSuffix(vidName, filepath.Ext(vidName)) + subExt
		oldSubPath, newSubPath := filepath.Join(subDir, subName), filepath.Join(subDir, newSubName)

		slog.Debug("Rename subtitles", "old_path", oldSubPath, "new_path", newSubPath)

		err = os.Rename(oldSubPath, newSubPath)
		if err != nil {
			return fmt.Errorf("failed to rename subtitle file: %w", err)
		}
	}

	if cpy {
		slog.Info("Copying subs...")

		if err := copySubs(vidDir, subDir); err != nil {
			return fmt.Errorf("copy subs: %w", err)
		}
	}

	slog.Info("Success!")
	return nil
}

func parseEpisodes(dir string) (map[int]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %q, %w", dir, err)
	}

	if len(entries) < minFileNum {
		return nil, errors.New("number of files must be greater than 2")
	}

	filteredEntries := filterFiles(entries)
	if len(filteredEntries) < minFileNum {
		return nil, errors.New("number of files must be greater than 2")
	}

	epStartIndex, err := getEpPosInName(filteredEntries[0].Name(), filteredEntries[1].Name())
	if err != nil {
		return nil, fmt.Errorf("failed to get episode position in file name: %q, %w", filteredEntries[0].Name(), err)
	}

	nameEpMap := make(map[int]string, len(filteredEntries))
	for _, entry := range filteredEntries {
		fileName := entry.Name()
		if epStartIndex > len(fileName)-1 {
			slog.Warn("Cannot get episode number, skip file", "file_name", fileName)
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

	for i := 0; i < len(numMatchIndex1); i++ {
		subMatch1, subMatch2 := numMatchIndex1[i], numMatchIndex2[i]
		num1 := fileName1[subMatch1[0]:subMatch1[1]]
		num2 := fileName2[subMatch2[0]:subMatch2[1]]

		if num1 != num2 {
			return subMatch1[0], nil
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

func filterFiles(entries []os.DirEntry) []os.DirEntry {
	filteredEntries := util.SliceFilter(entries, func(e os.DirEntry) bool {
		return !e.IsDir()
	})

	return filteredEntries
}

func copyFile(dstPath, srcPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open file %q: %w", srcPath, err)
	}
	defer fClose(src)

	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("create file: %q: %w", dstPath, err)
	}
	defer fClose(dst)

	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("copy %q to %q: %w", srcPath, dstPath, err)
	}

	return nil
}

func fClose(f *os.File) {
	_ = f.Close()
}

func copySubs(vidDir, subDir string) error {
	entries, err := os.ReadDir(subDir)
	if err != nil {
		return fmt.Errorf("read directory %q: %w", subDir, err)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()
		err := copyFile(filepath.Join(vidDir, name), filepath.Join(subDir, name))
		if err != nil {
			return err
		}
	}

	return nil
}
