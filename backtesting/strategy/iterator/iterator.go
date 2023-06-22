package iterator

import (
	"github.com/five-aces-research/toolkit/backtesting/strategy"
	"github.com/five-aces-research/toolkit/backtesting/ta"
)

type Iterator interface {
	Addresses() (lenghts []*int, srcs []*ta.Series, fns []*func(src ta.Series, l int) ta.Series) // Addresses returns pointers on the struct values, which are needed to iterate in a fast, low ram eating way
	Calculation() (buy, sell ta.Condition)                                                       // The Calculation
	Parameters() string                                                                          // Value of the Parameter Represented as string
}

type Iter struct {
	it         Iterator
	registered []iterator
	Parameter  []*int
	Srcs       []*ta.Series
	Fns        []*func(src ta.Series, l int) ta.Series
}

/*
Iter are for loops, that help you to interate your strategy with different parameters
First the Iterator interface needs to be implemented for your strategy, than parameter options have to be added
See the example folder for an easy example
*/
func New(it Iterator) *Iter {
	iter := new(Iter)
	iter.Parameter, iter.Srcs, iter.Fns = it.Addresses()
	iter.it = it
	return iter
}

type iterator interface {
	Next() bool
	Iterate()
	Reset()
}

func (it *Iter) Run(b strategy.Backtester) {
	it.run(b, it.registered)
}

func (it *Iter) run(b strategy.Backtester, iters []iterator) {
	if len(iters) == 1 {
		for iters[0].Next() {
			buy, sell := it.it.Calculation()
			b.AddStrategy(buy, sell, it.it.Parameters())
			iters[0].Iterate()
		}
	} else {
		for iters[0].Next() {
			it.run(b, iters[1:])
			iters[0].Iterate()
			iters[1].Reset()
		}
	}
}

func (it *Iter) RegisterInt(position, start, end, add int) {
	if position >= len(it.Parameter) || add == 0 {
		return
	}
	if start > end && add > 0 {
		return
	}
	if start < end && add < 0 {
		return
	}

	*it.Parameter[position] = start
	ii := &intIterator{
		val:   it.Parameter[position],
		start: start,
		end:   end,
		add:   add,
	}
	it.registered = append(it.registered, ii)
}

func (it *Iter) RegisterSeries(position int, src ...ta.Series) {
	if position >= len(it.Parameter) || len(src) == 0 {
		return
	}

	*it.Srcs[position] = src[0]
	ii := &seriesIterator{
		index: position,
		val:   it.Srcs[position],
		ss:    src,
	}
	it.registered = append(it.registered, ii)
}

func (it *Iter) RegisterFunctions(position int, fns ...func(src ta.Series, l int) ta.Series) {
	if position >= len(it.Parameter) || len(fns) == 0 {
		return
	}
	*it.Fns[position] = fns[0]
	ii := &funcIterator{
		index: position,
		val:   it.Fns[position],
		ss:    fns,
	}
	it.registered = append(it.registered, ii)
}
