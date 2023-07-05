package builder

import (
	"time"

	"github.com/five-aces-research/toolkit/backtesting/plot"
	"github.com/five-aces-research/toolkit/backtesting/strategy"
	"github.com/five-aces-research/toolkit/backtesting/ta"
	"github.com/five-aces-research/toolkit/fas/pg_wrapper"
)

/*
builder helps analyse strategies


*/

type Config struct {
	Exchange  string
	St        time.Time // Starttime
	Et        time.Time // Endtime
	Res       int32     // Resolution in Minutes
	Indicator []strategy.Indicator
}

func OneStratMultiTicker(db *pg_wrapper.Pgx, filename string, ticker []string, cfg Config, algo strategy.AlgoFunc, TE strategy.TradeExecution, paras strategy.Parameter, filter ...strategy.Filter) error {

	mt := strategy.NewMultiTicker("", algo, TE, paras)
	charts := make([]ta.Chart, 0, len(ticker))
	for _, v := range ticker {
		ch, err := db.Klines(cfg.Exchange, v, cfg.St, cfg.Et, cfg.Res)
		if err != nil {
			return err
		}
		charts = append(charts, ta.NewChart(v, ch))
	}
	if cfg.Indicator != nil {
		mt.AddIndicators(cfg.Indicator...)
	}

	mt.AddTickers(charts...)

	return plot.SimplePnl(filename, nil, paras.Balance, mt.Results)
}

func OneTickerMultiTradeExecution(db *pg_wrapper.Pgx, filename string, ticker string, cfg Config, algo strategy.AlgoFunc, paras strategy.Parameter, TE ...strategy.TradeExecution) error {

	ch, err := db.Klines(cfg.Exchange, ticker, cfg.St, cfg.Et, cfg.Res)
	if err != nil {
		return err
	}

	kline := ta.NewChart(ticker, ch)
	var bt []*strategy.BackTestStrategy
	b, s := algo(kline)
	for _, v := range TE {
		strat := strategy.NewBacktest(kline, v, paras)

		strat.AddStrategy(b, s, v.GetInfo().Name+v.GetInfo().Info)
		bt = append(bt, strat.Results...)
	}

	return plot.SimplePnl(filename, nil, paras.Balance, bt)
}

type Strat struct {
	Name string
	Algo func(ch ta.Chart) (buy, sell ta.Condition)
	Res  int32
}

func OneTickerMultipleStrat(db *pg_wrapper.Pgx, filename string, ticker string, cfg Config, TE strategy.TradeExecution, paras strategy.Parameter, algos ...Strat) error {
	var bt []*strategy.BackTestStrategy
	var ch *ta.Kline
	for _, v := range algos {
		data, err := db.Klines("BYBIT", ticker, cfg.St, cfg.Et, v.Res)
		if err != nil {
			return err
		}
		ch = ta.NewChart(ticker, data)
		btest := strategy.NewBacktest(ch, TE, paras)
		b, s := v.Algo(ch)
		btest.AddStrategy(b, s, v.Name)
		bt = append(bt, btest.Results...)
	}

	return plot.SimplePnl(filename, nil, paras.Balance, bt)
}

type Filter struct {
	Name   string
	Filter strategy.Filter
}

func OneStrategyMultipleFilters(ch ta.Chart, filename string, cfg Config, algo strategy.AlgoFunc, TE strategy.TradeExecution, paras strategy.Parameter, indis []ta.Series, filter ...Filter) error {
	var filtered []*strategy.BackTestStrategy

	bt := strategy.NewBacktest(ch, TE, paras)
	bt.SetIndicators(indis)

	b, s := algo(ch)
	bt.AddStrategy(b, s, "")

	for _, v := range filter {
		temp := bt.Filter(v.Name, v.Filter)
		filtered = append(filtered, temp...)
	}

	return plot.DrawPnlDistributionColumn(filename, bt.Results[0].Trade(), filtered)
}
