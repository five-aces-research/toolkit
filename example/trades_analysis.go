package main

import "github.com/five-aces-research/toolkit/backtesting/strategy"

type TradeDistribution struct {
	Name         string
	Distribution []float64
}

func TradeAnalysis(results []*strategy.BackTestStrategy) []TradeDistribution {
	var out []TradeDistribution
	for _, vv := range results {
		tt := make([]float64, 11, 11)
		for _, t := range vv.Trade() {
			value := TradePnl(t)
			switch {
			case value < -8:
				tt[0] = tt[0] + 1
			case value >= -8 && value < -6:
				tt[1] = tt[1] + 1
			case value >= -6 && value < -4:
				tt[2] = tt[2] + 1
			case value >= -4 && value < -2:
				tt[3] = tt[3] + 1
			case value >= -2 && value < -0.5:
				tt[4] = tt[4] + 1
			case value >= -0.5 && value < 0.5:
				tt[5] = tt[5] + 1
			case value >= 0.5 && value < 2:
				tt[6] = tt[6] + 1
			case value >= 2 && value < 4:
				tt[7] = tt[7] + 1
			case value >= 4 && value < 6:
				tt[8] = tt[8] + 1
			case value >= 6 && value < 8:
				tt[9] = tt[9] + 1
			case value >= 8:
				tt[10] = tt[10] + 1
			default:
			}
		}
		out = append(out, TradeDistribution{Name: vv.Name, Distribution: tt})

	}
	return out
}

func TradePnl(t *strategy.Trade) float64 {
	var x float64
	if t.Side {
		x = (t.AvgSell - t.AvgBuy) / t.AvgBuy
	} else {
		x = -1 * (t.AvgBuy - t.AvgSell) / t.AvgBuy
	}
	return x * 100
}
