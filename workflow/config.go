package workflow

import (
	"encoding/json"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type Config struct {
	GlobalMixer    GlobalMixerConfig    `json:"GlobalMixer"`
	GilsonPipetMax GilsonPipetMaxConfig `json:"GilsonPipetMax"`
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

type DeviceInstanceID string

type GilsonPipetMaxConfig struct {
	Defaults *GilsonPipetMaxInstanceConfig                      `json:"Defaults,omitempty"`
	Devices  map[DeviceInstanceID]*GilsonPipetMaxInstanceConfig `json:"Devices"`
}

type GilsonPipetMaxInstanceConfig struct {
	GenericDeviceConfig
}

type GenericDeviceConfig struct {
	Connection string `json:"Connection,omitempty"`
	Data       []byte `json:"-"`
}

func (gdc *GenericDeviceConfig) UnmarshalJSON(bs []byte) error {
	if string(bs) == "null" {
		return nil
	}
	connOnly := struct {
		Connection string `json:"Connection:omitempty"`
	}{}
	if err := json.Unmarshal(bs, &connOnly); err != nil {
		return err
	}
	gdc.Connection = connOnly.Connection
	// from the encoding/json docs: one must copy thy bs array if one
	// wishes to hang on to it:
	bsCopy := make([]byte, len(bs))
	copy(bsCopy, bs)
	gdc.Data = bsCopy
	return nil
}

func (gdc *GenericDeviceConfig) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	if err := json.Unmarshal(gdc.Data, &m); err != nil {
		return nil, err
	}
	// belt and braces - Go is essentially not case sensitive on JSON keys.
	for key := range m {
		if strings.ToLower(key) == "connection" {
			delete(m, key)
		}
	}
	if gdc.Connection != "" { // implement omitempty
		m["Connection"] = gdc.Connection
	}
	return json.Marshal(m)
}
