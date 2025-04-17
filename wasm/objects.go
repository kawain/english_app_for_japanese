//go:build js && wasm

package main

import (
	"english_app_for_japanese/wasm/objects"
	"strconv"
	"syscall/js"
)

// CreateObjects はJavaScript から文字列の2次元配列を受け取る
func CreateObjects(this js.Value, args []js.Value) any {
	consoleLog.Invoke(js.ValueOf("Go関数(CreateObjects)が呼び出されました"))

	if len(args) < 1 || args[0].Type() != js.TypeObject || args[0].Length() == 0 {
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
		if data == nil {
			continue
		}
		if len(data) >= 8 {
			appData.AddData(objects.NewDatum(data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7]))
		}
	}

	consoleLog.Invoke(js.ValueOf("Go関数(CreateObjects)でappData作成Length:"), js.ValueOf(len(appData.Data)))

	return nil
}

// CreateStorage はJavaScript から文字列の1次元配列を受け取る
func CreateStorage(this js.Value, args []js.Value) any {
	consoleLog.Invoke(js.ValueOf("Go関数(CreateStorage)が呼び出されました"))

	if len(args) < 1 || args[0].Type() != js.TypeObject {
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
		appData.AddStorage(id)
	}

	consoleLog.Invoke(js.ValueOf("Go関数(CreateStorage)で作成したLocalStorageの長さ:"), js.ValueOf(len(appData.LocalStorage)))

	return nil
}

// AddStorage はJavaScript から数値を受け取る
func AddStorage(this js.Value, args []js.Value) any {
	consoleLog.Invoke(js.ValueOf("Go関数(AddStorage)が呼び出されました"))

	consoleLog.Invoke(js.ValueOf("Go関数(AddStorage)で追加前のLocalStorageの長さ:"), js.ValueOf(len(appData.LocalStorage)))

	if len(args) < 1 {
		return nil
	}

	id := args[0].Int()
	appData.AddStorage(id)

	consoleLog.Invoke(js.ValueOf("Go関数(AddStorage)で追加後のLocalStorageの長さ:"), js.ValueOf(len(appData.LocalStorage)))

	return nil
}
