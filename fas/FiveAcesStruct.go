package fas

import "time"

/*
fas stands for FiveAcesStructs
*/

type Order struct {
	Id     string // If is Int64 gets converted to string
	Side   bool   // true = buy, false = sell
	Ticker string // Ticker, Instrument, Symbol Name. Add additional Info if its not unique
	//e.g. S.BTCUSD or I.BTCUSD when there are multiple Ticker with same Quote & Base pair
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
	Ticker           string    //
	AvgPrice         float64   // Price is AvgPrice
	Size             float64   // Size in coin value
	NotionalSize     float64   // Size in USD if avaible
	LiquidationPrice float64   // If Avaible
	PNL              float64   // Realised PNL
	UPNL             float64   // Unrealised PNL
	Created          time.Time // Time when order got created
}

type Orderbook struct {
	Ticker    string       // Symbol name
	Bid       [][2]float64 // Bid Price, Bid Size
	Ask       [][2]float64 // Ask Price, Ask Size
	Timestamp time.Time
}

type TickerInfo struct {
	Ticker     string  // Ticker name
	BaseCoin   string  // Base Coin e.g. BTC
	QuoteCoin  string  // Base Coin e.g. USDT
	SettleCoin string  // Coin thats used for fees and funding fees
	TickSize   float64 // The TickSize describes the smallest price increments in the Ticker price
	//e.g. 0.5 means that the price is increased in 0.5 steps. 20000 -> 20000,5 -> 20001
	//this can change over time
	QtyStep float64 // The QtyStep gets info of the smallest size increment in the Ticker
	// e.g 1 means e.g. the Buy Size has to be 5,6,7 ...
	MinOrderQty float64 // MinOrderQty, The Minimum Order  Quantity to buy an asset, usually equally to QtyStep
}

type FundingRate struct {
	Ticker    string  // Ticker name
	Rate      float64 // Rate in %
	Timestamp time.Time
}

type FundingPayment struct {
	Ticker     string  //Ticker name
	Rate       float64 // Rate in %
	SettleCoin string  // The Coin that is used for the fee
	Fee        float64 // Actually paid amount in the Settle Currency
	Timestamp  time.Time
}

type Coin struct {
	Coin     string  `json:"coin"`
	Symbol   string  `json:"symbol"`
	Equity   float64 `json:"equity"`
	UsdValue float64 `json:"usd_value"`
}

type Wallet struct {
	TotalEquity float64 // Total Equity in USD
	FreeEquity  float64 // Total Free Equity in USD
	Coins       []Coin  // Your Wallet Balance
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

// A Transfer describes a transaction of a coin from 1 to another address/account
type Transfer struct {
	TransferId string       // TransferId used internally in the Exchange
	Type       TransferType // The Type, Deposit, WithDraw or Account-To-Account Transfer
	Created    time.Time
	Coin       string  // The Coin Used
	Amount     float64 // The Amount used

	Chain string  // The Chain for the Transfer, use exchange or empty if its account transfer
	Fee   float64 // Fee Paid
	TxId  string  // TxId if its done on exchange
	From  string  // Accountname if its internal, else public adress
	To    string  // Accountname if its internal, else public adress
}

type FeeRate struct {
	Ticker     string  // The Ticker
	SettleCoin string  // The Coin that is used to burn the Fee
	Maker      float64 // in %
	Taker      float64 // in %
}

type Ticker struct {
	Ticker string
	Bid1   [2]float64 // Bid Price, Bid Size
	Ask1   [2]float64 // Ask Price, Ask Size
	Last   [2]float64 // Last Price, Last Size
}

// OHCLV4 calculates the avg of all 4 Price points
func (c Candle) OHCL4() float64 {
	return (c.Open + c.Close + c.High + c.Low) / 4
}

type WsCandle struct {
	Data     Candle    // OHCLV Data
	End      time.Time // EndTime of the Candle
	Ticker   string    // Tickername
	Finished bool      // Finished just Closed
}

func (c WsCandle) ToCandle() Candle {
	return c.Data
}
