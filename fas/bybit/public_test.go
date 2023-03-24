package bybit

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	bb := NewPublic()
	l, err := bb.GetMarketPrice("l.BTCUSDT")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(l)
	l, err = bb.GetMarketPrice("l.BTCUSDT")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(l)

}
