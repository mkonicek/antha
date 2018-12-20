package effects

import (
	"sync"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/ast"
)

// a CommandInst is a generic intrinsic instruction
type CommandInst struct {
	// Arguments to this command. Used to determine command dependencies.
	Args []*wtype.Liquid
	// Components created by this command. Returned back to user code
	Result  []*wtype.Liquid
	Command *ast.Command
}

type Trace struct {
	lock         sync.Mutex
	instructions []*CommandInst
}

func NewTrace() *Trace {
	return &Trace{}
}

// Issue an instruction - this records the instruction into the trace.
func (tr *Trace) Issue(instruction *CommandInst) {
	tr.lock.Lock()
	defer tr.lock.Unlock()

	tr.instructions = append(tr.instructions, instruction)
}

// Returns a (shallow) copy (to avoid data races) of the issued
// instructions.
func (tr *Trace) Instructions() []*CommandInst {
	tr.lock.Lock()
	defer tr.lock.Unlock()

	clone := make([]*CommandInst, len(tr.instructions))
	copy(clone, tr.instructions)
	return clone
}
