// This element will design a pair of primers to cover a specified region of a sequence for colonyPCR.
// But it's not finished yet!, please finish it off by designing the reverse primer
// go to cd $GOPATH/src/github.com/antha-lang/antha/antha/examples/workflows/AnthaAcademy/Exercises/PrimerDesignExercise
// make antharun return correct primerpairs for the three cases shown
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

// this needs to be changed to PrimerPair [2]oligo.Primer

// Physical inputs to this protocol

// Physical outputs from this protocol

func _PrimerDesign_ColonyPCRRequirements() {

}

// Actions to perform before protocol itself
func _PrimerDesign_ColonyPCRSetup(_ctx context.Context, _input *PrimerDesign_ColonyPCRInput) {

}

// Core process of the protocol: steps to be performed for each input
func _PrimerDesign_ColonyPCRSteps(_ctx context.Context, _input *PrimerDesign_ColonyPCRInput, _output *PrimerDesign_ColonyPCROutput) {

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

	if _input.FlankTargetSequence {
		region = oligos.DNAregion(fulldnaseqs[0], regionstart, regionend+150)
	} else {
		region = oligos.DNAregion(fulldnaseqs[0], regionstart, regionend)
	}
	fmt.Println("6")

	_input.Seqstoavoid = append(_input.Seqstoavoid, _output.FWDPrimer.Seq)

	_output.REVPrimer, _output.Warnings = oligos.REVOligoSeq(region, _input.Maxgc, _input.Minlength, _input.Maxlength, _input.Mintemp, _input.Maxtemp, _input.Seqstoavoid, _input.PermittednucleotideOverlapBetweenPrimers)

	if _output.Warnings != nil {
		fmt.Println("REVOligoSeqFail")
		execute.Errorf(_ctx, _output.Warnings.Error())
	}

	fmt.Println(text.Print("REVPrimer:", _output.REVPrimer))

}

// Actions to perform after steps block to analyze data
func _PrimerDesign_ColonyPCRAnalysis(_ctx context.Context, _input *PrimerDesign_ColonyPCRInput, _output *PrimerDesign_ColonyPCROutput) {

}

func _PrimerDesign_ColonyPCRValidation(_ctx context.Context, _input *PrimerDesign_ColonyPCRInput, _output *PrimerDesign_ColonyPCROutput) {

}
func _PrimerDesign_ColonyPCRRun(_ctx context.Context, input *PrimerDesign_ColonyPCRInput) *PrimerDesign_ColonyPCROutput {
	output := &PrimerDesign_ColonyPCROutput{}
	_PrimerDesign_ColonyPCRSetup(_ctx, input)
	_PrimerDesign_ColonyPCRSteps(_ctx, input, output)
	_PrimerDesign_ColonyPCRAnalysis(_ctx, input, output)
	_PrimerDesign_ColonyPCRValidation(_ctx, input, output)
	return output
}

func PrimerDesign_ColonyPCRRunSteps(_ctx context.Context, input *PrimerDesign_ColonyPCRInput) *PrimerDesign_ColonyPCRSOutput {
	soutput := &PrimerDesign_ColonyPCRSOutput{}
	output := _PrimerDesign_ColonyPCRRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PrimerDesign_ColonyPCRNew() interface{} {
	return &PrimerDesign_ColonyPCRElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PrimerDesign_ColonyPCRInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PrimerDesign_ColonyPCRRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PrimerDesign_ColonyPCRInput{},
			Out: &PrimerDesign_ColonyPCROutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PrimerDesign_ColonyPCRElement struct {
	inject.CheckedRunner
}

type PrimerDesign_ColonyPCRInput struct {
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

type PrimerDesign_ColonyPCROutput struct {
	FWDPrimer oligos.Primer
	REVPrimer oligos.Primer
	Warnings  error
}

type PrimerDesign_ColonyPCRSOutput struct {
	Data struct {
		FWDPrimer oligos.Primer
		REVPrimer oligos.Primer
		Warnings  error
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PrimerDesign_ColonyPCR",
		Constructor: PrimerDesign_ColonyPCRNew,
		Desc: component.ComponentDesc{
			Desc: "This element will design a pair of primers to cover a specified region of a sequence for colonyPCR.\nBut it's not finished yet!, please finish it off by designing the reverse primer\ngo to cd $GOPATH/src/github.com/antha-lang/antha/antha/examples/workflows/AnthaAcademy/Exercises/PrimerDesignExercise\nmake antharun return correct primerpairs for the three cases shown\nDesign criteria such as maximum gc content, acceptable ranges of melting temperatures and primer length may be specified by the user.\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/DNA/PrimerDesign/PrimerDesign_ColonyPCR.an",
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
				{Name: "FWDPrimer", Desc: "this needs to be changed to PrimerPair [2]oligo.Primer\n", Kind: "Data"},
				{Name: "REVPrimer", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
