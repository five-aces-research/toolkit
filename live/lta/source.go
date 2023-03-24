package lta

import (
	"github.com/five-aces-research/toolkit/backtesting/ta"
	"github.com/five-aces-research/toolkit/fas"
)

type source struct {
	URS[float64]
	src Chart
	fn  func(candle fas.Candle) float64
}

func Source(o Chart, op func(candle fas.Candle) float64, name string) Series {
	s := new(source)
	s.src = o

	src := ta.Source(o, op, name)
	s.name = name
	s.st, s.res = src.StartTime(), src.Resolution()

	s.data = Array(src.Data())

	s.ug = o.GetUpdateGroup()
	s.ug.Add(s)

	return s
}

func (s *source) OnTick(new bool) {
	if new {
		src0, src1 := s.src.V(0), s.src.V(1)
		s.data.SetValue(0, s.fn(src1))
		s.data.Append(s.fn(src0))
	} else {
		src0 := s.src.V(0)
		s.data.SetValue(0, s.fn(src0))
	}
}

func (s *source) SetLimit(limit int) {
	if limit > s.limit {
		s.limit = limit
		s.src.SetLimit(s.limit)
	}
}

func Open(c Chart) Series {
	fn := func(e fas.Candle) float64 {
		return e.Open
	}
	return Source(c, fn, "Open")
}

func Close(c Chart) Series {
	fn := func(e fas.Candle) float64 {
		return e.Close
	}
	return Source(c, fn, "Close")
}

func High(c Chart) Series {
	fn := func(e fas.Candle) float64 {
		return e.High
	}
	return Source(c, fn, "High")
}

func Low(c Chart) Series {
	fn := func(e fas.Candle) float64 {
		return e.Low
	}
	return Source(c, fn, "Low")
}

func Volume(c Chart) Series {
	fn := func(e fas.Candle) float64 {
		return e.Volume
	}
	return Source(c, fn, "Volume")
}
