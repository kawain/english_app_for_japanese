//go:build js && wasm

package main

import (
	"syscall/js"
)

// CreateQuiz は問題の正解を得る
func CreateQuiz(this js.Value, args []js.Value) any {
	// Promiseを返すためのハンドラ
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1] // reject関数を取得

		// 非同期処理
		go func() {
			if appData.Data == nil {
				reject.Invoke(js.ValueOf("Go関数(CreateQuiz)エラー: appDataが初期化されていません。CreateObjectを先に呼び出してください。"))
				return
			}
			if len(args) != 2 {
				reject.Invoke(js.ValueOf("Go関数(CreateQuiz)エラー: 引数は2つ必要です"))
				return
			}
			if args[0].Type() != js.TypeNumber {
				reject.Invoke(js.ValueOf("Go関数(CreateQuiz)エラー: 引数0は数値である必要があります。"))
				return
			}
			if args[1].Type() != js.TypeNumber {
				reject.Invoke(js.ValueOf("Go関数(CreateQuiz)エラー: 引数1は数値である必要があります。"))
				return
			}

			level := args[0].Int()
			choiceCount := args[1].Int()

			consoleLog.Invoke(js.ValueOf("Go関数(CreateQuiz)で使用したレベル:"), js.ValueOf(level))

			// もしもquizDataにQuizDataがない、またはレベルが変更されていたら
			if quizData.FilteredArray == nil || quizData.Level != level {
				quizData.Init(&appData, level, choiceCount)
			}

			// 次の問題へ(最初の問題含む)
			quizData.Next()
			if quizData.CorrectAnswer == nil {
				reject.Invoke(js.ValueOf("Go関数(CreateQuiz)エラー: 次の問題の取得に失敗しました。データがない可能性があります。"))
				return
			}

			// 結果をJavaScriptのオブジェクトとして返す
			jsResult := map[string]interface{}{
				"id":  quizData.CorrectAnswer.ID,
				"en":  quizData.CorrectAnswer.En,
				"jp":  quizData.CorrectAnswer.Jp,
				"en2": quizData.CorrectAnswer.En2,
				"jp2": quizData.CorrectAnswer.Jp2,
			}

			resolve.Invoke(jsResult)
		}()
		return nil
	})

	// JavaScriptのPromiseを生成
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// CreateQuizChoices は問題の選択肢を得る
func CreateQuizChoices(this js.Value, args []js.Value) any {
	// Promiseを返すためのハンドラ
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1] // reject関数を取得

		// 非同期処理
		go func() {
			if quizData.OptionsArray == nil {
				reject.Invoke(js.ValueOf("Go関数(CreateQuizChoices)エラー: 選択肢の取得に失敗しました。CreateQuizを先に呼び出してください。"))
				return
			}

			// JavaScriptのオブジェクトの配列に変換
			jsResult := make([]interface{}, len(quizData.OptionsArray))

			for i, choice := range quizData.OptionsArray {
				// 各選択肢を map[string]interface{} (JavaScriptオブジェクトに対応) に変換
				choiceObj := map[string]interface{}{
					"id": choice.ID, // IDは数値のまま渡す
					"jp": choice.Jp, // 日本語訳
				}
				jsResult[i] = choiceObj // map を interface{} としてスライスに追加
			}

			resolve.Invoke(jsResult)
		}()
		return nil
	})
	// JavaScriptのPromiseを生成
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}
