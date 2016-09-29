// This element will design a sequencing primer to target amplification of a target region within a sequence file
// Design criteria such as maximum gc content, acceptable ranges of melting temperatures and primer length may be specified by the user.
package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/oligos"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"

// Input parameters for this protocol

// string // genbank file (.gb)

// as a proportion of 1, i.e. 1 == 100%

// number of nucleotides which primers can overlap by

// Data which is returned from this protocol

// Physical inputs to this protocol

// Physical outputs from this protocol

func _PrimerDesign_FWD_wtypeRequirements() {

}

// Actions to perform before protocol itself
func _PrimerDesign_FWD_wtypeSetup(_ctx context.Context, _input *PrimerDesign_FWD_wtypeInput) {

}

// Core process of the protocol: steps to be performed for each input
func _PrimerDesign_FWD_wtypeSteps(_ctx context.Context, _input *PrimerDesign_FWD_wtypeInput, _output *PrimerDesign_FWD_wtypeOutput) {

	var region wtype.DNASequence

	regionstart, regionend, err := oligos.FindPositioninSequence(_input.FullDNASeq, _input.RegionSequence)

	if err != nil {
		fmt.Println("FindPositioninoligoFail")
		fmt.Println(err.Error())
		_output.Warnings = err
		execute.Errorf(_ctx, _output.Warnings.Error())
	}

	// if true then the start point to design primers is moved back 150bp to ensure full region is covered
	if _input.FlankTargetSequence && regionstart-150 >= 0 {
		region = oligos.DNAregion(_input.FullDNASeq, regionstart-150, regionend)
	} else if _input.FlankTargetSequence && regionstart-150 < 0 && regionstart-_input.Maxlength >= 0 {
		region = oligos.DNAregion(_input.FullDNASeq, 0, regionend)
	} else if _input.FlankTargetSequence && regionstart-150 < 0 && regionstart-_input.Maxlength < 0 && _input.FullDNASeq.Plasmid {
		region = oligos.DNAregion(_input.FullDNASeq, len(_input.FullDNASeq.Seq)-150, regionend)
	} else {
		region = oligos.DNAregion(_input.FullDNASeq, regionstart, regionend)
	}

	_output.FWDPrimer, _output.Warnings = oligos.FWDOligoSeq(region, _input.Maxgc, _input.Minlength, _input.Maxlength, _input.Mintemp, _input.Maxtemp, _input.Seqstoavoid, _input.PermittednucleotideOverlapBetweenPrimers)

	if _output.Warnings != nil {
		fmt.Println("FWDOligoSeqFail")
		errstr := _output.Warnings.Error()
		fmt.Println(errstr)
		execute.Errorf(_ctx, errstr)
	}

	fmt.Println(text.Print("FWDPrimer:", _output.FWDPrimer))
}

// Actions to perform after steps block to analyze data
func _PrimerDesign_FWD_wtypeAnalysis(_ctx context.Context, _input *PrimerDesign_FWD_wtypeInput, _output *PrimerDesign_FWD_wtypeOutput) {

}

func _PrimerDesign_FWD_wtypeValidation(_ctx context.Context, _input *PrimerDesign_FWD_wtypeInput, _output *PrimerDesign_FWD_wtypeOutput) {

}
func _PrimerDesign_FWD_wtypeRun(_ctx context.Context, input *PrimerDesign_FWD_wtypeInput) *PrimerDesign_FWD_wtypeOutput {
	output := &PrimerDesign_FWD_wtypeOutput{}
	_PrimerDesign_FWD_wtypeSetup(_ctx, input)
	_PrimerDesign_FWD_wtypeSteps(_ctx, input, output)
	_PrimerDesign_FWD_wtypeAnalysis(_ctx, input, output)
	_PrimerDesign_FWD_wtypeValidation(_ctx, input, output)
	return output
}

func PrimerDesign_FWD_wtypeRunSteps(_ctx context.Context, input *PrimerDesign_FWD_wtypeInput) *PrimerDesign_FWD_wtypeSOutput {
	soutput := &PrimerDesign_FWD_wtypeSOutput{}
	output := _PrimerDesign_FWD_wtypeRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PrimerDesign_FWD_wtypeNew() interface{} {
	return &PrimerDesign_FWD_wtypeElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PrimerDesign_FWD_wtypeInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PrimerDesign_FWD_wtypeRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PrimerDesign_FWD_wtypeInput{},
			Out: &PrimerDesign_FWD_wtypeOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PrimerDesign_FWD_wtypeElement struct {
	inject.CheckedRunner
}

type PrimerDesign_FWD_wtypeInput struct {
	FlankTargetSequence                      bool
	FullDNASeq                               wtype.DNASequence
	Maxgc                                    float64
	Maxlength                                int
	Maxtemp                                  wunit.Temperature
	Minlength                                int
	Mintemp                                  wunit.Temperature
	PermittednucleotideOverlapBetweenPrimers int
	RegionSequence                           wtype.DNASequence
	Seqstoavoid                              []string
}

type PrimerDesign_FWD_wtypeOutput struct {
	FWDPrimer oligos.Primer
	Warnings  error
}

type PrimerDesign_FWD_wtypeSOutput struct {
	Data struct {
		FWDPrimer oligos.Primer
		Warnings  error
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PrimerDesign_FWD_wtype",
		Constructor: PrimerDesign_FWD_wtypeNew,
		Desc: component.ComponentDesc{
			Desc: "This element will design a sequencing primer to target amplification of a target region within a sequence file\nDesign criteria such as maximum gc content, acceptable ranges of melting temperatures and primer length may be specified by the user.\n",
			Path: "antha/component/an/Data/DNA/PrimerDesign/PrimerDesign_FWD_wtype.an",
			Params: []component.ParamDesc{
				{Name: "FlankTargetSequence", Desc: "", Kind: "Parameters"},
				{Name: "FullDNASeq", Desc: "string // genbank file (.gb)\n", Kind: "Parameters"},
				{Name: "Maxgc", Desc: "as a proportion of 1, i.e. 1 == 100%\n", Kind: "Parameters"},
				{Name: "Maxlength", Desc: "", Kind: "Parameters"},
				{Name: "Maxtemp", Desc: "", Kind: "Parameters"},
				{Name: "Minlength", Desc: "", Kind: "Parameters"},
				{Name: "Mintemp", Desc: "", Kind: "Parameters"},
				{Name: "PermittednucleotideOverlapBetweenPrimers", Desc: "number of nucleotides which primers can overlap by\n", Kind: "Parameters"},
				{Name: "RegionSequence", Desc: "", Kind: "Parameters"},
				{Name: "Seqstoavoid", Desc: "", Kind: "Parameters"},
				{Name: "FWDPrimer", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
