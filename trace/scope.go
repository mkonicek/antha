package trace

import (
	"context"
	"sync"
)

type scopeKey int

const theScopeKey scopeKey = 0

func withScope(parent context.Context) context.Context {
	pscope, _ := parent.Value(theScopeKey).(*Scope)
	var s *Scope
	if pscope == nil {
		s = &Scope{}
	} else {
		s = pscope.MakeScope()
	}

	return context.WithValue(parent, theScopeKey, s)
}

func getScope(ctx context.Context) *Scope {
	s := ctx.Value(theScopeKey).(*Scope)
	if s == nil {
		panic("trace: scope not defined")
	}
	return s
}

// A Name is an identifier
type Name struct {
	scope *Scope
	idx   int
	desc  string
}

// Scope for a name
type Scope struct {
	lock   sync.Mutex
	parent *Scope
	pidx   int
	count  int
}

// MakeScope creates a child name scope
func (a *Scope) MakeScope() *Scope {
	a.lock.Lock()
	defer a.lock.Unlock()
	s := &Scope{parent: a, pidx: a.count}
	a.count++
	return s
}

// MakeName makes a name in this scope
func (a *Scope) MakeName(desc string) Name {
	a.lock.Lock()
	defer a.lock.Unlock()
	n := Name{scope: a, idx: a.count, desc: desc}
	a.count++
	return n
}

// MakeValue creates a value that has no relation to any existing values.
func MakeValue(ctx context.Context, desc string, v interface{}) Value {
	return &basicValue{
		name: getScope(ctx).MakeName(desc),
		v:    v,
	}
}

// MakeValueFrom creates a value that is a function of some existing values.
// This information is used to track value dependencies across instructions.
func MakeValueFrom(ctx context.Context, desc string, v interface{}, from ...Value) Value {
	return &fromValue{
		name: getScope(ctx).MakeName(desc),
		v:    v,
		from: from,
	}
}
