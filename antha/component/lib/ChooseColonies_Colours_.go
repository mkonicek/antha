package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image/pick"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
	"image/color"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _ChooseColonies_ColoursRequirements() {

}

// Conditions to run on startup
func _ChooseColonies_ColoursSetup(_ctx context.Context, _input *ChooseColonies_ColoursInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _ChooseColonies_ColoursSteps(_ctx context.Context, _input *ChooseColonies_ColoursInput, _output *ChooseColonies_ColoursOutput) {

	_output.Wellstopick, _output.MyColourPaletteMap, _output.Error = pick.PickAndExportWelltoColourJSON(_input.Imagefile, _input.ExportFileName, _input.PlateForCoordinates, _input.NumbertoPick, _input.Setplateperimeterfirst, _input.Rotate)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _ChooseColonies_ColoursAnalysis(_ctx context.Context, _input *ChooseColonies_ColoursInput, _output *ChooseColonies_ColoursOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _ChooseColonies_ColoursValidation(_ctx context.Context, _input *ChooseColonies_ColoursInput, _output *ChooseColonies_ColoursOutput) {

}
func _ChooseColonies_ColoursRun(_ctx context.Context, input *ChooseColonies_ColoursInput) *ChooseColonies_ColoursOutput {
	output := &ChooseColonies_ColoursOutput{}
	_ChooseColonies_ColoursSetup(_ctx, input)
	_ChooseColonies_ColoursSteps(_ctx, input, output)
	_ChooseColonies_ColoursAnalysis(_ctx, input, output)
	_ChooseColonies_ColoursValidation(_ctx, input, output)
	return output
}

func ChooseColonies_ColoursRunSteps(_ctx context.Context, input *ChooseColonies_ColoursInput) *ChooseColonies_ColoursSOutput {
	soutput := &ChooseColonies_ColoursSOutput{}
	output := _ChooseColonies_ColoursRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ChooseColonies_ColoursNew() interface{} {
	return &ChooseColonies_ColoursElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ChooseColonies_ColoursInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ChooseColonies_ColoursRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ChooseColonies_ColoursInput{},
			Out: &ChooseColonies_ColoursOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type ChooseColonies_ColoursElement struct {
	inject.CheckedRunner
}

type ChooseColonies_ColoursInput struct {
	ExportFileName         string
	Imagefile              string
	NumbertoPick           int
	PlateForCoordinates    *wtype.LHPlate
	Rotate                 bool
	Setplateperimeterfirst bool
}

type ChooseColonies_ColoursOutput struct {
	Error              error
	MyColourPaletteMap map[string]color.Color
	Wellstopick        []string
}

type ChooseColonies_ColoursSOutput struct {
	Data struct {
		Error              error
		MyColourPaletteMap map[string]color.Color
		Wellstopick        []string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "ChooseColonies_Colours",
		Constructor: ChooseColonies_ColoursNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Data/choosecolonies/ChooseColonies_Colours.an",
			Params: []component.ParamDesc{
				{Name: "ExportFileName", Desc: "", Kind: "Parameters"},
				{Name: "Imagefile", Desc: "", Kind: "Parameters"},
				{Name: "NumbertoPick", Desc: "", Kind: "Parameters"},
				{Name: "PlateForCoordinates", Desc: "", Kind: "Inputs"},
				{Name: "Rotate", Desc: "", Kind: "Parameters"},
				{Name: "Setplateperimeterfirst", Desc: "", Kind: "Parameters"},
				{Name: "Error", Desc: "", Kind: "Data"},
				{Name: "MyColourPaletteMap", Desc: "", Kind: "Data"},
				{Name: "Wellstopick", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
