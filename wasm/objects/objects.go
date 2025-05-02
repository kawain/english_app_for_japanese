package objects

import (
	"errors"
	"math/rand/v2"
	"strconv"
	"strings"
)

// Datum は単語とその関連情報を保持する構造体です。
type Datum struct {
	ID           int
	Word         string
	DefinitionEn string
	DefinitionJa string
	ExampleEn    string
	ExampleJa    string
	Kana         string
	Level        int
	SimilarIDs   []int
}

// NewDatum は文字列形式のデータから新しい Datum オブジェクトを生成します。
func NewDatum(idText, word, definitionEn, definitionJa, exampleEn, exampleJa, kana, levelText, similarText string) Datum {
	idText = strings.TrimSpace(idText)
	word = strings.TrimSpace(word)
	definitionEn = strings.TrimSpace(definitionEn)
	definitionJa = strings.TrimSpace(definitionJa)
	exampleEn = strings.TrimSpace(exampleEn)
	exampleJa = strings.TrimSpace(exampleJa)
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
		ID:           id,
		Word:         word,
		DefinitionEn: definitionEn,
		DefinitionJa: definitionJa,
		ExampleEn:    exampleEn,
		ExampleJa:    exampleJa,
		Kana:         kana,
		Level:        level,
		SimilarIDs:   similar,
	}
}

// AppData はアプリケーション全体のデータ（単語データとローカルストレージ情報）を保持します。
type AppData struct {
	Data         []Datum // すべての単語データのスライス
	LocalStorage []int   // ローカルストレージに保存されている（学習済みなどの）単語IDのスライス (重複なし)
}

// AddData は AppData の Data スライスに新しい Datum を追加します。
//
// 引数:
//   - datum: 追加する Datum オブジェクト。
func (a *AppData) AddData(datum Datum) {
	a.Data = append(a.Data, datum)
}

// AddStorage は AppData の LocalStorage スライスに新しい単語IDを追加します。
// すでにIDが存在する場合は、何も行いません（重複を防ぐ）。
//
// 引数:
//   - id: 追加する単語ID。
func (a *AppData) AddStorage(id int) {
	// 既存のIDをチェック
	for _, existingID := range a.LocalStorage {
		if existingID == id {
			return // すでに存在するので何もしない
		}
	}
	// 存在しない場合のみ追加
	a.LocalStorage = append(a.LocalStorage, id)
}

// RemoveStorage は AppData の LocalStorage スライスから指定されたIDを削除します。
// 指定されたIDが存在しない場合、スライスは変更されません。
//
// 引数:
//   - idToRemove: 削除する単語ID。
func (a *AppData) RemoveStorage(idToRemove int) {
	newStorage := make([]int, 0, len(a.LocalStorage))
	for _, id := range a.LocalStorage {
		if id != idToRemove {
			newStorage = append(newStorage, id)
		}
	}
	a.LocalStorage = newStorage
}

// ClearStorage は AppData の LocalStorage スライスを空にします。
func (a *AppData) ClearStorage() {
	a.LocalStorage = make([]int, 0)
}

// FilterNotInStorage は AppData の Data スライスから、
// LocalStorage に ID が含まれて *いない* Datum のみをフィルタリングして新しいスライスとして返します。
// 元の Data スライスは変更されません。
//
// 戻り値:
//   - LocalStorage に含まれていない Datum のスライス。
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

// FilterInStorage は AppData の Data スライスから、
// LocalStorage に ID が含まれて *いる* Datum のみをフィルタリングして新しいスライスとして返します。
// 元の Data スライスは変更されません。
//
// 戻り値:
//   - LocalStorage に含まれている Datum のスライス。
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

// ShuffleCopy は与えられたスライスのコピーを作成し、そのコピーをシャッフルして返します。
// 元のスライスは変更されません。ジェネリック関数です。
// 要素数が1以下の場合は、単にコピーを返します。
//
// 引数:
//   - original: シャッフルしたい元のスライス。
//
// 戻り値:
//   - シャッフルされた新しいスライス。
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

// GetRandomElement は与えられたスライスからランダムな要素を1つ取得して返します。
// ジェネリック関数です。
// スライスが空の場合、型に応じたゼロ値とエラーを返します。
//
// 引数:
//   - slice: 要素を取得したいスライス。
//
// 戻り値:
//   - ランダムに選択された要素。
//   - スライスが空だった場合のエラー。
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

// FilterByLevel は Datum のスライスを受け取り、指定された level に一致する要素のみを
// フィルタリングして新しいスライスとして返します。
// 元のスライスは変更されません。
//
// 引数:
//   - data: フィルタリング対象の Datum スライス。
//   - level: フィルタリング条件となるレベル値。
//
// 戻り値:
//   - 指定された level に一致する Datum のスライス。
func FilterByLevel(data []Datum, level int) []Datum {
	results := make([]Datum, 0)
	for _, obj := range data {
		if obj.Level == level {
			results = append(results, obj)
		}
	}
	return results
}
