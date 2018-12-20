package target

import (
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/laboratory"
)

// A Device is a scheduling plugin
type Device interface {
	CanCompile(ast.Request) bool // Can this device compile this request
	MoveCost(from Device) int64  // A non-negative cost to move to this device

	// Compile produces a single-entry, single-exit DAG of instructions where
	// insts[0] is the entry point and insts[len(insts)-1] is the exit point
	Compile(labBuild *laboratory.LaboratoryBuilder, cmds []ast.Node) (insts []Inst, err error)
}
