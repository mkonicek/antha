package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/factory"
	inplate "github.com/antha-lang/antha/target/mixer"
	"golang.org/x/net/context"
	"strconv"
)

// Input parameters for this protocol (data)

// PCRprep parameters

// e.g. ["left homology arm"]:"fwdprimer","revprimer"

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _AutoColonyPCRRequirements() {
}

// Conditions to run on startup
func _AutoColonyPCRSetup(_ctx context.Context, _input *AutoColonyPCRInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _AutoColonyPCRSteps(_ctx context.Context, _input *AutoColonyPCRInput, _output *AutoColonyPCROutput) {

	// make map to cross reference starting well to destination well
	_output.OriginalWelltoNewWell = make(map[string]string)

	// initialise some counters
	platenum := 1
	counter := 0

	// set low volume of colony to pick
	colonypickvol := wunit.NewVolume(1.0, "ul")

	// stuff we need later
	//wellpositionarray := make([]string,0)
	numberofcolonies := 0

	// parse colony locations from file
	inputplate, err := inplate.ParseInputPlateFile(_input.WellstopickCSV)

	if err != nil {
		execute.Errorf(_ctx, "Error parsing inputplate csv file")
	}

	// count number of colonies from inplate

	for _, wellcontents := range inputplate.Wellcoords {
		if wellcontents.Empty() == false {
			numberofcolonies = numberofcolonies + 1
		}
	}

	// reset before adding colonies
	platenum = 1
	counter = 0

	_output.Reactions = make([]*wtype.LHComponent, 0)
	volumes := make([]wunit.Volume, 0)
	welllocations := make([]string, 0)

	// add colonies
	for originalwell, wellcontents := range inputplate.Wellcoords {

		if wellcontents.Empty() == false {

			if counter == (_input.Plate.WlsX * _input.Plate. /*+NumberofBlanks*/ WlsY) {
				fmt.Println("plate full, counter = ", counter)
				platenum++
				//reset counter
				counter = 0
			}

			colonyComponent := wellcontents.WContents

			wellposition := _input.Plate.AllWellPositions(wtype.BYCOLUMN)[counter]

			result := ColonyPCR_vol_mmxRunSteps(_ctx, &ColonyPCR_vol_mmxInput{RecoveryWaterVolume: wunit.NewVolume(10, "ul"),
				MasterMixVolume:                   wunit.NewVolume(17, "ul"),
				FwdPrimerName:                     _input.Reactiontoprimerpair[wellcontents.WContents.CName][0],
				RevPrimerName:                     _input.Reactiontoprimerpair[wellcontents.WContents.CName][1],
				TemplateName:                      wellcontents.WContents.CName, //+"_"+inputplate.PlateName+"_"+originalwell, // colony starting well position
				ReactionName:                      wellcontents.WContents.CName + "_" + inputplate.PlateName + "_" + originalwell,
				PolymeraseVolume:                  wunit.NewVolume(0, "ul"),
				FwdPrimerVol:                      wunit.NewVolume(1, "ul"),
				RevPrimerVol:                      wunit.NewVolume(1, "ul"),
				Templatevolume:                    colonypickvol,
				Numberofcycles:                    30,
				InitDenaturationtime:              wunit.NewTime(30, "s"),
				Denaturationtime:                  wunit.NewTime(5, "s"),
				Annealingtime:                     wunit.NewTime(10, "s"),
				AnnealingTemp:                     wunit.NewTemperature(72, "C"), // Should be calculated from primer and template binding
				Extensiontime:                     wunit.NewTime(60, "s"),        // should be calculated from template length and polymerase rate
				Finalextensiontime:                wunit.NewTime(180, "s"),
				WellPosition:                      wellposition,
				PolymeraseAlreadyaddedtoMastermix: true,

				FwdPrimer:     _input.FwdPrimertype,
				RevPrimer:     _input.RevPrimertype,
				MasterMix:     factory.GetComponentByType("Q5mastermix"),
				PCRPolymerase: factory.GetComponentByType("Q5Polymerase"),
				RecoveryWater: factory.GetComponentByType("water"),
				Template:      colonyComponent,
				OutPlate:      _input.Plate,
				RecoveryPlate: _input.RecoveryPlate},
			)

			_output.Reactions = append(_output.Reactions, result.Outputs.Reaction)
			volumes = append(volumes, result.Outputs.Reaction.Volume())
			welllocations = append(welllocations, wellposition)

			// add info to output map for well location cross referencing
			_output.OriginalWelltoNewWell[inputplate.PlateName+"_"+originalwell] = strconv.Itoa(platenum) + wellposition

			counter++

		}

	}
	_output.NumberofColonies = numberofcolonies

	_output.Error = wtype.ExportPlateCSV(_input.Projectname+".csv", _input.Plate, _input.Projectname+"outputPlate", welllocations, _output.Reactions, volumes)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AutoColonyPCRAnalysis(_ctx context.Context, _input *AutoColonyPCRInput, _output *AutoColonyPCROutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AutoColonyPCRValidation(_ctx context.Context, _input *AutoColonyPCRInput, _output *AutoColonyPCROutput) {
}
func _AutoColonyPCRRun(_ctx context.Context, input *AutoColonyPCRInput) *AutoColonyPCROutput {
	output := &AutoColonyPCROutput{}
	_AutoColonyPCRSetup(_ctx, input)
	_AutoColonyPCRSteps(_ctx, input, output)
	_AutoColonyPCRAnalysis(_ctx, input, output)
	_AutoColonyPCRValidation(_ctx, input, output)
	return output
}

func AutoColonyPCRRunSteps(_ctx context.Context, input *AutoColonyPCRInput) *AutoColonyPCRSOutput {
	soutput := &AutoColonyPCRSOutput{}
	output := _AutoColonyPCRRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AutoColonyPCRNew() interface{} {
	return &AutoColonyPCRElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AutoColonyPCRInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AutoColonyPCRRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AutoColonyPCRInput{},
			Out: &AutoColonyPCROutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type AutoColonyPCRElement struct {
	inject.CheckedRunner
}

type AutoColonyPCRInput struct {
	FwdPrimertype        *wtype.LHComponent
	Plate                *wtype.LHPlate
	Projectname          string
	Reactiontoprimerpair map[string][]string
	RecoveryPlate        *wtype.LHPlate
	RevPrimertype        *wtype.LHComponent
	Templatetype         *wtype.LHComponent
	WellstopickCSV       string
}

type AutoColonyPCROutput struct {
	Error                 error
	NumberofColonies      int
	OriginalWelltoNewWell map[string]string
	Reactions             []*wtype.LHComponent
}

type AutoColonyPCRSOutput struct {
	Data struct {
		Error                 error
		NumberofColonies      int
		OriginalWelltoNewWell map[string]string
	}
	Outputs struct {
		Reactions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AutoColonyPCR",
		Constructor: AutoColonyPCRNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/PCR/AutoColonyPCR.an",
			Params: []component.ParamDesc{
				{Name: "FwdPrimertype", Desc: "", Kind: "Inputs"},
				{Name: "Plate", Desc: "", Kind: "Inputs"},
				{Name: "Projectname", Desc: "PCRprep parameters\n", Kind: "Parameters"},
				{Name: "Reactiontoprimerpair", Desc: "e.g. [\"left homology arm\"]:\"fwdprimer\",\"revprimer\"\n", Kind: "Parameters"},
				{Name: "RecoveryPlate", Desc: "", Kind: "Inputs"},
				{Name: "RevPrimertype", Desc: "", Kind: "Inputs"},
				{Name: "Templatetype", Desc: "", Kind: "Inputs"},
				{Name: "WellstopickCSV", Desc: "", Kind: "Parameters"},
				{Name: "Error", Desc: "", Kind: "Data"},
				{Name: "NumberofColonies", Desc: "", Kind: "Data"},
				{Name: "OriginalWelltoNewWell", Desc: "", Kind: "Data"},
				{Name: "Reactions", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
