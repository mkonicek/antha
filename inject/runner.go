package inject

import (
	"fmt"

	"context"
)

// Basic signature of injectable functions
type RunFunc func(context.Context, Value) (Value, error)

// An injectable function
type Runner interface {
	Run(context.Context, Value) (Value, error) // Run the function and return results
}

// Untyped injectable function
type FuncRunner struct {
	RunFunc
}

func (a *FuncRunner) Run(ctx context.Context, value Value) (Value, error) {
	return a.RunFunc(ctx, value)
}

type TypedRunner interface {
	Runner
	Input() interface{}
	Output() interface{}
}

// Typed injectable function. Check if input parameter is assignable to In and
// output parameter is assignable to Out.
type CheckedRunner struct {
	RunFunc
	In  interface{}
	Out interface{}
}

func (a *CheckedRunner) Input() interface{} {
	return a.In
}

func (a *CheckedRunner) Output() interface{} {
	return a.Out
}

func (a *CheckedRunner) Run(ctx context.Context, value Value) (Value, error) {
	inT := a.In
	if err := AssignableTo(value, inT); err != nil {
		return nil, fmt.Errorf("input value not assignable to %T: %s", inT, err)
	}

	out, err := a.RunFunc(ctx, value)

	if err != nil {
		return out, err
	}

	outT := a.Out
	if err := AssignableTo(out, outT); err != nil {
		return nil, fmt.Errorf("output value not assignable to %T: %s", outT, err)
	}

	return out, nil
}
