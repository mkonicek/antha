package plates

import (
	"fmt"
	"sync"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/laboratory/effects/id"
)

type Inventory struct {
	lock  sync.Mutex
	idGen *id.IDGenerator

	plateTypeByType wtype.PlateTypes
}

func NewInventory(idGen *id.IDGenerator) *Inventory {
	return &Inventory{
		idGen:           idGen,
		plateTypeByType: make(wtype.PlateTypes),
	}
}

func (inv *Inventory) LoadLibrary() {
	inv.SetPlateTypes(makePlateTypes(inv.idGen))
}

func (inv *Inventory) SetPlateTypes(pts wtype.PlateTypes) {
	inv.lock.Lock()
	defer inv.lock.Unlock()
	inv.plateTypeByType = pts
}

func (inv *Inventory) NewPlate(typ wtype.PlateTypeName) (*wtype.Plate, error) {
	if pt, err := inv.NewPlateType(typ); err != nil {
		return nil, err
	} else {
		return wtype.LHPlateFromType(inv.idGen, pt), nil
	}
}

func (inv *Inventory) NewPlateType(typ wtype.PlateTypeName) (*wtype.PlateType, error) {
	inv.lock.Lock()
	defer inv.lock.Unlock()

	if pt, found := inv.plateTypeByType[typ]; !found {
		return nil, fmt.Errorf("Unknown plate type: %v", typ)
	} else {
		return pt, nil
	}
}
