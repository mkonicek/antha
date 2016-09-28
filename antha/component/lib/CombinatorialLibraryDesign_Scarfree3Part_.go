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

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// parts to order

// desired sequence to end up with after assembly

// Input Requirement specification
func _CombinatorialLibraryDesign_Scarfree3PartRequirements() {
	// e.g. are MoClo types valid?
}

// Conditions to run on startup
func _CombinatorialLibraryDesign_Scarfree3PartSetup(_ctx context.Context, _input *CombinatorialLibraryDesign_Scarfree3PartInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _CombinatorialLibraryDesign_Scarfree3PartSteps(_ctx context.Context, _input *CombinatorialLibraryDesign_Scarfree3PartInput, _output *CombinatorialLibraryDesign_Scarfree3PartOutput) {
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
		for k := range _input.Part1s {
			for l := range _input.Part2s {
				for m := range _input.Part3s {

					key := _input.ProjectName + "_" + "Vector" + strconv.Itoa(j+1) + "_1." + strconv.Itoa(k+1) + "_2." + strconv.Itoa(l+1) + "_3." + strconv.Itoa(m+1)
					assembly := Scarfree_siteremove_orfcheckRunSteps(_ctx, &Scarfree_siteremove_orfcheckInput{Constructname: key,
						Seqsinorder: []string{_input.Part1s[k], _input.Part2s[l], _input.Part3s[m]},
						Enzymename:  _input.EnzymeName,
						Vector:      _input.Vectors[j],
						OtherEnzymeSitesToRemove:      _input.SitesToRemove,
						ORFstoConfirm:                 _input.ORFStoconfirm, // enter each as amino acid sequence
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

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _CombinatorialLibraryDesign_Scarfree3PartAnalysis(_ctx context.Context, _input *CombinatorialLibraryDesign_Scarfree3PartInput, _output *CombinatorialLibraryDesign_Scarfree3PartOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _CombinatorialLibraryDesign_Scarfree3PartValidation(_ctx context.Context, _input *CombinatorialLibraryDesign_Scarfree3PartInput, _output *CombinatorialLibraryDesign_Scarfree3PartOutput) {
}
func _CombinatorialLibraryDesign_Scarfree3PartRun(_ctx context.Context, input *CombinatorialLibraryDesign_Scarfree3PartInput) *CombinatorialLibraryDesign_Scarfree3PartOutput {
	output := &CombinatorialLibraryDesign_Scarfree3PartOutput{}
	_CombinatorialLibraryDesign_Scarfree3PartSetup(_ctx, input)
	_CombinatorialLibraryDesign_Scarfree3PartSteps(_ctx, input, output)
	_CombinatorialLibraryDesign_Scarfree3PartAnalysis(_ctx, input, output)
	_CombinatorialLibraryDesign_Scarfree3PartValidation(_ctx, input, output)
	return output
}

func CombinatorialLibraryDesign_Scarfree3PartRunSteps(_ctx context.Context, input *CombinatorialLibraryDesign_Scarfree3PartInput) *CombinatorialLibraryDesign_Scarfree3PartSOutput {
	soutput := &CombinatorialLibraryDesign_Scarfree3PartSOutput{}
	output := _CombinatorialLibraryDesign_Scarfree3PartRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func CombinatorialLibraryDesign_Scarfree3PartNew() interface{} {
	return &CombinatorialLibraryDesign_Scarfree3PartElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &CombinatorialLibraryDesign_Scarfree3PartInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _CombinatorialLibraryDesign_Scarfree3PartRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &CombinatorialLibraryDesign_Scarfree3PartInput{},
			Out: &CombinatorialLibraryDesign_Scarfree3PartOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type CombinatorialLibraryDesign_Scarfree3PartElement struct {
	inject.CheckedRunner
}

type CombinatorialLibraryDesign_Scarfree3PartInput struct {
	BlastSearchSeqs bool
	EnzymeName      string
	ORFStoconfirm   []string
	Part1s          []string
	Part2s          []string
	Part3s          []string
	ProjectName     string
	SitesToRemove   []string
	Vectors         []string
}

type CombinatorialLibraryDesign_Scarfree3PartOutput struct {
	EndreportMap          map[string]string
	PartswithOverhangsMap map[string][]wtype.DNASequence
	PassMap               map[string]bool
	PositionReportMap     map[string][]string
	PrimerMap             map[string]oligos.Primer
	SeqsMap               map[string]wtype.DNASequence
	StatusMap             map[string]string
}

type CombinatorialLibraryDesign_Scarfree3PartSOutput struct {
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
	if err := addComponent(component.Component{Name: "CombinatorialLibraryDesign_Scarfree3Part",
		Constructor: CombinatorialLibraryDesign_Scarfree3PartNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Data/DNA/TypeIISAssembly_design/CombinatorialLibraryDesign_Scarfree.an",
			Params: []component.ParamDesc{
				{Name: "BlastSearchSeqs", Desc: "", Kind: "Parameters"},
				{Name: "EnzymeName", Desc: "", Kind: "Parameters"},
				{Name: "ORFStoconfirm", Desc: "", Kind: "Parameters"},
				{Name: "Part1s", Desc: "", Kind: "Parameters"},
				{Name: "Part2s", Desc: "", Kind: "Parameters"},
				{Name: "Part3s", Desc: "", Kind: "Parameters"},
				{Name: "ProjectName", Desc: "", Kind: "Parameters"},
				{Name: "SitesToRemove", Desc: "", Kind: "Parameters"},
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
