//go:build js && wasm

package main

import "syscall/js"

func GetListeningQuestion(this js.Value, args []js.Value) any {
	consoleLog.Invoke(js.ValueOf("Go関数(GetListeningQuestion)が呼び出されました"))

	if len(args) < 1 {
		return nil
	}

	if args[0].Type() != js.TypeNumber {
		consoleLog.Invoke(js.ValueOf("Go関数(GetListeningQuestion)エラー: 引数0 (level) は数値である必要があります。"))
		return nil
	}

	level := args[0].Int()

	if listeningData.FilteredArray == nil || listeningData.Level != level {
		if appData.Data == nil {
			consoleLog.Invoke(js.ValueOf("Go関数(GetListeningQuestion)エラー: appDataが初期化されていません。CreateObjectsを先に呼び出してください。"))
			return nil
		}
		listeningData.Init(&appData, level)
	}

	consoleLog.Invoke(js.ValueOf("Go関数(GetListeningQuestion)で使用したレベル:"), js.ValueOf(level))

	listeningData.Next()
	if listeningData.CurrentData == nil {
		consoleLog.Invoke(js.ValueOf("Go関数(GetListeningQuestion)エラー: 次の問題の取得に失敗しました。データがない可能性があります。"))
		return nil
	}

	// Similarをどうするかは後で考えることにする

	// 結果をJavaScriptのオブジェクトとして返す
	result := map[string]interface{}{
		"id":  listeningData.CurrentData.ID,
		"en":  listeningData.CurrentData.En,
		"jp":  listeningData.CurrentData.Jp,
		"en2": listeningData.CurrentData.En2,
	}

	return js.ValueOf(result)
}
