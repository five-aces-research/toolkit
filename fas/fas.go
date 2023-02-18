package fas

import "time"

/*
fas stands for FiveAcesStructs
*/
type OrderState int

const (
	OPEN OrderState = iota
	CANCELED
	FILLED
)

const LONG = true
const FALSE = false

type Order struct {
	Id     string // If is Int64 gets converted to string
	Side   bool   // true = buy, false = sell
	Ticker string // Ticker, Instrument, Symbol Name. Add additional Info if its not unique
	//e.g. S.BTCUSD or I.BTCUSD when there are multiple Ticker
	Size         float64 // Size in coin value
	NotionalSize float64 // Size in USD if avaible
	Price        float64 // Price
	ReduceOnly   float64
	State        OrderState // Orderstate
	Conditional  bool       // true if its conditional order
}

type Position struct {
	Id               string // If is Int64 gets converted to string
	Side             bool
	Ticker           string
	EntryPrice       float64 // Price is AvgPrice
	Size             float64 // Size in coin value
	NotionalSize     float64 // Size in USD if avaible
	LiquidationPrice float64 // If Avaible
	PNL              float64
	UPNL             float64
}

// A Candle represents the OHCLV data of a market
type Candle struct {
	Close     float64   `json:"close"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Open      float64   `json:"open"`
	Volume    float64   `json:"volume"`
	StartTime time.Time `json:"startTime"`
}

type Orderbook struct {
	Ticker    string       // Symbol name
	Bid       [][2]float64 // Bid Price, Bid Size
	Ask       [][2]float64 // Ask Price, Ask Size
	Timestamp time.Time
}

type TickerInfo struct {
	Ticker      string
	BaseCoin    string
	QuoteCoin   string
	TickSize    float64
	QtyStep     float64
	MinOrderQty float64
}

type FundingRate struct {
	Rate      float64
	Timestamp time.Time
}

type FundingPayment struct {
	Ticker    string
	Rate      float64
	Fee       float64
	Timestamp time.Time
}
