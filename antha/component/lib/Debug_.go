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

func _DebugRequirements() {
}

// Conditions to run on startup
func _DebugSetup(_ctx context.Context, _input *DebugInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _DebugSteps(_ctx context.Context, _input *DebugInput, _output *DebugOutput) {
	v := wunit.NewVolume(25, "ul")
	samples := make([]*wtype.LHComponent, 0, 1)

	samples = append(samples, mixer.Sample(_input.A, v))
	samples = append(samples, mixer.Sample(_input.B, v))
	samples = append(samples, mixer.Sample(_input.C, v))
	samples = append(samples, mixer.Sample(_input.D, v))
	samples = append(samples, mixer.Sample(_input.E, v))

	execute.Mix(_ctx, samples...)
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _DebugAnalysis(_ctx context.Context, _input *DebugInput, _output *DebugOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _DebugValidation(_ctx context.Context, _input *DebugInput, _output *DebugOutput) {
}
func _DebugRun(_ctx context.Context, input *DebugInput) *DebugOutput {
	output := &DebugOutput{}
	_DebugSetup(_ctx, input)
	_DebugSteps(_ctx, input, output)
	_DebugAnalysis(_ctx, input, output)
	_DebugValidation(_ctx, input, output)
	return output
}

func DebugRunSteps(_ctx context.Context, input *DebugInput) *DebugSOutput {
	soutput := &DebugSOutput{}
	output := _DebugRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func DebugNew() interface{} {
	return &DebugElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &DebugInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _DebugRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &DebugInput{},
			Out: &DebugOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type DebugElement struct {
	inject.CheckedRunner
}

type DebugInput struct {
	A *wtype.LHComponent
	B *wtype.LHComponent
	C *wtype.LHComponent
	D *wtype.LHComponent
	E *wtype.LHComponent
}

type DebugOutput struct {
}

type DebugSOutput struct {
	Data struct {
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Debug",
		Constructor: DebugNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Test/debug/debug.an",
			Params: []component.ParamDesc{
				{Name: "A", Desc: "", Kind: "Inputs"},
				{Name: "B", Desc: "", Kind: "Inputs"},
				{Name: "C", Desc: "", Kind: "Inputs"},
				{Name: "D", Desc: "", Kind: "Inputs"},
				{Name: "E", Desc: "", Kind: "Inputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
