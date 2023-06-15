package cparser

import (
	"fmt"
	"github.com/five-aces-research/toolkit/live/cle"
	"github.com/five-aces-research/toolkit/live/cle/clexer"
	"strconv"
)

type Stop struct {
	Side   bool
	Ticker string
	A      Amount
	P      StopPrice
}

//	stop buy btc-perp [-o, -p]

func ParseStop(tk []clexer.Token) (o *Stop, err error) {
	o = new(Stop)
	if len(tk) < 3 {
		return nil, nerr(empty, "Error Parse Stop, Input not enough arguments")
	}

	if tk[0].Type != clexer.SIDE {
		return nil, nerr(empty, "Error Parse Stop, After Stop buy/sell has to be provided")
	}
	o.Side = tk[0].Value == "buy"
	if tk[1].Type == clexer.VARIABLE {
		o.Ticker = tk[1].Value
	} else {
		return nil, nerr(empty, fmt.Sprintf("Error Parse Stop, %s is no Ticker", tk[0].Value))
	}

	var a Amount
	a.Ticker = a.Ticker
	a.side = o.Side

	switch tk[2].Type {
	case clexer.FLOAT:
		a.Type = FIAT
	case clexer.UFLOAT:
		a.Type = COIN
	case clexer.POSITION:
		a.Type = POSITION
	case clexer.POSITIONORDER:
		a.Type = POSITIONORDER
	default:
		return nil, nerr(empty, fmt.Sprintf("Error Parse Stop, false Stop Size of type %v", tk[2].Value))
	}
	a.Value, err = strconv.ParseFloat(tk[2].Value, 64)
	if err != nil {
		return nil, nerr(err, fmt.Sprintf("Parse Error Wrong Value should be a Float is %s", tk[1].Value))
	}

	var p StopPrice
	switch tk[3].Type {
	case clexer.FLOAT: // 30000 places order at $30000
		p.Type = PRICE
	case clexer.DFLOAT: // -300 places order $300 below/above the marketprice
		p.Type = DIFFERENCE
	case clexer.PERCENT: // 2% places order 2% above/below the marketprice
		p.Type = PERCENTPRICE
	default:
		return nil, nerr(empty, fmt.Sprintf("Error Price Parsing, %v %s is not a valid price", tk[0].Type, tk[0].Value))
	}

	p.Value, err = strconv.ParseFloat(tk[3].Value, 64)
	if err != nil {
		return nil, nerr(err, fmt.Sprintf("Error Price Parsing %s is not a Number", tk[0].Value))
	}
	o.A = a
	o.P = p

	return o, err
}

func (s *Stop) Evaluate(f cle.CLEIO, w Communicator) error {
	size, err := s.A.GetAmount(f)
	if err != nil {
		return err
	}
	return s.P.Execute(f, w, s.Side, s.Ticker, size)
}

type StopPrice struct {
	Type  PriceType
	Value float64
}

func (s *StopPrice) Execute(f cle.CLEIO, w Communicator, side bool, ticker string, size float64) error {
	var price float64
	factor := 1.0 // Factor has to be opposite!
	if side {
		factor = -1.0
	}

	switch s.Type {
	case PRICE:
		price = s.Value
	case DIFFERENCE:
		mp, err := f.GetMarketPrice(ticker)
		if err != nil {
			return nil
		}
		price = mp - s.Value*factor
	case PERCENTPRICE:
		mp, err := f.GetMarketPrice(ticker)
		if err != nil {
			return nil
		}
		price = mp - mp*s.Value/100*factor
	}

	fmt.Println(side, ticker, price, size, "trigger", true)

	or, err := f.SetTriggerOrder(side, ticker, price, size, "trigger", true)
	if err == nil {
		w.Write([]byte(fmt.Sprintf("placed stop order %s @ %2.f %2.f ", or.Ticker, or.Price, or.Size)))
	}
	return err
}
