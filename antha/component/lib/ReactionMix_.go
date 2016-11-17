package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

func _ReactionMixRequirements() {}

// Conditions to run on startup
func _ReactionMixSetup(_ctx context.Context, _input *ReactionMixInput) {}

func _ReactionMixSteps(_ctx context.Context, _input *ReactionMixInput, _output *ReactionMixOutput) {

	fmt.Println("Components:", _input.Components)

	samples := make([]*wtype.LHComponent, 0)

	VectorS := mixer.Sample(_input.Vector, _input.VectorV)
	BufferS := mixer.Sample(_input.Buffer, _input.BufferV)
	LigaseS := mixer.Sample(_input.Ligase, _input.LigaseV)
	ATPS := mixer.Sample(_input.ATP, _input.ATPV)
	RES := mixer.Sample(_input.RE, _input.REV)

	//com := []wunit.Volume{VectorV, BufferV, LigaseV, ATPV, REV}
	//WaterS := mixer.TopUpVolume(Water, com, ReactionVolume)
	WaterS := mixer.SampleForTotalVolume(_input.Water, _input.ReactionVolume)

	ComponentsS := mixer.MultiSample(_input.Components, _input.ComponentsV)

	samples = append(samples, VectorS, BufferS, LigaseS, ATPS, RES, WaterS)
	samples = append(samples, ComponentsS...)

	// Incubate
	_output.Reaction = execute.Incubate(_ctx, mixer.MixInto(_input.OutPlate, "", samples...), _input.ReactionTemp, _input.ReactionTime, false)

}

func _ReactionMixAnalysis(_ctx context.Context, _input *ReactionMixInput, _output *ReactionMixOutput) {
}

func _ReactionMixValidation(_ctx context.Context, _input *ReactionMixInput, _output *ReactionMixOutput) {
}
func _ReactionMixRun(_ctx context.Context, input *ReactionMixInput) *ReactionMixOutput {
	output := &ReactionMixOutput{}
	_ReactionMixSetup(_ctx, input)
	_ReactionMixSteps(_ctx, input, output)
	_ReactionMixAnalysis(_ctx, input, output)
	_ReactionMixValidation(_ctx, input, output)
	return output
}

func ReactionMixRunSteps(_ctx context.Context, input *ReactionMixInput) *ReactionMixSOutput {
	soutput := &ReactionMixSOutput{}
	output := _ReactionMixRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ReactionMixNew() interface{} {
	return &ReactionMixElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ReactionMixInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ReactionMixRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ReactionMixInput{},
			Out: &ReactionMixOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type ReactionMixElement struct {
	inject.CheckedRunner
}

type ReactionMixInput struct {
	ATP            *wtype.LHComponent
	ATPV           wunit.Volume
	Buffer         *wtype.LHComponent
	BufferV        wunit.Volume
	Components     []*wtype.LHComponent
	ComponentsV    []wunit.Volume
	InPlate        *wtype.LHPlate
	Ligase         *wtype.LHComponent
	LigaseV        wunit.Volume
	OutPlate       *wtype.LHPlate
	RE             *wtype.LHComponent
	REV            wunit.Volume
	ReactionTemp   wunit.Temperature
	ReactionTime   wunit.Time
	ReactionVolume wunit.Volume
	Vector         *wtype.LHComponent
	VectorV        wunit.Volume
	Water          *wtype.LHComponent
}

type ReactionMixOutput struct {
	Reaction *wtype.LHComponent
}

type ReactionMixSOutput struct {
	Data struct {
	}
	Outputs struct {
		Reaction *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "ReactionMix",
		Constructor: ReactionMixNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/DNA/GeneDesign/ReactionMix.an",
			Params: []component.ParamDesc{
				{Name: "ATP", Desc: "", Kind: "Inputs"},
				{Name: "ATPV", Desc: "", Kind: "Parameters"},
				{Name: "Buffer", Desc: "", Kind: "Inputs"},
				{Name: "BufferV", Desc: "", Kind: "Parameters"},
				{Name: "Components", Desc: "", Kind: "Inputs"},
				{Name: "ComponentsV", Desc: "", Kind: "Parameters"},
				{Name: "InPlate", Desc: "", Kind: "Inputs"},
				{Name: "Ligase", Desc: "", Kind: "Inputs"},
				{Name: "LigaseV", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "RE", Desc: "", Kind: "Inputs"},
				{Name: "REV", Desc: "", Kind: "Parameters"},
				{Name: "ReactionTemp", Desc: "", Kind: "Parameters"},
				{Name: "ReactionTime", Desc: "", Kind: "Parameters"},
				{Name: "ReactionVolume", Desc: "", Kind: "Parameters"},
				{Name: "Vector", Desc: "", Kind: "Inputs"},
				{Name: "VectorV", Desc: "", Kind: "Parameters"},
				{Name: "Water", Desc: "", Kind: "Inputs"},
				{Name: "Reaction", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
