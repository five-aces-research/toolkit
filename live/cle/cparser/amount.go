package cparser

import "github.com/five-aces-research/toolkit/live/cle"

type AmountType int

const (
	COIN AmountType = iota
	FIAT
	ACCOUNT
	POSITION
	POSITIONORDER
)

type Amount struct {
	Ticker  string
	Type    AmountType
	Value   float64
	side    bool
	trigger float64
}

func (a Amount) GetAmount(f cle.CLEIO) (float64, error) {
	switch a.Type {
	case FIAT:
		return a.Value, nil
	case COIN:
		mp, err := f.GetMarketPrice(a.Ticker)
		return a.Value * mp, err
	case ACCOUNT:
		_, free, err := f.Collateral(a.Ticker)
		return free * a.Value / 100, err
	case POSITION:
		p, err := f.Position(a.Ticker)
		return p.NotionalSize, err
	case POSITIONORDER:
		p, err := f.Position(a.Ticker)
		if err != nil {
			return 0, err
		}
		totalSize := p.NotionalSize
		var side int = 1
		if a.side {
			side = -1
		}

		or, err := f.OpenOrders(side, a.Ticker)
		for _, v := range or {
			totalSize += v.NotionalSize
		}

		return totalSize, err
	}
	return 0, nerr(empty, "Error Evaluating Amount of Order")
}
