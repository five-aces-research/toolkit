package lta

import (
	"github.com/five-aces-research/toolkit/fas"
	"log"
	"time"
)

type Chart interface {
	V(index int) fas.Candle
	Data() []fas.Candle
	SetLimit(limit int)
	ExecuteLimit()
	GetUpdateGroup() UpdateGroup
	ResolutionStartTime
}

type KLINE struct {
	URS[fas.Candle]
	stream chan fas.WsCandle
	start  chan struct{}
	name   string
	res    int64
	st     int64
	ug     []Updater
	limit  int
}

func (o *KLINE) Add(u Updater) {
	o.ug = append(o.ug, u)
}

func Kline(data fas.Streamer, ticker string, resolution int64, limit int64, parameters ...any) *KLINE {
	//Implementiert auch Updater, must be high
	o := new(KLINE)
	o.name = ticker
	end := time.Now()
	start := end.Add(-time.Minute * time.Duration(resolution*limit))
	ch, err := data.Kline(ticker, resolution, start, end)
	if err != nil {
		log.Panicln(err)
		return nil
	}
	o.data = Array(ch)

	o.stream, err = data.LiveKline(ticker, resolution, parameters...)
	o.start = make(chan struct{}) //Create Start Channel
	if err != nil {
		log.Panicln(err)
		return nil
	}

	go o.receive()

	return o
}

func (o *KLINE) receive() {
	<-o.start
	o.ExecuteLimit()
	for _, v := range o.ug {
		v.ExecuteLimit()
	}
	log.Println("starting", o.name)
	for v := range o.stream {
		if v.Finished {
			o.data.SetValue(0, v.Data)
			o.data.Append(newCandle(v.Data.Close, v.Data.StartTime, o.res))
		} else {
			o.data.SetValue(0, v.Data)
		}
		for _, vv := range o.ug {
			vv.OnTick(v.Finished)
		}
	}
}

func (o *KLINE) Start() {
	o.start <- struct{}{}
}

func (o *KLINE) V(index int) fas.Candle {
	return o.data.V(index)
}

func (o *KLINE) Data() []fas.Candle {
	return o.data.Data()

}

func (o *KLINE) StartTime() int64 {
	return o.st
}

func (o *KLINE) Resolution() int64 {
	return o.res
}

func (o *KLINE) Name() string {
	return o.name
}

func (o *KLINE) SetName(s string) {
	o.name = s
}

func (o *KLINE) GetUpdateGroup() UpdateGroup {
	return o
}

func (o *KLINE) SetLimit(limit int) {
	if limit > o.limit {
		limit = o.limit
	}
}

func (o *KLINE) ExecuteLimit() {
	o.data.SetLimit(o.limit)
}

func newCandle(c float64, st time.Time, res int64) fas.Candle {
	return fas.Candle{
		Close:     c,
		High:      c,
		Low:       c,
		Open:      c,
		Volume:    0,
		StartTime: st.Add(time.Duration(res) * time.Minute),
	}
}
