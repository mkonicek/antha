package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

func _InoculateDOERequirements() {
}

func _InoculateDOESetup(_ctx context.Context, _input *InoculateDOEInput) {
}

func _InoculateDOESteps(_ctx context.Context, _input *InoculateDOEInput, _output *InoculateDOEOutput) {
	inocsample := mixer.Sample(_input.Inoculum, _input.InoculumVolume)
	fmt.Println("Inoculum ", _input.Inoculum.CName, "Inoculum Volume ", _input.InoculumVolume.ToString(), "Medium", _input.Medium.CName)
	_output.Seed = execute.Mix(_ctx, _input.Medium, inocsample)
	fmt.Println("Seed:", _output.Seed.CName)
}

func _InoculateDOEAnalysis(_ctx context.Context, _input *InoculateDOEInput, _output *InoculateDOEOutput) {
}

func _InoculateDOEValidation(_ctx context.Context, _input *InoculateDOEInput, _output *InoculateDOEOutput) {
}
func _InoculateDOERun(_ctx context.Context, input *InoculateDOEInput) *InoculateDOEOutput {
	output := &InoculateDOEOutput{}
	_InoculateDOESetup(_ctx, input)
	_InoculateDOESteps(_ctx, input, output)
	_InoculateDOEAnalysis(_ctx, input, output)
	_InoculateDOEValidation(_ctx, input, output)
	return output
}

func InoculateDOERunSteps(_ctx context.Context, input *InoculateDOEInput) *InoculateDOESOutput {
	soutput := &InoculateDOESOutput{}
	output := _InoculateDOERun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func InoculateDOENew() interface{} {
	return &InoculateDOEElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &InoculateDOEInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _InoculateDOERun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &InoculateDOEInput{},
			Out: &InoculateDOEOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type InoculateDOEElement struct {
	inject.CheckedRunner
}

type InoculateDOEInput struct {
	Inoculum       *wtype.LHComponent
	InoculumVolume wunit.Volume
	Medium         *wtype.LHComponent
}

type InoculateDOEOutput struct {
	Seed *wtype.LHComponent
}

type InoculateDOESOutput struct {
	Data struct {
	}
	Outputs struct {
		Seed *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "InoculateDOE",
		Constructor: InoculateDOENew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/GrowthAndAssay/inoculate.an",
			Params: []component.ParamDesc{
				{Name: "Inoculum", Desc: "", Kind: "Inputs"},
				{Name: "InoculumVolume", Desc: "", Kind: "Parameters"},
				{Name: "Medium", Desc: "", Kind: "Inputs"},
				{Name: "Seed", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
