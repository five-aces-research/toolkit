package deribit

import (
	"github.com/five-aces-research/toolkit/fas"
	"github.com/frankrap/deribit-api/models"
	"strings"
	"time"
)

func (p *Private) buyMarket(ticker string, qtySize float64, reduceOnly bool) (fas.Order, error) {
	o, err := p.d.Buy(&models.BuyParams{
		InstrumentName: ticker,
		Amount:         qtySize,
		Type:           "market",
		ReduceOnly:     reduceOnly,
	})
	if err != nil {
		return fas.Order{}, err
	}
	return deribitOrderToOrder(o.Order), nil
}

func (p *Private) buyLimit(ticker string, qtyPrice, qtySize float64, reduceOnly, postOnly bool) (fas.Order, error) {
	o, err := p.d.Buy(&models.BuyParams{
		InstrumentName: ticker,
		Amount:         qtySize,
		Type:           "limit",
		Price:          qtyPrice,
		PostOnly:       postOnly,
		ReduceOnly:     reduceOnly,
	})
	if err != nil {
		return fas.Order{}, err
	}
	return deribitOrderToOrder(o.Order), nil
}

func (p *Private) sellMarket(ticker string, qtySize float64, reduceOnly bool) (fas.Order, error) {
	o, err := p.d.Sell(&models.SellParams{
		InstrumentName: ticker,
		Amount:         qtySize,
		Type:           "market",
		ReduceOnly:     reduceOnly,
	})
	if err != nil {
		return fas.Order{}, err
	}
	return deribitOrderToOrder(o.Order), nil
}

func (p *Private) sellLimit(ticker string, qtyPrice, qtySize float64, reduceOnly, postOnly bool) (fas.Order, error) {
	o, err := p.d.Sell(&models.SellParams{
		InstrumentName: ticker,
		Amount:         qtySize,
		Type:           "limit",
		Price:          qtyPrice,
		PostOnly:       postOnly,
		ReduceOnly:     reduceOnly,
	})
	if err != nil {
		return fas.Order{}, err
	}
	return deribitOrderToOrder(o.Order), nil
}

func deribitOrderToOrder(o models.Order) fas.Order {
	var os fas.OrderState
	//Order state: "open", "filled", "rejected", "cancelled", "untriggered"
	switch o.OrderState {
	case "open", "untriggered":
		os = fas.OPEN
	case "filled":
		os = fas.FILLED
	case "rejected", "cancelled":
		os = fas.CANCELED
	}
	var Side bool
	if o.Direction == "buy" {
		Side = true
	}

	//Order type: "limit", "market", "stop_limit", "stop_market"
	var conditionalOrder bool
	if strings.HasPrefix(o.OrderType, "stop") {
		conditionalOrder = true
	}

	return fas.Order{
		Id:           o.OrderID,
		Side:         Side,
		Ticker:       o.InstrumentName,
		Size:         o.Amount / o.Price.ToFloat64(),
		NotionalSize: o.Amount,
		Price:        o.Price.ToFloat64(),
		ReduceOnly:   o.ReduceOnly,
		State:        os,
		Conditional:  conditionalOrder,
		Created:      time.Unix(o.CreationTimestamp/1000, 0),
	}
}

func (p *Private) buyTrigger(ticker string, price float64, size float64, reduceOnly bool) (fas.Order, error) {
	o, err := p.d.Buy(&models.BuyParams{
		InstrumentName: ticker,
		Amount:         size,
		Type:           "stop_market",
		ReduceOnly:     reduceOnly,
		Trigger:        "mark_price",
		StopPrice:      price,
	})
	if err != nil {
		return fas.Order{}, err
	}
	return deribitOrderToTriggerOrder(o.Order), nil
}

func (p *Private) sellTrigger(ticker string, price float64, size float64, reduceOnly bool) (fas.Order, error) {
	o, err := p.d.Sell(&models.SellParams{
		InstrumentName: ticker,
		Amount:         size,
		Type:           "stop_market",
		ReduceOnly:     reduceOnly,
		Trigger:        "mark_price",
		StopPrice:      price,
	})
	if err != nil {
		return fas.Order{}, err
	}
	return deribitOrderToTriggerOrder(o.Order), nil
}

func deribitOrderToTriggerOrder(o models.Order) fas.Order {
	var os fas.OrderState
	//Order state: "open", "filled", "rejected", "cancelled", "untriggered"
	switch o.OrderState {
	case "open", "untriggered":
		os = fas.OPEN
	case "filled":
		os = fas.FILLED
	case "rejected", "cancelled":
		os = fas.CANCELED
	}
	var Side bool
	if o.Direction == "buy" {
		Side = true
	}

	return fas.Order{
		Id:           o.OrderID,
		Side:         Side,
		Ticker:       o.InstrumentName,
		Size:         o.Amount / o.StopPrice,
		NotionalSize: o.Amount,
		Price:        o.StopPrice,
		ReduceOnly:   o.ReduceOnly,
		State:        os,
		Conditional:  true,
		Created:      time.Unix(o.CreationTimestamp/1000, 0),
	}
}
