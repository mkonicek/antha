package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/export"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
	"path/filepath"
)

//"github.com/antha-lang/antha/microArch/factory"

// Input parameters for this protocol (data)

// PCRprep parameters

// e.g. ["left homology arm"]:"templatename"

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _AutoAssembly2Requirements() {
}

// Conditions to run on startup
func _AutoAssembly2Setup(_ctx context.Context, _input *AutoAssembly2Input) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _AutoAssembly2Steps(_ctx context.Context, _input *AutoAssembly2Input, _output *AutoAssembly2Output) {

	// set up a counter to use as an index for increasing well position
	var counter int

	// set up some empty slices to fill as we iterate through the reactions

	_output.Reactions = make([]*wtype.LHComponent, 0)
	_output.ReactionsMap = make(map[string]*wtype.LHComponent)
	_output.SequencesMap = make(map[string]wtype.DNASequence)
	volumes := make([]wunit.Volume, 0)
	OutputLocations := make([]string, 0)

	// range through the Reaction to template map
	for reactionname, constructname := range _input.Reactiontonames {

		// use counter to find next available well position in plate
		wellposition := _input.OutPlate.AllWellPositions(wtype.BYCOLUMN)[counter]

		// Run TypeIISConstructAssemblyMMX_forscreen element
		result := TypeIISConstructAssemblyMMX_forscreenRunSteps(_ctx, &TypeIISConstructAssemblyMMX_forscreenInput{InactivationTemp: wunit.NewTemperature(40, "C"),
			InactivationTime:    wunit.NewTime(60, "s"),
			MasterMixVolume:     wunit.NewVolume(5, "ul"),
			PartVols:            _input.Reactiontoassemblyvolumes[reactionname],
			PartSeqs:            _input.Reactiontopartseqs[reactionname],
			OutputReactionName:  reactionname,
			OutputConstructName: constructname,
			ReactionTemp:        wunit.NewTemperature(25, "C"),
			ReactionTime:        wunit.NewTime(1800, "s"),
			ReactionVolume:      wunit.NewVolume(20, "ul"),
			OutputLocation:      wellposition,
			EnzymeName:          _input.EnzymeName,
			LHPolicyName:        _input.LHPolicyName,

			MasterMix: _input.MasterMixtype,
			Water:     _input.Watertype,
			Parts:     _input.Reactiontoparttypes[reactionname],
			OutPlate:  _input.OutPlate},
		)

		// add result to reactions slice
		_output.Reactions = append(_output.Reactions, result.Outputs.Reaction)
		_output.ReactionsMap[reactionname] = result.Outputs.Reaction
		_output.SequencesMap[reactionname] = result.Data.Sequence
		_output.Sequences = append(_output.Sequences, result.Data.Sequence)
		volumes = append(volumes, result.Outputs.Reaction.Volume())
		OutputLocations = append(OutputLocations, wellposition)
		// increase counter by 1 ready for next iteration of loop
		counter++

	}

	// export simulated assembled sequences to file
	export.Makefastaserial2(export.LOCAL, filepath.Join(_input.Projectname, "AssembledSequences"), _output.Sequences)

	// once all values of loop have been completed, export the plate contents as a csv file
	_output.Error = wtype.ExportPlateCSV(_input.Projectname+".csv", _input.OutPlate, _input.Projectname+"outputPlate", OutputLocations, _output.Reactions, volumes)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AutoAssembly2Analysis(_ctx context.Context, _input *AutoAssembly2Input, _output *AutoAssembly2Output) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AutoAssembly2Validation(_ctx context.Context, _input *AutoAssembly2Input, _output *AutoAssembly2Output) {
}
func _AutoAssembly2Run(_ctx context.Context, input *AutoAssembly2Input) *AutoAssembly2Output {
	output := &AutoAssembly2Output{}
	_AutoAssembly2Setup(_ctx, input)
	_AutoAssembly2Steps(_ctx, input, output)
	_AutoAssembly2Analysis(_ctx, input, output)
	_AutoAssembly2Validation(_ctx, input, output)
	return output
}

func AutoAssembly2RunSteps(_ctx context.Context, input *AutoAssembly2Input) *AutoAssembly2SOutput {
	soutput := &AutoAssembly2SOutput{}
	output := _AutoAssembly2Run(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AutoAssembly2New() interface{} {
	return &AutoAssembly2Element{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AutoAssembly2Input{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AutoAssembly2Run(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AutoAssembly2Input{},
			Out: &AutoAssembly2Output{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type AutoAssembly2Element struct {
	inject.CheckedRunner
}

type AutoAssembly2Input struct {
	EnzymeName                string
	LHPolicyName              string
	MasterMixtype             *wtype.LHComponent
	OutPlate                  *wtype.LHPlate
	OutputPlateNum            int
	Projectname               string
	Reactiontoassemblyvolumes map[string][]wunit.Volume
	Reactiontonames           map[string]string
	Reactiontopartseqs        map[string][]wtype.DNASequence
	Reactiontoparttypes       map[string][]*wtype.LHComponent
	Watertype                 *wtype.LHComponent
}

type AutoAssembly2Output struct {
	Error        error
	Reactions    []*wtype.LHComponent
	ReactionsMap map[string]*wtype.LHComponent
	Sequences    []wtype.DNASequence
	SequencesMap map[string]wtype.DNASequence
}

type AutoAssembly2SOutput struct {
	Data struct {
		Error        error
		Sequences    []wtype.DNASequence
		SequencesMap map[string]wtype.DNASequence
	}
	Outputs struct {
		Reactions    []*wtype.LHComponent
		ReactionsMap map[string]*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AutoAssembly2",
		Constructor: AutoAssembly2New,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/PooledLibrary/playground/Refactored2/AutoAssembly.an",
			Params: []component.ParamDesc{
				{Name: "EnzymeName", Desc: "", Kind: "Parameters"},
				{Name: "LHPolicyName", Desc: "", Kind: "Parameters"},
				{Name: "MasterMixtype", Desc: "", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "OutputPlateNum", Desc: "", Kind: "Parameters"},
				{Name: "Projectname", Desc: "PCRprep parameters\n", Kind: "Parameters"},
				{Name: "Reactiontoassemblyvolumes", Desc: "e.g. [\"left homology arm\"]:\"templatename\"\n", Kind: "Parameters"},
				{Name: "Reactiontonames", Desc: "", Kind: "Parameters"},
				{Name: "Reactiontopartseqs", Desc: "", Kind: "Parameters"},
				{Name: "Reactiontoparttypes", Desc: "", Kind: "Inputs"},
				{Name: "Watertype", Desc: "", Kind: "Inputs"},
				{Name: "Error", Desc: "", Kind: "Data"},
				{Name: "Reactions", Desc: "", Kind: "Outputs"},
				{Name: "ReactionsMap", Desc: "", Kind: "Outputs"},
				{Name: "Sequences", Desc: "", Kind: "Data"},
				{Name: "SequencesMap", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
