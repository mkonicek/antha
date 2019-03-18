package mixer

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/workflow"
)

var (
	_ effects.Device = (*GilsonPipetMaxInstance)(nil)
)

type GilsonPipetMaxInstance struct {
	ID                   workflow.DeviceInstanceID
	MaxPlates            float64
	MaxWells             float64
	ResidualVolumeWeight float64

	global *GlobalMixerConfig
	*workflow.GilsonPipetMaxInstanceConfig

	*BaseMixer
}

func NewGilsonPipetMaxInstances(logger *logger.Logger, tgt *target.Target, inv *inventory.Inventory, global *GlobalMixerConfig, config workflow.GilsonPipetMaxConfig) error {
	defaultsWF := config.Defaults
	if defaultsWF == nil {
		defaultsWF = &workflow.GilsonPipetMaxInstanceConfig{}
	}

	var (
		defaultMaxPlates            = 4.5
		defaultMaxWells             = 278.0
		defaultResidualVolumeWeight = 1.0
	)

	defaults := &GilsonPipetMaxInstance{
		MaxPlates:                    floatValue(defaultsWF.MaxPlates, &defaultMaxPlates),
		MaxWells:                     floatValue(defaultsWF.MaxWells, &defaultMaxWells),
		ResidualVolumeWeight:         floatValue(defaultsWF.ResidualVolumeWeight, &defaultResidualVolumeWeight),
		GilsonPipetMaxInstanceConfig: defaultsWF,
	}
	if err := defaults.Validate(inv); err != nil {
		return err
	}

	for id, instWF := range config.Devices {
		instance := &GilsonPipetMaxInstance{
			ID:                           id,
			MaxPlates:                    floatValue(instWF.MaxPlates, &defaults.MaxPlates),
			MaxWells:                     floatValue(instWF.MaxWells, &defaults.MaxWells),
			ResidualVolumeWeight:         floatValue(instWF.MaxPlates, &defaults.ResidualVolumeWeight),
			global:                       global,
			GilsonPipetMaxInstanceConfig: instWF,
			BaseMixer:                    NewBaseMixer(logger, id, instWF.ParsedConnection, target.GilsonPipetmaxSubType),
		}
		if err := instance.Validate(inv); err != nil {
			return err
		} else if err := tgt.AddDevice(instance); err != nil {
			return err
		}
	}

	return nil
}

func (inst *GilsonPipetMaxInstance) Validate(inv *inventory.Inventory) error {
	switch {
	case inst.MaxPlates <= 0:
		return errors.New("Validation error: MaxPlates must be > 0")
	case inst.MaxWells <= 0:
		return errors.New("Validation error: MaxWells must be > 0")
	case inst.ResidualVolumeWeight < 0:
		return errors.New("Validation error: ResidualVolumeWeight must be >= 0")
	}

	// TODO: add extra validation here!
	for _, ptns := range [][]wtype.PlateTypeName{inst.InputPlateTypes, inst.OutputPlateTypes} {
		for _, ptn := range ptns {
			if _, err := inv.PlateTypes.NewPlateType(ptn); err != nil {
				return err
			}
		}
	}

	for _, tt := range inst.TipTypes {
		if _, err := inv.TipBoxes.FetchTipbox(tt); err != nil {
			return err
		}
	}

	return nil
}

func (inst *GilsonPipetMaxInstance) Connect(wf *workflow.Workflow) error {
	if inst.properties == nil {
		if err := inst.connect(wf, nil); err != nil {
			return err
		} else if err := inst.properties.ApplyUserPreferences(inst.LayoutPreferences); err != nil {
			inst.Close()
			return err
		}
	}
	return nil
}

func (inst *GilsonPipetMaxInstance) Compile(labEffects *effects.LaboratoryEffects, dir string, nodes []effects.Node) (effects.Insts, error) {
	instrs, err := checkInstructions(nodes)
	if err != nil {
		return nil, err
	}

	mix, err := mixOpts{
		Device:     inst,
		Base:       inst.BaseMixer,
		LabEffects: labEffects,
		Global:     inst.global,
		Instrs:     instrs,
		InputWeights: map[string]float64{
			"MAX_N_PLATES":           inst.MaxPlates,
			"MAX_N_WELLS":            inst.MaxWells,
			"RESIDUAL_VOLUME_WEIGHT": inst.ResidualVolumeWeight,
		},
		InputPlateTypes:  inst.InputPlateTypes,
		OutputPlateTypes: inst.OutputPlateTypes,
		TipTypes:         inst.TipTypes,
	}.mix()
	if err != nil {
		return nil, err
	}

	tarballPath := filepath.Join(dir, fmt.Sprintf("%v.tar.gz", inst.ID))
	contentPath := fmt.Sprintf("%v-%v.sqlite", inst.ID, instrs[0].BlockID)

	if err != nil {
		return nil, err
	} else if err := writeToTarball(tarballPath, contentPath, mix.Files.Tarball); err != nil {
		return nil, err
	} else {
		return []effects.Inst{mix}, nil
	}
}
