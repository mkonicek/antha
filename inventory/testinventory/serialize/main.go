package main

import (
	"encoding/json"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory/testinventory"
	"sort"
	"strings"
)

const (
	xStartOffset = 14.28
	yStartOffset = 11.24
	zStartOffset = 0.7
)

func main() {
	platesForSerializing := make([]testinventory.PlateForSerializing, 0, 1)

	thePlates := makePlates()

	plateNames := make([]string, 0, len(thePlates))

	thePlateMap := make(map[string]*wtype.Plate)

	for _, p := range thePlates {
		plateNames = append(plateNames, p.Type)
		thePlateMap[p.Type] = p
	}

	sort.Strings(plateNames)

	for _, p := range plateNames {
		plate := thePlateMap[p]
		sPlate := testinventory.NewPlateForSerializing(plate)

		if !strings.Contains(sPlate.PlateType, "FromSpec") {
			// add offset values to WellX,Y,ZStart
			sPlate = reviseWellStarts(sPlate, xStartOffset, yStartOffset, zStartOffset)
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

//		sPlate = reviseWellStarts(sPlate, xStartOffset, yStartOffset, zStartOffset)
func reviseWellStarts(sPlate testinventory.PlateForSerializing, xStartOffset, yStartOffset, zStartOffset float64) testinventory.PlateForSerializing {
	sPlate.WellXStart += xStartOffset
	sPlate.WellYStart += yStartOffset
	sPlate.WellZStart += zStartOffset

	return sPlate
}
