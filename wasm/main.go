//go:build js && wasm

package main

import (
	"english_app_for_japanese/wasm/listening"
	"english_app_for_japanese/wasm/objects"
	"english_app_for_japanese/wasm/quiz"
	"english_app_for_japanese/wasm/typing"
	"syscall/js"
)

var consoleLog js.Value
var appData objects.AppData
var quizData quiz.Quiz
var typingData typing.Typing
var listeningData listening.Listening

func init() {
	consoleLog = js.Global().Get("console").Get("log")
	appData = objects.AppData{}
	quizData = quiz.Quiz{}
	typingData = typing.Typing{}
	listeningData = listening.Listening{}
}

func SearchAndReturnData(this js.Value, args []js.Value) any {
	consoleLog.Invoke(js.ValueOf("Go関数(SearchAndReturnData)が呼び出されました"))

	if len(args) < 1 || args[0].Type() != js.TypeNumber {
		return nil
	}

	mode := args[0].Int()

	if appData.Data == nil {
		consoleLog.Invoke(js.ValueOf("Go関数(SearchAndReturnData)エラー: appDataが初期化されていません。CreateObjectsを先に呼び出してください。"))
		return nil
	}

	var results []objects.Datum

	switch mode {
	case 0:
		// LocalStorageにあるデータのランダム
		results = appData.FilterInStorage()
		results = objects.ShuffleCopy(results)

	case 1:
		// レベル1からLocalStorageにあるデータを除外したデータのランダム
		results = appData.FilterNotInStorage()
		results = objects.FilterByLevel(results, 1)
		results = objects.ShuffleCopy(results)

	case 2:
		// レベル2からLocalStorageにあるデータを除外したデータのランダム
		results = appData.FilterNotInStorage()
		results = objects.FilterByLevel(results, 2)
		results = objects.ShuffleCopy(results)

	default:
		return nil
	}

	// JavaScriptのデータに変える
	jsResult := make([]interface{}, len(results))
	for i, v := range results {
		// オブジェクト配列にする
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

	return js.ValueOf(jsResult)
}

func main() {
	// ラッパー関数を登録
	js.Global().Set("CreateObjects", js.FuncOf(CreateObjects))
	js.Global().Set("CreateStorage", js.FuncOf(CreateStorage))
	js.Global().Set("AddStorage", js.FuncOf(AddStorage))
	js.Global().Set("CreateQuiz", js.FuncOf(CreateQuiz))
	js.Global().Set("CreateQuizChoices", js.FuncOf(CreateQuizChoices))
	js.Global().Set("GetTypingQuestion", js.FuncOf(GetTypingQuestion))
	js.Global().Set("GetTypingQuestionSlice", js.FuncOf(GetTypingQuestionSlice))
	js.Global().Set("GetListeningQuestion", js.FuncOf(GetListeningQuestion))
	js.Global().Set("SearchAndReturnData", js.FuncOf(SearchAndReturnData))

	// プログラムを終了させない
	select {}
}
