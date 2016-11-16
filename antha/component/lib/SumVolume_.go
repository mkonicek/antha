package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

//"github.com/antha-lang/antha/antha/anthalib/wunit"

// Input parameters for this protocol

//D Concentration
//E float64

// Data which is returned from this protocol

//DmolarConc wunit.MolarConcentration

// Physical inputs to this protocol

// Physical outputs from this protocol

func _SumVolumeRequirements() {

}

// Actions to perform before protocol itself
func _SumVolumeSetup(_ctx context.Context, _input *SumVolumeInput) {

}

// Core process of the protocol: steps to be performed for each input
func _SumVolumeSteps(_ctx context.Context, _input *SumVolumeInput, _output *SumVolumeOutput) {
	//var Dmassconc wunit.MassConcentration = D

	/*	molarmass := wunit.NewAmount(E,"M")

		var Dnew = wunit.MoleculeConcentration{D,E}

		mass := wunit.NewMass(1,"g")

		DmolarConc = Dnew.AsMolar(mass)
	*/
	_output.Sum = (wunit.CopyVolume(_input.A))
	_output.Sum.Add(_input.B)
	_output.Status = fmt.Sprintln(
		"Sum of", _input.A.ToString(), "and", _input.B.ToString(), "=", _output.Sum.ToString(), "Temp=", _input.C.ToString(),
	) //"D Concentration in g/l", D, "D concentration in M/l", DmolarConc)
}

// Actions to perform after steps block to analyze data
func _SumVolumeAnalysis(_ctx context.Context, _input *SumVolumeInput, _output *SumVolumeOutput) {

}

func _SumVolumeValidation(_ctx context.Context, _input *SumVolumeInput, _output *SumVolumeOutput) {

}
func _SumVolumeRun(_ctx context.Context, input *SumVolumeInput) *SumVolumeOutput {
	output := &SumVolumeOutput{}
	_SumVolumeSetup(_ctx, input)
	_SumVolumeSteps(_ctx, input, output)
	_SumVolumeAnalysis(_ctx, input, output)
	_SumVolumeValidation(_ctx, input, output)
	return output
}

func SumVolumeRunSteps(_ctx context.Context, input *SumVolumeInput) *SumVolumeSOutput {
	soutput := &SumVolumeSOutput{}
	output := _SumVolumeRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SumVolumeNew() interface{} {
	return &SumVolumeElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SumVolumeInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SumVolumeRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SumVolumeInput{},
			Out: &SumVolumeOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type SumVolumeElement struct {
	inject.CheckedRunner
}

type SumVolumeInput struct {
	A wunit.Volume
	B wunit.Volume
	C wunit.Temperature
}

type SumVolumeOutput struct {
	Status string
	Sum    wunit.Volume
}

type SumVolumeSOutput struct {
	Data struct {
		Status string
		Sum    wunit.Volume
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "SumVolume",
		Constructor: SumVolumeNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/sumVolume/Sum.an",
			Params: []component.ParamDesc{
				{Name: "A", Desc: "", Kind: "Parameters"},
				{Name: "B", Desc: "", Kind: "Parameters"},
				{Name: "C", Desc: "", Kind: "Parameters"},
				{Name: "Status", Desc: "", Kind: "Data"},
				{Name: "Sum", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
