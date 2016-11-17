// This element will design a sequencing primer to target amplification of a target region within a sequence file
// Design criteria such as maximum gc content, acceptable ranges of melting temperatures and primer length may be specified by the user.
package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/oligos"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol

// genbank file (.gb)

// as a proportion of 1, i.e. 1 == 100%

// number of nucleotides which primers can overlap by

// Data which is returned from this protocol

// Physical inputs to this protocol

// Physical outputs from this protocol

func _PrimerDesign_FWDRequirements() {

}

// Actions to perform before protocol itself
func _PrimerDesign_FWDSetup(_ctx context.Context, _input *PrimerDesign_FWDInput) {

}

// Core process of the protocol: steps to be performed for each input
func _PrimerDesign_FWDSteps(_ctx context.Context, _input *PrimerDesign_FWDInput, _output *PrimerDesign_FWDOutput) {

	var region wtype.DNASequence

	fulldnaseqs, err := parser.DNAFiletoDNASequence(_input.DNASeqfile, _input.Plasmid)

	if err != nil {
		fmt.Println("ParseSeqFail")
		_output.Warnings = err
		execute.Errorf(_ctx, _output.Warnings.Error())
	}
	fmt.Println("1")
	if len(fulldnaseqs) != 1 {
		_output.Warnings = fmt.Errorf("more than one matching dna sequence found in target sequence")
		execute.Errorf(_ctx, _output.Warnings.Error())
	}
	fmt.Println("2")
	RegionSequence := wtype.MakeLinearDNASequence("region", _input.RegionSequenceString)
	fmt.Println("3")
	fmt.Println("fulldnaseqs[0]", fulldnaseqs[0])
	fmt.Println("RegionSequence", RegionSequence)
	regionstart, regionend, err := oligos.FindPositioninSequence(fulldnaseqs[0], RegionSequence)
	fmt.Println("4")
	if err != nil {
		fmt.Println("FindPositioninoligoFail")
		_output.Warnings = err
		execute.Errorf(_ctx, _output.Warnings.Error())
	}
	fmt.Println("5")
	// if true then the start oint ot design primers is moved back 150bp to ensure full region is covered
	if _input.FlankTargetSequence {
		region = oligos.DNAregion(fulldnaseqs[0], regionstart-150, regionend)
	} else {
		region = oligos.DNAregion(fulldnaseqs[0], regionstart, regionend)
	}
	fmt.Println("6")

	_output.FWDPrimer, _output.Warnings = oligos.FWDOligoSeq(region, _input.Maxgc, _input.Minlength, _input.Maxlength, _input.Mintemp, _input.Maxtemp, _input.Seqstoavoid, _input.PermittednucleotideOverlapBetweenPrimers)

	if _output.Warnings != nil {
		fmt.Println("FWDOligoSeqFail")
		execute.Errorf(_ctx, _output.Warnings.Error())
	}

	fmt.Println(text.Print("FWDPrimer:", _output.FWDPrimer))
}

// Actions to perform after steps block to analyze data
func _PrimerDesign_FWDAnalysis(_ctx context.Context, _input *PrimerDesign_FWDInput, _output *PrimerDesign_FWDOutput) {

}

func _PrimerDesign_FWDValidation(_ctx context.Context, _input *PrimerDesign_FWDInput, _output *PrimerDesign_FWDOutput) {

}
func _PrimerDesign_FWDRun(_ctx context.Context, input *PrimerDesign_FWDInput) *PrimerDesign_FWDOutput {
	output := &PrimerDesign_FWDOutput{}
	_PrimerDesign_FWDSetup(_ctx, input)
	_PrimerDesign_FWDSteps(_ctx, input, output)
	_PrimerDesign_FWDAnalysis(_ctx, input, output)
	_PrimerDesign_FWDValidation(_ctx, input, output)
	return output
}

func PrimerDesign_FWDRunSteps(_ctx context.Context, input *PrimerDesign_FWDInput) *PrimerDesign_FWDSOutput {
	soutput := &PrimerDesign_FWDSOutput{}
	output := _PrimerDesign_FWDRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PrimerDesign_FWDNew() interface{} {
	return &PrimerDesign_FWDElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PrimerDesign_FWDInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PrimerDesign_FWDRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PrimerDesign_FWDInput{},
			Out: &PrimerDesign_FWDOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PrimerDesign_FWDElement struct {
	inject.CheckedRunner
}

type PrimerDesign_FWDInput struct {
	DNASeqfile                               string
	FlankTargetSequence                      bool
	Maxgc                                    float64
	Maxlength                                int
	Maxtemp                                  wunit.Temperature
	Minlength                                int
	Mintemp                                  wunit.Temperature
	PermittednucleotideOverlapBetweenPrimers int
	Plasmid                                  bool
	RegionSequenceString                     string
	Seqstoavoid                              []string
}

type PrimerDesign_FWDOutput struct {
	FWDPrimer oligos.Primer
	Warnings  error
}

type PrimerDesign_FWDSOutput struct {
	Data struct {
		FWDPrimer oligos.Primer
		Warnings  error
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PrimerDesign_FWD",
		Constructor: PrimerDesign_FWDNew,
		Desc: component.ComponentDesc{
			Desc: "This element will design a sequencing primer to target amplification of a target region within a sequence file\nDesign criteria such as maximum gc content, acceptable ranges of melting temperatures and primer length may be specified by the user.\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/DNA/PrimerDesign/PrimerDesign_FWD.an",
			Params: []component.ParamDesc{
				{Name: "DNASeqfile", Desc: "genbank file (.gb)\n", Kind: "Parameters"},
				{Name: "FlankTargetSequence", Desc: "", Kind: "Parameters"},
				{Name: "Maxgc", Desc: "as a proportion of 1, i.e. 1 == 100%\n", Kind: "Parameters"},
				{Name: "Maxlength", Desc: "", Kind: "Parameters"},
				{Name: "Maxtemp", Desc: "", Kind: "Parameters"},
				{Name: "Minlength", Desc: "", Kind: "Parameters"},
				{Name: "Mintemp", Desc: "", Kind: "Parameters"},
				{Name: "PermittednucleotideOverlapBetweenPrimers", Desc: "number of nucleotides which primers can overlap by\n", Kind: "Parameters"},
				{Name: "Plasmid", Desc: "", Kind: "Parameters"},
				{Name: "RegionSequenceString", Desc: "", Kind: "Parameters"},
				{Name: "Seqstoavoid", Desc: "", Kind: "Parameters"},
				{Name: "FWDPrimer", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
