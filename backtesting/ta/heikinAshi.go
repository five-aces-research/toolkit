package ta

import (
	"github.com/five-aces-research/toolkit/fas"
	"math"
)

type heikinAshi struct {
	ERS[fas.Candle]
}

func HeikinAshi(src *Kline) *heikinAshi {
	var s heikinAshi
	s.st = src.StartTime()
	s.res = src.Resolution()
	s.name = "HA of " + s.Name()

	f := src.Data()
	d := make([]fas.Candle, 0, len(f)+1)
	var c fas.Candle

	c.Open = src.ch[0].Open
	c.Close = src.ch[0].OHCL4()
	c.Low = math.Min(src.ch[0].Open, c.Open)
	c.High = math.Max(src.ch[0].High, c.Open)
	c.Volume = src.ch[0].Volume
	d = append(d, c)
	for _, v := range src.ch[1:] {
		c.Open = (c.Open + c.Close) / 2
		c.Close = v.OHCL4()
		c.Low = math.Min(v.Open, c.Open)
		c.High = math.Min(v.Open, c.Open)
		c.Volume = v.Volume
		c.StartTime = v.StartTime
		d = append(d, c)
	}

	s.data = d
	return &s
}
