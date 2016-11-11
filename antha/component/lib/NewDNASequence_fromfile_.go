// Protocol for creating a DNASequence from a sequence file format. // Supported formats: .gdx .fasta .gb
package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
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

func _NewDNASequence_fromfileRequirements() {

}

// Actions to perform before protocol itself
func _NewDNASequence_fromfileSetup(_ctx context.Context, _input *NewDNASequence_fromfileInput) {

}

// Core process of the protocol: steps to be performed for each input
func _NewDNASequence_fromfileSteps(_ctx context.Context, _input *NewDNASequence_fromfileInput, _output *NewDNASequence_fromfileOutput) {

	seqs, err := parser.DNAFiletoDNASequence(_input.Filename, _input.Plasmid)

	if len(seqs) == 1 && err == nil {

		_output.DNA = seqs[0]

		_output.DNA.Nm = _input.Gene_name

	}

	orfs := sequences.FindallORFs(_output.DNA.Seq)

	if len(_output.DNA.Features) == 0 {
		features := sequences.ORFs2Features(orfs)

		_output.DNA = wtype.Annotate(_output.DNA, features)
	}

	_output.Status = fmt.Sprintln(
		text.Print("DNA_Seq: ", _output.DNA),
		text.Print("ORFs: ", orfs),
	)

	_output.Warnings = err

}

// Actions to perform after steps block to analyze data
func _NewDNASequence_fromfileAnalysis(_ctx context.Context, _input *NewDNASequence_fromfileInput, _output *NewDNASequence_fromfileOutput) {

}

func _NewDNASequence_fromfileValidation(_ctx context.Context, _input *NewDNASequence_fromfileInput, _output *NewDNASequence_fromfileOutput) {

}
func _NewDNASequence_fromfileRun(_ctx context.Context, input *NewDNASequence_fromfileInput) *NewDNASequence_fromfileOutput {
	output := &NewDNASequence_fromfileOutput{}
	_NewDNASequence_fromfileSetup(_ctx, input)
	_NewDNASequence_fromfileSteps(_ctx, input, output)
	_NewDNASequence_fromfileAnalysis(_ctx, input, output)
	_NewDNASequence_fromfileValidation(_ctx, input, output)
	return output
}

func NewDNASequence_fromfileRunSteps(_ctx context.Context, input *NewDNASequence_fromfileInput) *NewDNASequence_fromfileSOutput {
	soutput := &NewDNASequence_fromfileSOutput{}
	output := _NewDNASequence_fromfileRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func NewDNASequence_fromfileNew() interface{} {
	return &NewDNASequence_fromfileElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &NewDNASequence_fromfileInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _NewDNASequence_fromfileRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &NewDNASequence_fromfileInput{},
			Out: &NewDNASequence_fromfileOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type NewDNASequence_fromfileElement struct {
	inject.CheckedRunner
}

type NewDNASequence_fromfileInput struct {
	Filename       string
	Gene_name      string
	Linear         bool
	Plasmid        bool
	SingleStranded bool
}

type NewDNASequence_fromfileOutput struct {
	DNA      wtype.DNASequence
	Status   string
	Warnings error
}

type NewDNASequence_fromfileSOutput struct {
	Data struct {
		DNA      wtype.DNASequence
		Status   string
		Warnings error
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "NewDNASequence_fromfile",
		Constructor: NewDNASequence_fromfileNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol for creating a DNASequence from a sequence file format. // Supported formats: .gdx .fasta .gb\n",
			Path: "antha/component/an/AnthaAcademy/Lesson6_DNA/B_NewDNASequence_fromfile.an",
			Params: []component.ParamDesc{
				{Name: "Filename", Desc: "", Kind: "Parameters"},
				{Name: "Gene_name", Desc: "", Kind: "Parameters"},
				{Name: "Linear", Desc: "", Kind: "Parameters"},
				{Name: "Plasmid", Desc: "", Kind: "Parameters"},
				{Name: "SingleStranded", Desc: "", Kind: "Parameters"},
				{Name: "DNA", Desc: "", Kind: "Data"},
				{Name: "Status", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
