package lib

import (
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

func _GrowthAndAssayRequirements() {
}

func _GrowthAndAssaySetup(_ctx context.Context, _input *GrowthAndAssayInput) {
}

func _GrowthAndAssaySteps(_ctx context.Context, _input *GrowthAndAssayInput, _output *GrowthAndAssayOutput) {
}
func _GrowthAndAssayAnalysis(_ctx context.Context, _input *GrowthAndAssayInput, _output *GrowthAndAssayOutput) {
}

func _GrowthAndAssayValidation(_ctx context.Context, _input *GrowthAndAssayInput, _output *GrowthAndAssayOutput) {
}
func _GrowthAndAssayRun(_ctx context.Context, input *GrowthAndAssayInput) *GrowthAndAssayOutput {
	output := &GrowthAndAssayOutput{}
	_GrowthAndAssaySetup(_ctx, input)
	_GrowthAndAssaySteps(_ctx, input, output)
	_GrowthAndAssayAnalysis(_ctx, input, output)
	_GrowthAndAssayValidation(_ctx, input, output)
	return output
}

func GrowthAndAssayRunSteps(_ctx context.Context, input *GrowthAndAssayInput) *GrowthAndAssaySOutput {
	soutput := &GrowthAndAssaySOutput{}
	output := _GrowthAndAssayRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func GrowthAndAssayNew() interface{} {
	return &GrowthAndAssayElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &GrowthAndAssayInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _GrowthAndAssayRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &GrowthAndAssayInput{},
			Out: &GrowthAndAssayOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type GrowthAndAssayElement struct {
	inject.CheckedRunner
}

type GrowthAndAssayInput struct {
}

type GrowthAndAssayOutput struct {
}

type GrowthAndAssaySOutput struct {
	Data struct {
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "GrowthAndAssay",
		Constructor: GrowthAndAssayNew,
		Desc: component.ComponentDesc{
			Desc:   "",
			Path:   "src/github.com/antha-lang/antha/antha/component/an/GrowthAndAssay/default.an",
			Params: []component.ParamDesc{},
		},
	}); err != nil {
		panic(err)
	}
}
