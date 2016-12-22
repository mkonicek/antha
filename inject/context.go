// Package inject implements delayed binding of function calls to runtime.
package inject

import (
	"context"
	"errors"
)

type injectKey int

const theInjectKey injectKey = 0

var noRegistry = errors.New("no registry found")
var funcNotFound = errors.New("not found")

// Create a new inject context
func NewContext(parent context.Context) context.Context {
	return context.WithValue(parent, theInjectKey, &registry{parent: parent})
}

func getRegistry(parent context.Context) *registry {
	r, ok := parent.Value(theInjectKey).(*registry)
	if !ok {
		return nil
	}
	return r
}

// Add a function to the inject context
func Add(parent context.Context, name Name, runner Runner) error {
	reg := getRegistry(parent)
	if reg == nil {
		return noRegistry
	}
	return reg.Add(name, runner)
}

func Find(parent context.Context, query NameQuery) (Runner, error) {
	type result struct {
		runner Runner
		level  int
	}

	ctx := parent
	level := 0
	reg := getRegistry(ctx)
	var results []result
	for reg != nil {
		runners, err := reg.Find(query)
		if err != nil {
			return nil, err
		}
		for _, runner := range runners {
			results = append(results, result{level: level, runner: runner})
		}
		level += 1
		ctx = reg.parent
		reg = getRegistry(ctx)
	}

	// TODO: better matching heuristics?
	for _, r := range results {
		return r.runner, nil
	}
	return nil, funcNotFound
}

// Call a function that satisfies the query
func Call(parent context.Context, query NameQuery, value Value) (Value, error) {
	r, err := Find(parent, query)
	if err != nil {
		return nil, err
	}
	return r.Run(parent, value)
}
