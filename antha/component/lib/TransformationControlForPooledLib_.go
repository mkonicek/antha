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

//PlatewithRecoveryMedia *wtype.LHPlate

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _TransformationControlForPooledLibRequirements() {
}

// Conditions to run on startup
func _TransformationControlForPooledLibSetup(_ctx context.Context, _input *TransformationControlForPooledLibInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _TransformationControlForPooledLibSteps(_ctx context.Context, _input *TransformationControlForPooledLibInput, _output *TransformationControlForPooledLibOutput) {

	Reaction := _input.ReactionsMap[_input.ReactionName]
	dnaSample := mixer.Sample(Reaction, _input.ReactionVolume)
	dnaSample.Type = wtype.LTDNAMIX

	for i, RecoveryPlateWell := range _input.RecoveryPlateWells {
		transformation := execute.MixTo(_ctx, _input.PlateWithCompetentCells.Type, _input.CompetentCellPlateWells[i], 1, dnaSample)

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
func _TransformationControlForPooledLibAnalysis(_ctx context.Context, _input *TransformationControlForPooledLibInput, _output *TransformationControlForPooledLibOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _TransformationControlForPooledLibValidation(_ctx context.Context, _input *TransformationControlForPooledLibInput, _output *TransformationControlForPooledLibOutput) {
}
func _TransformationControlForPooledLibRun(_ctx context.Context, input *TransformationControlForPooledLibInput) *TransformationControlForPooledLibOutput {
	output := &TransformationControlForPooledLibOutput{}
	_TransformationControlForPooledLibSetup(_ctx, input)
	_TransformationControlForPooledLibSteps(_ctx, input, output)
	_TransformationControlForPooledLibAnalysis(_ctx, input, output)
	_TransformationControlForPooledLibValidation(_ctx, input, output)
	return output
}

func TransformationControlForPooledLibRunSteps(_ctx context.Context, input *TransformationControlForPooledLibInput) *TransformationControlForPooledLibSOutput {
	soutput := &TransformationControlForPooledLibSOutput{}
	output := _TransformationControlForPooledLibRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func TransformationControlForPooledLibNew() interface{} {
	return &TransformationControlForPooledLibElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &TransformationControlForPooledLibInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _TransformationControlForPooledLibRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &TransformationControlForPooledLibInput{},
			Out: &TransformationControlForPooledLibOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type TransformationControlForPooledLibElement struct {
	inject.CheckedRunner
}

type TransformationControlForPooledLibInput struct {
	CompetentCellPlateWells     []string
	CompetentCellTransferVolume wunit.Volume
	PlateWithCompetentCells     *wtype.LHPlate
	PlatewithRecoveryMedia      *wtype.LHPlate
	PostPlasmidTemp             wunit.Temperature
	PostPlasmidTime             wunit.Time
	ReactionName                string
	ReactionVolume              wunit.Volume
	ReactionsMap                map[string]*wtype.LHComponent
	RecoveryPlateNumber         int
	RecoveryPlateWells          []string
	RecoveryTemp                wunit.Temperature
	RecoveryTime                wunit.Time
}

type TransformationControlForPooledLibOutput struct {
	RecoveredCells []*wtype.LHComponent
}

type TransformationControlForPooledLibSOutput struct {
	Data struct {
	}
	Outputs struct {
		RecoveredCells []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "TransformationControlForPooledLib",
		Constructor: TransformationControlForPooledLibNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/PooledLibrary/playground/Transformation control/Transformation.an",
			Params: []component.ParamDesc{
				{Name: "CompetentCellPlateWells", Desc: "PlatewithRecoveryMedia *wtype.LHPlate\n", Kind: "Parameters"},
				{Name: "CompetentCellTransferVolume", Desc: "", Kind: "Parameters"},
				{Name: "PlateWithCompetentCells", Desc: "", Kind: "Inputs"},
				{Name: "PlatewithRecoveryMedia", Desc: "", Kind: "Inputs"},
				{Name: "PostPlasmidTemp", Desc: "", Kind: "Parameters"},
				{Name: "PostPlasmidTime", Desc: "", Kind: "Parameters"},
				{Name: "ReactionName", Desc: "", Kind: "Parameters"},
				{Name: "ReactionVolume", Desc: "", Kind: "Parameters"},
				{Name: "ReactionsMap", Desc: "", Kind: "Inputs"},
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
