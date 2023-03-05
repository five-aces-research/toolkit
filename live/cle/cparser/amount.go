package cparser

type AmountType int

const (
	COIN AmountType = iota
	FIAT
	ACCOUNTSIZE
	POSITIONSIZE
	ALL
)

type Amount struct {
	Ticker string
	Type   AmountType
	Value  float64
}
