// Protocol for resuspending freeze dried DNA with a diluent
package lib

import

// we need to import the wtype package to use the LHComponent type
// the mixer package is required to use the Sample function
(
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

func _PairOligosRequirements() {
}

func _PairOligosSetup(_ctx context.Context, _input *PairOligosInput) {
}

func _PairOligosSteps(_ctx context.Context, _input *PairOligosInput, _output *PairOligosOutput) {

	diluentSample := mixer.SampleForTotalVolume(_input.Diluent, _input.TotalVolume)

	fwdOligoSample := mixer.Sample(_input.FwdOligo, _input.FWDOligoVolume)

	revOligoSample := mixer.Sample(_input.RevOligo, _input.REVOligoVolume)

	revOligoSample.Type = wtype.LTDNAMIX

	_output.OligoPairs = execute.MixNamed(_ctx, _input.Plate.Type, _input.Well, fmt.Sprint("OligoPlate", _input.PlateNumber), diluentSample, fwdOligoSample, revOligoSample)

	_output.OligoPairs = execute.Incubate(_ctx, _output.OligoPairs, _input.IncubationTemp, _input.IncubationTime, false)

}

func _PairOligosAnalysis(_ctx context.Context, _input *PairOligosInput, _output *PairOligosOutput) {
}

func _PairOligosValidation(_ctx context.Context, _input *PairOligosInput, _output *PairOligosOutput) {
}
func _PairOligosRun(_ctx context.Context, input *PairOligosInput) *PairOligosOutput {
	output := &PairOligosOutput{}
	_PairOligosSetup(_ctx, input)
	_PairOligosSteps(_ctx, input, output)
	_PairOligosAnalysis(_ctx, input, output)
	_PairOligosValidation(_ctx, input, output)
	return output
}

func PairOligosRunSteps(_ctx context.Context, input *PairOligosInput) *PairOligosSOutput {
	soutput := &PairOligosSOutput{}
	output := _PairOligosRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PairOligosNew() interface{} {
	return &PairOligosElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PairOligosInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PairOligosRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PairOligosInput{},
			Out: &PairOligosOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PairOligosElement struct {
	inject.CheckedRunner
}

type PairOligosInput struct {
	Diluent        *wtype.LHComponent
	FWDOligoVolume wunit.Volume
	FwdOligo       *wtype.LHComponent
	IncubationTemp wunit.Temperature
	IncubationTime wunit.Time
	Plate          *wtype.LHPlate
	PlateNumber    int
	REVOligoVolume wunit.Volume
	RevOligo       *wtype.LHComponent
	TotalVolume    wunit.Volume
	Well           string
}

type PairOligosOutput struct {
	OligoPairs *wtype.LHComponent
}

type PairOligosSOutput struct {
	Data struct {
	}
	Outputs struct {
		OligoPairs *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PairOligos",
		Constructor: PairOligosNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol for resuspending freeze dried DNA with a diluent\n",
			Path: "antha/component/an/Liquid_handling/ResuspendDNA/PairOligos.an",
			Params: []component.ParamDesc{
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "FWDOligoVolume", Desc: "", Kind: "Parameters"},
				{Name: "FwdOligo", Desc: "", Kind: "Inputs"},
				{Name: "IncubationTemp", Desc: "", Kind: "Parameters"},
				{Name: "IncubationTime", Desc: "", Kind: "Parameters"},
				{Name: "Plate", Desc: "", Kind: "Inputs"},
				{Name: "PlateNumber", Desc: "", Kind: "Parameters"},
				{Name: "REVOligoVolume", Desc: "", Kind: "Parameters"},
				{Name: "RevOligo", Desc: "", Kind: "Inputs"},
				{Name: "TotalVolume", Desc: "", Kind: "Parameters"},
				{Name: "Well", Desc: "", Kind: "Parameters"},
				{Name: "OligoPairs", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
