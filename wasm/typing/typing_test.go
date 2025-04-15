package typing

import (
	"fmt"
	"testing"
)

func TestKeyDown(t *testing.T) {
	a, b := KeyDown("o", "zzzkk", 0, []string{"ら", "っ", "こ"})
	fmt.Println(a)
	fmt.Println(b)

	a, b = KeyDown("i", "r", 0, []string{"り", "ん", "ご"})
	fmt.Println(a)
	fmt.Println(b)

	a, b = KeyDown("a", "xxx", 0, []string{"i", "s", "a"})
	fmt.Println(a)
	fmt.Println(b)
}
