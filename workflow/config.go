package workflow

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

type Config struct {
	GilsonPipetMax GilsonPipetMaxConfig `json:"GilsonPipetMax"`
	GlobalMixer    GlobalMixerConfig    `json:"GlobalMixer"`
}

type DeviceInstanceID string

type GilsonPipetMaxConfig struct {
	Defaults *GilsonPipetMaxInstanceConfig                      `json:"Defaults,omitempty"`
	Devices  map[DeviceInstanceID]*GilsonPipetMaxInstanceConfig `json:"Devices"`
}

type GilsonPipetMaxInstanceConfig struct {
	Connection           string                    `json:"connection,omitempty"`
	LayoutPreferences    *liquidhandling.LayoutOpt `json:"layoutPreferences,omitempty"`
	OutputFileName       string                    `json:"outputFileName,omitempty"` // Specify file name in the instruction stream of any driver generated file
	MaxPlates            *float64                  `json:"maxPlates,omitempty"`
	MaxWells             *float64                  `json:"maxWells,omitempty"`
	ResidualVolumeWeight *float64                  `json:"residualVolumeWeight,omitempty"`
	InputPlateTypes      []wtype.PlateTypeName     `json:"inputPlateTypes,omitempty"`
	OutputPlateTypes     []wtype.PlateTypeName     `json:"outputPlateTypes,omitempty"`
	TipTypes             []string                  `json:"tipTypes,omitempty"`
}

type GlobalMixerConfig struct {
	PrintInstructions        bool `json:"printInstructions"`
	UseDriverTipTracking     bool `json:"useDriverTipTracking"`
	IgnorePhysicalSimulation bool `json:"ignorePhysicalSimulation"` //ignore errors in physical simulation

	// Direct specification of input and output plates
	InputPlates  []wtype.Plate `json:"inputPlates,omitempty"`
	OutputPlates []wtype.Plate `json:"outputPlates,omitempty"`

	CustomPolicyRuleSet *wtype.LHPolicyRuleSet `json:"customPolicyRuleSet,omitempty"`
}
