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

func init() {
	consoleLog = js.Global().Get("console").Get("log")
	appData = objects.AppData{}
	quizData = quiz.Quiz{}
	typingData = typing.Typing{}
	listeningData = listening.Listening{}
}

// SearchData はレベルによってデータを返す関数
func SearchData(this js.Value, args []js.Value) any {
	// Promiseを返すためのハンドラ
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1] // reject関数を取得

		// 非同期処理
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

	// JavaScriptのPromiseを生成
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

func main() {
	// ラッパー関数を登録
	js.Global().Set("CreateObject", js.FuncOf(CreateObject))
	js.Global().Set("SetStorage", js.FuncOf(SetStorage))
	js.Global().Set("AddStorage", js.FuncOf(AddStorage))
	js.Global().Set("RemoveStorage", js.FuncOf(RemoveStorage))
	js.Global().Set("ClearStorage", js.FuncOf(ClearStorage))

	js.Global().Set("SearchData", js.FuncOf(SearchData))

	js.Global().Set("CreateQuiz", js.FuncOf(CreateQuiz))
	js.Global().Set("CreateQuizChoices", js.FuncOf(CreateQuizChoices))
	js.Global().Set("CreateTyping", js.FuncOf(CreateTyping))
	js.Global().Set("GetTypingQuestion", js.FuncOf(GetTypingQuestion))
	js.Global().Set("GetTypingQuestionSlice", js.FuncOf(GetTypingQuestionSlice))
	js.Global().Set("KeyDown", js.FuncOf(KeyDown))
	js.Global().Set("GetListeningQuestion", js.FuncOf(GetListeningQuestion))

	consoleLog.Invoke(js.ValueOf("Go WASM Initialized. Registered functions."))
	// プログラムを終了させない
	select {}
}
