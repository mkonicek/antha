package execute

import (
	"context"
	"sync"
)

// This is pretty gross. What we do here is add into a context a
// container into which we later inject instructions.  Eventually,
// when the stack unwinds, we grab all these instructions. This is an
// abuse of the intent of context, but to fix this will require
// changes to the code generation of elements.

type Trace struct {
	lock   sync.Mutex
	instrs []*commandInst
}

func (tr *Trace) Issue(instr *commandInst) {
	tr.lock.Lock()
	defer tr.lock.Unlock()

	tr.instrs = append(tr.instrs, instr)
}

func (tr *Trace) Instrs() []*commandInst {
	tr.lock.Lock()
	defer tr.lock.Unlock()

	clone := make([]*commandInst, len(tr.instrs))
	copy(clone, tr.instrs)
	return clone
}

func Issue(ctx context.Context, instr *commandInst) {
	getTrace(ctx).Issue(instr)
}

type traceKey int

const theTraceKey traceKey = 0

func getTrace(ctx context.Context) *Trace {
	tr, ok := ctx.Value(theTraceKey).(*Trace)
	if !ok || tr == nil {
		panic("trace: trace not found")
	}
	return tr
}

func WithTrace(parent context.Context) (context.Context, *Trace) {
	tr := &Trace{}
	return context.WithValue(parent, theTraceKey, tr), tr
}
