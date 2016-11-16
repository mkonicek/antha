// example protocol showing The MixTo command which allows a specifc plate type to be specified. i.e. platetype pcrplate_skirted
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

func _Aliquot_PlateTypeRequirements() {

}

// Conditions to run on startup
func _Aliquot_PlateTypeSetup(_ctx context.Context, _input *Aliquot_PlateTypeInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _Aliquot_PlateTypeSteps(_ctx context.Context, _input *Aliquot_PlateTypeInput, _output *Aliquot_PlateTypeOutput) {

	number := _input.SolutionVolume.SIValue() / _input.VolumePerAliquot.SIValue()
	possiblenumberofAliquots, _ := wutil.RoundDown(number)
	if possiblenumberofAliquots < _input.NumberofAliquots {
		execute.Errorf(_ctx, "Not enough solution for this many aliquots")
	}

	aliquots := make([]*wtype.LHComponent, 0)

	for i := 0; i < _input.NumberofAliquots; i++ {
		if _input.Solution.TypeName() == "dna" {
			_input.Solution.Type = wtype.LTDoNotMix
		}
		aliquotSample := mixer.Sample(_input.Solution, _input.VolumePerAliquot)
		// the MixTo command is used instead of Mix to specify the plate type (e.g. "greiner384" or "pcrplate_skirted")
		// the plate types can be found in antha-lang/antha/microArch/factory/make_plate_library.go
		// the four input fields to the MixTo command represent
		// 1. the platetype as a string: commonly the input to the antha element will actually be an LHPlate rather than a string so the type field can be accessed with OutPlate.Type
		// 2. well location as a  string e.g. "A1" (in this case leaving it blank "" will leave the well location up to the scheduler),
		// 3. the plate number,starting from 1 (not zero)
		// 4. the sample or array of samples to be mixed; in the case of an array you'd normally feed this in as samples...
		aliquot := execute.MixTo(_ctx, _input.OutPlatetype, "", 1, aliquotSample)
		aliquots = append(aliquots, aliquot)
	}
	_output.Aliquots = aliquots //
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Aliquot_PlateTypeAnalysis(_ctx context.Context, _input *Aliquot_PlateTypeInput, _output *Aliquot_PlateTypeOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _Aliquot_PlateTypeValidation(_ctx context.Context, _input *Aliquot_PlateTypeInput, _output *Aliquot_PlateTypeOutput) {

}
func _Aliquot_PlateTypeRun(_ctx context.Context, input *Aliquot_PlateTypeInput) *Aliquot_PlateTypeOutput {
	output := &Aliquot_PlateTypeOutput{}
	_Aliquot_PlateTypeSetup(_ctx, input)
	_Aliquot_PlateTypeSteps(_ctx, input, output)
	_Aliquot_PlateTypeAnalysis(_ctx, input, output)
	_Aliquot_PlateTypeValidation(_ctx, input, output)
	return output
}

func Aliquot_PlateTypeRunSteps(_ctx context.Context, input *Aliquot_PlateTypeInput) *Aliquot_PlateTypeSOutput {
	soutput := &Aliquot_PlateTypeSOutput{}
	output := _Aliquot_PlateTypeRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Aliquot_PlateTypeNew() interface{} {
	return &Aliquot_PlateTypeElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Aliquot_PlateTypeInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Aliquot_PlateTypeRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Aliquot_PlateTypeInput{},
			Out: &Aliquot_PlateTypeOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Aliquot_PlateTypeElement struct {
	inject.CheckedRunner
}

type Aliquot_PlateTypeInput struct {
	NumberofAliquots int
	OutPlatetype     string
	Solution         *wtype.LHComponent
	SolutionVolume   wunit.Volume
	VolumePerAliquot wunit.Volume
}

type Aliquot_PlateTypeOutput struct {
	Aliquots []*wtype.LHComponent
}

type Aliquot_PlateTypeSOutput struct {
	Data struct {
	}
	Outputs struct {
		Aliquots []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Aliquot_PlateType",
		Constructor: Aliquot_PlateTypeNew,
		Desc: component.ComponentDesc{
			Desc: "example protocol showing The MixTo command which allows a specifc plate type to be specified. i.e. platetype pcrplate_skirted\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/AnthaAcademy/Lesson2_mix/C_AliquotTo_PlateType.an",
			Params: []component.ParamDesc{
				{Name: "NumberofAliquots", Desc: "", Kind: "Parameters"},
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
