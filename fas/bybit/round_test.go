package bybit

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestRound(t *testing.T) {
	fmt.Println(roundValue(0.004248583, 0.001))

	for i := 0; i < 100; i++ {
		f := rand.Float64()
		s := roundValue(f, 0.005)
		fmt.Println(s)
	}
}
