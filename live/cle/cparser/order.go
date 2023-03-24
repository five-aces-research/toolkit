package cparser

import (
	"fmt"
	"github.com/five-aces-research/toolkit/live/cle"
	"github.com/five-aces-research/toolkit/live/cle/clexer"
	"strconv"
)

type Order struct {
	Side   bool
	Ticker string
	A      Amount
	P      Price
}

// ParseOrder parses Orders
func ParseOrder(Side string, tk []clexer.Token) (o *Order, err error) {
	o = &Order{}

	if Side == "buy" {
		o.Side = true
	}
	var a Amount

	if len(tk) < 3 {
		return nil, nerr(empty, "Error Parse Order, False Input")
	}

	if tk[0].Type == clexer.VARIABLE {
		o.Ticker = tk[0].Value
	} else {
		return nil, nerr(empty, fmt.Sprintf("Error no Ticker is %s", tk[0].Value))
	}

	switch tk[1].Type {
	case clexer.FLOAT: // 5.2 -> 5.2 Coins
		a.Type = FIAT
	case clexer.UFLOAT: // u500 -> 500 USD worth of the coin
		a.Type = COIN
	case clexer.PERCENT: // 100% -> 100% of your Free Collateral of the Coin
		a.Type = ACCOUNT
	case clexer.POSITION: // -position -> 100% of the Positions Size
		a.Type = POSITION
	default:
		return nil, nerr(empty, fmt.Sprintf("Error Parse Order, false Order Size of type"))
	}
	a.Ticker = o.Ticker
	a.Value, err = strconv.ParseFloat(tk[1].Value, 64)
	if err != nil {
		return nil, nerr(err, fmt.Sprintf("Parse Error Wrong Value should be a Float is %s", tk[1].Value))
	}

	o.A = a

	o.P, err = ParsePrice(tk[2:])
	return o, err
}

func (o *Order) Evaluate(f cle.CLEIO, w Communicator) error {
	size, err := o.A.GetAmount(f)
	if err != nil {
		return err
	}
	return o.P.Execute(f, w, o.Side, o.Ticker, size)
}
