package trace

import (
	"context"
)

type resolverKey int

const theResolverKey resolverKey = 0

// A Resolver returns values for some instructions
type Resolver func(ctx context.Context, insts []interface{}) (values map[int]interface{}, err error)

type resolverCtx struct {
	Parent   context.Context
	Resolver Resolver
}

// WithResolver creates a new resolver context. When all pool contexts are
// blocked on Reads, the Tracer calls Resolvers to compute the promised values.
// If there are no Resolvers, the promised values are nil.
func WithResolver(parent context.Context, resolver Resolver) context.Context {
	return context.WithValue(parent, theResolverKey, &resolverCtx{
		Parent:   parent,
		Resolver: resolver,
	})
}

func nilResolver(ctx context.Context, xs []interface{}) (map[int]interface{}, error) {
	m := make(map[int]interface{})
	for idx := range xs {
		m[idx] = nil
	}
	return m, nil
}

func resolve(ctx context.Context, instps []instp) error {
	// TODO: Deterministic sort of instructions
	origCtx := ctx
	for len(instps) != 0 {
		var insts []interface{}
		for _, v := range instps {
			insts = append(insts, v.inst)
		}

		resolver := nilResolver
		if v := ctx.Value(theResolverKey); v == nil {
		} else if rctx, ok := v.(*resolverCtx); ok {
			resolver = rctx.Resolver
			ctx = rctx.Parent
		}

		values, err := resolver(origCtx, insts)
		if err != nil {
			return err
		}

		for idx, v := range values {
			p := instps[idx].promise
			p.set(v)
			close(p.out)
		}

		var next []instp
		for idx, v := range instps {
			if _, seen := values[idx]; !seen {
				next = append(next, v)
			}
		}
		instps = next
	}
	return nil
}
