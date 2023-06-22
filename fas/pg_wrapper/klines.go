package pg_wrapper

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/five-aces-research/toolkit/fas"
	"github.com/five-aces-research/toolkit/fas/pg_wrapper/qq"
	"github.com/jackc/pgx/v5"
)

func (pg *Pgx) GetMinMax(exchange, ticker string, res int32) (min, max time.Time, err error) {
	exchangeId := getExchangeId(exchange)
	ticker = strings.ToUpper(ticker)
	id, err := pg.q.GetTickerId(ctx, qq.GetTickerIdParams{
		ExchangeID: exchangeId,
		Name:       ticker,
	})
	if err != nil {
		return
	}

	r, err := pg.q.MinAndMax(ctx, qq.MinAndMaxParams{
		TickerID:   id,
		Resolution: res,
	})
	return time.Unix(r.Min, 0), time.Unix(r.Max, 0), err
}

func dbchToCandle(dbch []qq.ReadOHCLRow) []fas.Candle {
	ch := make([]fas.Candle, 0, len(dbch))
	for _, v := range dbch {
		ch = append(ch, fas.Candle{
			Close:     v.Close,
			High:      v.High,
			Low:       v.Low,
			Open:      v.Open,
			Volume:    v.Volume,
			StartTime: time.Unix(v.Starttime, 0),
		})
	}
	return ch
}

func (pg *Pgx) CopyFromKline(exchangeId int32, ticker string, ch []fas.Candle) error {
	ticker = strings.ToUpper(ticker)
	id, err := pg.q.GetTickerId(ctx, qq.GetTickerIdParams{
		ExchangeID: exchangeId,
		Name:       ticker,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			if id, err = pg.q.CreateTicker(ctx, qq.CreateTickerParams{
				ExchangeID: exchangeId,
				Name:       ticker,
			}); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	resolution := (ch[1].StartTime.Unix() - ch[0].StartTime.Unix()) / 60

	dbch := candleToDBCandle(id, ch, int32(resolution))
	if dbch == nil {
		return errors.New("Error Candle")
	}

	amount, err := pg.q.WriteOHCLV(ctx, dbch)
	fmt.Printf("Added %d Candles of %s", amount, ticker)
	return err
}

func candleToDBCandle(ticker_id int32, ch []fas.Candle, resolution int32) (res []qq.WriteOHCLVParams) {
	if len(ch) == 0 {
		return nil
	}

	for _, v := range ch {
		res = append(res, qq.WriteOHCLVParams{
			TickerID:   ticker_id,
			Resolution: resolution,
			Starttime:  v.StartTime.Unix(),
			Open:       v.Open,
			High:       v.High,
			Close:      v.Close,
			Low:        v.Low,
			Volume:     v.Volume,
		})
	}
	return res
}

func (pg *Pgx) Update(exchange string, ticker string, res int32) error {
	_, max, err := pg.GetMinMax(exchange, ticker, res)
	if err != nil {
		return err
	}
	ex := loadExchanger(getExchangeId(exchange))
	ch, err := ex.Kline(ticker, int64(res), max.Add(time.Second), time.Now())
	if err != nil {
		return err
	}

	return pg.CopyFromKline(getExchangeId(exchange), ticker, ch)
}
