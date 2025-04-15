package objects

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
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

// ShuffleObjects は与えられた Object スライスのコピーをランダムにシャッフルして返します。
// 元のスライスは変更されません。
// 注意: WASM環境での time.Now() の挙動によっては、毎回同じシードになる可能性があります。
// より信頼性の高いランダム性が必要な場合は、JavaScript側からシード値を提供することを検討してください。
func ShuffleObjects(inputObjects []Object) []Object {
	if len(inputObjects) <= 1 {
		shuffled := make([]Object, len(inputObjects))
		copy(shuffled, inputObjects)
		return shuffled
	}

	shuffled := make([]Object, len(inputObjects))
	copy(shuffled, inputObjects)

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled
}

// FilterObjectsNotInStorage は、与えられた Object スライスから、
// グローバル変数 Storage スライスに含まれていない ID を持つ要素を抽出して
// 新しいスライスとして返します。
func FilterObjectsNotInStorage(inputObjects []Object) []Object {
	// Storage の ID を効率的に検索するためのセットを作成
	storageSet := make(map[int]struct{}, len(Storage))
	for _, id := range Storage {
		storageSet[id] = struct{}{}
	}

	// 結果を格納するためのスライス
	filteredObjects := make([]Object, 0)

	// 与えられた inputObjects スライスを走査し、Storage に含まれていないものを抽出
	for _, obj := range inputObjects { // ループ対象を inputObjects に変更
		// obj.ID が storageSet に存在するかどうかを確認
		if _, found := storageSet[obj.ID]; !found {
			// 存在しない場合（Storageに含まれていない場合）、結果のスライスに追加
			filteredObjects = append(filteredObjects, obj)
		}
	}

	// 抽出された Object のスライスを返す
	return filteredObjects
}

// FilterObjectsByLevel は、与えられた Object スライスから指定された level に一致する
// Level を持つ要素を抽出して新しいスライスとして返します。
func FilterObjectsByLevel(inputObjects []Object, level int) []Object {
	filteredObjects := make([]Object, 0)

	for _, obj := range inputObjects {
		if obj.Level == level {
			filteredObjects = append(filteredObjects, obj)
		}
	}

	return filteredObjects
}

// SelectRandomObject は与えられた Object スライスからランダムに1つの Object を選択して返します。
// スライスが空の場合はゼロ値の Object (全てのフィールドがデフォルト値) を返します。
// 注意: WASM環境での time.Now() の挙動によっては、毎回同じシードになる可能性があります。
// より信頼性の高いランダム性が必要な場合は、外部からシード値を提供することを検討してください。
func SelectRandomObject(inputObjects []Object) Object {
	// スライスの長さを確認
	n := len(inputObjects)

	// スライスが空の場合の処理
	if n == 0 {
		// 空のスライスが渡された場合は、ゼロ値のObjectを返す
		return Object{}
	}

	// ランダムなインデックスを生成するための準備
	// ShuffleObjects と同様に、WASM 環境での time.Now() の信頼性に注意
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	// 0 から n-1 の範囲でランダムなインデックスを生成
	randomIndex := r.Intn(n)

	// ランダムに選ばれたインデックスの Object を返す
	return inputObjects[randomIndex]
}

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

// GetObjectsForQuiz は出題を作成する
func GetObjectsForQuiz(level int, index int) {

}
