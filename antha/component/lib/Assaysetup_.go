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

func _AssaysetupRequirements() {
}

// Conditions to run on startup
func _AssaysetupSetup(_ctx context.Context, _input *AssaysetupInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _AssaysetupSteps(_ctx context.Context, _input *AssaysetupInput, _output *AssaysetupOutput) {

	reactions := make([]*wtype.LHComponent, 0)

	for i := 0; i < _input.NumberofReactions; i++ {
		eachreaction := make([]*wtype.LHComponent, 0)
		bufferSample := mixer.SampleForTotalVolume(_input.Buffer, _input.TotalVolume)
		eachreaction = append(eachreaction, bufferSample)
		subSample := mixer.Sample(_input.Substrate, _input.SubstrateVolume)
		eachreaction = append(eachreaction, subSample)
		enzSample := mixer.Sample(_input.Enzyme, _input.EnzymeVolume)
		eachreaction = append(eachreaction, enzSample)
		reaction := execute.MixInto(_ctx, _input.OutPlate, "", eachreaction...)
		reactions = append(reactions, reaction)

	}
	_output.Reactions = reactions

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AssaysetupAnalysis(_ctx context.Context, _input *AssaysetupInput, _output *AssaysetupOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AssaysetupValidation(_ctx context.Context, _input *AssaysetupInput, _output *AssaysetupOutput) {
}
func _AssaysetupRun(_ctx context.Context, input *AssaysetupInput) *AssaysetupOutput {
	output := &AssaysetupOutput{}
	_AssaysetupSetup(_ctx, input)
	_AssaysetupSteps(_ctx, input, output)
	_AssaysetupAnalysis(_ctx, input, output)
	_AssaysetupValidation(_ctx, input, output)
	return output
}

func AssaysetupRunSteps(_ctx context.Context, input *AssaysetupInput) *AssaysetupSOutput {
	soutput := &AssaysetupSOutput{}
	output := _AssaysetupRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AssaysetupNew() interface{} {
	return &AssaysetupElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AssaysetupInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AssaysetupRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AssaysetupInput{},
			Out: &AssaysetupOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type AssaysetupElement struct {
	inject.CheckedRunner
}

type AssaysetupInput struct {
	Buffer            *wtype.LHComponent
	Enzyme            *wtype.LHComponent
	EnzymeVolume      wunit.Volume
	NumberofReactions int
	OutPlate          *wtype.LHPlate
	Substrate         *wtype.LHComponent
	SubstrateVolume   wunit.Volume
	TotalVolume       wunit.Volume
}

type AssaysetupOutput struct {
	Reactions []*wtype.LHComponent
	Status    string
}

type AssaysetupSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		Reactions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Assaysetup",
		Constructor: AssaysetupNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/AssaySetUp/Assaysetup.an",
			Params: []component.ParamDesc{
				{Name: "Buffer", Desc: "", Kind: "Inputs"},
				{Name: "Enzyme", Desc: "", Kind: "Inputs"},
				{Name: "EnzymeVolume", Desc: "", Kind: "Parameters"},
				{Name: "NumberofReactions", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Substrate", Desc: "", Kind: "Inputs"},
				{Name: "SubstrateVolume", Desc: "", Kind: "Parameters"},
				{Name: "TotalVolume", Desc: "", Kind: "Parameters"},
				{Name: "Reactions", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
