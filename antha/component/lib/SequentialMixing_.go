// Protocol to Demonstrate how to perform sequential mixing using the example of
// making a serial dilution series from a solution and diluent
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

// e.g. 10 would take 1 part solution to 9 parts diluent for each dilution

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _SequentialMixingRequirements() {

}

// Conditions to run on startup
func _SequentialMixingSetup(_ctx context.Context, _input *SequentialMixingInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _SequentialMixingSteps(_ctx context.Context, _input *SequentialMixingInput, _output *SequentialMixingOutput) {

	dilutions := make([]*wtype.LHComponent, 0)

	var aliquot *wtype.LHComponent

	// calculate solution volume

	// create copy of TotalVolumeperDilution
	solutionVolume := (wunit.CopyVolume(_input.TotalVolumeperDilution))

	// use divideby method
	solutionVolume.DivideBy(float64(_input.DilutionFactor))

	// use same approach to work out diluent volume to add
	diluentVolume := (wunit.CopyVolume(_input.TotalVolumeperDilution))

	// this time using the substract method
	diluentVolume.Subtract(solutionVolume)

	// sample diluent
	diluentSample := mixer.Sample(_input.Diluent, diluentVolume)

	// Ensure liquid type set to Pre and Post Mix
	_input.Solution.Type = wtype.LTNeedToMix

	// sample solution
	solutionSample := mixer.Sample(_input.Solution, solutionVolume)

	// mix both samples to OutPlate
	aliquot = execute.MixTo(_ctx, _input.OutPlate.Type, "", 1, diluentSample, solutionSample)

	// add to dilutions array
	dilutions = append(dilutions, aliquot)

	// loop through NumberOfDilutions until all serial dilutions are made
	for k := 1; k < _input.NumberOfDilutions; k++ {

		// take next sample of diluent
		nextdiluentSample := mixer.Sample(_input.Diluent, diluentVolume)

		// Ensure liquid type set to Pre and Post Mix
		aliquot.Type = wtype.LTNeedToMix

		// sample from previous dilution sample
		nextSample := mixer.Sample(aliquot, solutionVolume)

		// Mix sample into nextdiluent sample
		nextaliquot := execute.Mix(_ctx, nextdiluentSample, nextSample)

		// add to dilutions array
		dilutions = append(dilutions, nextaliquot)
		// reset aliquot
		aliquot = nextaliquot
	}

	// export as Output
	_output.Dilutions = dilutions

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _SequentialMixingAnalysis(_ctx context.Context, _input *SequentialMixingInput, _output *SequentialMixingOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _SequentialMixingValidation(_ctx context.Context, _input *SequentialMixingInput, _output *SequentialMixingOutput) {

}
func _SequentialMixingRun(_ctx context.Context, input *SequentialMixingInput) *SequentialMixingOutput {
	output := &SequentialMixingOutput{}
	_SequentialMixingSetup(_ctx, input)
	_SequentialMixingSteps(_ctx, input, output)
	_SequentialMixingAnalysis(_ctx, input, output)
	_SequentialMixingValidation(_ctx, input, output)
	return output
}

func SequentialMixingRunSteps(_ctx context.Context, input *SequentialMixingInput) *SequentialMixingSOutput {
	soutput := &SequentialMixingSOutput{}
	output := _SequentialMixingRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SequentialMixingNew() interface{} {
	return &SequentialMixingElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SequentialMixingInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SequentialMixingRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SequentialMixingInput{},
			Out: &SequentialMixingOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type SequentialMixingElement struct {
	inject.CheckedRunner
}

type SequentialMixingInput struct {
	Diluent                *wtype.LHComponent
	DilutionFactor         int
	NumberOfDilutions      int
	OutPlate               *wtype.LHPlate
	Solution               *wtype.LHComponent
	TotalVolumeperDilution wunit.Volume
}

type SequentialMixingOutput struct {
	Dilutions []*wtype.LHComponent
}

type SequentialMixingSOutput struct {
	Data struct {
	}
	Outputs struct {
		Dilutions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "SequentialMixing",
		Constructor: SequentialMixingNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol to Demonstrate how to perform sequential mixing using the example of\nmaking a serial dilution series from a solution and diluent\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/AnthaAcademy/Lesson3_MixPart2/C_SequentialMixing.an",
			Params: []component.ParamDesc{
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "DilutionFactor", Desc: "e.g. 10 would take 1 part solution to 9 parts diluent for each dilution\n", Kind: "Parameters"},
				{Name: "NumberOfDilutions", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Solution", Desc: "", Kind: "Inputs"},
				{Name: "TotalVolumeperDilution", Desc: "", Kind: "Parameters"},
				{Name: "Dilutions", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
