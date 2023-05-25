package pg_wrapper

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQQ(t *testing.T) {
	pg, err := Connect("127.0.0.1", "toolkit", "postgres", "password", 5432)
	assert.Nil(t, err)
	res, err := pg.GetTickers()
	assert.Nil(t, err)

	for _, v := range res {
		fmt.Println(v)
	}

	min, max, err := pg.GetMinMax("bybit", "L.BTCUSDT", 360)
	assert.Nil(t, err)
	fmt.Println(min, max)

}
