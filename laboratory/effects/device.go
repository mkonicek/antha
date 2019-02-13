package effects

import (
	"github.com/antha-lang/antha/workflow"
)

// A Device is a scheduling plugin
type Device interface {
	CanCompile(Request) bool // Can this device compile this request

	// Compile produces a single-entry, single-exit DAG of instructions where
	// insts[0] is the entry point and insts[len(insts)-1] is the exit point
	Compile(labEffects *LaboratoryEffects, cmds []Node) (insts []Inst, err error)

	// Must be idempotent and thread safe
	Connect(*workflow.Workflow) error
	// Must be idempotent and thread safe
	Close()
}

// An Inst is a instruction
type Inst interface {
	// Device that this instruction was generated for
	Device() Device
	// DependsOn returns instructions that this instruction depends on
	DependsOn() []Inst
	// SetDependsOn sets to the list of dependencies to only the args
	SetDependsOn(...Inst)
	// AppendDependsOn adds to the args to the existing list of dependencies
	AppendDependsOn(...Inst)
}

type Insts []Inst

// SequentialOrder takes a slice of instructions and modifies them
// in-place, resetting to sequential dependencies.
func (insts Insts) SequentialOrder() {
	if len(insts) > 1 {
		prev := insts[0]
		for _, cur := range insts[1:] {
			cur.SetDependsOn(prev)
			prev = cur
		}
	}
}
