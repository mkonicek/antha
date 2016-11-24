package lib

import (
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

func _WorkingPlateInitializerRequirements() {

}

// Conditions to run on startup
func _WorkingPlateInitializerSetup(_ctx context.Context, _input *WorkingPlateInitializerInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _WorkingPlateInitializerSteps(_ctx context.Context, _input *WorkingPlateInitializerInput, _output *WorkingPlateInitializerOutput) {

	_input.InPlate.PlateName = _input.Name
	_output.Plate = _input.InPlate

	_output.Status = "Initialized working plate"

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _WorkingPlateInitializerAnalysis(_ctx context.Context, _input *WorkingPlateInitializerInput, _output *WorkingPlateInitializerOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _WorkingPlateInitializerValidation(_ctx context.Context, _input *WorkingPlateInitializerInput, _output *WorkingPlateInitializerOutput) {

}
func _WorkingPlateInitializerRun(_ctx context.Context, input *WorkingPlateInitializerInput) *WorkingPlateInitializerOutput {
	output := &WorkingPlateInitializerOutput{}
	_WorkingPlateInitializerSetup(_ctx, input)
	_WorkingPlateInitializerSteps(_ctx, input, output)
	_WorkingPlateInitializerAnalysis(_ctx, input, output)
	_WorkingPlateInitializerValidation(_ctx, input, output)
	return output
}

func WorkingPlateInitializerRunSteps(_ctx context.Context, input *WorkingPlateInitializerInput) *WorkingPlateInitializerSOutput {
	soutput := &WorkingPlateInitializerSOutput{}
	output := _WorkingPlateInitializerRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func WorkingPlateInitializerNew() interface{} {
	return &WorkingPlateInitializerElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &WorkingPlateInitializerInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _WorkingPlateInitializerRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &WorkingPlateInitializerInput{},
			Out: &WorkingPlateInitializerOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type WorkingPlateInitializerElement struct {
	inject.CheckedRunner
}

type WorkingPlateInitializerInput struct {
	InPlate *wtype.LHPlate
	Name    string
}

type WorkingPlateInitializerOutput struct {
	Plate  *wtype.LHPlate
	Status string
}

type WorkingPlateInitializerSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		Plate *wtype.LHPlate
	}
}

func init() {
	if err := addComponent(component.Component{Name: "WorkingPlateInitializer",
		Constructor: WorkingPlateInitializerNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/BS/WorkingPlateInitializer.an",
			Params: []component.ParamDesc{
				{Name: "InPlate", Desc: "", Kind: "Parameters"},
				{Name: "Name", Desc: "", Kind: "Parameters"},
				{Name: "Plate", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
