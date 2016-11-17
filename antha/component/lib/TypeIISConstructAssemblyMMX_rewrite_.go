package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

func _TypeIISConstructAssemblyMMX_rewriteRequirements() {}

// Conditions to run on startup
func _TypeIISConstructAssemblyMMX_rewriteSetup(_ctx context.Context, _input *TypeIISConstructAssemblyMMX_rewriteInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _TypeIISConstructAssemblyMMX_rewriteSteps(_ctx context.Context, _input *TypeIISConstructAssemblyMMX_rewriteInput, _output *TypeIISConstructAssemblyMMX_rewriteOutput) {
	var err error

	_output.ConstructName = _input.OutputConstructName

	last := len(_input.PartSeqs) - 1
	output, count, _, seq, err := enzymes.Assemblysimulator(enzymes.Assemblyparameters{
		Constructname: _output.ConstructName,
		Enzymename:    _input.EnzymeName,
		Vector:        _input.PartSeqs[last],
		Partsinorder:  _input.PartSeqs[:last],
	})
	if err != nil {
		_output.Errors = append(_output.Errors, fmt.Sprintf("%s: %s", output, err))
		return
	}
	if count != 1 {
		_output.Errors = append(_output.Errors, fmt.Sprintf("no successful assembly"))
		return
	}

	_output.Sequence = seq

	waterSample := mixer.SampleForTotalVolume(_input.Water, _input.ReactionVolume)
	reacmixture := execute.MixTo(_ctx, _input.OutPlate.Type, _input.OutputLocation, _input.OutputPlateNum, waterSample)

	for k, part := range _input.Parts {
		fmt.Println("creating dna part num ", k, " comp ", part.CName, " renamed to ", _input.PartSeqs[k].Nm, " vol ", _input.PartVols[k])

		part.Type, err = wtype.LiquidTypeFromString(_input.LHPolicyName)

		if err != nil {
			_output.Errors = append(_output.Errors, fmt.Sprintf("cannot find liquid type: %s", err))
			return
		}

		partSample := mixer.Sample(part, _input.PartVols[k])
		partSample.CName = _input.PartSeqs[k].Nm
		reacmixture = execute.Mix(_ctx, reacmixture, partSample)
	}

	mmxSample := mixer.Sample(_input.MasterMix, _input.MasterMixVolume)
	// ensure the last step is mixed
	mmxSample.Type = wtype.LTDNAMIX
	_output.Reaction = execute.Mix(_ctx, reacmixture, mmxSample)

	// incubate the reaction mixture
	// commented out pending changes to incubate
	execute.Incubate(_ctx, _output.Reaction, _input.ReactionTemp, _input.ReactionTime, false)
	// inactivate
	//Incubate(Reaction, InactivationTemp, InactivationTime, false)
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _TypeIISConstructAssemblyMMX_rewriteAnalysis(_ctx context.Context, _input *TypeIISConstructAssemblyMMX_rewriteInput, _output *TypeIISConstructAssemblyMMX_rewriteOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _TypeIISConstructAssemblyMMX_rewriteValidation(_ctx context.Context, _input *TypeIISConstructAssemblyMMX_rewriteInput, _output *TypeIISConstructAssemblyMMX_rewriteOutput) {
}
func _TypeIISConstructAssemblyMMX_rewriteRun(_ctx context.Context, input *TypeIISConstructAssemblyMMX_rewriteInput) *TypeIISConstructAssemblyMMX_rewriteOutput {
	output := &TypeIISConstructAssemblyMMX_rewriteOutput{}
	_TypeIISConstructAssemblyMMX_rewriteSetup(_ctx, input)
	_TypeIISConstructAssemblyMMX_rewriteSteps(_ctx, input, output)
	_TypeIISConstructAssemblyMMX_rewriteAnalysis(_ctx, input, output)
	_TypeIISConstructAssemblyMMX_rewriteValidation(_ctx, input, output)
	return output
}

func TypeIISConstructAssemblyMMX_rewriteRunSteps(_ctx context.Context, input *TypeIISConstructAssemblyMMX_rewriteInput) *TypeIISConstructAssemblyMMX_rewriteSOutput {
	soutput := &TypeIISConstructAssemblyMMX_rewriteSOutput{}
	output := _TypeIISConstructAssemblyMMX_rewriteRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func TypeIISConstructAssemblyMMX_rewriteNew() interface{} {
	return &TypeIISConstructAssemblyMMX_rewriteElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &TypeIISConstructAssemblyMMX_rewriteInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _TypeIISConstructAssemblyMMX_rewriteRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &TypeIISConstructAssemblyMMX_rewriteInput{},
			Out: &TypeIISConstructAssemblyMMX_rewriteOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type TypeIISConstructAssemblyMMX_rewriteElement struct {
	inject.CheckedRunner
}

type TypeIISConstructAssemblyMMX_rewriteInput struct {
	EnzymeName          string
	InactivationTemp    wunit.Temperature
	InactivationTime    wunit.Time
	LHPolicyName        string
	MasterMix           *wtype.LHComponent
	MasterMixVolume     wunit.Volume
	OutPlate            *wtype.LHPlate
	OutputConstructName string
	OutputLocation      string
	OutputPlateNum      int
	OutputReactionName  string
	PartSeqs            []wtype.DNASequence
	PartVols            []wunit.Volume
	Parts               []*wtype.LHComponent
	ReactionTemp        wunit.Temperature
	ReactionTime        wunit.Time
	ReactionVolume      wunit.Volume
	Water               *wtype.LHComponent
}

type TypeIISConstructAssemblyMMX_rewriteOutput struct {
	ConstructName string
	Errors        []string
	Reaction      *wtype.LHComponent
	Sequence      wtype.DNASequence
}

type TypeIISConstructAssemblyMMX_rewriteSOutput struct {
	Data struct {
		ConstructName string
		Errors        []string
		Sequence      wtype.DNASequence
	}
	Outputs struct {
		Reaction *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "TypeIISConstructAssemblyMMX_rewrite",
		Constructor: TypeIISConstructAssemblyMMX_rewriteNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/TypeIIsAssembly/TypeIISConstructAssemblyMMX_rewrite/TypeIISConstructAssemblyMMX.an",
			Params: []component.ParamDesc{
				{Name: "EnzymeName", Desc: "", Kind: "Parameters"},
				{Name: "InactivationTemp", Desc: "", Kind: "Parameters"},
				{Name: "InactivationTime", Desc: "", Kind: "Parameters"},
				{Name: "LHPolicyName", Desc: "", Kind: "Parameters"},
				{Name: "MasterMix", Desc: "", Kind: "Inputs"},
				{Name: "MasterMixVolume", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "OutputConstructName", Desc: "", Kind: "Parameters"},
				{Name: "OutputLocation", Desc: "", Kind: "Parameters"},
				{Name: "OutputPlateNum", Desc: "", Kind: "Parameters"},
				{Name: "OutputReactionName", Desc: "", Kind: "Parameters"},
				{Name: "PartSeqs", Desc: "", Kind: "Parameters"},
				{Name: "PartVols", Desc: "", Kind: "Parameters"},
				{Name: "Parts", Desc: "", Kind: "Inputs"},
				{Name: "ReactionTemp", Desc: "", Kind: "Parameters"},
				{Name: "ReactionTime", Desc: "", Kind: "Parameters"},
				{Name: "ReactionVolume", Desc: "", Kind: "Parameters"},
				{Name: "Water", Desc: "", Kind: "Inputs"},
				{Name: "ConstructName", Desc: "", Kind: "Data"},
				{Name: "Errors", Desc: "", Kind: "Data"},
				{Name: "Reaction", Desc: "", Kind: "Outputs"},
				{Name: "Sequence", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
