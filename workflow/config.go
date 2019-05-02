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
	Tecan          TecanConfig          `json:"Tecan"`
	CyBio          CyBioConfig          `json:"CyBio"`
	Labcyte        LabcyteConfig        `json:"Labcyte"`
	Hamilton       HamiltonConfig       `json:"Hamilton"`

	QPCR            QPCRConfig            `json:"QPCR"`
	ShakerIncubator ShakerIncubatorConfig `json:"ShakerIncubator"`
	PlateReader     PlateReaderConfig     `json:"PlateReader"`
}

func EmptyConfig() Config {
	return Config{
		GilsonPipetMax: GilsonPipetMaxConfig{
			Devices: make(map[DeviceInstanceID]*GilsonPipetMaxInstanceConfig),
		},
		Tecan: TecanConfig{
			Devices: make(map[DeviceInstanceID]*TecanInstanceConfig),
		},
		CyBio: CyBioConfig{
			Devices: make(map[DeviceInstanceID]*CyBioInstanceConfig),
		},
		Labcyte: LabcyteConfig{
			Devices: make(map[DeviceInstanceID]*LabcyteInstanceConfig),
		},
		Hamilton: HamiltonConfig{
			Devices: make(map[DeviceInstanceID]*HamiltonInstanceConfig),
		},

		QPCR: QPCRConfig{
			Devices: make(map[DeviceInstanceID]struct{}),
		},
		ShakerIncubator: ShakerIncubatorConfig{
			Devices: make(map[DeviceInstanceID]struct{}),
		},
		PlateReader: PlateReaderConfig{
			Devices: make(map[DeviceInstanceID]struct{}),
		},
	}
}

type GlobalMixerConfig struct {
	PrintInstructions        bool `json:"printInstructions"`
	UseDriverTipTracking     bool `json:"useDriverTipTracking"`
	IgnorePhysicalSimulation bool `json:"ignorePhysicalSimulation"` //ignore errors in physical simulation

	// Direct specification of input plates
	InputPlates []*wtype.Plate `json:"inputPlates,omitempty"`

	CustomPolicyRuleSet *wtype.LHPolicyRuleSet `json:"customPolicyRuleSet,omitempty"`
}

type DeviceInstanceID string

// Gilson
type GilsonPipetMaxConfig struct {
	Defaults *GilsonPipetMaxInstanceConfig                      `json:"Defaults,omitempty"`
	Devices  map[DeviceInstanceID]*GilsonPipetMaxInstanceConfig `json:"Devices,omitempty"`
}

type GilsonPipetMaxInstanceConfig struct {
	commonMixerInstanceConfig
	tipsOnly
}

func (cfg *GilsonPipetMaxInstanceConfig) MarshalJSON() ([]byte, error) {
	return MergeToMapAndMarshal(&cfg.commonMixerInstanceConfig, &cfg.tipsOnly)
}

func (cfg *GilsonPipetMaxInstanceConfig) UnmarshalJSON(bs []byte) error {
	return UnmarshalMapsMerged(bs, &cfg.commonMixerInstanceConfig, &cfg.tipsOnly)
}

// Tecan
type TecanConfig struct {
	Defaults *TecanInstanceConfig                      `json:"Defaults,omitempty"`
	Devices  map[DeviceInstanceID]*TecanInstanceConfig `json:"Devices,omitempty"`
}

type TecanInstanceConfig struct {
	ModelOnly
	commonMixerInstanceConfig
	tipsOnly
	GeneralConfig
}

func (cfg *TecanInstanceConfig) MarshalJSON() ([]byte, error) {
	return MergeToMapAndMarshal(&cfg.GeneralConfig, &cfg.commonMixerInstanceConfig, &cfg.ModelOnly, &cfg.tipsOnly)
}

func (cfg *TecanInstanceConfig) UnmarshalJSON(bs []byte) error {
	return UnmarshalMapsMerged(bs, &cfg.GeneralConfig, &cfg.commonMixerInstanceConfig, &cfg.ModelOnly, &cfg.tipsOnly)
}

// CyBio
type CyBioConfig struct {
	Defaults *CyBioInstanceConfig                      `json:"Defaults,omitempty"`
	Devices  map[DeviceInstanceID]*CyBioInstanceConfig `json:"Devices,omitempty"`
}

type CyBioInstanceConfig struct {
	ModelOnly
	commonMixerInstanceConfig
	tipsOnly
}

func (cfg *CyBioInstanceConfig) MarshalJSON() ([]byte, error) {
	return MergeToMapAndMarshal(&cfg.commonMixerInstanceConfig, &cfg.ModelOnly, &cfg.tipsOnly)
}

func (cfg *CyBioInstanceConfig) UnmarshalJSON(bs []byte) error {
	return UnmarshalMapsMerged(bs, &cfg.commonMixerInstanceConfig, &cfg.ModelOnly, &cfg.tipsOnly)
}

// Labcyte
type LabcyteConfig struct {
	Defaults *LabcyteInstanceConfig                      `json:"Defaults,omitempty"`
	Devices  map[DeviceInstanceID]*LabcyteInstanceConfig `json:"Devices,omitempty"`
}

type LabcyteInstanceConfig struct {
	ModelOnly
	commonMixerInstanceConfig
}

func (cfg *LabcyteInstanceConfig) MarshalJSON() ([]byte, error) {
	return MergeToMapAndMarshal(&cfg.commonMixerInstanceConfig, &cfg.ModelOnly)
}

func (cfg *LabcyteInstanceConfig) UnmarshalJSON(bs []byte) error {
	return UnmarshalMapsMerged(bs, &cfg.commonMixerInstanceConfig, &cfg.ModelOnly)
}

// HamiltonConfig
type HamiltonConfig struct {
	Defaults *HamiltonInstanceConfig                      `json:"Defaults,omitempty"`
	Devices  map[DeviceInstanceID]*HamiltonInstanceConfig `json:"Devices,omitempty"`
}

type HamiltonInstanceConfig struct {
	commonMixerInstanceConfig
	// tipsOnly - specifying tip types will be supported in later releases
}

func (cfg *HamiltonInstanceConfig) MarshalJSON() ([]byte, error) {
	return MergeToMapAndMarshal(&cfg.commonMixerInstanceConfig)
}

func (cfg *HamiltonInstanceConfig) UnmarshalJSON(bs []byte) error {
	return UnmarshalMapsMerged(bs, &cfg.commonMixerInstanceConfig)
}

type commonMixerInstanceConfig struct {
	Connection string `json:"Connection,omitempty"`

	LayoutPreferences    *LayoutOpt            `json:"layoutPreferences,omitempty"`
	MaxPlates            *float64              `json:"maxPlates,omitempty"`
	MaxWells             *float64              `json:"maxWells,omitempty"`
	ResidualVolumeWeight *float64              `json:"residualVolumeWeight,omitempty"`
	InputPlateTypes      []wtype.PlateTypeName `json:"inputPlateTypes,omitempty"`
	OutputPlateTypes     []wtype.PlateTypeName `json:"outputPlateTypes,omitempty"`

	ParsedConnection `json:"-"`
}

type ParsedConnection struct {
	HostPort      string `json:"-"`
	ExecFile      string `json:"-"`
	CompileAndRun string `json:"-"`
}

type tipsOnly struct {
	TipTypes []string `json:"tipTypes,omitempty"`
}

type ModelOnly struct {
	Model string `json:"Model"`
}

type GeneralConfig struct {
	Config string `json:"Config"`
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

func MergeToMapAndMarshal(components ...interface{}) ([]byte, error) {
	m := make(map[string]json.RawMessage)
	for _, com := range components {
		if bs, err := json.Marshal(com); err != nil {
			return nil, err
		} else if err := json.Unmarshal(bs, &m); err != nil {
			return nil, err
		}
	}
	return json.Marshal(m)
}

func UnmarshalMapsMerged(bs []byte, components ...interface{}) error {
	for _, com := range components {
		if err := json.Unmarshal(bs, com); err != nil {
			return err
		}
	}
	return nil
}

type QPCRConfig struct {
	Devices map[DeviceInstanceID]struct{} `json:"Devices,omitempty"`
}

type ShakerIncubatorConfig struct {
	Devices map[DeviceInstanceID]struct{} `json:"Devices,omitempty"`
}

type PlateReaderConfig struct {
	Devices map[DeviceInstanceID]struct{} `json:"Devices,omitempty"`
}
