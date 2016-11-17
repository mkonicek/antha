package lib

import (
	"fmt"
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

func _PrintnameRequirements() {

}

// Actions to perform before protocol itself
func _PrintnameSetup(_ctx context.Context, _input *PrintnameInput) {

}

// Core process of the protocol: steps to be performed for each input
func _PrintnameSteps(_ctx context.Context, _input *PrintnameInput, _output *PrintnameOutput) {
	if _input.Name == "Michael Jackson" {
		_output.Fullname = fmt.Sprintln(_input.Name)
	} else {
		_output.Fullname = "there's only one Michael Jackson, we accept no imitators"
	}
}

// Actions to perform after steps block to analyze data
func _PrintnameAnalysis(_ctx context.Context, _input *PrintnameInput, _output *PrintnameOutput) {

}

func _PrintnameValidation(_ctx context.Context, _input *PrintnameInput, _output *PrintnameOutput) {

}
func _PrintnameRun(_ctx context.Context, input *PrintnameInput) *PrintnameOutput {
	output := &PrintnameOutput{}
	_PrintnameSetup(_ctx, input)
	_PrintnameSteps(_ctx, input, output)
	_PrintnameAnalysis(_ctx, input, output)
	_PrintnameValidation(_ctx, input, output)
	return output
}

func PrintnameRunSteps(_ctx context.Context, input *PrintnameInput) *PrintnameSOutput {
	soutput := &PrintnameSOutput{}
	output := _PrintnameRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PrintnameNew() interface{} {
	return &PrintnameElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PrintnameInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PrintnameRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PrintnameInput{},
			Out: &PrintnameOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PrintnameElement struct {
	inject.CheckedRunner
}

type PrintnameInput struct {
	Name string
}

type PrintnameOutput struct {
	Fullname string
}

type PrintnameSOutput struct {
	Data struct {
		Fullname string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Printname",
		Constructor: PrintnameNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/Printname/Printname.an",
			Params: []component.ParamDesc{
				{Name: "Name", Desc: "", Kind: "Parameters"},
				{Name: "Fullname", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
