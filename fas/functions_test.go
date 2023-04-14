package fas

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestRounding(t *testing.T) {
	fmt.Println(RoundValue(0.004248583, 0.001))

	for i := 0; i < 100; i++ {
		f := rand.Float64()
		s := RoundValue(f, 0.005)
		fmt.Println(s)
	}
}
