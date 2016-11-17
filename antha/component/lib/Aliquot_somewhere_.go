// example protocol showing the highest level antha mix command which does not specify a plate type, therefore leaving it up to the scheduler to decide
package lib

import

// we can import code libraries and use functions and types from these libraries
(
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype" // the LHComponent type is imported from this library
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// the Sample function is imported from mixer

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _Aliquot_somewhereRequirements() {

}

// Conditions to run on startup
func _Aliquot_somewhereSetup(_ctx context.Context, _input *Aliquot_somewhereInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _Aliquot_somewhereSteps(_ctx context.Context, _input *Aliquot_somewhereInput, _output *Aliquot_somewhereOutput) {

	// First check that we can make enough aliquots of this volume
	number := _input.SolutionVolume.SIValue() / _input.VolumePerAliquot.SIValue()
	possiblenumberofAliquots, _ := wutil.RoundDown(number)
	if possiblenumberofAliquots < _input.NumberofAliquots {
		execute.Errorf(_ctx, "Not enough solution for this many aliquots")
	}

	// make a slice of components which we'll fill with aliquots;
	// same as we would for an array of samples but this time we won't mix together
	aliquots := make([]*wtype.LHComponent, 0)

	// this is golang syntax for a for loop
	// variable i is initialised at 0 and will increase with each loop whilst i < NumberofAliquots is still true
	for i := 0; i < _input.NumberofAliquots; i++ {

		// this is golang syntax for if statements
		// here we're checking if the liquid type is "dna" and if so we're changing the type
		// to ensure risk of cross contamination is completely avoided the dna liquid type does not allow multipipetting
		// in this case where we're just aliquoting the same dna into multiple destinations we can override this by changing the liquid type
		if _input.Solution.TypeName() == "dna" {
			_input.Solution.Type = wtype.LTDoNotMix
		}
		aliquotSample := mixer.Sample(_input.Solution, _input.VolumePerAliquot)
		aliquot := execute.Mix(_ctx, aliquotSample)

		// this time we append the slice of components after mixing
		aliquots = append(aliquots, aliquot)
	}

	// Now we assign our temporary variable aliqouts to export as a variable as specified in Outputs
	// In Antha the first letter of a variablename must be uppercase to allow the variable to be exported
	_output.Aliquots = aliquots
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Aliquot_somewhereAnalysis(_ctx context.Context, _input *Aliquot_somewhereInput, _output *Aliquot_somewhereOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _Aliquot_somewhereValidation(_ctx context.Context, _input *Aliquot_somewhereInput, _output *Aliquot_somewhereOutput) {

}
func _Aliquot_somewhereRun(_ctx context.Context, input *Aliquot_somewhereInput) *Aliquot_somewhereOutput {
	output := &Aliquot_somewhereOutput{}
	_Aliquot_somewhereSetup(_ctx, input)
	_Aliquot_somewhereSteps(_ctx, input, output)
	_Aliquot_somewhereAnalysis(_ctx, input, output)
	_Aliquot_somewhereValidation(_ctx, input, output)
	return output
}

func Aliquot_somewhereRunSteps(_ctx context.Context, input *Aliquot_somewhereInput) *Aliquot_somewhereSOutput {
	soutput := &Aliquot_somewhereSOutput{}
	output := _Aliquot_somewhereRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Aliquot_somewhereNew() interface{} {
	return &Aliquot_somewhereElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Aliquot_somewhereInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Aliquot_somewhereRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Aliquot_somewhereInput{},
			Out: &Aliquot_somewhereOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Aliquot_somewhereElement struct {
	inject.CheckedRunner
}

type Aliquot_somewhereInput struct {
	NumberofAliquots int
	Solution         *wtype.LHComponent
	SolutionVolume   wunit.Volume
	VolumePerAliquot wunit.Volume
}

type Aliquot_somewhereOutput struct {
	Aliquots []*wtype.LHComponent
}

type Aliquot_somewhereSOutput struct {
	Data struct {
	}
	Outputs struct {
		Aliquots []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Aliquot_somewhere",
		Constructor: Aliquot_somewhereNew,
		Desc: component.ComponentDesc{
			Desc: "example protocol showing the highest level antha mix command which does not specify a plate type, therefore leaving it up to the scheduler to decide\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/AnthaAcademy/Lesson2_mix/A_Aliquot_somewhereorother.an",
			Params: []component.ParamDesc{
				{Name: "NumberofAliquots", Desc: "", Kind: "Parameters"},
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
