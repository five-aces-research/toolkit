package cparser

import (
	"errors"
	"fmt"
	"github.com/five-aces-research/toolkit/live/cle"
	"github.com/five-aces-research/toolkit/live/cle/clexer"
)

/*
	Cancel cancels orders

syntax: cancel [buy | sell] ?-stop ticker... the
The Order doesn't matter
So cancel buy -stop
or cancel -stop buy is equal
if no Side get provided, all get canceled
*/
type Cancel struct {
	Side    int      // Side -1 Sell 0 All 1 buy
	Ticker  []string // array of all tickers
	trigger bool     // stops
}

func ParseCancel(tk []clexer.Token) (c *Cancel, err error) {
	c = new(Cancel)
	if len(tk) == 0 {
		return c, errors.New("parse error: empty input")
	}

	for _, v := range tk {
		switch v.Type {
		case clexer.SIDE:
			if v.Value == "buy" {
				c.Side = 1
			} else {
				c.Side = -1
			}
		case clexer.VARIABLE:
			c.Ticker = append(c.Ticker, v.Value)
		case clexer.FLAG:
			if v.Value == "stop" {
				c.trigger = true
			}
		default:
			return nil, fmt.Errorf("type %s not supported", v.Stringer())
		}
	}
	if len(c.Ticker) == 0 {
		return nil, errors.New("error no ticker provided")
	}

	return
}

func (c *Cancel) Evaluate(f cle.CLEIO, w Communicator) (errr error) {
	for _, v := range c.Ticker {
		if c.trigger {
			err := f.CancelTrigger(c.Side, v)
			if err != nil {
				errr = fmt.Errorf("%s | %s", errr.Error(), err)
			} else {
				w.Write([]byte(fmt.Sprintf("%s Trigger Orders  cancelled succesfully", v)))
			}
		} else {
			err := f.Cancel(c.Side, v)
			if err != nil {
				errr = fmt.Errorf("%s | %s", errr.Error(), err)
			} else {
				w.Write([]byte(fmt.Sprintf("%s Orders  cancelled succesfully", v)))
			}
		}
	}
	return errr
}
