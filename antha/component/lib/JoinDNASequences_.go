// Demo protocol of how to create an array of dna types from parsing user inputs of various types
// scenarios handled:
// Biobrick IDS
// genbank files
// raw sequence
// inventory lookup
package lib

import (
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol

// Data which is returned from this protocol

// Physical inputs to this protocol

// Physical outputs from this protocol

func _JoinDNASequencesRequirements() {

}

// Actions to perform before protocol itself
func _JoinDNASequencesSetup(_ctx context.Context, _input *JoinDNASequencesInput) {

}

// Core process of the protocol: steps to be performed for each input
func _JoinDNASequencesSteps(_ctx context.Context, _input *JoinDNASequencesInput, _output *JoinDNASequencesOutput) {

	newSeq := _input.Seqsinorder[0]
	seqnames := make([]string, 0)

	for i, seq := range _input.Seqsinorder {
		if i != 0 {
			newSeq.Append(seq.Seq)
		}
		seqnames = append(seqnames, seq.Nm)
	}

	newSeq.Nm = strings.Join(seqnames, "_")

	_output.Seq = newSeq

}

// Actions to perform after steps block to analyze data
func _JoinDNASequencesAnalysis(_ctx context.Context, _input *JoinDNASequencesInput, _output *JoinDNASequencesOutput) {

}

func _JoinDNASequencesValidation(_ctx context.Context, _input *JoinDNASequencesInput, _output *JoinDNASequencesOutput) {

}
func _JoinDNASequencesRun(_ctx context.Context, input *JoinDNASequencesInput) *JoinDNASequencesOutput {
	output := &JoinDNASequencesOutput{}
	_JoinDNASequencesSetup(_ctx, input)
	_JoinDNASequencesSteps(_ctx, input, output)
	_JoinDNASequencesAnalysis(_ctx, input, output)
	_JoinDNASequencesValidation(_ctx, input, output)
	return output
}

func JoinDNASequencesRunSteps(_ctx context.Context, input *JoinDNASequencesInput) *JoinDNASequencesSOutput {
	soutput := &JoinDNASequencesSOutput{}
	output := _JoinDNASequencesRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func JoinDNASequencesNew() interface{} {
	return &JoinDNASequencesElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &JoinDNASequencesInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _JoinDNASequencesRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &JoinDNASequencesInput{},
			Out: &JoinDNASequencesOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type JoinDNASequencesElement struct {
	inject.CheckedRunner
}

type JoinDNASequencesInput struct {
	Seqsinorder []wtype.DNASequence
}

type JoinDNASequencesOutput struct {
	Seq wtype.DNASequence
}

type JoinDNASequencesSOutput struct {
	Data struct {
		Seq wtype.DNASequence
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "JoinDNASequences",
		Constructor: JoinDNASequencesNew,
		Desc: component.ComponentDesc{
			Desc: "Demo protocol of how to create an array of dna types from parsing user inputs of various types\nscenarios handled:\nBiobrick IDS\ngenbank files\nraw sequence\ninventory lookup\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/AnthaAcademy/Lesson4_DNA/E_JoinDNASequences.an",
			Params: []component.ParamDesc{
				{Name: "Seqsinorder", Desc: "", Kind: "Parameters"},
				{Name: "Seq", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
