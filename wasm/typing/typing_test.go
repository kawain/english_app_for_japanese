package typing

import (
	"english_app_for_japanese/wasm/objects"
	"fmt"
	"testing"
)

// equalStringSlice は2つの文字列スライスが等しいか比較します。
func equalStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestCreateCurrentDataArray(t *testing.T) {
	typingInstance := Typing{}
	testCases := []struct {
		name     string
		en2      string
		kana     string
		expected []string
	}{
		{"Hiragana", "", "あきょう", []string{"あ", "きょ", "う"}},
		{"Mixed", "hello ", "こんにちは", []string{"h", "e", "l", "l", "o", " ", "こ", "ん", "に", "ち", "は"}},
		{"Empty", "", "", []string{}},
		{"Small Kana", "", "しゃしん", []string{"しゃ", "し", "ん"}},
		{"Complex Small Kana", "", "ちぇるのびゅいりゅ", []string{"ちぇ", "る", "の", "びゅ", "い", "りゅ"}},
		// カタカナは未対応
		{"With Symbols", "Go!", "ゴー！", []string{"G", "o", "!", " ", "ゴ", "ー", "！"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			typingInstance.CurrentData = &objects.Datum{En2: tc.en2, Kana: tc.kana}
			typingInstance.createCurrentDataArray()
			got := typingInstance.CurrentDataArray
			if equalStringSlice(got, tc.expected) {
				t.Errorf("createCurrentDataArray() for '%s %s' failed: expected %v, got %v", tc.en2, tc.kana, tc.expected, got)
			}
			fmt.Println(tc.name)
		})
	}
}

func TestKeyDown(t *testing.T) {
	typingInstance := Typing{} // Typing インスタンスを作成

	testCases := []struct {
		name          string
		questionSlice []string
		userInput     string
		index         int
		expectedIndex int
	}{
		{"Normal Hiragana", []string{"ら", "こ"}, "ra", 0, 1},
		{"Normal Alphabet", []string{"i", "s", "a"}, "a", 2, 3},
		{"Sokuon (xtu)", []string{"ら", "っ", "こ"}, "xtu", 1, 2},                 // 「っ」だけ入力
		{"Sokuon (ltu)", []string{"ら", "っ", "こ"}, "ltu", 1, 2},                 // 「っ」だけ入力
		{"Sokuon + Consonant (kko)", []string{"ら", "っ", "こ"}, "zkko", 1, 3},    // 「っこ」を入力
		{"Sokuon + Consonant (ppa)", []string{"ら", "っ", "ぱ"}, "aaaappa", 1, 3}, // 「っぱ」を入力
		{"N + Consonant (n)", []string{"り", "ん", "ご"}, "n", 1, 2},              // 「ん」の後に子音 -> n
		{"N + Vowel (nn)", []string{"か", "ん", "い"}, "nn", 1, 2},                // 「ん」の後に母音 -> nn
		{"N + Vowel (n, fail)", []string{"か", "ん", "い"}, "n", 1, 1},            // 「ん」の後に母音 -> n ではダメ
		{"N + N sound (nn)", []string{"ほ", "ん", "な"}, "nn", 1, 2},              // 「ん」の後に「な行」 -> nn
		{"N + N sound (n, fail)", []string{"ほ", "ん", "な"}, "n", 1, 1},          // 「ん」の後に「な行」 -> n ではダメ
		{"N + Y sound (nn)", []string{"き", "ん", "よう"}, "nn", 1, 2},             // 「ん」の後に「や行」 -> nn
		{"N + Y sound (n, fail)", []string{"き", "ん", "よう"}, "n", 1, 1},         // 「ん」の後に「や行」 -> n ではダメ
		{"N at end (nn)", []string{"ぺ", "ん"}, "nn", 1, 2},                      // 文末の「ん」 -> nn
		{"N at end (n, fail)", []string{"ぺ", "ん"}, "n", 1, 1},                  // 文末の「ん」 -> n ではダメ
		{"Out of bounds index", []string{"a", "b"}, "c", 2, 2},
		{"Empty slice", []string{}, "a", 0, 0},
		{"Romaji variation (shi)", []string{"し"}, "shi", 0, 1},
		{"Romaji variation (si)", []string{"し"}, "si", 0, 1},
		{"Romaji variation (chi)", []string{"ち"}, "chi", 0, 1},
		{"Romaji variation (ti)", []string{"ち"}, "ti", 0, 1},
		{"Romaji variation (tsu)", []string{"つ"}, "tsu", 0, 1},
		{"Romaji variation (tu)", []string{"つ"}, "tu", 0, 1},
		{"Romaji variation (ja)", []string{"じゃ"}, "ja", 0, 1},
		{"Romaji variation (jya)", []string{"じゃ"}, "jya", 0, 1},
		{"Romaji variation (zya)", []string{"じゃ"}, "zya", 0, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			typingInstance.CurrentDataArray = tc.questionSlice
			result := typingInstance.KeyDown(tc.userInput, tc.index)
			fmt.Printf("Slice: %v, Input: '%s', Index: %d -> Result: %d (Expected: %d)\n", tc.questionSlice, tc.userInput, tc.index, result, tc.expectedIndex)
			if result != tc.expectedIndex {
				t.Errorf("KeyDown(%q, %d) with slice %v failed: expected %d, got %d", tc.userInput, tc.index, tc.questionSlice, tc.expectedIndex, result)
			}
		})
	}
}
