package algos

import "github.com/five-aces-research/toolkit/backtesting/ta"

func SaphirGenerator(src1, src2 func(ta.Chart) ta.Series, l1, l2 int, longPer float64, shortPer float64, lookback int) func(ta.Chart) (buy, sell ta.Condition) {
	return func(ch ta.Chart) (buy, sell ta.Condition) {
		s1, s2 := src1(ch), src2(ch)

		outR1 := ta.Sma(ta.Rsi(s1, l1), 2)
		outR2 := ta.Sma(ta.Rsi(s2, l1), 2)
		outR := ta.Div(ta.AddF(outR2, outR1, 2), 3)
		outB1 := ta.Sma(outR, l2)
		outB2 := ta.Sma(outB1, l2)
		outB := ta.SubF(outB1, outB2, 2)

		cc := ta.Sub(outR, outB)
		c := ta.Sma(cc, 2)

		c1, c2 := ta.OffS(c, 1), ta.OffS(c, 2)

		buy = ta.And(ta.Greater(c, c1), ta.Smaller(c1, c2))
		sell = ta.And(ta.Smaller(c, c1), ta.Greater(c1, c2))
		close, low, high := ta.Close(ch), ta.Low(ch), ta.High(ch)
		lc := longCon(longPer, lookback, close, low)
		sc := shortCon(shortPer, lookback, close, high)
		return ta.And(buy, lc), ta.And(sell, sc)
	}
}

func DefaultSaphir() func(ta.Chart) (buy, sell ta.Condition) {
	return SaphirGenerator(ta.OC2, ta.HL2, 12, 4, 0.4, 0.8, 4)
}
