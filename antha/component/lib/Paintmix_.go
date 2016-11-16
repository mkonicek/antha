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

func _PaintmixRequirements() {
}

// Conditions to run on startup
func _PaintmixSetup(_ctx context.Context, _input *PaintmixInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _PaintmixSteps(_ctx context.Context, _input *PaintmixInput, _output *PaintmixOutput) {

	reactions := make([]*wtype.LHComponent, 0)

	for i := 0; i < _input.Numberofcopies; i++ {
		eachreaction := make([]*wtype.LHComponent, 0)
		col1Sample := mixer.Sample(_input.Colour1, _input.Colour1vol)
		eachreaction = append(eachreaction, col1Sample)
		col2Sample := mixer.Sample(_input.Colour2, _input.Colour2vol)
		eachreaction = append(eachreaction, col2Sample)
		reaction := execute.MixInto(_ctx, _input.OutPlate, "", eachreaction...)
		reactions = append(reactions, reaction)

	}
	_output.NewColours = reactions

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PaintmixAnalysis(_ctx context.Context, _input *PaintmixInput, _output *PaintmixOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PaintmixValidation(_ctx context.Context, _input *PaintmixInput, _output *PaintmixOutput) {
}
func _PaintmixRun(_ctx context.Context, input *PaintmixInput) *PaintmixOutput {
	output := &PaintmixOutput{}
	_PaintmixSetup(_ctx, input)
	_PaintmixSteps(_ctx, input, output)
	_PaintmixAnalysis(_ctx, input, output)
	_PaintmixValidation(_ctx, input, output)
	return output
}

func PaintmixRunSteps(_ctx context.Context, input *PaintmixInput) *PaintmixSOutput {
	soutput := &PaintmixSOutput{}
	output := _PaintmixRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PaintmixNew() interface{} {
	return &PaintmixElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PaintmixInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PaintmixRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PaintmixInput{},
			Out: &PaintmixOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PaintmixElement struct {
	inject.CheckedRunner
}

type PaintmixInput struct {
	Colour1        *wtype.LHComponent
	Colour1vol     wunit.Volume
	Colour2        *wtype.LHComponent
	Colour2vol     wunit.Volume
	Numberofcopies int
	OutPlate       *wtype.LHPlate
}

type PaintmixOutput struct {
	NewColours []*wtype.LHComponent
	Status     string
}

type PaintmixSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		NewColours []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Paintmix",
		Constructor: PaintmixNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/Colourmix/Paintmix.an",
			Params: []component.ParamDesc{
				{Name: "Colour1", Desc: "", Kind: "Inputs"},
				{Name: "Colour1vol", Desc: "", Kind: "Parameters"},
				{Name: "Colour2", Desc: "", Kind: "Inputs"},
				{Name: "Colour2vol", Desc: "", Kind: "Parameters"},
				{Name: "Numberofcopies", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "NewColours", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
