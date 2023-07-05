package ta

import (
	"fmt"
	"github.com/five-aces-research/toolkit/fas/bybit"
	"testing"
	"time"
)

func TestMacd(t *testing.T) {
	by := bybit.NewPublic()
	ch, _ := by.Kline("i.BTCUSDT", 240, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), time.Now())
	kk := NewChart("i.BTCUSDT", ch)
	c := Close(kk)
	macd, signal, hist := MacdRelative(c, 12, 26, 9)

	for _, v := range []Series{macd, signal, hist} {
		d := v.Data()
		fmt.Println(d[len(d)-1])
	}

}
