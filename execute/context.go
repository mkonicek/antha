package execute

import (
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
)

type contextKey int

const theContextKey contextKey = 0

type withExecute struct {
	Id    string
	Maker *maker
}

func getMaker(ctx context.Context) *maker {
	return ctx.Value(theContextKey).(*withExecute).Maker
}

func getId(ctx context.Context) string {
	v, ok := ctx.Value(theContextKey).(*withExecute)
	if !ok {
		return ""
	}
	return v.Id
}

func WithId(parent context.Context, id string) context.Context {
	return context.WithValue(parent, theContextKey, &withExecute{
		Id:    id,
		Maker: newMaker(),
	})
}
