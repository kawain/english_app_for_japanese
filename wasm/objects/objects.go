package objects

import (
	"strconv"
	"strings"
)

type Object struct {
	ID      int
	En      string
	Jp      string
	En2     string
	Jp2     string
	Kana    string
	Level   int
	Similar []int
}

var Objects []Object
var Storage []int

func NewObjects(idText, en, jp, en2, jp2, kana, levelText, similarText string) Object {
	id, _ := strconv.Atoi(idText)
	level, _ := strconv.Atoi(levelText)
	similarSlice := strings.Split(similarText, ",")
	// 結果を格納するためのintスライスを準備
	similar := make([]int, 0, len(similarSlice))

	for _, strNum := range similarSlice {
		trimmedStr := strings.TrimSpace(strNum)
		// 空文字列になった場合はスキップする（例: "1,,2" のような入力）
		if trimmedStr == "" {
			continue
		}
		num, _ := strconv.Atoi(trimmedStr)
		similar = append(similar, num)
	}

	return Object{
		ID:      id,
		En:      en,
		Jp:      jp,
		En2:     en2,
		Jp2:     jp2,
		Kana:    kana,
		Level:   level,
		Similar: similar,
	}
}
