package ta

import (
	"fmt"
	"testing"
	"time"

	"github.com/five-aces-research/toolkit/fas/bybit"
)

func TestKetler(t *testing.T) {

	ch, err := bybit.NewPublic().Kline("i.BTCUSDT", 1440, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), time.Now())
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	kk := NewChart("", ch)
	_, h, c, l := kk.GetSources()
	upper, lower, ma := KetlerChannels(c, c, h, l, 20, 2.0, Ema)

	u, ld, m := upper.Data(), lower.Data(), ma.Data()

	for i := len(m) - 10; i < len(m)-1; i++ {
		fmt.Println(ld[i], m[i], u[i])
	}
}
