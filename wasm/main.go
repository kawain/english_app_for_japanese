//go:build js && wasm

package main

import (
	"english_app_for_japanese/wasm/listening"
	"english_app_for_japanese/wasm/objects"
	"english_app_for_japanese/wasm/quiz"
	"english_app_for_japanese/wasm/typing"
	"fmt"
	"syscall/js"
)

var consoleLog js.Value
var appData objects.AppData
var quizData quiz.Quiz
var typingData typing.Typing
var listeningData listening.Listening

// init はGoプログラムの初期化関数です。
// JavaScriptの `console.log` 関数への参照を取得し、
// アプリケーションデータ、クイズデータ、タイピングデータ、リスニングデータの
// 各構造体を初期化します。
func init() {
	consoleLog = js.Global().Get("console").Get("log")
	appData = objects.AppData{}
	quizData = quiz.Quiz{}
	typingData = typing.Typing{}
	listeningData = listening.Listening{}
}

// SearchData はJavaScriptから呼び出され、指定されたレベルに基づいてデータを検索し、
// シャッフルされた結果をJavaScriptのオブジェクト配列として返します。
//
// 引数:
//   - args[0]: level (数値型)
//   - 0: ローカルストレージ（学習済みなど）に含まれるデータを検索
//   - 1: ローカルストレージに含まれないレベル1のデータを検索
//   - 2: ローカルストレージに含まれないレベル2のデータを検索
//
// 戻り値:
//   - JavaScriptのPromiseオブジェクト。成功時には検索結果のオブジェクト配列、
//     失敗時にはエラーメッセージで解決または拒否されます。
//
// 処理内容:
//  1. appDataが初期化されているか確認します。
//  2. 引数の数と型を検証します。
//  3. 指定されたlevelに基づいてデータをフィルタリングします。
//     - level 0: appData.FilterInStorage() を使用します。
//     - level 1, 2: appData.FilterNotInStorage() と objects.FilterByLevel() を使用します。
//  4. フィルタリングされた結果を objects.ShuffleCopy() でシャッフルします。
//  5. 検索結果の各DatumオブジェクトをJavaScriptで扱いやすい形式 (map[string]interface{}) に変換します。
//  6. 変換されたオブジェクトの配列をPromiseのresolve関数に渡して返します。
//  7. エラーが発生した場合は、Promiseのreject関数にエラーメッセージを渡します。
func SearchData(this js.Value, args []js.Value) any {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]
		go func() {
			if appData.Data == nil {
				reject.Invoke(js.ValueOf("Go関数(SearchData)エラー: appDataが初期化されていません。CreateObjectを先に呼び出してください。"))
				return
			}
			if len(args) != 1 {
				reject.Invoke(js.ValueOf("Go関数(SearchData)エラー: 引数は1つ必要です"))
				return
			}
			if args[0].Type() != js.TypeNumber {
				reject.Invoke(js.ValueOf("Go関数(SearchData)エラー: 引数は数値型である必要があります"))
				return
			}
			level := args[0].Int()
			var results []objects.Datum
			switch level {
			case 0:
				r := appData.FilterInStorage()
				results = objects.ShuffleCopy(r)
			case 1:
				r := appData.FilterNotInStorage()
				r = objects.FilterByLevel(r, 1)
				results = objects.ShuffleCopy(r)
			case 2:
				r := appData.FilterNotInStorage()
				r = objects.FilterByLevel(r, 2)
				results = objects.ShuffleCopy(r)
			default:
				reject.Invoke(js.ValueOf(fmt.Sprintf("Go関数(SearchData)エラー: 無効なlevel値です: %d", level)))
				return
			}
			consoleLog.Invoke(js.ValueOf("Go関数(SearchData)で検索したデータの長さ:"), js.ValueOf(len(results)))
			// --- JavaScriptのデータに変換 ---
			jsResult := make([]interface{}, len(results))
			for i, v := range results {
				obj := map[string]interface{}{
					"id":    v.ID,
					"en":    v.En,
					"jp":    v.Jp,
					"en2":   v.En2,
					"jp2":   v.Jp2,
					"level": v.Level,
				}
				jsResult[i] = obj
			}
			resolve.Invoke(jsResult)
		}()
		return nil
	})
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// main はWASMモジュールのエントリーポイントです。
// Goで実装された各種機能をJavaScriptのグローバルスコープに登録し、
// JavaScript側から呼び出せるようにします。
// 登録後、プログラムが終了しないように `select {}` で待機します。
func main() {
	// データ管理関連の関数を登録
	js.Global().Set("CreateObject", js.FuncOf(CreateObject))
	js.Global().Set("SetStorage", js.FuncOf(SetStorage))
	js.Global().Set("AddStorage", js.FuncOf(AddStorage))
	js.Global().Set("RemoveStorage", js.FuncOf(RemoveStorage))
	js.Global().Set("ClearStorage", js.FuncOf(ClearStorage))

	// データ検索関数を登録
	js.Global().Set("SearchData", js.FuncOf(SearchData))

	// クイズ関連の関数を登録
	js.Global().Set("CreateQuiz", js.FuncOf(CreateQuiz))
	js.Global().Set("CreateQuizChoices", js.FuncOf(CreateQuizChoices))

	// リスニング関連の関数を登録
	js.Global().Set("GetListeningData", js.FuncOf(GetListeningData))

	// タイピング関連の関数を登録
	js.Global().Set("CreateTyping", js.FuncOf(CreateTyping))
	js.Global().Set("GetTypingQuestion", js.FuncOf(GetTypingQuestion))
	js.Global().Set("GetTypingQuestionSlice", js.FuncOf(GetTypingQuestionSlice))
	js.Global().Set("TypingKeyDown", js.FuncOf(TypingKeyDown))

	// 初期化完了をコンソールに出力
	consoleLog.Invoke(js.ValueOf("Go WASM Initialized. Registered functions."))

	// Goプログラムが終了しないように待機
	select {}
}
