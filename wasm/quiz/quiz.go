package quiz

import (
	"english_app_for_japanese/wasm/objects"
	"fmt"
)

// 現在のクイズの正解を保持
var CorrectAnswer objects.Object

// 現在のクイズの選択肢を保持
var QuizChoices []objects.Object

// 現在のレベルで出題可能な問題リスト（シャッフル済み）
var QuizObjects []objects.Object

// 現在のクイズレベル
var QuizLevel int

// QuizObjects のどこまで出題したかのインデックス
var QuizIndex int

// NewQuiz は新しいクイズ（問題と選択肢）を作成し、CurrentQuiz に設定する
// choiceCount を引数に追加 (例: 4択なら4)
func NewQuiz(level int, choiceCount int) {
	// レベルが変わったか、または最初の呼び出しの場合、問題リストを準備
	if QuizLevel != level || len(QuizObjects) == 0 {
		tmp := objects.FilterObjectsNotInStorage(objects.Objects) // Storageにないものを抽出
		// level0は全部
		if level != 0 {
			tmp = objects.FilterObjectsByLevel(tmp, level) // 指定レベルでフィルタリング
		}
		QuizObjects = objects.ShuffleObjects(tmp) // 問題リストをシャッフル
		QuizLevel = level
		QuizIndex = 0 // 新しいリストになったのでインデックスをリセット
	}

	// QuizObjects は既にシャッフルされているので、QuizIndex を使って順番に選択する
	if QuizIndex >= len(QuizObjects) {
		// 全ての問題を一度出題し終えた場合、再度シャッフルしてインデックスをリセット
		QuizObjects = objects.ShuffleObjects(QuizObjects)
		QuizIndex = 0
	}

	CorrectAnswer = QuizObjects[QuizIndex]
	QuizIndex++ // 次の問題に進むためにインデックスを増やす

	// --- ここから選択肢作成ロジック ---
	// 正解を含めて選択肢n個をobjects.Objectsの中から決める

	// 選択肢スライスを初期化し、まず正解を追加
	QuizChoices = []objects.Object{CorrectAnswer}
	// 既に追加された選択肢のIDを記録するためのマップ
	selectedIDs := make(map[int]struct{})
	selectedIDs[CorrectAnswer.ID] = struct{}{} // 正解のIDを記録

	// 必要な選択肢の数になるまでループ (正解は既に追加済み)
	// 安全のための試行回数上限を設定 (例: 全オブジェクト数の2倍)
	maxAttempts := len(objects.Objects) * 2
	attempts := 0
	for len(QuizChoices) < choiceCount && attempts < maxAttempts {
		attempts++
		// 全オブジェクトリストからランダムに1つ選択
		candidate := objects.SelectRandomObject(objects.Objects)
		// 候補が有効か (IDが0でないか)、かつ、まだ選択肢に追加されていないかを確認
		if candidate.ID != 0 {
			if _, exists := selectedIDs[candidate.ID]; !exists {
				// 新しいユニークな選択肢なので追加
				QuizChoices = append(QuizChoices, candidate)
				// 追加した選択肢のIDを記録
				selectedIDs[candidate.ID] = struct{}{}
			}
		} else if len(objects.Objects) == 0 {
			break // 全オブジェクトリストが空ならループを抜ける
		}
	}

	// ループ後、必要な数の選択肢が集まらなかった場合の警告
	if len(QuizChoices) < choiceCount {
		fmt.Println("Warning: Could not find enough unique choices. Found:", len(QuizChoices), "needed:", choiceCount)
	}

	// 最終的な選択肢リストをシャッフルする
	QuizChoices = objects.ShuffleObjects(QuizChoices)
}
