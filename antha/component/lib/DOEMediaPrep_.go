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

func _DOEMediaPrepRequirements() {
}

func _DOEMediaPrepSetup(_ctx context.Context, _input *DOEMediaPrepInput) {
}

func _DOEMediaPrepSteps(_ctx context.Context, _input *DOEMediaPrepInput, _output *DOEMediaPrepOutput) {
	medsample := mixer.SampleForTotalVolume(_input.BaseMedium, _input.TotalVolume)
	medinplate := execute.MixTo(_ctx, _input.OutputPlateType, "", 1, medsample)
	yesample := mixer.SampleForConcentration(_input.YeastExtract, _input.YeastExtractConc)
	ye_med := execute.Mix(_ctx, medinplate, yesample)
	trysample := mixer.SampleForConcentration(_input.Tryptone, _input.TryptoneConc)
	try_ye_med := execute.Mix(_ctx, ye_med, trysample)
	glysample := mixer.SampleForConcentration(_input.Glycerol, _input.GlycerolConc)
	_output.GrowthMedium = execute.Mix(_ctx, try_ye_med, glysample)
}

func _DOEMediaPrepAnalysis(_ctx context.Context, _input *DOEMediaPrepInput, _output *DOEMediaPrepOutput) {
}

func _DOEMediaPrepValidation(_ctx context.Context, _input *DOEMediaPrepInput, _output *DOEMediaPrepOutput) {
}
func _DOEMediaPrepRun(_ctx context.Context, input *DOEMediaPrepInput) *DOEMediaPrepOutput {
	output := &DOEMediaPrepOutput{}
	_DOEMediaPrepSetup(_ctx, input)
	_DOEMediaPrepSteps(_ctx, input, output)
	_DOEMediaPrepAnalysis(_ctx, input, output)
	_DOEMediaPrepValidation(_ctx, input, output)
	return output
}

func DOEMediaPrepRunSteps(_ctx context.Context, input *DOEMediaPrepInput) *DOEMediaPrepSOutput {
	soutput := &DOEMediaPrepSOutput{}
	output := _DOEMediaPrepRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func DOEMediaPrepNew() interface{} {
	return &DOEMediaPrepElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &DOEMediaPrepInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _DOEMediaPrepRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &DOEMediaPrepInput{},
			Out: &DOEMediaPrepOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type DOEMediaPrepElement struct {
	inject.CheckedRunner
}

type DOEMediaPrepInput struct {
	BaseMedium       *wtype.LHComponent
	Glycerol         *wtype.LHComponent
	GlycerolConc     wunit.Concentration
	OutputPlateType  string
	TotalVolume      wunit.Volume
	Tryptone         *wtype.LHComponent
	TryptoneConc     wunit.Concentration
	YeastExtract     *wtype.LHComponent
	YeastExtractConc wunit.Concentration
}

type DOEMediaPrepOutput struct {
	GrowthMedium *wtype.LHComponent
}

type DOEMediaPrepSOutput struct {
	Data struct {
	}
	Outputs struct {
		GrowthMedium *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "DOEMediaPrep",
		Constructor: DOEMediaPrepNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/GrowthAndAssay/mediaprep.an",
			Params: []component.ParamDesc{
				{Name: "BaseMedium", Desc: "", Kind: "Inputs"},
				{Name: "Glycerol", Desc: "", Kind: "Inputs"},
				{Name: "GlycerolConc", Desc: "", Kind: "Parameters"},
				{Name: "OutputPlateType", Desc: "", Kind: "Parameters"},
				{Name: "TotalVolume", Desc: "", Kind: "Parameters"},
				{Name: "Tryptone", Desc: "", Kind: "Inputs"},
				{Name: "TryptoneConc", Desc: "", Kind: "Parameters"},
				{Name: "YeastExtract", Desc: "", Kind: "Inputs"},
				{Name: "YeastExtractConc", Desc: "", Kind: "Parameters"},
				{Name: "GrowthMedium", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
