package main

import (
	"fmt"
	"github.com/five-aces-research/toolkit/backtesting/strategy"
	"github.com/five-aces-research/toolkit/backtesting/strategy/builder"
	"github.com/five-aces-research/toolkit/backtesting/strategy/iterator"
	"github.com/five-aces-research/toolkit/backtesting/ta"
	"github.com/five-aces-research/toolkit/example/algos"
	"github.com/five-aces-research/toolkit/fas/pg_wrapper"
	"sort"
)

func SolapeIter(pg *pg_wrapper.Pgx, ticker string, cfg builder.Config, TE strategy.TradeExecution, paras strategy.Parameter) error {
	ch, err := pg.Klines("BYBIT", ticker, cfg.St, cfg.Et, cfg.Res)
	if err != nil {
		return err
	}

	kk := ta.NewChart(ticker, ch)
	o, h, c, l := kk.GetSources()

	solape := builder.NewIter(algos.Solape)

	bt := strategy.NewBacktest(kk, TE, paras)

	it := iterator.New(solape)
	it.RegisterInt(0, 4, 20, 2)
	it.RegisterInt(1, 4, 20, 2)
	it.RegisterSeries(0, o, h, c, l, ta.OC2(kk))
	it.RegisterFunctions(0, ta.Sma, ta.Ema, ta.Rma, ta.WrappedRsi)

	it.Run(bt)

	for _, v := range bt.Results {
		v.CalculatePNL()
	}

	sort.Sort(strategy.BackTestStrategies(bt.Results))

	for _, v := range bt.Results[len(bt.Results)-20:] {
		fmt.Println(v.Name, v.TotalPnl, v.Winrate)
	}

	return nil
}

func StochIter(pg *pg_wrapper.Pgx, ticker string, cfg builder.Config, TE strategy.TradeExecution, paras strategy.Parameter) error {
	ch, err := pg.Klines("BYBIT", ticker, cfg.St, cfg.Et, cfg.Res)
	if err != nil {
		return err
	}

	kk := ta.NewChart(ticker, ch)
	_, h, c, l := kk.GetSources()
	fn := algos.MfiGen(ta.Volume(kk))

	bt := strategy.NewBacktest(kk, TE, paras)

	itgirls := builder.NewIter(fn)

	it := iterator.New(itgirls)
	it.RegisterInt(0, 4, 30, 2)
	it.RegisterInt(1, 4, 30, 2)
	it.RegisterSeries(0, c, h, l, ta.HL2(kk), ta.OC2(kk))
	it.RegisterFunctions(0, ta.Sma, ta.Ema, ta.Rma, ta.Wma)

	it.Run(bt)

	for _, v := range bt.Results {
		v.CalculatePNL()
	}

	sort.Sort(strategy.BackTestStrategies(bt.Results))

	for _, v := range bt.Results[len(bt.Results)-20:] {
		fmt.Println(v.Name, len(v.Trade()), v.TotalPnl, v.Winrate)
	}

	return nil

}
