package wtype

import (
	"fmt"
	"strconv"
	"strings"
)

// These types are used by composer and directly define the definition
// of plate types in the workflow.

type PlateTypes map[PlateTypeName]*PlateType

type PlateTypes2 []PlateType2

type PlateTypeName string

type PlateType2 struct {
	ID                   string  `json:"id"`
	Type                 string  `json:"type"`
	Manufacturer         string  `json:"manufacturer"`
	Name                 string  `json:"name"`
	CatalogNumber        string  `json:"catalog_number"`
	Accessory            string  `json:"accessory"`
	Format               string  `json:"format"`
	Function             string  `json:"function"`
	Columns              string  `json:"columns"`
	Rows                 string  `json:"rows"`
	WellType             string  `json:"well_type"`
	WellBottom           string  `json:"well_bottom"`
	WellVolumeUl         float64 `json:"well_volume_ul"`
	WellVolumeResidualUl float64 `json:"well_volume_residual_ul"`
	Revision             string  `json:"revision"`
	WellDimensionMM      struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	} `json:"well_dimension_mm"`
	WellBottomZOffsetMM float64 `json:"well_bottom_z_offset_mm"`
	WellOffsetMM        struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"well_offset_mm"`
	WellStartMM struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	} `json:"well_start_mm"`
	DimensionMM struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	} `json:"dimension_mm"`
	CatalogURL string `json:"catalog_url"`
	HasAdaptor bool   `json:"has_adaptor"`
}

type PlateType struct {
	Name          PlateTypeName          // name of plate type, potentially including riser
	Manufacturer  string                 // name of plate manufacturer
	CatalogNumber string                 // catalog number
	WellShape     string                 // Name of well shape, one of "cylinder", "box", "trapezoid"
	WellH         float64                // size of well in X direction (long side of plate)
	WellW         float64                // size of well in Y direction (short side of plate)
	WellD         float64                // size of well in Z direction (vertical from plane of plate)
	MaxVol        float64                // maximum volume well can hold in microlitres
	MinVol        float64                // residual volume of well in microlitres
	BottomType    WellBottomType         // shape of well bottom, one of "flat","U", "V"
	BottomH       float64                // offset from well bottom to rest of well in mm (i.e. height of U or V - 0 if flat)
	WellX         float64                // size of well in X direction (long side of plate)
	WellY         float64                // size of well in Y direction (short side of plate)
	WellZ         float64                // size of well in Z direction (vertical from plane of plate)
	ColSize       int                    // number of wells in a column
	RowSize       int                    // number of wells in a row
	Height        float64                // size of plate in Z direction (vertical from plane of plate)
	WellXOffset   float64                // distance between adjacent well centres in X direction (long side)
	WellYOffset   float64                // distance between adjacent well centres in Y direction (short side)
	WellXStart    float64                // offset from top-left corner of plate to centre of top-leftmost well in X direction (long side)
	WellYStart    float64                // offset from top-left corner of plate to centre of top-leftmost well in Y direction (short side)
	WellZStart    float64                // offset from top of plate to well bottom
	Extra         map[string]interface{} // container for additional well properties such as constraints
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

//FIXME !!
func WellBottomTypeFromString(in string) WellBottomType {
	if strings.Contains(in, "flat") {
		return 0
	}
	if strings.Contains(in, "U") {
		return 1
	}
	if strings.Contains(in, "V") {
		return 2
	}
	// default
	return 0
}

func (bt WellBottomType) String() string {
	return WellBottomNames[bt]
}

func (in PlateType2) ConvertToPlateType() PlateType {
	var out PlateType
	// Name is Type obviously -- oh man !
	out.Name = PlateTypeName(in.Type)
	out.Manufacturer = in.Manufacturer
	out.CatalogNumber = in.CatalogNumber
	out.WellShape = in.WellType
	out.WellH = in.WellDimensionMM.X
	out.WellW = in.WellDimensionMM.Y
	out.WellD = in.WellDimensionMM.Z
	out.MaxVol = in.WellVolumeUl
	out.MinVol = in.WellVolumeResidualUl
	out.BottomType = WellBottomTypeFromString(in.WellBottom)
	//check BottomH
	out.BottomH = in.WellBottomZOffsetMM
	out.WellX = in.WellDimensionMM.X
	out.WellY = in.WellDimensionMM.Y
	out.WellZ = in.WellDimensionMM.Z
	// yeah maybe need error handling
	out.ColSize, _ = strconv.Atoi(in.Columns)
	out.RowSize, _ = strconv.Atoi(in.Rows)
	// check Height
	out.Height = in.DimensionMM.Z
	out.WellXOffset = in.WellOffsetMM.X
	out.WellYOffset = in.WellOffsetMM.Y
	out.WellXStart = in.WellStartMM.X
	out.WellYStart = in.WellStartMM.Y
	out.WellZStart = in.WellStartMM.Z
	// what about extra stuff ?
	return out
}

func (a *PlateTypes) Merge(b PlateTypes) error {
	if *a == nil {
		*a = make(PlateTypes)
	}
	aMap := *a
	// May need revising: currently we error if there's any
	// overlap. Equality between PlateTypes can't be based on simple
	// structural equality due to the Extra field being a map.
	for ptn, pt := range b {
		if _, found := aMap[ptn]; found {
			return fmt.Errorf("Cannot merge: PlateType '%v' is redefined", ptn)
		}
		aMap[ptn] = pt
	}
	return nil
}

func (p PlateTypes) Validate() error {
	// not a validation step really, more just consistency
	for ptn, pt := range p {
		pt.Name = ptn
	}
	return nil
}
