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

//	"fmt"

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

func _Ligation_2partsintovectorRequirements() {}

// Conditions to run on startup
func _Ligation_2partsintovectorSetup(_ctx context.Context, _input *Ligation_2partsintovectorInput) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _Ligation_2partsintovectorSteps(_ctx context.Context, _input *Ligation_2partsintovectorInput, _output *Ligation_2partsintovectorOutput) {
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

	// part 1
	//fmt.Println("creating dna part 1 ", " comp ", Part1.CName, " renamed to ", PartNames[0], " vol ", PartVols[0])
	partSample := mixer.Sample(_input.Part1, _input.PartVols[0])
	//partSample.CName = PartNames[0]
	samples = append(samples, partSample)

	// part 2
	//fmt.Println("creating dna part 2 ", " comp ", Part2.CName, " renamed to ", PartNames[1], " vol ", PartVols[1])
	partSample = mixer.Sample(_input.Part2, _input.PartVols[1])
	//partSample.CName = PartNames[1]
	samples = append(samples, partSample)

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
func _Ligation_2partsintovectorAnalysis(_ctx context.Context, _input *Ligation_2partsintovectorInput, _output *Ligation_2partsintovectorOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _Ligation_2partsintovectorValidation(_ctx context.Context, _input *Ligation_2partsintovectorInput, _output *Ligation_2partsintovectorOutput) {
}
func _Ligation_2partsintovectorRun(_ctx context.Context, input *Ligation_2partsintovectorInput) *Ligation_2partsintovectorOutput {
	output := &Ligation_2partsintovectorOutput{}
	_Ligation_2partsintovectorSetup(_ctx, input)
	_Ligation_2partsintovectorSteps(_ctx, input, output)
	_Ligation_2partsintovectorAnalysis(_ctx, input, output)
	_Ligation_2partsintovectorValidation(_ctx, input, output)
	return output
}

func Ligation_2partsintovectorRunSteps(_ctx context.Context, input *Ligation_2partsintovectorInput) *Ligation_2partsintovectorSOutput {
	soutput := &Ligation_2partsintovectorSOutput{}
	output := _Ligation_2partsintovectorRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Ligation_2partsintovectorNew() interface{} {
	return &Ligation_2partsintovectorElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Ligation_2partsintovectorInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Ligation_2partsintovectorRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Ligation_2partsintovectorInput{},
			Out: &Ligation_2partsintovectorOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Ligation_2partsintovectorElement struct {
	inject.CheckedRunner
}

type Ligation_2partsintovectorInput struct {
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
	Part1              *wtype.LHComponent
	Part2              *wtype.LHComponent
	PartVols           []wunit.Volume
	ReactionTemp       wunit.Temperature
	ReactionTime       wunit.Time
	ReactionVolume     wunit.Volume
	Vector             *wtype.LHComponent
	VectorVol          wunit.Volume
	Water              *wtype.LHComponent
}

type Ligation_2partsintovectorOutput struct {
	Reaction *wtype.LHComponent
}

type Ligation_2partsintovectorSOutput struct {
	Data struct {
	}
	Outputs struct {
		Reaction *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Ligation_2partsintovector",
		Constructor: Ligation_2partsintovectorNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/Ligation/Ligation.an",
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
				{Name: "Part1", Desc: "", Kind: "Inputs"},
				{Name: "Part2", Desc: "", Kind: "Inputs"},
				{Name: "PartVols", Desc: "", Kind: "Parameters"},
				{Name: "ReactionTemp", Desc: "", Kind: "Parameters"},
				{Name: "ReactionTime", Desc: "", Kind: "Parameters"},
				{Name: "ReactionVolume", Desc: "", Kind: "Parameters"},
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
