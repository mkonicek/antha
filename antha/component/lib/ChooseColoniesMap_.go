package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image/pick"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _ChooseColoniesMapRequirements() {

}

// Conditions to run on startup
func _ChooseColoniesMapSetup(_ctx context.Context, _input *ChooseColoniesMapInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _ChooseColoniesMapSteps(_ctx context.Context, _input *ChooseColoniesMapInput, _output *ChooseColoniesMapOutput) {

	_output.Wellstopick, _output.Error = pick.PickAndExportCSVMap(_input.Imagefile, _input.ExportFileName, _input.PlateForCoordinates, _input.ReactiontoNumberMap, _input.Setplateperimeterfirst, _input.Rotate)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _ChooseColoniesMapAnalysis(_ctx context.Context, _input *ChooseColoniesMapInput, _output *ChooseColoniesMapOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _ChooseColoniesMapValidation(_ctx context.Context, _input *ChooseColoniesMapInput, _output *ChooseColoniesMapOutput) {

}
func _ChooseColoniesMapRun(_ctx context.Context, input *ChooseColoniesMapInput) *ChooseColoniesMapOutput {
	output := &ChooseColoniesMapOutput{}
	_ChooseColoniesMapSetup(_ctx, input)
	_ChooseColoniesMapSteps(_ctx, input, output)
	_ChooseColoniesMapAnalysis(_ctx, input, output)
	_ChooseColoniesMapValidation(_ctx, input, output)
	return output
}

func ChooseColoniesMapRunSteps(_ctx context.Context, input *ChooseColoniesMapInput) *ChooseColoniesMapSOutput {
	soutput := &ChooseColoniesMapSOutput{}
	output := _ChooseColoniesMapRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ChooseColoniesMapNew() interface{} {
	return &ChooseColoniesMapElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ChooseColoniesMapInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ChooseColoniesMapRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ChooseColoniesMapInput{},
			Out: &ChooseColoniesMapOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type ChooseColoniesMapElement struct {
	inject.CheckedRunner
}

type ChooseColoniesMapInput struct {
	ExportFileName         string
	Imagefile              string
	PlateForCoordinates    *wtype.LHPlate
	ReactiontoNumberMap    map[string]int
	Rotate                 bool
	Setplateperimeterfirst bool
}

type ChooseColoniesMapOutput struct {
	Error       error
	Wellstopick []string
}

type ChooseColoniesMapSOutput struct {
	Data struct {
		Error       error
		Wellstopick []string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "ChooseColoniesMap",
		Constructor: ChooseColoniesMapNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Data/choosecolonies/ChooseColoniesMap.an",
			Params: []component.ParamDesc{
				{Name: "ExportFileName", Desc: "", Kind: "Parameters"},
				{Name: "Imagefile", Desc: "", Kind: "Parameters"},
				{Name: "PlateForCoordinates", Desc: "", Kind: "Inputs"},
				{Name: "ReactiontoNumberMap", Desc: "", Kind: "Parameters"},
				{Name: "Rotate", Desc: "", Kind: "Parameters"},
				{Name: "Setplateperimeterfirst", Desc: "", Kind: "Parameters"},
				{Name: "Error", Desc: "", Kind: "Data"},
				{Name: "Wellstopick", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
