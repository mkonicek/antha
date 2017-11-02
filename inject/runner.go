package inject

import (
	"fmt"

	"context"
)

// RunFunc is the signature of injectable functions
type RunFunc func(context.Context, Value) (Value, error)

// A Runner is an injectable function
type Runner interface {
	// Run the function and return results
	Run(context.Context, Value) (Value, error)
}

// A FuncRunner is an untyped injectable function
type FuncRunner struct {
	RunFunc
}

// Run implements a Runner
func (a *FuncRunner) Run(ctx context.Context, value Value) (Value, error) {
	return a.RunFunc(ctx, value)
}

// A TypedRunner is a typed injectable function
type TypedRunner interface {
	Runner
	Input() interface{}
	Output() interface{}
}

// A CheckedRunner is a typed injectable function. It checks if input parameter
// is assignable to In and output parameter is assignable to Out.
type CheckedRunner struct {
	RunFunc
	In  interface{}
	Out interface{}
}

// Input returns an example of an input to this Runner
func (a *CheckedRunner) Input() interface{} {
	return a.In
}

// Output returns an example of an output of this Runner
func (a *CheckedRunner) Output() interface{} {
	return a.Out
}

// Run implements a Runner
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
