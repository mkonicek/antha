package workflow

import (
	"encoding/json"
	"net"
	"net/url"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type Config struct {
	GlobalMixer    GlobalMixerConfig    `json:"GlobalMixer"`
	GilsonPipetMax GilsonPipetMaxConfig `json:"GilsonPipetMax"`
	Labcyte        LabcyteConfig        `json:"Labcyte"`
}

type GlobalMixerConfig struct {
	PrintInstructions        bool `json:"printInstructions"`
	UseDriverTipTracking     bool `json:"useDriverTipTracking"`
	IgnorePhysicalSimulation bool `json:"ignorePhysicalSimulation"` //ignore errors in physical simulation

	// Direct specification of input and output plates
	InputPlates  []*wtype.Plate `json:"inputPlates,omitempty"`
	OutputPlates []*wtype.Plate `json:"outputPlates,omitempty"`

	CustomPolicyRuleSet *wtype.LHPolicyRuleSet `json:"customPolicyRuleSet,omitempty"`
}

type DeviceInstanceID string

// Gilson
type GilsonPipetMaxConfig struct {
	Defaults *GilsonPipetMaxInstanceConfig                      `json:"Defaults,omitempty"`
	Devices  map[DeviceInstanceID]*GilsonPipetMaxInstanceConfig `json:"Devices"`
}

type GilsonPipetMaxInstanceConfig struct {
	commonMixerInstanceConfig
}

// Labcyte
type LabcyteConfig struct {
	Defaults *LabcyteInstanceConfig                      `json:"Defaults,omitempty"`
	Devices  map[DeviceInstanceID]*LabcyteInstanceConfig `json:"Devices"`
}

type LabcyteInstanceConfig struct {
	modelOnly
	commonMixerInstanceConfig
}

type modelOnly struct {
	Model string `json:model`
}

func (cfg *LabcyteInstanceConfig) MarshalJSON() ([]byte, error) {
	if commonBs, err := json.Marshal(&cfg.commonMixerInstanceConfig); err != nil {
		return nil, err
	} else if modelBs, err := json.Marshal(&cfg.modelOnly); err != nil {
		return nil, err

	} else {
		m := make(map[string]json.RawMessage)
		if err := json.Unmarshal(commonBs, &m); err != nil {
			return nil, err
		} else if err := json.Unmarshal(modelBs, &m); err != nil {
			return nil, err
		} else {
			return json.Marshal(m)
		}
	}
}

func (cfg *LabcyteInstanceConfig) UnmarshalJSON(bs []byte) error {
	if err := json.Unmarshal(bs, &cfg.commonMixerInstanceConfig); err != nil {
		return err
	} else if err := json.Unmarshal(bs, &cfg.modelOnly); err != nil {
		return err
	} else {
		return nil
	}
}

type commonMixerInstanceConfig struct {
	Connection string `json:"Connection,omitempty"`

	LayoutPreferences    *LayoutOpt            `json:"layoutPreferences,omitempty"`
	MaxPlates            *float64              `json:"maxPlates,omitempty"`
	MaxWells             *float64              `json:"maxWells,omitempty"`
	ResidualVolumeWeight *float64              `json:"residualVolumeWeight,omitempty"`
	InputPlateTypes      []wtype.PlateTypeName `json:"inputPlateTypes,omitempty"`
	OutputPlateTypes     []wtype.PlateTypeName `json:"outputPlateTypes,omitempty"`
	TipTypes             []string              `json:"tipTypes,omitempty"`

	ParsedConnection `json:"-"`
}

type ParsedConnection struct {
	HostPort      string `json:"-"`
	ExecFile      string `json:"-"`
	CompileAndRun string `json:"-"`
}

// type aliases do not inherit methods, so this is a cheap way to
// avoid infinite recursion:
type commonMixerInstanceConfigNoCustomMarshal commonMixerInstanceConfig

func (cfg *commonMixerInstanceConfig) MarshalJSON() ([]byte, error) {
	switch {
	case cfg.HostPort != "":
		cfg.Connection = cfg.HostPort
	case cfg.ExecFile != "":
		cfg.Connection = "file://" + cfg.ExecFile
	case cfg.CompileAndRun != "":
		cfg.Connection = "go://" + cfg.CompileAndRun
	}
	cfg2 := (*commonMixerInstanceConfigNoCustomMarshal)(cfg)
	return json.Marshal(cfg2)
}

func (cfg *commonMixerInstanceConfig) UnmarshalJSON(bs []byte) error {
	cfg2 := commonMixerInstanceConfigNoCustomMarshal{}
	if err := json.Unmarshal(bs, &cfg2); err != nil {
		return err
	}
	*cfg = commonMixerInstanceConfig(cfg2)

	if u, err := url.Parse(cfg.Connection); err == nil && u.Scheme == "go" {
		cfg.CompileAndRun = u.Host + u.Path
	} else if err == nil && u.Scheme == "file" {
		cfg.ExecFile = u.Host + u.Path // have to include Host to cope with PATH-based lookups, or relative paths
	} else if _, _, err := net.SplitHostPort(cfg.Connection); err == nil {
		cfg.HostPort = cfg.Connection
	} else {
		cfg.ExecFile = cfg.Connection
	}

	cfg.Connection = "" // wipe it out to make sure we don't accidentally use it.

	return nil
}
