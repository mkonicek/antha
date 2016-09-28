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

func _AutoPCR_Gradient_mmx_multiplereactionsRequirements() {
}

// Conditions to run on startup
func _AutoPCR_Gradient_mmx_multiplereactionsSetup(_ctx context.Context, _input *AutoPCR_Gradient_mmx_multiplereactionsInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _AutoPCR_Gradient_mmx_multiplereactionsSteps(_ctx context.Context, _input *AutoPCR_Gradient_mmx_multiplereactionsInput, _output *AutoPCR_Gradient_mmx_multiplereactionsOutput) {
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

	/*
	   // add step to make mastermix first
	   mastermix := RunSteps(MakePCRmmx,
	                  Parameters{
	                       WaterVolume: wunit.NewVolume(10,"ul"),
	   					ReactionVolume: wunit.NewVolume(25,"ul"),
	             BufferConcinX: 5,
	       FwdPrimerName: Reactiontoprimerpair[reactionname][0],
	       RevPrimerName: Reactiontoprimerpair[reactionname][1],
	       	TemplateName: templatename,
	       ReactionName: reactionname,
	             FwdPrimerVol: wunit.NewVolume(1,"ul"),
	             RevPrimerVol: wunit.NewVolume(1,"ul"),
	             AdditiveVols: []wunit.Volume{wunit.NewVolume(5,"ul")},
	             Templatevolume: wunit.NewVolume(1,"ul"),
	             PolymeraseVolume: wunit.NewVolume(1,"ul"),
	             DNTPVol:wunit.NewVolume(1,"ul"),
	             Numberofcycles: 30 ,
	             InitDenaturationtime: wunit.NewTime(30,"s"),
	             Denaturationtime: wunit.NewTime(5,"s"),
	             Annealingtime: wunit.NewTime(10,"s"),
	             AnnealingTemp: wunit.NewTemperature(72,"C"), // Should be calculated from primer and template binding
	             Extensiontime: wunit.NewTime(60,"s"), // should be calculated from template length and polymerase rate
	             Finalextensiontime: wunit.NewTime(180,"s"),
	             Hotstart: false,
	             AddPrimerstoMasterMix: false,
	   		WellPosition: wellposition,
	                   }, Inputs{
	             FwdPrimer:FwdPrimertype,
	             RevPrimer: RevPrimertype,
	             DNTPS: factory.GetComponentByType("DNTPs") ,
	             PCRPolymerase:factory.GetComponentByType("Q5Polymerase"),
	             Buffer:factory.GetComponentByType("Q5buffer"),
	             Water:factory.GetComponentByType("water"),
	             Template: Templatetype,
	             Additives: []*wtype.LHComponent{factory.GetComponentByType("GCenhancer")} ,
	             OutPlate: Plate,

	                   })

	*/

	for reactionname, templatename := range _input.Reactiontotemplate {

		//wellposition := Plate.AllWellPositions(wtype.BYCOLUMN)[counter]

		wellpositions := make([]string, 0)

		for j := 0; j < Samplenumber; j++ {
			//	for i:= 0;i < len(Reactions);i++{

			var wellcoords = wtype.WellCoords{X: counter, Y: j}

			if _input.RowGradientRatherthanColumn {
				wellcoords = wtype.WellCoords{X: j, Y: counter}
			} else {
				wellcoords = wtype.WellCoords{X: counter, Y: j}
			}

			wellpositions = append(wellpositions, wellcoords.FormatA1())

		}

		result := PCR_vol_mmx_multiplereactionsRunSteps(_ctx, &PCR_vol_mmx_multiplereactionsInput{

			//       WaterVolume: wunit.NewVolume(10,"ul"),
			MasterMixVolume:                   wunit.NewVolume(17, "ul"),
			FwdPrimerName:                     _input.Reactiontoprimerpair[reactionname][0],
			RevPrimerName:                     _input.Reactiontoprimerpair[reactionname][1],
			TemplateName:                      templatename,
			ReactionName:                      reactionname,
			FwdPrimerVol:                      wunit.NewVolume(1, "ul"),
			RevPrimerVol:                      wunit.NewVolume(1, "ul"),
			Templatevolume:                    wunit.NewVolume(1, "ul"),
			PolymeraseVolume:                  wunit.NewVolume(1, "ul"),
			Numberofcycles:                    30,
			InitDenaturationtime:              wunit.NewTime(30, "s"),
			Denaturationtime:                  wunit.NewTime(5, "s"),
			Annealingtime:                     wunit.NewTime(10, "s"),
			AnnealingTemp:                     wunit.NewTemperature(72, "C"), // Should be calculated from primer and template binding
			Extensiontime:                     wunit.NewTime(60, "s"),        // should be calculated from template length and polymerase rate
			Finalextensiontime:                wunit.NewTime(180, "s"),
			PrimersalreadyAddedtoMasterMix:    false,
			PolymeraseAlreadyaddedtoMastermix: true,
			WellPositions:                     wellpositions,

			FwdPrimer:     _input.FwdPrimertype,
			RevPrimer:     _input.RevPrimertype,
			PCRPolymerase: factory.GetComponentByType("Q5Polymerase"),
			Template:      _input.Templatetype,
			OutPlate:      _input.Plate,
			MasterMix:     factory.GetComponentByType("Q5mastermix")},
		)

		for m, reaction := range result.Outputs.Reactions {
			_output.Reactions = append(_output.Reactions, reaction)
			volumes = append(volumes, reaction.Volume())
			welllocations = append(welllocations, wellpositions[m])
		}

		counter++
	}
	_output.NumberOfReactions = len(_output.Reactions)
	_output.Error = wtype.ExportPlateCSV(_input.Projectname+".csv", _input.Plate, _input.Projectname+"outputPlate", welllocations, _output.Reactions, volumes)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AutoPCR_Gradient_mmx_multiplereactionsAnalysis(_ctx context.Context, _input *AutoPCR_Gradient_mmx_multiplereactionsInput, _output *AutoPCR_Gradient_mmx_multiplereactionsOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AutoPCR_Gradient_mmx_multiplereactionsValidation(_ctx context.Context, _input *AutoPCR_Gradient_mmx_multiplereactionsInput, _output *AutoPCR_Gradient_mmx_multiplereactionsOutput) {
}
func _AutoPCR_Gradient_mmx_multiplereactionsRun(_ctx context.Context, input *AutoPCR_Gradient_mmx_multiplereactionsInput) *AutoPCR_Gradient_mmx_multiplereactionsOutput {
	output := &AutoPCR_Gradient_mmx_multiplereactionsOutput{}
	_AutoPCR_Gradient_mmx_multiplereactionsSetup(_ctx, input)
	_AutoPCR_Gradient_mmx_multiplereactionsSteps(_ctx, input, output)
	_AutoPCR_Gradient_mmx_multiplereactionsAnalysis(_ctx, input, output)
	_AutoPCR_Gradient_mmx_multiplereactionsValidation(_ctx, input, output)
	return output
}

func AutoPCR_Gradient_mmx_multiplereactionsRunSteps(_ctx context.Context, input *AutoPCR_Gradient_mmx_multiplereactionsInput) *AutoPCR_Gradient_mmx_multiplereactionsSOutput {
	soutput := &AutoPCR_Gradient_mmx_multiplereactionsSOutput{}
	output := _AutoPCR_Gradient_mmx_multiplereactionsRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AutoPCR_Gradient_mmx_multiplereactionsNew() interface{} {
	return &AutoPCR_Gradient_mmx_multiplereactionsElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AutoPCR_Gradient_mmx_multiplereactionsInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AutoPCR_Gradient_mmx_multiplereactionsRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AutoPCR_Gradient_mmx_multiplereactionsInput{},
			Out: &AutoPCR_Gradient_mmx_multiplereactionsOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type AutoPCR_Gradient_mmx_multiplereactionsElement struct {
	inject.CheckedRunner
}

type AutoPCR_Gradient_mmx_multiplereactionsInput struct {
	FwdPrimertype               *wtype.LHComponent
	Plate                       *wtype.LHPlate
	Projectname                 string
	Reactiontoprimerpair        map[string][]string
	Reactiontotemplate          map[string]string
	RevPrimertype               *wtype.LHComponent
	RowGradientRatherthanColumn bool
	Templatetype                *wtype.LHComponent
}

type AutoPCR_Gradient_mmx_multiplereactionsOutput struct {
	Error             error
	NumberOfReactions int
	Reactions         []*wtype.LHComponent
}

type AutoPCR_Gradient_mmx_multiplereactionsSOutput struct {
	Data struct {
		Error             error
		NumberOfReactions int
	}
	Outputs struct {
		Reactions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AutoPCR_Gradient_mmx_multiplereactions",
		Constructor: AutoPCR_Gradient_mmx_multiplereactionsNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/PCR/AutoGradientPCRmmx_multiplereactions.an",
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
