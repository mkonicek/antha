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

func _SequentialMixTestProtocolRequirements() {}

// Conditions to run on startup
func _SequentialMixTestProtocolSetup(_ctx context.Context, _input *SequentialMixTestProtocolInput) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _SequentialMixTestProtocolSteps(_ctx context.Context, _input *SequentialMixTestProtocolInput, _output *SequentialMixTestProtocolOutput) {
	mmxSample := mixer.SampleForTotalVolume(_input.MasterMix, _input.ReactionVolume)
	mixture := execute.MixTo(_ctx, _input.OutPlate.Type, _input.OutputLocation, _input.OutputPlateNum, mmxSample)

	for k, part := range _input.Parts {
		partSample := mixer.Sample(part, _input.PartVols[k])
		partSample.CName = _input.PartNames[k]
		mixture = execute.Mix(_ctx, mixture, partSample) // mix in place
	}

	_output.Reaction = mixture

	//Reaction = MixTo(OutPlate.Type, OutputLocation, OutputPlateNum, samples...)

	// incubate the reaction mixture
	// commented out pending changes to incubate
	//Incubate(Reaction, ReactionTemp, ReactionTime, false)
	// inactivate
	//Incubate(Reaction, InactivationTemp, InactivationTime, false)
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _SequentialMixTestProtocolAnalysis(_ctx context.Context, _input *SequentialMixTestProtocolInput, _output *SequentialMixTestProtocolOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _SequentialMixTestProtocolValidation(_ctx context.Context, _input *SequentialMixTestProtocolInput, _output *SequentialMixTestProtocolOutput) {
}
func _SequentialMixTestProtocolRun(_ctx context.Context, input *SequentialMixTestProtocolInput) *SequentialMixTestProtocolOutput {
	output := &SequentialMixTestProtocolOutput{}
	_SequentialMixTestProtocolSetup(_ctx, input)
	_SequentialMixTestProtocolSteps(_ctx, input, output)
	_SequentialMixTestProtocolAnalysis(_ctx, input, output)
	_SequentialMixTestProtocolValidation(_ctx, input, output)
	return output
}

func SequentialMixTestProtocolRunSteps(_ctx context.Context, input *SequentialMixTestProtocolInput) *SequentialMixTestProtocolSOutput {
	soutput := &SequentialMixTestProtocolSOutput{}
	output := _SequentialMixTestProtocolRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SequentialMixTestProtocolNew() interface{} {
	return &SequentialMixTestProtocolElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SequentialMixTestProtocolInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SequentialMixTestProtocolRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SequentialMixTestProtocolInput{},
			Out: &SequentialMixTestProtocolOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type SequentialMixTestProtocolElement struct {
	inject.CheckedRunner
}

type SequentialMixTestProtocolInput struct {
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

type SequentialMixTestProtocolOutput struct {
	Reaction *wtype.LHComponent
}

type SequentialMixTestProtocolSOutput struct {
	Data struct {
	}
	Outputs struct {
		Reaction *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "SequentialMixTestProtocol",
		Constructor: SequentialMixTestProtocolNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Test/SequentialMixTestProtocol/SequentialMixTestProtocol.an",
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
