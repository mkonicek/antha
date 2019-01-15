package inventory

import (
	"context"
	"errors"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

var (
	// ErrUnknownType is returned if type is not in inventory
	ErrUnknownType = errors.New("unknown type requested from inventory")

	errCannotListPlates = errors.New("cannot list plates")
)

const (
	// WaterType is the component type of water
	WaterType = "water"
)

const (
	theCtxKey ctxKey = "inventory"
)

type ctxKey string

func fromContext(ctx context.Context) Inventory {
	return ctx.Value(theCtxKey).(Inventory)
}

// An Inventory returns items by name
type Inventory interface {
	NewComponent(typ string) (*wtype.Liquid, error)
	NewPlate(typ string) (*wtype.Plate, error)
}

// NewContext returns a context with the given inventory
func NewContext(ctx context.Context, inv Inventory) context.Context {
	return context.WithValue(ctx, theCtxKey, inv)
}

// GetInventory returns an Inventory instance from Context
func GetInventory(ctx context.Context) Inventory {
	return fromContext(ctx)
}

// NewComponent returns a new component of the given type
func NewComponent(ctx context.Context, typ string) (*wtype.Liquid, error) {
	return fromContext(ctx).NewComponent(typ)
}

// NewPlate returns a new plate of the given type
func NewPlate(ctx context.Context, typ string) (*wtype.Plate, error) {
	return fromContext(ctx).NewPlate(typ)
}
