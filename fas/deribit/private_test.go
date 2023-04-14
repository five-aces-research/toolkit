package deribit

import (
	"fmt"
	"testing"
)

func TestPrivate_SetOrder(t *testing.T) {
	d := NewPrivate("Bubu", "", "", false)
	fmt.Println(d.Position("BTC"))
}
