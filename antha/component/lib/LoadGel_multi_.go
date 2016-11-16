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
	"strconv"
)

//    RunVoltage      Int
//    RunLength       Time

//preload well with 10uL of water
//protein samples for running
//96 well plate with water, marker and samples
//Gel to load ie OutPlate

//Run length in cm, and protein band height and pixed density after digital scanning

func _LoadGel_multiSetup(_ctx context.Context, _input *LoadGel_multiInput) {
}

func _LoadGel_multiSteps(_ctx context.Context, _input *LoadGel_multiInput, _output *LoadGel_multiOutput) {

	// work out well coordinates for any plate
	wellpositionarray := make([]string, 0)

	//alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphabet := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J",
		"K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X",
		"Y", "Z", "AA", "BB", "CC", "DD", "EE", "FF"}
	//k := 0
	for j := 0; j < _input.GelPlate.WlsY; j++ {
		for i := 0; i < _input.GelPlate.WlsX; i++ { //countingfrom1iswhatmakesushuman := j + 1
			//k = k + 1
			wellposition := string(alphabet[j]) + strconv.Itoa(i+1)
			//fmt.Println(wellposition, k)
			wellpositionarray = append(wellpositionarray, wellposition)
		}

	}

	_output.RunSolutions = make([]*wtype.LHComponent, 0)

	var RunSolution *wtype.LHComponent

	for k, SampleName := range _input.SampleNames {
		samples := make([]*wtype.LHComponent, 0)
		waterSample := mixer.Sample(_input.Water, _input.WaterVolume)
		waterSample.CName = _input.WaterName
		samples = append(samples, waterSample)

		loadSample := mixer.Sample(_input.Protein, _input.LoadVolume)
		loadSample.CName = SampleName
		samples = append(samples, loadSample)
		fmt.Println("This is a list of samples for loading:", samples)

		RunSolution = execute.MixTo(_ctx, _input.GelPlate.Type, wellpositionarray[k], 0, samples...)
		_output.RunSolutions = append(_output.RunSolutions, RunSolution)
	}
}

func _LoadGel_multiAnalysis(_ctx context.Context, _input *LoadGel_multiInput, _output *LoadGel_multiOutput) {
}

func _LoadGel_multiValidation(_ctx context.Context, _input *LoadGel_multiInput, _output *LoadGel_multiOutput) {
}
func _LoadGel_multiRun(_ctx context.Context, input *LoadGel_multiInput) *LoadGel_multiOutput {
	output := &LoadGel_multiOutput{}
	_LoadGel_multiSetup(_ctx, input)
	_LoadGel_multiSteps(_ctx, input, output)
	_LoadGel_multiAnalysis(_ctx, input, output)
	_LoadGel_multiValidation(_ctx, input, output)
	return output
}

func LoadGel_multiRunSteps(_ctx context.Context, input *LoadGel_multiInput) *LoadGel_multiSOutput {
	soutput := &LoadGel_multiSOutput{}
	output := _LoadGel_multiRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func LoadGel_multiNew() interface{} {
	return &LoadGel_multiElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &LoadGel_multiInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _LoadGel_multiRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &LoadGel_multiInput{},
			Out: &LoadGel_multiOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type LoadGel_multiElement struct {
	inject.CheckedRunner
}

type LoadGel_multiInput struct {
	GelPlate    *wtype.LHPlate
	InPlate     *wtype.LHPlate
	LoadVolume  wunit.Volume
	Protein     *wtype.LHComponent
	SampleNames []string
	Water       *wtype.LHComponent
	WaterName   string
	WaterVolume wunit.Volume
}

type LoadGel_multiOutput struct {
	RunSolutions []*wtype.LHComponent
	Status       string
}

type LoadGel_multiSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		RunSolutions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "LoadGel_multi",
		Constructor: LoadGel_multiNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/LoadGel/LoadGel_multi.an",
			Params: []component.ParamDesc{
				{Name: "GelPlate", Desc: "Gel to load ie OutPlate\n", Kind: "Inputs"},
				{Name: "InPlate", Desc: "96 well plate with water, marker and samples\n", Kind: "Inputs"},
				{Name: "LoadVolume", Desc: "", Kind: "Parameters"},
				{Name: "Protein", Desc: "protein samples for running\n", Kind: "Inputs"},
				{Name: "SampleNames", Desc: "", Kind: "Parameters"},
				{Name: "Water", Desc: "preload well with 10uL of water\n", Kind: "Inputs"},
				{Name: "WaterName", Desc: "", Kind: "Parameters"},
				{Name: "WaterVolume", Desc: "", Kind: "Parameters"},
				{Name: "RunSolutions", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
