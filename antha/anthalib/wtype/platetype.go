package wtype

type PlateTypes map[PlateTypeName]*PlateType

type PlateTypeName string

type PlateType struct {
	Name         PlateTypeName          // name of plate type, potentially including riser
	Manufacturer string                 // name of plate manufacturer
	WellShape    string                 // Name of well shape, one of "cylinder", "box", "trapezoid"
	WellH        float64                // size of well in X direction (long side of plate)
	WellW        float64                // size of well in Y direction (short side of plate)
	WellD        float64                // size of well in Z direction (vertical from plane of plate)
	MaxVol       float64                // maximum volume well can hold in microlitres
	MinVol       float64                // residual volume of well in microlitres
	BottomType   WellBottomType         // shape of well bottom, one of "flat","U", "V"
	BottomH      float64                // offset from well bottom to rest of well in mm (i.e. height of U or V - 0 if flat)
	WellX        float64                // size of well in X direction (long side of plate)
	WellY        float64                // size of well in Y direction (short side of plate)
	WellZ        float64                // size of well in Z direction (vertical from plane of plate)
	ColSize      int                    // number of wells in a column
	RowSize      int                    // number of wells in a row
	Height       float64                // size of plate in Z direction (vertical from plane of plate)
	WellXOffset  float64                // distance between adjacent well centres in X direction (long side)
	WellYOffset  float64                // distance between adjacent well centres in Y direction (short side)
	WellXStart   float64                // offset from top-left corner of plate to centre of top-leftmost well in X direction (long side)
	WellYStart   float64                // offset from top-left corner of plate to centre of top-leftmost well in Y direction (short side)
	WellZStart   float64                // offset from top of plate to well bottom
	Extra        map[string]interface{} // container for additional well properties such as constraints
}

type WellBottomType uint8

const (
	FlatWellBottom WellBottomType = iota
	UWellBottom
	VWellBottom
)

var WellBottomNames []string = []string{
	FlatWellBottom: "flat",
	UWellBottom:    "U",
	VWellBottom:    "V",
}

func (bt WellBottomType) String() string {
	return WellBottomNames[bt]
}
