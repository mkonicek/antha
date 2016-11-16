// This protocol is intended to design assembly parts using a specified enzyme.
// overhangs are added to complement the adjacent parts and leave no scar.
// parts can be entered as genbank (.gb) files, sequences or biobrick IDs
// If assembly simulation fails after overhangs are added. In order to help the user
// diagnose the reason, a report of the part overhangs
// is returned to the user along with a list of cut sites in each part.

package lib

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

// constructname to sequence combination
//MoClo
// of assembly standard

// labels e.g. pro = promoter

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// parts to order

// desired sequence to end up with after assembly

// Input Requirement specification
func _CombinatorialLibraryDesign_PRO_RBS_CDSRequirements() {
	// e.g. are MoClo types valid?
}

// Conditions to run on startup
func _CombinatorialLibraryDesign_PRO_RBS_CDSSetup(_ctx context.Context, _input *CombinatorialLibraryDesign_PRO_RBS_CDSInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _CombinatorialLibraryDesign_PRO_RBS_CDSSteps(_ctx context.Context, _input *CombinatorialLibraryDesign_PRO_RBS_CDSInput, _output *CombinatorialLibraryDesign_PRO_RBS_CDSOutput) {
	_output.StatusMap = make(map[string]string)
	_output.PartswithOverhangsMap = make(map[string][]wtype.DNASequence) // parts to order
	_output.PassMap = make(map[string]bool)
	_output.SeqsMap = make(map[string]wtype.DNASequence) // desired sequence to end up with after assembly
	_output.EndreportMap = make(map[string]string)
	_output.PositionReportMap = make(map[string][]string)
	_output.StatusMap = make(map[string]string)

	for j := range _input.Vectors {
		for key, seqsinorder := range _input.Seqsinorder {
			assembly := AssemblyStandard_siteremove_orfcheckRunSteps(_ctx, &AssemblyStandard_siteremove_orfcheckInput{Constructname: key,
				Seqsinorder:                   seqsinorder,
				AssemblyStandard:              _input.Standard,
				Level:                         _input.StandardLevel, // of assembly standard
				Vector:                        _input.Vectors[j],
				PartMoClotypesinorder:         _input.PartTypesinOrder, // labels e.g. pro = promoter
				OtherEnzymeSitesToRemove:      _input.SitesToRemove,
				ORFstoConfirm:                 []string{}, // enter each as amino acid sequence
				RemoveproblemRestrictionSites: true,
				EndsAlreadyadded:              false,
				ExporttoFastaFile:             true},
			)
			key = key + _input.Vectors[j]
			_output.PartswithOverhangsMap[key] = assembly.Data.PartswithOverhangs // parts to order
			_output.PassMap[key] = assembly.Data.Simulationpass
			_output.EndreportMap[key] = assembly.Data.Endreport
			_output.PositionReportMap[key] = assembly.Data.PositionReport
			_output.SeqsMap[key] = assembly.Data.NewDNASequence
			_output.StatusMap[key] = assembly.Data.Status
		}
	}
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _CombinatorialLibraryDesign_PRO_RBS_CDSAnalysis(_ctx context.Context, _input *CombinatorialLibraryDesign_PRO_RBS_CDSInput, _output *CombinatorialLibraryDesign_PRO_RBS_CDSOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _CombinatorialLibraryDesign_PRO_RBS_CDSValidation(_ctx context.Context, _input *CombinatorialLibraryDesign_PRO_RBS_CDSInput, _output *CombinatorialLibraryDesign_PRO_RBS_CDSOutput) {
}
func _CombinatorialLibraryDesign_PRO_RBS_CDSRun(_ctx context.Context, input *CombinatorialLibraryDesign_PRO_RBS_CDSInput) *CombinatorialLibraryDesign_PRO_RBS_CDSOutput {
	output := &CombinatorialLibraryDesign_PRO_RBS_CDSOutput{}
	_CombinatorialLibraryDesign_PRO_RBS_CDSSetup(_ctx, input)
	_CombinatorialLibraryDesign_PRO_RBS_CDSSteps(_ctx, input, output)
	_CombinatorialLibraryDesign_PRO_RBS_CDSAnalysis(_ctx, input, output)
	_CombinatorialLibraryDesign_PRO_RBS_CDSValidation(_ctx, input, output)
	return output
}

func CombinatorialLibraryDesign_PRO_RBS_CDSRunSteps(_ctx context.Context, input *CombinatorialLibraryDesign_PRO_RBS_CDSInput) *CombinatorialLibraryDesign_PRO_RBS_CDSSOutput {
	soutput := &CombinatorialLibraryDesign_PRO_RBS_CDSSOutput{}
	output := _CombinatorialLibraryDesign_PRO_RBS_CDSRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func CombinatorialLibraryDesign_PRO_RBS_CDSNew() interface{} {
	return &CombinatorialLibraryDesign_PRO_RBS_CDSElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &CombinatorialLibraryDesign_PRO_RBS_CDSInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _CombinatorialLibraryDesign_PRO_RBS_CDSRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &CombinatorialLibraryDesign_PRO_RBS_CDSInput{},
			Out: &CombinatorialLibraryDesign_PRO_RBS_CDSOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type CombinatorialLibraryDesign_PRO_RBS_CDSElement struct {
	inject.CheckedRunner
}

type CombinatorialLibraryDesign_PRO_RBS_CDSInput struct {
	PartTypesinOrder []string
	Seqsinorder      map[string][]string
	SitesToRemove    []string
	Standard         string
	StandardLevel    string
	Vectors          []string
}

type CombinatorialLibraryDesign_PRO_RBS_CDSOutput struct {
	EndreportMap          map[string]string
	PartswithOverhangsMap map[string][]wtype.DNASequence
	PassMap               map[string]bool
	PositionReportMap     map[string][]string
	SeqsMap               map[string]wtype.DNASequence
	StatusMap             map[string]string
}

type CombinatorialLibraryDesign_PRO_RBS_CDSSOutput struct {
	Data struct {
		EndreportMap          map[string]string
		PartswithOverhangsMap map[string][]wtype.DNASequence
		PassMap               map[string]bool
		PositionReportMap     map[string][]string
		SeqsMap               map[string]wtype.DNASequence
		StatusMap             map[string]string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "CombinatorialLibraryDesign_PRO_RBS_CDS",
		Constructor: CombinatorialLibraryDesign_PRO_RBS_CDSNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/DNA/TypeIISAssembly_design/CombinatorialLibraryDesign.an",
			Params: []component.ParamDesc{
				{Name: "PartTypesinOrder", Desc: "labels e.g. pro = promoter\n", Kind: "Parameters"},
				{Name: "Seqsinorder", Desc: "constructname to sequence combination\n", Kind: "Parameters"},
				{Name: "SitesToRemove", Desc: "", Kind: "Parameters"},
				{Name: "Standard", Desc: "MoClo\n", Kind: "Parameters"},
				{Name: "StandardLevel", Desc: "of assembly standard\n", Kind: "Parameters"},
				{Name: "Vectors", Desc: "", Kind: "Parameters"},
				{Name: "EndreportMap", Desc: "", Kind: "Data"},
				{Name: "PartswithOverhangsMap", Desc: "parts to order\n", Kind: "Data"},
				{Name: "PassMap", Desc: "", Kind: "Data"},
				{Name: "PositionReportMap", Desc: "", Kind: "Data"},
				{Name: "SeqsMap", Desc: "desired sequence to end up with after assembly\n", Kind: "Data"},
				{Name: "StatusMap", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
