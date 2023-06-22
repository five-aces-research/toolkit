package pg_wrapper

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/five-aces-research/toolkit/fas"
	"github.com/five-aces-research/toolkit/fas/bybit"
	"github.com/five-aces-research/toolkit/fas/deribit"
	"github.com/five-aces-research/toolkit/fas/pg_wrapper/qq"
	"github.com/jackc/pgx/v5"
)

func (pg *Pgx) Klines(exchange string, ticker string, st, et time.Time, res int32) (ch []fas.Candle, err error) {
	exId := getExchangeId(exchange)
	ticker = strings.ToUpper(ticker)
	tickerId, err := pg.getTickerId(exId, ticker)
	if err != nil {
		return nil, err
	}

	st = fixStartTime(st, res)
	et = fixEndTime(et, res)

	min, max, err := pg.getMinMax(tickerId, res)

	if err != nil {
		return nil, err
	} else {
		if min.Unix() == max.Unix() {
			ch, err = downloadKlines(exId, ticker, st, et, res)
			if err != nil {
				return nil, err
			}
			dbch := candleToDBCandle(tickerId, ch, res)
			if dbch == nil {
				return nil, errors.New("error Candle")
			}

			amount, err := pg.q.WriteOHCLV(ctx, dbch)
			if err != nil {
				return nil, err
			}
			fmt.Printf("Added %d Candles of %s\n", amount, ticker)
			return ch, err
		}
	}

	if st.Before(min) {
		ch, err = downloadKlines(exId, ticker, st, min.Add(-1), res)
		if err != nil {
			log.Println(err)
		}

		if len(ch) == 0 {
			pg.q.WriteOHCLV(ctx, []qq.WriteOHCLVParams{{
				TickerID: tickerId, Resolution: res, Starttime: -1}})
			fmt.Println("lower limit already downloaded")
		} else {
			dbch := candleToDBCandle(tickerId, ch, res)
			if dbch == nil {
				return nil, errors.New("Error Candle")
			}

			amount, err := pg.q.WriteOHCLV(ctx, dbch)
			if err != nil {
				return nil, err
			}
			fmt.Printf("Added %d Candles of %s", amount, ticker)
		}
	}

	if et.After(max) {
		fmt.Println(max.Add(time.Second).UnixMilli(), et.UnixMilli())
		ch, err := downloadKlines(exId, ticker, max.Add(time.Second), et, res)
		if err != nil {
			log.Println(err)
		}
		for _, v := range ch {
			fmt.Println(v.StartTime.Unix())
		}

		if len(ch) == 0 {
			fmt.Println("Upper limit already downloaded")
		} else {
			dbch := candleToDBCandle(tickerId, ch, res)
			if dbch == nil {
				return nil, errors.New("Error Candle")
			}
			amount, err := pg.q.WriteOHCLV(ctx, dbch)
			if err != nil {
				return nil, err
			}
			fmt.Printf("Added %d Candles of %s", amount, ticker)
		}
	}

	dbch, err := pg.q.ReadOHCL(ctx, qq.ReadOHCLParams{
		St:         st.Unix(),
		TickerID:   tickerId,
		Resolution: res,
		Et:         et.Unix(),
	})
	if err != nil {
		return nil, err
	}

	return dbchToCandle(dbch), nil
}

func (pg *Pgx) getTickerId(exchangeId int32, ticker string) (int32, error) {
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
				return id, err
			}
		} else {
			return id, err
		}
	}
	return id, err
}

func (pg *Pgx) getMinMax(ticker, res int32) (min, max time.Time, err error) {
	r, err := pg.q.MinAndMax(ctx, qq.MinAndMaxParams{
		TickerID:   ticker,
		Resolution: res,
	})

	return time.Unix(r.Min, 0), time.Unix(r.Max, 0), err
}

func fixEndTime(et time.Time, res int32) time.Time {
	if time.Since(et) < time.Duration(res)*60 {
		et = time.Now()
	}

	resCut := int64(res) * 60
	tNow := et.Unix() / resCut
	return time.Unix(tNow*resCut, 0)
}

func fixStartTime(st time.Time, res int32) time.Time {
	resCut := int64(res) * 60
	ST := st.Unix() / resCut
	return time.Unix(ST*resCut+int64(res)*60, 0)
}

func downloadKlines(exId int32, ticker string, st, et time.Time, res int32) ([]fas.Candle, error) {
	ex := loadExchanger(exId)
	if ex == nil {
		return nil, fmt.Errorf("unknown Exchange with ID %d", exId)
	}

	return ex.Kline(ticker, int64(res), st, et)
}

func loadExchanger(id int32) fas.Public {
	switch id {
	case 1:
		return bybit.NewPublic()
	case 2:
		return deribit.NewPublic(false)
	default:
		return nil
	}
}
