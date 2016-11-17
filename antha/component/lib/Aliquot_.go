package lib

import (
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _AliquotRequirements() {

}

// Conditions to run on startup
func _AliquotSetup(_ctx context.Context, _input *AliquotInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _AliquotSteps(_ctx context.Context, _input *AliquotInput, _output *AliquotOutput) {

	// check there's enough liquid for the aliquots
	number := _input.SolutionVolume.SIValue() / _input.VolumePerAliquot.SIValue()
	possiblenumberofAliquots, _ := wutil.RoundDown(number)
	if possiblenumberofAliquots < _input.NumberofAliquots {
		panic("Not enough solution for this many aliquots")
	}

	aliquots := make([]*wtype.LHComponent, 0)

	for i := 0; i < _input.NumberofAliquots; i++ {
		if _input.Solution.TypeName() == "dna" {
			_input.Solution.Type = wtype.LTDoNotMix
		}
		aliquotSample := mixer.Sample(_input.Solution, _input.VolumePerAliquot)
		aliquot := execute.MixTo(_ctx, _input.OutPlatetype, "", 1, aliquotSample)
		aliquots = append(aliquots, aliquot)
	}
	_output.Aliquots = aliquots
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AliquotAnalysis(_ctx context.Context, _input *AliquotInput, _output *AliquotOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _AliquotValidation(_ctx context.Context, _input *AliquotInput, _output *AliquotOutput) {

}
func _AliquotRun(_ctx context.Context, input *AliquotInput) *AliquotOutput {
	output := &AliquotOutput{}
	_AliquotSetup(_ctx, input)
	_AliquotSteps(_ctx, input, output)
	_AliquotAnalysis(_ctx, input, output)
	_AliquotValidation(_ctx, input, output)
	return output
}

func AliquotRunSteps(_ctx context.Context, input *AliquotInput) *AliquotSOutput {
	soutput := &AliquotSOutput{}
	output := _AliquotRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AliquotNew() interface{} {
	return &AliquotElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AliquotInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AliquotRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AliquotInput{},
			Out: &AliquotOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type AliquotElement struct {
	inject.CheckedRunner
}

type AliquotInput struct {
	NumberofAliquots int
	OutPlate         *wtype.LHPlate
	OutPlatetype     string
	Solution         *wtype.LHComponent
	SolutionVolume   wunit.Volume
	VolumePerAliquot wunit.Volume
}

type AliquotOutput struct {
	Aliquots []*wtype.LHComponent
}

type AliquotSOutput struct {
	Data struct {
	}
	Outputs struct {
		Aliquots []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Aliquot",
		Constructor: AliquotNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/Aliquot/Aliquot.an",
			Params: []component.ParamDesc{
				{Name: "NumberofAliquots", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "OutPlatetype", Desc: "", Kind: "Parameters"},
				{Name: "Solution", Desc: "", Kind: "Inputs"},
				{Name: "SolutionVolume", Desc: "", Kind: "Parameters"},
				{Name: "VolumePerAliquot", Desc: "", Kind: "Parameters"},
				{Name: "Aliquots", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
