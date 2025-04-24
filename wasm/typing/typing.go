package typing

import (
	"english_app_for_japanese/wasm/objects"
	"strings"
)

// RomajiMap はひらがな（および一部記号）と対応するローマ字入力の組み合わせを保持するマップです。
// 1つのひらがなに対して複数のローマ字入力（例: し -> si, shi, ci）が存在する場合も考慮されています。
// 拗音（きゃ、きゅ、きょなど）や促音（っ）、長音（ー）なども含まれます。
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

// Typing はタイピングゲームのデータと状態を管理する構造体です。
type Typing struct {
	appData           *objects.AppData // アプリケーション全体のデータへのポインタ
	FilteredArray     []objects.Datum  // シャッフルされた問題データのスライス
	CurrentData       *objects.Datum   // 現在表示中の問題データへのポインタ
	CurrentDataArrayE []string         // 現在の問題の英語例文 (En2) を文字単位に分割したスライス
	CurrentDataArrayJ []string         // 現在の問題の日本語かな (Kana) を文字単位（拗音含む）に分割したスライス
}

// Init は Typing 構造体を初期化します。
// アプリケーションデータ全体をシャッフルし、FilteredArray に格納します。
//
// 引数:
//   - appData: アプリケーション全体のデータ (objects.AppData) へのポインタ。
func (t *Typing) Init(appData *objects.AppData) {
	t.appData = appData
	// アプリケーションデータ全体をシャッフルしてタイピング問題リストとする
	t.FilteredArray = objects.ShuffleCopy(t.appData.Data)
}

// SetData は指定されたインデックスに対応する問題データを設定します。
// FilteredArray から該当する Datum を CurrentData に設定し、
// その Datum の En2 (英語例文) と Kana (日本語かな) を
// それぞれ文字単位に分割したスライス (CurrentDataArrayE, CurrentDataArrayJ) を生成します。
// インデックスが範囲外の場合は、有効な範囲内に調整されます。
//
// 引数:
//   - index: FilteredArray 内の問題データのインデックス。
func (t *Typing) SetData(index int) {
	allQuestions := len(t.FilteredArray)
	// インデックスが負の場合は0にする
	if index < 0 {
		index = 0
	}
	// インデックスが配列の長さを超える場合は最後の要素にする
	if index >= allQuestions {
		index = allQuestions - 1
	}
	// 現在の問題データを設定
	t.CurrentData = &t.FilteredArray[index]
	// 英語と日本語の文字配列を生成
	t.createCurrentDataArrayE()
	t.createCurrentDataArrayJ()
}

// createDataArray は与えられたテキスト文字列を、タイピングに適した文字単位のスライスに分割します。
// 特に日本語の場合、RomajiMap を参照して拗音（例: "きゃ"）などを1つの要素として扱います。
// 英語や記号は基本的に1文字ずつ分割されます。
//
// 引数:
//   - text: 分割対象の文字列。
//
// 戻り値:
//   - 分割された文字列のスライス。
func (t *Typing) createDataArray(text string) []string {
	tmp := make([]string, 0, len([]rune(text))) // rune スライスのおおよその長さで初期化
	runes := []rune(text)                       // 文字列を rune スライスに変換
	i := 0
	for i < len(runes) {
		// 次の文字が存在する場合
		if i < len(runes)-1 {
			// 現在の文字と次の文字を結合して2文字の文字列を作成
			twoChars := string(runes[i : i+2])
			// RomajiMap に2文字の組み合わせが存在するか確認 (拗音などのチェック)
			if _, exists := RomajiMap[twoChars]; exists {
				// 存在すれば2文字を1要素としてスライスに追加し、インデックスを2進める
				tmp = append(tmp, twoChars)
				i += 2
			} else {
				// 存在しなければ現在の1文字を要素としてスライスに追加し、インデックスを1進める
				tmp = append(tmp, string(runes[i:i+1]))
				i++
			}
		} else {
			// 最後の1文字を要素としてスライスに追加し、ループを終了
			tmp = append(tmp, string(runes[i:i+1]))
			i++
		}
	}
	return tmp
}

// createCurrentDataArrayE は現在の問題データ (CurrentData) の英語例文 (En2) を
// 文字単位に分割し、CurrentDataArrayE に格納します。
func (t *Typing) createCurrentDataArrayE() {
	t.CurrentDataArrayE = t.createDataArray(t.CurrentData.En2)
}

// createCurrentDataArrayJ は現在の問題データ (CurrentData) の日本語かな (Kana) を
// 文字単位（拗音などを考慮）に分割し、CurrentDataArrayJ に格納します。
func (t *Typing) createCurrentDataArrayJ() {
	t.CurrentDataArrayJ = t.createDataArray(t.CurrentData.Kana)
}

// KeyDown はユーザーのキー入力 (userInput) を受け取り、現在のタイピング問題の
// 指定されたインデックス (index) の文字と比較して、正解かどうかを判定します。
// 正解の場合、次の文字に進むための新しいインデックスを返します。
// 不正解、または入力がまだ完了していない場合は、現在のインデックスを返します。
// ローマ字入力の判定（RomajiMapに基づく複数の入力パターン、促音「っ」、撥音「ん」の特殊処理）を行います。
//
// 引数:
//   - userInput: ユーザーが入力した現在の文字列全体。
//   - index: 現在判定対象となっている文字のインデックス (CurrentDataArrayE または CurrentDataArrayJ のインデックス)。
//   - mode: 判定対象の配列を指定するモード。
//   - 1: 英語 (CurrentDataArrayE) を対象とする。
//   - 2: 日本語 (CurrentDataArrayJ) を対象とする。
//
// 戻り値:
//   - 判定後の新しい文字インデックス。
//   - 正解の場合: `index + 1` または `index + 2` (促音「っ」＋子音の場合など)。
//   - 不正解または入力未完了の場合: `index`。
func (t *Typing) KeyDown(userInput string, index, mode int) int {
	var dataArry []string
	// mode に基づいて判定対象の配列を選択
	if mode == 1 {
		dataArry = t.CurrentDataArrayE
	} else if mode == 2 {
		dataArry = t.CurrentDataArrayJ
	} else {
		// 無効な mode の場合は現在のインデックスを返す
		return index
	}

	sliceLength := len(dataArry)

	// 配列が空、またはインデックスが範囲外の場合は現在のインデックスを返す
	if sliceLength == 0 || index >= sliceLength {
		return index
	}

	// 判定対象の文字（または文字列、例: "きゃ"）を取得
	targetElement := dataArry[index]

	// --- 通常の文字（促音「っ」、撥音「ん」以外）の判定 ---
	if targetElement != "っ" && targetElement != "ん" {
		// RomajiMap に対象文字が存在するか確認 (主に日本語ひらがな)
		if words, exists := RomajiMap[targetElement]; exists {
			// 対応するローマ字入力パターンを順にチェック
			for _, word := range words {
				// ユーザー入力がローマ字パターンで終わっているか確認
				if strings.HasSuffix(userInput, word) {
					// 正解ならインデックスを1進めて返す
					return index + 1
				}
			}
		} else if strings.HasSuffix(userInput, targetElement) {
			// RomajiMap に存在しない文字（英語アルファベット、記号など）の場合、
			// ユーザー入力が対象文字で終わっているか確認
			return index + 1
		}
		// どのパターンにも一致しない場合は現在のインデックスを返す (入力途中または不正解)
		return index
	}

	// --- 特殊文字（促音「っ」、撥音「ん」）の判定 ---
	nextElement := ""
	// 次の要素が存在するか確認
	if index < sliceLength-1 {
		nextElement = dataArry[index+1]
	}

	// --- 促音「っ」の判定 ---
	if targetElement == "っ" {
		// 文末が「っ」の場合（通常は稀）
		if nextElement == "" {
			// "xtu" または "ltu" で終わっていれば正解
			for _, word := range []string{"xtu", "ltu"} {
				if strings.HasSuffix(userInput, word) {
					return index + 1
				}
			}
		} else {
			// 次の文字が存在する場合 (例: 「っこ」「っぱ」)
			// 次の文字に対応するローマ字パターンを取得
			if words, exists := RomajiMap[nextElement]; exists {
				// 判定用パターンリストを初期化 ("xtu", "ltu" は「っ」単独入力用)
				patterns := []string{"xtu", "ltu"}
				// 次の文字のローマ字パターンから促音パターンを生成
				for _, word := range words {
					// 例: word が "ko" なら firstChar は "k"
					firstChar := word[:1]
					// 例: "k" + "ko" -> "kko" をパターンに追加
					patterns = append(patterns, firstChar+word)
				}
				// 生成したパターンとユーザー入力を比較
				for _, pattern := range patterns {
					if strings.HasSuffix(userInput, pattern) {
						// "xtu" または "ltu" なら「っ」のみ入力完了 -> インデックスを1進める
						if pattern == "xtu" || pattern == "ltu" {
							return index + 1
						}
						// それ以外のパターン（例: "kko"）なら「っ」＋次の文字まで入力完了 -> インデックスを2進める
						return index + 2
					}
				}
			}
		}
	} else if targetElement == "ん" { // --- 撥音「ん」の判定 ---
		// 文末が「ん」の場合
		if nextElement == "" {
			// 文末の「ん」は "nn" で入力する必要がある
			if strings.HasSuffix(userInput, "nn") {
				return index + 1
			}
		} else {
			// 次の文字が存在する場合
			// 次の文字に対応するローマ字パターンを取得
			if words, exists := RomajiMap[nextElement]; exists {
				needsNN := false // "nn" が必要かどうかのフラグ
				// 次の文字のローマ字パターンをチェック
				for _, word := range words {
					// パターンの最初の文字を取得
					firstChar := string(word[0])
					// 次の文字が母音(a,i,u,e,o)、n、y で始まる場合、"nn" が必要
					if strings.Contains("aiueony", firstChar) {
						needsNN = true
						break
					}
				}
				// "nn" が必要で、ユーザー入力が "nn" で終わっている場合
				if needsNN && strings.HasSuffix(userInput, "nn") {
					return index + 1
				}
				// "nn" が不要で、ユーザー入力が "n" で終わっている場合
				if !needsNN && strings.HasSuffix(userInput, "n") {
					return index + 1
				}
			} else if strings.HasSuffix(userInput, "nn") {
				// 次の文字が RomajiMap にない場合 (例: んa)、"nn" で入力する必要がある
				return index + 1
			}
		}
	}

	// どの条件にも一致しない場合は現在のインデックスを返す (入力途中または不正解)
	return index
}
