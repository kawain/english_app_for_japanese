//go:build js && wasm

package main

import (
	"fmt"
	"syscall/js"
)

func GetListeningData(this js.Value, args []js.Value) any {
	// Promiseを返すためのハンドラ
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1] // reject関数を取得

		// 非同期処理
		go func() {
			if appData.Data == nil {
				reject.Invoke(js.ValueOf("Go関数(GetListeningData)エラー: appDataが初期化されていません。InitializeAppDataを先に呼び出してください。"))
				return
			}
			if len(args) != 1 {
				reject.Invoke(js.ValueOf("Go関数(GetListeningData)エラー: 引数は1つ必要です"))
				return
			}
			if args[0].Type() != js.TypeNumber {
				reject.Invoke(js.ValueOf("Go関数(GetListeningData)エラー: 引数は数値型である必要があります"))
				return
			}

			count := args[0].Int()
			consoleLog.Invoke(js.ValueOf(fmt.Sprintf("Go関数(GetListeningData) - count: %d", count)))

			// 指定されたレベルでデータを初期化（フィルタリングとシャッフル）
			// listeningData.Init に取得件数 (count) を渡す
			listeningData.Init(&appData, count)

			if len(listeningData.FilteredArray) == 0 {
				// フィルタリングの結果、データがない場合はエラーとして扱う
				reject.Invoke(js.ValueOf("Go関数(GetListeningData)エラー: フィルタリングされたデータが見つかりませんでした。"))
				return
			}

			// FilteredArrayをJavaScriptのオブジェクト配列に変換
			jsResult := make([]interface{}, len(listeningData.FilteredArray))
			for i, item := range listeningData.FilteredArray {
				jsResult[i] = map[string]interface{}{
					"id":  item.ID,
					"en":  item.Word,
					"jp":  item.DefinitionJa,
					"en2": item.ExampleEn,
					"jp2": item.ExampleJa,
				}
			}

			resolve.Invoke(jsResult)
		}()
		return nil
	})

	// JavaScriptのPromiseを生成
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}
