package tipwastes

import (
	"fmt"
	"sync"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/laboratory/effects/id"
)

type Inventory struct {
	lock           sync.Mutex
	idGen          *id.IDGenerator
	tipwasteByType map[string]*wtype.LHTipwaste
}

func NewInventory(idGen *id.IDGenerator) *Inventory {
	inv := &Inventory{
		idGen:          idGen,
		tipwasteByType: make(map[string]*wtype.LHTipwaste),
	}
	for _, tw := range makeTipwastes(idGen) {
		if _, found := inv.tipwasteByType[tw.Type]; found {
			panic(fmt.Sprintf("tipwaste %s already added", tw.Type))
		}
		inv.tipwasteByType[tw.Type] = tw
	}
	return inv
}

func (inv *Inventory) NewTipwaste(typ string) (*wtype.LHTipwaste, error) {
	inv.lock.Lock()
	defer inv.lock.Unlock()

	if tw, found := inv.tipwasteByType[typ]; !found {
		return nil, fmt.Errorf("Unknown tip waste type: '%s'", typ)
	} else {
		return tw.Dup(inv.idGen), nil
	}
}
