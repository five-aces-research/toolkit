package bybit

import (
	"fmt"
	"github.com/DawnKosmos/bybit-go5/models"
	"github.com/five-aces-research/toolkit/fas"
	"math"
	"strconv"
	"strings"
)

// SetOrder be aware that linear, spot execute Unit, inverse exectutes amount in USD
func (p *Private) SetOrder(side bool, Ticker string, price float64, sizeInUsd float64, marketOrder, postOnly, reduceOnly bool) (fas.Order, error) {
	Side := "Sell"
	if side {
		Side = "Buy"
	}
	OrderType := "Limit"
	var err error
	if marketOrder {
		OrderType = "Market"
		price, err = p.GetMarketPrice(Ticker)
		if err != nil {
			return fas.Order{}, err
		}
	}
	tickerInfo, err := p.GetTickerInfo(Ticker)
	if err != nil {
		return fas.Order{}, err
	}
	cat, ticker := categoryTicker(Ticker)
	var qty float64
	if cat == "inverse" {
		qty = sizeInUsd
	} else {
		qty = sizeInUsd / price
	}

	qtyAsString := roundValue(qty, tickerInfo.QtyStep)
	priceAsString := roundValue(price, tickerInfo.TickSize)

	res, err := p.by.PlaceOrder(models.PlaceOrderRequest{
		Category:   cat,
		Symbol:     ticker,
		Side:       Side,
		OrderType:  OrderType,
		Qty:        qtyAsString,
		Price:      priceAsString,
		ReduceOnly: reduceOnly,
	})
	if err != nil {
		return fas.Order{}, err
	}

	orders, err := p.by.GetOpenOrders(models.GetOpenOrdersRequest{
		Category: cat,
		Symbol:   ticker,
		OrderId:  res.OrderId,
	})
	if err != nil {
		return fas.Order{}, err
	}

	o := orders.List[0]
	side = false
	if o.Side == "Buy" {
		side = true
	}

	return apiOrderToFasOrder(o.OrderId, o.Side, Ticker, o.Qty, o.Price, o.ReduceOnly, o.OrderStatus, false), nil
}

func apiOrderToFasOrder(id string, Side string, Ticker string, Size string, Price string, ReduceOnly bool, orderStatus string, Condional bool) fas.Order {
	var side bool
	if Side == "Buy" {
		side = true
	}

	price, _ := strconv.ParseFloat(Price, 64)
	qty, _ := strconv.ParseFloat(Size, 64)
	return fas.Order{
		Id:           id,
		Side:         side,
		Ticker:       Ticker,
		Size:         qty,
		NotionalSize: qty * price,
		Price:        price,
		ReduceOnly:   ReduceOnly,
		State:        orderStatusToStatus(orderStatus),
		Conditional:  Condional,
	}

}

func roundValue(price float64, tickSize float64) string {
	rounded := math.Round(price/tickSize) * tickSize
	precision := decimalPlaces(tickSize)
	formatted := fmt.Sprintf("%.*f", precision, rounded)
	return strings.TrimRight(strings.TrimRight(formatted, "0"), ".")
}

func decimalPlaces(value float64) int {
	str := fmt.Sprintf("%f", value)
	decimals := len(str) - 1 - strings.Index(str, ".")
	return decimals
}

func orderStatusToStatus(status string) fas.OrderState {
	switch status {
	case "Created", "New":
		return fas.OPEN
	case "Rejected", "Cancelled":
		return fas.CANCELED
	default:
		return fas.FILLED
	}
}
