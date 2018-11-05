package target

import (
	"context"

	"github.com/antha-lang/antha/ast"
)

// A Device is a scheduling plugin
type Device interface {
	CanCompile(ast.Request) bool // Can this device compile this request
	MoveCost(from Device) int64  // A non-negative cost to move to this device

	// Compile produces a single-entry, single-exit DAG of instructions where
	// insts[0] is the entry point and insts[len(insts)-1] is the exit point
	Compile(ctx context.Context, cmds []ast.Node) (insts []Inst, err error)
}
