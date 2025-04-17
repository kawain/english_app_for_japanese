//go:build js && wasm

package main

import (
	"syscall/js"
)

// CreateQuiz は問題の正解を得る
func CreateQuiz(this js.Value, args []js.Value) any {
	consoleLog.Invoke(js.ValueOf("Go関数(CreateQuiz)が呼び出されました"))

	// 引数の数をチェック (最低2つ必要)
	if len(args) < 2 {
		consoleLog.Invoke(js.ValueOf("Go関数(CreateQuiz)エラー: 引数が不足しています。levelとchoiceCountの2つが必要です。"))
		return nil
	}

	// 引数0 (level) の型をチェック (数値型か？)
	if args[0].Type() != js.TypeNumber {
		consoleLog.Invoke(js.ValueOf("Go関数(CreateQuiz)エラー: 引数0 (level) は数値である必要があります。"))
		return nil
	}
	level := args[0].Int()

	// 引数1 (choiceCount) の型をチェック (数値型か？)
	if args[1].Type() != js.TypeNumber {
		consoleLog.Invoke(js.ValueOf("Go関数(CreateQuiz)エラー: 引数1 (choiceCount) は数値である必要があります。"))
		return nil
	}
	choiceCount := args[1].Int()

	// もしもquizDataにQuizDataがない、またはレベルが変更されていたら
	if quizData.FilteredArray == nil || quizData.Level != level {
		if appData.Data == nil {
			consoleLog.Invoke(js.ValueOf("Go関数(CreateQuiz)エラー: appDataが初期化されていません。CreateObjectsを先に呼び出してください。"))
			return nil
		}
		quizData.Init(&appData, level, choiceCount)
	}

	// 次の問題へ(最初の問題含む)
	quizData.Next()
	if quizData.CorrectAnswer == nil {
		consoleLog.Invoke(js.ValueOf("Go関数(CreateQuiz)エラー: 次の問題の取得に失敗しました。データがない可能性があります。"))
		return nil
	}

	consoleLog.Invoke(js.ValueOf("Go関数(CreateQuiz)で使用したレベル:"), js.ValueOf(level))

	// 結果をJavaScriptのオブジェクトとして返す
	result := map[string]interface{}{
		"id":  quizData.CorrectAnswer.ID,
		"en":  quizData.CorrectAnswer.En,
		"jp":  quizData.CorrectAnswer.Jp,
		"en2": quizData.CorrectAnswer.En2,
		"jp2": quizData.CorrectAnswer.Jp2,
	}

	return js.ValueOf(result)
}

// CreateQuizChoices は問題の選択肢を得る
func CreateQuizChoices(this js.Value, args []js.Value) any {
	consoleLog.Invoke(js.ValueOf("Go関数(CreateQuizChoices)が呼び出されました"))

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

	// js.ValueOfを使ってJavaScriptの配列値に変換
	return js.ValueOf(jsResult)
}
