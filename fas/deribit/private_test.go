package deribit

import (
	"fmt"
	"testing"
)

func TestPrivate_SetOrder(t *testing.T) {
	d := NewPrivate("Main", "dJ_gOAyI", "", false)
	res, err := d.AccountInformation()
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	fmt.Println(res.TotalEquity, res.FreeEquity)

	for _, v := range res.Coins {
		fmt.Println(v)
	}

}
