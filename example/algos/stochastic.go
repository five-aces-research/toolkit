package algos

import (
	"github.com/five-aces-research/toolkit/backtesting/ta"
)

func StochGen(h, c, l ta.Series, lenStoch, lenSignal int) func(src ta.Series, ma ta.MaFunc, l1, l2 int) (buy, sell ta.Condition) {
	lSmooth := 3
	return func(src ta.Series, ma ta.MaFunc, l1, l2 int) (ta.Condition, ta.Condition) {
		hS1, lS1 := ta.Ema(h, lSmooth), ta.Ema(l, lSmooth)
		hS := ta.SubF(hS1, ta.Ema(hS1, lSmooth), 2)
		lS := ta.SubF(lS1, ta.Ema(lS1, lSmooth), 2)
		stoch := ta.Stoch(c, hS, lS, lenStoch)

		//fmt.Println(stoch.Data())

		perD := ma(stoch, lenSignal)
		slow := ma(perD, l1)

		signal1 := ta.Ema(slow, l2)
		signal := ta.SubF(signal1, ta.Ema(signal1, l2), 2)

		buy := ta.Crossover(slow, signal)
		sell := ta.Crossunder(slow, signal)
		return buy, sell
	}
}

func MfiGen(vol ta.Series) func(src ta.Series, ma ta.MaFunc, l1, l2 int) (buy, sell ta.Condition) {
	return func(src ta.Series, ma ta.MaFunc, l1, l2 int) (buy, sell ta.Condition) {
		mfi := ta.MFI(src, vol, l1)
		s := ta.Sma(mfi, 2)
		s2 := ma(s, l2)
		ss := ta.SubF(s2, ma(s2, l2), 2)

		buy = ta.Crossover(s, ss)
		sell = ta.Crossunder(s, ss)
		return buy, sell
	}
}
