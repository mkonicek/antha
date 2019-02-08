package tipboxes

import (
	"fmt"
	"sync"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/laboratory/effects/id"
)

type Inventory struct {
	lock  sync.Mutex
	idGen *id.IDGenerator

	tipboxByType map[string]*wtype.LHTipbox
}

func NewInventory(idGen *id.IDGenerator) *Inventory {
	return &Inventory{
		idGen:        idGen,
		tipboxByType: make(map[string]*wtype.LHTipbox),
	}
}

func (inv *Inventory) LoadLibrary() {
	inv.lock.Lock()
	defer inv.lock.Unlock()
	for _, tb := range makeTipboxes(inv.idGen) {
		if _, found := inv.tipboxByType[tb.Type]; found {
			panic(fmt.Sprintf("tipbox %s already added", tb.Type))
		}
		if _, found := inv.tipboxByType[tb.Tiptype.Type]; found {
			panic(fmt.Sprintf("tipbox %s already added", tb.Tiptype.Type))
		}
		inv.tipboxByType[tb.Type] = tb
		inv.tipboxByType[tb.Tiptype.Type] = tb
	}
}

func (inv *Inventory) NewTipbox(typ string) (*wtype.LHTipbox, error) {
	if tb, err := inv.FetchTipbox(typ); err != nil {
		return nil, err
	} else {
		return tb.Dup(inv.idGen), nil
	}
}

// Same as NewTipbox but does not Dup the result. Only use for
// read-only access to the tip box inventory.
func (inv *Inventory) FetchTipbox(typ string) (*wtype.LHTipbox, error) {
	inv.lock.Lock()
	defer inv.lock.Unlock()

	if tb, found := inv.tipboxByType[typ]; !found {
		return nil, fmt.Errorf("Unknown tip box type: '%s'", typ)
	} else {
		return tb, nil
	}
}
