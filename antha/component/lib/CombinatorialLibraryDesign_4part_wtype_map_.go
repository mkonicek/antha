// This protocol is intended to design assembly parts using a specified enzyme.
// overhangs are added to complement the adjacent parts and leave no scar.
// parts can be entered as genbank (.gb) files, sequences or biobrick IDs
// If assembly simulation fails after overhangs are added. In order to help the user
// diagnose the reason, a report of the part overhangs
// is returned to the user along with a list of cut sites in each part.
package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/export"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/oligos"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"path/filepath"
	//"strconv"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

//Seqsinorder					map[string][]string // constructname to sequence combination

//MoClo
// of assembly standard

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// parts to order

// desired sequence to end up with after assembly

// Input Requirement specification
func _CombinatorialLibraryDesign_4part_wtype_mapRequirements() {
	// e.g. are MoClo types valid?
}

// Conditions to run on startup
func _CombinatorialLibraryDesign_4part_wtype_mapSetup(_ctx context.Context, _input *CombinatorialLibraryDesign_4part_wtype_mapInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _CombinatorialLibraryDesign_4part_wtype_mapSteps(_ctx context.Context, _input *CombinatorialLibraryDesign_4part_wtype_mapInput, _output *CombinatorialLibraryDesign_4part_wtype_mapOutput) {
	_output.StatusMap = make(map[string]string)
	_output.PartswithOverhangsMap = make(map[string][]wtype.DNASequence) // parts to order
	_output.PassMap = make(map[string]bool)
	_output.SeqsMap = make(map[string]wtype.DNASequence) // desired sequence to end up with after assembly
	_output.EndreportMap = make(map[string]string)
	_output.PositionReportMap = make(map[string][]string)
	_output.StatusMap = make(map[string]string)
	_output.PrimerMap = make(map[string]oligos.Primer)

	_output.Sequences = make([]wtype.DNASequence, 0)
	_output.Parts = make([][]wtype.DNASequence, 0)
	_output.SequencingPrimers = make([][]wtype.DNASequence, 0)
	//uniqueFWDPrimers := make([]wtype.DNASequence,0)
	//uniqueREVPrimers := make([]wtype.DNASequence,0)

	var counter int = 1

	for j := range _input.Vectors {
		for k := range _input.PROs {
			for l := range _input.RBSs {
				for m := range _input.CDSs {
					for n := range _input.TERs {
						key := _input.ProjectName + _input.Vectors[j].Nm + "_" + _input.PROs[k].Nm + "_" + _input.RBSs[l].Nm + "_" + _input.CDSs[m].Nm
						assembly := AssemblyStandard_siteremove_orfcheck_wtypeRunSteps(_ctx, &AssemblyStandard_siteremove_orfcheck_wtypeInput{Constructname: key,
							Seqsinorder:                   []wtype.DNASequence{_input.PROs[k], _input.RBSs[l], _input.CDSs[m], _input.TERs[n]},
							AssemblyStandard:              _input.Standard,
							Level:                         _input.StandardLevel, // of assembly standard
							Vector:                        _input.Vectors[j],
							PartMoClotypesinorder:         []string{_input.StandardLabels[0], _input.StandardLabels[1], _input.StandardLabels[2], _input.StandardLabels[3]},
							OtherEnzymeSitesToRemove:      _input.SitesToRemove,
							ORFstoConfirm:                 []string{}, // enter each as amino acid sequence
							RemoveproblemRestrictionSites: true,
							OnlyRemovesitesinORFs:         false,
							EndsAlreadyadded:              false,
							ExporttoFastaFile:             _input.FolderPerConstruct,
							BlastSeqswithNoName:           _input.BlastSearchSeqs},
						)
						key = key                                                             //+ Vectors[j]
						_output.PartswithOverhangsMap[key] = assembly.Data.PartswithOverhangs // parts to order
						_output.Parts = append(_output.Parts, assembly.Data.PartswithOverhangs)
						_output.PassMap[key] = assembly.Data.Simulationpass
						_output.EndreportMap[key] = assembly.Data.Endreport
						_output.PositionReportMap[key] = assembly.Data.PositionReport
						_output.SeqsMap[key] = assembly.Data.NewDNASequence
						_output.Sequences = append(_output.Sequences, assembly.Data.NewDNASequence)
						_output.StatusMap[key] = assembly.Data.Status

						// for each vector we'll also design sequencing primers

						primer := PrimerDesign_ColonyPCR_wtypeRunSteps(_ctx, &PrimerDesign_ColonyPCR_wtypeInput{FullDNASeq: assembly.Data.NewDNASequence,
							Maxtemp:                                  wunit.NewTemperature(72, "C"),
							Mintemp:                                  wunit.NewTemperature(50, "C"),
							Maxgc:                                    0.7,
							Minlength:                                12,
							Maxlength:                                30,
							Seqstoavoid:                              []string{},
							PermittednucleotideOverlapBetweenPrimers: 10,                                     // number of nucleotides which primers can overlap by
							RegionSequence:                           assembly.Data.PartsWithSitesRemoved[0], // first part
							FlankTargetSequence:                      true},
						)

						// rename primers
						primer.Data.FWDPrimer.Nm = primer.Data.FWDPrimer.Nm + _input.ProjectName + _input.Vectors[j].Nm + "_FWD"
						primer.Data.REVPrimer.Nm = primer.Data.REVPrimer.Nm + _input.ProjectName + _input.Vectors[j].Nm + "_REV"

						_output.PrimerMap[key+"_FWD"] = primer.Data.FWDPrimer
						_output.PrimerMap[key+"_REV"] = primer.Data.REVPrimer
						_output.SequencingPrimers = append(_output.SequencingPrimers, []wtype.DNASequence{primer.Data.FWDPrimer.DNASequence, primer.Data.REVPrimer.DNASequence})

						//uniqueFWDPrimers = append(uniqueFWDPrimers,primer.Data.FWDPrimer.Seq)
						//uniqueREVPrimers = append(uniqueREVPrimers,primer.Data.REVPrimer.Seq)

						counter++
					}
				}
			}
		}
	}

	//uniqueFWDPrimers = search.RemoveDuplicateSequences(uniqueFWDPrimers)
	//uniqueREVPrimers = search.RemoveDuplicateSequences(uniqueREVPrimers)

	// export sequence to fasta
	if _input.FolderPerProject {

		// export simulated sequences to file
		export.Makefastaserial2(export.LOCAL, filepath.Join(_input.ProjectName, "AssembledSequences"), _output.Sequences)

		// add fasta file for each set of parts with overhangs
		labels := _input.StandardLabels //[]string{"Promoters","RBSs","CDSs","Ters"}

		refactoredparts := make(map[string][]wtype.DNASequence)

		newparts := make([]wtype.DNASequence, 0)

		for _, parts := range _output.Parts {

			for j := range parts {
				newparts = refactoredparts[labels[j]]
				newparts = append(newparts, parts[j])
				refactoredparts[labels[j]] = newparts
			}
		}

		for key, value := range refactoredparts {

			duplicateremoved := search.RemoveDuplicateSequences(value)
			csv := make([]string, 0)
			revcompseqs := make([]wtype.DNASequence, 0)

			for _, seq := range duplicateremoved {
				revcomp := wtype.RevComp(seq.Seq)
				revcompseqs = append(revcompseqs, wtype.MakeSingleStrandedDNASequence(seq.Nm+"_ReverseCompliment", revcomp))
				output := fmt.Sprintln(seq.Nm, ",", seq.Seq, ",", seq.Nm+"_ReverseCompliment", ",", revcomp)
				csv = append(csv, output)
			}

			export.Makefastaserial2(export.LOCAL, filepath.Join(_input.ProjectName, key), duplicateremoved)
			export.Makefastaserial2(export.LOCAL, filepath.Join(_input.ProjectName, key+"_RevComp"), revcompseqs)
			err := export.ExporttoTextFile(filepath.Join(_input.ProjectName, key)+".csv", csv)

			if err != nil {
				panic(err.Error())
			}
		}

		// add fasta file for each set of primers
		labels = []string{"FWDPrimers", "REVPrimers"}

		refactoredparts = make(map[string][]wtype.DNASequence)
		csv := make([]string, 0)
		newparts = make([]wtype.DNASequence, 0)

		for i, parts := range _output.SequencingPrimers {

			for j := range parts {
				newparts = refactoredparts[labels[j]]
				newparts = append(newparts, parts[j])
				refactoredparts[labels[j]] = newparts
			}
			output := fmt.Sprintln(i+1, ",", "FWD: ", parts[0].Nm, ",", parts[0].Seq, ",", "REV: ", parts[1].Nm, ",", parts[1].Seq)
			csv = append(csv, output)
		}

		for key, value := range refactoredparts {

			duplicateremoved := search.RemoveDuplicateSequences(value)

			export.Makefastaserial2(export.LOCAL, filepath.Join(_input.ProjectName, key), duplicateremoved)
		}

		err := export.ExporttoTextFile(filepath.Join(_input.ProjectName, "Primers")+".csv", csv)

		if err != nil {
			panic(err.Error())
		}

	}
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _CombinatorialLibraryDesign_4part_wtype_mapAnalysis(_ctx context.Context, _input *CombinatorialLibraryDesign_4part_wtype_mapInput, _output *CombinatorialLibraryDesign_4part_wtype_mapOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _CombinatorialLibraryDesign_4part_wtype_mapValidation(_ctx context.Context, _input *CombinatorialLibraryDesign_4part_wtype_mapInput, _output *CombinatorialLibraryDesign_4part_wtype_mapOutput) {
}
func _CombinatorialLibraryDesign_4part_wtype_mapRun(_ctx context.Context, input *CombinatorialLibraryDesign_4part_wtype_mapInput) *CombinatorialLibraryDesign_4part_wtype_mapOutput {
	output := &CombinatorialLibraryDesign_4part_wtype_mapOutput{}
	_CombinatorialLibraryDesign_4part_wtype_mapSetup(_ctx, input)
	_CombinatorialLibraryDesign_4part_wtype_mapSteps(_ctx, input, output)
	_CombinatorialLibraryDesign_4part_wtype_mapAnalysis(_ctx, input, output)
	_CombinatorialLibraryDesign_4part_wtype_mapValidation(_ctx, input, output)
	return output
}

func CombinatorialLibraryDesign_4part_wtype_mapRunSteps(_ctx context.Context, input *CombinatorialLibraryDesign_4part_wtype_mapInput) *CombinatorialLibraryDesign_4part_wtype_mapSOutput {
	soutput := &CombinatorialLibraryDesign_4part_wtype_mapSOutput{}
	output := _CombinatorialLibraryDesign_4part_wtype_mapRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func CombinatorialLibraryDesign_4part_wtype_mapNew() interface{} {
	return &CombinatorialLibraryDesign_4part_wtype_mapElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &CombinatorialLibraryDesign_4part_wtype_mapInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _CombinatorialLibraryDesign_4part_wtype_mapRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &CombinatorialLibraryDesign_4part_wtype_mapInput{},
			Out: &CombinatorialLibraryDesign_4part_wtype_mapOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type CombinatorialLibraryDesign_4part_wtype_mapElement struct {
	inject.CheckedRunner
}

type CombinatorialLibraryDesign_4part_wtype_mapInput struct {
	BlastSearchSeqs    bool
	CDSs               []wtype.DNASequence
	FolderPerConstruct bool
	FolderPerProject   bool
	PROs               []wtype.DNASequence
	ProjectName        string
	RBSs               []wtype.DNASequence
	SitesToRemove      []string
	Standard           string
	StandardLabels     []string
	StandardLevel      string
	TERs               []wtype.DNASequence
	Vectors            []wtype.DNASequence
}

type CombinatorialLibraryDesign_4part_wtype_mapOutput struct {
	EndreportMap          map[string]string
	Parts                 [][]wtype.DNASequence
	PartswithOverhangsMap map[string][]wtype.DNASequence
	PassMap               map[string]bool
	PositionReportMap     map[string][]string
	PrimerMap             map[string]oligos.Primer
	SeqsMap               map[string]wtype.DNASequence
	Sequences             []wtype.DNASequence
	SequencingPrimers     [][]wtype.DNASequence
	StatusMap             map[string]string
}

type CombinatorialLibraryDesign_4part_wtype_mapSOutput struct {
	Data struct {
		EndreportMap          map[string]string
		Parts                 [][]wtype.DNASequence
		PartswithOverhangsMap map[string][]wtype.DNASequence
		PassMap               map[string]bool
		PositionReportMap     map[string][]string
		PrimerMap             map[string]oligos.Primer
		SeqsMap               map[string]wtype.DNASequence
		Sequences             []wtype.DNASequence
		SequencingPrimers     [][]wtype.DNASequence
		StatusMap             map[string]string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "CombinatorialLibraryDesign_4part_wtype_map",
		Constructor: CombinatorialLibraryDesign_4part_wtype_mapNew,
		Desc: component.ComponentDesc{
			Desc: "This protocol is intended to design assembly parts using a specified enzyme.\noverhangs are added to complement the adjacent parts and leave no scar.\nparts can be entered as genbank (.gb) files, sequences or biobrick IDs\nIf assembly simulation fails after overhangs are added. In order to help the user\ndiagnose the reason, a report of the part overhangs\nis returned to the user along with a list of cut sites in each part.\n",
			Path: "antha/component/an/Data/DNA/TypeIISAssembly_design/CombinatorialLibraryDesign4partmap.an",
			Params: []component.ParamDesc{
				{Name: "BlastSearchSeqs", Desc: "", Kind: "Parameters"},
				{Name: "CDSs", Desc: "", Kind: "Parameters"},
				{Name: "FolderPerConstruct", Desc: "", Kind: "Parameters"},
				{Name: "FolderPerProject", Desc: "", Kind: "Parameters"},
				{Name: "PROs", Desc: "", Kind: "Parameters"},
				{Name: "ProjectName", Desc: "Seqsinorder\t\t\t\t\tmap[string][]string // constructname to sequence combination\n", Kind: "Parameters"},
				{Name: "RBSs", Desc: "", Kind: "Parameters"},
				{Name: "SitesToRemove", Desc: "", Kind: "Parameters"},
				{Name: "Standard", Desc: "MoClo\n", Kind: "Parameters"},
				{Name: "StandardLabels", Desc: "", Kind: "Parameters"},
				{Name: "StandardLevel", Desc: "of assembly standard\n", Kind: "Parameters"},
				{Name: "TERs", Desc: "", Kind: "Parameters"},
				{Name: "Vectors", Desc: "", Kind: "Parameters"},
				{Name: "EndreportMap", Desc: "", Kind: "Data"},
				{Name: "Parts", Desc: "", Kind: "Data"},
				{Name: "PartswithOverhangsMap", Desc: "parts to order\n", Kind: "Data"},
				{Name: "PassMap", Desc: "", Kind: "Data"},
				{Name: "PositionReportMap", Desc: "", Kind: "Data"},
				{Name: "PrimerMap", Desc: "", Kind: "Data"},
				{Name: "SeqsMap", Desc: "desired sequence to end up with after assembly\n", Kind: "Data"},
				{Name: "Sequences", Desc: "", Kind: "Data"},
				{Name: "SequencingPrimers", Desc: "", Kind: "Data"},
				{Name: "StatusMap", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
