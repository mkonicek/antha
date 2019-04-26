package instructions

import (
	"github.com/antha-lang/antha/laboratory/effects/id"
)

// An Inst is a instruction. These are "low-level" - i.e. the results
// of calls to Device.Compile(), *not* instructions issued by
// intrinsics.
type Inst interface {
	Id() string
	// Idempotent - will not change an id once set.
	SetId(*id.IDGenerator)
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

type DependsMixin struct {
	Depends []Inst
}

func (a *DependsMixin) DependsOn() []Inst {
	return a.Depends
}

func (a *DependsMixin) SetDependsOn(x ...Inst) {
	a.Depends = x
}

func (a *DependsMixin) AppendDependsOn(x ...Inst) {
	a.Depends = append(a.Depends, x...)
}

type IdMixin struct {
	id string
}

func (a IdMixin) Id() string {
	return a.id
}

func (a *IdMixin) SetId(idGen *id.IDGenerator) {
	if a.id == "" {
		a.id = idGen.NextID()
	}
}
