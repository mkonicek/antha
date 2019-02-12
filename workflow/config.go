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

type GilsonPipetMaxConfig struct {
	Defaults *GilsonPipetMaxInstanceConfig                      `json:"Defaults,omitempty"`
	Devices  map[DeviceInstanceID]*GilsonPipetMaxInstanceConfig `json:"Devices"`
}

type GilsonPipetMaxInstanceConfig struct {
	Connection           string                `json:"Connection,omitempty"`
	LayoutPreferences    *LayoutOpt            `json:"layoutPreferences,omitempty"`
	MaxPlates            *float64              `json:"maxPlates,omitempty"`
	MaxWells             *float64              `json:"maxWells,omitempty"`
	ResidualVolumeWeight *float64              `json:"residualVolumeWeight,omitempty"`
	InputPlateTypes      []wtype.PlateTypeName `json:"inputPlateTypes,omitempty"`
	OutputPlateTypes     []wtype.PlateTypeName `json:"outputPlateTypes,omitempty"`
	TipTypes             []string              `json:"tipTypes,omitempty"`

	ParsedConnection ParsedConnection `json:"-"`
}

type ParsedConnection struct {
	HostPort      string `json:"-"`
	ExecFile      string `json:"-"`
	CompileAndRun string `json:"-"`
}

// type aliases do not inherit methods, so this is a cheap way to
// avoid infinite recursion:
type gilsonPipetMaxInstanceConfigNoCustomMarshal GilsonPipetMaxInstanceConfig

func (cfg *GilsonPipetMaxInstanceConfig) MarshalJSON() ([]byte, error) {
	switch {
	case cfg.ParsedConnection.HostPort != "":
		cfg.Connection = cfg.ParsedConnection.HostPort
	case cfg.ParsedConnection.ExecFile != "":
		cfg.Connection = "file://" + cfg.ParsedConnection.ExecFile
	case cfg.ParsedConnection.CompileAndRun != "":
		cfg.Connection = "go://" + cfg.ParsedConnection.CompileAndRun
	}
	cfg2 := (*gilsonPipetMaxInstanceConfigNoCustomMarshal)(cfg)
	return json.Marshal(cfg2)
}

func (cfg *GilsonPipetMaxInstanceConfig) UnmarshalJSON(bs []byte) error {
	cfg2 := gilsonPipetMaxInstanceConfigNoCustomMarshal{}
	if err := json.Unmarshal(bs, &cfg2); err != nil {
		return err
	}
	*cfg = GilsonPipetMaxInstanceConfig(cfg2)

	if u, err := url.Parse(cfg.Connection); err == nil && u.Scheme == "go" {
		cfg.ParsedConnection.CompileAndRun = u.Host + u.Path
	} else if err == nil && u.Scheme == "file" {
		cfg.ParsedConnection.ExecFile = u.Host + u.Path
	} else if _, _, err := net.SplitHostPort(cfg.Connection); err == nil {
		cfg.ParsedConnection.HostPort = cfg.Connection
	} else {
		cfg.ParsedConnection.ExecFile = cfg.Connection
	}

	cfg.Connection = "" // wipe it out to make sure we don't accidentally use it.

	return nil
}
