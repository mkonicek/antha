package testinventory

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
)

type testInventory struct {
	componentByName map[string]*wtype.Liquid
	plateByType     map[string]PlateForSerializing
	tipboxByType    map[string]*wtype.LHTipbox
	tipwasteByType  map[string]*wtype.LHTipwaste
}

func (i *testInventory) NewComponent(ctx context.Context, name string) (*wtype.Liquid, error) {
	c, ok := i.componentByName[name]
	if !ok {
		return nil, fmt.Errorf("%s: invalid solution: %s", inventory.ErrUnknownType, name)
	}
	// Cp is required here to ensure component IDs are unique
	return c.Cp(), nil
}

func (i *testInventory) NewPlate(ctx context.Context, typ string) (*wtype.Plate, error) {
	p, ok := i.plateByType[typ]
	if !ok {
		return nil, fmt.Errorf("%s: invalid plate: %s", inventory.ErrUnknownType, typ)
	}
	return p.LHPlate(), nil
}
func (i *testInventory) NewTipbox(ctx context.Context, typ string) (*wtype.LHTipbox, error) {
	tb, ok := i.tipboxByType[typ]
	if !ok {
		return nil, inventory.ErrUnknownType
	}
	return tb.Dup(), nil
}

func (i *testInventory) NewTipwaste(ctx context.Context, typ string) (*wtype.LHTipwaste, error) {
	tw, ok := i.tipwasteByType[typ]
	if !ok {
		return nil, inventory.ErrUnknownType
	}
	return tw.Dup(), nil
}

func (i *testInventory) XXXGetPlates(ctx context.Context) ([]*wtype.Plate, error) {
	plates := GetPlates(ctx)
	return plates, nil
}

// NewContext creates a new test inventory context
func NewContext(ctx context.Context) context.Context {
	inv := &testInventory{
		componentByName: make(map[string]*wtype.Liquid),
		plateByType:     make(map[string]PlateForSerializing),
		tipboxByType:    make(map[string]*wtype.LHTipbox),
		tipwasteByType:  make(map[string]*wtype.LHTipwaste),
	}

	for _, c := range makeComponents() {
		if _, seen := inv.componentByName[c.CName]; seen {
			panic(fmt.Sprintf("component %s already added", c.CName))
		}
		inv.componentByName[c.CName] = c
	}

	serialPlateArr, err := getPlatesFromSerial()

	if err != nil {
		panic(err)
	}

	for _, p := range serialPlateArr {
		if _, seen := inv.plateByType[p.PlateType]; seen {
			panic(fmt.Sprintf("plate %s already added", p.PlateType))
		}
		inv.plateByType[p.PlateType] = p
	}

	for _, tb := range makeTipboxes() {
		if _, seen := inv.tipboxByType[tb.Type]; seen {
			panic(fmt.Sprintf("tipbox %s already added", tb.Type))
		}
		if _, seen := inv.tipboxByType[tb.Tiptype.Type]; seen {
			panic(fmt.Sprintf("tipbox %s already added", tb.Tiptype.Type))
		}
		inv.tipboxByType[tb.Type] = tb
		inv.tipboxByType[tb.Tiptype.Type] = tb
	}

	for _, tw := range makeTipwastes() {
		if _, seen := inv.tipwasteByType[tw.Type]; seen {
			panic(fmt.Sprintf("tipwaste %s already added", tw.Type))
		}
		inv.tipwasteByType[tw.Type] = tw
	}

	return inventory.NewContext(ctx, inv)
}

// GetTipboxes returns the tipboxes in a test inventory context
func GetTipboxes(ctx context.Context) []*wtype.LHTipbox {
	inv := inventory.GetInventory(ctx).(*testInventory)
	var tbs []*wtype.LHTipbox
	for _, tb := range inv.tipboxByType {
		tbs = append(tbs, tb)
	}

	sort.Slice(tbs, func(i, j int) bool {
		return tbs[i].Type < tbs[j].Type
	})

	return tbs
}

// GetPlates returns the plates in a test inventory context
func GetPlates(ctx context.Context) []*wtype.Plate {
	inv := inventory.GetInventory(ctx).(*testInventory)
	var ps []*wtype.Plate
	for _, p := range inv.plateByType {
		ps = append(ps, p.LHPlate())
	}

	sort.Slice(ps, func(i, j int) bool {
		return ps[i].Type < ps[j].Type
	})

	return ps
}

// GetComponents returns the components in a test inventory context
func GetComponents(ctx context.Context) []*wtype.Liquid {
	inv := inventory.GetInventory(ctx).(*testInventory)
	var cs []*wtype.Liquid
	for _, c := range inv.componentByName {
		cs = append(cs, c)
	}

	sort.Slice(cs, func(i, j int) bool {
		return cs[i].Type < cs[j].Type
	})

	return cs
}

func getPlatesFromSerial() ([]PlateForSerializing, error) {
	var pltArr []PlateForSerializing

	err := json.Unmarshal(plateBytes, &pltArr)

	if err != nil {
		return nil, err
	}

	return pltArr, nil
}
