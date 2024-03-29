package strategy

import (
	"fmt"

	"github.com/five-aces-research/toolkit/backtesting/ta"
)

type MultiTicker struct {
	AlgoName   string
	Algo       func(ch ta.Chart) (buy ta.Condition, sell ta.Condition)
	TE         TradeExecution
	Indicators []Indicator
	Results    []*BackTestStrategy
	Parameters Parameter
}

type AlgoFunc func(ch ta.Chart) (buy ta.Condition, sell ta.Condition)

type Indicator func(ch ta.Chart) ta.Series

func NewMultiTicker(AlgoName string, Algo AlgoFunc, TE TradeExecution, parameter Parameter) *MultiTicker {
	return &MultiTicker{
		AlgoName:   AlgoName,
		Algo:       Algo,
		TE:         TE,
		Parameters: parameter,
		Indicators: make([]Indicator, 0),
	}
}

func (mt *MultiTicker) AddIndicators(indis ...Indicator) {
	mt.Indicators = append(mt.Indicators, indis...)
}

func (mt *MultiTicker) AddTickers(ch ...ta.Chart) {
	for _, v := range ch {
		bt := BackTest{
			ch:         v,
			TE:         mt.TE,
			Parameters: mt.Parameters,
		}

		//Set Indicators
		if len(mt.Indicators) != 0 {
			indicators := make([]ta.Series, 0, len(mt.Indicators))
			for _, vv := range mt.Indicators {
				indicators = append(indicators, vv(v))
			}
			bt.SetIndicators(indicators)
		}

		b, s := mt.Algo(v)
		bt.AddStrategy(b, s, fmt.Sprintf("%s %s", v.Name(), mt.AlgoName))
		mt.Results = append(mt.Results, bt.Results[0])
	}
}

func (mt *MultiTicker) Filter(ConditionName string, op Filter) []*BackTestStrategy {
	var bb []*BackTestStrategy

	for _, vv := range mt.Results {
		var tt []*Trade
		for _, v := range vv.tr {
			if op(v.Indicator) {
				tt = append(tt, v)
			}
		}
		if len(tt) > 0 {
			bb = append(bb, &BackTestStrategy{
				Name:       vv.Name + "\t" + ConditionName + " true",
				tr:         tt,
				Parameters: vv.Parameters,
				sortAlgo:   vv.sortAlgo,
			})
		}
	}
	return bb
}
