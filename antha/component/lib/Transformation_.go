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

//= 50.(uL)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _TransformationRequirements() {
}

// Conditions to run on startup
func _TransformationSetup(_ctx context.Context, _input *TransformationInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _TransformationSteps(_ctx context.Context, _input *TransformationInput, _output *TransformationOutput) {

	competetentcellmix := mixer.Sample(_input.ReadyCompCells, _input.CompetentCellvolumeperassembly)
	transformationmix := make([]*wtype.LHComponent, 0)
	transformationmix = append(transformationmix, competetentcellmix)
	DNAsample := mixer.Sample(_input.Reaction, _input.Reactionvolume)
	transformationmix = append(transformationmix, DNAsample)

	transformedcells := execute.MixInto(_ctx, _input.OutPlate, "", transformationmix...)

	_output.Transformedcells = execute.Incubate(_ctx, transformedcells, _input.Postplasmidtemp, _input.Postplasmidtime, false)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _TransformationAnalysis(_ctx context.Context, _input *TransformationInput, _output *TransformationOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _TransformationValidation(_ctx context.Context, _input *TransformationInput, _output *TransformationOutput) {
}
func _TransformationRun(_ctx context.Context, input *TransformationInput) *TransformationOutput {
	output := &TransformationOutput{}
	_TransformationSetup(_ctx, input)
	_TransformationSteps(_ctx, input, output)
	_TransformationAnalysis(_ctx, input, output)
	_TransformationValidation(_ctx, input, output)
	return output
}

func TransformationRunSteps(_ctx context.Context, input *TransformationInput) *TransformationSOutput {
	soutput := &TransformationSOutput{}
	output := _TransformationRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func TransformationNew() interface{} {
	return &TransformationElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &TransformationInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _TransformationRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &TransformationInput{},
			Out: &TransformationOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type TransformationElement struct {
	inject.CheckedRunner
}

type TransformationInput struct {
	CompetentCellvolumeperassembly wunit.Volume
	OutPlate                       *wtype.LHPlate
	Postplasmidtemp                wunit.Temperature
	Postplasmidtime                wunit.Time
	Reaction                       *wtype.LHComponent
	Reactionvolume                 wunit.Volume
	ReadyCompCells                 *wtype.LHComponent
}

type TransformationOutput struct {
	Transformedcells *wtype.LHComponent
}

type TransformationSOutput struct {
	Data struct {
	}
	Outputs struct {
		Transformedcells *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Transformation",
		Constructor: TransformationNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/Transformation/Transformation.an",
			Params: []component.ParamDesc{
				{Name: "CompetentCellvolumeperassembly", Desc: "= 50.(uL)\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Postplasmidtemp", Desc: "", Kind: "Parameters"},
				{Name: "Postplasmidtime", Desc: "", Kind: "Parameters"},
				{Name: "Reaction", Desc: "", Kind: "Inputs"},
				{Name: "Reactionvolume", Desc: "", Kind: "Parameters"},
				{Name: "ReadyCompCells", Desc: "", Kind: "Inputs"},
				{Name: "Transformedcells", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
