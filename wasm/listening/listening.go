package listening

import "english_app_for_japanese/wasm/objects"

type Listening struct {
	appData       *objects.AppData
	FilteredArray []objects.Datum
	index         int
	Level         int
	CurrentData   *objects.Datum
}

func (l *Listening) Init(appData *objects.AppData, level int) {
	l.appData = appData
	l.Level = level
	tmp := l.appData.FilterNotInStorage()
	if l.Level != 0 {
		tmp = objects.FilterByLevel(tmp, l.Level)
	}
	l.FilteredArray = objects.ShuffleCopy(tmp)
	l.index = 0
}

func (l *Listening) Next() {
	l.CurrentData = &l.FilteredArray[l.index]
	l.index++
	if l.index >= len(l.FilteredArray) {
		l.index = 0
	}
}
