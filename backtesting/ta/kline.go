package ta

import (
	"github.com/five-aces-research/toolkit/fas"
	"log"
)

type Chart interface {
	Data() []fas.Candle
	ResolutionStartTime
}

type Kline struct {
	ch   []fas.Candle
	st   int64
	res  int64
	name string
}

// NewChart returns a Chart, returns 0 if error occured
func NewChart(ticker string, ch []fas.Candle) *Kline {
	o := new(Kline)
	if len(ch) == 0 {
		log.Println("error empty input")
		return nil
	}
	o.st = ch[0].StartTime.Unix()
	o.res = ch[1].StartTime.Unix() - o.st
	o.ch = ch
	o.name = ticker
	return o
}

func (o *Kline) SetName(name string) {
	o.name = name
}

func (o *Kline) Data() []fas.Candle {
	return o.ch
}

func (o *Kline) StartTime() int64 {
	return o.st
}

func (o *Kline) Resolution() int64 {
	return o.res
}

func (o *Kline) Name() string {
	return o.name
}
