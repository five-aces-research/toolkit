package bybit

import (
	"github.com/DawnKosmos/bybit-go5/models"
	"github.com/DawnKosmos/bybit-go5/ws"
	"github.com/five-aces-research/toolkit/fas"
	"strconv"
	"time"
)

type Streamer struct {
	wss *ws.Stream
	pp  *Public
}

func (s *Streamer) Kline(ticker string, resolution int64, start time.Time, end time.Time) ([]fas.Candle, error) {
	return s.pp.Kline(ticker, resolution, start, end)
}

func (s *Streamer) LiveKline(ticker string, resolution int64, parameters ...any) (chan fas.WsCandle, error) {
	ch := make(chan fas.WsCandle)
	_, ticker = categoryTicker(ticker)

	err := s.wss.Kline(ticker, resolutionToString(resolution), func(e *models.WsKline) {
		for _, v := range e.Data {
			o, _ := strconv.ParseFloat(v.Open, 64)
			h, _ := strconv.ParseFloat(v.High, 64)
			c, _ := strconv.ParseFloat(v.Close, 64)
			l, _ := strconv.ParseFloat(v.Low, 64)
			vol, _ := strconv.ParseFloat(v.Volume, 64)

			ch <- fas.WsCandle{
				Data: fas.Candle{
					Close:     c,
					High:      h,
					Low:       l,
					Open:      o,
					Volume:    vol,
					StartTime: time.Unix(v.Start/1000, 0),
				},
				End:      time.Unix(v.End/1000, 0),
				Ticker:   e.Topic,
				Finished: v.Confirm,
			}
		}
	})
	return ch, err
}

func (s *Streamer) Ping() error {
	return nil
}

func NewStreamer(category ws.WsLink) *Streamer {
	pp := NewPublic()
	wss := ws.New(ws.Config{
		Endpoint:      category,
		A:             nil,
		AutoReconnect: true,
		Debug:         false,
		TestNet:       false,
	})
	return &Streamer{
		wss: wss,
		pp:  pp,
	}
}
