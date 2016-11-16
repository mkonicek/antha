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

func _ProtocolName_from_an_fileRequirements() {

}

// Conditions to run on startup
func _ProtocolName_from_an_fileSetup(_ctx context.Context, _input *ProtocolName_from_an_fileInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _ProtocolName_from_an_fileSteps(_ctx context.Context, _input *ProtocolName_from_an_fileInput, _output *ProtocolName_from_an_fileOutput) {

	_output.OutputData = make([]string, 0)

	for i := 0; i < _input.ParameterVariableAsint; i++ {
		output := _input.ParameterVariableAsValuewithunit.ToString() + "of" + _input.ParameterVariablestring
		_output.OutputData = append(_output.OutputData, output)
	}
	sample := mixer.Sample(_input.InputVariable, _input.ParameterVariableAsValuewithunit)
	_output.PhysicalOutput = execute.MixInto(_ctx, _input.OutPlate, "", sample)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _ProtocolName_from_an_fileAnalysis(_ctx context.Context, _input *ProtocolName_from_an_fileInput, _output *ProtocolName_from_an_fileOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _ProtocolName_from_an_fileValidation(_ctx context.Context, _input *ProtocolName_from_an_fileInput, _output *ProtocolName_from_an_fileOutput) {

}
func _ProtocolName_from_an_fileRun(_ctx context.Context, input *ProtocolName_from_an_fileInput) *ProtocolName_from_an_fileOutput {
	output := &ProtocolName_from_an_fileOutput{}
	_ProtocolName_from_an_fileSetup(_ctx, input)
	_ProtocolName_from_an_fileSteps(_ctx, input, output)
	_ProtocolName_from_an_fileAnalysis(_ctx, input, output)
	_ProtocolName_from_an_fileValidation(_ctx, input, output)
	return output
}

func ProtocolName_from_an_fileRunSteps(_ctx context.Context, input *ProtocolName_from_an_fileInput) *ProtocolName_from_an_fileSOutput {
	soutput := &ProtocolName_from_an_fileSOutput{}
	output := _ProtocolName_from_an_fileRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ProtocolName_from_an_fileNew() interface{} {
	return &ProtocolName_from_an_fileElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ProtocolName_from_an_fileInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ProtocolName_from_an_fileRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ProtocolName_from_an_fileInput{},
			Out: &ProtocolName_from_an_fileOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type ProtocolName_from_an_fileElement struct {
	inject.CheckedRunner
}

type ProtocolName_from_an_fileInput struct {
	InputVariable                    *wtype.LHComponent
	OutPlate                         *wtype.LHPlate
	ParameterVariableAsValuewithunit wunit.Volume
	ParameterVariableAsint           int
	ParameterVariablestring          string
}

type ProtocolName_from_an_fileOutput struct {
	OutputData     []string
	PhysicalOutput *wtype.LHComponent
}

type ProtocolName_from_an_fileSOutput struct {
	Data struct {
		OutputData []string
	}
	Outputs struct {
		PhysicalOutput *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "ProtocolName_from_an_file",
		Constructor: ProtocolName_from_an_fileNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/default.an",
			Params: []component.ParamDesc{
				{Name: "InputVariable", Desc: "", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "ParameterVariableAsValuewithunit", Desc: "", Kind: "Parameters"},
				{Name: "ParameterVariableAsint", Desc: "", Kind: "Parameters"},
				{Name: "ParameterVariablestring", Desc: "", Kind: "Parameters"},
				{Name: "OutputData", Desc: "", Kind: "Data"},
				{Name: "PhysicalOutput", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
