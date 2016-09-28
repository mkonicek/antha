package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

/*type Polymerase struct {
	wtype.LHComponent
	Rate_BPpers float64
	Fidelity_errorrate float64 // could dictate how many colonies are checked in validation!
	Extensiontemp Temperature
	Hotstart bool
	StockConcentration Concentration // this is normally in U?
	TargetConcentration Concentration
	// this is also a glycerol solution rather than a watersolution!
}
*/

// Input parameters for this protocol (data)

// PCRprep parameters:

/*
	// let's be ambitious and try this as part of type polymerase Polymeraseconc Volume

	//Templatetype string  // e.g. colony, genomic, pure plasmid... will effect efficiency. We could get more sophisticated here later on...
	//FullTemplatesequence string // better to use Sid's type system here after proof of concept
	//FullTemplatelength int	// clearly could be calculated from the sequence... Sid will have a method to do this already so check!
	//TargetTemplatesequence string // better to use Sid's type system here after proof of concept
	//TargetTemplatelengthinBP int
*/
// Reaction parameters: (could be a entered as thermocycle parameters type possibly?)

//Denaturationtemp Temperature

// Should be calculated from primer and template binding
// should be calculated from template length and polymerase rate

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _PCR_vol_mmx_multiplereactionsRequirements() {
}

// Conditions to run on startup
func _PCR_vol_mmx_multiplereactionsSetup(_ctx context.Context, _input *PCR_vol_mmx_multiplereactionsInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _PCR_vol_mmx_multiplereactionsSteps(_ctx context.Context, _input *PCR_vol_mmx_multiplereactionsInput, _output *PCR_vol_mmx_multiplereactionsOutput) {

	_output.Reactions = make([]*wtype.LHComponent, 0)

	// rename components

	multiplier := float64(len(_input.WellPositions))

	_input.Template.CName = _input.TemplateName
	_input.FwdPrimer.CName = _input.FwdPrimerName
	_input.RevPrimer.CName = _input.RevPrimerName

	// Make a mastermix

	mmxSample := mixer.Sample(_input.MasterMix, wunit.MultiplyVolume(_input.MasterMixVolume, multiplier))

	// pipette out to make mastermix
	mastermix := execute.MixInto(_ctx, _input.OutPlate, _input.WellPositions[0], mmxSample)

	// rest samples to zero
	samples := make([]*wtype.LHComponent, 0)

	// if this is false do stuff inside {}

	// add primers

	if !_input.PrimersalreadyAddedtoMasterMix {
		FwdPrimerSample := mixer.Sample(_input.FwdPrimer, wunit.MultiplyVolume(_input.FwdPrimerVol, multiplier))
		samples = append(samples, FwdPrimerSample)
		RevPrimerSample := mixer.Sample(_input.RevPrimer, wunit.MultiplyVolume(_input.RevPrimerVol, multiplier))
		samples = append(samples, RevPrimerSample)
	}

	// add template
	templateSample := mixer.Sample(_input.Template, wunit.MultiplyVolume(_input.Templatevolume, multiplier))
	samples = append(samples, templateSample)

	for j := range samples {
		mastermix = execute.Mix(_ctx, mastermix, samples[j])
	}
	reaction := mastermix

	// this needs to go after an initial denaturation!
	if !_input.PolymeraseAlreadyaddedtoMastermix {
		polySample := mixer.Sample(_input.PCRPolymerase, wunit.MultiplyVolume(_input.PolymeraseVolume, multiplier))

		reaction = execute.Mix(_ctx, reaction, polySample)
	}

	reaction.CName = _input.ReactionName
	// aliquot out to well positions
	for position := range _input.WellPositions {
		if position != 0 {
			totalVolume := wunit.AddVolumes([]wunit.Volume{_input.PolymeraseVolume, _input.Templatevolume, _input.RevPrimerVol, _input.FwdPrimerVol, _input.MasterMixVolume})
			_output.Reactions = append(_output.Reactions, execute.MixInto(_ctx, _input.OutPlate, _input.WellPositions[position], mixer.Sample(reaction, totalVolume)))
		}
	}

	// thermocycle parameters called from enzyme lookup:

	polymerase := _input.PCRPolymerase.CName

	extensionTemp := enzymes.DNApolymerasetemps[polymerase]["extensiontemp"]
	meltingTemp := enzymes.DNApolymerasetemps[polymerase]["meltingtemp"]

	// initial Denaturation

	r1 := execute.Incubate(_ctx, reaction, meltingTemp, _input.InitDenaturationtime, false)

	for i := 0; i < _input.Numberofcycles; i++ {

		// Denature

		r1 = execute.Incubate(_ctx, r1, meltingTemp, _input.Denaturationtime, false)

		// Anneal
		r1 = execute.Incubate(_ctx, r1, _input.AnnealingTemp, _input.Annealingtime, false)

		//extensiontime := TargetTemplatelengthinBP/PCRPolymerase.RateBPpers // we'll get type issues here so leave it out for now

		// Extend
		r1 = execute.Incubate(_ctx, r1, extensionTemp, _input.Extensiontime, false)

	}
	// Final Extension
	r1 = execute.Incubate(_ctx, r1, extensionTemp, _input.Finalextensiontime, false)

	// all done
	//Reaction = r1

	reaction.CName = _input.ReactionName
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PCR_vol_mmx_multiplereactionsAnalysis(_ctx context.Context, _input *PCR_vol_mmx_multiplereactionsInput, _output *PCR_vol_mmx_multiplereactionsOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PCR_vol_mmx_multiplereactionsValidation(_ctx context.Context, _input *PCR_vol_mmx_multiplereactionsInput, _output *PCR_vol_mmx_multiplereactionsOutput) {
}
func _PCR_vol_mmx_multiplereactionsRun(_ctx context.Context, input *PCR_vol_mmx_multiplereactionsInput) *PCR_vol_mmx_multiplereactionsOutput {
	output := &PCR_vol_mmx_multiplereactionsOutput{}
	_PCR_vol_mmx_multiplereactionsSetup(_ctx, input)
	_PCR_vol_mmx_multiplereactionsSteps(_ctx, input, output)
	_PCR_vol_mmx_multiplereactionsAnalysis(_ctx, input, output)
	_PCR_vol_mmx_multiplereactionsValidation(_ctx, input, output)
	return output
}

func PCR_vol_mmx_multiplereactionsRunSteps(_ctx context.Context, input *PCR_vol_mmx_multiplereactionsInput) *PCR_vol_mmx_multiplereactionsSOutput {
	soutput := &PCR_vol_mmx_multiplereactionsSOutput{}
	output := _PCR_vol_mmx_multiplereactionsRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PCR_vol_mmx_multiplereactionsNew() interface{} {
	return &PCR_vol_mmx_multiplereactionsElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PCR_vol_mmx_multiplereactionsInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PCR_vol_mmx_multiplereactionsRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PCR_vol_mmx_multiplereactionsInput{},
			Out: &PCR_vol_mmx_multiplereactionsOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PCR_vol_mmx_multiplereactionsElement struct {
	inject.CheckedRunner
}

type PCR_vol_mmx_multiplereactionsInput struct {
	AnnealingTemp                     wunit.Temperature
	Annealingtime                     wunit.Time
	Denaturationtime                  wunit.Time
	Extensiontime                     wunit.Time
	Finalextensiontime                wunit.Time
	FwdPrimer                         *wtype.LHComponent
	FwdPrimerName                     string
	FwdPrimerVol                      wunit.Volume
	InitDenaturationtime              wunit.Time
	MasterMix                         *wtype.LHComponent
	MasterMixVolume                   wunit.Volume
	Numberofcycles                    int
	OutPlate                          *wtype.LHPlate
	PCRPolymerase                     *wtype.LHComponent
	PolymeraseAlreadyaddedtoMastermix bool
	PolymeraseVolume                  wunit.Volume
	PrimersalreadyAddedtoMasterMix    bool
	ReactionName                      string
	RevPrimer                         *wtype.LHComponent
	RevPrimerName                     string
	RevPrimerVol                      wunit.Volume
	Template                          *wtype.LHComponent
	TemplateName                      string
	Templatevolume                    wunit.Volume
	WellPositions                     []string
}

type PCR_vol_mmx_multiplereactionsOutput struct {
	Reactions []*wtype.LHComponent
}

type PCR_vol_mmx_multiplereactionsSOutput struct {
	Data struct {
	}
	Outputs struct {
		Reactions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PCR_vol_mmx_multiplereactions",
		Constructor: PCR_vol_mmx_multiplereactionsNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/PCR/pcr_vol_mmx_multiplereactions.an",
			Params: []component.ParamDesc{
				{Name: "AnnealingTemp", Desc: "Should be calculated from primer and template binding\n", Kind: "Parameters"},
				{Name: "Annealingtime", Desc: "Denaturationtemp Temperature\n", Kind: "Parameters"},
				{Name: "Denaturationtime", Desc: "", Kind: "Parameters"},
				{Name: "Extensiontime", Desc: "should be calculated from template length and polymerase rate\n", Kind: "Parameters"},
				{Name: "Finalextensiontime", Desc: "", Kind: "Parameters"},
				{Name: "FwdPrimer", Desc: "", Kind: "Inputs"},
				{Name: "FwdPrimerName", Desc: "", Kind: "Parameters"},
				{Name: "FwdPrimerVol", Desc: "", Kind: "Parameters"},
				{Name: "InitDenaturationtime", Desc: "", Kind: "Parameters"},
				{Name: "MasterMix", Desc: "", Kind: "Inputs"},
				{Name: "MasterMixVolume", Desc: "PCRprep parameters:\n", Kind: "Parameters"},
				{Name: "Numberofcycles", Desc: "\t// let's be ambitious and try this as part of type polymerase Polymeraseconc Volume\n\n\t//Templatetype string  // e.g. colony, genomic, pure plasmid... will effect efficiency. We could get more sophisticated here later on...\n\t//FullTemplatesequence string // better to use Sid's type system here after proof of concept\n\t//FullTemplatelength int\t// clearly could be calculated from the sequence... Sid will have a method to do this already so check!\n\t//TargetTemplatesequence string // better to use Sid's type system here after proof of concept\n\t//TargetTemplatelengthinBP int\n\nReaction parameters: (could be a entered as thermocycle parameters type possibly?)\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "PCRPolymerase", Desc: "", Kind: "Inputs"},
				{Name: "PolymeraseAlreadyaddedtoMastermix", Desc: "", Kind: "Parameters"},
				{Name: "PolymeraseVolume", Desc: "", Kind: "Parameters"},
				{Name: "PrimersalreadyAddedtoMasterMix", Desc: "", Kind: "Parameters"},
				{Name: "ReactionName", Desc: "", Kind: "Parameters"},
				{Name: "RevPrimer", Desc: "", Kind: "Inputs"},
				{Name: "RevPrimerName", Desc: "", Kind: "Parameters"},
				{Name: "RevPrimerVol", Desc: "", Kind: "Parameters"},
				{Name: "Template", Desc: "", Kind: "Inputs"},
				{Name: "TemplateName", Desc: "", Kind: "Parameters"},
				{Name: "Templatevolume", Desc: "", Kind: "Parameters"},
				{Name: "WellPositions", Desc: "", Kind: "Parameters"},
				{Name: "Reactions", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
