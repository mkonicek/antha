package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/platereader"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

func _AbsorbanceMeasurementRequirements() {
}

func _AbsorbanceMeasurementSetup(_ctx context.Context, _input *AbsorbanceMeasurementInput) {
}

func _AbsorbanceMeasurementSteps(_ctx context.Context, _input *AbsorbanceMeasurementInput, _output *AbsorbanceMeasurementOutput) {

	// dilute sample
	diluentSample := mixer.Sample(_input.Diluent, _input.DilutionVolume)

	sampleforreading := mixer.SampleAll(_input.SampleForReading)

	dilutedSample := execute.MixTo(_ctx, _input.Plate.Type, "", 1, sampleforreading, diluentSample)
	//dilutedSample:=Mix(sampleforreading, diluentSample)

	// read
	abs := platereader.ReadAbsorbance(_input.Plate, dilutedSample, _input.AbsorbanceWavelength.RawValue())

	// prepare blank and read
	blankSample := execute.MixTo(_ctx, _input.Plate.Type, "", 1, mixer.Sample(_input.Diluent, dilutedSample.Volume()))

	blankabs := platereader.ReadAbsorbance(_input.Plate, blankSample, _input.AbsorbanceWavelength.RawValue())

	// blank correct
	blankcorrected, err := platereader.Blankcorrect(blankabs, abs)
	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	// estimate pathlength
	pathlength, err := platereader.EstimatePathLength(_input.Plate, dilutedSample.Volume())

	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	// pathlength correct
	pathlengthcorrected := platereader.PathlengthCorrect(pathlength, blankcorrected)

	// calculate actual conc based on extinction coefficient
	actualconc := platereader.Concentration(pathlengthcorrected, _input.ExtinctionCoefficient)

	_output.ActualConcentration = actualconc

	_output.AbsorbanceMeasurement = abs.Reading

}

func _AbsorbanceMeasurementAnalysis(_ctx context.Context, _input *AbsorbanceMeasurementInput, _output *AbsorbanceMeasurementOutput) {
}

func _AbsorbanceMeasurementValidation(_ctx context.Context, _input *AbsorbanceMeasurementInput, _output *AbsorbanceMeasurementOutput) {
}
func _AbsorbanceMeasurementRun(_ctx context.Context, input *AbsorbanceMeasurementInput) *AbsorbanceMeasurementOutput {
	output := &AbsorbanceMeasurementOutput{}
	_AbsorbanceMeasurementSetup(_ctx, input)
	_AbsorbanceMeasurementSteps(_ctx, input, output)
	_AbsorbanceMeasurementAnalysis(_ctx, input, output)
	_AbsorbanceMeasurementValidation(_ctx, input, output)
	return output
}

func AbsorbanceMeasurementRunSteps(_ctx context.Context, input *AbsorbanceMeasurementInput) *AbsorbanceMeasurementSOutput {
	soutput := &AbsorbanceMeasurementSOutput{}
	output := _AbsorbanceMeasurementRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AbsorbanceMeasurementNew() interface{} {
	return &AbsorbanceMeasurementElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AbsorbanceMeasurementInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AbsorbanceMeasurementRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AbsorbanceMeasurementInput{},
			Out: &AbsorbanceMeasurementOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type AbsorbanceMeasurementElement struct {
	inject.CheckedRunner
}

type AbsorbanceMeasurementInput struct {
	AbsorbanceWavelength  wunit.Length
	Diluent               *wtype.LHComponent
	DilutionVolume        wunit.Volume
	ExtinctionCoefficient float64
	Plate                 *wtype.LHPlate
	SampleForReading      *wtype.LHComponent
}

type AbsorbanceMeasurementOutput struct {
	AbsorbanceMeasurement float64
	ActualConcentration   wunit.Concentration
}

type AbsorbanceMeasurementSOutput struct {
	Data struct {
		AbsorbanceMeasurement float64
		ActualConcentration   wunit.Concentration
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AbsorbanceMeasurement",
		Constructor: AbsorbanceMeasurementNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/GrowthAndAssay/absorbanceassay.an",
			Params: []component.ParamDesc{
				{Name: "AbsorbanceWavelength", Desc: "", Kind: "Parameters"},
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "DilutionVolume", Desc: "", Kind: "Parameters"},
				{Name: "ExtinctionCoefficient", Desc: "", Kind: "Parameters"},
				{Name: "Plate", Desc: "", Kind: "Inputs"},
				{Name: "SampleForReading", Desc: "", Kind: "Inputs"},
				{Name: "AbsorbanceMeasurement", Desc: "", Kind: "Data"},
				{Name: "ActualConcentration", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
