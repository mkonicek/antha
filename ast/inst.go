package ast

import (
	"context"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

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

// A Device is a scheduling plugin
type Device interface {
	CanCompile(Request) bool // Can this device compile this request

	// Compile produces a single-entry, single-exit DAG of instructions where
	// insts[0] is the entry point and insts[len(insts)-1] is the exit point
	Compile(ctx context.Context, cmds []Node) (insts []Inst, err error)
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

// An IncubateInst is a high-level command to incubate a component
type IncubateInst struct {
	// Time for which to incubate component
	Time wunit.Time
	// Temperature at which to incubate component
	Temp wunit.Temperature
	// Rate at which to shake incubator (force is device dependent)
	ShakeRate wunit.Rate
	// Radius at which ShakeRate is defined
	ShakeRadius wunit.Length

	// Time for which to pre-heat incubator
	PreTemp wunit.Temperature
	// Temperature at which to pre-heat incubator
	PreTime wunit.Time
	// Rate at which to pre-heat incubator
	PreShakeRate wunit.Rate
	// Radius at which PreShakeRate is defined
	PreShakeRadius wunit.Length
}

// A PromptInst is a high-level command to prompt a human
type PromptInst struct {
	Message string
}
