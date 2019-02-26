package mixer

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/meta"
)

var (
	defaultMaxPlates            = 4.5
	defaultMaxWells             = 278.0
	defaultResidualVolumeWeight = 1.0

	// DefaultOpt is the default Mixer Opt
	DefaultOpt = Opt{
		MaxPlates:            &defaultMaxPlates,
		MaxWells:             &defaultMaxWells,
		ResidualVolumeWeight: &defaultResidualVolumeWeight,
		PlanningVersion:      "ep2",
		LegacyVolume:         true,
		FixVolumes:           true,
	}
)

// Opt are options for a Mixer
type Opt struct {
	MaxPlates            *float64 `json:"maxPlates,omitempty"`
	MaxWells             *float64 `json:"maxWells,omitempty"`
	ResidualVolumeWeight *float64 `json:"residualVolumeWeight,omitempty"`
	InputPlateTypes      []string `json:"inputPlateTypes,omitempty"`
	OutputPlateTypes     []string `json:"outputPlateTypes,omitempty"`
	TipTypes             []string `json:"tipTypes,omitempty"`
	PlanningVersion      string   `json:"executionPlannerVersion,omitempty"`

	// Two methods of populating input plates
	InputPlateData [][]byte       `json:"inputPlateData,omitempty"` // From contents of files
	InputPlates    []*wtype.Plate `json:"inputPlates,omitempty"`    // Directly

	// Direct specification of output plates
	OutputPlates []*wtype.Plate `json:"outputPlates,omitempty"`

	// Specify file name in the instruction stream of any driver generated file
	DriverOutputFileName string `json:"driverOutputFileName,omitempty"`

	// Driver specific options. Semantics are not stable. Will need to be
	// revised when multi device execution is supported.
	DriverSpecificInputPreferences    []string `json:"driverSpecificInputPreferences,omitempty"`
	DriverSpecificOutputPreferences   []string `json:"driverSpecificOutputPreferences,omitempty"`
	DriverSpecificTipPreferences      []string `json:"driverSpecificTipPreferences,omitempty"` // Driver specific position names (e.g., position_1 or A2)
	DriverSpecificTipWastePreferences []string `json:"driverSpecificTipWastePreferences,omitempty"`
	DriverSpecificWashPreferences     []string `json:"driverSpecificWashPreferences,omitempty"`

	ModelEvaporation         bool `json:"modelEvaporation"`
	OutputSort               bool `json:"outputSort"`
	PrintInstructions        bool `json:"printInstructions"`
	UseDriverTipTracking     bool `json:"useDriverTipTracking"`
	LegacyVolume             bool `json:"legacyVolume"`             // Don't track volumes for intermediates
	FixVolumes               bool `json:"fixVolumes"`               // Aim to revise requested volumes to service requirements
	IgnorePhysicalSimulation bool `json:"ignorePhysicalSimulation"` //ignore errors in physical simulation

	// Two ways to set user liquid policies rule set
	CustomPolicyData    map[string]wtype.LHPolicy `json:"customPolicyData,omitempty"`    // Set rule set from policies
	CustomPolicyRuleSet *wtype.LHPolicyRuleSet    `json:"customPolicyRuleSet,omitempty"` // Directly
}

// Merge two configs together and return the result. Values in the argument
// override those in the receiver.
func (a Opt) Merge(x *Opt) Opt {
	if x == nil {
		return a
	}

	obj, err := meta.ShallowMerge(a, *x)
	if err != nil {
		panic(err)
	}
	return obj.(Opt)
}
