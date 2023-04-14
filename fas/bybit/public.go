package bybit

import (
	"fmt"
	"github.com/DawnKosmos/bybit-go5"
	"github.com/DawnKosmos/bybit-go5/models"
	"github.com/five-aces-research/toolkit/fas"
	"log"
	"strconv"
	"strings"
	"time"
)

type Public struct {
	by         *bybit.Client
	tickerInfo map[string]fas.TickerInfo
	cache      Cache
}

func NewPublic() *Public {
	cl, err := bybit.New(nil, bybit.URL, nil, false)
	if err != nil {
		log.Println(err)
	}
	return &Public{by: cl, tickerInfo: make(map[string]fas.TickerInfo), cache: Cache{data: make(map[string]CacheEntry)}}
}

func (b *Public) Kline(ticker string, resolution int64, start time.Time, endTime time.Time) ([]fas.Candle, error) {
	tNow := time.Now()
	if endTime.After(tNow) {
		endTime = tNow
	}

	st, et := start.UnixMilli(), endTime.UnixMilli()
	end := et
	newRes := checkResolution(resolution)
	var hp []fas.Candle
	cat, ticker := categoryTicker(ticker)

	resString := resolutionToString(newRes)
	for {
		c, err := b.kline(cat, ticker, resString, st, end)
		if err != nil {
			log.Println(st, end)
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

func (b *Public) kline(cat, ticker string, resolution string, start int64, end int64) ([]fas.Candle, error) {
	res, err := b.by.GetKline(models.GetKlineRequest{
		Category: cat,
		Symbol:   ticker,
		Interval: resolution,
		Start:    start,
		End:      end,
		Limit:    200,
	})
	if err != nil {
		return nil, err
	}
	return bybitKlineToFas(res), nil
}

func (b *Public) GetMarketPrice(ticker string) (float64, error) {
	cat, ticker := categoryTicker(ticker)
	req := models.GetTickersRequest{
		Category: cat,
		Symbol:   ticker,
	}
	switch cat {
	case "linear", "inverse":
		var res models.GetTickersLinearResponse
		err := b.Check(2, &res, b.by.GetTickersLinear, req)
		//res, err := b.by.GetTickersLinear(req)
		if err != nil {
			return 0, err
		}
		fmt.Println(res)
		l := res.List[0]
		return median(l.Ask1Price, l.Bid1Price, l.LastPrice), nil
	case "spot":
		res, err := b.by.GetTickersSpot(req)
		if err != nil {
			return 0, err
		}
		l := res.List[0]
		return median(l.Ask1Price, l.Bid1Price, l.LastPrice), nil
	}
	return 0, nil
}

func (b *Public) GetOrderbook(ticker string, limit int) (fas.Orderbook, error) {
	cat, ticker := categoryTicker(ticker)
	var o fas.Orderbook
	res, err := b.by.GetOrderbook(models.GetOrderbookRequest{
		Category: cat,
		Symbol:   ticker,
		Limit:    int64(limit),
	})
	if err != nil {
		return o, nil
	}

	o.Ticker = ticker
	for _, v := range res.B {
		price, _ := strconv.ParseFloat(v[0], 64)
		size, _ := strconv.ParseFloat(v[1], 64)
		o.Bid = append(o.Bid, [2]float64{price, size})
	}

	for _, v := range res.A {
		price, _ := strconv.ParseFloat(v[0], 64)
		size, _ := strconv.ParseFloat(v[1], 64)
		o.Ask = append(o.Ask, [2]float64{price, size})
	}

	return o, nil
}

func (b *Public) GetTickerInfo(Ticker string) (fas.TickerInfo, error) {
	cat, ticker := categoryTicker(Ticker)
	ti, ok := b.tickerInfo[strings.ToUpper(Ticker)]
	if ok {
		return ti, nil
	}

	req := models.GetInstrumentsInfoRequest{
		Category: cat,
		Symbol:   ticker,
	}

	switch cat {
	case "linear", "inverse":
		res, err := b.by.GetInstrumentsInfoLinear(req)
		if err != nil {
			return fas.TickerInfo{}, nil
		}
		l := res.List[0]
		QtyStep, _ := strconv.ParseFloat(l.LotSizeFilter.QtyStep, 64)
		minOrder, _ := strconv.ParseFloat(l.LotSizeFilter.MinOrderQty, 64)
		tickSize, _ := strconv.ParseFloat(l.PriceFilter.TickSize, 64)

		ti = fas.TickerInfo{
			Ticker:      ticker,
			BaseCoin:    l.BaseCoin,
			QuoteCoin:   l.QuoteCoin,
			TickSize:    tickSize,
			QtyStep:     QtyStep,
			MinOrderQty: minOrder,
		}
	case "spot":
		res, err := b.by.GetInstrumentsInfoSpot(req)
		if err != nil {
			return fas.TickerInfo{}, nil
		}
		l := res.List[0]
		QtyStep, _ := strconv.ParseFloat(l.LotSizeFilter.MinOrderQty, 64)
		minOrder, _ := strconv.ParseFloat(l.LotSizeFilter.MinOrderQty, 64)
		tickSize, _ := strconv.ParseFloat(l.PriceFilter.TickSize, 64)

		ti = fas.TickerInfo{
			Ticker:      ticker,
			BaseCoin:    l.BaseCoin,
			QuoteCoin:   l.QuoteCoin,
			TickSize:    tickSize,
			QtyStep:     QtyStep,
			MinOrderQty: minOrder,
		}
	default:
		return ti, fmt.Errorf("not supported category %s", cat)
	}

	b.tickerInfo[strings.ToUpper(Ticker)] = ti
	return ti, nil
}

func (b *Public) GetFundingRate(ticker string, start, endTime time.Time) ([]fas.FundingRate, error) {
	tNow := time.Now()
	if endTime.After(tNow) {
		endTime = tNow
	}

	st, et := start.UnixMilli(), endTime.UnixMilli()
	end := et
	var hp []fas.FundingRate
	cat, ticker := categoryTicker(ticker)

	for {
		c, err := b.getFundingRate(cat, ticker, st, end)
		if err != nil {
			log.Println(st, end)
			return hp, err
		}
		if len(c) < 2 {
			break
		}

		hp = append(c, hp...)
		end = hp[0].Timestamp.UnixMilli() - 1000
	}

	return hp, nil
}

func (b *Public) getFundingRate(cat, ticker string, start, end int64) ([]fas.FundingRate, error) {
	res, err := b.by.GetFundingRateHistory(models.GetFundingRateHistoryRequest{
		Category:  cat,
		Symbol:    ticker,
		StartTime: start,
		EndTime:   end,
		Limit:     200,
	})
	if err != nil {
		return nil, err
	}

	o := make([]fas.FundingRate, len(res.List), len(res.List))
	i := len(res.List) - 1
	for _, v := range res.List {
		ff, _ := strconv.ParseFloat(v.FundingRate, 64)
		fr := fas.FundingRate{
			Rate:      ff,
			Timestamp: time.Time{},
		}
		o[i] = fr
		i--
	}

	return o, nil
}

func (b *Public) GetOpenInterest(ticker string, resolution int64, start, endTime time.Time) ([]fas.OpenInterest, error) {
	tNow := time.Now()
	if endTime.After(tNow) {
		endTime = tNow
	}
	st, et := start.UnixMilli(), endTime.UnixMilli()
	end := et
	cat, ticker := categoryTicker(ticker)
	interval := openInterestInterval(resolution)
	var hp []fas.OpenInterest

	for {
		c, err := b.getOpenInterest(cat, ticker, interval, st, end)
		if err != nil {
			log.Println(st, end)
			return hp, err
		}
		if len(c) < 2 {
			break
		}
		fmt.Println(c[0].Timestamp)

		hp = append(c, hp...)
		end = hp[0].Timestamp.UnixMilli() - 1000
	}
	return hp, nil
}

func (b *Public) getOpenInterest(cat, ticker string, resolution string, start, end int64) ([]fas.OpenInterest, error) {
	res, err := b.by.GetOpenInterest(models.GetOpenInterestRequest{
		Category:     cat,
		Symbol:       ticker,
		IntervalTime: resolution,
		StartTime:    start,
		EndTime:      end,
		Limit:        200,
		Cursor:       "",
	})
	if err != nil {
		return nil, err
	}

	o := make([]fas.OpenInterest, len(res.List), len(res.List))
	i := len(res.List) - 1
	for _, v := range res.List {
		ff, _ := strconv.ParseFloat(v.OpenInterest, 64)
		t, _ := strconv.ParseInt(v.Timestamp, 10, 64)

		fr := fas.OpenInterest{
			OI:        ff,
			Timestamp: time.Unix(t/1000, 0),
		}
		o[i] = fr
		i--
	}

	return o, nil
}

func bybitKlineToFas(r *models.GetKlineResponse) []fas.Candle {
	ch := make([]fas.Candle, len(r.List), len(r.List))
	i := len(r.List) - 1
	for _, v := range r.List {
		var c fas.Candle
		c.Open, _ = strconv.ParseFloat(v[1], 64)
		c.High, _ = strconv.ParseFloat(v[2], 64)
		c.Low, _ = strconv.ParseFloat(v[3], 64)
		c.Close, _ = strconv.ParseFloat(v[4], 64)
		c.Volume, _ = strconv.ParseFloat(v[5], 64)
		t, _ := strconv.ParseInt(v[0], 10, 64)
		c.StartTime = time.Unix(t/1000, 0)

		ch[i] = c
		i--
	}
	return ch
}

var fnRes = fas.GenerateResolutionFunc(10080, 1440, 720, 360, 240, 120, 60, 30, 15, 5, 3, 1)

func checkResolution(res int64) int64 {
	return fnRes(res)
}

func resolutionToString(i int64) string {
	switch i {
	case 1440:
		return "D"
	case 10080:
		return "W"
	default:
		return strconv.FormatInt(i, 10)
	}
}

func median(ask, bid, last string) float64 {
	a, _ := strconv.ParseFloat(ask, 64)
	b, _ := strconv.ParseFloat(bid, 64)
	c, _ := strconv.ParseFloat(last, 64)
	if (a >= b && a <= c) || (a <= b && a >= c) {
		return a
	} else if (b >= a && b <= c) || (b <= a && b >= c) {
		return b
	} else {
		return c
	}
}

func openInterestInterval(res int64) string {
	switch {
	case res > 240:
		return "1d"
	case res > 60:
		return "4h"
	case res > 30:
		return "1h"
	case res > 15:
		return "30min"
	case res > 5:
		return "15min"
	case res > 1:
		return "5min"
	default:
		return "1d"
	}
}
