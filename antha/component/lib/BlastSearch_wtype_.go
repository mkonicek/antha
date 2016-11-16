// Example element demonstrating how to perform a BLAST search using the megablast algorithm

package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/blast"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	biogo "github.com/biogo/ncbi/blast"
	"golang.org/x/net/context"
)

// Input parameters for this protocol

// Data which is returned from this protocol; output data

//AllHits []biogo.Hit

// Physical inputs to this protocol

// Physical outputs from this protocol

func _BlastSearch_wtypeRequirements() {

}

// Actions to perform before protocol itself
func _BlastSearch_wtypeSetup(_ctx context.Context, _input *BlastSearch_wtypeInput) {

}

// Core process of the protocol: steps to be performed for each input
func _BlastSearch_wtypeSteps(_ctx context.Context, _input *BlastSearch_wtypeInput, _output *BlastSearch_wtypeOutput) {

	var err error
	var hits []biogo.Hit
	var hitsummary string
	var identity float64
	var coverage float64
	var besthitsummary string

	_output.AnthaSeq = _input.DNA

	// look for orfs
	orf, orftrue := sequences.FindORF(_output.AnthaSeq.Seq)

	if orftrue == true && len(orf.DNASeq) == len(_output.AnthaSeq.Seq) {
		// if open reading frame is detected, we'll perform a blastP search'
		fmt.Println("ORF detected:", "full sequence length: ", len(_output.AnthaSeq.Seq), "ORF length: ", len(orf.DNASeq))
		hits, err = blast.MegaBlastP(orf.ProtSeq)
	} else {
		// otherwise we'll blast the nucleotide sequence
		hits, err = _output.AnthaSeq.Blast()
	}
	if err != nil {
		fmt.Println(err.Error())

	}

	_output.ExactHits, hitsummary, err = blast.AllExactMatches(hits)

	if len(_output.ExactHits) == 0 {
		hitsummary, err = blast.HitSummary(hits, 10, 10)
	}
	_output.BestHit, identity, coverage, besthitsummary, err = blast.FindBestHit(hits)

	//	AllHits = hits
	_output.Hitssummary = hitsummary
	fmt.Println(hitsummary)
	fmt.Println(besthitsummary)
	// Rename Sequence with ID of top blast hit

	if coverage == 100 && identity == 100 {
		_output.AnthaSeq.Nm = _output.BestHit.Id
	}
	_output.Warning = err
	_output.Identity = identity
	_output.Coverage = coverage

}

// Actions to perform after steps block to analyze data
func _BlastSearch_wtypeAnalysis(_ctx context.Context, _input *BlastSearch_wtypeInput, _output *BlastSearch_wtypeOutput) {

}

func _BlastSearch_wtypeValidation(_ctx context.Context, _input *BlastSearch_wtypeInput, _output *BlastSearch_wtypeOutput) {

}
func _BlastSearch_wtypeRun(_ctx context.Context, input *BlastSearch_wtypeInput) *BlastSearch_wtypeOutput {
	output := &BlastSearch_wtypeOutput{}
	_BlastSearch_wtypeSetup(_ctx, input)
	_BlastSearch_wtypeSteps(_ctx, input, output)
	_BlastSearch_wtypeAnalysis(_ctx, input, output)
	_BlastSearch_wtypeValidation(_ctx, input, output)
	return output
}

func BlastSearch_wtypeRunSteps(_ctx context.Context, input *BlastSearch_wtypeInput) *BlastSearch_wtypeSOutput {
	soutput := &BlastSearch_wtypeSOutput{}
	output := _BlastSearch_wtypeRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func BlastSearch_wtypeNew() interface{} {
	return &BlastSearch_wtypeElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &BlastSearch_wtypeInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _BlastSearch_wtypeRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &BlastSearch_wtypeInput{},
			Out: &BlastSearch_wtypeOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type BlastSearch_wtypeElement struct {
	inject.CheckedRunner
}

type BlastSearch_wtypeInput struct {
	DNA wtype.DNASequence
}

type BlastSearch_wtypeOutput struct {
	AnthaSeq    wtype.DNASequence
	BestHit     biogo.Hit
	Coverage    float64
	ExactHits   []biogo.Hit
	Hitssummary string
	Identity    float64
	Warning     error
}

type BlastSearch_wtypeSOutput struct {
	Data struct {
		AnthaSeq    wtype.DNASequence
		BestHit     biogo.Hit
		Coverage    float64
		ExactHits   []biogo.Hit
		Hitssummary string
		Identity    float64
		Warning     error
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "BlastSearch_wtype",
		Constructor: BlastSearch_wtypeNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/DNA/BlastSearch/BlastSearch_wtype.an",
			Params: []component.ParamDesc{
				{Name: "DNA", Desc: "", Kind: "Parameters"},
				{Name: "AnthaSeq", Desc: "", Kind: "Data"},
				{Name: "BestHit", Desc: "", Kind: "Data"},
				{Name: "Coverage", Desc: "", Kind: "Data"},
				{Name: "ExactHits", Desc: "AllHits []biogo.Hit\n", Kind: "Data"},
				{Name: "Hitssummary", Desc: "", Kind: "Data"},
				{Name: "Identity", Desc: "", Kind: "Data"},
				{Name: "Warning", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
