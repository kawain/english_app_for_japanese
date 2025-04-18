package objects

import (
	"errors"
	"math/rand/v2"
	"strconv"
	"strings"
)

type Datum struct {
	ID      int
	En      string
	Jp      string
	En2     string
	Jp2     string
	Kana    string
	Level   int
	Similar []int
}

func NewDatum(idText, en, jp, en2, jp2, kana, levelText, similarText string) Datum {
	idText = strings.TrimSpace(idText)
	en = strings.TrimSpace(en)
	jp = strings.TrimSpace(jp)
	en2 = strings.TrimSpace(en2)
	jp2 = strings.TrimSpace(jp2)
	kana = strings.TrimSpace(kana)
	levelText = strings.TrimSpace(levelText)
	similarText = strings.TrimSpace(similarText)

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

	return Datum{
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

type AppData struct {
	Data         []Datum
	LocalStorage []int
}

func (a *AppData) AddData(datum Datum) {
	a.Data = append(a.Data, datum)
}

func (a *AppData) AddStorage(id int) {
	a.LocalStorage = append(a.LocalStorage, id)
}

func (a *AppData) RemoveStorage(idToRemove int) {
	newStorage := make([]int, 0, len(a.LocalStorage))
	for _, id := range a.LocalStorage {
		if id != idToRemove {
			newStorage = append(newStorage, id)
		}
	}
	a.LocalStorage = newStorage
}

func (a *AppData) ClearStorage() {
	a.LocalStorage = make([]int, 0)
}

// FilterNotInStorage は、LocalStorageに含まれていない ID を持つ要素を抽出して
// 新しいスライスとして返します
func (a *AppData) FilterNotInStorage() []Datum {
	storageSet := make(map[int]bool, len(a.LocalStorage))
	for _, id := range a.LocalStorage {
		storageSet[id] = true
	}
	results := make([]Datum, 0)
	for _, obj := range a.Data {
		if !storageSet[obj.ID] {
			results = append(results, obj)
		}
	}
	return results
}

// FilterInStorage は、LocalStorageに含まれている ID を持つ要素を抽出して
// 新しいスライスとして返します
func (a *AppData) FilterInStorage() []Datum {
	storageSet := make(map[int]bool, len(a.LocalStorage))
	for _, id := range a.LocalStorage {
		storageSet[id] = true
	}
	results := make([]Datum, 0)
	for _, obj := range a.Data {
		if storageSet[obj.ID] {
			results = append(results, obj)
		}
	}
	return results
}

// ShuffleCopy は元のスライスを変更せず、シャッフルされた新しいスライスを返します。
// ジェネリクスを使用して任意の型のスライスに対応します。
func ShuffleCopy[T any](original []T) []T {
	n := len(original)
	// 要素数が1以下の場合はシャッフルの必要がない（コピーだけ行う）
	if n <= 1 {
		newSlice := make([]T, n)
		copy(newSlice, original)
		return newSlice
	}

	// 1. 元のスライスのコピーを作成
	shuffled := make([]T, n)
	copy(shuffled, original) // original の内容を shuffled にコピー

	// 2. コピーしたスライス (shuffled) をシャッフル
	// Go 1.22 以降 (math/rand/v2)
	rand.Shuffle(n, func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	// 3. シャッフルされた新しいスライスを返す
	return shuffled
}

// スライスからランダムに1つの要素を返す関数
func GetRandomElement[T any](slice []T) (T, error) {
	n := len(slice)
	if n == 0 {
		var zero T // 型に応じたゼロ値を返すため
		return zero, errors.New("cannot get random element from an empty slice")
	}

	// 0 から n-1 の範囲でランダムなインデックスを取得
	// Go 1.22 以降 (math/rand/v2)
	randomIndex := rand.IntN(n) // 0 <= randomIndex < n

	return slice[randomIndex], nil
}

// FilterByLevel は与えられた配列からlevelで抽出した結果を返す
func FilterByLevel(data []Datum, level int) []Datum {
	results := make([]Datum, 0)
	for _, obj := range data {
		if obj.Level == level {
			results = append(results, obj)
		}
	}
	return results
}
