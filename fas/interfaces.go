package fas

import "time"

type Exchanger interface {
	//Name returns the Exchange Name and Subaccount Name if available
	Name() string
	Public
	Private
}

type OpenInterest struct {
	OI        float64
	Timestamp time.Time
}

// Public interface are all function that are needed to grab data from an exchange
type Public interface {
	//Name returns the Name of the exchange
	Name() string
	//Kline returns candlestickdata, ordered ascending, resolution in minutes
	Kline(ticker string, resolution int64, start time.Time, end time.Time) ([]Candle, error)
	//MarketPrice return the Market Price of the asked Ticker
	GetMarketPrice(ticker string) (float64, error)
	//GetOrderbook returns the orderbook
	GetOrderbook(ticker string, limit int) (Orderbook, error)
	//GetTicker return ticker/instrument information
	GetTickerInfo(ticker string) (TickerInfo, error)
	//GetFundingRate, ordered ascending
	GetFundingRate(ticker string, start, end time.Time) ([]FundingRate, error)
	//GetOpenInterest Returns open interest
	GetOpenInterest(ticker string, resolution int64, start, end time.Time) ([]OpenInterest, error)
}

// Private interface is a collection of functions that are used for all cases where authentification is needed
type Private interface {
	//SetOrder sets an Order. Returns the set order. Size always in USD if available
	SetOrder(side bool, ticker string, price float64, size float64, marketOrder, postOnly, reduceOnly bool) (Order, error)
	//OpenOrders All=0, Buy=1 Sell=-1 Returns open orders for given ticker, "" return all open orders
	OpenOrders(side int, ticker string) ([]Order, error)
	//SetTriggerOrder set an TriggerOrder
	SetTriggerOrder(side bool, ticker string, price float64, size float64, orderType string, reduceOnly bool) (Order, error)
	//Cancel All=0, Buy=1 Sell=-1 orders on given ticker. No ticker means all orders get cancelled. Return is the amount of orders that got cancelled
	Cancel(Side int, Ticker string) error
	//CancelTrigger All=0, Buy=1 Sell=-1 orders on given ticker. No ticker means all orders get cancelled. Return is the amount of orders that got cancelled
	CancelTrigger(Side int, Ticker string) error
	//Collateral returns the amount of free collatal in USD terms
	Collateral(ticker string) (total float64, free float64, err error)
	//OpenPositions returns all Open positions
	OpenPositions() ([]Position, error)
	//Position return given Position nil, if does not exist
	Position(ticker string) (*Position, error)
	//FundingHistory No Ticker or nil equal all Coins
	FundingHistory(ticker []string, start, end time.Time) ([]FundingPayment, error)
	//GetOrderHistory Returns Historical Orders. Orders that got cancelled filled; No Ticker or nil equal all Coins
	GetOrderHistory(ticker []string, start, end time.Time) ([]Order, error)
	//AccountInformation returns the Total Wallet Balance and a Sum of all Coins owned
	AccountInformation() (Wallet, error)
	//GetTransfers returns the latest transfers done by the account
	GetTransfers(ticker string, st time.Time, et time.Time, OptionalType ...TransferType) ([]Transfer, error)
	//GetFeeRate get the trading fee rate
	GetFeeRate(ticker ...string) ([]FeeRate, error)
	BlockOrder(side bool, ticker string, b bool, prices [][2]float64, only bool) ([]Order, error)
}

type Streamer interface {
	Kline(ticker string, resolution int64, start time.Time, end time.Time) ([]Candle, error)
	LiveKline(ticker string, resolution int64, parameters ...any) (chan WsCandle, error)
	Ticker(ticker string) (chan Ticker, error)
	Ping() error
}
