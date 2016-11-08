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

func _ResuspendDNARequirements() {
}

func _ResuspendDNASetup(_ctx context.Context, _input *ResuspendDNAInput) {
}

func _ResuspendDNASteps(_ctx context.Context, _input *ResuspendDNAInput, _output *ResuspendDNAOutput) {

	targetconcgperL := _input.TargetConc.GramPerL(_input.MolecularWeight).SIValue()

	dnamassG := _input.DNAMass.SIValue()

	if _input.DNAMass.Unit().BaseSIUnit() == "kg" {
		dnamassG = dnamassG * 1000
		fmt.Println("Base Unit correction; Base unit of mass = ", _input.DNAMass.Unit().BaseSIUnit(), " therfore multiplying by 1000 to convert to grams")
	}

	volumetoadd := wunit.NewVolume(dnamassG/targetconcgperL, "L")

	diluentSample := mixer.Sample(_input.Diluent, volumetoadd)

	diluentSample.Type = wtype.LTDNAMIX

	_output.ResuspendedDNA = execute.MixInto(_ctx, _input.DNAPlate, _input.Well, diluentSample)

}

func _ResuspendDNAAnalysis(_ctx context.Context, _input *ResuspendDNAInput, _output *ResuspendDNAOutput) {
}

func _ResuspendDNAValidation(_ctx context.Context, _input *ResuspendDNAInput, _output *ResuspendDNAOutput) {
}
func _ResuspendDNARun(_ctx context.Context, input *ResuspendDNAInput) *ResuspendDNAOutput {
	output := &ResuspendDNAOutput{}
	_ResuspendDNASetup(_ctx, input)
	_ResuspendDNASteps(_ctx, input, output)
	_ResuspendDNAAnalysis(_ctx, input, output)
	_ResuspendDNAValidation(_ctx, input, output)
	return output
}

func ResuspendDNARunSteps(_ctx context.Context, input *ResuspendDNAInput) *ResuspendDNASOutput {
	soutput := &ResuspendDNASOutput{}
	output := _ResuspendDNARun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ResuspendDNANew() interface{} {
	return &ResuspendDNAElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ResuspendDNAInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ResuspendDNARun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ResuspendDNAInput{},
			Out: &ResuspendDNAOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type ResuspendDNAElement struct {
	inject.CheckedRunner
}

type ResuspendDNAInput struct {
	DNAMass         wunit.Mass
	DNAPlate        *wtype.LHPlate
	Diluent         *wtype.LHComponent
	MolecularWeight float64
	TargetConc      wunit.Concentration
	Well            string
}

type ResuspendDNAOutput struct {
	ResuspendedDNA *wtype.LHComponent
}

type ResuspendDNASOutput struct {
	Data struct {
	}
	Outputs struct {
		ResuspendedDNA *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "ResuspendDNA",
		Constructor: ResuspendDNANew,
		Desc: component.ComponentDesc{
			Desc: "Protocol for resuspending freeze dried DNA with a diluent\n",
			Path: "antha/component/an/Liquid_handling/ResuspendDNA/ResuspendDNA.an",
			Params: []component.ParamDesc{
				{Name: "DNAMass", Desc: "", Kind: "Parameters"},
				{Name: "DNAPlate", Desc: "", Kind: "Inputs"},
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "MolecularWeight", Desc: "", Kind: "Parameters"},
				{Name: "TargetConc", Desc: "", Kind: "Parameters"},
				{Name: "Well", Desc: "", Kind: "Parameters"},
				{Name: "ResuspendedDNA", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
