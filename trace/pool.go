package trace

import (
	"context"
	"errors"
	"runtime/debug"
	"sync"
)

type poolKey int

const thePoolKey poolKey = 0

var poolDoneError error = errors.New("trace: done error") // Dummy error to signal sucessful execution

type poolCtx struct {
	context.Context
	lock    sync.Mutex
	alive   int
	blocked map[*Promise]bool
	err     error
	done    chan struct{}
}

func (a *poolCtx) Value(key interface{}) interface{} {
	if key == thePoolKey {
		return a
	}
	return a.Context.Value(key)
}

func (a *poolCtx) Err() error {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.err == nil {
		return a.Context.Err()
	}
	if a.err == poolDoneError {
		return nil
	}
	return a.err
}

func (a *poolCtx) Done() <-chan struct{} {
	return a.done
}

func (a *poolCtx) cancelWithLock(err error) {
	if a.err == nil {
		a.err = err
		close(a.done)
	}
}

func (a *poolCtx) cancel(err error) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.cancelWithLock(err)
}

func (a *poolCtx) remove(p *Promise) {
	a.lock.Lock()
	defer a.lock.Unlock()
	delete(a.blocked, p)
}

type DoneFunc func() <-chan struct{}

// Create a pool context. It is done when all Go()-created routines return
// normally or when any return an error.
func WithPool(parent context.Context) (context.Context, context.CancelFunc, DoneFunc) {
	tr := getTrace(parent)
	done := make(chan struct{})
	pctx := &poolCtx{
		Context: parent,
		done:    done,
		alive:   1,
		blocked: make(map[*Promise]bool),
	}
	rootPromise := &Promise{}
	pctx.blocked[rootPromise] = true

	if parent.Done() != nil {
		go func() {
			select {
			case <-parent.Done():
				pctx.cancel(parent.Err())
			case <-done:
			}
		}()
	}

	return pctx,
		func() { pctx.cancel(context.Canceled) },
		func() <-chan struct{} {
			pctx.remove(rootPromise)
			decrement(pctx, tr, 1, nil)
			return done
		}
}

func getPool(ctx context.Context) *poolCtx {
	c := ctx.Value(thePoolKey).(*poolCtx)
	if c == nil {
		panic("trace: pool not found")
	}
	return c
}

func tryUnblock(tr *trace, pctx *poolCtx) {
	if err := tr.signal(pctx); err != nil {
		pctx.cancelWithLock(err)
	}
}

func decrement(pctx *poolCtx, tr *trace, delta int, err error) {
	pctx.lock.Lock()
	defer pctx.lock.Unlock()
	pctx.alive -= delta
	tryUnblock(tr, pctx)

	var cancel error
	if err != nil {
		cancel = err
	}
	if pctx.alive == 0 && cancel == nil {
		cancel = poolDoneError
	}
	if cancel != nil {
		pctx.cancelWithLock(cancel)
	}
}

// Create a new goroutine in the pool context
func Go(parent context.Context, fn func(ctx context.Context) error) {
	pctx := getPool(parent)
	ctx := withScope(parent)
	tr := getTrace(parent)

	pctx.lock.Lock()
	defer pctx.lock.Unlock()

	pctx.alive += 1

	go func() {
		var err error
		defer func() {
			decrement(pctx, tr, 1, err)
		}()
		defer func() {
			if res := recover(); res != nil {
				err = &Error{BaseError: res, Stack: debug.Stack()}
			}
		}()
		err = fn(ctx)
	}()
}
