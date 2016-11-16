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
//ComponentNames []string

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// optional if nil this is ignored

// Physical outputs from this protocol with types

func _MastermixRequirements() {
}

// Conditions to run on startup
func _MastermixSetup(_ctx context.Context, _input *MastermixInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _MastermixSteps(_ctx context.Context, _input *MastermixInput, _output *MastermixOutput) {
	if len(_input.OtherComponents) != len(_input.OtherComponentVolumes) {
		execute.Errorf(_ctx, "%d != %d", len(_input.OtherComponents), len(_input.OtherComponentVolumes))
	}

	mastermixes := make([]*wtype.LHComponent, 0)

	if _input.AliquotbyRow {
		execute.Errorf(_ctx, "MixTo based method coming soon!")
	} else {
		for i := 0; i < _input.NumberofMastermixes; i++ {

			eachmastermix := make([]*wtype.LHComponent, 0)

			if _input.Buffer != nil {
				bufferSample := mixer.SampleForTotalVolume(_input.Buffer, _input.TotalVolumeperMastermix)
				eachmastermix = append(eachmastermix, bufferSample)
			}

			for k, component := range _input.OtherComponents {
				if k == len(_input.OtherComponents) {
					component.Type = wtype.LTNeedToMix //"NeedToMix"
				}
				componentSample := mixer.Sample(component, _input.OtherComponentVolumes[k])
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
func _MastermixAnalysis(_ctx context.Context, _input *MastermixInput, _output *MastermixOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _MastermixValidation(_ctx context.Context, _input *MastermixInput, _output *MastermixOutput) {
}
func _MastermixRun(_ctx context.Context, input *MastermixInput) *MastermixOutput {
	output := &MastermixOutput{}
	_MastermixSetup(_ctx, input)
	_MastermixSteps(_ctx, input, output)
	_MastermixAnalysis(_ctx, input, output)
	_MastermixValidation(_ctx, input, output)
	return output
}

func MastermixRunSteps(_ctx context.Context, input *MastermixInput) *MastermixSOutput {
	soutput := &MastermixSOutput{}
	output := _MastermixRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MastermixNew() interface{} {
	return &MastermixElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MastermixInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MastermixRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MastermixInput{},
			Out: &MastermixOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type MastermixElement struct {
	inject.CheckedRunner
}

type MastermixInput struct {
	AliquotbyRow            bool
	Buffer                  *wtype.LHComponent
	Inplate                 *wtype.LHPlate
	NumberofMastermixes     int
	OtherComponentVolumes   []wunit.Volume
	OtherComponents         []*wtype.LHComponent
	OutPlate                *wtype.LHPlate
	TotalVolumeperMastermix wunit.Volume
}

type MastermixOutput struct {
	Mastermixes []*wtype.LHComponent
	Status      string
}

type MastermixSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		Mastermixes []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Mastermix",
		Constructor: MastermixNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/MakeMastermix/Mastermix.an",
			Params: []component.ParamDesc{
				{Name: "AliquotbyRow", Desc: "", Kind: "Parameters"},
				{Name: "Buffer", Desc: "optional if nil this is ignored\n", Kind: "Inputs"},
				{Name: "Inplate", Desc: "", Kind: "Inputs"},
				{Name: "NumberofMastermixes", Desc: "", Kind: "Parameters"},
				{Name: "OtherComponentVolumes", Desc: "ComponentNames []string\n", Kind: "Parameters"},
				{Name: "OtherComponents", Desc: "", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "TotalVolumeperMastermix", Desc: "if buffer is being added\n", Kind: "Parameters"},
				{Name: "Mastermixes", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
