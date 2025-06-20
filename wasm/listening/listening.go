package listening

import "english_app_for_japanese/wasm/objects"

// Listening はリスニング学習モードのデータと状態を管理する構造体です。
type Listening struct {
	appData       *objects.AppData // アプリケーション全体のデータへのポインタ
	FilteredArray []objects.Datum  // フィルタリングおよびシャッフルされた問題データのスライス
}

// Init は Listening 構造体を初期化します。
// 指定されたレベルに基づいて、アプリケーションデータから未学習の問題をフィルタリングし、
// シャッフルして内部の FilteredArray に格納します。
//
// 引数:
//   - appData: アプリケーション全体のデータ (objects.AppData) へのポインタ。
//   - count: 取得する問題の最大数。
func (l *Listening) Init(appData *objects.AppData, count int) {
	l.appData = appData
	// LocalStorageに含まれていない（未学習の）データを取得
	tmp := l.appData.FilterNotInStorage()
	// フィルタリングされたデータをシャッフルして格納
	l.FilteredArray = objects.ShuffleCopy(tmp)
	// FilteredArrayの要素数を指定されたcountに制限
	if count > 0 && len(l.FilteredArray) > count {
		l.FilteredArray = l.FilteredArray[:count]
	}
}
