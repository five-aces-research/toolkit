package bybit

import (
	"fmt"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	bb := NewPublic()
	l, err := bb.GetOpenInterest("l.btcusdt", 360,
		time.Unix(1642326016, 0), time.Now())
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	for _, v := range l {
		fmt.Println(v)
	}

}
