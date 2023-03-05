package clexer

import (
	"fmt"
	"testing"
)

func TestLexer(t *testing.T) {
	tt, err := Lexer("x = func(a,b,c) buy a 100% -le 5 b c")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	for _, t := range tt {
		fmt.Println(t.Stringer(), t.Value)
	}
}
