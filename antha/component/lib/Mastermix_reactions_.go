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

// if buffer is being added

// add as many as possible option e.g. if == -1

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// optional if nil this is ignored

// Physical outputs from this protocol with types

func _Mastermix_reactionsRequirements() {
}

// Conditions to run on startup
func _Mastermix_reactionsSetup(_ctx context.Context, _input *Mastermix_reactionsInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _Mastermix_reactionsSteps(_ctx context.Context, _input *Mastermix_reactionsInput, _output *Mastermix_reactionsOutput) {

	// work out volume to top up to in each case (per reaction) in l:
	topupVolumeperreacttion := _input.TotalVolumeperreaction.SIValue() - _input.VolumetoLeaveforDNAperreaction.SIValue()

	// multiply by number of reactions per mastermix
	topupVolume := wunit.NewVolume(float64(_input.Reactionspermastermix)*topupVolumeperreacttion, "l")

	if len(_input.Components) != len(_input.ComponentVolumesperReaction) {
		panic("len(Components) != len(OtherComponentVolumes)")
	}

	mastermixes := make([]*wtype.LHComponent, 0)

	if _input.AliquotbyRow {
		panic("MixTo based method coming soon!")
	} else {
		for i := 0; i < _input.NumberofMastermixes; i++ {

			eachmastermix := make([]*wtype.LHComponent, 0)

			if _input.TopUpBuffer != nil {
				bufferSample := mixer.SampleForTotalVolume(_input.TopUpBuffer, topupVolume)
				eachmastermix = append(eachmastermix, bufferSample)
			}

			for k, component := range _input.Components {
				if k == len(_input.Components) {
					component.Type = wtype.LTNeedToMix //"NeedToMix"
				}

				// multiply volume of each component by number of reactions per mastermix
				adjustedvol := wunit.NewVolume(float64(_input.Reactionspermastermix)*_input.ComponentVolumesperReaction[k].SIValue(), "l")

				componentSample := mixer.Sample(component, adjustedvol)
				eachmastermix = append(eachmastermix, componentSample)
			}

			mastermix := execute.MixInto(_ctx, _input.OutPlate, "", eachmastermix...)
			mastermixes = append(mastermixes, mastermix)

		}

	}
	_output.Mastermixes = mastermixes

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Mastermix_reactionsAnalysis(_ctx context.Context, _input *Mastermix_reactionsInput, _output *Mastermix_reactionsOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _Mastermix_reactionsValidation(_ctx context.Context, _input *Mastermix_reactionsInput, _output *Mastermix_reactionsOutput) {
}
func _Mastermix_reactionsRun(_ctx context.Context, input *Mastermix_reactionsInput) *Mastermix_reactionsOutput {
	output := &Mastermix_reactionsOutput{}
	_Mastermix_reactionsSetup(_ctx, input)
	_Mastermix_reactionsSteps(_ctx, input, output)
	_Mastermix_reactionsAnalysis(_ctx, input, output)
	_Mastermix_reactionsValidation(_ctx, input, output)
	return output
}

func Mastermix_reactionsRunSteps(_ctx context.Context, input *Mastermix_reactionsInput) *Mastermix_reactionsSOutput {
	soutput := &Mastermix_reactionsSOutput{}
	output := _Mastermix_reactionsRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Mastermix_reactionsNew() interface{} {
	return &Mastermix_reactionsElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Mastermix_reactionsInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Mastermix_reactionsRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Mastermix_reactionsInput{},
			Out: &Mastermix_reactionsOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Mastermix_reactionsElement struct {
	inject.CheckedRunner
}

type Mastermix_reactionsInput struct {
	AliquotbyRow                   bool
	ComponentVolumesperReaction    []wunit.Volume
	Components                     []*wtype.LHComponent
	Inplate                        *wtype.LHPlate
	NumberofMastermixes            int
	OutPlate                       *wtype.LHPlate
	Reactionspermastermix          int
	TopUpBuffer                    *wtype.LHComponent
	TotalVolumeperreaction         wunit.Volume
	VolumetoLeaveforDNAperreaction wunit.Volume
}

type Mastermix_reactionsOutput struct {
	Mastermixes []*wtype.LHComponent
	Status      string
}

type Mastermix_reactionsSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		Mastermixes []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Mastermix_reactions",
		Constructor: Mastermix_reactionsNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/MakeMastermix/Mastermix_reactions.an",
			Params: []component.ParamDesc{
				{Name: "AliquotbyRow", Desc: "", Kind: "Parameters"},
				{Name: "ComponentVolumesperReaction", Desc: "", Kind: "Parameters"},
				{Name: "Components", Desc: "", Kind: "Inputs"},
				{Name: "Inplate", Desc: "", Kind: "Inputs"},
				{Name: "NumberofMastermixes", Desc: "add as many as possible option e.g. if == -1\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Reactionspermastermix", Desc: "", Kind: "Parameters"},
				{Name: "TopUpBuffer", Desc: "optional if nil this is ignored\n", Kind: "Inputs"},
				{Name: "TotalVolumeperreaction", Desc: "if buffer is being added\n", Kind: "Parameters"},
				{Name: "VolumetoLeaveforDNAperreaction", Desc: "", Kind: "Parameters"},
				{Name: "Mastermixes", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
