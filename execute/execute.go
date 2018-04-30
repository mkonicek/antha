// Package execute connects Antha elements to the trace execution
// infrastructure.
package execute

import (
	"context"

	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/trace"
	"github.com/antha-lang/antha/workflow"
)

// Result of executing a workflow.
type Result struct {
	Workflow *workflow.Workflow
	Input    []ast.Node
	Insts    []target.Inst
}

// An Opt are options for Run.
type Opt struct {
	// Target machine configuration
	Target *target.Target
	// Deprecated for separate assignment of values to workflow. Raw workflow.
	Workflow *workflow.Desc
	// Deprecated for separate assignment of values to workflow. Raw parameters.
	Params *RawParams
	// Job ID.
	ID string
	// Deprecated for separate assignment of values to workflow. If true, read
	// content for each wtype.File from file of the same name in the current
	// directory.
	TransitionalReadLocalFiles bool
}

// Run is a simple entrypoint for one-shot execution of workflows.
func Run(parent context.Context, opt Opt) (*Result, error) {
	ctx := target.WithTarget(withID(parent, opt.ID), opt.Target)

	w, err := workflow.New(workflow.Opt{FromDesc: opt.Workflow})
	if err != nil {
		return nil, err
	}

	if _, err := setParams(ctx, w, opt.Params, opt.TransitionalReadLocalFiles); err != nil {
		return nil, err
	}

	r := &resolver{}

	err = w.Run(trace.WithResolver(ctx, func(ctx context.Context, insts []interface{}) (map[int]interface{}, error) {
		return r.resolve(ctx, insts)
	}))

	if err == nil {
		return &Result{
			Workflow: w,
			Input:    r.nodes,
			Insts:    r.insts,
		}, nil
	}

	// Unwrap execute.Error
	if terr, ok := err.(*trace.Error); ok {
		if unwrapped, ok := unwrapError(terr.BaseError); ok {
			err = unwrapped
		}
	}

	return nil, err
}
