// This element will design primers to cover a specified region of a sequence at the interval specified by the user (e.g. every 800 bp).
// Design criteria such as maximum gc content, acceptable ranges of melting temperatures and primer length may be specified by the user.
package lib

import (
	//"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/oligos"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
	"strings"
)

// Input parameters for this protocol

// genbank file (.gb)
//Plasmid 	bool

// as a proportion of 1, i.e. 1 == 100%

// number of nucleotides which primers can overlap by

// permissable values: "byFeaturename", "byPositions", "bySequence"

// Data which is returned from this protocol

// Physical inputs to this protocol

// Physical outputs from this protocol

func _PrimerDesign_coverRegionRequirements() {

}

// Actions to perform before protocol itself
func _PrimerDesign_coverRegionSetup(_ctx context.Context, _input *PrimerDesign_coverRegionInput) {

}

// Core process of the protocol: steps to be performed for each input
func _PrimerDesign_coverRegionSteps(_ctx context.Context, _input *PrimerDesign_coverRegionInput, _output *PrimerDesign_coverRegionOutput) {
	var plasmid wtype.DNASequence
	var allprimers []oligos.Primer

	plasmid, _ = parser.GenbanktoAnnotatedSeq(_input.DNASeqfile)

	/*if len(plasmids)>0 {
	plasmid = plasmids[0]
	}
	if len(plasmids)>1{
		Warnings = fmt.Errorf("Warning! more than one sequence in sequence file! Only used first sequence for primer design")
	}*/

	if strings.Contains(strings.ToUpper(_input.Method), "POSITIONS") {
		allprimers = oligos.DesignFWDPRimerstoCoverRegion(plasmid, _input.RegionStart, _input.RegionEnd, _input.PrimereveryXnucleotides, _input.Maxgc, _input.Minlength, _input.Maxlength, _input.Mintemp, _input.Maxtemp, _input.Seqstoavoid, _input.PermittednucleotideOverlapBetweenPrimers)
	} else if strings.Contains(strings.ToUpper(_input.Method), "NAME") {
		allprimers = oligos.DesignFWDPRimerstoCoverFeature(plasmid, _input.RegionName, _input.PrimereveryXnucleotides, _input.Maxgc, _input.Minlength, _input.Maxlength, _input.Mintemp, _input.Maxtemp, _input.Seqstoavoid, _input.PermittednucleotideOverlapBetweenPrimers)

	} else if strings.Contains(strings.ToUpper(_input.Method), "SEQUENCE") {
		allprimers = oligos.DesignFWDPRimerstoCoverSequence(plasmid, _input.RegionSequence, _input.PrimereveryXnucleotides, _input.Maxgc, _input.Minlength, _input.Maxlength, _input.Mintemp, _input.Maxtemp, _input.Seqstoavoid, _input.PermittednucleotideOverlapBetweenPrimers)

	}
	_output.AllPrimers = allprimers

}

// Actions to perform after steps block to analyze data
func _PrimerDesign_coverRegionAnalysis(_ctx context.Context, _input *PrimerDesign_coverRegionInput, _output *PrimerDesign_coverRegionOutput) {

}

func _PrimerDesign_coverRegionValidation(_ctx context.Context, _input *PrimerDesign_coverRegionInput, _output *PrimerDesign_coverRegionOutput) {

}
func _PrimerDesign_coverRegionRun(_ctx context.Context, input *PrimerDesign_coverRegionInput) *PrimerDesign_coverRegionOutput {
	output := &PrimerDesign_coverRegionOutput{}
	_PrimerDesign_coverRegionSetup(_ctx, input)
	_PrimerDesign_coverRegionSteps(_ctx, input, output)
	_PrimerDesign_coverRegionAnalysis(_ctx, input, output)
	_PrimerDesign_coverRegionValidation(_ctx, input, output)
	return output
}

func PrimerDesign_coverRegionRunSteps(_ctx context.Context, input *PrimerDesign_coverRegionInput) *PrimerDesign_coverRegionSOutput {
	soutput := &PrimerDesign_coverRegionSOutput{}
	output := _PrimerDesign_coverRegionRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PrimerDesign_coverRegionNew() interface{} {
	return &PrimerDesign_coverRegionElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PrimerDesign_coverRegionInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PrimerDesign_coverRegionRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PrimerDesign_coverRegionInput{},
			Out: &PrimerDesign_coverRegionOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PrimerDesign_coverRegionElement struct {
	inject.CheckedRunner
}

type PrimerDesign_coverRegionInput struct {
	DNASeqfile                               string
	Maxgc                                    float64
	Maxlength                                int
	Maxtemp                                  wunit.Temperature
	Method                                   string
	Minlength                                int
	Mintemp                                  wunit.Temperature
	PermittednucleotideOverlapBetweenPrimers int
	PrimereveryXnucleotides                  int
	RegionEnd                                int
	RegionName                               string
	RegionSequence                           string
	RegionStart                              int
	Seqstoavoid                              []string
}

type PrimerDesign_coverRegionOutput struct {
	AllPrimers []oligos.Primer
	Warnings   error
}

type PrimerDesign_coverRegionSOutput struct {
	Data struct {
		AllPrimers []oligos.Primer
		Warnings   error
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PrimerDesign_coverRegion",
		Constructor: PrimerDesign_coverRegionNew,
		Desc: component.ComponentDesc{
			Desc: "This element will design primers to cover a specified region of a sequence at the interval specified by the user (e.g. every 800 bp).\nDesign criteria such as maximum gc content, acceptable ranges of melting temperatures and primer length may be specified by the user.\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/DNA/PrimerDesign/PrimerDesign_coverRegion.an",
			Params: []component.ParamDesc{
				{Name: "DNASeqfile", Desc: "genbank file (.gb)\n", Kind: "Parameters"},
				{Name: "Maxgc", Desc: "as a proportion of 1, i.e. 1 == 100%\n", Kind: "Parameters"},
				{Name: "Maxlength", Desc: "", Kind: "Parameters"},
				{Name: "Maxtemp", Desc: "Plasmid \tbool\n", Kind: "Parameters"},
				{Name: "Method", Desc: "permissable values: \"byFeaturename\", \"byPositions\", \"bySequence\"\n", Kind: "Parameters"},
				{Name: "Minlength", Desc: "", Kind: "Parameters"},
				{Name: "Mintemp", Desc: "", Kind: "Parameters"},
				{Name: "PermittednucleotideOverlapBetweenPrimers", Desc: "number of nucleotides which primers can overlap by\n", Kind: "Parameters"},
				{Name: "PrimereveryXnucleotides", Desc: "", Kind: "Parameters"},
				{Name: "RegionEnd", Desc: "", Kind: "Parameters"},
				{Name: "RegionName", Desc: "", Kind: "Parameters"},
				{Name: "RegionSequence", Desc: "", Kind: "Parameters"},
				{Name: "RegionStart", Desc: "", Kind: "Parameters"},
				{Name: "Seqstoavoid", Desc: "", Kind: "Parameters"},
				{Name: "AllPrimers", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
