package mixer

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
	"github.com/antha-lang/antha/workflow"
)

type LabcyteInstance struct {
	ID                   workflow.DeviceInstanceID
	MaxPlates            float64
	MaxWells             float64
	ResidualVolumeWeight float64

	global *GlobalMixerConfig
	*workflow.LabcyteInstanceConfig

	*BaseMixer
}

type LabcyteInstances []*LabcyteInstance

func NewLabcyteInstances(logger *logger.Logger, inv *inventory.Inventory, global *GlobalMixerConfig, config workflow.LabcyteConfig) (LabcyteInstances, error) {
	defaultsWF := config.Defaults
	if defaultsWF == nil {
		defaultsWF = &workflow.LabcyteInstanceConfig{}
	}

	var (
		defaultMaxPlates            = 4.5
		defaultMaxWells             = 278.0
		defaultResidualVolumeWeight = 1.0
	)

	defaults := &LabcyteInstance{
		MaxPlates:             floatValue(defaultsWF.MaxPlates, &defaultMaxPlates),
		MaxWells:              floatValue(defaultsWF.MaxWells, &defaultMaxWells),
		ResidualVolumeWeight:  floatValue(defaultsWF.ResidualVolumeWeight, &defaultResidualVolumeWeight),
		LabcyteInstanceConfig: defaultsWF,
	}
	if err := defaults.Validate(inv); err != nil {
		return nil, err
	}

	instances := make(LabcyteInstances, 0, len(config.Devices))

	for id, instWF := range config.Devices {
		instance := &LabcyteInstance{
			ID:                    id,
			MaxPlates:             floatValue(instWF.MaxPlates, &defaults.MaxPlates),
			MaxWells:              floatValue(instWF.MaxWells, &defaults.MaxWells),
			ResidualVolumeWeight:  floatValue(instWF.MaxPlates, &defaults.ResidualVolumeWeight),
			global:                global,
			LabcyteInstanceConfig: instWF,
			BaseMixer:             NewBaseMixer(logger, id, instWF.ParsedConnection, LabcyteSubType),
		}
		if err := instance.Validate(inv); err != nil {
			return nil, err
		} else {
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

func (inst *LabcyteInstance) Validate(inv *inventory.Inventory) error {
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

	return nil
}

func (inst *LabcyteInstance) Connect(wf *workflow.Workflow) error {
	if inst.properties == nil {
		if data, err := json.Marshal(inst.Model); err != nil {
			return err
		} else if err := inst.connect(wf, data); err != nil {
			return err
		} else if err := inst.properties.ApplyUserPreferences(inst.LayoutPreferences); err != nil {
			inst.Close()
			return err
		}
	}
	return nil
}

func (inst *LabcyteInstance) Compile(labEffects *effects.LaboratoryEffects, nodes []effects.Node) ([]effects.Inst, error) {
	instrs, err := checkInstructions(nodes)
	if err != nil {
		return nil, err
	}

	props := inst.properties.Dup(labEffects.IDGenerator)
	req := liquidhandling.NewLHRequest(labEffects.IDGenerator)
	req.BlockID = instrs[0].BlockID

	if err := inst.global.ApplyToLHRequest(req); err != nil {
		return nil, err
	}

	req.InputSetupWeights["MAX_N_PLATES"] = inst.MaxPlates
	req.InputSetupWeights["MAX_N_WELLS"] = inst.MaxWells
	req.InputSetupWeights["RESIDUAL_VOLUME_WEIGHT"] = inst.ResidualVolumeWeight

	for _, ptn := range inst.InputPlateTypes {
		if pt, err := labEffects.Inventory.PlateTypes.NewPlate(ptn); err != nil {
			return nil, err
		} else {
			req.InputPlatetypes = append(req.InputPlatetypes, pt)
		}
	}

	for _, ptn := range inst.OutputPlateTypes {
		if pt, err := labEffects.Inventory.PlateTypes.NewPlate(ptn); err != nil {
			return nil, err
		} else {
			req.OutputPlatetypes = append(req.OutputPlatetypes, pt)
		}
	}

	for _, ps := range [][]*wtype.Plate{inst.global.InputPlates, labEffects.SampleTracker.GetInputPlates()} {
		for _, p := range ps {
			if err := req.AddUserPlate(labEffects.IDGenerator, p); err != nil {
				return nil, err
			}
		}
	}

	if err := addCustomPolicies(instrs, req); err != nil {
		return nil, err
	}

	hasOutputPlate := func(typ wtype.PlateTypeName, id string) bool {
		for _, p := range req.OutputPlatetypes {
			if p.Type == typ && (id == "" || p.ID == id) {
				return true
			}
		}
		return false
	}

	for _, instr := range instrs {
		if instr.OutPlate != nil {
			if p, found := req.OutputPlates[instr.OutPlate.ID]; found && p != instr.OutPlate {
				return nil, fmt.Errorf("Mix setup error: Plate %s already requested in different state for mix.", p.ID)
			} else {
				req.OutputPlates[instr.OutPlate.ID] = instr.OutPlate
			}
		}

		if len(instr.Platetype) != 0 && !hasOutputPlate(instr.Platetype, instr.PlateID) {
			if pt, err := labEffects.Inventory.PlateTypes.NewPlate(instr.Platetype); err != nil {
				return nil, err
			} else {
				pt.ID = instr.PlateID
				req.OutputPlatetypes = append(req.OutputPlatetypes, pt)
			}
		}
		req.Add_instruction(instr)
	}

	planner := liquidhandling.Init(props)

	if err := planner.MakeSolutions(labEffects, req); err != nil {
		return nil, err
	}

	return nil, nil
}
