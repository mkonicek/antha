// Package execute connects Antha elements to the trace execution
// infrastructure.
package execute

import (
	"context"
	"errors"

	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/trace"
	"github.com/antha-lang/antha/workflow"
)

var (
	cannotConfigure = errors.New("cannot configure liquid handler")
)

// TODO(ddn): extend result when protocols can block

// Result of executing a workflow.
type Result struct {
	Workflow *workflow.Workflow
	Input    []ast.Node
	Insts    []target.Inst
}

type Opt struct {
	Workflow *workflow.Desc // Or workflow directly
	Params   *RawParams     // Or parameters directly
	Target   *target.Target // Target machine configuration
	Id       string         // Job Id
}

// Simple entrypoint for one-shot execution of workflows.
func Run(parent context.Context, opt Opt) (*Result, error) {
	w, err := workflow.New(workflow.Opt{FromDesc: opt.Workflow})
	if err != nil {
		return nil, err
	}

	if _, err := setParams(parent, opt.Params, w); err != nil {
		return nil, err
	}

	ctx := target.WithTarget(WithId(parent, opt.Id), opt.Target)

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
		if eerr, ok := terr.BaseError.(*Error); ok {
			err = eerr
		}
	}

	return nil, err
}
