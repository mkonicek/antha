package trace

import (
	"sync"
)

// A Promise is placeholder for value to be returned by an instruction
type Promise struct {
	lock      sync.Mutex
	construct func(interface{}) Value
	v         Value
	e         error
	out       chan interface{}
}

func (a *Promise) err() error {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.e
}

func (a *Promise) value() Value {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.v
}

func (a *Promise) set(v interface{}) {
	if a.construct == nil {
		return
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	if a.v != nil {
		a.v = a.construct(v)
	}
}
