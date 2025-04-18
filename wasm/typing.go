//go:build js && wasm

package main

import (
	"syscall/js"
)

// GetTypingQuestion
func GetTypingQuestion(this js.Value, args []js.Value) any {
	consoleLog.Invoke(js.ValueOf("Go関数(GetTypingQuestion)が呼び出されました"))

	if typingData.FilteredArray == nil {
		if appData.Data == nil {
			consoleLog.Invoke(js.ValueOf("Go関数(GetTypingQuestion)エラー: appDataが初期化されていません。CreateObjectsを先に呼び出してください。"))
			return nil
		}
		typingData.Init(&appData)
	}

	// 次の問題へ(最初の問題含む)
	typingData.Next()
	if typingData.CurrentData == nil {
		consoleLog.Invoke(js.ValueOf("Go関数(GetTypingQuestion)エラー: 次の問題の取得に失敗しました。データがない可能性があります。"))
		return nil
	}

	// 結果をJavaScriptのオブジェクトとして返す
	result := map[string]interface{}{
		"en2": typingData.CurrentData.En2,
		"jp2": typingData.CurrentData.Jp2,
	}
	return js.ValueOf(result)
}

// GetTypingQuestionSlice
func GetTypingQuestionSlice(this js.Value, args []js.Value) any {
	consoleLog.Invoke(js.ValueOf("Go関数(GetTypingQuestionSlice)が呼び出されました"))

	// JavaScriptの配列に変換
	jsArray := make([]interface{}, len(typingData.CurrentDataArray))
	for i, v := range typingData.CurrentDataArray {
		jsArray[i] = v
	}
	return js.ValueOf(jsArray)
}

// KeyDown
func KeyDown(this js.Value, args []js.Value) any {
	consoleLog.Invoke(js.ValueOf("Go関数(KeyDown)が呼び出されました"))

	if len(args) < 2 {
		return nil
	}

	if args[0].Type() != js.TypeString {
		consoleLog.Invoke(js.ValueOf("Go関数(KeyDown)エラー: 引数0 (userInput) は文字列である必要があります。"))
		return nil
	}
	userInput := args[0].String()

	if args[1].Type() != js.TypeNumber {
		consoleLog.Invoke(js.ValueOf("Go関数(KeyDown)エラー: 引数1 (index) は数値である必要があります。"))
		return nil
	}
	index := args[1].Int()

	return js.ValueOf(typingData.KeyDown(userInput, index))
}
