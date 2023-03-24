package bybit

import (
	"errors"
	"fmt"
	"github.com/DawnKosmos/bybit-go5"
	"github.com/DawnKosmos/bybit-go5/models"
	"github.com/five-aces-research/toolkit/fas"
	"log"
	"strconv"
	"time"
)

type Private struct {
	Public
}

func NewPrivate(name string, public string, private string, test bool) *Private {
	link := bybit.URL
	if test {
		link = bybit.TESTURL
	}

	bb, err := bybit.New(nil, link, &bybit.Account{
		PublicKey: public,
		SecretKey: private,
	}, false)
	if err != nil {
		log.Panicln(err)
	}
	var np Private
	np.by = bb
	np.tickerInfo = make(map[string]fas.TickerInfo)
	np.cache = Cache{data: make(map[string]CacheEntry)}
	return &np
}

func (p *Private) BlockOrder(side bool, Ticker string, trigger bool, priceSize [][2]float64, reduceOnly bool) ([]fas.Order, error) {
	var out []fas.Order

	for _, v := range priceSize {
		o, err := p.SetOrder(side, Ticker, v[0], v[1], false, true, false)
		if err != nil {
			log.Println(err)
		} else {
			out = append(out, o)
		}
	}
	return out, nil
}

func (p *Private) OpenOrders(side bool, Ticker string) ([]fas.Order, error) {
	cat, ticker := categoryTicker(Ticker)

	res, err := p.by.GetOpenOrders(models.GetOpenOrdersRequest{
		Category:    cat,
		Symbol:      ticker,
		OpenOnly:    0,
		OrderFilter: "Order",
	})
	if err != nil {
		return nil, err
	}

	var out []fas.Order

	for _, o := range res.List {
		var Side bool
		if o.Side == "Buy" {
			Side = true
		}

		if Side == side {
			oo := apiOrderToFasOrder(o.OrderId, o.Side, Ticker, o.Qty, o.Price, o.ReduceOnly, o.OrderStatus, false)
			out = append(out, oo)
		}
	}

	return out, nil
}

func (p *Private) SetTriggerOrder(side bool, ticker string, price float64, size float64, orderType string, reduceOnly bool) (fas.Order, error) {
	cat, ticker := categoryTicker(ticker)
	fmt.Println(cat, ticker)

	panic("implement me")
}

func (p *Private) Cancel(Side int, Ticker string) error {
	cat, ticker := categoryTicker(Ticker)

	res, err := p.by.GetOpenOrders(models.GetOpenOrdersRequest{
		Category:    cat,
		Symbol:      ticker,
		OpenOnly:    0,
		OrderFilter: "Order",
	})
	if err != nil {
		return err
	}

	var condition func(string) bool

	switch Side {
	case 0:
		condition = func(string) bool { return true }
	case -1:
		condition = func(ss string) bool {
			if ss == "Sell" {
				return true
			}
			return false
		}
	case 1:
		condition = func(ss string) bool {
			if ss == "Buy" {
				return true
			}
			return false
		}
	}

	go func() {
		for _, v := range res.List {
			if condition(v.Side) {
				_, err := p.by.CancelOrder(models.CancelOrderRequest{
					Category: cat,
					Symbol:   ticker,
					OrderId:  v.OrderId,
				})
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
	return nil
}

func (p *Private) CancelTrigger(Side int, ticker string) error {
	cat, ticker := categoryTicker(ticker)

	res, err := p.by.GetOpenOrders(models.GetOpenOrdersRequest{
		Category:    cat,
		Symbol:      ticker,
		OpenOnly:    0,
		OrderFilter: "StopOrder",
	})
	if err != nil {
		return err
	}

	var condition func(string) bool

	switch Side {
	case 0:
		condition = func(string) bool { return true }
	case -1:
		condition = func(ss string) bool {
			if ss == "Sell" {
				return true
			}
			return false
		}
	case 1:
		condition = func(ss string) bool {
			if ss == "Buy" {
				return true
			}
			return false
		}
	}

	go func() {
		for _, v := range res.List {
			if condition(v.Side) {
				_, err := p.by.CancelOrder(models.CancelOrderRequest{
					Category: cat,
					Symbol:   ticker,
					OrderId:  v.OrderId,
				})
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
	return nil
}

func (p *Private) Collateral(Ticker string) (total float64, free float64, err error) {
	cat, ticker := categoryTicker(Ticker)
	if "inverse" == cat {
		coin := ticker[:3]
		if coin == "MAN" {
			coin = "MANA"
		}
		res, err := p.by.GetWalletBalance(models.GetWalletBalanceRequest{
			AccountType: "CONTRACT",
			Coin:        coin,
		})
		if err != nil {
			return 0, 0, err
		}

		l := res.List[0].Coin[0]

		total, _ = strconv.ParseFloat(l.Equity, 64)
		free, _ = strconv.ParseFloat(l.AvailableToWithdraw, 64)

		if total == 0 {
			return 0, 0, err
		}
		price, err := p.GetMarketPrice(Ticker)
		if err != nil {
			return 0, 0, err
		}
		return price * total, price * free, nil
	}

	res, err := p.by.GetWalletBalance(models.GetWalletBalanceRequest{
		AccountType: "UNIFIED",
		Coin:        "BTC",
	})
	if err != nil {
		return 0, 0, err
	}

	l := res.List[0]
	total, _ = strconv.ParseFloat(l.TotalEquity, 64)
	free, _ = strconv.ParseFloat(l.TotalAvailableBalance, 64)

	return total, free, nil
}

func (p *Private) OpenPositions() ([]fas.Position, error) {
	panic("impomenet")
}

func (p *Private) Position(Ticker string) (*fas.Position, error) {
	cat, ticker := categoryTicker(Ticker)
	res, err := p.by.GetPositionInfo(models.GetPositionInfoRequest{
		Category: cat,
		Symbol:   ticker,
	})
	if err != nil {
		return nil, err
	}

	if len(res.List) == 0 {
		return nil, nil
	}
	pos := res.List[0]
	if cat != "invers" {
		return toFasPosition(pos.Side, Ticker, pos.Size, pos.PositionValue, pos.AvgPrice, pos.UnrealisedPnl, pos.LiqPrice, pos.CreatedTime), nil
	} else {
		return toFasPosition(pos.Side, Ticker, pos.PositionValue, pos.Size, pos.AvgPrice, pos.UnrealisedPnl, pos.LiqPrice, pos.CreatedTime), nil
	}

}

func (p *Private) FundingHistory(ticker []string, start, end time.Time) ([]fas.FundingPayment, error) {
	return nil, errors.New("not implemented")
}

func toFasPosition(Side string, Ticker string, size, notionalSize, avgPrice, pnl, liqPrice string, created string) *fas.Position {
	var side bool
	if Side == "Buy" {
		side = true
	}

	Size, _ := strconv.ParseFloat(size, 64)
	NotionalSize, _ := strconv.ParseFloat(notionalSize, 64)
	AvgPrice, _ := strconv.ParseFloat(avgPrice, 64)
	Pnl, _ := strconv.ParseFloat(pnl, 64)
	LiqPrice, _ := strconv.ParseFloat(liqPrice, 64)
	Created, _ := strconv.ParseInt(created, 10, 64)

	return &fas.Position{
		Side:             side,
		Ticker:           Ticker,
		AvgPrice:         AvgPrice,
		Size:             Size,
		NotionalSize:     NotionalSize,
		LiquidationPrice: LiqPrice,
		UPNL:             Pnl,
		Created:          time.Unix(Created/1000, 0),
	}
}
