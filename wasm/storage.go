//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

// saveLocalStorage は appData.LocalStorage の内容をブラウザの localStorage に保存します。
// エラーが発生した場合はエラーメッセージを返します。
func saveLocalStorage() string {
	localStorage := js.Global().Get("localStorage")
	jsonData, err := json.Marshal(appData.LocalStorage)
	if err != nil {
		errMsg := fmt.Sprintf("Go関数(saveLocalStorage)エラー: ローカルストレージデータのJSONエンコード失敗: %v", err)
		consoleLog.Invoke(errMsg)
		return errMsg // エラーメッセージを返す
	}
	localStorage.Call("setItem", localStorageKey, string(jsonData))
	consoleLog.Invoke(js.ValueOf(fmt.Sprintf("Go関数(saveLocalStorage): ローカルストレージに %d 個のIDを保存しました。", len(appData.LocalStorage))))
	return "" // エラーなし
}

// SetStorage はブラウザの localStorage データをappData.LocalStorageに保存します。
// ブラウザの localStorageにインポートした後に使用する想定。
func SetStorage(this js.Value, args []js.Value) any {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]
		go func() {
			if appData.Data == nil {
				reject.Invoke(js.ValueOf("Go関数(SetStorage)エラー: appDataが初期化されていません。InitializeAppDataを先に呼び出してください。"))
				return
			}

			// --- ローカルストレージ取得処理 ---
			localStorage := js.Global().Get("localStorage")
			storedValueJS := localStorage.Call("getItem", localStorageKey)
			if !storedValueJS.IsNull() && !storedValueJS.IsUndefined() {
				storedValueStr := storedValueJS.String()
				var loadedStorage []int
				err := json.Unmarshal([]byte(storedValueStr), &loadedStorage)
				if err != nil {
					// JSONデコード失敗はエラーとして扱う
					errMsg := fmt.Sprintf("Go関数(SetStorage)エラー: ローカルストレージのJSONデコード失敗: %v", err)
					consoleLog.Invoke(errMsg) // コンソールにもログを残す
					reject.Invoke(js.ValueOf(errMsg))
					return
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

				logMsg := fmt.Sprintf("Go関数(SetStorage): ローカルストレージから %d 個のIDを検証し、%d 個の有効なIDをロードしました。", len(loadedStorage), addedCount)
				if skippedCount > 0 {
					logMsg += fmt.Sprintf(" (%d 個のIDは現在のデータセットに存在しないためスキップ)", skippedCount)
				}
				consoleLog.Invoke(js.ValueOf(logMsg))

				// 重複のないデータをローカルストレージに代入
				jsonData, err := json.Marshal(appData.LocalStorage)
				if err != nil {
					errMsg := fmt.Sprintf("Go関数(SetStorage)エラー: ローカルストレージデータのJSONエンコード失敗: %v", err)
					consoleLog.Invoke(errMsg)
					return
				}
				localStorage.Call("setItem", localStorageKey, string(jsonData))
				consoleLog.Invoke(js.ValueOf(fmt.Sprintf("Go関数(SetStorage): ローカルストレージに %d 個のIDを保存しました。", len(appData.LocalStorage))))

			} else {
				consoleLog.Invoke(js.ValueOf("Go関数(SetStorage): ローカルストレージにデータが見つかりませんでした。"))
			}
			resolve.Invoke(len(appData.LocalStorage))
		}()
		return nil
	})
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// AddStorage はJavaScriptから呼び出され、指定された単語IDを
// アプリケーション内部のローカルストレージ情報 (appData.LocalStorage) に追加し、
// ブラウザの localStorage も更新します。
//
// 引数:
//   - args[0]: JavaScriptの数値。追加する単語のID (int) であることを期待します。
//
// 戻り値:
//   - JavaScriptのPromiseオブジェクト。
//   - 成功時: 追加後の `appData.LocalStorage` に含まれるIDの総数 (int) で解決されます。
//   - 失敗時: エラーメッセージ (string) で拒否されます。
func AddStorage(this js.Value, args []js.Value) any {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]
		go func() {
			if appData.Data == nil {
				reject.Invoke(js.ValueOf("Go関数(AddStorage)エラー: appDataが初期化されていません。InitializeAppDataを先に呼び出してください。"))
				return
			}
			if len(args) != 1 {
				reject.Invoke(js.ValueOf("Go関数(AddStorage)エラー: 引数は1つ必要です"))
				return
			}
			if args[0].Type() != js.TypeNumber {
				reject.Invoke(js.ValueOf("Go関数(AddStorage)エラー: 引数は数値型が必要です"))
				return
			}
			id := args[0].Int()
			consoleLog.Invoke(js.ValueOf("Go関数(AddStorage)で追加前の内部LocalStorageの長さ:"), js.ValueOf(len(appData.LocalStorage)))
			// 1. appData.LocalStorage を更新 (重複チェックはしない)
			appData.AddStorage(id)
			consoleLog.Invoke(js.ValueOf("Go関数(AddStorage)で追加後の内部LocalStorageの長さ:"), js.ValueOf(len(appData.LocalStorage)))
			// 2. ブラウザの localStorage を更新
			errMsg := saveLocalStorage()
			if errMsg != "" {
				reject.Invoke(js.ValueOf(errMsg)) // 保存失敗
				return
			}
			// 3. 成功：更新後の要素数を返す
			resolve.Invoke(len(appData.LocalStorage))
		}()
		return nil
	})
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// RemoveStorage はJavaScriptから呼び出され、指定された単語IDを
// アプリケーション内部のローカルストレージ情報 (appData.LocalStorage) から削除し、
// ブラウザの localStorage も更新します。
//
// 引数:
//   - args[0]: JavaScriptの数値。削除する単語のID (int) であることを期待します。
//
// 戻り値:
//   - JavaScriptのPromiseオブジェクト。
//   - 成功時: 削除後の `appData.LocalStorage` に含まれるIDの総数 (int) で解決されます。
//   - 失敗時: エラーメッセージ (string) で拒否されます。
func RemoveStorage(this js.Value, args []js.Value) any {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]
		go func() {
			if appData.Data == nil {
				reject.Invoke(js.ValueOf("Go関数(RemoveStorage)エラー: appDataが初期化されていません。InitializeAppDataを先に呼び出してください。"))
				return
			}
			if len(args) != 1 {
				reject.Invoke(js.ValueOf("Go関数(RemoveStorage)エラー: 引数は1つ必要です"))
				return
			}
			if args[0].Type() != js.TypeNumber {
				reject.Invoke(js.ValueOf("Go関数(RemoveStorage)エラー: 引数は数値型が必要です"))
				return
			}
			id := args[0].Int()
			consoleLog.Invoke(js.ValueOf("Go関数(RemoveStorage)で削除前の内部LocalStorageの長さ:"), js.ValueOf(len(appData.LocalStorage)))
			// 1. appData.LocalStorage を更新
			appData.RemoveStorage(id)
			consoleLog.Invoke(js.ValueOf("Go関数(RemoveStorage)で削除後の内部LocalStorageの長さ:"), js.ValueOf(len(appData.LocalStorage)))
			// 2. ブラウザの localStorage を更新
			errMsg := saveLocalStorage()
			if errMsg != "" {
				reject.Invoke(js.ValueOf(errMsg)) // 保存失敗
				return
			}
			// 3. 成功：更新後の要素数を返す
			resolve.Invoke(len(appData.LocalStorage))
		}()
		return nil
	})
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// ClearStorage はJavaScriptから呼び出され、
// アプリケーション内部のローカルストレージ情報 (appData.LocalStorage) をすべてクリアし、
// ブラウザの localStorage からも該当データを削除します。
//
// 引数:
//   - なし (args は使用されません)
//
// 戻り値:
//   - JavaScriptのPromiseオブジェクト。
//   - 成功時: クリア後の `appData.LocalStorage` の要素数 (常に 0) で解決されます。
//   - 失敗時: エラーメッセージ (string) で拒否されます。（現状、localStorage.removeItemはエラーを投げない想定）
func ClearStorage(this js.Value, args []js.Value) any {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]
		go func() {
			if appData.Data == nil {
				reject.Invoke(js.ValueOf("Go関数(ClearStorage)エラー: appDataが初期化されていません。InitializeAppDataを先に呼び出してください。"))
				return
			}
			consoleLog.Invoke(js.ValueOf("Go関数(ClearStorage)で削除前の内部LocalStorageの長さ:"), js.ValueOf(len(appData.LocalStorage)))
			// 1. appData.LocalStorage をクリア
			appData.ClearStorage()
			consoleLog.Invoke(js.ValueOf("Go関数(ClearStorage)で削除後の内部LocalStorageの長さ:"), js.ValueOf(len(appData.LocalStorage)))
			// 2. ブラウザの localStorage から削除
			localStorage := js.Global().Get("localStorage")
			localStorage.Call("removeItem", localStorageKey)
			consoleLog.Invoke(js.ValueOf(fmt.Sprintf("Go関数(ClearStorage): ブラウザのlocalStorageからキー '%s' を削除しました。", localStorageKey)))
			// 3. 成功：クリア後の要素数 (0) を返す
			resolve.Invoke(len(appData.LocalStorage))
		}()
		return nil
	})
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}
