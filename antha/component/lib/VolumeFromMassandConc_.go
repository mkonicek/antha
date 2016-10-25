// example of how to convert a concentration and mass to a volume
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

func _VolumeFromMassandConcRequirements() {

}

// Actions to perform before protocol itself
func _VolumeFromMassandConcSetup(_ctx context.Context, _input *VolumeFromMassandConcInput) {

}

// Core process of the protocol: steps to be performed for each input
func _VolumeFromMassandConcSteps(_ctx context.Context, _input *VolumeFromMassandConcInput, _output *VolumeFromMassandConcOutput) {
	_output.DNAVol, _output.Err = wunit.VolumeForTargetMass(_input.DNAMassperReaction, _input.DNAConc)
	if _output.Err != nil {
		execute.Errorf(_ctx, _output.Err.Error())
	}
}

// Actions to perform after steps block to analyze data
func _VolumeFromMassandConcAnalysis(_ctx context.Context, _input *VolumeFromMassandConcInput, _output *VolumeFromMassandConcOutput) {

}

func _VolumeFromMassandConcValidation(_ctx context.Context, _input *VolumeFromMassandConcInput, _output *VolumeFromMassandConcOutput) {

}
func _VolumeFromMassandConcRun(_ctx context.Context, input *VolumeFromMassandConcInput) *VolumeFromMassandConcOutput {
	output := &VolumeFromMassandConcOutput{}
	_VolumeFromMassandConcSetup(_ctx, input)
	_VolumeFromMassandConcSteps(_ctx, input, output)
	_VolumeFromMassandConcAnalysis(_ctx, input, output)
	_VolumeFromMassandConcValidation(_ctx, input, output)
	return output
}

func VolumeFromMassandConcRunSteps(_ctx context.Context, input *VolumeFromMassandConcInput) *VolumeFromMassandConcSOutput {
	soutput := &VolumeFromMassandConcSOutput{}
	output := _VolumeFromMassandConcRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func VolumeFromMassandConcNew() interface{} {
	return &VolumeFromMassandConcElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &VolumeFromMassandConcInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _VolumeFromMassandConcRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &VolumeFromMassandConcInput{},
			Out: &VolumeFromMassandConcOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type VolumeFromMassandConcElement struct {
	inject.CheckedRunner
}

type VolumeFromMassandConcInput struct {
	DNAConc            wunit.Concentration
	DNAMassperReaction wunit.Mass
}

type VolumeFromMassandConcOutput struct {
	DNAVol wunit.Volume
	Err    error
}

type VolumeFromMassandConcSOutput struct {
	Data struct {
		DNAVol wunit.Volume
		Err    error
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "VolumeFromMassandConc",
		Constructor: VolumeFromMassandConcNew,
		Desc: component.ComponentDesc{
			Desc: "example of how to convert a concentration and mass to a volume\n",
			Path: "antha/component/an/AnthaAcademy/Lesson5_Units2/C_VolumefromMassandConc.an",
			Params: []component.ParamDesc{
				{Name: "DNAConc", Desc: "", Kind: "Parameters"},
				{Name: "DNAMassperReaction", Desc: "", Kind: "Parameters"},
				{Name: "DNAVol", Desc: "", Kind: "Data"},
				{Name: "Err", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
