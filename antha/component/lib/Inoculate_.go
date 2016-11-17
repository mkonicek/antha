// Example inoculation protocol.
// Inoculates seed culture into fresh media (and logs conditions?)
// TODO: in progress from edited bradford protocol
package lib

import (
	// "liquid handler"
	//"labware"
	//"OD"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

//"github.com/antha-lang/antha/antha/anthalib/wunit"

// we do comments like this

// Input parameters for this protocol (data)

//= uL(25)
//= uL(475)
//= mgperml (100)
//= mgperml  (0.1)
//= 0 // Note: 1 replicate means experiment is in duplicate, etc.

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

//suspension which contains living cells

// type buffer which could have a concentration automatically?

// Physical outputs from this protocol with types

func _InoculateRequirements() {
	// None
}

func _InoculateSetup(_ctx context.Context, _input *InoculateInput) {
	//none
}

func _InoculateSteps(_ctx context.Context, _input *InoculateInput, _output *InoculateOutput) {
	//antibiotic_volume  := wunit.NewVolume(Media_volume.SIValue() * (Desiredantibioticconcentration.SIValue()/Antibioticstockconc.SIValue()),"l")

	media_with_antibiotic := _input.Media
	//media_with_antibiotic := mixer.Mix(mixer.Sample(Antibiotic,antibiotic_volume), mixer.Sample(Media,Media_volume))
	_output.Inoculated_culture = execute.MixInto(_ctx, _input.OutPlate, "", mixer.Sample(_input.Seed, _input.Seed_volume), mixer.Sample(media_with_antibiotic, _input.Media_volume))
}

//should the transfer to thermomixer/incubator command be included in this protocol or in a separate protocol
func _InoculateAnalysis(_ctx context.Context, _input *InoculateInput, _output *InoculateOutput) {
	//OD_at_inoculation = OD.Inoculated_culture // need to know signatures of protocol_OD I,O,Q - function signature
}

func _InoculateValidation(_ctx context.Context, _input *InoculateInput, _output *InoculateOutput) {
	/*
		if OD.sample_absorbance > 1 {
		panic("Sample likely needs further dilution")
		}
		if OD.sample_absorbance < 0.02 {
		warn("low inoculation OD")
		//could add visual (i.e. manual or camera based) validation
		// TODO: add test of replicate variance
		}*/
}
func _InoculateRun(_ctx context.Context, input *InoculateInput) *InoculateOutput {
	output := &InoculateOutput{}
	_InoculateSetup(_ctx, input)
	_InoculateSteps(_ctx, input, output)
	_InoculateAnalysis(_ctx, input, output)
	_InoculateValidation(_ctx, input, output)
	return output
}

func InoculateRunSteps(_ctx context.Context, input *InoculateInput) *InoculateSOutput {
	soutput := &InoculateSOutput{}
	output := _InoculateRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func InoculateNew() interface{} {
	return &InoculateElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &InoculateInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _InoculateRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &InoculateInput{},
			Out: &InoculateOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type InoculateElement struct {
	inject.CheckedRunner
}

type InoculateInput struct {
	Antibiotic                     *wtype.LHComponent
	Antibioticstockconc            wunit.Concentration
	Desiredantibioticconcentration wunit.Concentration
	Media                          *wtype.LHComponent
	Media_volume                   wunit.Volume
	OutPlate                       *wtype.LHPlate
	Replicate_count                int
	Seed                           *wtype.LHComponent
	Seed_volume                    wunit.Volume
}

type InoculateOutput struct {
	Inoculated_culture *wtype.LHComponent
	OD_at_inoculation  float64
}

type InoculateSOutput struct {
	Data struct {
		OD_at_inoculation float64
	}
	Outputs struct {
		Inoculated_culture *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Inoculate",
		Constructor: InoculateNew,
		Desc: component.ComponentDesc{
			Desc: "Example inoculation protocol.\nInoculates seed culture into fresh media (and logs conditions?)\nTODO: in progress from edited bradford protocol\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/Fermentation/Inoculate.an",
			Params: []component.ParamDesc{
				{Name: "Antibiotic", Desc: "type buffer which could have a concentration automatically?\n", Kind: "Inputs"},
				{Name: "Antibioticstockconc", Desc: "= mgperml (100)\n", Kind: "Parameters"},
				{Name: "Desiredantibioticconcentration", Desc: "= mgperml  (0.1)\n", Kind: "Parameters"},
				{Name: "Media", Desc: "", Kind: "Inputs"},
				{Name: "Media_volume", Desc: "= uL(475)\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "Replicate_count", Desc: "= 0 // Note: 1 replicate means experiment is in duplicate, etc.\n", Kind: "Parameters"},
				{Name: "Seed", Desc: "suspension which contains living cells\n", Kind: "Inputs"},
				{Name: "Seed_volume", Desc: "= uL(25)\n", Kind: "Parameters"},
				{Name: "Inoculated_culture", Desc: "", Kind: "Outputs"},
				{Name: "OD_at_inoculation", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
