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
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"

// Input parameters for this protocol

// as a proportion of 1, i.e. 1 == 100%

// number of nucleotides which primers can overlap by

// Data which is returned from this protocol

// this needs to be changed to PrimerPair [2]oligo.Primer

// Physical inputs to this protocol

// Physical outputs from this protocol

func _PrimerDesign_ColonyPCR_wtypeRequirements() {

}

// Actions to perform before protocol itself
func _PrimerDesign_ColonyPCR_wtypeSetup(_ctx context.Context, _input *PrimerDesign_ColonyPCR_wtypeInput) {

}

// Core process of the protocol: steps to be performed for each input
func _PrimerDesign_ColonyPCR_wtypeSteps(_ctx context.Context, _input *PrimerDesign_ColonyPCR_wtypeInput, _output *PrimerDesign_ColonyPCR_wtypeOutput) {

	var region wtype.DNASequence = _input.RegionSequence

	regionstart, regionend, err := oligos.FindPositioninSequence(_input.FullDNASeq, _input.RegionSequence)

	if err != nil {
		fmt.Println("FindPositioninoligoFail with ", _input.FullDNASeq.Nm, _input.FullDNASeq.Seq, " and ", _input.RegionSequence.Nm, _input.RegionSequence.Seq)
		_output.Warnings = err
		execute.Errorf(_ctx, _output.Warnings.Error())
	}

	// if true then the start oint ot design primers is moved back 150bp to ensure full region is covered
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
		execute.Errorf(_ctx, _output.Warnings.Error())
	}

	fmt.Println(text.Print("FWDPrimer:", _output.FWDPrimer))

	if _input.FlankTargetSequence {
		region = oligos.DNAregion(_input.FullDNASeq, regionstart, regionend+150)
	} else {
		region = oligos.DNAregion(_input.FullDNASeq, regionstart, regionend)
	}

	_input.Seqstoavoid = append(_input.Seqstoavoid, _output.FWDPrimer.Seq)

	_output.REVPrimer, _output.Warnings = oligos.REVOligoSeq(region, _input.Maxgc, _input.Minlength, _input.Maxlength, _input.Mintemp, _input.Maxtemp, _input.Seqstoavoid, _input.PermittednucleotideOverlapBetweenPrimers)

	if _output.Warnings != nil {
		fmt.Println("REVOligoSeqFail")
		execute.Errorf(_ctx, _output.Warnings.Error())
	}

	fmt.Println(text.Print("REVPrimer:", _output.REVPrimer))

}

// Actions to perform after steps block to analyze data
func _PrimerDesign_ColonyPCR_wtypeAnalysis(_ctx context.Context, _input *PrimerDesign_ColonyPCR_wtypeInput, _output *PrimerDesign_ColonyPCR_wtypeOutput) {

}

func _PrimerDesign_ColonyPCR_wtypeValidation(_ctx context.Context, _input *PrimerDesign_ColonyPCR_wtypeInput, _output *PrimerDesign_ColonyPCR_wtypeOutput) {

}
func _PrimerDesign_ColonyPCR_wtypeRun(_ctx context.Context, input *PrimerDesign_ColonyPCR_wtypeInput) *PrimerDesign_ColonyPCR_wtypeOutput {
	output := &PrimerDesign_ColonyPCR_wtypeOutput{}
	_PrimerDesign_ColonyPCR_wtypeSetup(_ctx, input)
	_PrimerDesign_ColonyPCR_wtypeSteps(_ctx, input, output)
	_PrimerDesign_ColonyPCR_wtypeAnalysis(_ctx, input, output)
	_PrimerDesign_ColonyPCR_wtypeValidation(_ctx, input, output)
	return output
}

func PrimerDesign_ColonyPCR_wtypeRunSteps(_ctx context.Context, input *PrimerDesign_ColonyPCR_wtypeInput) *PrimerDesign_ColonyPCR_wtypeSOutput {
	soutput := &PrimerDesign_ColonyPCR_wtypeSOutput{}
	output := _PrimerDesign_ColonyPCR_wtypeRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PrimerDesign_ColonyPCR_wtypeNew() interface{} {
	return &PrimerDesign_ColonyPCR_wtypeElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PrimerDesign_ColonyPCR_wtypeInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PrimerDesign_ColonyPCR_wtypeRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PrimerDesign_ColonyPCR_wtypeInput{},
			Out: &PrimerDesign_ColonyPCR_wtypeOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PrimerDesign_ColonyPCR_wtypeElement struct {
	inject.CheckedRunner
}

type PrimerDesign_ColonyPCR_wtypeInput struct {
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

type PrimerDesign_ColonyPCR_wtypeOutput struct {
	FWDPrimer oligos.Primer
	REVPrimer oligos.Primer
	Warnings  error
}

type PrimerDesign_ColonyPCR_wtypeSOutput struct {
	Data struct {
		FWDPrimer oligos.Primer
		REVPrimer oligos.Primer
		Warnings  error
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PrimerDesign_ColonyPCR_wtype",
		Constructor: PrimerDesign_ColonyPCR_wtypeNew,
		Desc: component.ComponentDesc{
			Desc: "This element will design a pair of primers to cover a specified region of a sequence for colonyPCR.\nBut it's not finished yet!, please finish it off by designing the reverse primer\ngo to cd $GOPATH/src/github.com/antha-lang/antha/antha/examples/workflows/AnthaAcademy/Exercises/PrimerDesignExercise\nmake antharun return correct primerpairs for the three cases shown\nDesign criteria such as maximum gc content, acceptable ranges of melting temperatures and primer length may be specified by the user.\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/DNA/PrimerDesign/PrimerDesign_ColonyPCR_wtype.an",
			Params: []component.ParamDesc{
				{Name: "FlankTargetSequence", Desc: "", Kind: "Parameters"},
				{Name: "FullDNASeq", Desc: "", Kind: "Parameters"},
				{Name: "Maxgc", Desc: "as a proportion of 1, i.e. 1 == 100%\n", Kind: "Parameters"},
				{Name: "Maxlength", Desc: "", Kind: "Parameters"},
				{Name: "Maxtemp", Desc: "", Kind: "Parameters"},
				{Name: "Minlength", Desc: "", Kind: "Parameters"},
				{Name: "Mintemp", Desc: "", Kind: "Parameters"},
				{Name: "PermittednucleotideOverlapBetweenPrimers", Desc: "number of nucleotides which primers can overlap by\n", Kind: "Parameters"},
				{Name: "RegionSequence", Desc: "", Kind: "Parameters"},
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
