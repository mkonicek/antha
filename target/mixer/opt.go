package mixer

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/workflow"
)

type GlobalMixerConfig struct {
	*workflow.GlobalMixerConfig
}

func (cfg *GlobalMixerConfig) validate(inv *inventory.Inventory) error {
	for _, plates := range [][]wtype.Plate{cfg.InputPlates, cfg.OutputPlates} {
		for _, plate := range plates {
			if _, err := inv.PlateTypes.NewPlateType(plate.Type); err != nil {
				return err
			}
		}
	}
	return nil
}

var (
	defaultMaxPlates            = 4.5
	defaultMaxWells             = 278.0
	defaultResidualVolumeWeight = 1.0
)

type GilsonPipetMaxInstances map[workflow.DeviceInstanceID]*GilsonPipetMaxInstanceConfig

type GilsonPipetMaxInstanceConfig struct {
	*GlobalMixerConfig

	MaxPlates            float64
	MaxWells             float64
	ResidualVolumeWeight float64

	*workflow.GilsonPipetMaxInstanceConfig
}

func GilsonPipetMaxInstancesFromWorkflow(wf *workflow.Workflow, inv *inventory.Inventory) (GilsonPipetMaxInstances, error) {
	global := &GlobalMixerConfig{
		GlobalMixerConfig: &wf.Config.GlobalMixer,
	}
	if err := global.validate(inv); err != nil {
		return nil, err
	}

	devices := wf.Config.GilsonPipetMax

	defaults := &GilsonPipetMaxInstanceConfig{
		MaxPlates:                    floatValue(devices.Defaults.MaxPlates, &defaultMaxPlates),
		MaxWells:                     floatValue(devices.Defaults.MaxWells, &defaultMaxWells),
		ResidualVolumeWeight:         floatValue(devices.Defaults.ResidualVolumeWeight, &defaultResidualVolumeWeight),
		GilsonPipetMaxInstanceConfig: devices.Defaults,
	}

	res := make(GilsonPipetMaxInstances, len(devices.Devices))
	for id, cfgWf := range devices.Devices {
		cfg := &GilsonPipetMaxInstanceConfig{
			GlobalMixerConfig: global,

			MaxPlates:            floatValue(cfgWf.MaxPlates, &defaults.MaxPlates),
			MaxWells:             floatValue(cfgWf.MaxWells, &defaults.MaxWells),
			ResidualVolumeWeight: floatValue(cfgWf.MaxPlates, &defaults.MaxPlates),

			GilsonPipetMaxInstanceConfig: cfgWf,
		}
		if err := cfg.validate(inv); err != nil {
			return nil, err
		}
		res[id] = cfg
	}
	return res, nil
}

func (cfg *GilsonPipetMaxInstanceConfig) validate(inv *inventory.Inventory) error {
	for _, ptns := range [][]wtype.PlateTypeName{cfg.InputPlateTypes, cfg.OutputPlateTypes} {
		for _, ptn := range ptns {
			if _, err := inv.PlateTypes.NewPlateType(ptn); err != nil {
				return err
			}
		}
	}
	return nil
}

func floatValue(a, b *float64) float64 {
	if a != nil {
		return *a
	} else {
		return *b
	}
}
