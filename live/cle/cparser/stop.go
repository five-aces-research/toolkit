package cparser

import (
	"fmt"
	"github.com/five-aces-research/toolkit/live/cle/clexer"
)

type Stop struct {
	Side   bool
	Ticker string
	A      Amount
	P      StopPrice
}

func ParseStop(tk []clexer.Token) (o *Stop, err error) {
	o = new(Stop)
	if len(tk) < 3 {
		return nil, nerr(empty, "Error Parse Stop, Input not enough arguments")
	}

	var a Amount
	if tk[0].Type != clexer.SIDE {
		return nil, nerr(empty, "Error Parse Stop, After Stop buy/sell has to be provided")
	}
	if tk[1].Type == clexer.VARIABLE {
		o.Ticker = tk[1].Value
	} else {
		return nil, nerr(empty, fmt.Sprintf("Error Parse Stop, %s is no Ticker", tk[0].Value))
	}

	switch tk[2].Type {
	case clexer.FLOAT:
		a.Type = FIAT
	case clexer.UFLOAT:
		a.Type = COIN
	case clexer.PERCENT:
		a.Type = ACCOUNT
	}
	return nil, err
}

type StopPrice struct {
	Type        PriceType
	PriceSource string
	Duration    int64
	Value       float64
}
