//go:build js && wasm

package main

import (
	"syscall/js"
)

// CreateTyping はJavaScriptから呼び出され、タイピングゲームで使用する単語データを初期化します。
// アプリケーションデータ (appData) をシャッフルし、タイピング用のデータセット (typingData.FilteredArray) を準備します。
//
// 引数:
//   - なし (args は使用されません)
//
// 戻り値:
//   - JavaScriptのPromiseオブジェクト。
//   - 成功時: 初期化されたタイピングデータの総数 (int) で解決されます。
//   - 失敗時: エラーメッセージ (string) で拒否されます。
//
// 処理内容:
//  1. appDataが初期化されているか確認します。
//  2. typingData.Init(&appData) を呼び出し、appData.DataをシャッフルしてtypingData.FilteredArrayに格納します。
//  3. FilteredArrayの要素数をPromiseのresolve関数に渡して返します。
//  4. エラーが発生した場合は、Promiseのreject関数にエラーメッセージを渡します。
func CreateTyping(this js.Value, args []js.Value) any {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]
		go func() {
			if appData.Data == nil {
				reject.Invoke(js.ValueOf("Go関数(CreateTyping)エラー: appDataが初期化されていません。CreateObjectを先に呼び出してください。"))
				return
			}
			typingData.Init(&appData)
			resolve.Invoke(len(typingData.FilteredArray))
		}()
		return nil
	})
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// GetTypingQuestion はJavaScriptから呼び出され、指定されたインデックスに対応する
// タイピングの問題文（英語と日本語）を取得します。
// 内部で、指定されたインデックスのデータを typingData.CurrentData に設定し、
// タイピング判定用の文字配列 (CurrentDataArrayE, CurrentDataArrayJ) も生成します。
//
// 引数:
//   - args[0]: 問題のインデックス (数値型)。typingData.FilteredArray のインデックスに対応します。
//
// 戻り値:
//   - JavaScriptのPromiseオブジェクト。
//   - 成功時: 問題文の英語 (`en2`) と日本語 (`jp2`) を含むJavaScriptオブジェクト (`{en2, jp2}`) で解決されます。
//   - 失敗時: エラーメッセージ (string) で拒否されます。
//
// 処理内容:
//  1. typingDataが初期化されているか (CreateTypingが呼ばれているか) 確認します。
//  2. 引数の数と型を検証します。
//  3. 指定されたインデックスを取得します。
//  4. typingData.SetData(index) を呼び出し、現在の問題データと文字配列を設定します。
//  5. 問題データが正常に設定されたか確認します。
//  6. 問題文の英語 (En2) と日本語 (Jp2) を含むオブジェクトを作成し、Promiseのresolve関数に渡して返します。
//  7. エラーが発生した場合は、Promiseのreject関数にエラーメッセージを渡します。
func GetTypingQuestion(this js.Value, args []js.Value) any {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]
		go func() {
			if typingData.FilteredArray == nil {
				reject.Invoke(js.ValueOf("Go関数(GetTypingQuestion)エラー: typingDataが初期化されていません。CreateTypingを先に呼び出してください。"))
				return
			}
			if len(args) != 1 {
				reject.Invoke(js.ValueOf("Go関数(GetTypingQuestion)エラー: 引数は1つ必要です"))
				return
			}
			if args[0].Type() != js.TypeNumber {
				reject.Invoke(js.ValueOf("Go関数(GetTypingQuestion)エラー: 引数は数値型である必要があります"))
				return
			}
			index := args[0].Int()
			// 問題設定
			typingData.SetData(index)
			if typingData.CurrentData == nil && typingData.CurrentDataArrayE == nil && typingData.CurrentDataArrayJ == nil {
				reject.Invoke(js.ValueOf("Go関数(GetTypingQuestion)エラー: 問題の取得に失敗しました。データがない可能性があります。"))
				return
			}
			// 結果をJavaScriptのオブジェクトとして返す
			result := map[string]interface{}{
				"en2": typingData.CurrentData.ExampleEn,
				"jp2": typingData.CurrentData.ExampleJa,
			}
			resolve.Invoke(result)
		}()
		return nil
	})
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// GetTypingQuestionSlice はJavaScriptから呼び出され、現在設定されているタイピング問題の
// 文字単位に分割された配列（英語または日本語）を取得します。
// この関数は GetTypingQuestion が呼び出された後に使用されることを想定しています。
//
// 引数:
//   - args[0]: モード (数値型)。
//   - 1: 英語の問題文の文字配列 (typingData.CurrentDataArrayE) を取得します。
//   - 2: 日本語（かな）の問題文の文字配列 (typingData.CurrentDataArrayJ) を取得します。
//
// 戻り値:
//   - JavaScriptのPromiseオブジェクト。
//   - 成功時: 指定されたモードに対応する文字配列 (JavaScriptの文字列配列) で解決されます。
//   - 失敗時: エラーメッセージ (string) で拒否されます。
//
// 処理内容:
//  1. typingDataの文字配列 (CurrentDataArrayE, CurrentDataArrayJ) が初期化されているか確認します。
//  2. 引数の数と型、およびモードの値 (1 or 2) を検証します。
//  3. 指定されたモードに基づいて、対応するGoの文字列スライスを選択します。
//  4. 選択されたGoのスライスをJavaScriptの配列に変換します。
//  5. 変換された配列をPromiseのresolve関数に渡して返します。
//  6. エラーが発生した場合は、Promiseのreject関数にエラーメッセージを渡します。
func GetTypingQuestionSlice(this js.Value, args []js.Value) any {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]
		// 非同期処理
		go func() {
			if typingData.CurrentDataArrayE == nil || typingData.CurrentDataArrayJ == nil {
				reject.Invoke(js.ValueOf("Go関数(GetTypingQuestionSlice)エラー: typingDataが初期化されていません。GetTypingQuestionを先に呼び出してください。"))
				return
			}
			if len(args) != 1 {
				reject.Invoke(js.ValueOf("Go関数(GetTypingQuestionSlice)エラー: 引数は1つ必要です"))
				return
			}
			if args[0].Type() != js.TypeNumber {
				reject.Invoke(js.ValueOf("Go関数(GetTypingQuestionSlice)エラー: 引数は数値型である必要があります"))
				return
			}
			mode := args[0].Int()
			var dataArry []string
			if mode == 1 {
				dataArry = typingData.CurrentDataArrayE
			} else if mode == 2 {
				dataArry = typingData.CurrentDataArrayJ
			} else {
				reject.Invoke(js.ValueOf("Go関数(GetTypingQuestionSlice)エラー: 引数は1か2である必要があります"))
				return
			}
			// JavaScriptの配列に変換
			jsArray := make([]interface{}, len(dataArry))
			for i, v := range dataArry {
				jsArray[i] = v
			}
			resolve.Invoke(jsArray)
		}()
		return nil
	})
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// TypingKeyDown はJavaScriptから呼び出され、ユーザーのキー入力を受け取り、
// 現在のタイピング問題と比較して、正解であれば次の文字に進むための新しいインデックスを返します。
// ローマ字入力の判定（促音「っ」や撥音「ん」の特殊処理を含む）を行います。
//
// 引数:
//   - args[0]: ユーザーが入力した現在の文字列全体 (文字列型)。
//   - args[1]: 現在判定対象となっている文字のインデックス (数値型)。
//   - args[2]: モード (数値型)。
//   - 1: 英語の問題文 (typingData.CurrentDataArrayE) に対して判定します。
//   - 2: 日本語（かな）の問題文 (typingData.CurrentDataArrayJ) に対して判定します。
//
// 戻り値:
//   - JavaScriptのPromiseオブジェクト。
//   - 成功時: 判定後の新しい文字インデックス (int) で解決されます。
//   - 入力が正しく、次の文字に進む場合は `index + 1` または `index + 2` (促音の場合など)。
//   - 入力が不正解、またはまだ文字入力が完了していない場合は現在の `index`。
//   - 失敗時: エラーメッセージ (string) で拒否されます。
//
// 処理内容:
//  1. typingDataの文字配列 (CurrentDataArrayE, CurrentDataArrayJ) が初期化されているか確認します。
//  2. 引数の数と型を検証します。
//  3. 引数からユーザー入力、現在のインデックス、モードを取得します。
//  4. typingData.KeyDown(userInput, index, mode) を呼び出し、入力判定ロジックを実行します。
//  5. KeyDown関数から返された新しいインデックスをPromiseのresolve関数に渡して返します。
//  6. エラーが発生した場合は、Promiseのreject関数にエラーメッセージを渡します。
func TypingKeyDown(this js.Value, args []js.Value) any {
	handler := js.FuncOf(func(this js.Value, promiseArgs []js.Value) interface{} {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]
		go func() {
			if typingData.CurrentDataArrayE == nil || typingData.CurrentDataArrayJ == nil {
				reject.Invoke(js.ValueOf("Go関数(TypingKeyDown)エラー: typingDataが初期化されていません。GetTypingQuestionを先に呼び出してください。"))
				return
			}
			if len(args) != 3 {
				reject.Invoke(js.ValueOf("Go関数(TypingKeyDown)エラー: 引数は3つ必要です"))
				return
			}
			if args[0].Type() != js.TypeString {
				reject.Invoke(js.ValueOf("Go関数(TypingKeyDown)エラー: 引数0は文字列型である必要があります"))
				return
			}
			if args[1].Type() != js.TypeNumber {
				reject.Invoke(js.ValueOf("Go関数(TypingKeyDown)エラー: 引数1は数値型である必要があります"))
				return
			}
			if args[2].Type() != js.TypeNumber {
				reject.Invoke(js.ValueOf("Go関数(TypingKeyDown)エラー: 引数2は数値型である必要があります"))
				return
			}
			userInput := args[0].String()
			index := args[1].Int()
			mode := args[2].Int()
			// typingパッケージのKeyDown関数を呼び出して判定し、新しいインデックスを取得
			resolve.Invoke(typingData.KeyDown(userInput, index, mode))
		}()
		return nil
	})
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}
