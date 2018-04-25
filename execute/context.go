package execute

import (
	"context"
)

type idContextKey int

const theIDContextKey idContextKey = 0

type withExecute struct {
	ID    string
	Maker *maker
}

type elementNameKey int

const theElementNameKey elementNameKey = 0

type withElementName struct {
	Name string
}

func getMaker(ctx context.Context) *maker {
	return ctx.Value(theIDContextKey).(*withExecute).Maker
}

func getID(ctx context.Context) string {
	v, ok := ctx.Value(theIDContextKey).(*withExecute)
	if !ok {
		return ""
	}
	return v.ID
}

func withID(parent context.Context, id string) context.Context {
	return context.WithValue(parent, theIDContextKey, &withExecute{
		ID:    id,
		Maker: newMaker(),
	})
}

// WithElementName returns a new context that stores the current element name
func WithElementName(parent context.Context, name string) context.Context {
	return context.WithValue(parent, theElementNameKey, &withElementName{
		Name: name,
	})
}

func getElementName(ctx context.Context) string {
	v, ok := ctx.Value(theElementNameKey).(*withElementName)
	if !ok {
		return ""
	}
	return v.Name
}
