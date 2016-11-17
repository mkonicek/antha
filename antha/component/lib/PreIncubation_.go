package lib

import (
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

//= 50.(uL)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _PreIncubationRequirements() {
}

// Conditions to run on startup
func _PreIncubationSetup(_ctx context.Context, _input *PreIncubationInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _PreIncubationSteps(_ctx context.Context, _input *PreIncubationInput, _output *PreIncubationOutput) {
	competentcells := make([]*wtype.LHComponent, 0)
	competentcellsample := mixer.Sample(_input.CompetentCells, _input.CompetentCellvolumeperassembly)
	competentcells = append(competentcells, competentcellsample)
	readycompetentcells := execute.MixInto(_ctx, _input.OutPlate, "", competentcells...)

	_output.ReadyCompCells = execute.Incubate(_ctx, readycompetentcells, _input.Preplasmidtemp, _input.Preplasmidtime, false)
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PreIncubationAnalysis(_ctx context.Context, _input *PreIncubationInput, _output *PreIncubationOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PreIncubationValidation(_ctx context.Context, _input *PreIncubationInput, _output *PreIncubationOutput) {
}
func _PreIncubationRun(_ctx context.Context, input *PreIncubationInput) *PreIncubationOutput {
	output := &PreIncubationOutput{}
	_PreIncubationSetup(_ctx, input)
	_PreIncubationSteps(_ctx, input, output)
	_PreIncubationAnalysis(_ctx, input, output)
	_PreIncubationValidation(_ctx, input, output)
	return output
}

func PreIncubationRunSteps(_ctx context.Context, input *PreIncubationInput) *PreIncubationSOutput {
	soutput := &PreIncubationSOutput{}
	output := _PreIncubationRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PreIncubationNew() interface{} {
	return &PreIncubationElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PreIncubationInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PreIncubationRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PreIncubationInput{},
			Out: &PreIncubationOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PreIncubationElement struct {
	inject.CheckedRunner
}

type PreIncubationInput struct {
	CompetentCells                 *wtype.LHComponent
	CompetentCellvolumeperassembly wunit.Volume
	OutPlate                       *wtype.LHPlate
	Preplasmidtemp                 wunit.Temperature
	Preplasmidtime                 wunit.Time
}

type PreIncubationOutput struct {
	ReadyCompCells *wtype.LHComponent
}

type PreIncubationSOutput struct {
	Data struct {
	}
	Outputs struct {
		ReadyCompCells *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PreIncubation",
		Constructor: PreIncubationNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/Transformation/PreIncubation.an",
			Params: []component.ParamDesc{
				{Name: "CompetentCells", Desc: "", Kind: "Inputs"},
				{Name: "CompetentCellvolumeperassembly", Desc: "= 50.(uL)\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Preplasmidtemp", Desc: "", Kind: "Parameters"},
				{Name: "Preplasmidtime", Desc: "", Kind: "Parameters"},
				{Name: "ReadyCompCells", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
