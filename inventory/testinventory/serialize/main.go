package main

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory/testinventory"
	"sort"
)

func main() {
	platesForSerializing := make([]testinventory.PlateForSerializing, 0, 1)

	thePlates := makePlates()

	plateNames := make([]string, 0, len(thePlates))

	thePlateMap := make(map[string]*wtype.LHPlate)

	for _, p := range thePlates {
		plateNames = append(plateNames, p.Type)
		thePlateMap[p.Type] = p
	}

	sort.Strings(plateNames)

	for _, p := range plateNames {
		plate := thePlateMap[p]
		sPlate := testinventory.PlateForSerializing{
			PlateType:    p,
			Manufacturer: plate.Mnfr,
			WellShape:    plate.Welltype.Shape().ShapeName,
			WellH:        plate.Welltype.Shape().H,
			WellW:        plate.Welltype.Shape().W,
			WellD:        plate.Welltype.Shape().D,
			MaxVol:       plate.Welltype.MaxVol,
			MinVol:       plate.Welltype.Rvol,
			BottomType:   plate.Welltype.Bottom,
			BottomH:      plate.Welltype.Bottomh,
			WellX:        plate.Welltype.Bounds.Size.X,
			WellY:        plate.Welltype.Bounds.Size.Y,
			WellZ:        plate.Welltype.Bounds.Size.Z,
			ColSize:      plate.WellsY(),
			RowSize:      plate.WellsX(),
			Height:       plate.Height(),
			WellXOffset:  plate.WellXOffset,
			WellYOffset:  plate.WellYOffset,
			WellXStart:   plate.WellXStart,
			WellYStart:   plate.WellYStart,
			WellZStart:   plate.WellZStart,
			Extra:        plate.Welltype.Extra,
		}
		platesForSerializing = append(platesForSerializing, sPlate)
	}

	s, err := json.MarshalIndent(platesForSerializing, "", " ")

	if err != nil {
		panic(fmt.Sprint("serialize error ", err))
	}

	fmt.Println("package testinventory")
	fmt.Println()
	fmt.Println("var plateBytes = []byte(`")
	fmt.Println(string(s))
	fmt.Println("`)")
}
