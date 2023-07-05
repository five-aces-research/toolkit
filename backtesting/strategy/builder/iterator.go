package builder

import (
	"fmt"
	"github.com/five-aces-research/toolkit/backtesting/ta"
)

//https://github.com/DawnKosmos/metapine/blob/main/backend/series/ta/backtest/iterator/iterator.go
// https://github.com/DawnKosmos/metapine/blob/main/example/iterator_example.go

type Iter struct {
	src    ta.Series
	maFast func(series ta.Series, l int) ta.Series
	l1     int
	l2     int
	fn     func(src ta.Series, maFunc ta.MaFunc, l1, l2 int) (buy, sell ta.Condition)
}

// StructsAdresse return every Pointer to Parameter(Struct Field), you want to change, in arrays
func (it *Iter) Addresses() ([]*int, []*ta.Series, []*func(src ta.Series, l int) ta.Series) {
	return []*int{&it.l1, &it.l2}, []*ta.Series{&it.src}, []*func(src ta.Series, l int) ta.Series{&it.maFast}
}

func (it *Iter) Calculation() (buy, sell ta.Condition) {
	return it.fn(it.src, it.maFast, it.l1, it.l2)
}

// Parameters returns a String of the Parameters Value, which is needed to separate the different Iterations
func (it *Iter) Parameters() string {
	//Every Indicator has its own Name(), but maFast(func(s1 Series, l int) Series) does not have it
	//Therefore we have to Wrap this function to have  a way to differentiate. See below
	return fmt.Sprintf("%s %s %d %d", ta.MaFunc(it.maFast).Name(), it.src.Name(), it.l1, it.l2)
}

func NewIter(fn func(src ta.Series, maFunc ta.MaFunc, l1, l2 int) (buy, sell ta.Condition)) *Iter {
	return &Iter{
		fn: fn,
	}
}
