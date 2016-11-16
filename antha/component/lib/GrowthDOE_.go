package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

//AngularVelocity

func _GrowthDOERequirements() {
}

func _GrowthDOESetup(_ctx context.Context, _input *GrowthDOEInput) {
}

func _GrowthDOESteps(_ctx context.Context, _input *GrowthDOEInput, _output *GrowthDOEOutput) {
	// TODO add RPM here
	fmt.Println("StrainInMedium =", _input.StrainInMedium, "growthtime = ", _input.Growthtime.ToString())
	_output.Culture = execute.Incubate(_ctx, _input.StrainInMedium, _input.Growthtemp, _input.Growthtime, true)
	fmt.Println("Culture =", _output.Culture, "StrainInMedium =", _input.StrainInMedium)
}

func _GrowthDOEAnalysis(_ctx context.Context, _input *GrowthDOEInput, _output *GrowthDOEOutput) {
}

func _GrowthDOEValidation(_ctx context.Context, _input *GrowthDOEInput, _output *GrowthDOEOutput) {
}
func _GrowthDOERun(_ctx context.Context, input *GrowthDOEInput) *GrowthDOEOutput {
	output := &GrowthDOEOutput{}
	_GrowthDOESetup(_ctx, input)
	_GrowthDOESteps(_ctx, input, output)
	_GrowthDOEAnalysis(_ctx, input, output)
	_GrowthDOEValidation(_ctx, input, output)
	return output
}

func GrowthDOERunSteps(_ctx context.Context, input *GrowthDOEInput) *GrowthDOESOutput {
	soutput := &GrowthDOESOutput{}
	output := _GrowthDOERun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func GrowthDOENew() interface{} {
	return &GrowthDOEElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &GrowthDOEInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _GrowthDOERun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &GrowthDOEInput{},
			Out: &GrowthDOEOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type GrowthDOEElement struct {
	inject.CheckedRunner
}

type GrowthDOEInput struct {
	Growthtemp     wunit.Temperature
	Growthtime     wunit.Time
	ShakingSpeed   wunit.Rate
	StrainInMedium *wtype.LHComponent
}

type GrowthDOEOutput struct {
	Culture *wtype.LHComponent
}

type GrowthDOESOutput struct {
	Data struct {
	}
	Outputs struct {
		Culture *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "GrowthDOE",
		Constructor: GrowthDOENew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/GrowthAndAssay/growth.an",
			Params: []component.ParamDesc{
				{Name: "Growthtemp", Desc: "", Kind: "Parameters"},
				{Name: "Growthtime", Desc: "", Kind: "Parameters"},
				{Name: "ShakingSpeed", Desc: "AngularVelocity\n", Kind: "Parameters"},
				{Name: "StrainInMedium", Desc: "", Kind: "Inputs"},
				{Name: "Culture", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
