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
const SHORT = false

type Order struct {
	Id     string // If is Int64 gets converted to string
	Side   bool   // true = buy, false = sell
	Ticker string // Ticker, Instrument, Symbol Name. Add additional Info if its not unique
	//e.g. S.BTCUSD or I.BTCUSD when there are multiple Ticker
	Size         float64 // Size in coin value
	NotionalSize float64 // Size in USD if avaible
	Price        float64 // Price
	ReduceOnly   bool
	State        OrderState // Orderstate
	Conditional  bool       // true if its conditional order
	Created      time.Time
}

type Position struct {
	Id               string // If is Int64 gets converted to string
	Side             bool
	Ticker           string
	AvgPrice         float64 // Price is AvgPrice
	Size             float64 // Size in coin value
	NotionalSize     float64 // Size in USD if avaible
	LiquidationPrice float64 // If Avaible
	PNL              float64
	UPNL             float64
	Created          time.Time
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

// A Candle represents the OHCLV data of a market
type Candle struct {
	Close     float64   `json:"close"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Open      float64   `json:"open"`
	Volume    float64   `json:"volume"`
	StartTime time.Time `json:"startTime"`
}

func (c Candle) OHCL4() float64 {
	return (c.Open + c.Close + c.High + c.Low) / 4
}

type WsCandle struct {
	Data     Candle
	End      time.Time
	Ticker   string
	Finished bool
}

func (c WsCandle) ToCandle() Candle {
	return c.Data
}

/*
func (c *Client) Request(request types.RequestStruct) (*types.ResponseStruct, error) {
	var resp models.CancelOrderResponse
	err := c.POST("/v5/order/cancel", request, &respBody)
	if err != nil {
		return nil, err
	}
	return resp, nil
}


*/
