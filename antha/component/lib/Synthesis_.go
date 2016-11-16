package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/factory"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

func _SynthesisRequirements() {}

// Conditions to run on startup
func _SynthesisSetup(_ctx context.Context, _input *SynthesisInput) {}

func _SynthesisSteps(_ctx context.Context, _input *SynthesisInput, _output *SynthesisOutput) {
	// Element with mock synthesises DNA. Converts DNA sequence type to LHComponent.
	//var dna *wtype.LHComponent
	components := make([]*wtype.LHComponent, 0)

	fmt.Println("Parts:", _input.PartsWithOverhangs)
	for i, part := range _input.PartsWithOverhangs {

		dna := factory.GetComponentByType("dna")
		dna.CName = part.Nm
		fmt.Println("dna:", i, dna)
		components = append(components, dna)
	}
	_output.Components = components

	fmt.Println("Components source:", _output.Components)
}

func _SynthesisAnalysis(_ctx context.Context, _input *SynthesisInput, _output *SynthesisOutput) {
}

func _SynthesisValidation(_ctx context.Context, _input *SynthesisInput, _output *SynthesisOutput) {
}
func _SynthesisRun(_ctx context.Context, input *SynthesisInput) *SynthesisOutput {
	output := &SynthesisOutput{}
	_SynthesisSetup(_ctx, input)
	_SynthesisSteps(_ctx, input, output)
	_SynthesisAnalysis(_ctx, input, output)
	_SynthesisValidation(_ctx, input, output)
	return output
}

func SynthesisRunSteps(_ctx context.Context, input *SynthesisInput) *SynthesisSOutput {
	soutput := &SynthesisSOutput{}
	output := _SynthesisRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SynthesisNew() interface{} {
	return &SynthesisElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SynthesisInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SynthesisRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SynthesisInput{},
			Out: &SynthesisOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type SynthesisElement struct {
	inject.CheckedRunner
}

type SynthesisInput struct {
	PartsWithOverhangs []wtype.DNASequence
}

type SynthesisOutput struct {
	Components []*wtype.LHComponent
}

type SynthesisSOutput struct {
	Data struct {
	}
	Outputs struct {
		Components []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Synthesis",
		Constructor: SynthesisNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/DNA/GeneDesign/Synthesis.an",
			Params: []component.ParamDesc{
				{Name: "PartsWithOverhangs", Desc: "", Kind: "Parameters"},
				{Name: "Components", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
