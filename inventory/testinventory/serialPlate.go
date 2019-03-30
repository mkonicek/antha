package testinventory

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

var (
	vunit = "ul" // unit used for volumes specified in this structure
	lunit = "mm" //  unit used for lengths specified in this structure
)

// PlateForSerializing contains required measurements and other properties for creating a plate of a given type
// but doesn't represent each well individually
type PlateForSerializing struct {
	PlateType    string // name of plate type, potentially including tiser
	Manufacturer string // name of plate manufacturer
	// shape "box", "mm", 8.2, 8.2, 41.3 	-- defines well shape
	WellShape string  // Name of well shape
	WellH     float64 // size of well in X direction (long side of plate)
	WellW     float64 // size of well in Y direction (short side of plate)
	WellD     float64 // size of well in Z direction (vertical from plane of plate)
	// well  "ul", 2000, 420, swshp, wtype.VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm"
	MaxVol     float64              // maximum volume well can hold in microlitres
	MinVol     float64              // residual volume of well in microlitres
	BottomType wtype.WellBottomType // shape of well bottom, one of "flat","U", "V"
	BottomH    float64              // offset from well bottom to rest of well in mm (i.e. height of U or V - 0 if flat)
	WellX      float64              // size of well in X direction (long side of plate)
	WellY      float64              // size of well in Y direction (short side of plate)
	WellZ      float64              // size of well in Z direction (vertical from plane of plate)
	// "Nunc96DeepWell", "Unknown", 8, 12, makePlateCoords(43.6), nunc96deepwell, 9, 9, -1.0, 0.0, 6.5
	ColSize     int                    // number of wells in a column
	RowSize     int                    // number of wells in a row
	Height      float64                // size of plate in Z direction (vertical from plane of plate)
	WellXOffset float64                // distance between adjacent well centres in X direction (long side)
	WellYOffset float64                // distance between adjacent well centres in Y direction (short side)
	WellXStart  float64                // offset from top-left corner of plate to centre of top-leftmost well in X direction (long side)
	WellYStart  float64                // offset from top-left corner of plate to centre of top-leftmost well in Y direction (short side)
	WellZStart  float64                // offset from top of plate to well bottom
	Extra       map[string]interface{} // container for additional well properties such as constraints
}

// NewPlateForSerializing get an easily serializable version of the plate
func NewPlateForSerializing(plate *wtype.LHPlate) PlateForSerializing {
	return PlateForSerializing{
		PlateType:    plate.Type,
		Manufacturer: plate.Mnfr,
		WellShape:    plate.Welltype.Shape().Type.String(),
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
}

// LHPlate returns an initialized, empty, LHPlate of the type corresponding to this PlateForSerializing
func (pt PlateForSerializing) LHPlate() *wtype.Plate {
	newWellShape := wtype.NewShape(wtype.ShapeTypeFromName(pt.WellShape), lunit, pt.WellH, pt.WellW, pt.WellD)

	newWelltype := wtype.NewLHWell(vunit, pt.MaxVol, pt.MinVol, newWellShape, pt.BottomType, pt.WellX, pt.WellY, pt.WellZ, pt.BottomH, lunit)

	plate := wtype.NewLHPlate(pt.PlateType, pt.Manufacturer, pt.ColSize, pt.RowSize, makePlateCoords(pt.Height), newWelltype, pt.WellXOffset, pt.WellYOffset, pt.WellXStart, pt.WellYStart, pt.WellZStart)

	plate.Welltype.Extra = pt.Extra

	return plate
}
func makePlateCoords(height float64) wtype.Coordinates3D {
	//standard X/Y size for 96 well plates
	return wtype.Coordinates3D{X: 127.76, Y: 85.48, Z: height}
}
