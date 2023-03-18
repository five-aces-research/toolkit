package lta

// lta stands for Live TA

type Series interface {
	Val(index int) float64 //Val(0) returns the most resent value, Val(1) the last etc.
	Data() []float64       //Returns all Values, used for init
	Updater
}

type Updater interface {
	ResolutionStartTime
	OnTick(NewTick bool) //Updates the latest Tick, When Update(true) adds a Tick
	SetLimit(i int)      //Sets the Limit that needs to be allocated for the indicator to work
	ExecuteLimit()       //Gets called once
	GetUpdateGroup() *UpdateGroup
}

type ResolutionStartTime interface {
	StartTime() int64
	Resolution() int64
	Name() string
	SetName(s string)
}

type UpdateGroup struct {
	Name string
	ug   []Updater
	re   int64
	exit bool
}
