package episode

import (
	"errors"
	"os"
	"regexp"
	"strconv"
)

var ()

func ParseEpisodes(dir string) (map[int]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	if len(entries) < 2 {
		return nil, errors.New("number of file must be greater than 2")
	}

	epStartIndex, err := getEpPosInName(entries[0].Name(), entries[1].Name())
	if err != nil {
		return nil, err
	}

	nameEpMap := make(map[int]string, len(entries))
	for _, entry := range entries {
		fileName := entry.Name()
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
	for isDigit(fileName[end]) {
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
