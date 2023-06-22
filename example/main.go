package main

import (
	"log"
	"os"
	"time"

	"github.com/five-aces-research/toolkit/backtesting/strategy"
	"github.com/five-aces-research/toolkit/backtesting/strategy/builder"
	"github.com/five-aces-research/toolkit/backtesting/strategy/mode"
	"github.com/five-aces-research/toolkit/backtesting/strategy/size"
	"github.com/five-aces-research/toolkit/backtesting/strategy/te"
	"github.com/five-aces-research/toolkit/backtesting/ta"
	"github.com/five-aces-research/toolkit/example/algos"
	"github.com/five-aces-research/toolkit/fas/pg_wrapper"
)

var market = te.Market(1.0)
var parasBybit = strategy.Parameter{
	Modus:      mode.ALL,
	Pyramiding: 1,
	Fee:        &strategy.Fee{Maker: 0.0000, Taker: 0.0000, Slippage: 0},
	Balance:    10000,
	SizeType:   size.Account,
}

var tickers = []string{"i.BTCUSDT", "i.ETHUSDT", "i.SOLUSDT", "i.ADAUSDT", "i.ETCUSDT", "i.LTCUSDT", "i.XRPUSDT", "i.MATICUSDT", "i.FLMUSDT", "i.ARBUSDT", "i.STXUSDT", "i.DOGEUSDT", "i.INJUSDT", "i.FTMUSDT"}

func main() {
	db, err := pg_wrapper.Connect("127.0.0.1", "toolkit", "postgres", "password", 5432)
	if err != nil {
		os.Exit(1)
	}

	ss := algos.SolapeGenerator(5, 9, ta.Roc, false, 0.6, 0.8, 6)
	//ss := algos.KetlerChannelDivergenceSell(20, 2.0, ta.Sma, 8, 13)
	cfg := builder.Config{
		St:       Year(2021),
		Et:       time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
		Res:      1440,
		Exchange: "BYBIT",
	}

	builder.OneStratMultiTicker(db, "ketler_channel.html", tickers, cfg, ss, market, parasBybit)
	if err != nil {
		log.Panicln(err)
	}

	os.Exit(0)
}

func Year(year int) time.Time {
	return time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
}
