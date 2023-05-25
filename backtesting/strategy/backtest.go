package strategy

import (
	"github.com/five-aces-research/toolkit/backtesting/strategy/mode"
	"github.com/five-aces-research/toolkit/backtesting/strategy/size"
	"github.com/five-aces-research/toolkit/backtesting/ta"
	"sort"
)

// BackTest is what it says a struct to backtest
type BackTest struct {
	ch ta.Chart //Kline Data needed for calculation
	//TE describes how Trades get executed. You can choose simple market orders, or set multiple different limit orders or set stops
	TE         TradeExecution
	Parameters Parameter //Parameters of the strategy
	//Indicators are saved in a [][]SafeFloat Array and synched to the lenght of the OHCLV data
	Indicators [][]SafeFloat
	Results    []*BackTestStrategy // The Results of a Strategy Used on  a Chart
}

// NewBacktest returns the Backtest Struct, you have to change
func NewBacktest(ch ta.Chart, TE TradeExecution, p Parameter) *BackTest {
	return &BackTest{
		ch:         ch,
		TE:         TE,
		Parameters: p,
	}
}

// AddStrategy , Strategies are created with buy and sell signals(just bool) They open a Position and close the  contrarian one.
func (bt *BackTest) AddStrategy(buy, sell ta.Condition, Name string) {
	var b = new(BackTestStrategy)
	//OHCLV, buy and sell have to match up the same size
	ch, l, s := bt.ch.Data(), buy.Data(), sell.Data()
	sl, _ := ta.MinInt(len(ch), len(l), len(s))
	ch = ch[len(ch)-sl:]
	l = l[len(l)-sl:]
	s = s[len(s)-sl:]
	b.sortAlgo = LessPnl

	//Check if PnlGraph is supported
	if bt.TE.GetInfo().CandlePnlSupport && bt.Parameters.PnlGraph {
		b.Pnl = make([]float64, len(ch), len(ch))
		bt.Parameters.PnlGraph = false
	}

	//Check if Indicators were added
	var indicators [][]SafeFloat
	if bt.Indicators != nil {
		indicators = bt.Indicators[len(bt.Indicators)-sl:]
	} else {
		indicators = make([][]SafeFloat, sl, sl)
	}
	var indexLong, indexShort []int
	var tr []*Trade

	balance := bt.Parameters.Balance
	parameters := bt.Parameters
	var tempBalance float64

	//Trades get Created here, this is a Simple Backtest. It does not support, having buy and sell strategies simultan next to each other.
	for j := 0; j < len(ch)-1; j++ {
		if l[j] {
			for i := 0; i < min(len(indexShort), parameters.Pyramiding); i++ {
				index := indexShort[i]
				t, err := bt.TE.CreateTrade(SHORT, ch[index+1:], j-index, indicators[index], balance, *bt.Parameters.Fee, b.Parameters.PnlGraph)
				if err != nil {
					continue
				}
				tr = append(tr, t)
				if parameters.SizeType == size.Account {
					tempBalance += t.RealisedPNL()
				}
			}
			if parameters.SizeType == size.Account {
				balance += tempBalance
				tempBalance = 0
			}

			indexShort = indexShort[:0]
			if parameters.Modus != mode.OnlySHORT {
				indexLong = append(indexLong, j)
			}
		}
		if s[j] {
			for i := 0; i < min(len(indexLong), parameters.Pyramiding); i++ {
				index := indexLong[i]
				t, err := bt.TE.CreateTrade(LONG, ch[index+1:], j-index, indicators[index], balance, *parameters.Fee, b.Parameters.PnlGraph)
				if err != nil {
					//fmt.Println("Create Longs at", j, err)
					continue
				}
				tr = append(tr, t)
				if parameters.SizeType == size.Account {
					tempBalance += t.RealisedPNL()
				}
			}
			if parameters.SizeType == size.Account {
				balance += tempBalance
				tempBalance = 0
			}

			indexLong = indexLong[:0]
			if parameters.Modus != mode.OnlyLONG {
				indexShort = append(indexShort, j)
			}
		}
	}

	sort.Sort(Trades(tr)) // Sort Trades by EntrySignalTime
	b.tr = tr
	b.Name = Name
	bt.Results = append(bt.Results, b)
}

/*
SetIndicator, fills the [][]SafeFloat with Series, Close etc is also a Series
Be aware that indicators can only be set before strategies get added.
*/
func (bt *BackTest) SetIndicators(indicators []ta.Series) *BackTest {
	if len(indicators) == 0 {
		return bt
	}

	d := bt.ch.Data()
	indi := make([][]SafeFloat, 0, len(d)) //init array
	f := indicators[0].Data()              // Data of first array
	l1 := len(indicators)

	for i := 0; i < len(d)-len(f); i++ {
		init := make([]SafeFloat, l1, l1) //init t
		indi = append(indi, init)
	}

	var j int = len(d) - len(f)
	for _, v := range f {
		init := make([]SafeFloat, l1, l1)
		init[0] = SafeFloat{Safe: true, Value: v}
		indi = append(indi, init)
		j++
	}

	var i int = 1
	for _, vv := range indicators[1:] {
		f = vv.Data()
		j = len(d) - len(f)
		for _, v := range f {
			indi[j][i] = SafeFloat{Safe: true, Value: v}
			j++
		}
		i++
	}
	bt.Indicators = indi
	return bt
}
