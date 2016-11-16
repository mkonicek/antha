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

//OutPlate *wtype.LHPlate

// Physical outputs from this protocol with types

func _Assay_quenchRequirements() {

}

// Conditions to run on startup
func _Assay_quenchSetup(_ctx context.Context, _input *Assay_quenchInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _Assay_quenchSteps(_ctx context.Context, _input *Assay_quenchInput, _output *Assay_quenchOutput) {

	substrate := mixer.Sample(_input.Substrate, _input.SubstrateVolume)
	enzyme := mixer.Sample(_input.Enzyme, _input.EnzymeVolume)

	// MixTo(platetype string, address string, platenum int, components ...*wtype.LHComponent)

	reaction := execute.MixTo(_ctx, _input.OutPlate, "", 1, substrate, enzyme)

	finishedreaction := execute.Incubate(_ctx, reaction, _input.ReactionTemp, _input.ReactionTime, true)

	quench := mixer.Sample(_input.Quenchingagent, _input.QuenchingagentVolume)

	_output.QuenchedReaction = execute.Mix(_ctx, finishedreaction, quench)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Assay_quenchAnalysis(_ctx context.Context, _input *Assay_quenchInput, _output *Assay_quenchOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _Assay_quenchValidation(_ctx context.Context, _input *Assay_quenchInput, _output *Assay_quenchOutput) {

}
func _Assay_quenchRun(_ctx context.Context, input *Assay_quenchInput) *Assay_quenchOutput {
	output := &Assay_quenchOutput{}
	_Assay_quenchSetup(_ctx, input)
	_Assay_quenchSteps(_ctx, input, output)
	_Assay_quenchAnalysis(_ctx, input, output)
	_Assay_quenchValidation(_ctx, input, output)
	return output
}

func Assay_quenchRunSteps(_ctx context.Context, input *Assay_quenchInput) *Assay_quenchSOutput {
	soutput := &Assay_quenchSOutput{}
	output := _Assay_quenchRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Assay_quenchNew() interface{} {
	return &Assay_quenchElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Assay_quenchInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Assay_quenchRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Assay_quenchInput{},
			Out: &Assay_quenchOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Assay_quenchElement struct {
	inject.CheckedRunner
}

type Assay_quenchInput struct {
	Enzyme               *wtype.LHComponent
	EnzymeVolume         wunit.Volume
	OutPlate             string
	Quenchingagent       *wtype.LHComponent
	QuenchingagentVolume wunit.Volume
	ReactionTemp         wunit.Temperature
	ReactionTime         wunit.Time
	Substrate            *wtype.LHComponent
	SubstrateVolume      wunit.Volume
}

type Assay_quenchOutput struct {
	QuenchedReaction *wtype.LHComponent
}

type Assay_quenchSOutput struct {
	Data struct {
	}
	Outputs struct {
		QuenchedReaction *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Assay_quench",
		Constructor: Assay_quenchNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/AssaySetUp/QuenchedReaction.an",
			Params: []component.ParamDesc{
				{Name: "Enzyme", Desc: "", Kind: "Inputs"},
				{Name: "EnzymeVolume", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Parameters"},
				{Name: "Quenchingagent", Desc: "", Kind: "Inputs"},
				{Name: "QuenchingagentVolume", Desc: "", Kind: "Parameters"},
				{Name: "ReactionTemp", Desc: "", Kind: "Parameters"},
				{Name: "ReactionTime", Desc: "", Kind: "Parameters"},
				{Name: "Substrate", Desc: "", Kind: "Inputs"},
				{Name: "SubstrateVolume", Desc: "", Kind: "Parameters"},
				{Name: "QuenchedReaction", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
