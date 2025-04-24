package listening

import "english_app_for_japanese/wasm/objects"

// Listening はリスニング学習モードのデータと状態を管理する構造体です。
type Listening struct {
	appData       *objects.AppData // アプリケーション全体のデータへのポインタ
	FilteredArray []objects.Datum  // フィルタリングおよびシャッフルされた問題データのスライス
	index         int              // FilteredArray 内の現在の問題インデックス
	Level         int              // 現在選択されている問題のレベル (0 は全レベル)
	CurrentData   *objects.Datum   // 現在表示または再生中の問題データへのポインタ
}

// Init は Listening 構造体を初期化します。
// 指定されたレベルに基づいて、アプリケーションデータから未学習の問題をフィルタリングし、
// シャッフルして内部の FilteredArray に格納します。
//
// 引数:
//   - appData: アプリケーション全体のデータ (objects.AppData) へのポインタ。
//   - level: フィルタリングする問題のレベル。0 を指定するとレベルに関係なくフィルタリングします。
func (l *Listening) Init(appData *objects.AppData, level int) {
	l.appData = appData
	l.Level = level
	// LocalStorageに含まれていない（未学習の）データを取得
	tmp := l.appData.FilterNotInStorage()
	// level が 0 以外の場合、指定されたレベルでさらにフィルタリング
	if l.Level != 0 {
		tmp = objects.FilterByLevel(tmp, l.Level)
	}
	// フィルタリングされたデータをシャッフルして格納
	l.FilteredArray = objects.ShuffleCopy(tmp)
	// インデックスを初期化
	l.index = 0
}

// Next は次のリスニング問題に進みます。
// FilteredArray から現在のインデックスに対応する問題データを CurrentData に設定し、
// インデックスを次に進めます。配列の末尾に達した場合は、インデックスを 0 に戻してループさせます。
func (l *Listening) Next() {
	// FilteredArray が空でないことを確認（Init が呼ばれている前提）
	if len(l.FilteredArray) == 0 {
		l.CurrentData = nil // データがない場合は nil を設定
		return
	}
	// 現在のインデックスのデータを CurrentData に設定
	l.CurrentData = &l.FilteredArray[l.index]
	// インデックスを次に進める
	l.index++
	// インデックスが配列の範囲を超えたら 0 に戻す
	if l.index >= len(l.FilteredArray) {
		l.index = 0
	}
}
