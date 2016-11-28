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

func _TypeIISConstructAssemblyMMXInplateRequirements() {}

// Conditions to run on startup
func _TypeIISConstructAssemblyMMXInplateSetup(_ctx context.Context, _input *TypeIISConstructAssemblyMMXInplateInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _TypeIISConstructAssemblyMMXInplateSteps(_ctx context.Context, _input *TypeIISConstructAssemblyMMXInplateInput, _output *TypeIISConstructAssemblyMMXInplateOutput) {
	samples := make([]*wtype.LHComponent, 0)
	//	waterSample := mixer.SampleForTotalVolume(Water, ReactionVolume)
	mmxSample := mixer.SampleForTotalVolume(_input.MasterMix, _input.ReactionVolume)
	samples = append(samples, mmxSample)

	for k, part := range _input.Parts {
		//	fmt.Println("creating dna part num ", k, " comp ", part.CName, " renamed to ", PartNames[k], " vol ", PartVols[k])
		partSample := mixer.Sample(part, _input.PartVols[k])
		partSample.CName = _input.PartNames[k]
		samples = append(samples, partSample)
	}

	_output.Reaction = execute.MixTo(_ctx, _input.OutPlateType, _input.OutputLocation, _input.OutputPlateNum, samples...)

	// incubate the reaction mixture
	// commented out pending changes to incubate
	//Incubate(Reaction, ReactionTemp, ReactionTime, false)
	// inactivate
	//Incubate(Reaction, InactivationTemp, InactivationTime, false)
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _TypeIISConstructAssemblyMMXInplateAnalysis(_ctx context.Context, _input *TypeIISConstructAssemblyMMXInplateInput, _output *TypeIISConstructAssemblyMMXInplateOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _TypeIISConstructAssemblyMMXInplateValidation(_ctx context.Context, _input *TypeIISConstructAssemblyMMXInplateInput, _output *TypeIISConstructAssemblyMMXInplateOutput) {
}
func _TypeIISConstructAssemblyMMXInplateRun(_ctx context.Context, input *TypeIISConstructAssemblyMMXInplateInput) *TypeIISConstructAssemblyMMXInplateOutput {
	output := &TypeIISConstructAssemblyMMXInplateOutput{}
	_TypeIISConstructAssemblyMMXInplateSetup(_ctx, input)
	_TypeIISConstructAssemblyMMXInplateSteps(_ctx, input, output)
	_TypeIISConstructAssemblyMMXInplateAnalysis(_ctx, input, output)
	_TypeIISConstructAssemblyMMXInplateValidation(_ctx, input, output)
	return output
}

func TypeIISConstructAssemblyMMXInplateRunSteps(_ctx context.Context, input *TypeIISConstructAssemblyMMXInplateInput) *TypeIISConstructAssemblyMMXInplateSOutput {
	soutput := &TypeIISConstructAssemblyMMXInplateSOutput{}
	output := _TypeIISConstructAssemblyMMXInplateRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func TypeIISConstructAssemblyMMXInplateNew() interface{} {
	return &TypeIISConstructAssemblyMMXInplateElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &TypeIISConstructAssemblyMMXInplateInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _TypeIISConstructAssemblyMMXInplateRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &TypeIISConstructAssemblyMMXInplateInput{},
			Out: &TypeIISConstructAssemblyMMXInplateOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type TypeIISConstructAssemblyMMXInplateElement struct {
	inject.CheckedRunner
}

type TypeIISConstructAssemblyMMXInplateInput struct {
	InactivationTemp   wunit.Temperature
	InactivationTime   wunit.Time
	InputPlate         *wtype.LHPlate
	MasterMix          *wtype.LHComponent
	OutPlateType       string
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

type TypeIISConstructAssemblyMMXInplateOutput struct {
	Reaction *wtype.LHComponent
}

type TypeIISConstructAssemblyMMXInplateSOutput struct {
	Data struct {
	}
	Outputs struct {
		Reaction *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "TypeIISConstructAssemblyMMXInplate",
		Constructor: TypeIISConstructAssemblyMMXInplateNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/TypeIIsAssembly/TypeIISConstructAssemblyMMXInplate/TypeIISConstructAssemblyMMXInplate.an",
			Params: []component.ParamDesc{
				{Name: "InactivationTemp", Desc: "", Kind: "Parameters"},
				{Name: "InactivationTime", Desc: "", Kind: "Parameters"},
				{Name: "InputPlate", Desc: "", Kind: "Inputs"},
				{Name: "MasterMix", Desc: "", Kind: "Inputs"},
				{Name: "OutPlateType", Desc: "", Kind: "Inputs"},
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
