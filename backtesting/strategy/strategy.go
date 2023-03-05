package strategy

import (
	"github.com/five-aces-research/toolkit/backtesting/strategy/mode"
	"github.com/five-aces-research/toolkit/backtesting/strategy/size"
	"github.com/five-aces-research/toolkit/backtesting/ta"
)

/*
the strategy library can be used for simple Backtest like on Tradingview
The features come in after you have these Results, because you are able to Do Even more Analysis
*/

const LONG = true
const SHORT = false

type Filterer interface {
	Split(condition string, op Filter)
	Filter(condition string, op Filter)
}

type Backtester interface {
	AddStrategy(buy, sell ta.Condition, IdName string)
}

/*
SafeFloat
Pretty often indicator have no Value while others have it. To still store everything in an Array we have this struct.
SafeFloat are used for filters in Strategies
*/
type SafeFloat struct {
	Safe  bool
	Value float64
}

type Parameter struct {
	//Modus, means OnlyShort, onlyLongs or ALL
	Modus mode.Mode
	//On Default pyramiding is one
	Pyramiding int
	//Fees are described by Maker(Market Orders) Taker(limit orders) and Slippage(market orders)
	Fee *Fee
	//The Account Balance
	Balance float64
	//Size is either Dollar or AccountSize
	SizeType size.SizeBase
	//right now PnlGraph isnt implemented. It will be useful in visual representation
	PnlGraph bool
}

type Fee struct {
	Maker    float64
	Taker    float64
	Slippage float64
}
