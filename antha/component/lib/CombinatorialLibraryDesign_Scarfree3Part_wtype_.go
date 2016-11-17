// This protocol is intended to design assembly parts using a specified enzyme.
// overhangs are added to complement the adjacent parts and leave no scar.
// parts can be entered as genbank (.gb) files, sequences or biobrick IDs
// If assembly simulation fails after overhangs are added. In order to help the user
// diagnose the reason, a report of the part overhangs
// is returned to the user along with a list of cut sites in each part.

package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/export"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/oligos"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
	"path/filepath"
	"strconv"
)

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// parts to order

// desired sequence to end up with after assembly

// Input Requirement specification
func _CombinatorialLibraryDesign_Scarfree3Part_wtypeRequirements() {
	// e.g. are MoClo types valid?
}

// Conditions to run on startup
func _CombinatorialLibraryDesign_Scarfree3Part_wtypeSetup(_ctx context.Context, _input *CombinatorialLibraryDesign_Scarfree3Part_wtypeInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _CombinatorialLibraryDesign_Scarfree3Part_wtypeSteps(_ctx context.Context, _input *CombinatorialLibraryDesign_Scarfree3Part_wtypeInput, _output *CombinatorialLibraryDesign_Scarfree3Part_wtypeOutput) {
	_output.StatusMap = make(map[string]string)
	_output.PartswithOverhangsMap = make(map[string][]wtype.DNASequence) // parts to order
	_output.PassMap = make(map[string]bool)
	_output.SeqsMap = make(map[string]wtype.DNASequence) // desired sequence to end up with after assembly
	_output.Sequences = make([]wtype.DNASequence, 0)
	_output.Parts = make([][]wtype.DNASequence, 0)
	_output.EndreportMap = make(map[string]string)
	_output.PositionReportMap = make(map[string][]string)
	_output.StatusMap = make(map[string]string)
	_output.PrimerMap = make(map[string]oligos.Primer)

	var counter int = 1

	for j := range _input.Vectors {
		for k := range _input.Part1s {
			for l := range _input.Part2s {
				for m := range _input.Part3s {

					key := _input.ProjectName + "_" + "Vector" + strconv.Itoa(j+1) + "_1." + strconv.Itoa(k+1) + "_2." + strconv.Itoa(l+1) + "_3." + strconv.Itoa(m+1)
					assembly := Scarfree_siteremove_orfcheck_wtypeRunSteps(_ctx, &Scarfree_siteremove_orfcheck_wtypeInput{Constructname: key,
						Seqsinorder: []wtype.DNASequence{_input.Part1s[k], _input.Part2s[l], _input.Part3s[m]},
						Enzymename:  _input.EnzymeName,
						Vector:      _input.Vectors[j],
						OtherEnzymeSitesToRemove:      _input.SitesToRemove,
						ORFstoConfirm:                 _input.ORFStoconfirm, // enter each as amino acid sequence
						RemoveproblemRestrictionSites: _input.RemoveproblemRestrictionSites,
						EndsAlreadyadded:              _input.EndsAlreadyadded,
						ExporttoFastaFile:             _input.FolderPerConstruct,
						BlastSeqswithNoName:           _input.BlastSearchSeqs},
					)
					key = key                                                             //+ Vectors[j]
					_output.PartswithOverhangsMap[key] = assembly.Data.PartswithOverhangs // parts to order
					_output.PassMap[key] = assembly.Data.Simulationpass
					_output.EndreportMap[key] = assembly.Data.Endreport
					_output.PositionReportMap[key] = assembly.Data.PositionReport
					_output.SeqsMap[key] = assembly.Data.NewDNASequence
					_output.Sequences = append(_output.Sequences, assembly.Data.NewDNASequence)
					_output.Parts = append(_output.Parts, assembly.Data.PartswithOverhangs)
					_output.StatusMap[key] = assembly.Data.Status

					// for each vector we'll also design sequencing primers

					primer := PrimerDesign_FWD_wtypeRunSteps(_ctx, &PrimerDesign_FWD_wtypeInput{FullDNASeq: _input.Vectors[j], // design sequencing primers in original vector assembly.Data.NewDNASequence,
						Maxtemp:                                  wunit.NewTemperature(60, "C"),
						Mintemp:                                  wunit.NewTemperature(50, "C"),
						Maxgc:                                    0.6,
						Minlength:                                15,
						Maxlength:                                30,
						Seqstoavoid:                              []string{},
						PermittednucleotideOverlapBetweenPrimers: 10,                                                                                                                                                                                                                                                                                                                                                                                                                                                        // number of nucleotides which primers can overlap by
						RegionSequence:                           wtype.MakeLinearDNASequence("PartofSacBPromoter", "gatattatgatattttctgaattgtgattaaaaaggcaactttatgcccatgcaacagaaactataaaaaatacagagaatgaaaagaaacagatagattttttagttctttaggcccgtagtctgcaaatccttttatgattttctatcaaacaaaagaggaaaatagaccagttgcaatccaaacgagagtctaatagaatgaggtcgaaaagtaaatcgcgcgggtttgttactgataaagcaggcaagacctaaaatgtgtaaagggcaaagtgtatactttggcgtcaccccttacatattttaggtctttttttattgtgcgtaactaacttgccatcttcaaacaggagggctggaagaagcagaccgctaacacagtaca"), // SacB promoter sequence in vector //PartsWithSitesRemoved[0], // first part
						FlankTargetSequence:                      true},
					)
					_output.PrimerMap[key] = primer.Data.FWDPrimer

					counter++

				}
			}
		}
		// export sequence to fasta
		if _input.FolderPerProject {

			// export simulated assembled sequences to file
			export.Makefastaserial2(export.LOCAL, filepath.Join(_input.ProjectName, "AssembledSequences"), _output.Sequences)

			// add fasta files of all parts with overhangs
			labels := []string{"Device1", "Device2", "Device3"}
			for j := range labels {
				for i := range _output.Parts {
					fmt.Println(i, len(labels), len(_output.Parts))
					export.Makefastaserial2(export.LOCAL, filepath.Join(_input.ProjectName, labels[j]), _output.Parts[i])
				}
			}
		}

	}
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _CombinatorialLibraryDesign_Scarfree3Part_wtypeAnalysis(_ctx context.Context, _input *CombinatorialLibraryDesign_Scarfree3Part_wtypeInput, _output *CombinatorialLibraryDesign_Scarfree3Part_wtypeOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _CombinatorialLibraryDesign_Scarfree3Part_wtypeValidation(_ctx context.Context, _input *CombinatorialLibraryDesign_Scarfree3Part_wtypeInput, _output *CombinatorialLibraryDesign_Scarfree3Part_wtypeOutput) {
}
func _CombinatorialLibraryDesign_Scarfree3Part_wtypeRun(_ctx context.Context, input *CombinatorialLibraryDesign_Scarfree3Part_wtypeInput) *CombinatorialLibraryDesign_Scarfree3Part_wtypeOutput {
	output := &CombinatorialLibraryDesign_Scarfree3Part_wtypeOutput{}
	_CombinatorialLibraryDesign_Scarfree3Part_wtypeSetup(_ctx, input)
	_CombinatorialLibraryDesign_Scarfree3Part_wtypeSteps(_ctx, input, output)
	_CombinatorialLibraryDesign_Scarfree3Part_wtypeAnalysis(_ctx, input, output)
	_CombinatorialLibraryDesign_Scarfree3Part_wtypeValidation(_ctx, input, output)
	return output
}

func CombinatorialLibraryDesign_Scarfree3Part_wtypeRunSteps(_ctx context.Context, input *CombinatorialLibraryDesign_Scarfree3Part_wtypeInput) *CombinatorialLibraryDesign_Scarfree3Part_wtypeSOutput {
	soutput := &CombinatorialLibraryDesign_Scarfree3Part_wtypeSOutput{}
	output := _CombinatorialLibraryDesign_Scarfree3Part_wtypeRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func CombinatorialLibraryDesign_Scarfree3Part_wtypeNew() interface{} {
	return &CombinatorialLibraryDesign_Scarfree3Part_wtypeElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &CombinatorialLibraryDesign_Scarfree3Part_wtypeInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _CombinatorialLibraryDesign_Scarfree3Part_wtypeRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &CombinatorialLibraryDesign_Scarfree3Part_wtypeInput{},
			Out: &CombinatorialLibraryDesign_Scarfree3Part_wtypeOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type CombinatorialLibraryDesign_Scarfree3Part_wtypeElement struct {
	inject.CheckedRunner
}

type CombinatorialLibraryDesign_Scarfree3Part_wtypeInput struct {
	BlastSearchSeqs               bool
	EndsAlreadyadded              bool
	EnzymeName                    string
	FolderPerConstruct            bool
	FolderPerProject              bool
	ORFStoconfirm                 []string
	Part1s                        []wtype.DNASequence
	Part2s                        []wtype.DNASequence
	Part3s                        []wtype.DNASequence
	ProjectName                   string
	RemoveproblemRestrictionSites bool
	SitesToRemove                 []string
	Vectors                       []wtype.DNASequence
}

type CombinatorialLibraryDesign_Scarfree3Part_wtypeOutput struct {
	EndreportMap          map[string]string
	Parts                 [][]wtype.DNASequence
	PartswithOverhangsMap map[string][]wtype.DNASequence
	PassMap               map[string]bool
	PositionReportMap     map[string][]string
	PrimerMap             map[string]oligos.Primer
	SeqsMap               map[string]wtype.DNASequence
	Sequences             []wtype.DNASequence
	StatusMap             map[string]string
}

type CombinatorialLibraryDesign_Scarfree3Part_wtypeSOutput struct {
	Data struct {
		EndreportMap          map[string]string
		Parts                 [][]wtype.DNASequence
		PartswithOverhangsMap map[string][]wtype.DNASequence
		PassMap               map[string]bool
		PositionReportMap     map[string][]string
		PrimerMap             map[string]oligos.Primer
		SeqsMap               map[string]wtype.DNASequence
		Sequences             []wtype.DNASequence
		StatusMap             map[string]string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "CombinatorialLibraryDesign_Scarfree3Part_wtype",
		Constructor: CombinatorialLibraryDesign_Scarfree3Part_wtypeNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/DNA/TypeIISAssembly_design/CombinatorialLibraryDesign_Scarfree_wtype.an",
			Params: []component.ParamDesc{
				{Name: "BlastSearchSeqs", Desc: "", Kind: "Parameters"},
				{Name: "EndsAlreadyadded", Desc: "", Kind: "Parameters"},
				{Name: "EnzymeName", Desc: "", Kind: "Parameters"},
				{Name: "FolderPerConstruct", Desc: "", Kind: "Parameters"},
				{Name: "FolderPerProject", Desc: "", Kind: "Parameters"},
				{Name: "ORFStoconfirm", Desc: "", Kind: "Parameters"},
				{Name: "Part1s", Desc: "", Kind: "Parameters"},
				{Name: "Part2s", Desc: "", Kind: "Parameters"},
				{Name: "Part3s", Desc: "", Kind: "Parameters"},
				{Name: "ProjectName", Desc: "", Kind: "Parameters"},
				{Name: "RemoveproblemRestrictionSites", Desc: "", Kind: "Parameters"},
				{Name: "SitesToRemove", Desc: "", Kind: "Parameters"},
				{Name: "Vectors", Desc: "", Kind: "Parameters"},
				{Name: "EndreportMap", Desc: "", Kind: "Data"},
				{Name: "Parts", Desc: "", Kind: "Data"},
				{Name: "PartswithOverhangsMap", Desc: "parts to order\n", Kind: "Data"},
				{Name: "PassMap", Desc: "", Kind: "Data"},
				{Name: "PositionReportMap", Desc: "", Kind: "Data"},
				{Name: "PrimerMap", Desc: "", Kind: "Data"},
				{Name: "SeqsMap", Desc: "desired sequence to end up with after assembly\n", Kind: "Data"},
				{Name: "Sequences", Desc: "", Kind: "Data"},
				{Name: "StatusMap", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
