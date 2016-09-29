package lib

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/factory"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

// PCRprep parameters

// e.g. ["left homology arm"]:"templatename"
// e.g. ["left homology arm"]:"fwdprimer","revprimer"

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _AutoPCR_primerbindRequirements() {
}

// Conditions to run on startup
func _AutoPCR_primerbindSetup(_ctx context.Context, _input *AutoPCR_primerbindInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _AutoPCR_primerbindSteps(_ctx context.Context, _input *AutoPCR_primerbindInput, _output *AutoPCR_primerbindOutput) {

	var counter int

	_output.Reactions = make([]*wtype.LHComponent, 0)
	volumes := make([]wunit.Volume, 0)
	welllocations := make([]string, 0)

	for reactionname, templatename := range _input.Reactiontotemplate {

		wellposition := _input.Plate.AllWellPositions(wtype.BYCOLUMN)[counter]

		result := PCR_vol_mmx_primerbindRunSteps(_ctx, &PCR_vol_mmx_primerbindInput{Targetsequence: "TAATGACCCCAAGGGCGACACCCCCTAATTAGCCCGGGCGAAAGGCCCAGTCTTTCGACTGAGCCTTTCGTTTTATTTGATGCCTGGCAGTTCCCTACTCTCGCATGGGGAGTCCCCACACTACCATCGGCGCTACGGCGTTTCACTTCTGAGTTCGGCATGGGGTCAGGTGGGACCACCGCGCTACTGCCGCCAGGCAAACAAGGGGTGTTATGAGCCATATTCAGGTATAAATGGGCTCGCGATAATGTTCAGAATTGGTTAATTGGTTGTAACACTGACCCCTATTTGTTTATTTTTCTAAATACATTCAAATATGTATCCGCTCATGAGACAATAACCCTGATAAATGCTTCAATAATATTGAAAAAGGAAGAATATGAGCCATATTCAACGGGAAACGTCGAGGCCGCGATTAAATTCCAACATGGATGCTGATTTATATGGGTATAAATGGGCTCGCGATAATGTCGGGCAATCAGGTGCGACAATCTATCGCTTGTATGGGAAGCCCGATGCGCCAGAGTTGTTTCTGAAACATGGCAAAGGTAGCGTTGCCAATGATGTTACAGATGAGATGGTCAGACTAAACTGGCTGACGGAATTTATGCCACTTCCGACCATCAAGCATTTTATCCGTACTCCTGATGATGCATGGTTACTCACCACTGCGATCCCCGGAAAAACAGCGTTCCAGGTATTAGAAGAATATCCTGATTCAGGTGAAAATATTGTTGATGCGCTGGCAGTGTTCCTGCGCCGGTTGCACTCGATTCCTGTTTGTAATTGTCCTTTTAACAGCGATCGCGTATTTCGCCTCGCTCAGGCGCAATCACGAATGAATAACGGTTTGGTTGATGCGAGTGATTTTGATGACGAGCGTAATGGCTGGCCTGTTGAACAAGTCTGGAAAGAAATGCATAAACTTTTGCCATTCTCACCGGATTCAGTCGTCACTCATGGTGATTTCTCACTTGATAACCTTATTTTTGACGAGGGGAAATTAATAGGTTGTATTGATGTTGGACGAGTCGGAATCGCAGACCGATACCAGGATCTTGCCATCCTATGGAACTGCCTCGGTGAGTTTTCTCCTTCATTACAGAAACGGCTTTTTCAAAAATATGGTATTGATAATCCTGATATGAATAAATTGCAGTTTCATTTGATGCTCGATGAGTTTTTCTAAGCGGCGCGCCATCGAATGGCGCAAAACCTTTCGCGGTATGGCATGATAGCGCCCGGAAGAGAGTCAATTCAGGGTGGTGAATATGAAACCAGTAACGTTATACGATGTCGCAGAGTATGCCGGTGTCTCTTATCAGACCGTTTCCCGCGTGGTGAACCAGGCCAGCCACGTTTCTGCGAAAACGCGGGAAAAAGTGGAAGCGGCGATGGCGGAGCTGAATTACATTCCCAACCGCGTGGCACAACAACTGGCGGGCAAACAGTCGTTGCTGATTGGCGTTGCCACCTCCAGTCTGGCCCTGCACGCGCCGTCGCAAATTGTCGCGGCGATTAAATCTCGCGCCGATCAACTGGGTGCCAGCGTGGTGGTGTCGATGGTAGAACGAAGCGGCGTCGAAGCCTGTAAAGCGGCGGTGCACAATCTTCTCGCGCAACGCGTCAGTGGGCTGATCATTAACTATCCGCTGGATGACCAGGATGCCATTGCTGTGGAAGCTGCCTGCACTAATGTTCCGGCGTTATTTCTTGATGTCTCTGACCAGACACCCATCAACAGTATTATTTTCTCCCATGAGGACGGTACGCGACTGGGCGTGGAGCATCTGGTCGCATTGGGTCACCAGCAAATCGCGCTGTTAGCGGGCCCATTAAGTTCTGTCTCGGCGCGTCTGCGTCTGGCTGGCTGGCATAAATATCTCACTCGCAATCAAATTCAGCCGATAGCGGAACGGGAAGGCGACTGGAGTGCCATGTCCGGTTTTCAACAAACCATGCAAATGCTGAATGAGGGCATCGTTCCCACTGCGATGCTGGTTGCCAACGATCAGATGGCGCTGGGCGCAATGCGCGCCATTACCGAGTCCGGGCTGCGCGTTGGTGCGGATATCTCGGTAGTGGGATACGACGATACCGAAGATAGCTCATGTTATATCCCGCCGTTAACCACCATCAAACAGGATTTTCGCCTGCTGGGGCAAACCAGCGTGGACCGCTTGCTGCAACTCTCTCAGGGCCAGGCGGTGAAGGGCAATCAGCTGTTGCCAGTCTCACTGGTGAAAAGAAAAACCACCCTGGCGCCCAATACGCAAACCGCCTCTCCCCGCGCGTTGGCCGATTCATTAATGCAGCTGGCACGACAGGTTTCCCGACTGGAAAGCGGGCAGTGACTCATGACCAAAATCCCTTAACGTGAGTTACGCGCGCGTCGTTCCACTGAGCGTCAGACCCCGTAGAAAAGATCAAAGGATCTTCTTGAGATCCTTTTTTTCTGCGCGTAATCTGCTGCTTGCAAACAAAAAAACCACCGCTACCAGCGGTGGTTTGTTTGCCGGATCAAGAGCTACCAACTCTTTTTCCGAAGGTAACTGGCTTCAGCAGAGCGCAGATACCAAATACTGTTCTTCTAGTGTAGCCGTAGTTAGCCCACCACTTCAAGAACTCTGTAGCACCGCCTACATACCTCGCTCTGCTAATCCTGTTACCAGTGGCTGCTGCCAGTGGCGATAAGTCGTGTCTTACCGGGTTGGACTCAAGACGATAGTTACCGGATAAGGCGCAGCGGTCGGGCTGAACGGGGGGTTCGTGCACACAGCCCAGCTTGGAGCGAACGACCTACACCGAACTGAGATACCTACAGCGTGAGCTATGAGAAAGCGCCACGCTTCCCGAAGGGAGAAAGGCGGACAGGTATCCGGTAAGCGGCAGGGTCGGAACAGGAGAGCGCACGAGGGAGCTTCCAGGGGGAAACGCCTGGTATCTTTATAGTCCTGTCGGGTTTCGCCACCTCTGACTTGAGCGTCGATTTTTGTGATGCTCGTCAGGGGGGCGGAGCCTATGGAAAAACGCCAGCAACGCGGCCTTTTTACGGTTCCTGGCCTTTTGCTGGCCTTTTGCTCACATGTTCTTTCCTGCGTTATCCCCTGATTCTGTGGATAACCGTATTACCGCCTTTGAGTGAGCTGATACCGCTCGCCGCAGCCGAACGACCGAGCGCAGCGAGTCAGTGAGCGAGGAAGCGGAAGGCGAGAGTAGGGAACTGCCAGGCATCAAACTAAGCAGAAGGCCCCTGACGGATGGCCTTTTTGCGTTTCTACAAACTCTTTCTGTGTTGTAAAACGACGGCCAGTCTTAAGCTCGGGCCCCCTGGGCGGTTCTGATAACGAGTAATCGTTAATCCGCAAATAACGTAAAAACCCGCTTCGGCGGGTTTTTTTATGGGGGGAGTTTAGGGAAAGAGCATTTGTCAGAATATTTAAGGGCGCCTGTCACTTTGCTTGATATATGAGAATTATTTAACCTTATAAATGAGAAAAAAGCAACGCACTTTAAATAAGATACGTTGCTTTTTCGATTGATGAACACCTATAATTAAACTATTCATCTATTATTTATGATTTTTTGTATATACAATATTTCTAGTTTGTTAAAGAGAATTAAGAAAATAAATCTCGAAAATAATAAAGGGAAAATCAGTTTTTGATATCAAAATTATACATGTCAACGATAATACAAAATATAATACAAACTATAAGATGTTATCAGTATTTATTATGCATTTAGAATAAATTTTGTGTCGGCTCTTCATGTCCACAATTCAGCAAATTGTGAACATCATCACGTTCATCTTTCCCTGGTTGCCAATGGCCCATTTTCCTGTCAGTAACGAGAAGGTCGCGTATTCAGGCGCTTTTTAGACTGGTCGTAATGAATACTGAATTCATTAAAGAGGAGAAAGGTCATATGAATAGCCTGATTAAAGAGAATATGCACATGAAGCTGTACATGGAAGGCACGGTGAATAACCACCACTTCAAATGCACCAGCGAGGGTGAGGGTAAACCGTATGAAGGCACCCAAACGATGCGTATCAAAGTTGTTGAGGGTGGCCCGTTGCCGTTTGCGTTCGACATTTTAGCGACGAGCTTTATGTATGGCTCTCGTACGTTTATCAAGTACCCGAAGGGTATTCCGGACTTTTTCAAACAATCTTTTCCAGAGGGTTTCACCTGGGAGCGCGTGACTCGCTACGAAGATGGCGGCGTCGTGACCGCAACGCAGGATACCTCCCTGGAAGATGGCTGCCTGGTCTACCACGTTCAGGTCCGTGGTGTCAATTTCCCGAGCAATGGTCCGGTTATGCAGAAGAAAACCCTGGGTTGGGAACCGAACACCGAGATGTTGTATCCTGCAGATGGTGGCCTGGAAGGTCGCAGCGACATGGCATTGAAACTGGTCGGTGGCGGCCATCTGAGCTGTAGCTTCGTGACCACGTATCGTTCGAAGAAAACGGTCGGTAACATCAAAATGCCGGGTATTCACGCGGTTGACCACCGTCTGGTGCGCATTAAAGAAGCCGACAAAGAGACTTACGTGGAGCAACATGAAGTAGCCGTTGCGAAATTTGCTGGTTTGGGCGGTGGTATGGACGAACTGTACAAAGCTTCCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGTTTTATCTGTTGTTTGTCGGTGAACGCTCTCTACTAGAGTCACACTGGCTCACCTTCGGGTGGGCCTTTCTGCGTTTATAGGTAGAAGAGC",
			FwdPrimerSeq:         "TCAGTAACGAGAAGGTCGCGTATTCAGGCGCTTT",
			RevPrimerSeq:         "ATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTT",
			FwdPrimerName:        _input.Reactiontoprimerpair[reactionname][0],
			RevPrimerName:        _input.Reactiontoprimerpair[reactionname][1],
			TemplateName:         templatename,
			ReactionName:         reactionname,
			FwdPrimerVol:         wunit.NewVolume(1, "ul"),
			RevPrimerVol:         wunit.NewVolume(1, "ul"),
			Templatevolume:       wunit.NewVolume(1, "ul"),
			PolymeraseVolume:     wunit.NewVolume(1, "ul"),
			Numberofcycles:       30,
			InitDenaturationtime: wunit.NewTime(30, "s"),
			Denaturationtime:     wunit.NewTime(5, "s"),
			Annealingtime:        wunit.NewTime(10, "s"),
			AnnealingTemp:        wunit.NewTemperature(72, "C"), // Should be calculated from primer and template binding
			Extensiontime:        wunit.NewTime(60, "s"),        // should be calculated from template length and polymerase rate
			Finalextensiontime:   wunit.NewTime(180, "s"),
			WellPosition:         wellposition,

			FwdPrimer: _input.FwdPrimertype,
			RevPrimer: _input.RevPrimertype,

			PCRPolymerase: factory.GetComponentByType("Q5Polymerase"),

			Template: _input.Templatetype,

			OutPlate: _input.Plate},
		)

		_output.Reactions = append(_output.Reactions, result.Outputs.Reaction)
		volumes = append(volumes, result.Outputs.Reaction.Volume())
		welllocations = append(welllocations, wellposition)
		counter++

	}

	_output.Error = wtype.ExportPlateCSV(_input.Projectname+".csv", _input.Plate, _input.Projectname+"outputPlate", welllocations, _output.Reactions, volumes)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AutoPCR_primerbindAnalysis(_ctx context.Context, _input *AutoPCR_primerbindInput, _output *AutoPCR_primerbindOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AutoPCR_primerbindValidation(_ctx context.Context, _input *AutoPCR_primerbindInput, _output *AutoPCR_primerbindOutput) {
}
func _AutoPCR_primerbindRun(_ctx context.Context, input *AutoPCR_primerbindInput) *AutoPCR_primerbindOutput {
	output := &AutoPCR_primerbindOutput{}
	_AutoPCR_primerbindSetup(_ctx, input)
	_AutoPCR_primerbindSteps(_ctx, input, output)
	_AutoPCR_primerbindAnalysis(_ctx, input, output)
	_AutoPCR_primerbindValidation(_ctx, input, output)
	return output
}

func AutoPCR_primerbindRunSteps(_ctx context.Context, input *AutoPCR_primerbindInput) *AutoPCR_primerbindSOutput {
	soutput := &AutoPCR_primerbindSOutput{}
	output := _AutoPCR_primerbindRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AutoPCR_primerbindNew() interface{} {
	return &AutoPCR_primerbindElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AutoPCR_primerbindInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AutoPCR_primerbindRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AutoPCR_primerbindInput{},
			Out: &AutoPCR_primerbindOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type AutoPCR_primerbindElement struct {
	inject.CheckedRunner
}

type AutoPCR_primerbindInput struct {
	FwdPrimertype        *wtype.LHComponent
	Plate                *wtype.LHPlate
	Projectname          string
	Reactiontoprimerpair map[string][]string
	Reactiontotemplate   map[string]string
	RevPrimertype        *wtype.LHComponent
	Templatetype         *wtype.LHComponent
}

type AutoPCR_primerbindOutput struct {
	Error     error
	Reactions []*wtype.LHComponent
}

type AutoPCR_primerbindSOutput struct {
	Data struct {
		Error error
	}
	Outputs struct {
		Reactions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AutoPCR_primerbind",
		Constructor: AutoPCR_primerbindNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/PCR/AutoPCR_primerbind.an",
			Params: []component.ParamDesc{
				{Name: "FwdPrimertype", Desc: "", Kind: "Inputs"},
				{Name: "Plate", Desc: "", Kind: "Inputs"},
				{Name: "Projectname", Desc: "PCRprep parameters\n", Kind: "Parameters"},
				{Name: "Reactiontoprimerpair", Desc: "e.g. [\"left homology arm\"]:\"fwdprimer\",\"revprimer\"\n", Kind: "Parameters"},
				{Name: "Reactiontotemplate", Desc: "e.g. [\"left homology arm\"]:\"templatename\"\n", Kind: "Parameters"},
				{Name: "RevPrimertype", Desc: "", Kind: "Inputs"},
				{Name: "Templatetype", Desc: "", Kind: "Inputs"},
				{Name: "Error", Desc: "", Kind: "Data"},
				{Name: "Reactions", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
