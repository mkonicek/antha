package cache

import (
	"context"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
)

const (
	theCtxKey ctxKey = "inventoryCache"
)

type ctxKey string

type Cache interface {
	inventory.Inventory
	ReturnObject(ctx context.Context, obj interface{}) error
	IsFromCache(ctx context.Context, obj interface{}) bool
}

func fromContext(ctx context.Context) Cache {
	return ctx.Value(theCtxKey).(Cache)
}

//NewContext returns a context with the new cache
func NewContext(ctx context.Context, pc Cache) context.Context {
	return context.WithValue(ctx, theCtxKey, pc)
}

//GetCache returns the plate cache from the context
func GetCache(ctx context.Context) Cache {
	return fromContext(ctx)
}

// NewComponent returns a new component of the given type
func NewComponent(ctx context.Context, typ string) (*wtype.Liquid, error) {
	return fromContext(ctx).NewComponent(ctx, typ)
}

// NewPlate returns a new plate of the given type
func NewPlate(ctx context.Context, typ string) (*wtype.Plate, error) {
	return fromContext(ctx).NewPlate(ctx, typ)
}

// NewTipwaste returns a new tipwaste of the given type
func NewTipwaste(ctx context.Context, typ string) (*wtype.LHTipwaste, error) {
	return fromContext(ctx).NewTipwaste(ctx, typ)
}

// NewTipbox returns a new tipbox of the given type
func NewTipbox(ctx context.Context, typ string) (*wtype.LHTipbox, error) {
	return fromContext(ctx).NewTipbox(ctx, typ)
}

// ReturnObject return an object to the cache to be cleaned
func ReturnObject(ctx context.Context, obj interface{}) error {
	return fromContext(ctx).ReturnObject(ctx, obj)
}

// IsFromCache returns true if the object is from the cache
func IsFromCache(ctx context.Context, obj interface{}) bool {
	return fromContext(ctx).IsFromCache(ctx, obj)
}
