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
	lock         sync.Mutex
	instructions []*commandInst
}

// Issue an instruction - this records the instruction into the trace.
func (tr *Trace) Issue(instruction *commandInst) {
	tr.lock.Lock()
	defer tr.lock.Unlock()

	tr.instructions = append(tr.instructions, instruction)
}

// Returns a (shallow) copy (to avoid data races) of the issued
// instructions.
func (tr *Trace) Instructions() []*commandInst {
	tr.lock.Lock()
	defer tr.lock.Unlock()

	clone := make([]*commandInst, len(tr.instructions))
	copy(clone, tr.instructions)
	return clone
}

// Issue an instruction - this records the instruction into the trace.
func Issue(ctx context.Context, instruction *commandInst) {
	getTrace(ctx).Issue(instruction)
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
