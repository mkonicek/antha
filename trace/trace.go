// Package trace provides deferred execution of instruction streams generated
// at runtime. Issue adds an instruction to be executed while the actual
// execution is deferred until all goroutines in a trace context are blocked on
// Reads of instruction results (i.e., Promises).
package trace

import (
	"context"
	"runtime/debug"
	"sync"
)

type instrKey int

const theInstrKey instrKey = 0

func withTrace(parent context.Context) context.Context {
	return context.WithValue(parent, theInstrKey, &trace{})
}

// NewContext creates a new trace and pool context.
//
// The general template for clients is:
//
//   ctx, cancel, allDone := NewContext(parent, ...)
//   defer cancel()
//   ...
//   Go(ctx, ...)
//   ...
//   select {
//     case <-allDone():
//       return ctx.Err()
//     ...
//   }
func NewContext(parent context.Context) (context.Context, context.CancelFunc, DoneFunc) {
	c, cfn, dfn := WithPool(withTrace(withScope(parent)))
	return c, cfn, dfn
}

func getTrace(ctx context.Context) *trace {
	t := ctx.Value(theInstrKey).(*trace)
	if t == nil {
		panic("trace: trace not found")
	}
	return t
}

// Read a promise value
func Read(ctx context.Context, p *Promise) (Value, error) {
	if v := p.value(); v != nil {
		return v, nil
	}

	pctx := getPool(ctx)
	tr := getTrace(ctx)

	pctx.lock.Lock()
	pctx.blocked[p] = true
	pctx.tryUnblockWithLock(tr)
	pctx.lock.Unlock()

	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case v, ok := <-p.out:
		if ok {
			p.set(v)
		}
	}

	if err == nil {
		err = p.err()
	}

	pctx.remove(p)

	return p.value(), err
}

// ReadAll reads all promises concurrently
func ReadAll(parent context.Context, ps ...*Promise) ([]Value, error) {
	var mapLock sync.Mutex
	vs := make(map[int]Value)

	pctx, cancel, allDone := WithPool(parent)
	defer cancel()

	for idx, promise := range ps {
		i := idx
		p := promise
		Go(pctx, func(ctx context.Context) error {
			v, err := Read(ctx, p)
			if err != nil {
				return err
			}

			mapLock.Lock()
			defer mapLock.Unlock()
			vs[i] = v

			return nil
		})
	}

	<-allDone()
	if err := pctx.Err(); err != nil {
		return nil, err
	}

	var ret []Value
	for i := 0; i < len(ps); i++ {
		ret = append(ret, vs[i])
	}
	return ret, nil
}

// Pair to track instructions with their promised return values
type instp struct {
	inst    interface{}
	promise *Promise
}

// Issue an instruction to execute in the trace context; return promise for
// return value.
func Issue(ctx context.Context, inst interface{}) *Promise {
	result := getScope(ctx).MakeName("")
	out := make(chan interface{})

	p := &Promise{
		construct: func(v interface{}) Value {
			return &promiseValue{
				name: result,
				inst: inst,
				v:    v,
			}
		},
		out: out,
	}

	t := getTrace(ctx)
	t.lock.Lock()
	defer t.lock.Unlock()

	t.issues = append(t.issues, instp{inst: inst, promise: p})

	return p
}

// Flush forces code generation and execution for current set of issued commands
func Flush(ctx context.Context) error {
	p := Issue(ctx, nil)
	_, err := Read(ctx, p)
	return err
}

type trace struct {
	// Given a list of pending instructions, return their values
	lock   sync.Mutex
	issues []instp
}

func (a *trace) signal(lockedPool *poolCtx) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.signalWithLock(lockedPool)
}

func (a *trace) signalWithLock(lockedPool *poolCtx) (err error) {
	defer func() {
		if res := recover(); res != nil {
			err = &Error{BaseError: res, Stack: debug.Stack()}
		}
	}()

	// Don't execute until all active all blocked
	if lockedPool.alive == len(lockedPool.blocked) {
		err = a.execute(lockedPool.Context)
	}
	return
}

func (a *trace) execute(ctx context.Context) error {
	if err := resolve(ctx, a.issues); err != nil {
		return err
	}

	a.issues = nil

	return nil
}
