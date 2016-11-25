package lib

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/target/mixer"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _InputPlateInitializerRequirements() {

}

// Conditions to run on startup
func _InputPlateInitializerSetup(_ctx context.Context, _input *InputPlateInitializerInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _InputPlateInitializerSteps(_ctx context.Context, _input *InputPlateInitializerInput, _output *InputPlateInitializerOutput) {

	plate, err := mixer.ParseInputPlateFile(_input.Filename)
	if err != nil {
		panic("Error reading input plate file:" + _input.Filename)
	}

	// required to allow use of these inputs
	execute.SetInputPlate(_ctx, plate)

	_output.Plate = plate

	_output.Status = "Initialized plate: " + _output.Plate.Name()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _InputPlateInitializerAnalysis(_ctx context.Context, _input *InputPlateInitializerInput, _output *InputPlateInitializerOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _InputPlateInitializerValidation(_ctx context.Context, _input *InputPlateInitializerInput, _output *InputPlateInitializerOutput) {

}
func _InputPlateInitializerRun(_ctx context.Context, input *InputPlateInitializerInput) *InputPlateInitializerOutput {
	output := &InputPlateInitializerOutput{}
	_InputPlateInitializerSetup(_ctx, input)
	_InputPlateInitializerSteps(_ctx, input, output)
	_InputPlateInitializerAnalysis(_ctx, input, output)
	_InputPlateInitializerValidation(_ctx, input, output)
	return output
}

func InputPlateInitializerRunSteps(_ctx context.Context, input *InputPlateInitializerInput) *InputPlateInitializerSOutput {
	soutput := &InputPlateInitializerSOutput{}
	output := _InputPlateInitializerRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func InputPlateInitializerNew() interface{} {
	return &InputPlateInitializerElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &InputPlateInitializerInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _InputPlateInitializerRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &InputPlateInitializerInput{},
			Out: &InputPlateInitializerOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type InputPlateInitializerElement struct {
	inject.CheckedRunner
}

type InputPlateInitializerInput struct {
	Filename string
}

type InputPlateInitializerOutput struct {
	Plate  *wtype.LHPlate
	Status string
}

type InputPlateInitializerSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		Plate *wtype.LHPlate
	}
}

func init() {
	if err := addComponent(component.Component{Name: "InputPlateInitializer",
		Constructor: InputPlateInitializerNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/BS/InputPlateInitializer.an",
			Params: []component.ParamDesc{
				{Name: "Filename", Desc: "", Kind: "Parameters"},
				{Name: "Plate", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
