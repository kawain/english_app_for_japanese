//go:build js && wasm

package main

import (
	"syscall/js"
)

// CreateQuiz はJavaScriptから呼び出され、指定されたレベルと選択肢の数に基づいて
// 新しいクイズの問題（正解データ）を準備し、そのデータを返します。
//
// 内部で quizData の初期化または更新、次の問題への遷移、選択肢の生成を行います。
//
// 引数:
//   - args[0]: level (数値型) - クイズの難易度レベル。
//   - 0: レベル指定なし（ローカルストレージに含まれない全データから出題）
//   - 1: レベル1のデータ（ローカルストレージに含まれないもの）から出題
//   - 2: レベル2のデータ（ローカルストレージに含まれないもの）から出題
//   - args[1]: choiceCount (数値型) - 生成する選択肢の数（正解を含む）。
//
// 戻り値:
//   - JavaScriptのPromiseオブジェクト。
//   - 成功時: 正解データの情報を含むJavaScriptオブジェクト (`{id, en, jp, en2, jp2}`) で解決されます。
//   - 失敗時: エラーメッセージで拒否されます。
//
// 処理内容:
//  1. appDataが初期化されているか確認します。
//  2. 引数の数と型を検証します。
//  3. 指定されたlevelとchoiceCountを取得します。
//  4. quizDataが未初期化、または指定されたlevelが前回と異なる場合、quizDataを初期化します。
//     (appDataから指定レベルの未学習データをフィルタリングし、シャッフルします)
//  5. quizData.Next()を呼び出し、次の問題（正解データ）を設定し、内部で選択肢も生成します。
//  6. 正解データが正常に取得できたか確認します。
//  7. 正解データをJavaScriptで扱いやすい形式 (map[string]interface{}) に変換します。
//  8. 変換された正解データをPromiseのresolve関数に渡して返します。
//  9. エラーが発生した場合は、Promiseのreject関数にエラーメッセージを渡します。
func CreateQuiz(this js.Value, args []js.Value) any {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]
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
				"en":  quizData.CorrectAnswer.Word,
				"jp":  quizData.CorrectAnswer.DefinitionJa,
				"en2": quizData.CorrectAnswer.ExampleEn,
				"jp2": quizData.CorrectAnswer.ExampleJa,
			}
			resolve.Invoke(jsResult)
		}()
		return nil
	})
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// CreateQuizChoices はJavaScriptから呼び出され、現在設定されているクイズ問題に対する
// 選択肢の配列（IDと日本語訳のみ）を返します。
//
// この関数は CreateQuiz が呼び出された後に使用されることを想定しています。
//
// 引数:
//   - なし (args は使用されません)
//
// 戻り値:
//   - JavaScriptのPromiseオブジェクト。
//   - 成功時: 選択肢の配列（各要素は `{id, jp}` のJavaScriptオブジェクト）で解決されます。
//   - 失敗時: エラーメッセージで拒否されます。
//
// 処理内容:
//  1. quizData.OptionsArray (CreateQuiz内で生成された選択肢配列) が存在するか確認します。
//  2. 存在する場合、各選択肢データ (objects.Datum) からIDと日本語訳 (Jp) を抽出し、
//     JavaScriptで扱いやすい形式 (map[string]interface{}) の配列に変換します。
//  3. 変換された選択肢の配列をPromiseのresolve関数に渡して返します。
//  4. quizData.OptionsArrayが存在しない場合（CreateQuizが未呼び出しなど）、
//     Promiseのreject関数にエラーメッセージを渡します。
func CreateQuizChoices(this js.Value, args []js.Value) any {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]
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
					"id": choice.ID,           // IDは数値のまま渡す
					"jp": choice.DefinitionJa, // 日本語訳
				}
				jsResult[i] = choiceObj // map を interface{} としてスライスに追加
			}
			resolve.Invoke(jsResult)
		}()
		return nil
	})
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}
