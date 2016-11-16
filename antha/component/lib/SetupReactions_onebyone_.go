// this protocol will set up a specified number of reactions one reaction at a time, i.e. in the following format:
// add all components into reaction 1 location,
// add all components into reaction 2 location,
// ...,
// add all components into reaction n location
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

func _SetupReactions_onebyoneRequirements() {
}

// Conditions to run on startup
func _SetupReactions_onebyoneSetup(_ctx context.Context, _input *SetupReactions_onebyoneInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _SetupReactions_onebyoneSteps(_ctx context.Context, _input *SetupReactions_onebyoneInput, _output *SetupReactions_onebyoneOutput) {

	reactions := make([]*wtype.LHComponent, 0)

	for i := 0; i < _input.NumberofReactions; i++ {
		// creating this eachreaction slice and appending with each sample is the key to ensuring a reaction is made one at a time
		// note that for each reaction this is reinitialised back to an empty slice
		eachreaction := make([]*wtype.LHComponent, 0)

		bufferSample := mixer.SampleForTotalVolume(_input.Buffer, _input.TotalVolume)
		eachreaction = append(eachreaction, bufferSample)

		subSample := mixer.Sample(_input.Substrate, _input.SubstrateVolume)
		eachreaction = append(eachreaction, subSample)

		enzSample := mixer.Sample(_input.Enzyme, _input.EnzymeVolume)
		eachreaction = append(eachreaction, enzSample)

		// the Mix command (in this case MixInto) is used once for all the samples
		// this ensures all components are mixed for reaction x before moving on to reaction x + 1
		reaction := execute.MixInto(_ctx, _input.OutPlate, "", eachreaction...)
		reactions = append(reactions, reaction)

	}
	_output.Reactions = reactions

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _SetupReactions_onebyoneAnalysis(_ctx context.Context, _input *SetupReactions_onebyoneInput, _output *SetupReactions_onebyoneOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _SetupReactions_onebyoneValidation(_ctx context.Context, _input *SetupReactions_onebyoneInput, _output *SetupReactions_onebyoneOutput) {
}
func _SetupReactions_onebyoneRun(_ctx context.Context, input *SetupReactions_onebyoneInput) *SetupReactions_onebyoneOutput {
	output := &SetupReactions_onebyoneOutput{}
	_SetupReactions_onebyoneSetup(_ctx, input)
	_SetupReactions_onebyoneSteps(_ctx, input, output)
	_SetupReactions_onebyoneAnalysis(_ctx, input, output)
	_SetupReactions_onebyoneValidation(_ctx, input, output)
	return output
}

func SetupReactions_onebyoneRunSteps(_ctx context.Context, input *SetupReactions_onebyoneInput) *SetupReactions_onebyoneSOutput {
	soutput := &SetupReactions_onebyoneSOutput{}
	output := _SetupReactions_onebyoneRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SetupReactions_onebyoneNew() interface{} {
	return &SetupReactions_onebyoneElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SetupReactions_onebyoneInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SetupReactions_onebyoneRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SetupReactions_onebyoneInput{},
			Out: &SetupReactions_onebyoneOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type SetupReactions_onebyoneElement struct {
	inject.CheckedRunner
}

type SetupReactions_onebyoneInput struct {
	Buffer            *wtype.LHComponent
	Enzyme            *wtype.LHComponent
	EnzymeVolume      wunit.Volume
	NumberofReactions int
	OutPlate          *wtype.LHPlate
	Substrate         *wtype.LHComponent
	SubstrateVolume   wunit.Volume
	TotalVolume       wunit.Volume
}

type SetupReactions_onebyoneOutput struct {
	Reactions []*wtype.LHComponent
	Status    string
}

type SetupReactions_onebyoneSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		Reactions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "SetupReactions_onebyone",
		Constructor: SetupReactions_onebyoneNew,
		Desc: component.ComponentDesc{
			Desc: "this protocol will set up a specified number of reactions one reaction at a time, i.e. in the following format:\nadd all components into reaction 1 location,\nadd all components into reaction 2 location,\n...,\nadd all components into reaction n location\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/AnthaAcademy/Lesson3_MixPart2/B_Assaysetup_reactionbyreaction.an",
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
