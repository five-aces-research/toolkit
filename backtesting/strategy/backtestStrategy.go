package strategy

type BackTestStrategy struct {
	Name       string
	Parameters Parameter

	//All the trades that got executed
	tr []*Trade

	//PNL is the PNLchart, right now not implemented
	Pnl      []float64
	TotalPnl float64 //Sum of Pnl of the Trades
	Winrate  float64 //Winrate of the Trades
	AvgTrade float64 //AvgTrade Gain

	sortAlgo func(v *BackTestStrategy) float64 //sortAlgo can be Changed, default is TotalPNL
}

func (bt *BackTestStrategy) ChangeSortAlgo(fn func(b *BackTestStrategy) float64) {
	bt.sortAlgo = fn
}

func LessPnl(b *BackTestStrategy) float64 {
	return b.TotalPnl
}

func LessWinrate(b *BackTestStrategy) float64 {
	return b.Winrate
}

func LessAvgTrade(b *BackTestStrategy) float64 {
	return b.AvgTrade
}

type BackTestStrategies []*BackTestStrategy

func (t BackTestStrategies) Less(i, j int) bool {
	return t[i].sortAlgo(t[i]) < t[j].sortAlgo(t[j])
}

func (t BackTestStrategies) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t BackTestStrategies) Len() int {
	return len(t)
}
