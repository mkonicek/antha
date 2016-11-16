// The lowest level example protocol showing The MixTo command being used to specify the specific wells to be aliquoted to;
// By doing this we are able to specify whether the aliqouts are pipetted by row or by column.
// In this case the user is still not specifying the well location (i.e. A1) in the parameters, although that would be possible to specify.
// We don't generally encourage this since Antha is designed to be prodiminantly a high level language which avoids the user specifying well locations but this possibility is there if necessary.
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
	"strconv"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// we're now going to aliquot multiple solutions at the same time (but not mixing them)

// Physical outputs from this protocol with types

func _AliquotToRequirements() {

}

// Conditions to run on startup
func _AliquotToSetup(_ctx context.Context, _input *AliquotToInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _AliquotToSteps(_ctx context.Context, _input *AliquotToInput, _output *AliquotToOutput) {

	number := _input.SolutionVolume.SIValue() / _input.VolumePerAliquot.SIValue()
	possiblenumberofAliquots, _ := wutil.RoundDown(number)
	if possiblenumberofAliquots < _input.NumberofAliquots {
		execute.Errorf(_ctx, "Not enough solution for this many aliquots")
	}

	aliquots := make([]*wtype.LHComponent, 0)

	// work out well coordinates for any plate
	wellpositionarray := make([]string, 0)

	alphabet := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J",
		"K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X",
		"Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF"}

	if _input.ByRow {
		// add well positions to the array based upon the number of wells per column (OutPlate.WlsX) and row (OutPlate.WlsY) of the plate type in question
		for j := 0; j < _input.OutPlate.WlsY; j++ {
			for i := 0; i < _input.OutPlate.WlsX; i++ {

				// antha, like golang upon which it is built, is a strongly type language so an int must be converted to a string using the strconv package
				// as shown here, strings can be concatenated using +
				// other types can sometimes be converted more directly.
				// In particular an int can be converted to a float64 like this:
				// var myInt int = 1
				// var myFloat float64
				// myFloat = float64(myInt)
				wellposition := alphabet[j] + strconv.Itoa(i+1)

				wellpositionarray = append(wellpositionarray, wellposition)
			}

		}
	} else {
		for j := 0; j < _input.OutPlate.WlsX; j++ {
			for i := 0; i < _input.OutPlate.WlsY; i++ {

				wellposition := alphabet[i] + strconv.Itoa(j+1)

				wellpositionarray = append(wellpositionarray, wellposition)
			}

		}
	}

	// initialise a counter
	var counter int // an int is initialised as zero therefore this is the same as counter := 0 or var counter = 0

	for _, Solution := range _input.Solutions {
		for k := 0; k < _input.NumberofAliquots; k++ {

			if Solution.TypeName() == "dna" {
				Solution.Type = wtype.LTDoNotMix
			}
			aliquotSample := mixer.Sample(Solution, _input.VolumePerAliquot)

			// this time we're using counter as an index to go through the wellpositionarray one position at a time and ensuring the next free position is chosen
			// the platenumber is hardcoded to 1 here so if we tried to specify too many aliquots in the parameters the protocol would fail
			// it would be better to create a platenumber variable of type int and use an if statement to increase platenumber by 1 if all well positions are filled up i.e.
			// if counter == len(wellpositionarray) {
			// 		platenumber++
			//}
			aliquot := execute.MixTo(_ctx, _input.OutPlate.Type, wellpositionarray[counter], 1, aliquotSample)
			aliquots = append(aliquots, aliquot)
			counter = counter + 1 // this is the same as using the more concise counter++
		}
		_output.Aliquots = aliquots

		// Exercise: refactor to use wtype.WellCoords instead of creating the well ids manually using alphabet and strconv
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AliquotToAnalysis(_ctx context.Context, _input *AliquotToInput, _output *AliquotToOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AliquotToValidation(_ctx context.Context, _input *AliquotToInput, _output *AliquotToOutput) {

}
func _AliquotToRun(_ctx context.Context, input *AliquotToInput) *AliquotToOutput {
	output := &AliquotToOutput{}
	_AliquotToSetup(_ctx, input)
	_AliquotToSteps(_ctx, input, output)
	_AliquotToAnalysis(_ctx, input, output)
	_AliquotToValidation(_ctx, input, output)
	return output
}

func AliquotToRunSteps(_ctx context.Context, input *AliquotToInput) *AliquotToSOutput {
	soutput := &AliquotToSOutput{}
	output := _AliquotToRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AliquotToNew() interface{} {
	return &AliquotToElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AliquotToInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AliquotToRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AliquotToInput{},
			Out: &AliquotToOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type AliquotToElement struct {
	inject.CheckedRunner
}

type AliquotToInput struct {
	ByRow            bool
	NumberofAliquots int
	OutPlate         *wtype.LHPlate
	SolutionVolume   wunit.Volume
	Solutions        []*wtype.LHComponent
	VolumePerAliquot wunit.Volume
}

type AliquotToOutput struct {
	Aliquots []*wtype.LHComponent
}

type AliquotToSOutput struct {
	Data struct {
	}
	Outputs struct {
		Aliquots []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AliquotTo",
		Constructor: AliquotToNew,
		Desc: component.ComponentDesc{
			Desc: "The lowest level example protocol showing The MixTo command being used to specify the specific wells to be aliquoted to;\nBy doing this we are able to specify whether the aliqouts are pipetted by row or by column.\nIn this case the user is still not specifying the well location (i.e. A1) in the parameters, although that would be possible to specify.\nWe don't generally encourage this since Antha is designed to be prodiminantly a high level language which avoids the user specifying well locations but this possibility is there if necessary.\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/AnthaAcademy/Lesson2_mix/E_AliquotTo_wellpositions.an",
			Params: []component.ParamDesc{
				{Name: "ByRow", Desc: "", Kind: "Parameters"},
				{Name: "NumberofAliquots", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "SolutionVolume", Desc: "", Kind: "Parameters"},
				{Name: "Solutions", Desc: "we're now going to aliquot multiple solutions at the same time (but not mixing them)\n", Kind: "Inputs"},
				{Name: "VolumePerAliquot", Desc: "", Kind: "Parameters"},
				{Name: "Aliquots", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
