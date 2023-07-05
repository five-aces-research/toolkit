package main

import (
	"fmt"
	"github.com/five-aces-research/toolkit/backtesting/ta"
	"github.com/five-aces-research/toolkit/example/algos"
	"os"
	"time"

	"github.com/five-aces-research/toolkit/backtesting/strategy"
	"github.com/five-aces-research/toolkit/backtesting/strategy/builder"
	"github.com/five-aces-research/toolkit/backtesting/strategy/mode"
	"github.com/five-aces-research/toolkit/backtesting/strategy/size"
	"github.com/five-aces-research/toolkit/backtesting/strategy/te"
	"github.com/five-aces-research/toolkit/fas/pg_wrapper"
)

var market = te.Market(1.0)
var parasBybit = strategy.Parameter{
	Modus:      mode.ALL,
	Pyramiding: 1,
	Fee:        &strategy.Fee{Maker: 0.0000, Taker: 0.0000, Slippage: 0},
	Balance:    10000,
	SizeType:   size.Dollar,
}

var tickers = []string{"l.BTCUSDT", "l.ETHUSDT", "l.SOLUSDT", "l.ADAUSDT", "l.ETCUSDT", "l.LTCUSDT", "l.XRPUSDT", "l.MATICUSDT", "l.FLMUSDT", "l.ARBUSDT", "l.STXUSDT", "l.DOGEUSDT", "l.INJUSDT", "l.FTMUSDT"}

func main() {
	db, err := pg_wrapper.Connect("127.0.0.1", "toolkit", "postgres", "password", 5432)
	if err != nil {
		os.Exit(1)
	}
	alg1 := algos.SolapeGenerator(ta.Open, 10, 12, ta.Sma, true, 0.0, 0.0, 4)
	alg2 := algos.SolapeGenerator(ta.OC2, 16, 6, ta.Sma, true, 0.0, 0.0, 4)
	alg3 := algos.SolapeGenerator(ta.OC2, 4, 6, ta.Sma, false, 0.0, 0.0, 4)
	d1 := algos.SolapeGenerator(ta.High, 20, 10, ta.Sma, true, 0.0, 0.0, 4)
	d2 := algos.SolapeGenerator(ta.High, 20, 4, ta.Rma, true, 0.0, 0.0, 4)

	//ss := algos.KetlerChannelDivergenceSell(20, 2.0, ta.Sma, 8, 13)
	cfg := builder.Config{
		St:       Year(2020),
		Et:       time.Date(2023, 6, 20, 0, 0, 0, 0, time.UTC),
		Res:      1440,
		Exchange: "BYBIT",
	}

	bs := []builder.Strat{{
		Name: "4h, Open, 10, 12, SMA",
		Algo: alg1,
		Res:  240,
	}, {
		Name: "4h, OC2, 16, 6, SMA",
		Algo: alg2,
		Res:  240,
	}, {Name: "4h, OC2, 4, 6, SMA",
		Algo: alg3,
		Res:  240,
	}, {Name: "D, High, 20, 10, SMA",
		Algo: d1,
		Res:  1440,
	}, {
		Name: "D, High, 20, 4, SMA",
		Algo: d2,
		Res:  1440,
	},
	}
	fmt.Println(len(bs))
	/*
		ch, err := db.Klines("BYBIT", "l.ETHUSDT", cfg.St, cfg.Et, cfg.Res)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		k := ta.NewChart("", ch)
		_, _, c, _ := k.GetSources()

		macd, signal, hist := ta.MacdRelative(c, 12, 26, 9)

		var indis []ta.Series = []ta.Series{ta.Rsi(c, 14), macd, signal, hist, c, ta.Sma(c, 20)}

		filters := []builder.Filter{{
			Name:   "RSI < 30",
			Filter: filter.SmallerAs(0, 30),
		}, {
			Name:   "MACD < -2",
			Filter: filter.SmallerAs(1, -2),
		}, {
			Name:   "MACD > 2",
			Filter: filter.GreaterAs(1, 2),
		},
		}
	*/
	//err = builder.OneTickerMultipleStrat(db, "solape/all_eth_with_vol.html", "l.ETHUSDT", cfg, te.Limit(0.2), parasBybit, bs...)
	//err = builder.OneStrategyMultipleFilters(k, "solape/eth_filter.html", cfg, alg3, te.Limit(3), parasBybit, indis, filters...)
	err = StochIter(db, "L.BTCUSDT", cfg, te.Limit(3), parasBybit)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func Year(year int) time.Time {
	return time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
}
