package bybit

import (
	"errors"
	"fmt"
	"github.com/DawnKosmos/bybit-go5"
	"github.com/DawnKosmos/bybit-go5/models"
	"github.com/five-aces-research/toolkit/fas"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Private struct {
	Public
}

func (p *Private) AccountInformation() (fas.Wallet, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Private) GetTransfers(ticker string, st time.Time, et time.Time, OptionalType ...fas.TransferType) ([]fas.Transfer, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Private) GetFeeRate(ticker ...string) ([]fas.FeeRate, error) {
	//TODO implement me
	panic("implement me")
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
		o, err := p.SetOrder(side, Ticker, v[0], v[1], false, true, reduceOnly)
		if err != nil {
			log.Println(err)
		} else {
			out = append(out, o)
		}
	}
	return out, nil
}

func (p *Private) OpenOrders(side int, Ticker string) ([]fas.Order, error) {
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

	filter := func(OrderSide string) bool {
		return true
	}
	if side >= 1 {
		filter = func(s string) bool {
			return s == "Buy"
		}
	} else if side <= -1 {
		filter = func(s string) bool {
			return s == "Sell"
		}
	}

	for _, o := range res.List {
		if filter(o.Side) {
			oo := apiOrderToFasOrder(o.OrderId, o.Side, Ticker, o.Qty, o.Price, o.ReduceOnly, o.OrderStatus, false)
			out = append(out, oo)
		}
	}

	return out, nil
}

func (p *Private) SetTriggerOrder(side bool, ticker string, price float64, size float64, orderType string, reduceOnly bool) (fas.Order, error) {
	cat, ticker := categoryTicker(ticker)
	fmt.Println(cat)
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
			return ss == "Sell"
		}
	case 1:
		condition = func(ss string) bool {
			return ss == "Buy"
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
		Coin:        "BTC", // BTC is choose, so that Bybit API is not giving us a list of all coins.
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

func (p *Private) GetOrderHistory(ticker []string, start, end time.Time) ([]fas.Order, error) {
	var or []fas.Order

	if ticker == nil || len(ticker) == 0 {
		res, err := p.getOrderHistory("spot", "", start, end)
		if err != nil {
			return or, err
		}
		or = append(or, res...)
		res, err = p.getOrderHistory("inverse", "", start, end)
		if err != nil {
			return or, err
		}
		or = append(or, res...)
		res, err = p.getOrderHistory("linear", "", start, end)
		if err != nil {
			return or, err
		}
		or = append(or, res...)
	} else {
		for _, v := range ticker {
			cat, tt := categoryTicker(v)
			res, err := p.getOrderHistory(cat, tt, start, end)
			if err != nil {
				return or, err
			}
			or = append(or, res...)
		}
	}
	sort.Sort(Orders(or))
	return or, nil
}

func (p *Private) getOrderHistory(category string, ticker string, start, end time.Time) ([]fas.Order, error) {
	var cursor string
	var or []fas.Order

	res, err := p.by.GetOrderHistory(models.GetOrderHistoryRequest{
		Category: category,
		Symbol:   ticker,
		Limit:    50,
	})
	if err != nil {
		return nil, err
	}
	cursor = res.NextPageCursor

	or = toFasOrderFromHistoryOrder(category, res)
	for cursor != "" {
		res, err := p.by.GetOrderHistory(models.GetOrderHistoryRequest{
			Category: category,
			Symbol:   ticker,
			Limit:    50,
			Cursor:   cursor,
		})
		if err != nil {
			return nil, err
		}
		cursor = res.NextPageCursor

		or = append(or, toFasOrderFromHistoryOrder(category, res)...)
	}

	return or, nil
}

func toFasOrderFromHistoryOrder(category string, request *models.GetOrderHistoryResponse) []fas.Order {
	or := make([]fas.Order, 0, len(request.List))
	catIdentifier := category[:1] + "."

	for _, v := range request.List {
		Side := v.Side == "buy"
		Size, _ := strconv.ParseFloat(v.Qty, 64)
		Price, _ := strconv.ParseFloat(v.Price, 64)
		CreateTime, _ := strconv.ParseInt(v.CreatedTime, 10, 64)

		or = append(or, fas.Order{ //TODO maybe adder trigger price
			Id:           v.OrderId,
			Side:         Side,
			Ticker:       strings.ToUpper(catIdentifier + v.Symbol),
			Size:         Size,
			NotionalSize: Price * Size,
			Price:        Price,
			ReduceOnly:   v.ReduceOnly,
			State:        orderStatusToStatus(v.OrderStatus),
			Conditional:  v.CloseOnTrigger,
			Created:      time.Unix(CreateTime/1000, 0),
		})
	}
	return or
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

type Orders []fas.Order

func (o Orders) Len() int {
	return len(o)
}

func (o Orders) Less(i, j int) bool {
	return o[i].Created.Unix() < o[j].Created.Unix()
}

func (o Orders) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}
