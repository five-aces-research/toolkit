package algos

import "github.com/five-aces-research/toolkit/backtesting/ta"

func KetlerChannelBreakoutGenerator(l1 int, mult float64, ma func(src ta.Series, l int) ta.Series) func(ta.Chart) (buy, sell ta.Condition) {
	return func(ch ta.Chart) (buy, sell ta.Condition) {
		c, l, h := ta.Close(ch), ta.Low(ch), ta.High(ch)
		upper, _, mid := ta.KetlerChannels(c, c, h, l, l1, mult, ma)
		c1, upper1 := ta.OffS(c, 1), ta.OffS(upper, 1)
		buy = ta.And(ta.Greater(c, upper), ta.SmallerEqual(c1, upper1))
		sell = ta.Greater(mid, c)

		return buy, sell
	}

}

func KetlerChannelDivergenceSell(l1 int, mult float64, ma func(src ta.Series, l int) ta.Series, l2, l3 int) func(ta.Chart) (buy, sell ta.Condition) {
	return func(ch ta.Chart) (buy, sell ta.Condition) {
		c, l, h := ta.Close(ch), ta.Low(ch), ta.High(ch)
		upper, _, _ := ta.KetlerChannels(c, c, h, l, l1, mult, ma)
		c1, upper1 := ta.OffS(c, 1), ta.OffS(upper, 1)
		buy = ta.And(ta.Greater(c, upper), ta.SmallerEqual(c1, upper1))

		oc2 := ta.OC2(ch)

		outR := ta.Sma(ta.Roc(oc2, l2), 2)
		outB1 := ma(outR, 3)
		outB2 := ma(outB1, l3)
		outB := ta.SubF(outB1, outB2, 2.0)
		cc := ta.Sub(outR, outB)
		c1 = ta.Sma(cc, 2)
		c2, c3 := ta.OffS(c1, 1), ta.OffS(c1, 2)
		sell = ta.And(ta.Smaller(c1, c2), ta.Greater(c2, c3))
		div := ta.ShortDiv(h, cc, sell)
		sell = ta.And(sell, ta.Equal(div, 1))

		return buy, sell
	}

}
