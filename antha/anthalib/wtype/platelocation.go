package wtype

import "strings"

type PlateLocation struct {
	ID     string
	Coords WellCoords
}

func (pc PlateLocation) ToString() string {
	return pc.ID + ":" + pc.Coords.FormatA1()
}

func PlateLocationFromString(s string) PlateLocation {
	pl := PlateLocation{}
	tx := strings.Split(s, ":")

	if len(tx) != 2 {
		return pl
	}

	return PlateLocation{tx[0], MakeWellCoords(tx[1])}
}

func (pc PlateLocation) Equals(opc PlateLocation) bool {
	if !(pc.ID == opc.ID && pc.Coords.FormatA1() == opc.Coords.FormatA1()) {
		return false
	}
	return true
}
