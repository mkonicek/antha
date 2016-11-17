package lib

import (
	"fmt"
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

func _JoinDNASequences_mapRequirements() {

}

// Actions to perform before protocol itself
func _JoinDNASequences_mapSetup(_ctx context.Context, _input *JoinDNASequences_mapInput) {

}

// Core process of the protocol: steps to be performed for each input
func _JoinDNASequences_mapSteps(_ctx context.Context, _input *JoinDNASequences_mapInput, _output *JoinDNASequences_mapOutput) {

	// make a lookup table of []DNASequences; in go these are called maps
	// in this case the map is called sequencepairs; the key is of type string; the values returned will be an array of DNASequences
	// e.g. we can add two sequences eGFP and sfGFP to the map like so:
	// sequencepairs["GFPSequences"] = wtype.DNASequence{eGFP,sfGFP}
	// we can return the two stored sequences to a variable called gfpseqscopy using the key "GFPSequences" by typing
	// gfpseqscopy := sequencepairs["GFPSequences"]

	var sequencepairs = make(map[string][]wtype.DNASequence)

	// Step 1. Call antha element for turning a string array of parts into []DNASequence for each map entry
	// We can range through maps in a similar fashion to arrays; the important difference is after for the key value is used rather than the index as is the case with array
	for key, values := range _input.Pairs {

		// this is how we call an antha element from within an element
		seqs := NewDNASequencesRunSteps(_ctx, &NewDNASequencesInput{Seqsinorder: values,
			BlastSeqswithNoName: _input.BlastSeqswithNoName,
			Vectors:             _input.Vectors},
		)

		sequencepairs[key] = seqs.Data.Parts
	}

	//  make an array of seqs to export for each map combination
	var seqstoexport = make([]wtype.DNASequence, 0)

	// Step 2. Range through the map created in step 1.
	for key, Seqsinorder := range sequencepairs {

		newSeq := Seqsinorder[0]
		fmt.Println("seq?", Seqsinorder[0].Nm)
		//seqnames := make([]string,0)

		// Step 2a. Each set of sequences we'll range through and concatenate the sequence with the next sequence
		for i, seq := range Seqsinorder {
			fmt.Println("seq[i]?", Seqsinorder[i].Nm)
			if i != 0 {
				newSeq.Append(seq.Seq)
			}
			//seqnames = append(seqnames,seq.Nm)
		}

		// Step 2b. Name the new DNAParts using map key. This could also name by concatenating but we'll use key for now
		newSeq.Nm = key //strings.Join(seqnames,"_")
		seqstoexport = append(seqstoexport, newSeq)
	}

	_output.JoinedSeqs = seqstoexport

}

// Actions to perform after steps block to analyze data
func _JoinDNASequences_mapAnalysis(_ctx context.Context, _input *JoinDNASequences_mapInput, _output *JoinDNASequences_mapOutput) {

}

func _JoinDNASequences_mapValidation(_ctx context.Context, _input *JoinDNASequences_mapInput, _output *JoinDNASequences_mapOutput) {

}
func _JoinDNASequences_mapRun(_ctx context.Context, input *JoinDNASequences_mapInput) *JoinDNASequences_mapOutput {
	output := &JoinDNASequences_mapOutput{}
	_JoinDNASequences_mapSetup(_ctx, input)
	_JoinDNASequences_mapSteps(_ctx, input, output)
	_JoinDNASequences_mapAnalysis(_ctx, input, output)
	_JoinDNASequences_mapValidation(_ctx, input, output)
	return output
}

func JoinDNASequences_mapRunSteps(_ctx context.Context, input *JoinDNASequences_mapInput) *JoinDNASequences_mapSOutput {
	soutput := &JoinDNASequences_mapSOutput{}
	output := _JoinDNASequences_mapRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func JoinDNASequences_mapNew() interface{} {
	return &JoinDNASequences_mapElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &JoinDNASequences_mapInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _JoinDNASequences_mapRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &JoinDNASequences_mapInput{},
			Out: &JoinDNASequences_mapOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type JoinDNASequences_mapElement struct {
	inject.CheckedRunner
}

type JoinDNASequences_mapInput struct {
	BlastSeqswithNoName bool
	Pairs               map[string][]string
	Vectors             bool
}

type JoinDNASequences_mapOutput struct {
	JoinedSeqs []wtype.DNASequence
}

type JoinDNASequences_mapSOutput struct {
	Data struct {
		JoinedSeqs []wtype.DNASequence
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "JoinDNASequences_map",
		Constructor: JoinDNASequences_mapNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/AnthaAcademy/Lesson4_DNA/G_JoinSequencePairs.an",
			Params: []component.ParamDesc{
				{Name: "BlastSeqswithNoName", Desc: "", Kind: "Parameters"},
				{Name: "Pairs", Desc: "", Kind: "Parameters"},
				{Name: "Vectors", Desc: "", Kind: "Parameters"},
				{Name: "JoinedSeqs", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
