// This element will design primers to cover the full length of a sequence at the interval specified by the user (e.g. every 800 bp).
// Design criteria such as maximum gc content, acceptable ranges of melting temperatures and primer length may be specified by the user.
package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/oligos"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol

// as a proportion of 1, i.e. 1 == 100%

// number of nucleotides which primers can overlap by

// Data which is returned from this protocol

// Physical inputs to this protocol

// Physical outputs from this protocol

func _PrimerDesign_coverfullsequenceRequirements() {

}

// Actions to perform before protocol itself
func _PrimerDesign_coverfullsequenceSetup(_ctx context.Context, _input *PrimerDesign_coverfullsequenceInput) {

}

// Core process of the protocol: steps to be performed for each input
func _PrimerDesign_coverfullsequenceSteps(_ctx context.Context, _input *PrimerDesign_coverfullsequenceInput, _output *PrimerDesign_coverfullsequenceOutput) {
	var plasmid wtype.DNASequence

	plasmids, _ := parser.DNAFiletoDNASequence(_input.DNASeqfile, _input.Plasmid)

	if len(plasmids) > 0 {
		plasmid = plasmids[0]
	}
	if len(plasmids) > 1 {
		_output.Warnings = fmt.Errorf("Warning! more than one sequence in sequence file! Only used first sequence for primer design")
	}

	allprimers := oligos.DesignFWDPRimerstoCoverFullSequence(plasmid, _input.PrimereveryXnucleotides, _input.Maxgc, _input.Minlength, _input.Maxlength, _input.Mintemp, _input.Maxtemp, _input.Seqstoavoid, _input.PermittednucleotideOverlapBetweenPrimers)

	_output.AllPrimers = allprimers
}

// Actions to perform after steps block to analyze data
func _PrimerDesign_coverfullsequenceAnalysis(_ctx context.Context, _input *PrimerDesign_coverfullsequenceInput, _output *PrimerDesign_coverfullsequenceOutput) {

}

func _PrimerDesign_coverfullsequenceValidation(_ctx context.Context, _input *PrimerDesign_coverfullsequenceInput, _output *PrimerDesign_coverfullsequenceOutput) {

}
func _PrimerDesign_coverfullsequenceRun(_ctx context.Context, input *PrimerDesign_coverfullsequenceInput) *PrimerDesign_coverfullsequenceOutput {
	output := &PrimerDesign_coverfullsequenceOutput{}
	_PrimerDesign_coverfullsequenceSetup(_ctx, input)
	_PrimerDesign_coverfullsequenceSteps(_ctx, input, output)
	_PrimerDesign_coverfullsequenceAnalysis(_ctx, input, output)
	_PrimerDesign_coverfullsequenceValidation(_ctx, input, output)
	return output
}

func PrimerDesign_coverfullsequenceRunSteps(_ctx context.Context, input *PrimerDesign_coverfullsequenceInput) *PrimerDesign_coverfullsequenceSOutput {
	soutput := &PrimerDesign_coverfullsequenceSOutput{}
	output := _PrimerDesign_coverfullsequenceRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PrimerDesign_coverfullsequenceNew() interface{} {
	return &PrimerDesign_coverfullsequenceElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PrimerDesign_coverfullsequenceInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PrimerDesign_coverfullsequenceRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PrimerDesign_coverfullsequenceInput{},
			Out: &PrimerDesign_coverfullsequenceOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PrimerDesign_coverfullsequenceElement struct {
	inject.CheckedRunner
}

type PrimerDesign_coverfullsequenceInput struct {
	DNASeqfile                               string
	Maxgc                                    float64
	Maxlength                                int
	Maxtemp                                  wunit.Temperature
	Minlength                                int
	Mintemp                                  wunit.Temperature
	PermittednucleotideOverlapBetweenPrimers int
	Plasmid                                  bool
	PrimereveryXnucleotides                  int
	Seqstoavoid                              []string
}

type PrimerDesign_coverfullsequenceOutput struct {
	AllPrimers []oligos.Primer
	Warnings   error
}

type PrimerDesign_coverfullsequenceSOutput struct {
	Data struct {
		AllPrimers []oligos.Primer
		Warnings   error
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PrimerDesign_coverfullsequence",
		Constructor: PrimerDesign_coverfullsequenceNew,
		Desc: component.ComponentDesc{
			Desc: "This element will design primers to cover the full length of a sequence at the interval specified by the user (e.g. every 800 bp).\nDesign criteria such as maximum gc content, acceptable ranges of melting temperatures and primer length may be specified by the user.\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/DNA/PrimerDesign/PrimerDesign_coverfullsequence.an",
			Params: []component.ParamDesc{
				{Name: "DNASeqfile", Desc: "", Kind: "Parameters"},
				{Name: "Maxgc", Desc: "as a proportion of 1, i.e. 1 == 100%\n", Kind: "Parameters"},
				{Name: "Maxlength", Desc: "", Kind: "Parameters"},
				{Name: "Maxtemp", Desc: "", Kind: "Parameters"},
				{Name: "Minlength", Desc: "", Kind: "Parameters"},
				{Name: "Mintemp", Desc: "", Kind: "Parameters"},
				{Name: "PermittednucleotideOverlapBetweenPrimers", Desc: "number of nucleotides which primers can overlap by\n", Kind: "Parameters"},
				{Name: "Plasmid", Desc: "", Kind: "Parameters"},
				{Name: "PrimereveryXnucleotides", Desc: "", Kind: "Parameters"},
				{Name: "Seqstoavoid", Desc: "", Kind: "Parameters"},
				{Name: "AllPrimers", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
