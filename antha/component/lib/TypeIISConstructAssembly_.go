package lib

import (
	"fmt"
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

func _TypeIISConstructAssemblyRequirements() {}

// Conditions to run on startup
func _TypeIISConstructAssemblySetup(_ctx context.Context, _input *TypeIISConstructAssemblyInput) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _TypeIISConstructAssemblySteps(_ctx context.Context, _input *TypeIISConstructAssemblyInput, _output *TypeIISConstructAssemblyOutput) {
	samples := make([]*wtype.LHComponent, 0)
	waterSample := mixer.SampleForTotalVolume(_input.Water, _input.ReactionVolume)
	samples = append(samples, waterSample)

	bufferSample := mixer.Sample(_input.Buffer, _input.BufferVol)
	samples = append(samples, bufferSample)

	atpSample := mixer.Sample(_input.Atp, _input.AtpVol)
	samples = append(samples, atpSample)

	//vectorSample := mixer.Sample(Vector, VectorVol)
	vectorSample := mixer.Sample(_input.Vector, _input.VectorVol)
	samples = append(samples, vectorSample)

	for k, part := range _input.Parts {
		fmt.Println("creating dna part num ", k, " comp ", part.CName, " renamed to ", _input.PartNames[k], " vol ", _input.PartVols[k])
		partSample := mixer.Sample(part, _input.PartVols[k])
		partSample.CName = _input.PartNames[k]
		samples = append(samples, partSample)
	}

	reSample := mixer.Sample(_input.RestrictionEnzyme, _input.ReVol)
	samples = append(samples, reSample)

	ligSample := mixer.Sample(_input.Ligase, _input.LigVol)
	samples = append(samples, ligSample)

	out1 := execute.MixInto(_ctx, _input.OutPlate, "", samples...)

	// incubate the reaction mixture
	out2 := execute.Incubate(_ctx, out1, _input.ReactionTemp, _input.ReactionTime, false)
	// inactivate
	_output.Reaction = execute.Incubate(_ctx, out2, _input.InactivationTemp, _input.InactivationTime, false)
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _TypeIISConstructAssemblyAnalysis(_ctx context.Context, _input *TypeIISConstructAssemblyInput, _output *TypeIISConstructAssemblyOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _TypeIISConstructAssemblyValidation(_ctx context.Context, _input *TypeIISConstructAssemblyInput, _output *TypeIISConstructAssemblyOutput) {
}
func _TypeIISConstructAssemblyRun(_ctx context.Context, input *TypeIISConstructAssemblyInput) *TypeIISConstructAssemblyOutput {
	output := &TypeIISConstructAssemblyOutput{}
	_TypeIISConstructAssemblySetup(_ctx, input)
	_TypeIISConstructAssemblySteps(_ctx, input, output)
	_TypeIISConstructAssemblyAnalysis(_ctx, input, output)
	_TypeIISConstructAssemblyValidation(_ctx, input, output)
	return output
}

func TypeIISConstructAssemblyRunSteps(_ctx context.Context, input *TypeIISConstructAssemblyInput) *TypeIISConstructAssemblySOutput {
	soutput := &TypeIISConstructAssemblySOutput{}
	output := _TypeIISConstructAssemblyRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func TypeIISConstructAssemblyNew() interface{} {
	return &TypeIISConstructAssemblyElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &TypeIISConstructAssemblyInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _TypeIISConstructAssemblyRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &TypeIISConstructAssemblyInput{},
			Out: &TypeIISConstructAssemblyOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type TypeIISConstructAssemblyElement struct {
	inject.CheckedRunner
}

type TypeIISConstructAssemblyInput struct {
	Atp                *wtype.LHComponent
	AtpVol             wunit.Volume
	Buffer             *wtype.LHComponent
	BufferVol          wunit.Volume
	InPlate            *wtype.LHPlate
	InactivationTemp   wunit.Temperature
	InactivationTime   wunit.Time
	LigVol             wunit.Volume
	Ligase             *wtype.LHComponent
	OutPlate           *wtype.LHPlate
	OutputReactionName string
	PartNames          []string
	PartVols           []wunit.Volume
	Parts              []*wtype.LHComponent
	ReVol              wunit.Volume
	ReactionTemp       wunit.Temperature
	ReactionTime       wunit.Time
	ReactionVolume     wunit.Volume
	RestrictionEnzyme  *wtype.LHComponent
	Vector             *wtype.LHComponent
	VectorVol          wunit.Volume
	Water              *wtype.LHComponent
}

type TypeIISConstructAssemblyOutput struct {
	Reaction *wtype.LHComponent
}

type TypeIISConstructAssemblySOutput struct {
	Data struct {
	}
	Outputs struct {
		Reaction *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "TypeIISConstructAssembly",
		Constructor: TypeIISConstructAssemblyNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/TypeIIsAssembly/TypeIISConstructAssembly/TypeIISConstructAssembly.an",
			Params: []component.ParamDesc{
				{Name: "Atp", Desc: "", Kind: "Inputs"},
				{Name: "AtpVol", Desc: "", Kind: "Parameters"},
				{Name: "Buffer", Desc: "", Kind: "Inputs"},
				{Name: "BufferVol", Desc: "", Kind: "Parameters"},
				{Name: "InPlate", Desc: "", Kind: "Inputs"},
				{Name: "InactivationTemp", Desc: "", Kind: "Parameters"},
				{Name: "InactivationTime", Desc: "", Kind: "Parameters"},
				{Name: "LigVol", Desc: "", Kind: "Parameters"},
				{Name: "Ligase", Desc: "", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "OutputReactionName", Desc: "", Kind: "Parameters"},
				{Name: "PartNames", Desc: "", Kind: "Parameters"},
				{Name: "PartVols", Desc: "", Kind: "Parameters"},
				{Name: "Parts", Desc: "", Kind: "Inputs"},
				{Name: "ReVol", Desc: "", Kind: "Parameters"},
				{Name: "ReactionTemp", Desc: "", Kind: "Parameters"},
				{Name: "ReactionTime", Desc: "", Kind: "Parameters"},
				{Name: "ReactionVolume", Desc: "", Kind: "Parameters"},
				{Name: "RestrictionEnzyme", Desc: "", Kind: "Inputs"},
				{Name: "Vector", Desc: "", Kind: "Inputs"},
				{Name: "VectorVol", Desc: "", Kind: "Parameters"},
				{Name: "Water", Desc: "", Kind: "Inputs"},
				{Name: "Reaction", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
