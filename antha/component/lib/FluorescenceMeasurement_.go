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

func _FluorescenceMeasurementRequirements() {
}

func _FluorescenceMeasurementSetup(_ctx context.Context, _input *FluorescenceMeasurementInput) {
}

func _FluorescenceMeasurementSteps(_ctx context.Context, _input *FluorescenceMeasurementInput, _output *FluorescenceMeasurementOutput) {
	dilutionSample := mixer.Sample(_input.Diluent, _input.DilutionVolume)
	execute.Mix(_ctx, _input.SampleForReading, dilutionSample)
	//dilutedSample:=Mix(SampleForReading, dilutionSample)
	//	FluorescenceMeasurement = ReadEM(dilutedSample, ExcitationWavelength, EmissionWavelength)
	_output.FluorescenceMeasurement = 0.5
}

func _FluorescenceMeasurementAnalysis(_ctx context.Context, _input *FluorescenceMeasurementInput, _output *FluorescenceMeasurementOutput) {
}

func _FluorescenceMeasurementValidation(_ctx context.Context, _input *FluorescenceMeasurementInput, _output *FluorescenceMeasurementOutput) {
}
func _FluorescenceMeasurementRun(_ctx context.Context, input *FluorescenceMeasurementInput) *FluorescenceMeasurementOutput {
	output := &FluorescenceMeasurementOutput{}
	_FluorescenceMeasurementSetup(_ctx, input)
	_FluorescenceMeasurementSteps(_ctx, input, output)
	_FluorescenceMeasurementAnalysis(_ctx, input, output)
	_FluorescenceMeasurementValidation(_ctx, input, output)
	return output
}

func FluorescenceMeasurementRunSteps(_ctx context.Context, input *FluorescenceMeasurementInput) *FluorescenceMeasurementSOutput {
	soutput := &FluorescenceMeasurementSOutput{}
	output := _FluorescenceMeasurementRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func FluorescenceMeasurementNew() interface{} {
	return &FluorescenceMeasurementElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &FluorescenceMeasurementInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _FluorescenceMeasurementRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &FluorescenceMeasurementInput{},
			Out: &FluorescenceMeasurementOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type FluorescenceMeasurementElement struct {
	inject.CheckedRunner
}

type FluorescenceMeasurementInput struct {
	Diluent              *wtype.LHComponent
	DilutionVolume       wunit.Volume
	EmissionWavelength   wunit.Length
	ExcitationWavelength wunit.Length
	SampleForReading     *wtype.LHComponent
}

type FluorescenceMeasurementOutput struct {
	FluorescenceMeasurement float64
}

type FluorescenceMeasurementSOutput struct {
	Data struct {
		FluorescenceMeasurement float64
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "FluorescenceMeasurement",
		Constructor: FluorescenceMeasurementNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/GrowthAndAssay/fluorescenceassay.an",
			Params: []component.ParamDesc{
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "DilutionVolume", Desc: "", Kind: "Parameters"},
				{Name: "EmissionWavelength", Desc: "", Kind: "Parameters"},
				{Name: "ExcitationWavelength", Desc: "", Kind: "Parameters"},
				{Name: "SampleForReading", Desc: "", Kind: "Inputs"},
				{Name: "FluorescenceMeasurement", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
