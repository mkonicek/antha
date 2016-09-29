package lib

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	//"github.com/antha-lang/antha/antha/anthalib/mixer"
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

func _AutoPCR_Gradient_SplitRequirements() {
}

// Conditions to run on startup
func _AutoPCR_Gradient_SplitSetup(_ctx context.Context, _input *AutoPCR_Gradient_SplitInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _AutoPCR_Gradient_SplitSteps(_ctx context.Context, _input *AutoPCR_Gradient_SplitInput, _output *AutoPCR_Gradient_SplitOutput) {
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
	var repeatcounter int
	var wellcoords = wtype.WellCoords{X: counter, Y: repeatcounter}

	Reactions := make([]*wtype.LHComponent, 0)
	_output.NewReactions = make([]*wtype.LHComponent, 0)
	volumes := make([]wunit.Volume, 0)
	welllocations := make([]string, 0)

	for reactionname, templatename := range _input.Reactiontotemplate {

		//wellposition := Plate.AllWellPositions(wtype.BYCOLUMN)[counter]

		//	for i:= 0;i < len(Reactions);i++{

		if _input.RowGradientRatherthanColumn {
			wellcoords = wtype.WellCoords{X: repeatcounter, Y: counter}
		} else {
			wellcoords = wtype.WellCoords{X: counter, Y: repeatcounter}
		}

		wellposition := wellcoords.FormatA1()
		result := PCR_volRunSteps(_ctx, &PCR_volInput{WaterVolume: wunit.NewVolume(10*float64(Samplenumber), "ul"),
			ReactionVolume:        wunit.NewVolume(25*float64(Samplenumber), "ul"),
			BufferConcinX:         5,
			FwdPrimerName:         _input.Reactiontoprimerpair[reactionname][0],
			RevPrimerName:         _input.Reactiontoprimerpair[reactionname][1],
			TemplateName:          templatename,
			ReactionName:          reactionname,
			FwdPrimerVol:          wunit.NewVolume(1*float64(Samplenumber), "ul"),
			RevPrimerVol:          wunit.NewVolume(1*float64(Samplenumber), "ul"),
			AdditiveVols:          []wunit.Volume{wunit.NewVolume(5*float64(Samplenumber), "ul")},
			Templatevolume:        wunit.NewVolume(1*float64(Samplenumber), "ul"),
			PolymeraseVolume:      wunit.NewVolume(1*float64(Samplenumber), "ul"),
			DNTPVol:               wunit.NewVolume(1*float64(Samplenumber), "ul"),
			Numberofcycles:        0,
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

		Reactions = append(Reactions, result.Outputs.Reaction)

		/*
		      volumes = append(volumes,result.Outputs.Reaction.Volume())
		   	welllocations = append(welllocations,wellposition)
		*/

		counter++
	}

	var byrow bool
	var startrow int
	var startcolumn int

	if !_input.RowGradientRatherthanColumn {
		byrow = true
		startrow = 0
		startcolumn = 0 //1
	} else {
		byrow = false
		startrow = 0 //1
		startcolumn = 0
	}

	aliquots := AliquotStartatRowColumnRunSteps(_ctx, &AliquotStartatRowColumnInput{SolutionVolume: wunit.NewVolume(25*float64(Samplenumber), "ul"),
		VolumePerAliquot: wunit.NewVolume(25, "ul"),
		NumberofAliquots: Samplenumber,
		StartRow:         startrow,
		StartColumn:      startcolumn,
		ByRow:            byrow,
		PreMix:           true,

		Solutions: Reactions,
		OutPlate:  _input.Plate},
	)

	for i, aliquot := range aliquots.Outputs.Aliquots {
		_output.NewReactions = append(_output.NewReactions, aliquot)
		volumes = append(volumes, aliquot.Volume())
		welllocations = append(welllocations, aliquots.Data.WellPositions[i])
	}

	/*
		// reset
		counter = 0
		// go to next
		repeatcounter++
		if RowGradientRatherthanColumn{
			wellcoords  = wtype.WellCoords{X:repeatcounter,Y:counter}
		} else {
			wellcoords  = wtype.WellCoords{X:counter,Y:repeatcounter}
		}



		for k, reaction := range Reactions {

		for l := repeatcounter; l < Samplenumber;l++ {
		reactionSample := mixer.Sample(reaction,wunit.NewVolume(25,"ul"))

		next := 	MixTo(Plate.Type,wellcoords.FormatA1(),1,reactionSample)
		Reactions = append(Reactions, next)
		volumes = append(volumes,wunit.NewVolume(25,"ul"))
		welllocations = append(welllocations,wellcoords.FormatA1())

		if RowGradientRatherthanColumn{
			wellcoords  = wtype.WellCoords{X:l,Y:k}
		} else {
			wellcoords  = wtype.WellCoords{X:k,Y:l}
		}

		}
	}*/
	_output.NumberOfReactions = len(Reactions)
	_output.Error = wtype.ExportPlateCSV(_input.Projectname+".csv", _input.Plate, _input.Projectname+"outputPlate", welllocations, Reactions, volumes)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AutoPCR_Gradient_SplitAnalysis(_ctx context.Context, _input *AutoPCR_Gradient_SplitInput, _output *AutoPCR_Gradient_SplitOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AutoPCR_Gradient_SplitValidation(_ctx context.Context, _input *AutoPCR_Gradient_SplitInput, _output *AutoPCR_Gradient_SplitOutput) {
}
func _AutoPCR_Gradient_SplitRun(_ctx context.Context, input *AutoPCR_Gradient_SplitInput) *AutoPCR_Gradient_SplitOutput {
	output := &AutoPCR_Gradient_SplitOutput{}
	_AutoPCR_Gradient_SplitSetup(_ctx, input)
	_AutoPCR_Gradient_SplitSteps(_ctx, input, output)
	_AutoPCR_Gradient_SplitAnalysis(_ctx, input, output)
	_AutoPCR_Gradient_SplitValidation(_ctx, input, output)
	return output
}

func AutoPCR_Gradient_SplitRunSteps(_ctx context.Context, input *AutoPCR_Gradient_SplitInput) *AutoPCR_Gradient_SplitSOutput {
	soutput := &AutoPCR_Gradient_SplitSOutput{}
	output := _AutoPCR_Gradient_SplitRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AutoPCR_Gradient_SplitNew() interface{} {
	return &AutoPCR_Gradient_SplitElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AutoPCR_Gradient_SplitInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AutoPCR_Gradient_SplitRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AutoPCR_Gradient_SplitInput{},
			Out: &AutoPCR_Gradient_SplitOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type AutoPCR_Gradient_SplitElement struct {
	inject.CheckedRunner
}

type AutoPCR_Gradient_SplitInput struct {
	FwdPrimertype               *wtype.LHComponent
	Plate                       *wtype.LHPlate
	Projectname                 string
	Reactiontoprimerpair        map[string][]string
	Reactiontotemplate          map[string]string
	RevPrimertype               *wtype.LHComponent
	RowGradientRatherthanColumn bool
	Templatetype                *wtype.LHComponent
}

type AutoPCR_Gradient_SplitOutput struct {
	Error             error
	NewReactions      []*wtype.LHComponent
	NumberOfReactions int
}

type AutoPCR_Gradient_SplitSOutput struct {
	Data struct {
		Error             error
		NumberOfReactions int
	}
	Outputs struct {
		NewReactions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AutoPCR_Gradient_Split",
		Constructor: AutoPCR_Gradient_SplitNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/PCR/AutoGradientPCR_Split.an",
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
				{Name: "NewReactions", Desc: "", Kind: "Outputs"},
				{Name: "NumberOfReactions", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
