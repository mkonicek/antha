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

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _PlateOutRequirements() {
}

// Conditions to run on startup
func _PlateOutSetup(_ctx context.Context, _input *PlateOutInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _PlateOutSteps(_ctx context.Context, _input *PlateOutInput, _output *PlateOutOutput) {

	plateout := make([]*wtype.LHComponent, 0)

	if _input.Diluent != nil && _input.DilutionX > 1 {
		diluentsample := mixer.SampleForTotalVolume(_input.Diluent, _input.Plateoutvolume)
		plateout = append(plateout, diluentsample)
		// redeclare Plateoutvolume for adjusted volume to add of recovery mixture based on dilution ratio
		_input.Plateoutvolume = wunit.NewVolume(_input.Plateoutvolume.RawValue()/float64(_input.DilutionX), _input.Plateoutvolume.Unit().PrefixedSymbol())

	}
	plateoutsample := mixer.Sample(_input.RecoveredCells, _input.Plateoutvolume)
	plateout = append(plateout, plateoutsample)
	platedculture := execute.MixInto(_ctx, _input.AgarPlate, "", plateout...)
	platedculture = execute.Incubate(_ctx, platedculture, _input.IncubationTemp, _input.IncubationTime, false)
	_output.Platedculture = platedculture

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PlateOutAnalysis(_ctx context.Context, _input *PlateOutInput, _output *PlateOutOutput) {

}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PlateOutValidation(_ctx context.Context, _input *PlateOutInput, _output *PlateOutOutput) {

}
func _PlateOutRun(_ctx context.Context, input *PlateOutInput) *PlateOutOutput {
	output := &PlateOutOutput{}
	_PlateOutSetup(_ctx, input)
	_PlateOutSteps(_ctx, input, output)
	_PlateOutAnalysis(_ctx, input, output)
	_PlateOutValidation(_ctx, input, output)
	return output
}

func PlateOutRunSteps(_ctx context.Context, input *PlateOutInput) *PlateOutSOutput {
	soutput := &PlateOutSOutput{}
	output := _PlateOutRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PlateOutNew() interface{} {
	return &PlateOutElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PlateOutInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PlateOutRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PlateOutInput{},
			Out: &PlateOutOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PlateOutElement struct {
	inject.CheckedRunner
}

type PlateOutInput struct {
	AgarPlate      *wtype.LHPlate
	Diluent        *wtype.LHComponent
	DilutionX      int
	IncubationTemp wunit.Temperature
	IncubationTime wunit.Time
	Plateoutvolume wunit.Volume
	RecoveredCells *wtype.LHComponent
}

type PlateOutOutput struct {
	Platedculture *wtype.LHComponent
}

type PlateOutSOutput struct {
	Data struct {
	}
	Outputs struct {
		Platedculture *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PlateOut",
		Constructor: PlateOutNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/Transformation/Plateout.an",
			Params: []component.ParamDesc{
				{Name: "AgarPlate", Desc: "", Kind: "Inputs"},
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "DilutionX", Desc: "", Kind: "Parameters"},
				{Name: "IncubationTemp", Desc: "", Kind: "Parameters"},
				{Name: "IncubationTime", Desc: "", Kind: "Parameters"},
				{Name: "Plateoutvolume", Desc: "", Kind: "Parameters"},
				{Name: "RecoveredCells", Desc: "", Kind: "Inputs"},
				{Name: "Platedculture", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
