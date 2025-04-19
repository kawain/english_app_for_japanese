package typing

import (
	"english_app_for_japanese/wasm/objects"
	"strings"
)

var RomajiMap = map[string][]string{
	"あ":  {"a"},
	"い":  {"i"},
	"う":  {"u", "wu", "whu"},
	"え":  {"e"},
	"お":  {"o"},
	"か":  {"ka", "ca"},
	"き":  {"ki"},
	"く":  {"ku", "cu"},
	"け":  {"ke"},
	"こ":  {"ko", "co"},
	"さ":  {"sa"},
	"し":  {"si", "shi", "ci"},
	"す":  {"su"},
	"せ":  {"se", "ce"},
	"そ":  {"so"},
	"た":  {"ta"},
	"ち":  {"ti", "chi"},
	"つ":  {"tu", "tsu"},
	"て":  {"te"},
	"と":  {"to"},
	"な":  {"na"},
	"に":  {"ni"},
	"ぬ":  {"nu"},
	"ね":  {"ne"},
	"の":  {"no"},
	"は":  {"ha"},
	"ひ":  {"hi"},
	"ふ":  {"hu", "fu"},
	"へ":  {"he"},
	"ほ":  {"ho"},
	"ま":  {"ma"},
	"み":  {"mi"},
	"む":  {"mu"},
	"め":  {"me"},
	"も":  {"mo"},
	"や":  {"ya"},
	"ゆ":  {"yu"},
	"よ":  {"yo"},
	"ら":  {"ra"},
	"り":  {"ri"},
	"る":  {"ru"},
	"れ":  {"re"},
	"ろ":  {"ro"},
	"わ":  {"wa"},
	"を":  {"wo"},
	"ん":  {"n", "nn"},
	"が":  {"ga"},
	"ぎ":  {"gi"},
	"ぐ":  {"gu"},
	"げ":  {"ge"},
	"ご":  {"go"},
	"ざ":  {"za"},
	"じ":  {"zi", "ji"},
	"ず":  {"zu"},
	"ぜ":  {"ze"},
	"ぞ":  {"zo"},
	"だ":  {"da"},
	"ぢ":  {"di"},
	"づ":  {"du"},
	"で":  {"de"},
	"ど":  {"do"},
	"ば":  {"ba"},
	"び":  {"bi"},
	"ぶ":  {"bu"},
	"べ":  {"be"},
	"ぼ":  {"bo"},
	"ぱ":  {"pa"},
	"ぴ":  {"pi"},
	"ぷ":  {"pu"},
	"ぺ":  {"pe"},
	"ぽ":  {"po"},
	"ゔ":  {"vu"},
	"うぁ": {"wha", "uxa", "ula"},
	"うぃ": {"whi", "wi", "uxi", "uli"},
	"うぇ": {"whe", "we", "uxe", "ule"},
	"うぉ": {"who", "uxo", "ulo"},
	"ゔぁ": {"va", "vuxa", "vula"},
	"ゔぃ": {"vi", "vuxi", "vuli"},
	"ゔぇ": {"ve", "vuxe", "vule"},
	"ゔぉ": {"vo", "vuxo", "vulo"},
	"いぇ": {"ye", "ixe", "ile"},
	"きゃ": {"kixya", "kilya", "kya"},
	"きぃ": {"kixi", "kili", "kyi"},
	"きゅ": {"kixyu", "kilyu", "kyu"},
	"きぇ": {"kixe", "kye", "kile"},
	"きょ": {"kyo", "kilyo", "kixyo"},
	"ぎゃ": {"gilya", "gya", "gixya"},
	"ぎぃ": {"gyi", "gixi", "gili"},
	"ぎゅ": {"gilyu", "gixyu", "gyu"},
	"ぎぇ": {"gile", "gye", "gixe"},
	"ぎょ": {"gixyo", "gilyo", "gyo"},
	"しゃ": {"shilya", "shixya", "cixya", "silya", "sixya", "sya", "cilya", "sha"},
	"しぃ": {"cixi", "sili", "syi", "shixi", "shili", "cili", "sixi"},
	"しゅ": {"sixyu", "syu", "shixyu", "shilyu", "cixyu", "shu", "silyu", "cilyu"},
	"しぇ": {"sixe", "shile", "cile", "sye", "she", "sile", "shixe", "cixe"},
	"しょ": {"shixyo", "cixyo", "cilyo", "sixyo", "silyo", "sho", "shilyo", "syo"},
	"じゃ": {"zilya", "jya", "zixya", "zya", "jixya", "ja", "jilya"},
	"じぃ": {"zili", "jixi", "zixi", "zyi", "jili", "jyi"},
	"じゅ": {"jilyu", "jixyu", "ju", "jyu", "zilyu", "zyu", "zixyu"},
	"じぇ": {"zile", "jixe", "zye", "je", "jile", "zixe", "jye"},
	"じょ": {"zyo", "jyo", "jilyo", "zixyo", "jixyo", "jo", "zilyo"},
	"ちゃ": {"cya", "tya", "tixya", "chilya", "chixya", "tilya"},
	"ちぃ": {"cyi", "tili", "tixi", "tyi", "chixi", "chili"},
	"ちゅ": {"cyu", "tyu", "tilyu", "tixyu", "chilyu", "chixyu"},
	"ちぇ": {"cye", "tixe", "tile", "chile", "chixe", "tye"},
	"ちょ": {"cyo", "chixyo", "chilyo", "tixyo", "tilyo", "tyo"},
	"ぢゃ": {"dilya", "dya", "dixya"},
	"ぢぃ": {"dyi", "dixi", "dili"},
	"ぢゅ": {"dyu", "dixyu", "dilyu"},
	"ぢぇ": {"dile", "dixe", "dye"},
	"ぢょ": {"dyo", "dixyo", "dilyo"},
	"てゃ": {"texya", "telya", "tha"},
	"てぃ": {"teli", "thi", "texi"},
	"てゅ": {"texyu", "telyu", "thu"},
	"てぇ": {"tele", "texe", "the"},
	"てょ": {"texyo", "telyo", "tho"},
	"でゃ": {"dha", "dexya", "delya"},
	"でぃ": {"dhi", "deli", "dexi"},
	"でゅ": {"dexyu", "delyu", "dhu"},
	"でぇ": {"dele", "dexe", "dhe"},
	"でょ": {"dho", "delyo", "dexyo"},
	"にゃ": {"nilya", "nya", "nixya"},
	"にぃ": {"nili", "nyi", "nixi"},
	"にゅ": {"nixyu", "nyu", "nilyu"},
	"にぇ": {"nye", "nixe", "nile"},
	"にょ": {"nilyo", "nixyo", "nyo"},
	"ひゃ": {"hilya", "hya", "hixya"},
	"ひぃ": {"hyi", "hili", "hixi"},
	"ひゅ": {"hixyu", "hilyu", "hyu"},
	"ひぇ": {"hye", "hile", "hixe"},
	"ひょ": {"hixyo", "hyo", "hilyo"},
	"びゃ": {"bya", "bixya", "bilya"},
	"びぃ": {"byi", "bili", "bixi"},
	"びゅ": {"bixyu", "bilyu", "byu"},
	"びぇ": {"bye", "bile", "bixe"},
	"びょ": {"byo", "bixyo", "bilyo"},
	"ぴゃ": {"pilya", "pixya", "pya"},
	"ぴぃ": {"pixi", "pili", "pyi"},
	"ぴゅ": {"pilyu", "pyu", "pixyu"},
	"ぴぇ": {"pye", "pixe", "pile"},
	"ぴょ": {"pyo", "pilyo", "pixyo"},
	"ふぁ": {"hula", "fula", "huxa", "fuxa", "fa"},
	"ふぃ": {"fuxi", "huxi", "fuli", "huli", "fi"},
	"ふぇ": {"fuxe", "huxe", "fe", "fule", "hule"},
	"ふぉ": {"fo", "fuxo", "hulo", "fulo", "huxo"},
	"ふゃ": {"hulya", "fya", "fuxya", "fulya", "huxya"},
	"ふょ": {"fulyo", "huxyo", "fuxyo", "fyo", "hulyo"},
	"みゃ": {"mixya", "milya", "mya"},
	"みぃ": {"myi", "mili", "mixi"},
	"みゅ": {"myu", "milyu", "mixyu"},
	"みぇ": {"mile", "mixe", "mye"},
	"みょ": {"myo", "mixyo", "milyo"},
	"りゃ": {"rilya", "rixya", "rya"},
	"りぃ": {"rili", "rixi", "ryi"},
	"りゅ": {"rilyu", "rixyu", "ryu"},
	"りぇ": {"rixe", "rile", "rye"},
	"りょ": {"ryo", "rilyo", "rixyo"},
	"とぅ": {"twu", "tolu", "toxu"},
	"どぅ": {"dwu", "dolu", "doxu"},
	"ぁ":  {"xa", "la"},
	"ぃ":  {"xi", "li"},
	"ぅ":  {"xu", "lu"},
	"ぇ":  {"xe", "le"},
	"ぉ":  {"xo", "lo"},
	"ゃ":  {"xya", "lya"},
	"ゅ":  {"xyu", "lyu"},
	"ょ":  {"xyo", "lyo"},
	"ゎ":  {"xwa", "lwa"},
	"っ":  {"xtu", "ltu"},
	"ー":  {"-"},
	"、":  {","},
	"。":  {"."},
	"？":  {"?"},
	"！":  {"!"},
	"〜":  {"~"},
	"（":  {"("},
	"）":  {")"},
	"「":  {"["},
	"」":  {"]"},
	"・":  {"/"},
	"；":  {";"},
	"：":  {":"},
	"　":  {" "},
}

type Typing struct {
	appData           *objects.AppData
	FilteredArray     []objects.Datum
	index             int
	CurrentData       *objects.Datum
	CurrentDataArrayE []string
	CurrentDataArrayJ []string
}

func (t *Typing) Init(appData *objects.AppData) {
	t.appData = appData
	t.FilteredArray = objects.ShuffleCopy(t.appData.Data)
	t.index = 0
}

func (t *Typing) Next() {
	t.CurrentData = &t.FilteredArray[t.index]
	t.index++
	if t.index >= len(t.FilteredArray) {
		t.index = 0
	}
	t.createCurrentDataArrayE()
	t.createCurrentDataArrayJ()
}

func (t *Typing) createDataArray(text string) []string {
	tmp := make([]string, 0, len([]rune(text)))
	runes := []rune(text)
	i := 0
	for i < len(runes) {
		if i < len(runes)-1 {
			// 2文字を切り出し
			twoChars := string(runes[i : i+2])
			if _, exists := RomajiMap[twoChars]; exists {
				tmp = append(tmp, twoChars)
				i += 2
			} else {
				tmp = append(tmp, string(runes[i:i+1]))
				i++
			}
		} else {
			tmp = append(tmp, string(runes[i:i+1]))
			i++
		}
	}
	return tmp
}

func (t *Typing) createCurrentDataArrayE() {
	t.CurrentDataArrayE = t.createDataArray(t.CurrentData.En2)
}

func (t *Typing) createCurrentDataArrayJ() {
	t.CurrentDataArrayJ = t.createDataArray(t.CurrentData.Kana)
}

// KeyDown 引数modeは1なら英語、2なら日本語の配列を調べる
func (t *Typing) KeyDown(userInput string, index, mode int) int {
	var dataArry []string
	if mode == 1 {
		dataArry = t.CurrentDataArrayE
	} else if mode == 2 {
		dataArry = t.CurrentDataArrayJ
	} else {
		return index
	}

	sliceLength := len(dataArry)

	if sliceLength == 0 || index >= sliceLength {
		return index
	}

	targetElement := dataArry[index]

	if targetElement != "っ" && targetElement != "ん" {
		if words, exists := RomajiMap[targetElement]; exists {
			for _, word := range words {
				if strings.HasSuffix(userInput, word) {
					return index + 1
				}
			}
		} else if strings.HasSuffix(userInput, targetElement) {
			return index + 1
		}
		return index
	}

	nextElement := ""
	// 次の要素がある場合
	if index < sliceLength-1 {
		nextElement = dataArry[index+1]
	}

	if targetElement == "っ" {
		if nextElement == "" {
			// ほとんどないが「っ」で終わる文章の場合
			for _, word := range []string{"xtu", "ltu"} {
				if strings.HasSuffix(userInput, word) {
					return index + 1
				}
			}
		} else {
			// 例「っぱ」の場合
			if words, exists := RomajiMap[nextElement]; exists {
				// 促音のパターンを生成
				patterns := []string{"xtu", "ltu"}
				for _, word := range words {
					// もしも「っぱ」なら
					// 「ぱ」なら「p」を取得
					firstChar := word[:1]
					// 「ppa」にする
					// 「っぱ」になる
					patterns = append(patterns, firstChar+word)
				}
				// パターンマッチング
				for _, pattern := range patterns {
					if strings.HasSuffix(userInput, pattern) {
						if pattern == "xtu" || pattern == "ltu" {
							// この場合は「っ」のみ正解なので1進める
							return index + 1
						}
						// 「っぱ」クリアなので2つ進める
						return index + 2
					}
				}
			}
		}
	} else if targetElement == "ん" {
		if nextElement == "" {
			// 「ん」で終わる場合は「nn」
			if strings.HasSuffix(userInput, "nn") {
				return index + 1
			}
		} else {
			//
			if words, exists := RomajiMap[nextElement]; exists {
				needsNN := false
				for _, word := range words {
					firstChar := string(word[0])
					// aiueony
					// 「ん」＋「あいうえお」「な行」「や行」は
					// 「nn」で入力しないとおかしくなる
					// それ以外は「n」でいい
					if strings.Contains("aiueony", firstChar) {
						needsNN = true
						break
					}
				}
				if needsNN && strings.HasSuffix(userInput, "nn") {
					return index + 1
				}
				if !needsNN && strings.HasSuffix(userInput, "n") {
					return index + 1
				}
			} else if strings.HasSuffix(userInput, "nn") {
				// 「んa」のようなRomajiMapにない時も「nn」
				return index + 1
			}
		}
	}

	return index
}
