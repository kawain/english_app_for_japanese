package typing

import (
	"fmt"
	"testing"
)

func TestSplitTextForTyping(t *testing.T) {
	for _, v := range []string{"あきょう", "hello こんにちは", "", "しゃしん", "ちぇるのびゅいりゅ"} {
		got := SplitTextForTyping(v)
		fmt.Println(got)
	}
}

func TestKeyDown(t *testing.T) {
	CurrentTypingQuestionSlice = []string{"ら", "っ", "こ"}
	result := KeyDown("zkko", 1)
	fmt.Println(result)

	CurrentTypingQuestionSlice = []string{"ら", "っ", "ぱ"}
	result = KeyDown("aaaappa", 1)
	fmt.Println(result)

	CurrentTypingQuestionSlice = []string{"り", "ん", "ご"}
	result = KeyDown("n", 1)
	fmt.Println(result)

	CurrentTypingQuestionSlice = []string{"i", "s", "a"}
	result = KeyDown("a", 2)
	fmt.Println(result)
}
