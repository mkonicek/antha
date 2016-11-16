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

func _TypeIISConstructAssemblyMMXRequirements() {}

// Conditions to run on startup
func _TypeIISConstructAssemblyMMXSetup(_ctx context.Context, _input *TypeIISConstructAssemblyMMXInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _TypeIISConstructAssemblyMMXSteps(_ctx context.Context, _input *TypeIISConstructAssemblyMMXInput, _output *TypeIISConstructAssemblyMMXOutput) {
	samples := make([]*wtype.LHComponent, 0)

	//waterSample:=mixer.SampleForTotalVolume(Water,ReactionVolume)
	//	samples = append(samples, waterSample)

	mmxSample := mixer.SampleForTotalVolume(_input.MasterMix, _input.ReactionVolume)
	samples = append(samples, mmxSample)

	for k, part := range _input.Parts {
		part.Type, _ = wtype.LiquidTypeFromString(_input.LHPolicyName)

		partSample := mixer.Sample(part, _input.PartVols[k])
		partSample.CName = _input.PartNames[k]
		samples = append(samples, partSample)
	}

	// ensure the last step is mixed
	samples[len(samples)-1].Type = wtype.LTDNAMIX
	_output.Reaction = execute.MixTo(_ctx, _input.OutPlate.Type, _input.OutputLocation, _input.OutputPlateNum, samples...)

	// incubate the reaction mixture
	// commented out pending changes to incubate
	execute.Incubate(_ctx, _output.Reaction, _input.ReactionTemp, _input.ReactionTime, false)
	// inactivate
	//Incubate(Reaction, InactivationTemp, InactivationTime, false)
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _TypeIISConstructAssemblyMMXAnalysis(_ctx context.Context, _input *TypeIISConstructAssemblyMMXInput, _output *TypeIISConstructAssemblyMMXOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _TypeIISConstructAssemblyMMXValidation(_ctx context.Context, _input *TypeIISConstructAssemblyMMXInput, _output *TypeIISConstructAssemblyMMXOutput) {
}
func _TypeIISConstructAssemblyMMXRun(_ctx context.Context, input *TypeIISConstructAssemblyMMXInput) *TypeIISConstructAssemblyMMXOutput {
	output := &TypeIISConstructAssemblyMMXOutput{}
	_TypeIISConstructAssemblyMMXSetup(_ctx, input)
	_TypeIISConstructAssemblyMMXSteps(_ctx, input, output)
	_TypeIISConstructAssemblyMMXAnalysis(_ctx, input, output)
	_TypeIISConstructAssemblyMMXValidation(_ctx, input, output)
	return output
}

func TypeIISConstructAssemblyMMXRunSteps(_ctx context.Context, input *TypeIISConstructAssemblyMMXInput) *TypeIISConstructAssemblyMMXSOutput {
	soutput := &TypeIISConstructAssemblyMMXSOutput{}
	output := _TypeIISConstructAssemblyMMXRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func TypeIISConstructAssemblyMMXNew() interface{} {
	return &TypeIISConstructAssemblyMMXElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &TypeIISConstructAssemblyMMXInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _TypeIISConstructAssemblyMMXRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &TypeIISConstructAssemblyMMXInput{},
			Out: &TypeIISConstructAssemblyMMXOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type TypeIISConstructAssemblyMMXElement struct {
	inject.CheckedRunner
}

type TypeIISConstructAssemblyMMXInput struct {
	InactivationTemp   wunit.Temperature
	InactivationTime   wunit.Time
	LHPolicyName       string
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

type TypeIISConstructAssemblyMMXOutput struct {
	Reaction *wtype.LHComponent
}

type TypeIISConstructAssemblyMMXSOutput struct {
	Data struct {
	}
	Outputs struct {
		Reaction *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "TypeIISConstructAssemblyMMX",
		Constructor: TypeIISConstructAssemblyMMXNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/TypeIIsAssembly/TypeIISConstructAssemblyMMX/TypeIISConstructAssemblyMMX.an",
			Params: []component.ParamDesc{
				{Name: "InactivationTemp", Desc: "", Kind: "Parameters"},
				{Name: "InactivationTime", Desc: "", Kind: "Parameters"},
				{Name: "LHPolicyName", Desc: "", Kind: "Parameters"},
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
