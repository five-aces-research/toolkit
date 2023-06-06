package deribit

import (
	"fmt"
	"github.com/five-aces-research/toolkit/fas"
	"github.com/frankrap/deribit-api"
	"github.com/frankrap/deribit-api/models"
	"log"
	"strings"
	"time"
)

type Private struct {
	Public
}

func (p *Private) GetOrderHistory(ticker []string, start, end time.Time) ([]fas.Order, error) {
	//TODO implement me
	panic("there is no api endpoint for this")
}

func (p *Private) AccountInformation() (fas.Wallet, error) {
	var c fas.Wallet

	//BTC
	or, err := p.d.GetAccountSummary(&models.GetAccountSummaryParams{
		Currency: "BTC",
		Extended: false,
	})
	if err != nil {
		return fas.Wallet{}, err
	}

	res, err := p.d.GetIndex(&models.GetIndexParams{Currency: "BTC"})
	if err != nil {
		return c, err
	}
	usdVal := or.Equity * res.BTC
	c.TotalEquity += usdVal
	c.FreeEquity += or.AvailableFunds * res.BTC

	c.Coins = append(c.Coins, fas.Coin{
		Coin:     "BITCOIN",
		Symbol:   or.Currency,
		Equity:   or.Equity,
		UsdValue: usdVal,
	})
	//ETH
	or, err = p.d.GetAccountSummary(&models.GetAccountSummaryParams{
		Currency: "ETH",
		Extended: false,
	})
	if err != nil {
		return c, err
	}

	res, err = p.d.GetIndex(&models.GetIndexParams{Currency: "ETH"})
	if err != nil {
		return c, err
	}
	usdVal = or.Equity * res.ETH
	c.TotalEquity += usdVal
	c.FreeEquity += or.AvailableFunds * res.ETH

	c.Coins = append(c.Coins, fas.Coin{
		Coin:     "ETHEREUM",
		Symbol:   or.Currency,
		Equity:   or.Equity,
		UsdValue: usdVal,
	})

	//ETH
	or, err = p.d.GetAccountSummary(&models.GetAccountSummaryParams{
		Currency: "USDC",
		Extended: false,
	})
	if err != nil {
		return c, err
	}

	c.TotalEquity += or.Equity
	c.FreeEquity += or.AvailableFunds

	c.Coins = append(c.Coins, fas.Coin{
		Coin:     "USDC",
		Symbol:   "USDC",
		Equity:   or.AvailableFunds,
		UsdValue: or.AvailableFunds,
	})

	return c, nil
}

func (p *Private) GetTransfers(ticker string, st time.Time, et time.Time, OptionalType ...fas.TransferType) ([]fas.Transfer, error) {
	//TODO implement me
	panic("there is no api endpoint for this")
}

func (p *Private) GetFeeRate(ticker ...string) ([]fas.FeeRate, error) {
	//TODO implement me
	panic("there is no api endpoint for this")
}

func NewPrivate(name string, public string, private string, testnet bool) *Private {
	url := deribit.RealBaseURL
	if testnet {
		url = deribit.TestBaseURL
	}
	var np Private
	np.d = deribit.New(&deribit.Configuration{
		Addr:          url,
		ApiKey:        public,
		SecretKey:     private,
		AutoReconnect: true,
		DebugMode:     false,
	})
	np.tickerInfo = make(map[string]fas.TickerInfo)

	return &np
}

func (p *Private) SetOrder(side bool, Ticker string, price float64, sizeInUsd float64, marketOrder, postOnly, reduceOnly bool) (fas.Order, error) {
	ticker := toPerpetual(Ticker)

	tickerInfo, err := p.GetTickerInfo(ticker)
	if err != nil {
		return fas.Order{}, err
	}

	if tickerInfo.SettleCoin == "USDC" || tickerInfo.SettleCoin == "" {
		sizeInUsd = sizeInUsd / price
	}

	qtySize := fas.RoundValue(sizeInUsd, tickerInfo.QtyStep)

	if marketOrder {
		if side {
			return p.buyMarket(ticker, qtySize, reduceOnly)
		} else {
			return p.sellMarket(ticker, qtySize, reduceOnly)
		}
	}

	qtyPrice := fas.RoundValue(price, tickerInfo.TickSize)

	fmt.Println(qtyPrice, qtySize)
	if side {
		return p.buyLimit(ticker, qtyPrice, qtySize, reduceOnly, postOnly)
	} else {
		return p.sellLimit(ticker, qtyPrice, qtySize, reduceOnly, postOnly)
	}
}

func (p *Private) BlockOrder(side bool, ticker string, trigger bool, priceSize [][2]float64, reduceOnly bool) ([]fas.Order, error) {
	var out []fas.Order

	for _, v := range priceSize {
		o, err := p.SetOrder(side, ticker, v[0], v[1], false, true, reduceOnly)
		if err != nil {
			log.Println(err)
		} else {
			out = append(out, o)
		}
	}
	return out, nil
}

func (p *Private) OpenOrders(side bool, ticker string) ([]fas.Order, error) {
	ticker = toPerpetual(ticker)
	var o []fas.Order
	or, err := p.d.GetOpenOrdersByInstrument(&models.GetOpenOrdersByInstrumentParams{
		InstrumentName: ticker,
		Type:           "limit",
	})
	if err != nil {
		return nil, err
	}

	Side := "sell"
	if side {
		Side = "buy"
	}

	for _, v := range or {
		if v.Direction == Side {
			o = append(o, deribitOrderToOrder(v))
		}
	}

	return o, nil
}

func (p *Private) SetTriggerOrder(side bool, Ticker string, price float64, sizeInUsd float64, orderType string, reduceOnly bool) (fas.Order, error) {
	ticker := toPerpetual(Ticker)

	tickerInfo, err := p.GetTickerInfo(ticker)
	if err != nil {
		return fas.Order{}, err
	}

	if tickerInfo.SettleCoin == "USDC" || tickerInfo.SettleCoin == "" {
		sizeInUsd = sizeInUsd / price
	}

	qtySize := fas.RoundValue(sizeInUsd, tickerInfo.QtyStep)

	qtyPrice := fas.RoundValue(price, tickerInfo.TickSize)

	fmt.Println(qtyPrice, qtySize)
	if side {
		return p.buyTrigger(ticker, qtyPrice, qtySize, reduceOnly)
	} else {
		return p.sellTrigger(ticker, qtyPrice, qtySize, reduceOnly)
	}
}

func (p *Private) Cancel(Side int, ticker string) error {
	ticker = toPerpetual(ticker)
	if Side == 0 {
		_, err := p.d.CancelAllByInstrument(&models.CancelAllByInstrumentParams{
			InstrumentName: ticker,
			Type:           "limit",
		})
		return err
	}

	var condition func(string) bool = func(s string) bool {
		return s == "sell"
	}
	if Side > 0 {
		condition = func(s string) bool {
			return s == "buy"
		}
	}

	or, err := p.d.GetOpenOrdersByInstrument(&models.GetOpenOrdersByInstrumentParams{
		InstrumentName: ticker,
		Type:           "limit",
	})
	if err != nil {
		return err
	}
	go func() {
		for _, v := range or {
			if condition(v.Direction) {
				_, err := p.d.Cancel(&models.CancelParams{OrderID: v.OrderID})
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
	return nil
}

func (p *Private) CancelTrigger(Side int, ticker string) error {
	ticker = toPerpetual(ticker)
	if Side == 0 {
		_, err := p.d.CancelAllByInstrument(&models.CancelAllByInstrumentParams{
			InstrumentName: ticker,
			Type:           "trigger_all",
		})
		return err
	}

	var condition func(string) bool = func(s string) bool {
		return s == "sell"
	}
	if Side > 0 {
		condition = func(s string) bool {
			return s == "buy"
		}
	}

	or, err := p.d.GetOpenOrdersByInstrument(&models.GetOpenOrdersByInstrumentParams{
		InstrumentName: ticker,
		Type:           "trigger_all",
	})
	if err != nil {
		return err
	}
	go func() {
		for _, v := range or {
			if condition(v.Direction) {
				_, err := p.d.Cancel(&models.CancelParams{OrderID: v.OrderID})
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
	return nil
}

func (p *Private) Collateral(ticker string) (total float64, free float64, err error) {
	ticker = strings.ToUpper(ticker)

	var currency string
	var price float64
	switch {
	case strings.HasPrefix(ticker, "USDC"):
		currency = "USDC"
		price = 1
	case strings.HasSuffix(ticker, "BTC"):
		currency = "BTC"
		price, err = p.GetMarketPrice("BTC")
		if err != nil {
			return 0, 0, err
		}
	case strings.HasSuffix(ticker, "ETH"):
		currency = "ETH"
		price, err = p.GetMarketPrice("ETH")
		if err != nil {
			return 0, 0, err
		}
	default:
		return 0, 0, fmt.Errorf("ticker %s not known", ticker)
	}

	res, err := p.d.GetAccountSummary(&models.GetAccountSummaryParams{
		Currency: currency,
	})

	return res.Balance * price, res.AvailableFunds * price, nil
}

func (p *Private) OpenPositions() ([]fas.Position, error) {
	var pos []fas.Position

	pBtc, err := p.d.GetPositions(&models.GetPositionsParams{
		Currency: "BTC",
	})
	if err != nil {
		return pos, err
	}
	for _, v := range pBtc {
		pos = append(pos, deribitPositiontoPosition(v))
	}
	pEth, err := p.d.GetPositions(&models.GetPositionsParams{
		Currency: "ETH",
	})
	if err != nil {
		return pos, err
	}
	for _, v := range pEth {
		pos = append(pos, deribitPositiontoPosition(v))
	}
	pUsdc, err := p.d.GetPositions(&models.GetPositionsParams{
		Currency: "USDC",
	})
	if err != nil {
		return pos, err
	}
	for _, v := range pUsdc {
		pos = append(pos, deribitPositiontoPosition(v))
	}
	return pos, nil
}

func (p *Private) Position(ticker string) (*fas.Position, error) {
	ticker = toPerpetual(ticker)

	res, err := p.d.GetPosition(&models.GetPositionParams{InstrumentName: ticker})
	if err != nil {
		return nil, err
	}
	if res.Size == 0 {
		return nil, nil
	}
	pos := deribitPositiontoPosition(res)
	return &pos, nil
}

func (p *Private) FundingHistory(ticker []string, start, end time.Time) ([]fas.FundingPayment, error) {
	//TODO rewrite with pagination
	return nil, nil
}

func deribitPositiontoPosition(p models.Position) fas.Position {
	Side := p.Direction == "buy"

	return fas.Position{
		Id:               p.InstrumentName + p.Kind,
		Side:             Side,
		Ticker:           p.InstrumentName,
		AvgPrice:         p.AveragePrice,
		Size:             p.SizeCurrency,
		NotionalSize:     p.Size,
		LiquidationPrice: p.EstimatedLiquidationPrice,
		PNL:              p.RealizedProfitLoss,
		UPNL:             p.TotalProfitLoss,
		Created:          time.Now(),
	}
}
