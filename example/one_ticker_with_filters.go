package main

/*
func OneTickerWithFilter(db *pg_wrapper.Pgx, ticker string, cfg AlgoConfig, strat strategy.AlgoFunc, TE strategy.TradeExecution, paras strategy.Parameter, indis ...strategy.Indicator) (*strategy.BackTest, error) {
	klines, err := db.Klines("BYBIT", ticker, cfg.St, cfg.Et, cfg.Res)
	if err != nil {
		return nil, err
	}
	ch := ta.NewChart(ticker, klines)

	var indicator []ta.Series
	for _, v := range indis {
		indicator = append(indicator, v(ch))
	}
	bt := strategy.NewBacktest(ch, TE, paras).SetIndicators(indicator)

	b, s := strat(ch)
	bt.AddStrategy(b, s, "")
	return bt, nil
}
*/
