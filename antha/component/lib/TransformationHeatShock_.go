package lib

import (
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _TransformationHeatShockRequirements() {
}

// Conditions to run on startup
func _TransformationHeatShockSetup(_ctx context.Context, _input *TransformationHeatShockInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _TransformationHeatShockSteps(_ctx context.Context, _input *TransformationHeatShockInput, _output *TransformationHeatShockOutput) {
	var err error
	samples := make([]*wtype.LHComponent, 0)

	for k, plasmid := range _input.Plasmids {
		plasmid.Type, err = wtype.LiquidTypeFromString(_input.LHPolicyName)

		if err != nil {
			execute.Errorf(_ctx, "cannot find liquid type: %s", err)
		}
		dnaSample := mixer.Sample(plasmid, _input.PlasmidVolumes[k])
		dnaSample.CName = _input.PlasmidNames[k]
		dnaSample.Type = wtype.LTDNAMIX
		samples = append(samples, dnaSample)
	}

	for i, RecoveryPlateWell := range _input.RecoveryPlateWells {
		transformation := execute.MixInto(_ctx, _input.PlateWithCompetentCells, _input.CompetentCellPlateWells[i], samples...)

		transformationSample := mixer.Sample(transformation, _input.CompetentCellTransferVolume)

		// change liquid type to mix cells with SOC Media
		transformationSample.Type = wtype.LTPostMix

		HeatShock := execute.MixNamed(_ctx, _input.HeatShockPlate.Type, _input.HeatShockPlateWells[i], "HeatShockPlate", transformationSample)
		HeatShockSample := mixer.Sample(HeatShock, _input.CompetentCellTransferVolume)
		ColdRecovery := execute.MixInto(_ctx, _input.PlateWithCompetentCells, _input.ColdRecoveryCellPlateWells[i], HeatShockSample)
		ColdRecoverySample := mixer.Sample(ColdRecovery, _input.CompetentCellTransferVolume)
		ColdRecoverySample.Type = wtype.LTPostMix

		Recovery := execute.MixNamed(_ctx, _input.PlatewithRecoveryMedia.Type, RecoveryPlateWell, "RecoveryPlate", ColdRecoverySample)
		_output.RecoveredCells = append(_output.RecoveredCells, Recovery)

		// incubate the reaction mixture
		// commented out pending changes to incubate
		execute.Incubate(_ctx, Recovery, _input.RecoveryTemp, _input.RecoveryTime, true)
		// inactivate
		//Incubate(Reaction, InactivationTemp, InactivationTime, false)
	}
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _TransformationHeatShockAnalysis(_ctx context.Context, _input *TransformationHeatShockInput, _output *TransformationHeatShockOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _TransformationHeatShockValidation(_ctx context.Context, _input *TransformationHeatShockInput, _output *TransformationHeatShockOutput) {
}
func _TransformationHeatShockRun(_ctx context.Context, input *TransformationHeatShockInput) *TransformationHeatShockOutput {
	output := &TransformationHeatShockOutput{}
	_TransformationHeatShockSetup(_ctx, input)
	_TransformationHeatShockSteps(_ctx, input, output)
	_TransformationHeatShockAnalysis(_ctx, input, output)
	_TransformationHeatShockValidation(_ctx, input, output)
	return output
}

func TransformationHeatShockRunSteps(_ctx context.Context, input *TransformationHeatShockInput) *TransformationHeatShockSOutput {
	soutput := &TransformationHeatShockSOutput{}
	output := _TransformationHeatShockRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func TransformationHeatShockNew() interface{} {
	return &TransformationHeatShockElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &TransformationHeatShockInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _TransformationHeatShockRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &TransformationHeatShockInput{},
			Out: &TransformationHeatShockOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type TransformationHeatShockElement struct {
	inject.CheckedRunner
}

type TransformationHeatShockInput struct {
	ColdRecoveryCellPlateWells  []string
	CompetentCellPlateWells     []string
	CompetentCellTransferVolume wunit.Volume
	HeatShockPlate              *wtype.LHPlate
	HeatShockPlateWells         []string
	LHPolicyName                string
	PlasmidNames                []string
	PlasmidVolumes              []wunit.Volume
	Plasmids                    []*wtype.LHComponent
	PlateWithCompetentCells     *wtype.LHPlate
	PlatewithRecoveryMedia      *wtype.LHPlate
	PostPlasmidTemp             wunit.Temperature
	PostPlasmidTime             wunit.Time
	RecoveryPlateNumber         int
	RecoveryPlateWells          []string
	RecoveryTemp                wunit.Temperature
	RecoveryTime                wunit.Time
}

type TransformationHeatShockOutput struct {
	RecoveredCells []*wtype.LHComponent
}

type TransformationHeatShockSOutput struct {
	Data struct {
	}
	Outputs struct {
		RecoveredCells []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "TransformationHeatShock",
		Constructor: TransformationHeatShockNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/PooledLibrary/playground/TransformationHeatShock/TransformationHeatShock.an",
			Params: []component.ParamDesc{
				{Name: "ColdRecoveryCellPlateWells", Desc: "", Kind: "Parameters"},
				{Name: "CompetentCellPlateWells", Desc: "", Kind: "Parameters"},
				{Name: "CompetentCellTransferVolume", Desc: "", Kind: "Parameters"},
				{Name: "HeatShockPlate", Desc: "", Kind: "Inputs"},
				{Name: "HeatShockPlateWells", Desc: "", Kind: "Parameters"},
				{Name: "LHPolicyName", Desc: "", Kind: "Parameters"},
				{Name: "PlasmidNames", Desc: "", Kind: "Parameters"},
				{Name: "PlasmidVolumes", Desc: "", Kind: "Parameters"},
				{Name: "Plasmids", Desc: "", Kind: "Inputs"},
				{Name: "PlateWithCompetentCells", Desc: "", Kind: "Inputs"},
				{Name: "PlatewithRecoveryMedia", Desc: "", Kind: "Inputs"},
				{Name: "PostPlasmidTemp", Desc: "", Kind: "Parameters"},
				{Name: "PostPlasmidTime", Desc: "", Kind: "Parameters"},
				{Name: "RecoveryPlateNumber", Desc: "", Kind: "Parameters"},
				{Name: "RecoveryPlateWells", Desc: "", Kind: "Parameters"},
				{Name: "RecoveryTemp", Desc: "", Kind: "Parameters"},
				{Name: "RecoveryTime", Desc: "", Kind: "Parameters"},
				{Name: "RecoveredCells", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
