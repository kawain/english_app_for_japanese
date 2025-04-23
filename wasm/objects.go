//go:build js && wasm

package main

import (
	"english_app_for_japanese/wasm/objects"
	"syscall/js"
)

// CreateObject はJavaScriptから文字列の2次元配列を受け取る
func CreateObject(this js.Value, args []js.Value) any {
	// Promiseを返すためのハンドラ
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1] // reject関数を取得

		// 非同期処理
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

	// JavaScriptのPromiseを生成
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// SetStorage はJavaScriptから文字列の1次元配列を受け取る
func SetStorage(this js.Value, args []js.Value) any {
	// Promiseを返すためのハンドラ
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1] // reject関数を取得

		// 非同期処理
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

	// JavaScriptのPromiseを生成
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// AddStorage はJavaScript から数値を受け取る
func AddStorage(this js.Value, args []js.Value) any {
	// Promiseを返すためのハンドラ
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1] // reject関数を取得

		// 非同期処理
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

	// JavaScriptのPromiseを生成
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// RemoveStorage はJavaScript から数値を受け取る
func RemoveStorage(this js.Value, args []js.Value) any {
	// Promiseを返すためのハンドラ
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1] // reject関数を取得

		// 非同期処理
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

	// JavaScriptのPromiseを生成
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// ClearStorage は全部消す
func ClearStorage(this js.Value, args []js.Value) any {
	// Promiseを返すためのハンドラ
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1] // reject関数を取得

		// 非同期処理
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

	// JavaScriptのPromiseを生成
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}
