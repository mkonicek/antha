package lib

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/factory"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

// PCRprep parameters

// e.g. ["left homology arm"]:"templatename"
// e.g. ["left homology arm"]:"fwdprimer","revprimer"

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _AutoAssemblyRequirements() {
}

// Conditions to run on startup
func _AutoAssemblySetup(_ctx context.Context, _input *AutoAssemblyInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _AutoAssemblySteps(_ctx context.Context, _input *AutoAssemblyInput, _output *AutoAssemblyOutput) {

	// set up a counter to use as an index for increasing well position
	var counter int

	// set up some empty slices to fill as we iterate through the reactions

	_output.Reactions = make([]*wtype.LHComponent, 0)
	volumes := make([]wunit.Volume, 0)
	welllocations := make([]string, 0)

	// range through the Reaction to template map
	for reactionname, constructname := range _input.Reactiontonames {

		// use counter to find next available well position in plate
		wellposition := _input.OutPlate.AllWellPositions(wtype.BYCOLUMN)[counter]

		// Run TypeIISConstructAssemblyMMX_forscreen element
		result := TypeIISConstructAssemblyMMX_forscreenRunSteps(_ctx, &TypeIISConstructAssemblyMMX_forscreenInput{InactivationTemp: wunit.NewTemperature(40, "C"),
			InactivationTime:    wunit.NewTime(60, "s"),
			MasterMixVolume:     wunit.NewVolume(5, "ul"),
			PartVols:            _input.Reactiontoassemblyvolumes,
			Partseqs:            _input.Reactiontopartseqs,
			OutputReactionName:  reactionname,
			OutputConstructName: constructname,
			ReactionTemp:        wunit.NewTemperature(25, "C"),
			ReactionTime:        wunit.NewTime(1800, "s"),
			ReactionVolume:      wunit.NewVolume(20, "ul"),
			OutputLocation:      wellposition,
			EnzymeName:          _input.EnzymeName,

			MasterMix: _input.MasterMixtype,
			Water:     _input.Watertype,
			Parts:     _input.Reactiontoparttypes},
		)

		// add result to reactions slice
		_output.Reactions = append(_output.Reactions, result.Outputs.Reaction)
		volumes = append(volumes, result.Outputs.Reaction.Volume())
		OutputLocation = append(OutputLocation, wellposition)
		// increase counter by 1 ready for next iteration of loop
		counter++

	}

	// once all values of loop have been completed, export the plate contents as a csv file
	_output.Error = wtype.ExportPlateCSV(_input.Projectname+".csv", Plate, _input.Projectname+"outputPlate", welllocations, _output.Reactions, volumes)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AutoAssemblyAnalysis(_ctx context.Context, _input *AutoAssemblyInput, _output *AutoAssemblyOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AutoAssemblyValidation(_ctx context.Context, _input *AutoAssemblyInput, _output *AutoAssemblyOutput) {
}
func _AutoAssemblyRun(_ctx context.Context, input *AutoAssemblyInput) *AutoAssemblyOutput {
	output := &AutoAssemblyOutput{}
	_AutoAssemblySetup(_ctx, input)
	_AutoAssemblySteps(_ctx, input, output)
	_AutoAssemblyAnalysis(_ctx, input, output)
	_AutoAssemblyValidation(_ctx, input, output)
	return output
}

func AutoAssemblyRunSteps(_ctx context.Context, input *AutoAssemblyInput) *AutoAssemblySOutput {
	soutput := &AutoAssemblySOutput{}
	output := _AutoAssemblyRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AutoAssemblyNew() interface{} {
	return &AutoAssemblyElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AutoAssemblyInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AutoAssemblyRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AutoAssemblyInput{},
			Out: &AutoAssemblyOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type AutoAssemblyElement struct {
	inject.CheckedRunner
}

type AutoAssemblyInput struct {
	EnzymeName                string
	LHPolicyName              string
	MasterMixtype             *wtype.LHComponent
	OutPlate                  *wtype.LHPlate
	OutputPlateNum            int
	Projectname               string
	Reactiontoassemblyparts   map[string][]string
	Reactiontoassemblyvolumes map[string][]wunit.Volume
	Reactiontonames           map[string][]string
	Reactiontopartseqs        map[string][]string
	Reactiontoparttypes       map[string][]*wtype.LHComponent
	Watertype                 *wtype.LHComponent
}

type AutoAssemblyOutput struct {
	Error     error
	Reactions []*wtype.LHComponent
}

type AutoAssemblySOutput struct {
	Data struct {
		Error error
	}
	Outputs struct {
		Reactions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AutoAssembly",
		Constructor: AutoAssemblyNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/PooledLibrary/playground/Refactored/AutoAssembly.an",
			Params: []component.ParamDesc{
				{Name: "EnzymeName", Desc: "", Kind: "Parameters"},
				{Name: "LHPolicyName", Desc: "", Kind: "Parameters"},
				{Name: "MasterMixtype", Desc: "", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "OutputPlateNum", Desc: "", Kind: "Parameters"},
				{Name: "Projectname", Desc: "PCRprep parameters\n", Kind: "Parameters"},
				{Name: "Reactiontoassemblyparts", Desc: "e.g. [\"left homology arm\"]:\"fwdprimer\",\"revprimer\"\n", Kind: "Parameters"},
				{Name: "Reactiontoassemblyvolumes", Desc: "e.g. [\"left homology arm\"]:\"templatename\"\n", Kind: "Parameters"},
				{Name: "Reactiontonames", Desc: "", Kind: "Parameters"},
				{Name: "Reactiontopartseqs", Desc: "", Kind: "Parameters"},
				{Name: "Reactiontoparttypes", Desc: "", Kind: "Inputs"},
				{Name: "Watertype", Desc: "", Kind: "Inputs"},
				{Name: "Error", Desc: "", Kind: "Data"},
				{Name: "Reactions", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
