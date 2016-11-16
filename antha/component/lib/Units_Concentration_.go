package lib

import (
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol

// Data which is returned from this protocol

// Physical inputs to this protocol

// Physical outputs from this protocol

func _Units_ConcentrationRequirements() {

}

// Actions to perform before protocol itself
func _Units_ConcentrationSetup(_ctx context.Context, _input *Units_ConcentrationInput) {

}

// Core process of the protocol: steps to be performed for each input
func _Units_ConcentrationSteps(_ctx context.Context, _input *Units_ConcentrationInput, _output *Units_ConcentrationOutput) {

	_output.ConcinMperL = _input.MyConc.MolPerL(_input.MolecularWeight)
	_output.ConcinGperL = _input.MyConc.GramPerL(_input.MolecularWeight)

}

// Actions to perform after steps block to analyze data
func _Units_ConcentrationAnalysis(_ctx context.Context, _input *Units_ConcentrationInput, _output *Units_ConcentrationOutput) {

}

func _Units_ConcentrationValidation(_ctx context.Context, _input *Units_ConcentrationInput, _output *Units_ConcentrationOutput) {

}
func _Units_ConcentrationRun(_ctx context.Context, input *Units_ConcentrationInput) *Units_ConcentrationOutput {
	output := &Units_ConcentrationOutput{}
	_Units_ConcentrationSetup(_ctx, input)
	_Units_ConcentrationSteps(_ctx, input, output)
	_Units_ConcentrationAnalysis(_ctx, input, output)
	_Units_ConcentrationValidation(_ctx, input, output)
	return output
}

func Units_ConcentrationRunSteps(_ctx context.Context, input *Units_ConcentrationInput) *Units_ConcentrationSOutput {
	soutput := &Units_ConcentrationSOutput{}
	output := _Units_ConcentrationRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Units_ConcentrationNew() interface{} {
	return &Units_ConcentrationElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Units_ConcentrationInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Units_ConcentrationRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Units_ConcentrationInput{},
			Out: &Units_ConcentrationOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Units_ConcentrationElement struct {
	inject.CheckedRunner
}

type Units_ConcentrationInput struct {
	MolecularWeight float64
	MyConc          wunit.Concentration
}

type Units_ConcentrationOutput struct {
	ConcinGperL wunit.Concentration
	ConcinMperL wunit.Concentration
}

type Units_ConcentrationSOutput struct {
	Data struct {
		ConcinGperL wunit.Concentration
		ConcinMperL wunit.Concentration
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Units_Concentration",
		Constructor: Units_ConcentrationNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/AnthaAcademy/Lesson0_Units/D_units_Conc.an",
			Params: []component.ParamDesc{
				{Name: "MolecularWeight", Desc: "", Kind: "Parameters"},
				{Name: "MyConc", Desc: "", Kind: "Parameters"},
				{Name: "ConcinGperL", Desc: "", Kind: "Data"},
				{Name: "ConcinMperL", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
