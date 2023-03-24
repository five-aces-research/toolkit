package lta

import (
	"sync"
)

type Dater[T any] interface {
	V(index int) T
	SetLimit(limit int)
	Append(v T)
	SetValue(index int, val T)
	Data() []T
}

// array is a generic container for an array of type T with custom methods.
type array[T any] struct {
	arr           []T
	limit         int
	allocateLimit int
	lock          sync.RWMutex
}

// Array is a constructor that creates a new array of type T and initializes it with the given slice.
func Array[T any](f []T) Dater[T] {
	d := new(array[T])
	d.arr = f
	d.limit = 1
	return d
}

// Append adds an element of type T to the end of the array.
func (d *array[T]) Append(v T) {
	d.arr = append(d.arr, v)
	if len(d.arr) == d.allocateLimit {
		newArr := make([]T, d.limit, d.allocateLimit)
		copy(newArr, d.arr[len(d.arr)-d.limit:])
		d.arr = newArr
	}
}

// V returns the value at the given index of the array.
func (d *array[T]) V(index int) T {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.arr[len(d.arr)-1-index]
}

// SetLimit sets the limit of the array and resizes the underlying array accordingly.
func (d *array[T]) SetLimit(limit int) {
	limit++
	d.limit = limit
	if limit < 10 {
		d.allocateLimit = limit + 10
	} else {
		d.allocateLimit = limit * 2
	}
	if len(d.arr) < limit {
		limit = len(d.arr)
	}

	newArr := make([]T, limit, d.allocateLimit)
	copy(newArr, d.arr[len(d.arr)-limit:])
	d.arr = newArr
}

// SetValue sets the value at the given index of the array to the specified value.
func (d *array[T]) SetValue(index int, val T) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.arr[len(d.arr)-1-index] = val
}

// Data returns the underlying array of type T.
func (d *array[T]) Data() []T {
	return d.arr
}
