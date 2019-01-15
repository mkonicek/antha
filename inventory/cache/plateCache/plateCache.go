package plateCache

import (
	"context"
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/cache"
)

type plateCache struct {
	platesByType    map[string][]*wtype.Plate
	platesFromCache map[string]bool
	inv             inventory.Inventory
}

func (p *plateCache) NewComponent(name string) (*wtype.Liquid, error) {
	return p.inv.NewComponent(name)
}

func (p *plateCache) NewPlate(typ string) (*wtype.Plate, error) {
	plateList, ok := p.platesByType[typ]
	if !ok {
		plateList = make([]*wtype.Plate, 0)
		p.platesByType[typ] = plateList
	}

	if len(plateList) > 0 {
		plate := plateList[0]
		p.platesByType[typ] = plateList[1:]
		return plate, nil
	}

	plate, err := p.inv.NewPlate(typ)
	if err != nil {
		return nil, err
	}
	p.platesFromCache[plate.ID] = true

	return plate, nil
}

func (p *plateCache) ReturnObject(obj interface{}) error {
	if !p.IsFromCache(obj) {
		return fmt.Errorf("cannont return non cache object %s", wtype.NameOf(obj))
	}
	plate, ok := obj.(*wtype.Plate)
	if !ok {
		return fmt.Errorf("cannot return object class %s to plate cache", wtype.ClassOf(obj))
	}
	plate.Clean()

	typ := wtype.TypeOf(plate)

	_, ok = p.platesByType[typ]
	if !ok {
		p.platesByType[typ] = make([]*wtype.Plate, 0, 1)
	}

	p.platesByType[typ] = append(p.platesByType[typ], plate)

	return nil
}

func (p *plateCache) IsFromCache(obj interface{}) bool {
	_, ok := p.platesFromCache[wtype.IDOf(obj)]
	return ok
}

// NewContext creates a new plateCache context
func NewContext(ctx context.Context) context.Context {
	pc := &plateCache{
		platesByType:    make(map[string][]*wtype.Plate),
		platesFromCache: make(map[string]bool),
		inv:             inventory.GetInventory(ctx),
	}

	return cache.NewContext(ctx, pc)
}
