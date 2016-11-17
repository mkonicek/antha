// This protocol is intended to design assembly parts using a specified enzyme.
// overhangs are added to complement the adjacent parts and leave no scar.
// parts can be entered as genbank (.gb) files, sequences or biobrick IDs
// If assembly simulation fails after overhangs are added. In order to help the user
// diagnose the reason, a report of the part overhangs
// is returned to the user along with a list of cut sites in each part.

package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/oligos"
	"github.com/antha-lang/antha/antha/anthalib/wtype"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
	"strconv"
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
func _CombinatorialLibraryDesign_PRO_RBS_CDS_mapRequirements() {
	// e.g. are MoClo types valid?
}

// Conditions to run on startup
func _CombinatorialLibraryDesign_PRO_RBS_CDS_mapSetup(_ctx context.Context, _input *CombinatorialLibraryDesign_PRO_RBS_CDS_mapInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _CombinatorialLibraryDesign_PRO_RBS_CDS_mapSteps(_ctx context.Context, _input *CombinatorialLibraryDesign_PRO_RBS_CDS_mapInput, _output *CombinatorialLibraryDesign_PRO_RBS_CDS_mapOutput) {
	_output.StatusMap = make(map[string]string)
	_output.PartswithOverhangsMap = make(map[string][]wtype.DNASequence) // parts to order
	_output.PassMap = make(map[string]bool)
	_output.SeqsMap = make(map[string]wtype.DNASequence) // desired sequence to end up with after assembly
	_output.EndreportMap = make(map[string]string)
	_output.PositionReportMap = make(map[string][]string)
	_output.StatusMap = make(map[string]string)
	_output.PrimerMap = make(map[string]oligos.Primer)

	var counter int = 1

	for j := range _input.Vectors {
		for k := range _input.PROs {
			for l := range _input.RBSs {
				for m := range _input.CDSs {
					for n := range _input.TERs {
						key := _input.ProjectName + strconv.Itoa(counter)
						assembly := AssemblyStandard_siteremove_orfcheckRunSteps(_ctx, &AssemblyStandard_siteremove_orfcheckInput{Constructname: key,
							Seqsinorder:                   []string{_input.PROs[k], _input.RBSs[l], _input.CDSs[m], _input.TERs[n]},
							AssemblyStandard:              _input.Standard,
							Level:                         _input.StandardLevel, // of assembly standard
							Vector:                        _input.Vectors[j],
							PartMoClotypesinorder:         []string{"Pro", "5U + NT1", "CDS1", "3U + Ter"},
							OtherEnzymeSitesToRemove:      _input.SitesToRemove,
							ORFstoConfirm:                 []string{}, // enter each as amino acid sequence
							RemoveproblemRestrictionSites: true,
							EndsAlreadyadded:              false,
							ExporttoFastaFile:             true,
							BlastSeqswithNoName:           _input.BlastSearchSeqs},
						)
						key = key                                                             //+ Vectors[j]
						_output.PartswithOverhangsMap[key] = assembly.Data.PartswithOverhangs // parts to order
						_output.PassMap[key] = assembly.Data.Simulationpass
						_output.EndreportMap[key] = assembly.Data.Endreport
						_output.PositionReportMap[key] = assembly.Data.PositionReport
						_output.SeqsMap[key] = assembly.Data.NewDNASequence
						_output.StatusMap[key] = assembly.Data.Status

						// for each vector we'll also design sequencing primers

						primer := PrimerDesign_FWD_wtypeRunSteps(_ctx, &PrimerDesign_FWD_wtypeInput{FullDNASeq: assembly.Data.NewDNASequence,
							Maxtemp:                                  wunit.NewTemperature(60, "C"),
							Mintemp:                                  wunit.NewTemperature(55, "C"),
							Maxgc:                                    0.6,
							Minlength:                                20,
							Maxlength:                                25,
							Seqstoavoid:                              []string{},
							PermittednucleotideOverlapBetweenPrimers: 10,                                     // number of nucleotides which primers can overlap by
							RegionSequence:                           assembly.Data.PartsWithSitesRemoved[0], // first part
							FlankTargetSequence:                      true},
						)
						_output.PrimerMap[key] = primer.Data.FWDPrimer
						counter++
					}
				}
			}
		}

	}
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _CombinatorialLibraryDesign_PRO_RBS_CDS_mapAnalysis(_ctx context.Context, _input *CombinatorialLibraryDesign_PRO_RBS_CDS_mapInput, _output *CombinatorialLibraryDesign_PRO_RBS_CDS_mapOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _CombinatorialLibraryDesign_PRO_RBS_CDS_mapValidation(_ctx context.Context, _input *CombinatorialLibraryDesign_PRO_RBS_CDS_mapInput, _output *CombinatorialLibraryDesign_PRO_RBS_CDS_mapOutput) {
}
func _CombinatorialLibraryDesign_PRO_RBS_CDS_mapRun(_ctx context.Context, input *CombinatorialLibraryDesign_PRO_RBS_CDS_mapInput) *CombinatorialLibraryDesign_PRO_RBS_CDS_mapOutput {
	output := &CombinatorialLibraryDesign_PRO_RBS_CDS_mapOutput{}
	_CombinatorialLibraryDesign_PRO_RBS_CDS_mapSetup(_ctx, input)
	_CombinatorialLibraryDesign_PRO_RBS_CDS_mapSteps(_ctx, input, output)
	_CombinatorialLibraryDesign_PRO_RBS_CDS_mapAnalysis(_ctx, input, output)
	_CombinatorialLibraryDesign_PRO_RBS_CDS_mapValidation(_ctx, input, output)
	return output
}

func CombinatorialLibraryDesign_PRO_RBS_CDS_mapRunSteps(_ctx context.Context, input *CombinatorialLibraryDesign_PRO_RBS_CDS_mapInput) *CombinatorialLibraryDesign_PRO_RBS_CDS_mapSOutput {
	soutput := &CombinatorialLibraryDesign_PRO_RBS_CDS_mapSOutput{}
	output := _CombinatorialLibraryDesign_PRO_RBS_CDS_mapRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func CombinatorialLibraryDesign_PRO_RBS_CDS_mapNew() interface{} {
	return &CombinatorialLibraryDesign_PRO_RBS_CDS_mapElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &CombinatorialLibraryDesign_PRO_RBS_CDS_mapInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _CombinatorialLibraryDesign_PRO_RBS_CDS_mapRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &CombinatorialLibraryDesign_PRO_RBS_CDS_mapInput{},
			Out: &CombinatorialLibraryDesign_PRO_RBS_CDS_mapOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type CombinatorialLibraryDesign_PRO_RBS_CDS_mapElement struct {
	inject.CheckedRunner
}

type CombinatorialLibraryDesign_PRO_RBS_CDS_mapInput struct {
	BlastSearchSeqs bool
	CDSs            []string
	PROs            []string
	ProjectName     string
	RBSs            []string
	SitesToRemove   []string
	Standard        string
	StandardLevel   string
	TERs            []string
	Vectors         []string
}

type CombinatorialLibraryDesign_PRO_RBS_CDS_mapOutput struct {
	EndreportMap          map[string]string
	PartswithOverhangsMap map[string][]wtype.DNASequence
	PassMap               map[string]bool
	PositionReportMap     map[string][]string
	PrimerMap             map[string]oligos.Primer
	SeqsMap               map[string]wtype.DNASequence
	StatusMap             map[string]string
}

type CombinatorialLibraryDesign_PRO_RBS_CDS_mapSOutput struct {
	Data struct {
		EndreportMap          map[string]string
		PartswithOverhangsMap map[string][]wtype.DNASequence
		PassMap               map[string]bool
		PositionReportMap     map[string][]string
		PrimerMap             map[string]oligos.Primer
		SeqsMap               map[string]wtype.DNASequence
		StatusMap             map[string]string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "CombinatorialLibraryDesign_PRO_RBS_CDS_map",
		Constructor: CombinatorialLibraryDesign_PRO_RBS_CDS_mapNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/DNA/TypeIISAssembly_design/CombinatorialLibraryDesign4part.an",
			Params: []component.ParamDesc{
				{Name: "BlastSearchSeqs", Desc: "", Kind: "Parameters"},
				{Name: "CDSs", Desc: "", Kind: "Parameters"},
				{Name: "PROs", Desc: "", Kind: "Parameters"},
				{Name: "ProjectName", Desc: "Seqsinorder\t\t\t\t\tmap[string][]string // constructname to sequence combination\n", Kind: "Parameters"},
				{Name: "RBSs", Desc: "", Kind: "Parameters"},
				{Name: "SitesToRemove", Desc: "", Kind: "Parameters"},
				{Name: "Standard", Desc: "MoClo\n", Kind: "Parameters"},
				{Name: "StandardLevel", Desc: "of assembly standard\n", Kind: "Parameters"},
				{Name: "TERs", Desc: "", Kind: "Parameters"},
				{Name: "Vectors", Desc: "", Kind: "Parameters"},
				{Name: "EndreportMap", Desc: "", Kind: "Data"},
				{Name: "PartswithOverhangsMap", Desc: "parts to order\n", Kind: "Data"},
				{Name: "PassMap", Desc: "", Kind: "Data"},
				{Name: "PositionReportMap", Desc: "", Kind: "Data"},
				{Name: "PrimerMap", Desc: "", Kind: "Data"},
				{Name: "SeqsMap", Desc: "desired sequence to end up with after assembly\n", Kind: "Data"},
				{Name: "StatusMap", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
