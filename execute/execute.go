// Package execute connects Antha elements to the trace execution
// infrastructure.
package execute

import (
	"context"
	"fmt"

	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/codegen"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/sampletracker"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/workflow"
)

// Result of executing a workflow.
type Result struct {
	Workflow *workflow.Workflow
	Input    []ast.Node
	Insts    []ast.Inst
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
func Run(parent context.Context, opt Opt) (res *Result, err error) {
	ctx := sampletracker.NewContext(target.WithTarget(withID(parent, opt.ID), opt.Target))

	w, err := workflow.New(workflow.Opt{FromDesc: opt.Workflow})
	if err != nil {
		return nil, err
	}

	if _, err := setParams(ctx, w, opt.Params, opt.TransitionalReadLocalFiles); err != nil {
		return nil, err
	}

	ctxTr, tr := WithTrace(ctx)
	defer func() {
		if res := recover(); res == nil {
			return
		} else if uErr, ok := res.(UserError); ok {
			// Errorf internally calls panic, which is *not* the Go
			// way. But until we fix that, to avoid full stack traces
			// popping out here, we catch this case, and we deliberately
			// do not attach a stack trace to it.
			err = uErr
		} else {
			err = fmt.Errorf("%s\n%s", res, inject.ElementStackTrace())
		}
	}()
	if err := w.Run(ctxTr); err != nil {
		return nil, err
	}

	t, err := target.GetTarget(ctx)
	if err != nil {
		return nil, err
	}

	nodes, err := getMaker(ctx).MakeNodes(tr.Instructions())
	if err != nil {
		return nil, err
	}

	instrs, err := codegen.Compile(ctx, t, nodes)
	if err != nil {
		return nil, err
	}

	return &Result{
		Workflow: w,
		Input:    nodes,
		Insts:    instrs,
	}, nil
}
