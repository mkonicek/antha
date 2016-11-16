// combine up to 8 sequences into array; for convenience in using workflow editor
package lib

import (
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

func _MakeSequenceArrayRequirements() {

}

// Actions to perform before protocol itself
func _MakeSequenceArraySetup(_ctx context.Context, _input *MakeSequenceArrayInput) {

}

// Core process of the protocol: steps to be performed for each input
func _MakeSequenceArraySteps(_ctx context.Context, _input *MakeSequenceArrayInput, _output *MakeSequenceArrayOutput) {

	seqs := make([]wtype.DNASequence, 0)

	if len(_input.Seq1.Seq) > 0 {
		seqs = append(seqs, _input.Seq1)
	}
	if len(_input.Seq2.Seq) > 0 {
		seqs = append(seqs, _input.Seq2)
	}
	if len(_input.Seq3.Seq) > 0 {
		seqs = append(seqs, _input.Seq3)
	}
	if len(_input.Seq4.Seq) > 0 {
		seqs = append(seqs, _input.Seq4)
	}
	if len(_input.Seq5.Seq) > 0 {
		seqs = append(seqs, _input.Seq5)
	}
	if len(_input.Seq6.Seq) > 0 {
		seqs = append(seqs, _input.Seq6)
	}
	if len(_input.Seq7.Seq) > 0 {
		seqs = append(seqs, _input.Seq7)
	}
	if len(_input.Seq8.Seq) > 0 {
		seqs = append(seqs, _input.Seq8)
	}
	_output.Seqs = seqs

}

// Actions to perform after steps block to analyze data
func _MakeSequenceArrayAnalysis(_ctx context.Context, _input *MakeSequenceArrayInput, _output *MakeSequenceArrayOutput) {

}

func _MakeSequenceArrayValidation(_ctx context.Context, _input *MakeSequenceArrayInput, _output *MakeSequenceArrayOutput) {

}
func _MakeSequenceArrayRun(_ctx context.Context, input *MakeSequenceArrayInput) *MakeSequenceArrayOutput {
	output := &MakeSequenceArrayOutput{}
	_MakeSequenceArraySetup(_ctx, input)
	_MakeSequenceArraySteps(_ctx, input, output)
	_MakeSequenceArrayAnalysis(_ctx, input, output)
	_MakeSequenceArrayValidation(_ctx, input, output)
	return output
}

func MakeSequenceArrayRunSteps(_ctx context.Context, input *MakeSequenceArrayInput) *MakeSequenceArraySOutput {
	soutput := &MakeSequenceArraySOutput{}
	output := _MakeSequenceArrayRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MakeSequenceArrayNew() interface{} {
	return &MakeSequenceArrayElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MakeSequenceArrayInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MakeSequenceArrayRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MakeSequenceArrayInput{},
			Out: &MakeSequenceArrayOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type MakeSequenceArrayElement struct {
	inject.CheckedRunner
}

type MakeSequenceArrayInput struct {
	Seq1 wtype.DNASequence
	Seq2 wtype.DNASequence
	Seq3 wtype.DNASequence
	Seq4 wtype.DNASequence
	Seq5 wtype.DNASequence
	Seq6 wtype.DNASequence
	Seq7 wtype.DNASequence
	Seq8 wtype.DNASequence
}

type MakeSequenceArrayOutput struct {
	Seqs []wtype.DNASequence
}

type MakeSequenceArraySOutput struct {
	Data struct {
		Seqs []wtype.DNASequence
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MakeSequenceArray",
		Constructor: MakeSequenceArrayNew,
		Desc: component.ComponentDesc{
			Desc: "combine up to 8 sequences into array; for convenience in using workflow editor\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/AnthaAcademy/Lesson4_DNA/F_MakeDNASequenceArray.an",
			Params: []component.ParamDesc{
				{Name: "Seq1", Desc: "", Kind: "Parameters"},
				{Name: "Seq2", Desc: "", Kind: "Parameters"},
				{Name: "Seq3", Desc: "", Kind: "Parameters"},
				{Name: "Seq4", Desc: "", Kind: "Parameters"},
				{Name: "Seq5", Desc: "", Kind: "Parameters"},
				{Name: "Seq6", Desc: "", Kind: "Parameters"},
				{Name: "Seq7", Desc: "", Kind: "Parameters"},
				{Name: "Seq8", Desc: "", Kind: "Parameters"},
				{Name: "Seqs", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
