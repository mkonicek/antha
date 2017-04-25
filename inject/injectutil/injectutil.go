// Package injectutil contains functionality built on top of the inject package
package injectutil

import (
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Multidimensial product space. Each map field is a finite dimension defined
// by its elements
type ProductSpace map[string][]interface{}

func (a ProductSpace) Elements() (r []inject.Value) {
	// Fix iteration order for dimensions
	var keys []string
	for k := range a {
		keys = append(keys, k)
	}

	if len(keys) == 0 {
		return
	}

	// Basic idea: successively append dimensions
	k := keys[0]
	for _, v := range a[k] {
		r = append(r, inject.Value{k: v})
	}

	for _, k := range keys[1:] {
		var rnew []inject.Value
		for _, v1 := range r {
			for _, v2 := range a[k] {
				// Cannot have an error here since ks are unique
				n, _ := v1.Concat(inject.Value{k: v2})
				rnew = append(rnew, n)
			}
		}
		r = rnew
	}

	return r
}

// Output of calling function on a product space.
type Product struct {
	Input, Output inject.Value
}

func CallCartesianProduct(ctx context.Context, runner inject.Runner, value inject.Value, space ProductSpace) ([]Product, error) {
	var ps []Product
	for _, elem := range space.Elements() {
		v := inject.MakeValue(elem)
		if in, err := value.Concat(v); err != nil {
			return nil, err
		} else if out, err := runner.Run(ctx, in); err != nil {
			return nil, err
		} else {
			ps = append(ps, Product{Input: in, Output: out})
		}
	}
	return ps, nil
}
