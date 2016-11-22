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

func _TransformationCotransRequirements() {
}

// Conditions to run on startup
func _TransformationCotransSetup(_ctx context.Context, _input *TransformationCotransInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _TransformationCotransSteps(_ctx context.Context, _input *TransformationCotransInput, _output *TransformationCotransOutput) {
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
		transformation := execute.MixTo(_ctx, _input.PlateWithCompetentCells.Type, _input.CompetentCellPlateWells[i], 1, samples...)

		transformationSample := mixer.Sample(transformation, _input.CompetentCellTransferVolume)

		// change liquid type to mix cells with SOC Media
		transformationSample.Type = wtype.LTPostMix

		Recovery := execute.MixNamed(_ctx, _input.PlatewithRecoveryMedia.Type, RecoveryPlateWell, "RecoveryPlate", transformationSample)
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
func _TransformationCotransAnalysis(_ctx context.Context, _input *TransformationCotransInput, _output *TransformationCotransOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _TransformationCotransValidation(_ctx context.Context, _input *TransformationCotransInput, _output *TransformationCotransOutput) {
}
func _TransformationCotransRun(_ctx context.Context, input *TransformationCotransInput) *TransformationCotransOutput {
	output := &TransformationCotransOutput{}
	_TransformationCotransSetup(_ctx, input)
	_TransformationCotransSteps(_ctx, input, output)
	_TransformationCotransAnalysis(_ctx, input, output)
	_TransformationCotransValidation(_ctx, input, output)
	return output
}

func TransformationCotransRunSteps(_ctx context.Context, input *TransformationCotransInput) *TransformationCotransSOutput {
	soutput := &TransformationCotransSOutput{}
	output := _TransformationCotransRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func TransformationCotransNew() interface{} {
	return &TransformationCotransElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &TransformationCotransInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _TransformationCotransRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &TransformationCotransInput{},
			Out: &TransformationCotransOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type TransformationCotransElement struct {
	inject.CheckedRunner
}

type TransformationCotransInput struct {
	CompetentCellPlateWells     []string
	CompetentCellTransferVolume wunit.Volume
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

type TransformationCotransOutput struct {
	RecoveredCells []*wtype.LHComponent
}

type TransformationCotransSOutput struct {
	Data struct {
	}
	Outputs struct {
		RecoveredCells []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "TransformationCotrans",
		Constructor: TransformationCotransNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/PooledLibrary/playground/Cotransformation/TransformationCotrans.an",
			Params: []component.ParamDesc{
				{Name: "CompetentCellPlateWells", Desc: "", Kind: "Parameters"},
				{Name: "CompetentCellTransferVolume", Desc: "", Kind: "Parameters"},
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
