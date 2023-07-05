package algos

import (
	"fmt"
	"github.com/five-aces-research/toolkit/backtesting/ta"
	"github.com/five-aces-research/toolkit/fas/bybit"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	by := bybit.NewPublic()
	ch, err := by.Kline("l.BTCUSDT", 240, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), time.Now())
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	kk := ta.NewChart("", ch)

	_, h, c, l := kk.GetSources()
	fn := StochGen(h, c, l, 14, 6)

	b, s := fn(nil, ta.Sma, 6, 6)
	fmt.Println(len(b.Data()), len(s.Data()))
}
