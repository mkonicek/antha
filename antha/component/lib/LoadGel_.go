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

//    RunVoltage      Int
//    RunLength       Time

//preload well with 10uL of water
//protein samples for running
//96 well plate with water, marker and samples
//Gel to load ie OutPlate

//Run length in cm, and protein band height and pixed density after digital scanning

func _LoadGelSetup(_ctx context.Context, _input *LoadGelInput) {
}

func _LoadGelSteps(_ctx context.Context, _input *LoadGelInput, _output *LoadGelOutput) {

	samples := make([]*wtype.LHComponent, 0)
	waterSample := mixer.Sample(_input.Water, _input.WaterVolume)
	waterSample.CName = _input.WaterName
	samples = append(samples, waterSample)

	loadSample := mixer.Sample(_input.Protein, _input.LoadVolume)
	loadSample.CName = _input.SampleName
	samples = append(samples, loadSample)
	fmt.Println("This is a list of samples for loading:", samples)

	_output.RunSolution = execute.MixInto(_ctx, _input.GelPlate, "", samples...)
}

func _LoadGelAnalysis(_ctx context.Context, _input *LoadGelInput, _output *LoadGelOutput) {
}

func _LoadGelValidation(_ctx context.Context, _input *LoadGelInput, _output *LoadGelOutput) {
}
func _LoadGelRun(_ctx context.Context, input *LoadGelInput) *LoadGelOutput {
	output := &LoadGelOutput{}
	_LoadGelSetup(_ctx, input)
	_LoadGelSteps(_ctx, input, output)
	_LoadGelAnalysis(_ctx, input, output)
	_LoadGelValidation(_ctx, input, output)
	return output
}

func LoadGelRunSteps(_ctx context.Context, input *LoadGelInput) *LoadGelSOutput {
	soutput := &LoadGelSOutput{}
	output := _LoadGelRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func LoadGelNew() interface{} {
	return &LoadGelElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &LoadGelInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _LoadGelRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &LoadGelInput{},
			Out: &LoadGelOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type LoadGelElement struct {
	inject.CheckedRunner
}

type LoadGelInput struct {
	GelPlate    *wtype.LHPlate
	InPlate     *wtype.LHPlate
	LoadVolume  wunit.Volume
	Protein     *wtype.LHComponent
	SampleName  string
	Water       *wtype.LHComponent
	WaterName   string
	WaterVolume wunit.Volume
}

type LoadGelOutput struct {
	RunSolution *wtype.LHComponent
	Status      string
}

type LoadGelSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		RunSolution *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "LoadGel",
		Constructor: LoadGelNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/LoadGel/LoadGel.an",
			Params: []component.ParamDesc{
				{Name: "GelPlate", Desc: "Gel to load ie OutPlate\n", Kind: "Inputs"},
				{Name: "InPlate", Desc: "96 well plate with water, marker and samples\n", Kind: "Inputs"},
				{Name: "LoadVolume", Desc: "", Kind: "Parameters"},
				{Name: "Protein", Desc: "protein samples for running\n", Kind: "Inputs"},
				{Name: "SampleName", Desc: "", Kind: "Parameters"},
				{Name: "Water", Desc: "preload well with 10uL of water\n", Kind: "Inputs"},
				{Name: "WaterName", Desc: "", Kind: "Parameters"},
				{Name: "WaterVolume", Desc: "", Kind: "Parameters"},
				{Name: "RunSolution", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
