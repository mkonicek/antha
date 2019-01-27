package components

import (
	"fmt"
	"sync"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/laboratory/effects/id"
)

const (
	// WaterType is the component type of water
	WaterType = "water"
)

type ErrComponentNotExist struct {
	err error
}

func (e ErrComponentNotExist) Error() string {
	return fmt.Sprintf("Component does not exist: %s", e.err.Error())
}

func IsComponentNotExist(err error) bool {
	_, ok := err.(ErrComponentNotExist)
	return ok
}

type Inventory struct {
	lock            sync.Mutex
	idGen           *id.IDGenerator
	componentByName map[string]*wtype.Liquid
}

func NewInventory(idGen *id.IDGenerator) *Inventory {
	inv := &Inventory{
		idGen:           idGen,
		componentByName: make(map[string]*wtype.Liquid),
	}

	for _, liq := range makeComponents(idGen) {
		if _, found := inv.componentByName[liq.CName]; found {
			panic(fmt.Sprintf("component %s already added", liq.CName))
		}
		inv.componentByName[liq.CName] = liq
	}

	return inv
}

func (inv *Inventory) NewComponent(name string) (*wtype.Liquid, error) {
	inv.lock.Lock()
	defer inv.lock.Unlock()

	if liq, found := inv.componentByName[name]; !found {
		return nil, ErrComponentNotExist{err: fmt.Errorf("invalid solution: '%s'", name)}
	} else {
		// Cp is required here to ensure component IDs are unique
		return liq.Cp(inv.idGen), nil
	}
}
