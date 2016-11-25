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

func _MixIntoSameRequirements() {

}

// Conditions to run on startup
func _MixIntoSameSetup(_ctx context.Context, _input *MixIntoSameInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _MixIntoSameSteps(_ctx context.Context, _input *MixIntoSameInput, _output *MixIntoSameOutput) {
	sample1 := mixer.Sample(_input.Liquid, _input.Volume1)
	platedSample := execute.MixInto(_ctx, _input.OutPlate, "A1", sample1)
	sample2 := mixer.Sample(platedSample, _input.Volume2)
	execute.MixInto(_ctx, _input.OutPlate, "A2", sample2)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MixIntoSameAnalysis(_ctx context.Context, _input *MixIntoSameInput, _output *MixIntoSameOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _MixIntoSameValidation(_ctx context.Context, _input *MixIntoSameInput, _output *MixIntoSameOutput) {

}
func _MixIntoSameRun(_ctx context.Context, input *MixIntoSameInput) *MixIntoSameOutput {
	output := &MixIntoSameOutput{}
	_MixIntoSameSetup(_ctx, input)
	_MixIntoSameSteps(_ctx, input, output)
	_MixIntoSameAnalysis(_ctx, input, output)
	_MixIntoSameValidation(_ctx, input, output)
	return output
}

func MixIntoSameRunSteps(_ctx context.Context, input *MixIntoSameInput) *MixIntoSameSOutput {
	soutput := &MixIntoSameSOutput{}
	output := _MixIntoSameRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MixIntoSameNew() interface{} {
	return &MixIntoSameElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MixIntoSameInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MixIntoSameRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MixIntoSameInput{},
			Out: &MixIntoSameOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type MixIntoSameElement struct {
	inject.CheckedRunner
}

type MixIntoSameInput struct {
	Liquid   *wtype.LHComponent
	OutPlate *wtype.LHPlate
	Volume1  wunit.Volume
	Volume2  wunit.Volume
}

type MixIntoSameOutput struct {
	Status string
}

type MixIntoSameSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MixIntoSame",
		Constructor: MixIntoSameNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/BS/MixIntoSame.an",
			Params: []component.ParamDesc{
				{Name: "Liquid", Desc: "", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Volume1", Desc: "", Kind: "Parameters"},
				{Name: "Volume2", Desc: "", Kind: "Parameters"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
