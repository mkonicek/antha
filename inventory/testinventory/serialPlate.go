package testinventory

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

var (
	vunit = "ul"
	lunit = "mm"
)

type PlateForSerializing struct {
	PlateType    string
	Manufacturer string
	// shape "box", "mm", 8.2, 8.2, 41.3 	-- defines well shape
	WellShape string
	WellH     float64
	WellW     float64
	WellD     float64
	// well  "ul", 2000, 420, swshp, wtype.VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm"
	MaxVol     float64
	MinVol     float64
	BottomType wtype.WellBottomType
	BottomH    float64
	WellX      float64
	WellY      float64
	WellZ      float64
	// "Nunc96DeepWell", "Unknown", 8, 12, makePlateCoords(43.6), nunc96deepwell, 9, 9, -1.0, 0.0, 6.5
	ColSize     int
	RowSize     int
	Height      float64
	WellXOffset float64
	WellYOffset float64
	WellXStart  float64
	WellYStart  float64
	WellZStart  float64
	Extra       map[string]interface{}
}

func (pt PlateForSerializing) LHPlate() *wtype.LHPlate {
	newWellShape := wtype.NewShape(pt.WellShape, lunit, pt.WellH, pt.WellW, pt.WellD)

	newWelltype := wtype.NewLHWell(vunit, pt.MaxVol, pt.MinVol, newWellShape, pt.BottomType, pt.WellX, pt.WellY, pt.WellZ, pt.BottomH, lunit)

	plate := wtype.NewLHPlate(pt.PlateType, pt.Manufacturer, pt.ColSize, pt.RowSize, makePlateCoords(pt.Height), newWelltype, pt.WellXOffset, pt.WellYOffset, pt.WellXStart, pt.WellYStart, pt.WellZStart)

	plate.Welltype.Extra = pt.Extra

	return plate
}
func makePlateCoords(height float64) wtype.Coordinates {
	//standard X/Y size for 96 well plates
	return wtype.Coordinates{X: 127.76, Y: 85.48, Z: height}
}
