package example

import (
	"fmt"
	"testing"
	"time"

	"github.com/five-aces-research/toolkit/fas/pg_wrapper"
	"github.com/stretchr/testify/assert"
)

func TestKlines(t *testing.T) {

	pg, err := pg_wrapper.Connect("127.0.0.1", "toolkit", "postgres", "password", 5432)
	assert.Nil(t, err)
	assert.Nil(t, pg.Ping())
	tnow := time.Now()

	ch, err := pg.KlinesNew("BYBIT", "i.BTCUSDT", Year(2022), tnow, 240)

	assert.Nil(t, err)
	assert.Condition(t, func() bool {
		return ch != nil || len(ch) > 0
	})

	for _, v := range ch {
		fmt.Println(v)
	}

}

func Year(year int) time.Time {
	return time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
}
