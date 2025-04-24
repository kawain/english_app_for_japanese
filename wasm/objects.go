//go:build js && wasm

package main

import (
	"english_app_for_japanese/wasm/objects"
	"syscall/js"
)

// CreateObject はJavaScriptから呼び出され、単語データの2次元配列を受け取り、
// アプリケーション内部のデータ構造 (appData.Data) を初期化・構築します。
//
// 引数:
//   - args[0]: JavaScriptの2次元配列。各内部配列は以下の文字列要素を持つことを期待します:
//     [ID, En, Jp, En2, Jp2, Kana, Level, Similar]
//
// 戻り値:
//   - JavaScriptのPromiseオブジェクト。
//   - 成功時: 正常に処理されたデータ（Datumオブジェクト）の総数 (int) で解決されます。
//   - 失敗時: エラーメッセージ (string) で拒否されます。
//
// 処理内容:
//  1. 引数の数と型（JavaScriptの配列）を検証します。
//  2. JavaScriptの2次元配列をGoの `[][]string` に変換します。
//  3. 各内部配列を `objects.NewDatum` を使用して `objects.Datum` 構造体に変換します。
//  4. 変換された `Datum` オブジェクトを `appData.AddData` を使用して `appData.Data` スライスに追加します。
//  5. 最終的に `appData.Data` に追加された要素の数をPromiseのresolve関数に渡して返します。
//  6. 途中でエラーが発生した場合は、Promiseのreject関数にエラーメッセージを渡します。
func CreateObject(this js.Value, args []js.Value) any {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]
		go func() {
			if len(args) != 1 {
				reject.Invoke(js.ValueOf("Go関数(CreateObject)エラー: 引数は1つ必要です"))
				return
			}
			if args[0].Type() != js.TypeObject {
				reject.Invoke(js.ValueOf("Go関数(CreateObject)エラー: 引数は配列である必要があります"))
				return
			}
			if args[0].Length() == 0 {
				reject.Invoke(js.ValueOf("Go関数(CreateObject)エラー: 引数の配列が空です"))
				return
			}
			jsData := args[0]
			dataLen := jsData.Length()
			consoleLog.Invoke(js.ValueOf("Go関数(CreateObject)で受け取ったデータの長さ:"), js.ValueOf(dataLen))
			// Goの2次元スライスを作成
			go2DSlice := make([][]string, dataLen)
			for i := 0; i < dataLen; i++ {
				// 内側の配列を取得 (js.Value)
				jsInnerArray := jsData.Index(i)
				// 型チェック: 内側の要素が配列であることを確認
				if jsInnerArray.Type() != js.TypeObject {
					go2DSlice[i] = nil
					continue
				}
				innerLength := jsInnerArray.Length()
				// 内側のGoスライスを作成
				goInnerSlice := make([]string, innerLength)
				for j := 0; j < innerLength; j++ {
					// 要素を取得 (js.Value)
					item := jsInnerArray.Index(j)
					if item.Type() == js.TypeString {
						// 文字列に変換して格納
						goInnerSlice[j] = item.String()
					} else {
						// 文字列でない要素が含まれていた場合の処理
						goInnerSlice[j] = ""
					}
				}
				// 完成した内側スライスを外側スライスに格納
				go2DSlice[i] = goInnerSlice
			}
			for _, data := range go2DSlice {
				if data == nil {
					continue
				}
				if len(data) >= 8 {
					appData.AddData(objects.NewDatum(data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7]))
				}
			}
			consoleLog.Invoke(js.ValueOf("Go関数(CreateObject)でappData作成Length:"), js.ValueOf(len(appData.Data)))
			resolve.Invoke(len(appData.Data))
		}()
		return nil
	})
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// SetStorage はJavaScriptから呼び出され、学習済みなどの単語IDの配列を受け取り、
// アプリケーション内部のローカルストレージ情報 (appData.LocalStorage) を設定します。
// この関数を呼び出すと、既存のLocalStorage情報はクリアされ、新しい配列で上書きされます。
//
// 引数:
//   - args[0]: JavaScriptの数値配列。各要素は単語のID (int) であることを期待します。
//
// 戻り値:
//   - JavaScriptのPromiseオブジェクト。
//   - 成功時: 設定後の `appData.LocalStorage` に含まれるIDの総数 (int) で解決されます。
//   - 失敗時: エラーメッセージ (string) で拒否されます。
//
// 処理内容:
//  1. `appData.Data` が初期化されているか確認します。
//  2. 引数の数と型（JavaScriptの配列）を検証します。
//  3. `appData.LocalStorage` を空のスライスで初期化します。
//  4. JavaScript配列の各要素（数値）を `appData.AddStorage` を使用して追加します。
//  5. 最終的に `appData.LocalStorage` に追加された要素の数をPromiseのresolve関数に渡して返します。
//  6. 途中でエラーが発生した場合は、Promiseのreject関数にエラーメッセージを渡します。
func SetStorage(this js.Value, args []js.Value) any {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]
		go func() {
			if appData.Data == nil {
				reject.Invoke(js.ValueOf("Go関数(SetStorage)エラー: appDataが初期化されていません。CreateObjectを先に呼び出してください。"))
				return
			}
			if len(args) != 1 {
				reject.Invoke(js.ValueOf("Go関数(SetStorage)エラー: 引数は1つ必要です"))
				return
			}
			if args[0].Type() != js.TypeObject {
				reject.Invoke(js.ValueOf("Go関数(SetStorage)エラー: 引数は配列である必要があります"))
				return
			}
			jsData := args[0]
			dataLen := jsData.Length()
			consoleLog.Invoke(js.ValueOf("Go関数(SetStorage)で受け取ったデータの長さ:"), js.ValueOf(dataLen))
			// 空にしてから代入（ファイルをインポートした場合の対策）
			appData.LocalStorage = make([]int, 0)
			for i := range dataLen {
				// 配列の要素を取得 (js.Value)
				item := jsData.Index(i)
				if item.Type() == js.TypeNumber {
					appData.AddStorage(item.Int())
				}
			}
			consoleLog.Invoke(js.ValueOf("Go関数(SetStorage)で作成したLocalStorageの長さ:"), js.ValueOf(len(appData.LocalStorage)))
			resolve.Invoke(len(appData.LocalStorage))
		}()
		return nil
	})
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// AddStorage はJavaScriptから呼び出され、指定された単語IDを
// アプリケーション内部のローカルストレージ情報 (appData.LocalStorage) に追加します。
//
// 引数:
//   - args[0]: JavaScriptの数値。追加する単語のID (int) であることを期待します。
//
// 戻り値:
//   - JavaScriptのPromiseオブジェクト。
//   - 成功時: 追加後の `appData.LocalStorage` に含まれるIDの総数 (int) で解決されます。
//   - 失敗時: エラーメッセージ (string) で拒否されます。
//
// 処理内容:
//  1. `appData.Data` が初期化されているか確認します。
//  2. 引数の数と型（JavaScriptの数値）を検証します。
//  3. `appData.AddStorage` を呼び出して、指定されたIDを追加します。
//  4. 追加後の `appData.LocalStorage` の要素数をPromiseのresolve関数に渡して返します。
//  5. 途中でエラーが発生した場合は、Promiseのreject関数にエラーメッセージを渡します。
func AddStorage(this js.Value, args []js.Value) any {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]
		go func() {
			if appData.Data == nil {
				reject.Invoke(js.ValueOf("Go関数(AddStorage)エラー: appDataが初期化されていません。CreateObjectを先に呼び出してください。"))
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
			consoleLog.Invoke(js.ValueOf("Go関数(AddStorage)で追加前のLocalStorageの長さ:"), js.ValueOf(len(appData.LocalStorage)))
			appData.AddStorage(id)
			consoleLog.Invoke(js.ValueOf("Go関数(AddStorage)で追加後のLocalStorageの長さ:"), js.ValueOf(len(appData.LocalStorage)))
			resolve.Invoke(len(appData.LocalStorage))
		}()
		return nil
	})
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// RemoveStorage はJavaScriptから呼び出され、指定された単語IDを
// アプリケーション内部のローカルストレージ情報 (appData.LocalStorage) から削除します。
//
// 引数:
//   - args[0]: JavaScriptの数値。削除する単語のID (int) であることを期待します。
//
// 戻り値:
//   - JavaScriptのPromiseオブジェクト。
//   - 成功時: 削除後の `appData.LocalStorage` に含まれるIDの総数 (int) で解決されます。
//   - 失敗時: エラーメッセージ (string) で拒否されます。
//
// 処理内容:
//  1. `appData.Data` が初期化されているか確認します。
//  2. 引数の数と型（JavaScriptの数値）を検証します。
//  3. `appData.RemoveStorage` を呼び出して、指定されたIDを削除します。
//  4. 削除後の `appData.LocalStorage` の要素数をPromiseのresolve関数に渡して返します。
//  5. 途中でエラーが発生した場合は、Promiseのreject関数にエラーメッセージを渡します。
func RemoveStorage(this js.Value, args []js.Value) any {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]
		go func() {
			if appData.Data == nil {
				reject.Invoke(js.ValueOf("Go関数(RemoveStorage)エラー: appDataが初期化されていません。CreateObjectを先に呼び出してください。"))
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
			consoleLog.Invoke(js.ValueOf("Go関数(RemoveStorage)で追加前のLocalStorageの長さ:"), js.ValueOf(len(appData.LocalStorage)))
			appData.RemoveStorage(id)
			consoleLog.Invoke(js.ValueOf("Go関数(RemoveStorage)で追加後のLocalStorageの長さ:"), js.ValueOf(len(appData.LocalStorage)))
			resolve.Invoke(len(appData.LocalStorage))
		}()
		return nil
	})
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// ClearStorage はJavaScriptから呼び出され、
// アプリケーション内部のローカルストレージ情報 (appData.LocalStorage) をすべてクリアします。
//
// 引数:
//   - なし (args は使用されません)
//
// 戻り値:
//   - JavaScriptのPromiseオブジェクト。
//   - 成功時: クリア後の `appData.LocalStorage` の要素数 (常に 0) で解決されます。
//   - 失敗時: エラーメッセージ (string) で拒否されます。
//
// 処理内容:
//  1. `appData.Data` が初期化されているか確認します。
//  2. `appData.ClearStorage` を呼び出して、LocalStorageの内容をクリアします。
//  3. クリア後の `appData.LocalStorage` の要素数 (0) をPromiseのresolve関数に渡して返します。
//  4. 途中でエラーが発生した場合は、Promiseのreject関数にエラーメッセージを渡します。
func ClearStorage(this js.Value, args []js.Value) any {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]
		go func() {
			if appData.Data == nil {
				reject.Invoke(js.ValueOf("Go関数(ClearStorage)エラー: appDataが初期化されていません。CreateObjectを先に呼び出してください。"))
				return
			}
			consoleLog.Invoke(js.ValueOf("Go関数(ClearStorage)で削除前のLocalStorageの長さ:"), js.ValueOf(len(appData.LocalStorage)))
			appData.ClearStorage()
			consoleLog.Invoke(js.ValueOf("Go関数(ClearStorage)で削除後のLocalStorageの長さ:"), js.ValueOf(len(appData.LocalStorage)))
			resolve.Invoke(len(appData.LocalStorage))
		}()
		return nil
	})
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}
