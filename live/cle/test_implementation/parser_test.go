package test_implementation

import (
	"fmt"
	"github.com/five-aces-research/toolkit/fas/deribit"
	"github.com/five-aces-research/toolkit/live/cle/clexer"
	"github.com/five-aces-research/toolkit/live/cle/cparser"
	"testing"
	"time"
)

const (
	buy100    = "buy BTC-PERPETUAL 100 -100"
	stop100   = "stop sell BTC-PERPETUAL 100 -200"
	cancel100 = "cancel BTC-PERPETUAL sell -stop"
)

func TestParse(t *testing.T) {
	pr := deribit.NewPrivate("ke", "dJ_gOAyI", "w4xY3ttbMUGyVci2uJ_NczsYX587nHOoz8yGuyes7bI", false)

	for _, v := range []string{stop100, cancel100} {
		tl, err := clexer.Lexer(v)
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

		fmt.Println(parser.Evaluate(pr, mk))
	}

	time.Sleep(1 * time.Second)
}
