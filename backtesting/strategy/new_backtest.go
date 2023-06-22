package strategy

import (
	"time"

	"github.com/five-aces-research/toolkit/backtesting/ta"
	"github.com/five-aces-research/toolkit/fas"
)

/*
ss := strategy.New(cfg, ch)
ss.Long("id", buy, Limit,Stop, Size)
ss.Close("id", closeSignal,)
ss.Short("id", sell, Limit, Stop,poi Size)

*/

// The Public Interface implements this Interface, Though I fetching the Data before and use a DB
type Exchange interface {
	GetTickerInfo(ticker string) (fas.TickerInfo, error)                           //Needed for accurute representation of Position Size and Slippage
	GetFundingRate(ticker string, start, end time.Time) ([]fas.FundingRate, error) //Needed for to get Fundingrate
}

type NewFee struct {
	Maker    float64
	Taker    float64
	Slippage int
}

type Config struct {
	Chart      ta.Chart
	Pyramiding int
	Exchange   Exchange
	Parameter  Parameter
	Fee        *NewFee
	Balance    float64
	PnlGraph   bool
}

type NewBacktester struct {
	cfg      Config
	orderMap map[string]ta.Condition
}

func New(cfg Config) *NewBacktester {

	return &NewBacktester{}
}

func (bt *NewBacktester) Long(id string, BuyCondition ta.Condition, Limit any, Stop any, Size float64) {
	/*
		add
		Limit can be a float64, func(), Series
		Stop
		Can be a float64, func(), Series
		Size is *Size
	*/
}

func (bt *NewBacktester) Short(id string, SellCondition ta.Condition, Limit any, Stop any, Size float64) {
	// add
}

func (bt *NewBacktester) Close(IdOfPosition string, Condition ta.Condition, Size float64) {
	//add
}

func (bt *NewBacktester) Execute() {
	//Match and Add
}
