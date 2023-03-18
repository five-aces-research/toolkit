package cparser

import "github.com/five-aces-research/toolkit/live/cle"

type AmountType int

const (
	COIN AmountType = iota
	FIAT
	ACCOUNT
	POSITION
	ALL
)

type Amount struct {
	Ticker string
	Type   AmountType
	Value  float64
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
	}
	return 0, nerr(empty, "Error Evaluating Amount of Order")
}
