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
	"strings"
)

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

func _TypeIISConstructAssembly_altRequirements() {}

// Conditions to run on startup
func _TypeIISConstructAssembly_altSetup(_ctx context.Context, _input *TypeIISConstructAssembly_altInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _TypeIISConstructAssembly_altSteps(_ctx context.Context, _input *TypeIISConstructAssembly_altInput, _output *TypeIISConstructAssembly_altOutput) {

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

	s := ""
	comments := make([]string, 0)
	var partSample *wtype.LHComponent

	for k, part := range _input.Parts {
		if _input.PartConcs[k].SIValue() <= 0.1 {
			s = fmt.Sprintln("creating dna part num ", k, " comp ", part.CName, " renamed to ", _input.PartNames[k], " vol ", _input.PartConcs[k].ToString())
			partSample = mixer.SampleForConcentration(part, _input.PartConcs[k])
		} else {
			s = fmt.Sprintln("Conc too low so minimum volume used", "creating dna part num ", k, " comp ", part.CName, " renamed to ", _input.PartNames[k], " vol ", _input.PartMinVol.ToString())
			partSample = mixer.Sample(part, _input.PartMinVol)
		}
		partSample.CName = _input.PartNames[k]
		samples = append(samples, partSample)
		comments = append(comments, s)

	}
	_output.S = strings.Join(comments, "")

	reSample := mixer.Sample(_input.RestrictionEnzyme, _input.ReVol)
	samples = append(samples, reSample)

	ligSample := mixer.Sample(_input.Ligase, _input.LigVol)
	samples = append(samples, ligSample)

	// incubate the reaction mixture
	out1 := execute.Incubate(_ctx, execute.MixInto(_ctx, _input.OutPlate, "", samples...), _input.ReactionTemp, _input.ReactionTime, false)
	// inactivate
	_output.Reaction = execute.Incubate(_ctx, out1, _input.InactivationTemp, _input.InactivationTime, false)
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _TypeIISConstructAssembly_altAnalysis(_ctx context.Context, _input *TypeIISConstructAssembly_altInput, _output *TypeIISConstructAssembly_altOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _TypeIISConstructAssembly_altValidation(_ctx context.Context, _input *TypeIISConstructAssembly_altInput, _output *TypeIISConstructAssembly_altOutput) {
}
func _TypeIISConstructAssembly_altRun(_ctx context.Context, input *TypeIISConstructAssembly_altInput) *TypeIISConstructAssembly_altOutput {
	output := &TypeIISConstructAssembly_altOutput{}
	_TypeIISConstructAssembly_altSetup(_ctx, input)
	_TypeIISConstructAssembly_altSteps(_ctx, input, output)
	_TypeIISConstructAssembly_altAnalysis(_ctx, input, output)
	_TypeIISConstructAssembly_altValidation(_ctx, input, output)
	return output
}

func TypeIISConstructAssembly_altRunSteps(_ctx context.Context, input *TypeIISConstructAssembly_altInput) *TypeIISConstructAssembly_altSOutput {
	soutput := &TypeIISConstructAssembly_altSOutput{}
	output := _TypeIISConstructAssembly_altRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func TypeIISConstructAssembly_altNew() interface{} {
	return &TypeIISConstructAssembly_altElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &TypeIISConstructAssembly_altInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _TypeIISConstructAssembly_altRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &TypeIISConstructAssembly_altInput{},
			Out: &TypeIISConstructAssembly_altOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type TypeIISConstructAssembly_altElement struct {
	inject.CheckedRunner
}

type TypeIISConstructAssembly_altInput struct {
	Atp               *wtype.LHComponent
	AtpVol            wunit.Volume
	Buffer            *wtype.LHComponent
	BufferVol         wunit.Volume
	InPlate           *wtype.LHPlate
	InactivationTemp  wunit.Temperature
	InactivationTime  wunit.Time
	LigVol            wunit.Volume
	Ligase            *wtype.LHComponent
	OutPlate          *wtype.LHPlate
	PartConcs         []wunit.Concentration
	PartMinVol        wunit.Volume
	PartNames         []string
	Parts             []*wtype.LHComponent
	ReVol             wunit.Volume
	ReactionTemp      wunit.Temperature
	ReactionTime      wunit.Time
	ReactionVolume    wunit.Volume
	RestrictionEnzyme *wtype.LHComponent
	Vector            *wtype.LHComponent
	VectorVol         wunit.Volume
	Water             *wtype.LHComponent
}

type TypeIISConstructAssembly_altOutput struct {
	Reaction *wtype.LHComponent
	S        string
}

type TypeIISConstructAssembly_altSOutput struct {
	Data struct {
		S string
	}
	Outputs struct {
		Reaction *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "TypeIISConstructAssembly_alt",
		Constructor: TypeIISConstructAssembly_altNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/TypeIIsAssembly/TypeIISConstructAssembly/TypeIISConstructAssembly_alt.an",
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
				{Name: "PartConcs", Desc: "", Kind: "Parameters"},
				{Name: "PartMinVol", Desc: "", Kind: "Parameters"},
				{Name: "PartNames", Desc: "", Kind: "Parameters"},
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
				{Name: "S", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
