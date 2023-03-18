package fas

import "time"

// Mock just is a Mock implementation of the Public/Private Interface
type Mock struct {
}

func (m *Mock) Kline(ticker string, resolution int, start time.Time, end time.Time) ([]Candle, error) {
	return []Candle{}, nil
}

func (m *Mock) GetMarketPrice(ticker string) (float64, error) {
	return 0, nil
}

func (m *Mock) GetOrderbook(ticker string, limit int) Orderbook {
	return Orderbook{}
}

func (m *Mock) GetTickerInfo(ticker string) (TickerInfo, error) {
	return TickerInfo{}, nil
}

func (m *Mock) GetFundingRate(ticker string, start, end time.Time) ([]FundingRate, error) {
	return []FundingRate{}, nil
}

func (m *Mock) GetOpenInterest(ticker string, resolution int64, start, end time.Time) ([]OpenInterest, error) {
	return []OpenInterest{}, nil
}

func (m *Mock) SetOrder(side bool, ticker string, price float64, size float64, marketOrder, postOnly, reduceOnly bool) (Order, error) {
	return Order{}, nil
}

func (m *Mock) BlockOrder(side bool, ticker string, trigger bool, priceSize [][2]float64, reduceOnly bool) ([]Order, error) {
	return []Order{}, nil
}

func (m *Mock) OpenOrders(side bool, ticker string) ([]Order, error) {
	return []Order{}, nil
}

func (m *Mock) SetTriggerOrder(side bool, ticker string, price float64, size float64, orderType string, reduceOnly bool) (Order, error) {
	return Order{}, nil
}

func (m *Mock) Cancel(Side int, Ticker string) error {
	return nil
}

func (m *Mock) CancelTrigger(Side int, Ticker string) error {
	return nil
}

func (m *Mock) Collateral(ticker string) (total float64, free float64, err error) {
	return 100, 100, nil
}

func (m *Mock) OpenPositions() ([]Position, error) {
	return []Position{}, nil
}

func (m *Mock) Position(ticker string) (Position, error) {
	return Position{}, nil
}

func (m *Mock) FundingHistory(ticker []string, start, end time.Time) ([]FundingPayment, error) {
	return []FundingPayment{}, nil
}
