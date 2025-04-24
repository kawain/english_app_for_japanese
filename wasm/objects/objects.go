package objects

import (
	"errors"
	"math/rand/v2"
	"strconv"
	"strings"
)

// Datum は単語とその関連情報を保持する構造体です。
type Datum struct {
	ID      int    // 単語の一意なID
	En      string // 英語の単語または短いフレーズ
	Jp      string // Enに対応する日本語訳
	En2     string // 英語の例文または長いフレーズ
	Jp2     string // En2に対応する日本語訳
	Kana    string // Jp2の読み仮名（ひらがな）
	Level   int    // 単語の難易度レベル (例: 1, 2)
	Similar []int  // 類似または関連する単語のIDのスライス
}

// NewDatum は文字列形式のデータから新しい Datum オブジェクトを生成します。
// 各文字列引数はトリムされ、idText と levelText は整数に、
// similarText はカンマ区切りの文字列から整数のスライスに変換されます。
// 変換に失敗した場合、対応するフィールドはゼロ値（数値の場合は0、スライスの場合は空）になります。
//
// 引数:
//   - idText: 単語IDを表す文字列。
//   - en: 英語の単語または短いフレーズ。
//   - jp: enに対応する日本語訳。
//   - en2: 英語の例文または長いフレーズ。
//   - jp2: en2に対応する日本語訳。
//   - kana: jp2の読み仮名（ひらがな）。
//   - levelText: 単語の難易度レベルを表す文字列。
//   - similarText: 類似または関連する単語のIDをカンマ区切りで表す文字列。
//
// 戻り値:
//   - 初期化された Datum オブジェクト。
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
