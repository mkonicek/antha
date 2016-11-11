package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image/pick"
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

func _ChooseColoniesRequirements() {

}

// Conditions to run on startup
func _ChooseColoniesSetup(_ctx context.Context, _input *ChooseColoniesInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _ChooseColoniesSteps(_ctx context.Context, _input *ChooseColoniesInput, _output *ChooseColoniesOutput) {

	_output.Wellstopick = pick.Pick(_input.Imagefile, _input.Plate, _input.NumbertoPick, _input.Setplateperimeterfirst, _input.Rotate)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _ChooseColoniesAnalysis(_ctx context.Context, _input *ChooseColoniesInput, _output *ChooseColoniesOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _ChooseColoniesValidation(_ctx context.Context, _input *ChooseColoniesInput, _output *ChooseColoniesOutput) {

}
func _ChooseColoniesRun(_ctx context.Context, input *ChooseColoniesInput) *ChooseColoniesOutput {
	output := &ChooseColoniesOutput{}
	_ChooseColoniesSetup(_ctx, input)
	_ChooseColoniesSteps(_ctx, input, output)
	_ChooseColoniesAnalysis(_ctx, input, output)
	_ChooseColoniesValidation(_ctx, input, output)
	return output
}

func ChooseColoniesRunSteps(_ctx context.Context, input *ChooseColoniesInput) *ChooseColoniesSOutput {
	soutput := &ChooseColoniesSOutput{}
	output := _ChooseColoniesRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ChooseColoniesNew() interface{} {
	return &ChooseColoniesElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ChooseColoniesInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ChooseColoniesRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ChooseColoniesInput{},
			Out: &ChooseColoniesOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type ChooseColoniesElement struct {
	inject.CheckedRunner
}

type ChooseColoniesInput struct {
	Imagefile              string
	NumbertoPick           int
	Plate                  *wtype.LHPlate
	Rotate                 bool
	Setplateperimeterfirst bool
}

type ChooseColoniesOutput struct {
	Wellstopick []string
}

type ChooseColoniesSOutput struct {
	Data struct {
		Wellstopick []string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "ChooseColonies",
		Constructor: ChooseColoniesNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Data/choosecolonies/ChooseColonies.an",
			Params: []component.ParamDesc{
				{Name: "Imagefile", Desc: "", Kind: "Parameters"},
				{Name: "NumbertoPick", Desc: "", Kind: "Parameters"},
				{Name: "Plate", Desc: "", Kind: "Inputs"},
				{Name: "Rotate", Desc: "", Kind: "Parameters"},
				{Name: "Setplateperimeterfirst", Desc: "", Kind: "Parameters"},
				{Name: "Wellstopick", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
