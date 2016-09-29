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
// if true, 12 replicates of each reaction will be set up, one set of reactions per row, else 8 reactions set up 1 set per column

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _AutoPCR_GradientRequirements() {
}

// Conditions to run on startup
func _AutoPCR_GradientSetup(_ctx context.Context, _input *AutoPCR_GradientInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _AutoPCR_GradientSteps(_ctx context.Context, _input *AutoPCR_GradientInput, _output *AutoPCR_GradientOutput) {
	var Samplenumber int

	// if RowGradientRatherthanColumn == true,
	//12 replicates of each reaction will be set up,
	//one set of reactions per row,
	//else 8 reactions set up 1 set per column
	if _input.RowGradientRatherthanColumn {
		Samplenumber = 12
	} else {
		Samplenumber = 8
	}

	var counter int

	_output.Reactions = make([]*wtype.LHComponent, 0)
	volumes := make([]wunit.Volume, 0)
	welllocations := make([]string, 0)

	for reactionname, templatename := range _input.Reactiontotemplate {

		//wellposition := Plate.AllWellPositions(wtype.BYCOLUMN)[counter]

		for j := 0; j < Samplenumber; j++ {
			//	for i:= 0;i < len(Reactions);i++{

			var wellcoords = wtype.WellCoords{X: counter, Y: j}

			if _input.RowGradientRatherthanColumn {
				wellcoords = wtype.WellCoords{X: j, Y: counter}
			} else {
				wellcoords = wtype.WellCoords{X: counter, Y: j}
			}

			wellposition := wellcoords.FormatA1()
			result := PCR_volRunSteps(_ctx, &PCR_volInput{WaterVolume: wunit.NewVolume(10, "ul"),
				ReactionVolume:        wunit.NewVolume(25, "ul"),
				BufferConcinX:         5,
				FwdPrimerName:         _input.Reactiontoprimerpair[reactionname][0],
				RevPrimerName:         _input.Reactiontoprimerpair[reactionname][1],
				TemplateName:          templatename,
				ReactionName:          reactionname,
				FwdPrimerVol:          wunit.NewVolume(1, "ul"),
				RevPrimerVol:          wunit.NewVolume(1, "ul"),
				AdditiveVols:          []wunit.Volume{wunit.NewVolume(5, "ul")},
				Templatevolume:        wunit.NewVolume(1, "ul"),
				PolymeraseVolume:      wunit.NewVolume(1, "ul"),
				DNTPVol:               wunit.NewVolume(1, "ul"),
				Numberofcycles:        30,
				InitDenaturationtime:  wunit.NewTime(30, "s"),
				Denaturationtime:      wunit.NewTime(5, "s"),
				Annealingtime:         wunit.NewTime(10, "s"),
				AnnealingTemp:         wunit.NewTemperature(72, "C"), // Should be calculated from primer and template binding
				Extensiontime:         wunit.NewTime(60, "s"),        // should be calculated from template length and polymerase rate
				Finalextensiontime:    wunit.NewTime(180, "s"),
				Hotstart:              false,
				AddPrimerstoMasterMix: false,
				WellPosition:          wellposition,

				FwdPrimer:     _input.FwdPrimertype,
				RevPrimer:     _input.RevPrimertype,
				DNTPS:         factory.GetComponentByType("DNTPs"),
				PCRPolymerase: factory.GetComponentByType("Q5Polymerase"),
				Buffer:        factory.GetComponentByType("Q5buffer"),
				Water:         factory.GetComponentByType("water"),
				Template:      _input.Templatetype,
				Additives:     []*wtype.LHComponent{factory.GetComponentByType("GCenhancer")},
				OutPlate:      _input.Plate},
			)

			_output.Reactions = append(_output.Reactions, result.Outputs.Reaction)
			volumes = append(volumes, result.Outputs.Reaction.Volume())
			welllocations = append(welllocations, wellposition)

		}
		counter++
	}
	_output.NumberOfReactions = len(_output.Reactions)
	_output.Error = wtype.ExportPlateCSV(_input.Projectname+".csv", _input.Plate, _input.Projectname+"outputPlate", welllocations, _output.Reactions, volumes)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AutoPCR_GradientAnalysis(_ctx context.Context, _input *AutoPCR_GradientInput, _output *AutoPCR_GradientOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AutoPCR_GradientValidation(_ctx context.Context, _input *AutoPCR_GradientInput, _output *AutoPCR_GradientOutput) {
}
func _AutoPCR_GradientRun(_ctx context.Context, input *AutoPCR_GradientInput) *AutoPCR_GradientOutput {
	output := &AutoPCR_GradientOutput{}
	_AutoPCR_GradientSetup(_ctx, input)
	_AutoPCR_GradientSteps(_ctx, input, output)
	_AutoPCR_GradientAnalysis(_ctx, input, output)
	_AutoPCR_GradientValidation(_ctx, input, output)
	return output
}

func AutoPCR_GradientRunSteps(_ctx context.Context, input *AutoPCR_GradientInput) *AutoPCR_GradientSOutput {
	soutput := &AutoPCR_GradientSOutput{}
	output := _AutoPCR_GradientRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AutoPCR_GradientNew() interface{} {
	return &AutoPCR_GradientElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AutoPCR_GradientInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AutoPCR_GradientRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AutoPCR_GradientInput{},
			Out: &AutoPCR_GradientOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type AutoPCR_GradientElement struct {
	inject.CheckedRunner
}

type AutoPCR_GradientInput struct {
	FwdPrimertype               *wtype.LHComponent
	Plate                       *wtype.LHPlate
	Projectname                 string
	Reactiontoprimerpair        map[string][]string
	Reactiontotemplate          map[string]string
	RevPrimertype               *wtype.LHComponent
	RowGradientRatherthanColumn bool
	Templatetype                *wtype.LHComponent
}

type AutoPCR_GradientOutput struct {
	Error             error
	NumberOfReactions int
	Reactions         []*wtype.LHComponent
}

type AutoPCR_GradientSOutput struct {
	Data struct {
		Error             error
		NumberOfReactions int
	}
	Outputs struct {
		Reactions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AutoPCR_Gradient",
		Constructor: AutoPCR_GradientNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/PCR/AutoGradientPCR.an",
			Params: []component.ParamDesc{
				{Name: "FwdPrimertype", Desc: "", Kind: "Inputs"},
				{Name: "Plate", Desc: "", Kind: "Inputs"},
				{Name: "Projectname", Desc: "PCRprep parameters\n", Kind: "Parameters"},
				{Name: "Reactiontoprimerpair", Desc: "e.g. [\"left homology arm\"]:\"fwdprimer\",\"revprimer\"\n", Kind: "Parameters"},
				{Name: "Reactiontotemplate", Desc: "e.g. [\"left homology arm\"]:\"templatename\"\n", Kind: "Parameters"},
				{Name: "RevPrimertype", Desc: "", Kind: "Inputs"},
				{Name: "RowGradientRatherthanColumn", Desc: "if true, 12 replicates of each reaction will be set up, one set of reactions per row, else 8 reactions set up 1 set per column\n", Kind: "Parameters"},
				{Name: "Templatetype", Desc: "", Kind: "Inputs"},
				{Name: "Error", Desc: "", Kind: "Data"},
				{Name: "NumberOfReactions", Desc: "", Kind: "Data"},
				{Name: "Reactions", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
