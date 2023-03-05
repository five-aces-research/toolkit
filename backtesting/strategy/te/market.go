package te

import (
	"errors"
	"github.com/five-aces-research/toolkit/backtesting/strategy"
	"github.com/five-aces-research/toolkit/backtesting/ta"
	"github.com/five-aces-research/toolkit/fas"
	"time"
)

type Market struct {
	size                float64
	stopLong, stopShort ta.Series
}

func NewMarketOrder(size float64) *Market {
	return &Market{size: size}
}

func (m *Market) Stop(long, short ta.Series) *Market {
	m.stopLong = long
	m.stopShort = short
	return m
}

func (m *Market) CreateTrade(Side bool, ch []fas.Candle, exitCandle int, indicators []strategy.SafeFloat, sizeInUsd float64, fee strategy.Fee, pnlgraph bool) (*strategy.Trade, error) {
	if exitCandle == 0 {
		return nil, errors.New("same candle")
	}

	if Side {
		fillSize := m.size * sizeInUsd
		t := strategy.NewTrade(strategy.Fill{
			Side:  Side,
			Type:  strategy.MARKET,
			Price: ch[0].Open + fee.Slippage,
			Size:  fillSize / (ch[0].Open + fee.Slippage),
			Time:  time.Time{},
			Fee:   fillSize * fee.Taker,
		})
		t.EntrySignalTime = ch[0].StartTime
		t.Indicator = indicators
		t.Close(ch[exitCandle].Open, fee.Slippage, ch[exitCandle].StartTime, strategy.MARKET, fee.Maker)
		return t, nil
	} else {
		fillSize := m.size * sizeInUsd
		t := strategy.NewTrade(strategy.Fill{
			Side:  Side,
			Type:  strategy.MARKET,
			Price: ch[0].Open - fee.Slippage,
			Size:  fillSize / (ch[0].Open - fee.Slippage),
			Time:  time.Time{},
			Fee:   fillSize * fee.Taker,
		})
		t.EntrySignalTime = ch[0].StartTime
		t.Indicator = indicators
		t.Close(ch[exitCandle].Open, fee.Slippage, ch[exitCandle].StartTime, strategy.MARKET, fee.Taker)
		return t, nil
	}

}

func (m *Market) GetInfo() strategy.TEInfo {
	return strategy.TEInfo{
		Name:             "Market Orders",
		Info:             "",
		CandlePnlSupport: true,
	}
}
