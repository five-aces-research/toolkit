package algos

import "github.com/five-aces-research/toolkit/backtesting/ta"

func SolapeGenerator(src func(ta.Chart) ta.Series, l1 int, l2 int, ma func(src ta.Series, l int) ta.Series, useVol bool, longPer float64, shortPer float64, lookback int) func(ta.Chart) (buy, sell ta.Condition) {

	return func(ch ta.Chart) (buy, sell ta.Condition) {
		oc2 := ta.OC2(ch)
		volume := ta.Volume(ch)
		close := ta.Close(ch)
		low := ta.Low(ch)
		high := ta.High(ch)

		outR := ta.Sma(ta.Roc(oc2, l1), 2)
		outB1 := ma(outR, l2)
		outB2 := ma(outB1, l2)
		outB := ta.SubF(outB1, outB2, 2.0)
		cc := ta.Sub(outR, outB)
		var c1 ta.Series
		if useVol {
			c1 = ta.Sma(cc, 2)
		} else {
			c1 = ta.Vwma(cc, volume, 2)
		}
		c2, c3 := ta.OffS(c1, 1), ta.OffS(c1, 2)
		buy = ta.And(ta.Greater(c1, c2), ta.Smaller(c2, c3))
		sell = ta.And(ta.Smaller(c1, c2), ta.Greater(c2, c3))

		lc := longCon(longPer, lookback, close, low)
		sc := shortCon(shortPer, lookback, close, high)
		return ta.And(buy, lc), ta.And(sell, sc)
	}
}

type MaFunc func(s ta.Series, i int) ta.Series

func (m MaFunc) Name() string {
	return m(ta.Constant([]float64{1, 2, 3, 4}, 0, 3600, ""), 2).Name()
}

func longCon(prozent float64, len int, close ta.Series, low ta.Series) ta.Condition {
	lowest := ta.Lowest(low, len)
	r1 := ta.Sub(ta.Div(close, lowest), 1)
	return ta.Greater(r1, prozent/100)
}

func shortCon(prozent float64, len int, close, high ta.Series) ta.Condition {
	highest := ta.Highest(high, len)
	r1 := ta.Sub(ta.Div(close, highest), 1)
	return ta.Smaller(r1, -prozent/100)
}

func Solape(src ta.Series, ma ta.MaFunc, l1, l2 int) (buy, sell ta.Condition) {
	outR := ta.Sma(ta.Roc(src, l1), 2)
	outB1 := ma(outR, l2)
	outB2 := ma(outB1, l2)
	outB := ta.SubF(outB1, outB2, 2.0)
	cc := ta.Sub(outR, outB)
	c1 := ta.Sma(cc, 2)
	c2, c3 := ta.OffS(c1, 1), ta.OffS(c1, 2)
	buy = ta.And(ta.Greater(c1, c2), ta.Smaller(c2, c3))
	sell = ta.And(ta.Smaller(c1, c2), ta.Greater(c2, c3))
	return
}
