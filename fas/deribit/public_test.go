package deribit

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestKline(t *testing.T) {
	public := NewPublic(false) // Use testnet

	ticker := "BTC"
	resolution := int64(60)
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	klines, err := public.Kline(ticker, resolution, start, end)

	assert.NoError(t, err)
	assert.NotEmpty(t, klines)
	c := klines[0]
	for _, k := range klines[1:] {
		assert.True(t, k.StartTime.Unix()-c.StartTime.Unix() == 3600)
		c = k
	}
}

func TestGetMarketPrice(t *testing.T) {
	public := NewPublic(false) // Use testnet

	ticker := "BTC"

	price, err := public.GetMarketPrice(ticker)

	assert.NoError(t, err)
	assert.True(t, price > 0)
}

func TestGetTickerInfo(t *testing.T) {
	public := NewPublic(false) // Use testnet

	ticker := "BTC"

	tickerInfo, err := public.GetTickerInfo(ticker)

	assert.NoError(t, err)
	assert.Equal(t, "BTC", tickerInfo.BaseCoin)
	assert.Equal(t, "USD", tickerInfo.QuoteCoin)
	assert.True(t, tickerInfo.TickSize > 0)
	assert.True(t, tickerInfo.QtyStep > 0)
	assert.True(t, tickerInfo.MinOrderQty > 0)
}

func TestGetOrderbook(t *testing.T) {
	public := NewPublic(false) // Use testnet

	ticker := "BTC"
	limit := 5

	orderbook, err := public.GetOrderbook(ticker, limit)

	assert.NoError(t, err)
	assert.Equal(t, "BTC-PERPETUAL", orderbook.Ticker)
	assert.Len(t, orderbook.Ask, limit)
	assert.Len(t, orderbook.Bid, limit)

	for _, ask := range orderbook.Ask {
		assert.True(t, ask[0] > 0) // Price
		assert.True(t, ask[1] > 0) // Quantity
	}

	for _, bid := range orderbook.Bid {
		assert.True(t, bid[0] > 0) // Price
		assert.True(t, bid[1] > 0) // Quantity
	}
}
