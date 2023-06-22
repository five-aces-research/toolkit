package te

func Limit(offset float64) *ScaledLimit {
	return NewScaledLimit(offset, offset, 1).Distribution(func(price, min, max float64, orderCount int) [][2]float64 {
		return [][2]float64{{1, price - price*(min/100)}}
	})
}
