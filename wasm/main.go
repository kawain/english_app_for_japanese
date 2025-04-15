//go:build js && wasm

package main

import (
	"english_app_for_japanese/wasm/objects"
	"english_app_for_japanese/wasm/quiz"
	"fmt"
	"strconv"
	"syscall/js"
)

var consoleLog js.Value

func init() {
	consoleLog = js.Global().Get("console").Get("log")
}

// JavaScript から文字列の2次元配列を受け取る
func CreateObjects(this js.Value, args []js.Value) any {
	consoleLog.Invoke(js.ValueOf("--- Go: CreateObjects called ---"))
	if len(args) == 0 || args[0].Type() != js.TypeObject || args[0].Length() == 0 {
		consoleLog.Invoke(js.ValueOf("Go: CreateObjects received invalid or empty data."))
		return nil
	}
	jsData := args[0]
	dataLen := jsData.Length()
	consoleLog.Invoke(js.ValueOf("Go: CreateObjects received data length:"), js.ValueOf(dataLen))

	// Goの2次元スライスを作成
	go2DSlice := make([][]string, dataLen)

	for i := 0; i < dataLen; i++ {
		// 内側の配列を取得 (js.Value)
		jsInnerArray := jsData.Index(i)
		// 型チェック: 内側の要素が配列であることを確認
		if jsInnerArray.Type() != js.TypeObject || jsInnerArray.Get("length").Type() != js.TypeNumber {
			fmt.Printf("Warning: Element at outer index %d is not an array (%s), skipping.\n", i, jsInnerArray.Type())
			go2DSlice[i] = nil // または空のスライス make([]string, 0)
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
				fmt.Printf("Warning: Element at index [%d][%d] is not a string (%s), using empty string.\n", i, j, item.Type())
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

	consoleLog.Invoke(js.ValueOf("Go: CreateObjects finished assignment. objects.Objects length:"), js.ValueOf(len(objects.Objects)))
	return nil
}

// JavaScript から文字列の1次元配列を受け取る
func CreateStorage(this js.Value, args []js.Value) any {
	consoleLog.Invoke(js.ValueOf("--- Go: CreateStorage called ---"))
	if len(args) == 0 || args[0].Type() != js.TypeObject { // JavaScriptの配列は TypeObject
		consoleLog.Invoke(js.ValueOf("Go: CreateStorage received invalid or empty data."))
		return nil
	}
	jsData := args[0]
	dataLen := jsData.Length()
	consoleLog.Invoke(js.ValueOf("Go: CreateStorage received data length:"), js.ValueOf(dataLen))

	// Goのスライスを作成
	goSlice := make([]string, dataLen)

	for i := 0; i < dataLen; i++ {
		// 配列の要素を取得 (js.Value)
		item := jsData.Index(i)
		if item.Type() == js.TypeString {
			// 文字列に変換してスライスに格納
			goSlice[i] = item.String()
		} else {
			// 文字列でない要素が含まれていた場合の処理 (エラーまたはデフォルト値)
			fmt.Printf("Warning: Element at index %d is not a string (%s), using empty string.\n", i, item.Type())
			goSlice[i] = ""
		}
	}

	for _, idText := range goSlice {
		id, _ := strconv.Atoi(idText)
		objects.Storage = append(objects.Storage, id)
	}

	consoleLog.Invoke(js.ValueOf("Go: CreateStorage finished assignment. objects.Storage length:"), js.ValueOf(len(objects.Storage)))

	return nil
}

func CreateQuiz(this js.Value, args []js.Value) any {
	if len(args) < 2 {
		return nil
	}
	level := args[0].Int()
	choiceCount := args[1].Int()
	quiz.NewQuiz(level, choiceCount)
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

func CreateQuizChoices(this js.Value, args []js.Value) any {
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

func test1(this js.Value, args []js.Value) any {
	objs := objects.ShuffleObjects(objects.Objects)

	// console.log 関数を取得
	consoleLog := js.Global().Get("console").Get("log")

	consoleLog.Invoke(js.ValueOf("--- test1 called ---")) // test1が呼ばれたことをログに出す

	// objs の内容を確認
	if len(objs) > 0 {
		o1 := objs[0] // 要素があれば最初の要素を取得
		consoleLog.Invoke(js.ValueOf("First object ID:"), js.ValueOf(o1.ID))
		consoleLog.Invoke(js.ValueOf("First object En:"), js.ValueOf(o1.En))
		consoleLog.Invoke(js.ValueOf("First object Jp:"), js.ValueOf(o1.Jp))
		// 必要に応じて他のフィールドも表示
	} else {
		// スライスが空の場合のメッセージ
		consoleLog.Invoke(js.ValueOf("objs is empty."))
	}

	// objects.Storage の内容を確認
	if len(objects.Storage) > 0 {
		o2 := objects.Storage[0] // 要素があれば最初の要素を取得
		consoleLog.Invoke(js.ValueOf("First storage value:"), js.ValueOf(o2))
	} else {
		// スライスが空の場合のメッセージ
		consoleLog.Invoke(js.ValueOf("objects.Storage is empty."))
	}

	consoleLog.Invoke(js.ValueOf("--- test1 finished ---")) // test1の終了をログに出す

	return nil
}

func main() {
	// ラッパー関数を登録
	js.Global().Set("CreateObjects", js.FuncOf(CreateObjects))
	js.Global().Set("CreateStorage", js.FuncOf(CreateStorage))
	js.Global().Set("CreateQuiz", js.FuncOf(CreateQuiz))
	js.Global().Set("CreateQuizChoices", js.FuncOf(CreateQuizChoices))

	js.Global().Set("test1", js.FuncOf(test1))

	// プログラムを終了させない
	select {}
}
