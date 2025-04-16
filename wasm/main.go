//go:build js && wasm

package main

import (
	"english_app_for_japanese/wasm/objects"
	"english_app_for_japanese/wasm/quiz"
	"english_app_for_japanese/wasm/typing"
	"strconv"
	"syscall/js"
)

var consoleLog js.Value

func init() {
	consoleLog = js.Global().Get("console").Get("log")
}

// CreateObjects はJavaScript から文字列の2次元配列を受け取る
func CreateObjects(this js.Value, args []js.Value) any {
	consoleLog.Invoke(js.ValueOf("Go関数(CreateObjects)が呼び出されました"))

	if len(args) == 0 || args[0].Type() != js.TypeObject || args[0].Length() == 0 {
		return nil
	}
	jsData := args[0]
	dataLen := jsData.Length()

	consoleLog.Invoke(js.ValueOf("Go関数(CreateObjects)で受け取ったデータの長さ:"), js.ValueOf(dataLen))

	// Goの2次元スライスを作成
	go2DSlice := make([][]string, dataLen)

	for i := 0; i < dataLen; i++ {
		// 内側の配列を取得 (js.Value)
		jsInnerArray := jsData.Index(i)
		// 型チェック: 内側の要素が配列であることを確認
		if jsInnerArray.Type() != js.TypeObject || jsInnerArray.Get("length").Type() != js.TypeNumber {
			go2DSlice[i] = nil
			continue
		}

		innerLength := jsInnerArray.Length()
		// 内側のGoスライスを作成
		goInnerSlice := make([]string, innerLength)

		for j := 0; j < innerLength; j++ {
			// 要素を取得 (js.Value)
			item := jsInnerArray.Index(j)
			if item.Type() == js.TypeString {
				// 文字列に変換して格納
				goInnerSlice[j] = item.String()
			} else {
				// 文字列でない要素が含まれていた場合の処理
				goInnerSlice[j] = ""
			}
		}
		// 完成した内側スライスを外側スライスに格納
		go2DSlice[i] = goInnerSlice
	}

	for _, data := range go2DSlice {
		if len(data) >= 8 {
			objects.Objects = append(objects.Objects, objects.NewObjects(data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7]))
		}
	}

	consoleLog.Invoke(js.ValueOf("Go関数(CreateObjects)で作成したobjects.Objectsの長さ:"), js.ValueOf(len(objects.Objects)))

	return nil
}

// CreateStorage はJavaScript から文字列の1次元配列を受け取る
func CreateStorage(this js.Value, args []js.Value) any {
	consoleLog.Invoke(js.ValueOf("Go関数(CreateStorage)が呼び出されました"))
	if len(args) == 0 || args[0].Type() != js.TypeObject { // JavaScriptの配列は TypeObject
		return nil
	}
	jsData := args[0]
	dataLen := jsData.Length()

	consoleLog.Invoke(js.ValueOf("Go関数(CreateStorage)で受け取ったデータの長さ:"), js.ValueOf(dataLen))

	// Goのスライスを作成
	goSlice := make([]string, dataLen)

	for i := 0; i < dataLen; i++ {
		// 配列の要素を取得 (js.Value)
		item := jsData.Index(i)
		if item.Type() == js.TypeString {
			// 文字列に変換してスライスに格納
			goSlice[i] = item.String()
		} else {
			// 文字列でない要素が含まれていた場合の処理
			goSlice[i] = ""
		}
	}

	for _, idText := range goSlice {
		id, _ := strconv.Atoi(idText)
		objects.Storage = append(objects.Storage, id)
	}

	consoleLog.Invoke(js.ValueOf("Go関数(CreateStorage)で作成したobjects.Storageの長さ:"), js.ValueOf(len(objects.Storage)))

	return nil
}

// AddStorage はJavaScript から数値を受け取る
func AddStorage(this js.Value, args []js.Value) any {
	consoleLog.Invoke(js.ValueOf("Go関数(AddStorage)が呼び出されました"))
	consoleLog.Invoke(js.ValueOf("Go関数(AddStorage)で追加前のobjects.Storageの長さ:"), js.ValueOf(len(objects.Storage)))

	if len(args) < 1 {
		return nil
	}

	id := args[0].Int()
	objects.Storage = append(objects.Storage, id)

	consoleLog.Invoke(js.ValueOf("Go関数(AddStorage)で追加後のobjects.Storageの長さ:"), js.ValueOf(len(objects.Storage)))

	return nil
}

// CreateQuiz は問題の正解を得る
func CreateQuiz(this js.Value, args []js.Value) any {
	consoleLog.Invoke(js.ValueOf("Go関数(CreateQuiz)が呼び出されました"))

	if len(args) < 2 {
		return nil
	}

	level := args[0].Int()
	choiceCount := args[1].Int()
	quiz.NewQuiz(level, choiceCount)

	consoleLog.Invoke(js.ValueOf("Go関数(CreateQuiz)で使用したレベル:"), js.ValueOf(level))
	consoleLog.Invoke(js.ValueOf("Go関数(CreateQuiz)で使用したquiz.QuizObjectsの長さ:"), js.ValueOf(len(quiz.QuizObjects)))

	// 結果をJavaScriptのオブジェクトとして返す
	result := map[string]interface{}{
		"id":  quiz.CorrectAnswer.ID,
		"en":  quiz.CorrectAnswer.En,
		"jp":  quiz.CorrectAnswer.Jp,
		"en2": quiz.CorrectAnswer.En2,
		"jp2": quiz.CorrectAnswer.Jp2,
	}

	return js.ValueOf(result)
}

// CreateQuizChoices は問題の選択肢を得る
func CreateQuizChoices(this js.Value, args []js.Value) any {
	consoleLog.Invoke(js.ValueOf("Go関数(CreateQuizChoices)が呼び出されました"))

	// JavaScriptのオブジェクトの配列に変換
	jsResult := make([]interface{}, len(quiz.QuizChoices)) // 外側のスライスは []interface{}

	for i, choice := range quiz.QuizChoices {
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

// GetTypingQuestion
func GetTypingQuestion(this js.Value, args []js.Value) any {
	consoleLog.Invoke(js.ValueOf("Go関数(GetTypingQuestion)が呼び出されました"))

	typing.PrepareQuestion()

	consoleLog.Invoke(js.ValueOf("Go関数(GetTypingQuestion)後のインデックス:"), js.ValueOf(typing.TypingQuestionIndex))

	// 結果をJavaScriptのオブジェクトとして返す
	result := map[string]interface{}{
		"en2": typing.CurrentTypingQuestion.En2,
		"jp2": typing.CurrentTypingQuestion.Jp2,
	}
	return js.ValueOf(result)
}

// GetTypingQuestionSlice
func GetTypingQuestionSlice(this js.Value, args []js.Value) any {
	consoleLog.Invoke(js.ValueOf("Go関数(GetTypingQuestionSlice)が呼び出されました"))
	// JavaScriptの配列に変換
	jsArray := make([]interface{}, len(typing.CurrentTypingQuestionSlice))
	for i, v := range typing.CurrentTypingQuestionSlice {
		jsArray[i] = v
	}
	return js.ValueOf(jsArray)
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

	// プログラムを終了させない
	select {}
}
