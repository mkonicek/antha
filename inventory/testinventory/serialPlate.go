package testinventory

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type plateForSerializing struct {
	PlateName    string
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
	// "Nunc96DeepWell", "Unknown", 8, 12, makePlateCoords(43.6), nunc96deepwell, 9, 9, -1.0, 0.0, 6.5
	ColSize     int
	RowSize     int
	Height      float64
	WellXOffset float64
	WellYOffset float64
	WellXStart  float64
	WellYStart  float64
	WellZStart  float64
	Special     bool
}
