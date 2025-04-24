package quiz

import (
	"english_app_for_japanese/wasm/objects"
)

// Quiz はクイズモードのデータと状態を管理する構造体です。
type Quiz struct {
	appData         *objects.AppData // アプリケーション全体のデータへのポインタ
	FilteredArray   []objects.Datum  // フィルタリングおよびシャッフルされた問題データのスライス
	index           int              // FilteredArray 内の現在の問題インデックス
	Level           int              // 現在選択されている問題のレベル (0 は全レベル)
	numberOfOptions int              // 各問題で表示する選択肢の数
	CorrectAnswer   *objects.Datum   // 現在の問題の正解データへのポインタ
	OptionsArray    []objects.Datum  // 現在の問題の選択肢（正解を含む）のスライス
}

// Init は Quiz 構造体を初期化します。
// 指定されたレベルに基づいて、アプリケーションデータから未学習の問題をフィルタリングし、
// シャッフルして内部の FilteredArray に格納します。また、選択肢の数を設定します。
//
// 引数:
//   - appData: アプリケーション全体のデータ (objects.AppData) へのポインタ。
//   - level: フィルタリングする問題のレベル。0 を指定するとレベルに関係なくフィルタリングします。
//   - choiceCount: 各問題で生成する選択肢の数（正解を含む）。
func (q *Quiz) Init(appData *objects.AppData, level int, choiceCount int) {
	q.appData = appData
	q.Level = level
	q.numberOfOptions = choiceCount
	q.index = 0 // インデックスを初期化
	// LocalStorageに含まれていない（未学習の）データを取得
	tmp := q.appData.FilterNotInStorage()
	// level が 0 以外の場合、指定されたレベルでさらにフィルタリング
	if q.Level != 0 {
		tmp = objects.FilterByLevel(tmp, q.Level)
	}
	// フィルタリングされたデータをシャッフルして格納
	q.FilteredArray = objects.ShuffleCopy(tmp)
}

// Next は次のクイズ問題に進みます。
// FilteredArray から現在のインデックスに対応する問題データを CorrectAnswer に設定し、
// インデックスを次に進めます。配列の末尾に達した場合は、インデックスを 0 に戻してループさせます。
// 最後に、新しい正解に対応する選択肢を生成するために CreateOptionsArray を呼び出します。
func (q *Quiz) Next() {
	// FilteredArray が空でないことを確認（Init が呼ばれている前提）
	if len(q.FilteredArray) == 0 {
		q.CorrectAnswer = nil
		q.OptionsArray = nil
		return
	}
	// 現在のインデックスのデータを CorrectAnswer に設定
	q.CorrectAnswer = &q.FilteredArray[q.index]
	// インデックスを次に進める
	q.index++
	// インデックスが配列の範囲を超えたら 0 に戻す
	if q.index >= len(q.FilteredArray) {
		q.index = 0
	}
	// 新しい正解に対する選択肢を生成する
	q.CreateOptionsArray()
}

// CreateOptionsArray は現在の正解 (CorrectAnswer) に対する選択肢の配列 (OptionsArray) を生成します。
// 正解データを含め、指定された numberOfOptions の数だけ、重複しないようにランダムな選択肢を
// アプリケーションデータ全体 (appData.Data) から選び出します。
// 生成された選択肢の配列は最後にシャッフルされます。
//
// 注意: appData.Data の要素数が numberOfOptions より少ない場合、
//
//	またはランダム選択の試行回数が上限に達した場合、
//	生成される選択肢の数が numberOfOptions より少なくなる可能性があります。
func (q *Quiz) CreateOptionsArray() {
	// 正解データが設定されていない場合は何もしない
	if q.CorrectAnswer == nil {
		q.OptionsArray = nil
		return
	}

	// 選択肢配列を正解データで初期化
	q.OptionsArray = []objects.Datum{*q.CorrectAnswer}
	// 選択済みのIDを記録するマップ
	selectedIDs := make(map[int]bool)
	selectedIDs[q.CorrectAnswer.ID] = true

	// 無限ループを防ぐための最大試行回数
	maxAttempts := len(q.appData.Data) * 2 // データ数の2倍を試行回数上限とする
	attempts := 0

	// 必要な選択肢の数に達するまで、または最大試行回数に達するまでループ
	for len(q.OptionsArray) < q.numberOfOptions && attempts < maxAttempts {
		attempts++
		// アプリケーションデータ全体からランダムに候補を選択
		candidate, err := objects.GetRandomElement(q.appData.Data)
		// エラーが発生した場合（データが空など）はループを抜ける
		if err != nil {
			break // もしくはエラーハンドリング
		}
		// 候補がまだ選択されていない場合
		if !selectedIDs[candidate.ID] {
			// 選択肢配列に追加し、選択済みIDマップにも記録
			q.OptionsArray = append(q.OptionsArray, candidate)
			selectedIDs[candidate.ID] = true
		}
	}
	// 最終的な選択肢配列をシャッフル
	q.OptionsArray = objects.ShuffleCopy(q.OptionsArray)
}
