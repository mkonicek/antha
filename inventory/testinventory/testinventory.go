package testinventory

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
	"sync"
)

type TestInventory struct {
	componentByName map[string]*wtype.Liquid
	plateByType     map[string]*wtype.Plate
	tipboxByType    map[string]*wtype.LHTipbox
	tipwasteByType  map[string]*wtype.LHTipwaste
	lock            *sync.Mutex
}

func NewTestInventory() *TestInventory {
	return &TestInventory{
		componentByName: getComponentsByType(),
		plateByType:     getPlatesByType(),
		tipboxByType:    getTipboxesByType(),
		tipwasteByType:  getTipwastesByType(),
		lock:            &sync.Mutex{},
	}
}

func (i *TestInventory) NewComponent(name string) (*wtype.Liquid, error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	c, ok := i.componentByName[name]
	if !ok {
		return nil, fmt.Errorf("%s: invalid solution: %s", inventory.ErrUnknownType, name)
	}
	// Cp is required here to ensure component IDs are unique
	return c.Cp(), nil
}

func (i *TestInventory) Components() []*wtype.Liquid {
	i.lock.Lock()
	defer i.lock.Unlock()

	ret := make([]*wtype.Liquid, 0, len(i.componentByName))
	for _, c := range i.componentByName {
		ret = append(ret, c.Cp())
	}
	return ret
}

func (i *TestInventory) NewPlate(typ string) (*wtype.Plate, error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	if p, ok := i.plateByType[typ]; !ok {
		return nil, fmt.Errorf("%s: invalid plate: %s", inventory.ErrUnknownType, typ)
	} else {
		return p.Dup(), nil
	}
}

func (i *TestInventory) Plates() []*wtype.Plate {
	i.lock.Lock()
	defer i.lock.Unlock()

	ret := make([]*wtype.Plate, 0, len(i.plateByType))
	for _, plate := range i.plateByType {
		ret = append(ret, plate.Dup())
	}
	return ret
}

func (i *TestInventory) NewTipbox(typ string) (*wtype.LHTipbox, error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	if p, ok := i.tipboxByType[typ]; !ok {
		return nil, fmt.Errorf("%s: invalid tipbox: %s", inventory.ErrUnknownType, typ)
	} else {
		return p.Dup(), nil
	}
}

func (i *TestInventory) NewTipwaste(typ string) (*wtype.LHTipwaste, error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	if p, ok := i.tipwasteByType[typ]; !ok {
		return nil, fmt.Errorf("%s: invalid tipwaste: %s", inventory.ErrUnknownType, typ)
	} else {
		return p.Dup(), nil
	}
}

// NewContext creates a new test inventory context
func NewContext(ctx context.Context) context.Context {
	return inventory.NewContext(ctx, NewTestInventory())
}

// invForTest a single inventory to be shared for testing, threadsafe and read only
var invForTest *TestInventory

func GetInventoryForTest() *TestInventory {
	if invForTest == nil {
		invForTest = NewTestInventory()
	}
	return invForTest
}

func NewContextForTest(ctx context.Context) context.Context {
	return inventory.NewContext(ctx, GetInventoryForTest())
}

// getTipboxesByType returns the test tipboxes
func getTipboxesByType() map[string]*wtype.LHTipbox {
	tbs := makeTipboxes()
	ret := make(map[string]*wtype.LHTipbox, len(tbs))
	for _, tb := range tbs {
		if _, seen := ret[tb.Type]; seen {
			panic(fmt.Sprintf("tipbox %s already added", tb.Type))
		} else if _, seen := ret[tb.Tiptype.Type]; seen {
			panic(fmt.Sprintf("tipbox %s already added", tb.Tiptype.Type))
		}
		ret[tb.Type] = tb
		ret[tb.Tiptype.Type] = tb
	}
	return ret
}

func getPlatesByType() map[string]*wtype.Plate {
	if serialPlateArr, err := getPlatesFromSerial(); err != nil {
		panic(err)
	} else {
		ret := make(map[string]*wtype.LHPlate, len(serialPlateArr))

		for _, p := range serialPlateArr {
			if _, seen := ret[p.PlateType]; seen {
				panic(fmt.Sprintf("plate %s already added", p.PlateType))
			}
			ret[p.PlateType] = p.LHPlate()
		}
		return ret
	}
}

// getComponentsByType returns the test components
func getComponentsByType() map[string]*wtype.Liquid {
	components := makeComponents()
	ret := make(map[string]*wtype.Liquid, len(components))
	for _, c := range components {
		if _, seen := ret[c.CName]; seen {
			panic(fmt.Sprintf("component %s already added", c.CName))
		}
		ret[c.CName] = c
	}
	return ret
}

// getTipwastesByType returns the test tipwastes
func getTipwastesByType() map[string]*wtype.LHTipwaste {
	tipwastes := makeTipwastes()
	ret := make(map[string]*wtype.LHTipwaste, len(tipwastes))
	for _, tw := range tipwastes {
		if _, seen := ret[tw.Type]; seen {
			panic(fmt.Sprintf("tipwaste %s already added", tw.Type))
		}
		ret[tw.Type] = tw
	}
	return ret
}

func getPlatesFromSerial() ([]PlateForSerializing, error) {
	var pltArr []PlateForSerializing

	err := json.Unmarshal(plateBytes, &pltArr)

	if err != nil {
		return nil, err
	}

	return pltArr, nil
}
