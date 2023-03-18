package test_implementation

import (
	"fmt"
	"github.com/five-aces-research/toolkit/fas/bybit"
	"github.com/five-aces-research/toolkit/live/cle/clexer"
	"github.com/five-aces-research/toolkit/live/cle/cparser"
	"testing"
)

func TestParse(t *testing.T) {
	pr := bybit.NewPrivate("ke", "kE5jRrPgSPSuOJikF6", "JsVyekm5hgVl9SPSwxbTTx9nW7XqwAmBEVKT", true)

	tk := "buy l.BTCUSDT 100 2%"
	tl, err := clexer.Lexer(tk)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	mk := NewMockCommunicator()
	parser, err := cparser.Parse(tl, mk)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	err = parser.Evaluate(pr, mk)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}
