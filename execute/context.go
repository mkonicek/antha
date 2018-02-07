package execute

import (
	"context"
)

type contextKey int

const theContextKey contextKey = 0

type withExecute struct {
	ID    string
	Maker *maker
}

func getMaker(ctx context.Context) *maker {
	return ctx.Value(theContextKey).(*withExecute).Maker
}

func getID(ctx context.Context) string {
	v, ok := ctx.Value(theContextKey).(*withExecute)
	if !ok {
		return ""
	}
	return v.ID
}

func withID(parent context.Context, id string) context.Context {
	return context.WithValue(parent, theContextKey, &withExecute{
		ID:    id,
		Maker: newMaker(),
	})
}
