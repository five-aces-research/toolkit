package pg_wrapper

import (
	"errors"
	"fmt"
	"github.com/five-aces-research/toolkit/fas"
	"github.com/five-aces-research/toolkit/fas/pg_wrapper/qq"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"strings"
	"time"
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
	return r.Min.Time, r.Max.Time, err
}

func (pg *Pgx) Klines(exchange, ticker string, st time.Time, et time.Time, res int32) ([]fas.Candle, error) {
	exchangeId := getExchangeId(exchange)
	ticker = strings.ToUpper(ticker)
	id, err := pg.q.GetTickerId(ctx, qq.GetTickerIdParams{
		ExchangeID: exchangeId,
		Name:       ticker,
	})
	if err != nil {
		return nil, err
	}

	dbch, err := pg.q.ReadOHCL(ctx, qq.ReadOHCLParams{
		St: pgtype.Timestamp{
			Time:             st,
			InfinityModifier: 0,
			Valid:            true,
		},
		TickerID:   id,
		Resolution: res,
		Et: pgtype.Timestamp{
			Time:             et,
			InfinityModifier: 0,
			Valid:            true,
		},
	})
	if err != nil {
		return nil, err
	}
	return dbchToCandle(dbch), nil
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
			StartTime: v.Starttime.Time,
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

	dbch := candleToDBCandle(id, ch)
	if dbch == nil {
		return errors.New("Error Candle")
	}

	amount, err := pg.q.WriteOHCLV(ctx, dbch)
	fmt.Printf("Added %d Candles of %s", amount, ticker)
	return err
}

func candleToDBCandle(ticker_id int32, ch []fas.Candle) (res []qq.WriteOHCLVParams) {
	if len(ch) < 3 {
		return nil
	}
	resolution := (ch[1].StartTime.Unix() - ch[0].StartTime.Unix()) / 60

	for _, v := range ch {
		res = append(res, qq.WriteOHCLVParams{
			TickerID:   ticker_id,
			Resolution: int32(resolution),
			Starttime: pgtype.Timestamp{
				Time:             v.StartTime,
				InfinityModifier: 0,
				Valid:            true,
			},
			Open:   v.Open,
			High:   v.High,
			Close:  v.Close,
			Low:    v.Low,
			Volume: v.Volume,
		})
	}
	return res
}
