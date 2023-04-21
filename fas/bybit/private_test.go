package bybit

import (
	"fmt"
	"testing"
	"time"
)

func TestHistoricalOrders(t *testing.T) {
	pr := NewPrivate("RoRo", "k1ofNVyJkAnoXtwj1l", "", false)

	res, err := pr.GetOrderHistory(nil, time.Unix(0, 0), time.Now())
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	for _, v := range res {
		fmt.Printf("%+v\n", v)
	}

}
