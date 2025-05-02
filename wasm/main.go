//go:build js && wasm

package main

import (
	"encoding/json"
	"english_app_for_japanese/wasm/listening"
	"english_app_for_japanese/wasm/objects"
	"english_app_for_japanese/wasm/quiz"
	"english_app_for_japanese/wasm/typing"
	"fmt"
	"strings"
	"syscall/js"
)

const localStorageKey = "excludedWords"

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

// InitializeAppData はJavaScriptから呼び出され、アプリケーションの初期化を行います。
// 指定されたURLから単語データを非同期で取得・パースし、
// アプリケーション内部のデータ構造 (appData.Data) に追加します。
// さらに、ブラウザのローカルストレージから学習済み単語IDリスト (localStorageKey) を読み込み、
// appData.LocalStorage に設定します。
//
// 引数:
//   - なし (args は使用されません)
//
// 戻り値:
//   - JavaScriptのPromiseオブジェクト。
//   - 成功時: true で解決されます。
//   - 失敗時: エラーメッセージ (string) で拒否されます。
//
// 処理内容:
//  1. Promiseハンドラ内で非同期処理を開始します。
//  2. JavaScriptの `fetch` APIを使用して "./word.csv" を取得します。
//  3. レスポンスをテキストとして取得します。
//  4. テキストデータを改行とタブで分割し、各行を Datum オブジェクトに変換して `appData.Data` に追加します。
//  5. ブラウザの `localStorage` から `localStorageKey` に対応する値を取得し、デコードして `appData.LocalStorage` に設定します。
//  6. すべての処理が成功した場合、Promiseを `true` で解決 (resolve) します。
//  7. いずれかのステップでエラーが発生した場合、Promiseをエラーメッセージで拒否 (reject) します。
func InitializeAppData(this js.Value, args []js.Value) any {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]

		// 非同期処理をゴルーチンで実行
		go func() {
			global := js.Global()
			fetch := global.Get("fetch")
			url := "./word.csv"

			consoleLog.Invoke(js.ValueOf(fmt.Sprintf("Go関数(InitializeAppData): %s からデータを取得しています...", url)))

			// fetchを非同期で実行
			promise := fetch.Invoke(url)

			// Promiseの結果を処理
			promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				response := args[0]
				// レスポンスが成功したか確認
				if !response.Get("ok").Bool() {
					status := response.Get("status").Int()
					statusText := response.Get("statusText").String()
					errMsg := fmt.Sprintf("Go関数(InitializeAppData)エラー: Fetch failed with status %d: %s", status, statusText)
					// エラーが発生したのでPromiseをreject
					reject.Invoke(js.ValueOf(errMsg))
					return js.Undefined() // Promiseチェーンを中断
				}
				consoleLog.Invoke(js.ValueOf("Go関数(InitializeAppData): データ取得成功、テキストを取得しています..."))
				return response.Call("text")
			})).Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				// 前のthenで reject が呼ばれた場合、または undefined が返された場合
				if len(args) == 0 || args[0].IsUndefined() {
					// すでにrejectされているか、エラーで中断されたので何もしない
					return nil
				}
				data := args[0].String()
				consoleLog.Invoke(js.ValueOf("Go関数(InitializeAppData): テキスト受信完了、データを解析しています..."))

				// --- CSVパース処理 ---
				lines := strings.Split(strings.TrimSpace(data), "\n")
				addedCount := 0
				initialCount := len(appData.Data)

				for i, line := range lines {
					if i == 0 {
						continue
					}
					fields := strings.Split(line, "\t")
					if len(fields) == 9 {
						for j := range fields {
							fields[j] = strings.TrimSpace(fields[j])
						}
						if fields[0] == "" {
							consoleLog.Invoke(js.ValueOf(fmt.Sprintf("Go関数(InitializeAppData): IDが空のため %d 行目をスキップします: %s", i+1, line)))
							continue
						}
						obj := objects.NewDatum(fields[0], fields[1], fields[2], fields[3], fields[4], fields[5], fields[6], fields[7], fields[8])
						appData.AddData(obj)
						addedCount++
					} else if len(strings.TrimSpace(line)) > 0 {
						consoleLog.Invoke(js.ValueOf(fmt.Sprintf("Go関数(InitializeAppData): 不正な行 %d をスキップします (期待されるフィールド数: 9, 実際のフィールド数: %d): %s", i+1, len(fields), line)))
					}
				}
				finalCount := len(appData.Data)
				consoleLog.Invoke(js.ValueOf(fmt.Sprintf("Go関数(InitializeAppData): CSVから %d 件のデータをロードしました。合計データ数: %d (以前: %d)。", addedCount, finalCount, initialCount)))

				// --- ローカルストレージ取得処理 ---
				localStorage := global.Get("localStorage")
				storedValueJS := localStorage.Call("getItem", localStorageKey)
				if !storedValueJS.IsNull() && !storedValueJS.IsUndefined() {
					storedValueStr := storedValueJS.String()
					var loadedStorage []int
					err := json.Unmarshal([]byte(storedValueStr), &loadedStorage)
					if err != nil {
						// JSONデコード失敗はエラーとして扱う
						errMsg := fmt.Sprintf("Go関数(InitializeAppData)エラー: ローカルストレージのJSONデコード失敗: %v", err)
						consoleLog.Invoke(errMsg) // コンソールにもログを残す
						reject.Invoke(js.ValueOf(errMsg))
						return nil // 処理中断
					}

					// appData.Data に存在するIDを効率的に検索するためのセットを作成
					validDataIDs := make(map[int]struct{}, len(appData.Data))
					for _, datum := range appData.Data {
						validDataIDs[datum.ID] = struct{}{}
					}

					// 既存のLocalStorageをクリアし、有効なIDのみを追加する
					appData.LocalStorage = make([]int, 0)
					addedCount := 0
					skippedCount := 0
					for _, id := range loadedStorage {
						// appData.Data に ID が存在するか確認
						if _, exists := validDataIDs[id]; exists {
							// 存在する場合のみ追加 (AddStorageは重複チェックを行う)
							appData.AddStorage(id)
							addedCount++
						} else {
							// 存在しないIDはスキップ
							skippedCount++
						}
					}

					logMsg := fmt.Sprintf("Go関数(InitializeAppData): ローカルストレージから %d 個のIDを検証し、%d 個の有効なIDをロードしました。", len(loadedStorage), addedCount)
					if skippedCount > 0 {
						logMsg += fmt.Sprintf(" (%d 個のIDは現在のデータセットに存在しないためスキップ)", skippedCount)
					}
					consoleLog.Invoke(js.ValueOf(logMsg))

				} else {
					consoleLog.Invoke(js.ValueOf("Go関数(InitializeAppData): ローカルストレージにデータが見つかりませんでした。"))
				}

				// すべての処理が成功したのでPromiseをtrueで解決
				resolve.Invoke(js.ValueOf(true))
				return nil
			})).Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				// fetch自体やthenの中でのJavaScriptレベルのエラーハンドリング
				errMsg := fmt.Sprintf("Go関数(InitializeAppData) JavaScriptエラー: %v", args[0])
				global.Get("console").Call("error", errMsg)
				// JavaScriptエラーが発生したのでPromiseをreject
				reject.Invoke(js.ValueOf(errMsg))
				return nil
			}))
		}() // ゴルーチン開始

		// Promiseハンドラは常にnilを返す
		return nil
	})

	// JavaScriptのPromiseコンストラクタを取得して新しいPromiseを生成・返す
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
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
				// InitializeAppDataが完了していないか、失敗した可能性
				reject.Invoke(js.ValueOf("Go関数(SearchData)エラー: appDataが初期化されていません。InitializeAppDataが正常に完了したか確認してください。"))
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
					"en":    v.Word,
					"jp":    v.DefinitionJa,
					"en2":   v.ExampleEn,
					"jp2":   v.ExampleJa,
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
	// アプリケーション初期化関数を登録 (CSV読み込み + ローカルストレージ読み込み)
	js.Global().Set("InitializeAppData", js.FuncOf(InitializeAppData))

	// データ管理関連の関数を登録
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
	consoleLog.Invoke(js.ValueOf("Go WASMが初期化され、関数が登録されました。"))

	// Goプログラムが終了しないように待機
	select {}
}
