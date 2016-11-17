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

func _PartiallySequentialMixTestProtocolRequirements() {}

// Conditions to run on startup
func _PartiallySequentialMixTestProtocolSetup(_ctx context.Context, _input *PartiallySequentialMixTestProtocolInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _PartiallySequentialMixTestProtocolSteps(_ctx context.Context, _input *PartiallySequentialMixTestProtocolInput, _output *PartiallySequentialMixTestProtocolOutput) {
	mmxSample := mixer.SampleForTotalVolume(_input.MasterMix, _input.ReactionVolume)
	mixture := execute.MixTo(_ctx, _input.OutPlate.Type, _input.OutputLocation, _input.OutputPlateNum, mmxSample)

	samples := make([]*wtype.LHComponent, 0, 1)
	samples = append(samples, mixture)

	for k, part := range _input.Parts {
		partSample := mixer.Sample(part, _input.PartVols[k])
		partSample.CName = _input.PartNames[k]
		samples = append(samples, partSample)
	}

	_output.Reaction = execute.Mix(_ctx, samples...)

	// incubate the reaction mixture
	// commented out pending changes to incubate
	//Incubate(Reaction, ReactionTemp, ReactionTime, false)
	// inactivate
	//Incubate(Reaction, InactivationTemp, InactivationTime, false)
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PartiallySequentialMixTestProtocolAnalysis(_ctx context.Context, _input *PartiallySequentialMixTestProtocolInput, _output *PartiallySequentialMixTestProtocolOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PartiallySequentialMixTestProtocolValidation(_ctx context.Context, _input *PartiallySequentialMixTestProtocolInput, _output *PartiallySequentialMixTestProtocolOutput) {
}
func _PartiallySequentialMixTestProtocolRun(_ctx context.Context, input *PartiallySequentialMixTestProtocolInput) *PartiallySequentialMixTestProtocolOutput {
	output := &PartiallySequentialMixTestProtocolOutput{}
	_PartiallySequentialMixTestProtocolSetup(_ctx, input)
	_PartiallySequentialMixTestProtocolSteps(_ctx, input, output)
	_PartiallySequentialMixTestProtocolAnalysis(_ctx, input, output)
	_PartiallySequentialMixTestProtocolValidation(_ctx, input, output)
	return output
}

func PartiallySequentialMixTestProtocolRunSteps(_ctx context.Context, input *PartiallySequentialMixTestProtocolInput) *PartiallySequentialMixTestProtocolSOutput {
	soutput := &PartiallySequentialMixTestProtocolSOutput{}
	output := _PartiallySequentialMixTestProtocolRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PartiallySequentialMixTestProtocolNew() interface{} {
	return &PartiallySequentialMixTestProtocolElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PartiallySequentialMixTestProtocolInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PartiallySequentialMixTestProtocolRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PartiallySequentialMixTestProtocolInput{},
			Out: &PartiallySequentialMixTestProtocolOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PartiallySequentialMixTestProtocolElement struct {
	inject.CheckedRunner
}

type PartiallySequentialMixTestProtocolInput struct {
	InactivationTemp   wunit.Temperature
	InactivationTime   wunit.Time
	MasterMix          *wtype.LHComponent
	OutPlate           *wtype.LHPlate
	OutputLocation     string
	OutputPlateNum     int
	OutputReactionName string
	PartNames          []string
	PartVols           []wunit.Volume
	Parts              []*wtype.LHComponent
	ReactionTemp       wunit.Temperature
	ReactionTime       wunit.Time
	ReactionVolume     wunit.Volume
}

type PartiallySequentialMixTestProtocolOutput struct {
	Reaction *wtype.LHComponent
}

type PartiallySequentialMixTestProtocolSOutput struct {
	Data struct {
	}
	Outputs struct {
		Reaction *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PartiallySequentialMixTestProtocol",
		Constructor: PartiallySequentialMixTestProtocolNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Test/PartiallySequentialMixTestProtocol/PartiallySequentialMixTestProtocol.an",
			Params: []component.ParamDesc{
				{Name: "InactivationTemp", Desc: "", Kind: "Parameters"},
				{Name: "InactivationTime", Desc: "", Kind: "Parameters"},
				{Name: "MasterMix", Desc: "", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "OutputLocation", Desc: "", Kind: "Parameters"},
				{Name: "OutputPlateNum", Desc: "", Kind: "Parameters"},
				{Name: "OutputReactionName", Desc: "", Kind: "Parameters"},
				{Name: "PartNames", Desc: "", Kind: "Parameters"},
				{Name: "PartVols", Desc: "", Kind: "Parameters"},
				{Name: "Parts", Desc: "", Kind: "Inputs"},
				{Name: "ReactionTemp", Desc: "", Kind: "Parameters"},
				{Name: "ReactionTime", Desc: "", Kind: "Parameters"},
				{Name: "ReactionVolume", Desc: "", Kind: "Parameters"},
				{Name: "Reaction", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
