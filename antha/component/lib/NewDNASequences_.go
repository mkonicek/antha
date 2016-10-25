// Demo protocol of how to create an array of dna types from parsing user inputs of various types
// scenarios handled:
// Biobrick IDS
// genbank files
// raw sequence
// inventory lookup
package lib

import (
	inventory "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Inventory"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/igem"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
	"strconv"
	"strings"
)

// Input parameters for this protocol

// Data which is returned from this protocol

// Physical inputs to this protocol

// Physical outputs from this protocol

func _NewDNASequencesRequirements() {

}

// Actions to perform before protocol itself
func _NewDNASequencesSetup(_ctx context.Context, _input *NewDNASequencesInput) {

}

// Core process of the protocol: steps to be performed for each input
func _NewDNASequencesSteps(_ctx context.Context, _input *NewDNASequencesInput, _output *NewDNASequencesOutput) {

	// set warnings reported back to user to none initially
	warnings := make([]string, 0)

	// make an empty array of DNA Sequences ready to fill
	partsinorder := make([]wtype.DNASequence, 0)

	var partDNA wtype.DNASequence
	var err error

	_output.Status = "all parts available"
	for i, part := range _input.Seqsinorder {
		// check if genbank feature
		if strings.Contains(part, ".gb") && strings.Contains(part, "Feature:") {

			split := strings.SplitAfter(part, ".gb")
			file := split[0]

			split2 := strings.Split(split[1], "Feature:")
			feature := split2[1]

			partDNA, err = parser.GenbankFeaturetoDNASequence(file, feature)

			if err != nil {
				execute.Errorf(_ctx, err.Error())
			}

			// check if genbank file
		} else if strings.Contains(part, ".gb") {

			partDNA, _ = parser.GenbanktoAnnotatedSeq(part)
			//check if biobrick
		} else if strings.Contains(part, "BBa_") {
			nm := "Part " + strconv.Itoa(i) + "_" + part
			part = igem.GetSequence(part)

			if _input.Vectors {
				partDNA = wtype.MakePlasmidDNASequence(nm, part)
			} else {
				partDNA = wtype.MakeLinearDNASequence(nm, part)
			}
			// check if in inventory
		} else if inventoryDNA, found := inventory.Partslist[part]; found {
			partDNA = inventoryDNA

			// else treat as DNA sequence and check
		} else {

			if _input.Vectors {
				partDNA = wtype.MakePlasmidDNASequence("Part "+strconv.Itoa(i), part)
			} else {
				partDNA = wtype.MakeLinearDNASequence("Part "+strconv.Itoa(i), part)
			}

			// test for illegal nucleotides
			pass, illegals, _ := sequences.Illegalnucleotides(partDNA)

			if !pass {
				var newstatus = make([]string, 0)
				for _, illegal := range illegals {

					newstatus = append(newstatus, "part: "+partDNA.Nm+" "+partDNA.Seq+": contains illegalnucleotides:"+illegal.ToString())
				}
				warnings = append(warnings, strings.Join(newstatus, ""))
				execute.Errorf(_ctx, strings.Join(newstatus, ""))
			}
			if pass && _input.BlastSeqswithNoName {
				// run a blast search on the sequence to get the name
				blastsearch := BlastSearch_wtypeRunSteps(_ctx, &BlastSearch_wtypeInput{DNA: partDNA})
				partDNA.Nm = blastsearch.Data.AnthaSeq.Nm
			}

		}
		partsinorder = append(partsinorder, partDNA)
	}

	_output.Parts = partsinorder

	_output.Warnings = warnings

}

// Actions to perform after steps block to analyze data
func _NewDNASequencesAnalysis(_ctx context.Context, _input *NewDNASequencesInput, _output *NewDNASequencesOutput) {

}

func _NewDNASequencesValidation(_ctx context.Context, _input *NewDNASequencesInput, _output *NewDNASequencesOutput) {

}
func _NewDNASequencesRun(_ctx context.Context, input *NewDNASequencesInput) *NewDNASequencesOutput {
	output := &NewDNASequencesOutput{}
	_NewDNASequencesSetup(_ctx, input)
	_NewDNASequencesSteps(_ctx, input, output)
	_NewDNASequencesAnalysis(_ctx, input, output)
	_NewDNASequencesValidation(_ctx, input, output)
	return output
}

func NewDNASequencesRunSteps(_ctx context.Context, input *NewDNASequencesInput) *NewDNASequencesSOutput {
	soutput := &NewDNASequencesSOutput{}
	output := _NewDNASequencesRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func NewDNASequencesNew() interface{} {
	return &NewDNASequencesElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &NewDNASequencesInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _NewDNASequencesRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &NewDNASequencesInput{},
			Out: &NewDNASequencesOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type NewDNASequencesElement struct {
	inject.CheckedRunner
}

type NewDNASequencesInput struct {
	BlastSeqswithNoName bool
	Seqsinorder         []string
	Vectors             bool
}

type NewDNASequencesOutput struct {
	Parts    []wtype.DNASequence
	Status   string
	Warnings []string
}

type NewDNASequencesSOutput struct {
	Data struct {
		Parts    []wtype.DNASequence
		Status   string
		Warnings []string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "NewDNASequences",
		Constructor: NewDNASequencesNew,
		Desc: component.ComponentDesc{
			Desc: "Demo protocol of how to create an array of dna types from parsing user inputs of various types\nscenarios handled:\nBiobrick IDS\ngenbank files\nraw sequence\ninventory lookup\n",
			Path: "antha/component/an/AnthaAcademy/Lesson6_DNA/D_NewDNASequences.an",
			Params: []component.ParamDesc{
				{Name: "BlastSeqswithNoName", Desc: "", Kind: "Parameters"},
				{Name: "Seqsinorder", Desc: "", Kind: "Parameters"},
				{Name: "Vectors", Desc: "", Kind: "Parameters"},
				{Name: "Parts", Desc: "", Kind: "Data"},
				{Name: "Status", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
