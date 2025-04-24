//go:build js && wasm

package main

import "syscall/js"

func GetListeningData(this js.Value, args []js.Value) any {
	// Promiseを返すためのハンドラ
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1] // reject関数を取得

		// 非同期処理
		go func() {
			if appData.Data == nil {
				reject.Invoke(js.ValueOf("Go関数(GetListeningData)エラー: appDataが初期化されていません。CreateObjectを先に呼び出してください。"))
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

			level := args[0].Int()
			consoleLog.Invoke(js.ValueOf("Go関数(GetListeningData)で使用したレベル:"), js.ValueOf(level))

			if listeningData.FilteredArray == nil || listeningData.Level != level {
				listeningData.Init(&appData, level)
			}

			listeningData.Next()
			if listeningData.CurrentData == nil {
				reject.Invoke(js.ValueOf("Go関数(GetListeningData)エラー: 次の問題の取得に失敗しました。データがない可能性があります。"))
				return
			}

			// 結果をJavaScriptのオブジェクトとして返す
			result := map[string]interface{}{
				"id":    listeningData.CurrentData.ID,
				"en":    listeningData.CurrentData.En,
				"jp":    listeningData.CurrentData.Jp,
				"en2":   listeningData.CurrentData.En2,
				"jp2":   listeningData.CurrentData.Jp2,
				"level": listeningData.CurrentData.Level,
			}

			resolve.Invoke(result)
		}()
		return nil
	})

	// JavaScriptのPromiseを生成
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)

}
