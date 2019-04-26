package effects

import (
	"github.com/antha-lang/antha/instructions"
	"github.com/antha-lang/antha/workflow"
)

type Device interface {
	CanCompile(instructions.Request) bool // Can this device compile this request

	// Compile produces a single-entry, single-exit DAG of instructions where
	// insts[0] is the entry point and insts[len(insts)-1] is the exit point
	Compile(labEffects *LaboratoryEffects, dir string, cmds []instructions.Node) (instructions.Insts, error)

	// Must be idempotent and thread safe
	Connect(*workflow.Workflow) error
	// Must be idempotent and thread safe
	Close()

	Id() workflow.DeviceInstanceID
}
