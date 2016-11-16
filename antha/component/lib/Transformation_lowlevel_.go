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

// Input parameters for this protocol (data)

//CompetentCellvolumeperassembly wunit.Volume //= 50.(uL)

//Coolplatepositions []string
//HotplatePositions []string
//RecoveryPositions []string

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

//RecoveryPlate *wtype.LHPlate
//CompcellPlate *wtype.LHPlate

// Physical outputs from this protocol with types

func _Transformation_lowlevelRequirements() {
}

// Conditions to run on startup
func _Transformation_lowlevelSetup(_ctx context.Context, _input *Transformation_lowlevelInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _Transformation_lowlevelSteps(_ctx context.Context, _input *Transformation_lowlevelInput, _output *Transformation_lowlevelOutput) {

	// declare variables for use later
	var transformations []*wtype.LHComponent
	var incubatedtransformations []*wtype.LHComponent
	var recoverymixes []*wtype.LHComponent

	// add dna to competent cell aliquots
	for i, reaction := range _input.Reactions {
		DNAsample := mixer.Sample(reaction, _input.Reactionvolume)

		transformationmix := execute.Mix(_ctx, _input.ReadyCompCells[i], DNAsample)

		transformations = append(transformations, transformationmix)

	}

	// wait
	for _, transformationmix := range transformations {
		incubated := execute.Incubate(_ctx, transformationmix, _input.Postplasmidtemp, _input.Postplasmidtime, false)
		incubatedtransformations = append(incubatedtransformations, incubated)
	}

	// add to recovery media
	for j, transformation := range incubatedtransformations {
		recovery := execute.Mix(_ctx, _input.RecoveryMediaAliquots[j], transformation)
		recoverymixes = append(recoverymixes, recovery)
	}

	// recovery
	for _, mix := range recoverymixes {
		incubated := execute.Incubate(_ctx, mix, _input.Recoverytemp, _input.Recoverytime, true)
		_output.Transformedcells = append(_output.Transformedcells, incubated)
	}

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Transformation_lowlevelAnalysis(_ctx context.Context, _input *Transformation_lowlevelInput, _output *Transformation_lowlevelOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _Transformation_lowlevelValidation(_ctx context.Context, _input *Transformation_lowlevelInput, _output *Transformation_lowlevelOutput) {
}
func _Transformation_lowlevelRun(_ctx context.Context, input *Transformation_lowlevelInput) *Transformation_lowlevelOutput {
	output := &Transformation_lowlevelOutput{}
	_Transformation_lowlevelSetup(_ctx, input)
	_Transformation_lowlevelSteps(_ctx, input, output)
	_Transformation_lowlevelAnalysis(_ctx, input, output)
	_Transformation_lowlevelValidation(_ctx, input, output)
	return output
}

func Transformation_lowlevelRunSteps(_ctx context.Context, input *Transformation_lowlevelInput) *Transformation_lowlevelSOutput {
	soutput := &Transformation_lowlevelSOutput{}
	output := _Transformation_lowlevelRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Transformation_lowlevelNew() interface{} {
	return &Transformation_lowlevelElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Transformation_lowlevelInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Transformation_lowlevelRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Transformation_lowlevelInput{},
			Out: &Transformation_lowlevelOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Transformation_lowlevelElement struct {
	inject.CheckedRunner
}

type Transformation_lowlevelInput struct {
	Postplasmidtemp       wunit.Temperature
	Postplasmidtime       wunit.Time
	Reactions             []*wtype.LHComponent
	Reactionvolume        wunit.Volume
	ReadyCompCells        []*wtype.LHComponent
	RecoveryMediaAliquots []*wtype.LHComponent
	Recoverytemp          wunit.Temperature
	Recoverytime          wunit.Time
}

type Transformation_lowlevelOutput struct {
	Transformedcells []*wtype.LHComponent
}

type Transformation_lowlevelSOutput struct {
	Data struct {
	}
	Outputs struct {
		Transformedcells []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Transformation_lowlevel",
		Constructor: Transformation_lowlevelNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/Transformation/Transformation_lowlevel.an",
			Params: []component.ParamDesc{
				{Name: "Postplasmidtemp", Desc: "", Kind: "Parameters"},
				{Name: "Postplasmidtime", Desc: "", Kind: "Parameters"},
				{Name: "Reactions", Desc: "", Kind: "Inputs"},
				{Name: "Reactionvolume", Desc: "CompetentCellvolumeperassembly wunit.Volume //= 50.(uL)\n", Kind: "Parameters"},
				{Name: "ReadyCompCells", Desc: "", Kind: "Inputs"},
				{Name: "RecoveryMediaAliquots", Desc: "", Kind: "Inputs"},
				{Name: "Recoverytemp", Desc: "", Kind: "Parameters"},
				{Name: "Recoverytime", Desc: "", Kind: "Parameters"},
				{Name: "Transformedcells", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
