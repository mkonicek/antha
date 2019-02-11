package mixer

import (
	"errors"
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/driver/liquidhandling/client"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/laboratory/effects"
	driver "github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/workflow"
)

type GilsonPipetMaxInstance struct {
	ID                   workflow.DeviceInstanceID
	MaxPlates            float64
	MaxWells             float64
	ResidualVolumeWeight float64

	global *GlobalMixerConfig
	*workflow.GilsonPipetMaxInstanceConfig

	base   *BaseMixer
	driver driver.LiquidhandlingDriver
	// these are the properties as returned by the driver
	properties *driver.LHProperties
}

var (
	defaultMaxPlates            = 4.5
	defaultMaxWells             = 278.0
	defaultResidualVolumeWeight = 1.0
)

type GilsonPipetMaxInstances []*GilsonPipetMaxInstance

func NewGilsonPipetMaxInstances(inv *inventory.Inventory, global *GlobalMixerConfig, config workflow.GilsonPipetMaxConfig) (GilsonPipetMaxInstances, error) {
	defaultsWF := config.Defaults
	if defaultsWF == nil {
		defaultsWF = &workflow.GilsonPipetMaxInstanceConfig{}
	}

	defaults := &GilsonPipetMaxInstance{
		MaxPlates:                    floatValue(defaultsWF.MaxPlates, &defaultMaxPlates),
		MaxWells:                     floatValue(defaultsWF.MaxWells, &defaultMaxWells),
		ResidualVolumeWeight:         floatValue(defaultsWF.ResidualVolumeWeight, &defaultResidualVolumeWeight),
		GilsonPipetMaxInstanceConfig: defaultsWF,
	}
	if err := defaults.Validate(inv); err != nil {
		return nil, err
	}

	instances := make(GilsonPipetMaxInstances, 0, len(config.Devices))

	for id, instWF := range config.Devices {
		instance := &GilsonPipetMaxInstance{
			ID:                           id,
			MaxPlates:                    floatValue(instWF.MaxPlates, &defaults.MaxPlates),
			MaxWells:                     floatValue(instWF.MaxWells, &defaults.MaxWells),
			ResidualVolumeWeight:         floatValue(instWF.MaxPlates, &defaults.ResidualVolumeWeight),
			global:                       global,
			GilsonPipetMaxInstanceConfig: instWF,
			base:                         NewBaseMixer(instWF.Connection, "GilsonPipetmax"),
		}
		if err := instance.Validate(inv); err != nil {
			return nil, err
		} else {
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

func (insts GilsonPipetMaxInstances) Connect(wf *workflow.Workflow) error {
	for _, inst := range insts {
		if err := inst.Connect(wf); err != nil {
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
	if inst.driver == nil {
		if conn, err := inst.base.ConnectInit(); err != nil {
			return err
		} else if conn != nil {
			driver := client.NewLowLevelClientFromConn(conn)
			if props, status := driver.Configure(wf.JobId, wf.Meta.Name, inst.ID); !status.Ok() {
				return status.GetError()
			} else if err := props.ApplyUserPreferences(inst.LayoutPreferences); err != nil {
				return err
			} else {
				inst.driver = driver
				props.Driver = driver
				inst.properties = props
			}
		}
	}
	return nil
}

func (inst *GilsonPipetMaxInstance) CanCompile(req effects.Request) bool {
	can := effects.Request{
		Selector: []effects.NameValue{
			target.DriverSelectorV1Mixer,
			target.DriverSelectorV1Prompter,
		},
	}
	if inst.properties.CanPrompt() {
		can.Selector = append(can.Selector, target.DriverSelectorV1Prompter)
	}
	return can.Contains(req)
}

func (inst *GilsonPipetMaxInstance) Compile(labEffects *effects.LaboratoryEffects, nodes []effects.Node) ([]effects.Inst, error) {
	var instrs []*wtype.LHInstruction
	for _, node := range nodes {
		if cmd, ok := node.(*effects.Command); !ok {
			return nil, fmt.Errorf("cannot compile %T", node)
		} else if instr, ok := cmd.Inst.(*wtype.LHInstruction); !ok {
			return nil, fmt.Errorf("cannot compile %T", cmd.Inst)
		} else {
			instrs = append(instrs, instr)
		}
	}

	return inst.mix(labEffects, instrs)
}

func (inst *GilsonPipetMaxInstance) mix(labEffects *effects.LaboratoryEffects, instrs []*wtype.LHInstruction) ([]effects.Inst, error) {
	if len(instrs) == 0 {
		return nil, errors.New("No instructions to mix!")
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

	for _, ttn := range inst.TipTypes {
		if tb, err := labEffects.Inventory.TipBoxes.NewTipbox(ttn); err != nil {
			return nil, err
		} else {
			req.TipBoxes = append(req.TipBoxes, tb)
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

func floatValue(a, b *float64) float64 {
	if a != nil {
		return *a
	} else {
		return *b
	}
}
