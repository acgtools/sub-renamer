package episode_test

import (
	"fmt"
	"log"
	"os"
	"testing"
)

const (
	testDir    = "../../test/"
	exampleDir = "../../test/example/"
)

func TestMain(m *testing.M) {
	genVidSubFiles()
	code := m.Run()
	os.Exit(code)
}

func genVidSubFiles() {
	const (
		vidNameFormat = "[VCB-Studio] FAIRY TAIL [%03d][Ma10p_1080p][x265_flac].mkv"
		subNameFormat = "[YYDM-11FANS][FAIRY_TAIL][%03d][BDRIP][720P][X264-10bit_AAC][B721D247].tc.ass"
	)

	originSub := exampleDir + "origin-sub/"
	vidDir := exampleDir + "vid/"
	sub := exampleDir + "sub/"

	for i := 1; i <= 175; i++ {
		_, err := os.Create(vidDir + fmt.Sprintf(vidNameFormat, i))
		if err != nil {
			log.Fatal("failed to generate video file")
		}

		_, err = os.Create(originSub + fmt.Sprintf(subNameFormat, i))
		if err != nil {
			log.Fatal("failed to generate original sub file")
		}

		_, err = os.Create(sub + fmt.Sprintf(subNameFormat, i))
		if err != nil {
			log.Fatal("failed to generate sub file")
		}
	}
}
