// demo protocol of how to create a dna type from user inputs
package lib

import (
	"fmt"
	//"math"
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

func _NewDNASequenceRequirements() {

}

// Actions to perform before protocol itself
func _NewDNASequenceSetup(_ctx context.Context, _input *NewDNASequenceInput) {

}

// Core process of the protocol: steps to be performed for each input
func _NewDNASequenceSteps(_ctx context.Context, _input *NewDNASequenceInput, _output *NewDNASequenceOutput) {

	// != is go syntax for not equal to
	if _input.Plasmid != _input.Linear {

		// equivalent to if Plasmid == true
		if _input.Plasmid {
			// different functions exist for making an antha DNA sequence based on the properties
			_output.DNA = wtype.MakePlasmidDNASequence(_input.Gene_name, _input.DNA_seq)

		} else if _input.Linear {

			_output.DNA = wtype.MakeLinearDNASequence(_input.Gene_name, _input.DNA_seq)

		} else if _input.SingleStranded {

			_output.DNA = wtype.MakeSingleStrandedDNASequence(_input.Gene_name, _input.DNA_seq)

		}

		// use FindallORFs from sequences library
		orfs := sequences.FindallORFs(_output.DNA.Seq)

		// convert those orfs to features
		features := sequences.ORFs2Features(orfs)

		// add annotations to sequence from features
		_output.DNAwithORFs = wtype.Annotate(_output.DNA, features)

		_output.Status = fmt.Sprintln(
			text.Print("DNA_Seq: ", _input.DNA_seq),
			text.Print("ORFs: ", _output.DNAwithORFs.Features),
		)

	} else {
		_output.Status = fmt.Sprintln("correct conditions not met")
	}

}

// Actions to perform after steps block to analyze data
func _NewDNASequenceAnalysis(_ctx context.Context, _input *NewDNASequenceInput, _output *NewDNASequenceOutput) {

}

func _NewDNASequenceValidation(_ctx context.Context, _input *NewDNASequenceInput, _output *NewDNASequenceOutput) {

}
func _NewDNASequenceRun(_ctx context.Context, input *NewDNASequenceInput) *NewDNASequenceOutput {
	output := &NewDNASequenceOutput{}
	_NewDNASequenceSetup(_ctx, input)
	_NewDNASequenceSteps(_ctx, input, output)
	_NewDNASequenceAnalysis(_ctx, input, output)
	_NewDNASequenceValidation(_ctx, input, output)
	return output
}

func NewDNASequenceRunSteps(_ctx context.Context, input *NewDNASequenceInput) *NewDNASequenceSOutput {
	soutput := &NewDNASequenceSOutput{}
	output := _NewDNASequenceRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func NewDNASequenceNew() interface{} {
	return &NewDNASequenceElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &NewDNASequenceInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _NewDNASequenceRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &NewDNASequenceInput{},
			Out: &NewDNASequenceOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type NewDNASequenceElement struct {
	inject.CheckedRunner
}

type NewDNASequenceInput struct {
	DNA_seq        string
	Gene_name      string
	Linear         bool
	Plasmid        bool
	SingleStranded bool
}

type NewDNASequenceOutput struct {
	DNA         wtype.DNASequence
	DNAwithORFs wtype.DNASequence
	Status      string
}

type NewDNASequenceSOutput struct {
	Data struct {
		DNA         wtype.DNASequence
		DNAwithORFs wtype.DNASequence
		Status      string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "NewDNASequence",
		Constructor: NewDNASequenceNew,
		Desc: component.ComponentDesc{
			Desc: "demo protocol of how to create a dna type from user inputs\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/AnthaAcademy/Lesson6_DNA/A_NewDNASequence.an",
			Params: []component.ParamDesc{
				{Name: "DNA_seq", Desc: "", Kind: "Parameters"},
				{Name: "Gene_name", Desc: "", Kind: "Parameters"},
				{Name: "Linear", Desc: "", Kind: "Parameters"},
				{Name: "Plasmid", Desc: "", Kind: "Parameters"},
				{Name: "SingleStranded", Desc: "", Kind: "Parameters"},
				{Name: "DNA", Desc: "", Kind: "Data"},
				{Name: "DNAwithORFs", Desc: "", Kind: "Data"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
