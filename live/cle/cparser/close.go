package cparser

import (
	"fmt"
	"github.com/five-aces-research/toolkit/live/cle"
	"github.com/five-aces-research/toolkit/live/cle/clexer"
)

/*
Close a position

must be of form close [buy | sell] ticker
*/
type Close struct {
	Ticker string
	Side   bool
}

func ParseClose(tk []clexer.Token) (o *Close, err error) {
	o = &Close{}
	if len(tk) < 2 {
		return o, nerr(empty, "Error Close Arguments missing")
	}
	if tk[0].Type == clexer.SIDE {
		if tk[0].Value == "buy" {
			o.Side = true
		}
	}

	if tk[1].Type == clexer.VARIABLE {
		o.Ticker = tk[1].Value
	} else {
		return nil, fmt.Errorf("Close wrong Type %s", tk[1].Stringer())
	}
	return
}

func (c *Close) Evaluate(f cle.CLEIO, w Communicator) error {
	pz, err := f.Position(c.Ticker)
	if err != nil {
		return err
	}

	if pz.Side == c.Side {
		v, err := f.SetOrder(!c.Side, c.Ticker, 0, pz.NotionalSize, true, false, true)
		if err != nil {
			return err
		}
		w.Write([]byte(fmt.Sprintf("Placed Order: %s %s %f %f", v.Side, v.Ticker, v.Size, v.Price)))
	}
	return nil
}
