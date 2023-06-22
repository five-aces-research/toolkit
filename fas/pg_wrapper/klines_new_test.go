package pg_wrapper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestKlines(t *testing.T) {
	pg, err := Connect("127.0.0.1", "toolkit", "postgres", "password", 5432)
	assert.Nil(t, err)
	assert.Nil(t, pg.Ping())
	tnow := time.Now()

	ch, err := pg.KlinesNew("BYBIT", "i.BTCUSDT", tnow.Add(-24*time.Hour), tnow, 60)

	assert.Nil(t, err)
	assert.Condition(t, func() bool {
		return ch != nil || len(ch) > 0
	})
}
