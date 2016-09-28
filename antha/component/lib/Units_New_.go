// demo of how to create units from raw values and unit strings
package lib

import (
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol

// Data which is returned from this protocol

// Physical inputs to this protocol

// Physical outputs from this protocol

func _Units_NewRequirements() {

}

// Actions to perform before protocol itself
func _Units_NewSetup(_ctx context.Context, _input *Units_NewInput) {

}

// Core process of the protocol: steps to be performed for each input
func _Units_NewSteps(_ctx context.Context, _input *Units_NewInput, _output *Units_NewOutput) {
	// this is importing the NewVolume function from the wunit package
	_output.MyVolume = wunit.NewVolume(_input.MyValue, _input.MyUnit)

	// Exercise: adjust the protocol to make a concentration instead
}

// Actions to perform after steps block to analyze data
func _Units_NewAnalysis(_ctx context.Context, _input *Units_NewInput, _output *Units_NewOutput) {

}

func _Units_NewValidation(_ctx context.Context, _input *Units_NewInput, _output *Units_NewOutput) {

}
func _Units_NewRun(_ctx context.Context, input *Units_NewInput) *Units_NewOutput {
	output := &Units_NewOutput{}
	_Units_NewSetup(_ctx, input)
	_Units_NewSteps(_ctx, input, output)
	_Units_NewAnalysis(_ctx, input, output)
	_Units_NewValidation(_ctx, input, output)
	return output
}

func Units_NewRunSteps(_ctx context.Context, input *Units_NewInput) *Units_NewSOutput {
	soutput := &Units_NewSOutput{}
	output := _Units_NewRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Units_NewNew() interface{} {
	return &Units_NewElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Units_NewInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Units_NewRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Units_NewInput{},
			Out: &Units_NewOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Units_NewElement struct {
	inject.CheckedRunner
}

type Units_NewInput struct {
	MyUnit  string
	MyValue float64
}

type Units_NewOutput struct {
	MyVolume wunit.Volume
}

type Units_NewSOutput struct {
	Data struct {
		MyVolume wunit.Volume
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Units_New",
		Constructor: Units_NewNew,
		Desc: component.ComponentDesc{
			Desc: "demo of how to create units from raw values and unit strings\n",
			Path: "antha/component/an/AnthaAcademy/Lesson0_Units/A_units_new.an",
			Params: []component.ParamDesc{
				{Name: "MyUnit", Desc: "", Kind: "Parameters"},
				{Name: "MyValue", Desc: "", Kind: "Parameters"},
				{Name: "MyVolume", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
