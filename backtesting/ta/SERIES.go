package ta

// A Series  is an Interface that Every Indicator needs to implement to Communicate its Value to other Indicators. OHCLV Data is also an Indicator
type Series interface {
	Data() []float64
	V(index int) float64 // return the last minus index value. So V(0) returns the latest Candle

	ResolutionStartTime
}

// A Condition is an Interface that Every []bool Series needs to implement. Such as And, Greater, OR, IFF
type Condition interface {
	Data() []bool
	V(index int) bool // return the last minus index value. So V(0) returns the latest Candle

	ResolutionStartTime
}

// ErrorResolutionStartime is needed to sync the Indicators in a fast way
type ResolutionStartTime interface {
	StartTime() int64    //Starttime in Unix Seconds
	Resolution() int64   // Resolution in Seconds
	Name() string        // Name is needed for identification
	SetName(name string) // SetName also used for identification
}

// ERS implements the ErrorResolutionStartTime and can be implemented in a Series and Condition
type ERS[T any] struct {
	st   int64
	res  int64
	data []T
	name string
}

func (e *ERS[T]) StartTime() int64 {
	return e.st
}

func (e *ERS[T]) Resolution() int64 {
	return e.res
}

func (e *ERS[T]) Data() []T {
	return e.data
}

func (e *ERS[T]) V(index int) T {
	return e.data[len(e.data)-1-index]
}
func (e *ERS[T]) Name() string {
	return e.name
}

func (e *ERS[T]) SetName(name string) {
	e.name = name
}
