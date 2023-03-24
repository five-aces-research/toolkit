package lta

import "fmt"

// lta stands for Live TA

type Series interface {
	V(index int) float64 //Val(0) returns the most resent value, Val(1) the last etc.
	Data() []float64     //Returns all Values, used for init
	Updater
}

type Updater interface {
	ResolutionStartTime
	OnTick(NewTick bool) //Updates the latest Tick, When Update(true) adds a Tick
	SetLimit(i int)      //Sets the Limit that needs to be allocated for the indicator to work
	ExecuteLimit()       //Gets called once
	GetUpdateGroup() UpdateGroup
}

type URS[T any] struct {
	st, res int64
	name    string
	data    Dater[T]
	ug      UpdateGroup
	limit   int
	recent  float64
}

func (e *URS[T]) StartTime() int64 {
	return e.st
}

func (e *URS[T]) Resolution() int64 {
	return e.res
}

func (e *URS[T]) Data() []T {
	return e.data.Data()
}

func (e *URS[T]) V(i int) T {
	return e.data.V(i)
}

func (e *URS[T]) SetLimit(limit int) {
	if limit > e.limit {
		e.limit = limit
	}
}

func (e *URS[T]) ExecuteLimit() {
	fmt.Println(e.Name(), e.limit)
	e.data.SetLimit(e.limit)
}

func (e *URS[T]) GetUpdateGroup() UpdateGroup {
	return e.ug
}

func (e *URS[T]) Name() string {
	return e.name
}

func (e *URS[T]) SetName(name string) {
	return
}

type ResolutionStartTime interface {
	StartTime() int64
	Resolution() int64
	Name() string
	SetName(s string)
}

type UpdateGroup interface {
	Add(u Updater)
}

/*
In Kline
for v range tick:
	if v.Finished:
		SetOHCLV, AddNewCandle with close price
	else:
		SetOHCLV

	for_, vv := range updategroup:
			vv.OnTick(v.Finished)


*/
