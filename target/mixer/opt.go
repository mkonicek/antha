package mixer

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/reflect"
)

var (
	defaultMaxPlates            = 4.5
	defaultMaxWells             = 278.0
	defaultResidualVolumeWeight = 1.0
	DefaultOpt                  = Opt{
		MaxPlates:            &defaultMaxPlates,
		MaxWells:             &defaultMaxWells,
		ResidualVolumeWeight: &defaultResidualVolumeWeight,
		InputPlateType:       []string{"pcrplate_skirted_riser20"},
		OutputPlateType:      []string{"pcrplate_skirted_riser20"},
		InputPlates:          []*wtype.LHPlate{},
		OutputPlates:         []*wtype.LHPlate{},
	}
)

type Opt struct {
	MaxPlates            *float64
	MaxWells             *float64
	ResidualVolumeWeight *float64
	InputPlateType       []string
	OutputPlateType      []string
	TipType              []string
	PlanningVersion      *int

	// Two methods of populating Opt.InputPlates
	InputPlateData [][]byte         // From contents of files
	InputPlates    []*wtype.LHPlate // Directly

	// Direct specification of Output plates
	OutputPlates []*wtype.LHPlate

	// Driver specific options. Semantics are not stable. Will need to be
	// revised when multi device execution is supported.
	DriverSpecificInputPreferences    []string
	DriverSpecificOutputPreferences   []string
	DriverSpecificTipPreferences      []string // Driver specific position names (e.g., position_1 or A2)
	DriverSpecificTipWastePreferences []string
	DriverSpecificWashPreferences     []string
	ModelEvaporation                  bool
}

// Merge two configs together and return the result. Values in the argument
// override those in the receiver.
func (a Opt) Merge(x *Opt) Opt {
	if x == nil {
		return a
	}
	obj, err := reflect.ShallowMerge(a, *x)
	if err != nil {
		panic(err)
	}
	return obj.(Opt)
}
