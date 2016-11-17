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

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

func _AnotherSequentialMixTestProtocolRequirements() {}

// Conditions to run on startup
func _AnotherSequentialMixTestProtocolSetup(_ctx context.Context, _input *AnotherSequentialMixTestProtocolInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _AnotherSequentialMixTestProtocolSteps(_ctx context.Context, _input *AnotherSequentialMixTestProtocolInput, _output *AnotherSequentialMixTestProtocolOutput) {
	var m *wtype.LHComponent

	for k := 0; k < len(_input.Vols); k++ {
		var s *wtype.LHComponent
		if !_input.Vols[k].IsZero() {
			s = mixer.Sample(_input.Parts[k], _input.Vols[k])
			if m != nil {
				m = execute.Mix(_ctx, m, s)
			} else {
				m = execute.MixTo(_ctx, _input.OutputPlateType, "", 1, s)
			}
		}
	}

	_output.Reaction = m
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AnotherSequentialMixTestProtocolAnalysis(_ctx context.Context, _input *AnotherSequentialMixTestProtocolInput, _output *AnotherSequentialMixTestProtocolOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AnotherSequentialMixTestProtocolValidation(_ctx context.Context, _input *AnotherSequentialMixTestProtocolInput, _output *AnotherSequentialMixTestProtocolOutput) {
}
func _AnotherSequentialMixTestProtocolRun(_ctx context.Context, input *AnotherSequentialMixTestProtocolInput) *AnotherSequentialMixTestProtocolOutput {
	output := &AnotherSequentialMixTestProtocolOutput{}
	_AnotherSequentialMixTestProtocolSetup(_ctx, input)
	_AnotherSequentialMixTestProtocolSteps(_ctx, input, output)
	_AnotherSequentialMixTestProtocolAnalysis(_ctx, input, output)
	_AnotherSequentialMixTestProtocolValidation(_ctx, input, output)
	return output
}

func AnotherSequentialMixTestProtocolRunSteps(_ctx context.Context, input *AnotherSequentialMixTestProtocolInput) *AnotherSequentialMixTestProtocolSOutput {
	soutput := &AnotherSequentialMixTestProtocolSOutput{}
	output := _AnotherSequentialMixTestProtocolRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AnotherSequentialMixTestProtocolNew() interface{} {
	return &AnotherSequentialMixTestProtocolElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AnotherSequentialMixTestProtocolInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AnotherSequentialMixTestProtocolRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AnotherSequentialMixTestProtocolInput{},
			Out: &AnotherSequentialMixTestProtocolOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type AnotherSequentialMixTestProtocolElement struct {
	inject.CheckedRunner
}

type AnotherSequentialMixTestProtocolInput struct {
	OutputPlateType string
	Parts           []*wtype.LHComponent
	Vols            []wunit.Volume
}

type AnotherSequentialMixTestProtocolOutput struct {
	Reaction *wtype.LHComponent
}

type AnotherSequentialMixTestProtocolSOutput struct {
	Data struct {
	}
	Outputs struct {
		Reaction *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AnotherSequentialMixTestProtocol",
		Constructor: AnotherSequentialMixTestProtocolNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Test/AnotherSequentialMixTestProtocol/AnotherSequentialMixTestProtocol.an",
			Params: []component.ParamDesc{
				{Name: "OutputPlateType", Desc: "", Kind: "Inputs"},
				{Name: "Parts", Desc: "", Kind: "Inputs"},
				{Name: "Vols", Desc: "", Kind: "Parameters"},
				{Name: "Reaction", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
