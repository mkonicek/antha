// row
package lib

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

//"strconv"

// Input parameters for this protocol (data)

//MediaVolume  Volume

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _PickColoniesRequirements() {

}

// Conditions to run on startup
func _PickColoniesSetup(_ctx context.Context, _input *PickColoniesInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _PickColoniesSteps(_ctx context.Context, _input *PickColoniesInput, _output *PickColoniesOutput) {

	components := make([]*wtype.LHComponent, 0)
	_output.ColoniesinRecoveryMedia = make([]*wtype.LHComponent, 0)
	var colonyinrecoverymedia *wtype.LHComponent

	for _, colonylocation := range _input.ColonyLocations {

		//togetplateid := mixer.Sample(Colonytype,ColonyVolume)

		//id := togetplateid.Loc

		colony := wtype.NewLHComponent()

		colony.Loc = colonylocation
		//colony.LContainer.Plateid= id

		realcolonysample := mixer.Sample(_input.Colonytype, _input.ColonyVolume)

		components = append(components, realcolonysample)

		colonyinrecoverymedia = execute.MixInto(_ctx, _input.OutPlatewithMedia, "", components...)

		_output.ColoniesinRecoveryMedia = append(_output.ColoniesinRecoveryMedia, colonyinrecoverymedia)
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PickColoniesAnalysis(_ctx context.Context, _input *PickColoniesInput, _output *PickColoniesOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PickColoniesValidation(_ctx context.Context, _input *PickColoniesInput, _output *PickColoniesOutput) {

}
func _PickColoniesRun(_ctx context.Context, input *PickColoniesInput) *PickColoniesOutput {
	output := &PickColoniesOutput{}
	_PickColoniesSetup(_ctx, input)
	_PickColoniesSteps(_ctx, input, output)
	_PickColoniesAnalysis(_ctx, input, output)
	_PickColoniesValidation(_ctx, input, output)
	return output
}

func PickColoniesRunSteps(_ctx context.Context, input *PickColoniesInput) *PickColoniesSOutput {
	soutput := &PickColoniesSOutput{}
	output := _PickColoniesRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PickColoniesNew() interface{} {
	return &PickColoniesElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PickColoniesInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PickColoniesRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PickColoniesInput{},
			Out: &PickColoniesOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PickColoniesElement struct {
	inject.CheckedRunner
}

type PickColoniesInput struct {
	ColonyLocations   []string
	ColonyVolume      wunit.Volume
	Colonytype        *wtype.LHComponent
	InPlate           *wtype.LHPlate
	OutPlatewithMedia *wtype.LHPlate
}

type PickColoniesOutput struct {
	ColoniesinRecoveryMedia []*wtype.LHComponent
}

type PickColoniesSOutput struct {
	Data struct {
	}
	Outputs struct {
		ColoniesinRecoveryMedia []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PickColonies",
		Constructor: PickColoniesNew,
		Desc: component.ComponentDesc{
			Desc: "row\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/PickColonies/PickColonies.an",
			Params: []component.ParamDesc{
				{Name: "ColonyLocations", Desc: "MediaVolume  Volume\n", Kind: "Parameters"},
				{Name: "ColonyVolume", Desc: "", Kind: "Parameters"},
				{Name: "Colonytype", Desc: "", Kind: "Inputs"},
				{Name: "InPlate", Desc: "", Kind: "Inputs"},
				{Name: "OutPlatewithMedia", Desc: "", Kind: "Inputs"},
				{Name: "ColoniesinRecoveryMedia", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
