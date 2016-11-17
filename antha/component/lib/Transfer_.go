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

func _TransferRequirements() {

}

// Conditions to run on startup
func _TransferSetup(_ctx context.Context, _input *TransferInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _TransferSteps(_ctx context.Context, _input *TransferInput, _output *TransferOutput) {

	sample := mixer.Sample(_input.Startingsolution, _input.LiquidVolume)
	_output.FinalSolution = execute.MixInto(_ctx, _input.OutPlate, "", sample)

	_output.Status = _input.LiquidVolume.ToString() + " of " + _input.Liquidname + " was mixed into " + _input.OutPlate.Type

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _TransferAnalysis(_ctx context.Context, _input *TransferInput, _output *TransferOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _TransferValidation(_ctx context.Context, _input *TransferInput, _output *TransferOutput) {

}
func _TransferRun(_ctx context.Context, input *TransferInput) *TransferOutput {
	output := &TransferOutput{}
	_TransferSetup(_ctx, input)
	_TransferSteps(_ctx, input, output)
	_TransferAnalysis(_ctx, input, output)
	_TransferValidation(_ctx, input, output)
	return output
}

func TransferRunSteps(_ctx context.Context, input *TransferInput) *TransferSOutput {
	soutput := &TransferSOutput{}
	output := _TransferRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func TransferNew() interface{} {
	return &TransferElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &TransferInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _TransferRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &TransferInput{},
			Out: &TransferOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type TransferElement struct {
	inject.CheckedRunner
}

type TransferInput struct {
	LiquidVolume     wunit.Volume
	Liquidname       string
	OutPlate         *wtype.LHPlate
	Startingsolution *wtype.LHComponent
}

type TransferOutput struct {
	FinalSolution *wtype.LHComponent
	Status        string
}

type TransferSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		FinalSolution *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Transfer",
		Constructor: TransferNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/Transfer/Transfer.an",
			Params: []component.ParamDesc{
				{Name: "LiquidVolume", Desc: "", Kind: "Parameters"},
				{Name: "Liquidname", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Startingsolution", Desc: "", Kind: "Inputs"},
				{Name: "FinalSolution", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
