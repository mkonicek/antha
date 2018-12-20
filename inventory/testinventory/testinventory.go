package testinventory

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

var ErrUnknownType = errors.New("unknown type requested from inventory")

const (
	// WaterType is the component type of water
	WaterType = "water"
)

type TestInventory struct {
	componentByName map[string]*wtype.Liquid
	plateByType     map[string]PlateForSerializing
	tipboxByType    map[string]*wtype.LHTipbox
	tipwasteByType  map[string]*wtype.LHTipwaste
}

func (i *TestInventory) NewComponent(name string) (*wtype.Liquid, error) {
	c, ok := i.componentByName[name]
	if !ok {
		return nil, fmt.Errorf("%s: invalid solution: %s", ErrUnknownType, name)
	}
	// Cp is required here to ensure component IDs are unique
	return c.Cp(), nil
}

func (i *TestInventory) NewPlate(typ string) (*wtype.Plate, error) {
	p, ok := i.plateByType[typ]
	if !ok {
		return nil, fmt.Errorf("%s: invalid plate: %s", ErrUnknownType, typ)
	}
	return p.LHPlate(), nil
}
func (i *TestInventory) NewTipbox(typ string) (*wtype.LHTipbox, error) {
	tb, ok := i.tipboxByType[typ]
	if !ok {
		return nil, ErrUnknownType
	}
	return tb.Dup(), nil
}

func (i *TestInventory) NewTipwaste(typ string) (*wtype.LHTipwaste, error) {
	tw, ok := i.tipwasteByType[typ]
	if !ok {
		return nil, ErrUnknownType
	}
	return tw.Dup(), nil
}

// NewContext creates a new test inventory context
func NewInventory() *TestInventory {
	inv := &TestInventory{
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

	return inv
}

// GetTipboxes returns the tipboxes in a test inventory context
func (inv *TestInventory) GetTipboxes() []*wtype.LHTipbox {
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
func (inv *TestInventory) GetPlates() []*wtype.Plate {
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
func (inv *TestInventory) GetComponents() []*wtype.Liquid {
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
