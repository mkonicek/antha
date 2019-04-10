package plateCache

import (
	"fmt"
	"sync"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory/plates"
)

type PlateCache struct {
	lock sync.Mutex

	inventory       *plates.Inventory
	platesByType    map[wtype.PlateTypeName][]*wtype.Plate
	platesFromCache map[string]struct{}
}

func NewPlateCache(inv *plates.Inventory) *PlateCache {
	return &PlateCache{
		inventory:       inv,
		platesByType:    make(map[wtype.PlateTypeName][]*wtype.Plate),
		platesFromCache: make(map[string]struct{}),
	}
}

func (pc *PlateCache) NewPlate(typ wtype.PlateTypeName) (*wtype.Plate, error) {
	pc.lock.Lock()
	defer pc.lock.Unlock()

	plates, found := pc.platesByType[typ]
	if !found {
		plates = make([]*wtype.Plate, 0, 1)
	}

	if len(plates) == 0 {
		plate, err := pc.inventory.NewPlate(typ)
		if err != nil {
			return nil, err
		}
		pc.platesFromCache[plate.ID] = struct{}{}
		plates = append(plates, plate)
	}

	plate := plates[0]
	pc.platesByType[typ] = plates[1:]
	return plate, nil
}

func (pc *PlateCache) ReturnPlate(plate *wtype.Plate) error {
	if !pc.IsFromCache(plate) {
		return fmt.Errorf("cannot return non-cache plate %v", plate)
	}

	pc.lock.Lock()
	defer pc.lock.Unlock()

	plate.Clean()

	typ := plate.Type

	if plates, found := pc.platesByType[typ]; !found {
		panic(fmt.Errorf("Impossible: plate is from cache, but plates slice not found! %v", typ))
	} else {
		pc.platesByType[typ] = append(plates, plate)
	}

	return nil
}

func (pc *PlateCache) IsFromCache(plate *wtype.Plate) bool {
	pc.lock.Lock()
	defer pc.lock.Unlock()

	_, found := pc.platesFromCache[plate.ID]
	return found
}
