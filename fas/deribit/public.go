package deribit

import (
	"errors"
	"github.com/frankrap/deribit-api"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/five-aces-research/toolkit/fas"
	"github.com/frankrap/deribit-api/models"
)

type Public struct {
	d          *deribit.Client
	tickerInfo map[string]fas.TickerInfo
}

func NewPublic(testnet bool) *Public {
	url := deribit.RealBaseURL
	if testnet {
		url = deribit.TestBaseURL
	}

	return &Public{d: deribit.New(&deribit.Configuration{
		Addr:          url,
		AutoReconnect: true,
		DebugMode:     false,
	}), tickerInfo: map[string]fas.TickerInfo{}}
}

func (p *Public) Kline(ticker string, resolution int64, start time.Time, endTime time.Time) ([]fas.Candle, error) {
	ticker = toPerpetual(ticker)

	var hp []fas.Candle

	st, et := start.UnixMilli(), endTime.UnixMilli()
	var end int64 = et
	if time.Now().UnixMilli() < end {
		end = time.Now().UnixMilli()
	}
	newRes := checkResolution(resolution)
	newResString := resolutionToString(newRes)
	for {
		c, err := p.kline(ticker, newResString, st, end)
		if err != nil {
			log.Printf("Error Kline Deribit %v", err)
			return hp, err
		}
		if len(c) < 2 {
			break
		}

		hp = append(c, hp...)
		end = hp[0].StartTime.UnixMilli() - 1000
	}

	return fas.ConvertChartResolution(newRes, resolution, hp)
}

func (p *Public) kline(ticker string, res string, start int64, end int64) ([]fas.Candle, error) {
	resp, err := p.d.GetTradingviewChartData(&models.GetTradingviewChartDataParams{InstrumentName: ticker,
		StartTimestamp: start,
		EndTimestamp:   end,
		Resolution:     res})
	if err != nil {
		return []fas.Candle{}, err
	}
	return deribitCandleToCandle(resp), nil

}

func (p *Public) GetMarketPrice(ticker string) (float64, error) {
	ticker = toPerpetual(ticker)

	resp, err := p.d.GetBookSummaryByInstrument(&models.GetBookSummaryByInstrumentParams{InstrumentName: ticker})
	if err != nil {
		return 0, err
	}
	if len(resp) == 0 {
		return 0, errors.New("len can't be zero")
	}

	v := resp[0]

	return median(v.AskPrice, v.BidPrice, v.Last), nil
}

func (p *Public) GetOrderbook(ticker string, limit int) (fas.Orderbook, error) {
	ticker = toPerpetual(ticker)
	var ff fas.Orderbook
	ff.Ticker = ticker

	resp, err := p.d.GetOrderBook(&models.GetOrderBookParams{InstrumentName: ticker, Depth: limit})
	if err != nil {
		return ff, err
	}
	for _, v := range resp.Asks {
		ff.Ask = append(ff.Ask, [2]float64{v[0], v[1]})
	}
	for _, v := range resp.Bids {
		ff.Bid = append(ff.Bid, [2]float64{v[0], v[1]})
	}

	ff.Timestamp = time.Unix(resp.Timestamp/1000, 0)
	return ff, nil
}

func (p *Public) GetTickerInfo(ticker string) (fas.TickerInfo, error) {
	ticker = toPerpetual(ticker)

	ti, ok := p.tickerInfo[strings.ToUpper(ticker)]
	if ok {
		return ti, nil
	}

	var res GetInstrumentResponse
	if err := p.d.Call("/public/get_instrument", GetInstrumentRequest{ticker}, &res); err != nil {
		return fas.TickerInfo{}, err
	}

	ti = fas.TickerInfo{
		Ticker:      ticker,
		BaseCoin:    res.BaseCurrency,
		QuoteCoin:   res.QuoteCurrency,
		TickSize:    res.TickSize,
		QtyStep:     res.MinTradeAmount,
		MinOrderQty: res.MinTradeAmount,
	}
	p.tickerInfo[ticker] = ti
	return ti, nil
}

var fnRes = fas.GenerateResolutionFunc(1440, 720, 360, 180, 120, 60, 30, 15, 10, 5, 3, 1)

func checkResolution(res int64) int64 {
	return fnRes(res)
}

func resolutionToString(i int64) string {
	switch i {
	case 1440:
		return "D"
	default:
		return strconv.FormatInt(i, 10)
	}
}

func deribitCandleToCandle(c models.GetTradingviewChartDataResponse) []fas.Candle {
	newChart := make([]fas.Candle, 0, len(c.Ticks))
	var ec fas.Candle
	for i := 0; i < len(c.Ticks); i++ {
		ec = fas.Candle{
			Close:     c.Close[i],
			Open:      c.Open[i],
			High:      c.High[i],
			Low:       c.Low[i],
			StartTime: time.Unix(c.Ticks[i]/1000, 0),
			Volume:    c.Volume[i],
		}
		newChart = append(newChart, ec)
	}
	return newChart
}

func median(a, b, c float64) float64 {
	if (a >= b && a <= c) || (a <= b && a >= c) {
		return a
	} else if (b >= a && b <= c) || (b <= a && b >= c) {
		return b
	} else {
		return c
	}
}

func toPerpetual(ticker string) (s string) {
	ticker = strings.ToUpper(ticker)
	switch ticker {
	case "BTC":
		s = "BTC-PERPETUAL"
	case "ETH":
		s = "ETH-PERPETUAL"
	default:
		// Add BTC-FUTURE which looking up the next future and sets it
		s = ticker
	}
	return s
}
