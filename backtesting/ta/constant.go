package ta

type emptySeries[v any] struct {
	ERS[v]
}

func Constant[T any](b []T, st, res int64, name string) *emptySeries[T] {
	s := new(emptySeries[T])
	s.st = st
	s.res = res
	s.data = b
	s.name = name
	return s
}

type constant[T any] struct {
	ERS[T]
}

func constantS(src Condition, a float64) Series {
	s := new(constant[float64])
	s.st = src.StartTime()
	s.res = src.Resolution()
	s.data = make([]float64, len(src.Data()), len(src.Data()))
	for i := range s.data {
		s.data[i] = a
	}
	return s
}

func constantB(src Condition, a bool) Condition {
	s := new(constant[bool])
	s.st = src.StartTime()
	s.res = src.Resolution()
	s.data = make([]bool, len(src.Data()), len(src.Data()))
	if a {
		for i := range s.data {
			s.data[i] = a
		}
	}
	return s
}

type MaFunc func(s Series, i int) Series

func (m MaFunc) Name() string {
	return m(Constant([]float64{1, 2, 3, 4}, 0, 3600, ""), 2).Name()
}
