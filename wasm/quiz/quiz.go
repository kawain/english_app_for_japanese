package quiz

import (
	"english_app_for_japanese/wasm/objects"
)

type Quiz struct {
	appData         *objects.AppData
	FilteredArray   []objects.Datum
	index           int
	Level           int
	numberOfOptions int
	CorrectAnswer   *objects.Datum
	OptionsArray    []objects.Datum
}

func (q *Quiz) Init(appData *objects.AppData, level int, choiceCount int) {
	q.appData = appData
	q.Level = level
	q.numberOfOptions = choiceCount
	q.index = 0
	// LocalStorageに含まれていない ID を持つ要素を抽出
	tmp := q.appData.FilterNotInStorage()
	if q.Level != 0 {
		// levelで抽出した結果を返す
		tmp = objects.FilterByLevel(tmp, q.Level)
	}
	// シャッフルされた新しいスライスを返す
	q.FilteredArray = objects.ShuffleCopy(tmp)
}

// Next は1問目からでもNext
func (q *Quiz) Next() {
	q.CorrectAnswer = &q.FilteredArray[q.index]
	q.index++
	if q.index >= len(q.FilteredArray) {
		q.index = 0
	}
	q.CreateOptionsArray()
}

func (q *Quiz) CreateOptionsArray() {
	q.OptionsArray = []objects.Datum{*q.CorrectAnswer}
	selectedIDs := make(map[int]bool)
	selectedIDs[q.CorrectAnswer.ID] = true
	maxAttempts := len(q.appData.Data) * 2
	attempts := 0
	for len(q.OptionsArray) < q.numberOfOptions && attempts < maxAttempts {
		attempts++
		candidate, _ := objects.GetRandomElement(q.appData.Data)
		if !selectedIDs[candidate.ID] {
			q.OptionsArray = append(q.OptionsArray, candidate)
			selectedIDs[candidate.ID] = true
		}
	}
	q.OptionsArray = objects.ShuffleCopy(q.OptionsArray)
}
