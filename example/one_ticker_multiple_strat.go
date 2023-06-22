package main

import (
	"github.com/five-aces-research/toolkit/backtesting/ta"
)

type Strat struct {
	Name string
	Algo func(ch ta.Chart) (buy, sell ta.Condition)
	Res  int32
}

/*



func OneTickerMultipleStrat(db *pg_wrapper.Pgx, ticker string, cfg AlgoConfig, strat []Strat, TE strategy.TradeExecution, paras strategy.Parameter) ([]*strategy.BackTestStrategy, error) {
	var res []*strategy.BackTestStrategy
	for _, v := range strat {
		data, err := db.Klines("BYBIT", ticker, cfg.St, cfg.Et, v.Res)
		if err != nil {
			return res, err
		}
		ch := ta.NewChart(ticker, data)
		btest := strategy.NewBacktest(ch, TE, paras)
		b, s := v.Algo(ch)
		btest.AddStrategy(b, s, ticker+v.Name)
		res = append(res, btest.Results...)
	}

	return res, nil
}

*/
