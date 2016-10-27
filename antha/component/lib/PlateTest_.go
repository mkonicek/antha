package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/factory"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _PlateTestRequirements() {

}

// Conditions to run on startup
func _PlateTestSetup(_ctx context.Context, _input *PlateTestInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _PlateTestSteps(_ctx context.Context, _input *PlateTestInput, _output *PlateTestOutput) {

	_output.FinalSolutions = make([]*wtype.LHComponent, 0)

	_output.WellsUsedPostRunPerPlate = make([]int, 0)

	platelist := factory.GetPlateList()

	if _input.WellsUsedperPlateTypeInorder == nil || len(_input.WellsUsedperPlateTypeInorder) == 0 {
		_input.WellsUsedperPlateTypeInorder = make([]int, len(_input.OutPlates))
		for l := range _input.OutPlates {
			_input.WellsUsedperPlateTypeInorder[l] = 0
		}
	}

	for k := range _input.OutPlates {

		wellpositionsarray := factory.GetPlateByType(_input.OutPlates[k]).AllWellPositions(wtype.BYCOLUMN)

		counter := _input.WellsUsedperPlateTypeInorder[k]

		for j := range _input.LiquidVolumes {
			for i := range _input.LiquidTypes {

				liquidtypestring, err := wtype.LiquidTypeFromString(_input.LiquidTypes[i])

				if err != nil {
					execute.Errorf(_ctx, "Liquid type issue with ", _input.LiquidTypes[i], err.Error())
				}

				_input.Startingsolution.Type = liquidtypestring

				sample := mixer.Sample(_input.Startingsolution, _input.LiquidVolumes[j])

				if !search.InSlice(_input.OutPlates[k], platelist) {
					execute.Errorf(_ctx, "No plate ", _input.OutPlates[k], " found in library ", platelist)
				}

				finalSolution := execute.MixNamed(_ctx, _input.OutPlates[k], wellpositionsarray[counter], _input.OutPlates[k], sample)
				_output.FinalSolutions = append(_output.FinalSolutions, finalSolution)

				_output.Status = _output.Status + fmt.Sprintln(_input.LiquidVolumes[j].ToString(), " of ", _input.Liquidname, "Liquid type ", _input.LiquidTypes[i], "was mixed into "+_input.OutPlates[k])

				counter++

			}
		}

		_output.WellsUsedPostRunPerPlate = append(_output.WellsUsedPostRunPerPlate, counter)

	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PlateTestAnalysis(_ctx context.Context, _input *PlateTestInput, _output *PlateTestOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PlateTestValidation(_ctx context.Context, _input *PlateTestInput, _output *PlateTestOutput) {

}
func _PlateTestRun(_ctx context.Context, input *PlateTestInput) *PlateTestOutput {
	output := &PlateTestOutput{}
	_PlateTestSetup(_ctx, input)
	_PlateTestSteps(_ctx, input, output)
	_PlateTestAnalysis(_ctx, input, output)
	_PlateTestValidation(_ctx, input, output)
	return output
}

func PlateTestRunSteps(_ctx context.Context, input *PlateTestInput) *PlateTestSOutput {
	soutput := &PlateTestSOutput{}
	output := _PlateTestRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PlateTestNew() interface{} {
	return &PlateTestElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PlateTestInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PlateTestRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PlateTestInput{},
			Out: &PlateTestOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PlateTestElement struct {
	inject.CheckedRunner
}

type PlateTestInput struct {
	LiquidTypes                  []string
	LiquidVolumes                []wunit.Volume
	Liquidname                   string
	OutPlates                    []string
	Startingsolution             *wtype.LHComponent
	WellsUsedperPlateTypeInorder []int
}

type PlateTestOutput struct {
	FinalSolutions           []*wtype.LHComponent
	Status                   string
	WellsUsedPostRunPerPlate []int
}

type PlateTestSOutput struct {
	Data struct {
		Status                   string
		WellsUsedPostRunPerPlate []int
	}
	Outputs struct {
		FinalSolutions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PlateTest",
		Constructor: PlateTestNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Utility/PlateHeightTest.an",
			Params: []component.ParamDesc{
				{Name: "LiquidTypes", Desc: "", Kind: "Parameters"},
				{Name: "LiquidVolumes", Desc: "", Kind: "Parameters"},
				{Name: "Liquidname", Desc: "", Kind: "Parameters"},
				{Name: "OutPlates", Desc: "", Kind: "Parameters"},
				{Name: "Startingsolution", Desc: "", Kind: "Inputs"},
				{Name: "WellsUsedperPlateTypeInorder", Desc: "", Kind: "Parameters"},
				{Name: "FinalSolutions", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
				{Name: "WellsUsedPostRunPerPlate", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
