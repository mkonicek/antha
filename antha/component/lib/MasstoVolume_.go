// example of how to convert a density and mass to a volume
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

func _MasstoVolumeRequirements() {

}

// Actions to perform before protocol itself
func _MasstoVolumeSetup(_ctx context.Context, _input *MasstoVolumeInput) {

}

// Core process of the protocol: steps to be performed for each input
func _MasstoVolumeSteps(_ctx context.Context, _input *MasstoVolumeInput, _output *MasstoVolumeOutput) {

	_output.Vol = wunit.MasstoVolume(_input.MyMass, _input.MyDensity)

	_output.BacktoMass = wunit.VolumetoMass(_output.Vol, _input.MyDensity)
}

// Actions to perform after steps block to analyze data
func _MasstoVolumeAnalysis(_ctx context.Context, _input *MasstoVolumeInput, _output *MasstoVolumeOutput) {

}

func _MasstoVolumeValidation(_ctx context.Context, _input *MasstoVolumeInput, _output *MasstoVolumeOutput) {

}
func _MasstoVolumeRun(_ctx context.Context, input *MasstoVolumeInput) *MasstoVolumeOutput {
	output := &MasstoVolumeOutput{}
	_MasstoVolumeSetup(_ctx, input)
	_MasstoVolumeSteps(_ctx, input, output)
	_MasstoVolumeAnalysis(_ctx, input, output)
	_MasstoVolumeValidation(_ctx, input, output)
	return output
}

func MasstoVolumeRunSteps(_ctx context.Context, input *MasstoVolumeInput) *MasstoVolumeSOutput {
	soutput := &MasstoVolumeSOutput{}
	output := _MasstoVolumeRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MasstoVolumeNew() interface{} {
	return &MasstoVolumeElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MasstoVolumeInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MasstoVolumeRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MasstoVolumeInput{},
			Out: &MasstoVolumeOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type MasstoVolumeElement struct {
	inject.CheckedRunner
}

type MasstoVolumeInput struct {
	MyDensity wunit.Density
	MyMass    wunit.Mass
}

type MasstoVolumeOutput struct {
	BacktoMass wunit.Mass
	Vol        wunit.Volume
}

type MasstoVolumeSOutput struct {
	Data struct {
		BacktoMass wunit.Mass
		Vol        wunit.Volume
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MasstoVolume",
		Constructor: MasstoVolumeNew,
		Desc: component.ComponentDesc{
			Desc: "example of how to convert a density and mass to a volume\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/AnthaAcademy/Lesson5_Units2/B_MasstoVolume.an",
			Params: []component.ParamDesc{
				{Name: "MyDensity", Desc: "", Kind: "Parameters"},
				{Name: "MyMass", Desc: "", Kind: "Parameters"},
				{Name: "BacktoMass", Desc: "", Kind: "Data"},
				{Name: "Vol", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
