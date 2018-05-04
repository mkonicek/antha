package plateCache

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/cache"
)

type plateCache struct {
	platesByType map[string][]*wtype.LHPlate
}

func (p *plateCache) NewComponent(ctx context.Context, name string) (*wtype.LHComponent, error) {
	return inventory.NewComponent(ctx, name)
}

func (p *plateCache) NewTipbox(ctx context.Context, typ string) (*wtype.LHTipbox, error) {
	return inventory.NewTipbox(ctx, typ)
}

func (p *plateCache) NewTipwaste(ctx context.Context, typ string) (*wtype.LHTipwaste, error) {
	return inventory.NewTipwaste(ctx, typ)
}

func (p *plateCache) NewPlate(ctx context.Context, typ string) (*wtype.LHPlate, error) {
	fmt.Printf("plateCache: Getting Plate type \"%s\"\n", typ)
	return inventory.NewPlate(ctx, typ)
}

func (p *plateCache) ReturnObject(ctx context.Context, obj interface{}) error {
	fmt.Printf("plateCache: Returning %s type \"%s\"\n", wtype.ClassOf(obj), wtype.TypeOf(obj))
	return nil
}

// NewContext creates a new plateCache context
func NewContext(ctx context.Context) context.Context {
	pc := &plateCache{
		platesByType: make(map[string][]*wtype.LHPlate),
	}

	return cache.NewContext(ctx, pc)
}
