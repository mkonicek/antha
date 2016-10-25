// Example protocol demonstrating the use of the Sample function
package lib

import

// we need to import the wtype package to use the LHComponent type
// the mixer package is required to use the Sample function
(
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

// antha, like golang is a strongly typed language in which the type of a variable must be declared.
// In this case we're creating a variable called SampleVolume which is of type Volume;
// the type system allows the antha compiler to catch many types of common errors before the programme is run
// the antha type system extends this to biological types such as volumes here.
// functions require inputs of particular types to be adhered to

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// the LHComponent is the principal liquidhandling type in antha
// the * signifies that this is a pointer to the component rather than the component itself
// most key antha functions such as Sample and Mix use *LHComponent rather than LHComponent
// since the type is imported from the wtype package we need to use  *wtype.LHComponent rather than simply *LHComponent

// Physical outputs from this protocol with types

func _SampleRequirements() {

}

// Conditions to run on startup
func _SampleSetup(_ctx context.Context, _input *SampleInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _SampleSteps(_ctx context.Context, _input *SampleInput, _output *SampleOutput) {

	// the Sample function is imported from the mixer library
	// in the mixer library the function signature can be found, here it is:
	// func Sample(l *wtype.LHComponent, v wunit.Volume) *wtype.LHComponent {
	// The function signature  shows that the function requires a *LHComponent and a Volume and returns an *LHComponent
	_output.Sample = mixer.Sample(_input.Solution, _input.SampleVolume)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _SampleAnalysis(_ctx context.Context, _input *SampleInput, _output *SampleOutput) {
}

// A block of tests to perform to validate that the sample was processed
//correctly. Optionally, destructive tests can be performed to validate
//results on a dipstick basis
func _SampleValidation(_ctx context.Context, _input *SampleInput, _output *SampleOutput) {

}
func _SampleRun(_ctx context.Context, input *SampleInput) *SampleOutput {
	output := &SampleOutput{}
	_SampleSetup(_ctx, input)
	_SampleSteps(_ctx, input, output)
	_SampleAnalysis(_ctx, input, output)
	_SampleValidation(_ctx, input, output)
	return output
}

func SampleRunSteps(_ctx context.Context, input *SampleInput) *SampleSOutput {
	soutput := &SampleSOutput{}
	output := _SampleRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SampleNew() interface{} {
	return &SampleElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SampleInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SampleRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SampleInput{},
			Out: &SampleOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type SampleElement struct {
	inject.CheckedRunner
}

type SampleInput struct {
	SampleVolume wunit.Volume
	Solution     *wtype.LHComponent
}

type SampleOutput struct {
	Sample *wtype.LHComponent
}

type SampleSOutput struct {
	Data struct {
	}
	Outputs struct {
		Sample *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Sample",
		Constructor: SampleNew,
		Desc: component.ComponentDesc{
			Desc: "Example protocol demonstrating the use of the Sample function\n",
			Path: "antha/component/an/AnthaAcademy/Lesson1_RunningWorkflows/A_Sample.an",
			Params: []component.ParamDesc{
				{Name: "SampleVolume", Desc: "antha, like golang is a strongly typed language in which the type of a variable must be declared.\nIn this case we're creating a variable called SampleVolume which is of type Volume;\nthe type system allows the antha compiler to catch many types of common errors before the programme is run\nthe antha type system extends this to biological types such as volumes here.\nfunctions require inputs of particular types to be adhered to\n", Kind: "Parameters"},
				{Name: "Solution", Desc: "the LHComponent is the principal liquidhandling type in antha\nthe * signifies that this is a pointer to the component rather than the component itself\nmost key antha functions such as Sample and Mix use *LHComponent rather than LHComponent\nsince the type is imported from the wtype package we need to use  *wtype.LHComponent rather than simply *LHComponent\n", Kind: "Inputs"},
				{Name: "Sample", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
