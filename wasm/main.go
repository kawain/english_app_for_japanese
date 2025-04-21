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

// Promiseを返す新しい関数
func SearchAndReturnDataPromise(this js.Value, args []js.Value) any {
	consoleLog.Invoke(js.ValueOf("Go関数(SearchAndReturnDataPromise)が呼び出されました"))

	// --- 引数チェック ---
	if len(args) < 1 || args[0].Type() != js.TypeNumber {
		errorMsg := "Go関数(SearchAndReturnDataPromise)エラー: 無効な引数です。数値型のmodeが必要です。"
		consoleLog.Invoke(js.ValueOf(errorMsg))
		// 引数エラーの場合は即座に拒否されたPromiseを返す
		jsError := js.Global().Get("Error").New(errorMsg)
		// Promise.reject(error) を返す
		return js.Global().Get("Promise").Call("reject", jsError)
	}
	mode := args[0].Int()
	// --- 引数チェック終了 ---

	// Promiseコンストラクタを取得
	promiseConstructor := js.Global().Get("Promise")

	// Promiseのexecutor関数を定義
	executor := js.FuncOf(func(this js.Value, promiseArgs []js.Value) any {
		// executor関数はPromiseコンストラクタによって一度だけ実行される
		// この関数自体はリソースなので、処理完了後に解放するのが望ましい
		// resolve/rejectが呼び出された後に解放する
		resolve := promiseArgs[0]
		reject := promiseArgs[1]

		// --- 本来の処理をここで行う ---
		// Goの処理自体は同期的だが、Promiseでラップする
		// エラーハンドリングを丁寧に行い、 reject を呼び出す

		if appData.Data == nil {
			errorMsg := "Go関数(SearchAndReturnDataPromise)エラー: appDataが初期化されていません。CreateObjectsを先に呼び出してください。"
			consoleLog.Invoke(js.ValueOf(errorMsg))
			jsError := js.Global().Get("Error").New(errorMsg)
			reject.Invoke(jsError) // Promiseをreject
			return nil             // executor関数自体は何も返さない (undefined)
		}

		var results []objects.Datum
		var processError error // 処理中のエラーを補足する変数

		// switch文を関数化してエラーハンドリングしやすくする (必須ではない)
		getData := func() ([]objects.Datum, error) {
			switch mode {
			case 0:
				r := appData.FilterInStorage()
				return objects.ShuffleCopy(r), nil // ShuffleCopyがエラーを返さない前提
			case 1:
				r := appData.FilterNotInStorage()
				r = objects.FilterByLevel(r, 1)
				return objects.ShuffleCopy(r), nil
			case 2:
				r := appData.FilterNotInStorage()
				r = objects.FilterByLevel(r, 2)
				return objects.ShuffleCopy(r), nil
			default:
				return nil, fmt.Errorf("無効なmode値です: %d", mode)
			}
		}

		results, processError = getData()

		if processError != nil {
			errorMsg := fmt.Sprintf("Go関数(SearchAndReturnDataPromise)処理エラー: %v", processError)
			consoleLog.Invoke(js.ValueOf(errorMsg))
			jsError := js.Global().Get("Error").New(errorMsg)
			reject.Invoke(jsError) // Promiseをreject
			return nil             // executor関数自体は何も返さない
		}

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
		// --- データ変換終了 ---

		consoleLog.Invoke(js.ValueOf("Go関数(SearchAndReturnDataPromise): 処理成功、resolveを呼び出します"))
		resolve.Invoke(js.ValueOf(jsResult)) // Promiseをresolve

		// executor関数自体は何も返さない (undefined)
		return nil
	}) // js.FuncOf の終わり

	// executor関数はPromiseコンストラクタに渡された後、不要になるので解放する
	// 注意: executor内の処理が非同期(goroutineなど)の場合、解放タイミングはresolve/reject呼び出し後になる
	// 今回は同期的処理なので、ここで解放しても問題ないはずだが、安全のためにはexecutor内で解放するのがより堅牢
	// defer executor.Release() // 同期的なのでこれでOKのはず

	// Promiseインスタンスを生成して返す
	return promiseConstructor.New(executor)

}

func main() {
	// ラッパー関数を登録
	js.Global().Set("CreateObjects", js.FuncOf(CreateObjects))
	js.Global().Set("CreateStorage", js.FuncOf(CreateStorage))
	js.Global().Set("AddStorage", js.FuncOf(AddStorage))
	js.Global().Set("RemoveStorage", js.FuncOf(RemoveStorage))
	js.Global().Set("ClearStorage", js.FuncOf(ClearStorage))
	js.Global().Set("CreateQuiz", js.FuncOf(CreateQuiz))
	js.Global().Set("CreateQuizChoices", js.FuncOf(CreateQuizChoices))
	js.Global().Set("CreateTyping", js.FuncOf(CreateTyping))
	js.Global().Set("GetTypingQuestion", js.FuncOf(GetTypingQuestion))
	js.Global().Set("GetTypingQuestionSlice", js.FuncOf(GetTypingQuestionSlice))
	js.Global().Set("KeyDown", js.FuncOf(KeyDown))
	js.Global().Set("GetListeningQuestion", js.FuncOf(GetListeningQuestion))
	js.Global().Set("SearchAndReturnData", js.FuncOf(SearchAndReturnData))
	js.Global().Set("SearchAndReturnDataPromise", js.FuncOf(SearchAndReturnDataPromise))

	consoleLog.Invoke(js.ValueOf("Go WASM Initialized. Registered functions."))
	// プログラムを終了させない
	select {}
}
