package plateCache

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory/testinventory"
)

type PlateCache struct {
	inventory       *testinventory.TestInventory
	platesByType    map[string][]*wtype.Plate
	platesFromCache map[string]bool
}

func (p *PlateCache) NewPlate(typ string) (*wtype.Plate, error) {
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

	plate, err := p.inventory.NewPlate(typ)
	if err != nil {
		return nil, err
	}
	p.platesFromCache[plate.ID] = true

	return plate, nil
}

func (p *PlateCache) ReturnObject(obj interface{}) error {
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

func (p *PlateCache) IsFromCache(obj interface{}) bool {
	_, ok := p.platesFromCache[wtype.IDOf(obj)]
	return ok
}

func NewPlateCache(inv *testinventory.TestInventory) *PlateCache {
	return &PlateCache{
		inventory:       inv,
		platesByType:    make(map[string][]*wtype.Plate),
		platesFromCache: make(map[string]bool),
	}
}
