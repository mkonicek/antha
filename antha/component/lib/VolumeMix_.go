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

func _VolumeMixRequirements() {

}

// Conditions to run on startup
func _VolumeMixSetup(_ctx context.Context, _input *VolumeMixInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _VolumeMixSteps(_ctx context.Context, _input *VolumeMixInput, _output *VolumeMixOutput) {
	sample := mixer.Sample(_input.Liquid, _input.Volume)
	execute.Mix(_ctx, sample)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _VolumeMixAnalysis(_ctx context.Context, _input *VolumeMixInput, _output *VolumeMixOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _VolumeMixValidation(_ctx context.Context, _input *VolumeMixInput, _output *VolumeMixOutput) {

}
func _VolumeMixRun(_ctx context.Context, input *VolumeMixInput) *VolumeMixOutput {
	output := &VolumeMixOutput{}
	_VolumeMixSetup(_ctx, input)
	_VolumeMixSteps(_ctx, input, output)
	_VolumeMixAnalysis(_ctx, input, output)
	_VolumeMixValidation(_ctx, input, output)
	return output
}

func VolumeMixRunSteps(_ctx context.Context, input *VolumeMixInput) *VolumeMixSOutput {
	soutput := &VolumeMixSOutput{}
	output := _VolumeMixRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func VolumeMixNew() interface{} {
	return &VolumeMixElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &VolumeMixInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _VolumeMixRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &VolumeMixInput{},
			Out: &VolumeMixOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type VolumeMixElement struct {
	inject.CheckedRunner
}

type VolumeMixInput struct {
	Liquid   *wtype.LHComponent
	OutPlate *wtype.LHPlate
	Volume   wunit.Volume
}

type VolumeMixOutput struct {
	Status string
}

type VolumeMixSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "VolumeMix",
		Constructor: VolumeMixNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/BS/VolumeMix.an",
			Params: []component.ParamDesc{
				{Name: "Liquid", Desc: "", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Volume", Desc: "", Kind: "Parameters"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
