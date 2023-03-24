package lta

import (
	"github.com/DawnKosmos/bybit-go5/ws"
	"github.com/five-aces-research/toolkit/fas/bybit"
	"testing"
)

func TestKline(t *testing.T) {
	ss := bybit.NewStreamer(ws.LINEAR)
	ch := Kline(ss, "l.BTCUSDT", 1, 150)
	c := Close(ch)
	Rsi(c, 14)

	ch.Start()

	for {

	}
}
