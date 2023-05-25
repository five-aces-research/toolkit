package ta

//Some Functions that are used multiple Times in this package

type calc interface {
	int | int64 | float64 | float32 | int32
}

func change[T calc](old, new T) T {
	return new - old
}

// ShortestLenOfArray Return the lenght of the shortest array
func ShortestLenOfArray[T any](f ...[]T) (int, int) {
	var l int = len(f[0])
	var position int = 0
	for i := 1; i < len(f); i++ {
		if len(f[i]) < l {
			l = len(f[i])
			position = i
		}
	}
	return l, position
}

func highest[T calc](f ...T) (T, int) {
	var high T = f[0]
	var position int = 0
	for i := 1; i < len(f); i++ {
		if f[i] > high {
			high = f[i]
			position = i
		}
	}
	return high, position
}

func MinInt(f ...int) (val int, position int) {
	val = f[0]
	for i, v := range f {
		if v < val {
			position = i
			val = v
		}
	}
	return
}

func lowest[T calc](f ...T) (T, int) {
	var low T = f[0]
	var position int = 0
	for i := 1; i < len(f); i++ {
		if f[i] < low {
			low = f[i]
			position = i
		}
	}
	return low, position
}

func summe[T calc](f ...T) T {
	var avg T = 0
	for _, v := range f {
		avg += v
	}
	return avg
}

func average[T calc](f ...T) T {
	return summe[T](f...) / (T)(len(f))
}

func ArrayOperation[T calc](op func(T, T) T, f1 []T, f2 []T) []T {
	if len(f1) != len(f2) {
		return nil
	}
	fOut := make([]T, 0, len(f1))
	for i := 0; i < len(f1); i++ {
		fOut = append(fOut, op(f1[i], f2[i]))
	}
	return fOut
}

func Mul[T calc](v1, v2 T) T {
	return v1 * v2
}
