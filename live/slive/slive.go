package slive

import (
	"github.com/five-aces-research/toolkit/backtesting/ta"
	"github.com/five-aces-research/toolkit/fas"
	"log"
	"sync"
	"time"
)

// simple live  trading, calculates the new Indicator every tick.
// simple live needs also some way to check the calculations and ticker

var DEFAULT = 200

type Strategy struct {
	Ticker           string
	Resolution       int64         // in minutes
	LookBack         time.Duration //nil uses default which is resolution*200
	Signal           func(o ta.Chart) (bool, bool)
	Buy              func(o ta.Chart) error
	Sell             func(o ta.Chart) error
	ReceiverFunction func(o []fas.Candle, c fas.WsCandle) // You have to use SetIsLooked(true) to  use this, remember to SetIsLooked(false) after used. else you do to many unnecessary calculations

	mu         sync.RWMutex
	isLookedAt bool
}

type LiveTrade struct {
	ee  fas.Private
	ews fas.Streamer
}

func New(ee fas.Private, ews fas.Streamer) (*LiveTrade, error) {
	lt := new(LiveTrade)

	// Check Private Valid.
	// Check Streamer

	lt.ee = ee
	lt.ews = ews

	return lt, nil
}

func (lt *LiveTrade) AddStrategy(str *Strategy, parameters ...any) error {
	tick, err := lt.ews.LiveKline(str.Ticker, str.Resolution, parameters)
	if err != nil {
		return nil
	}
	if str.LookBack == 0 {
		str.LookBack = time.Second * time.Duration(str.Resolution) * 200
	}

	start := time.Now().Add(-str.LookBack)
	ch, err := lt.ews.Kline(str.Ticker, str.Resolution, start, time.Now())

	go execFunc(lt.ee, ch, tick, str)
	return nil
}

func execFunc(ep fas.Private, ch []fas.Candle, tick chan fas.WsCandle, str *Strategy) {
	for {
		c := <-tick
		if !c.Finished {
			if str.isLooked() {
				str.ReceiverFunction(ch, c)
			}
			continue
		}
		ch = append(ch, c.ToCandle())
		o := ta.NewChart(c.Ticker, ch)
		buy, sell := str.Signal(o)

		if buy {
			err := str.Buy(o)
			if err != nil {
				log.Println(err)
			}
		}
		if sell {
			err := str.Sell(o)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (s *Strategy) SetIsLooked(b bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isLookedAt = b
}

func (s *Strategy) isLooked() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isLookedAt
}
