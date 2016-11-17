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

// e.g. DMSO

// Physical outputs from this protocol with types

func _Colony_PCRRequirements() {
}

// Conditions to run on startup
func _Colony_PCRSetup(_ctx context.Context, _input *Colony_PCRInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _Colony_PCRSteps(_ctx context.Context, _input *Colony_PCRInput, _output *Colony_PCROutput) {

	// Mix components
	samples := make([]*wtype.LHComponent, 0)
	bufferSample := mixer.SampleForTotalVolume(_input.Buffer, _input.ReactionVolume)
	samples = append(samples, bufferSample)
	templateSample := mixer.Sample(_input.Template, _input.Templatevolume)
	samples = append(samples, templateSample)
	dntpSample := mixer.SampleForConcentration(_input.DNTPS, _input.DNTPconc)
	samples = append(samples, dntpSample)
	FwdPrimerSample := mixer.SampleForConcentration(_input.FwdPrimer, _input.FwdPrimerConc)
	samples = append(samples, FwdPrimerSample)
	RevPrimerSample := mixer.SampleForConcentration(_input.RevPrimer, _input.RevPrimerConc)
	samples = append(samples, RevPrimerSample)

	for _, additive := range _input.Additives {
		additiveSample := mixer.SampleForConcentration(additive, _input.Additiveconc)
		samples = append(samples, additiveSample)
	}

	polySample := mixer.SampleForConcentration(_input.PCRPolymerase, _input.TargetpolymeraseConcentration)
	samples = append(samples, polySample)
	reaction := execute.MixInto(_ctx, _input.OutPlate, "", samples...)

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
	_output.Reaction = r1
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Colony_PCRAnalysis(_ctx context.Context, _input *Colony_PCRInput, _output *Colony_PCROutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _Colony_PCRValidation(_ctx context.Context, _input *Colony_PCRInput, _output *Colony_PCROutput) {
}
func _Colony_PCRRun(_ctx context.Context, input *Colony_PCRInput) *Colony_PCROutput {
	output := &Colony_PCROutput{}
	_Colony_PCRSetup(_ctx, input)
	_Colony_PCRSteps(_ctx, input, output)
	_Colony_PCRAnalysis(_ctx, input, output)
	_Colony_PCRValidation(_ctx, input, output)
	return output
}

func Colony_PCRRunSteps(_ctx context.Context, input *Colony_PCRInput) *Colony_PCRSOutput {
	soutput := &Colony_PCRSOutput{}
	output := _Colony_PCRRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Colony_PCRNew() interface{} {
	return &Colony_PCRElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Colony_PCRInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Colony_PCRRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Colony_PCRInput{},
			Out: &Colony_PCROutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Colony_PCRElement struct {
	inject.CheckedRunner
}

type Colony_PCRInput struct {
	Additiveconc                  wunit.Concentration
	Additives                     []*wtype.LHComponent
	AnnealingTemp                 wunit.Temperature
	Annealingtime                 wunit.Time
	Buffer                        *wtype.LHComponent
	DNTPS                         *wtype.LHComponent
	DNTPconc                      wunit.Concentration
	Denaturationtime              wunit.Time
	Extensiontemp                 wunit.Temperature
	Extensiontime                 wunit.Time
	Finalextensiontime            wunit.Time
	FwdPrimer                     *wtype.LHComponent
	FwdPrimerConc                 wunit.Concentration
	InitDenaturationtime          wunit.Time
	Numberofcycles                int
	OutPlate                      *wtype.LHPlate
	PCRPolymerase                 *wtype.LHComponent
	ReactionVolume                wunit.Volume
	RevPrimer                     *wtype.LHComponent
	RevPrimerConc                 wunit.Concentration
	TargetpolymeraseConcentration wunit.Concentration
	Template                      *wtype.LHComponent
	Templatevolume                wunit.Volume
}

type Colony_PCROutput struct {
	Reaction *wtype.LHComponent
}

type Colony_PCRSOutput struct {
	Data struct {
	}
	Outputs struct {
		Reaction *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Colony_PCR",
		Constructor: Colony_PCRNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/PCR/Colony_PCR.an",
			Params: []component.ParamDesc{
				{Name: "Additiveconc", Desc: "", Kind: "Parameters"},
				{Name: "Additives", Desc: "e.g. DMSO\n", Kind: "Inputs"},
				{Name: "AnnealingTemp", Desc: "Should be calculated from primer and template binding\n", Kind: "Parameters"},
				{Name: "Annealingtime", Desc: "Denaturationtemp Temperature\n", Kind: "Parameters"},
				{Name: "Buffer", Desc: "", Kind: "Inputs"},
				{Name: "DNTPS", Desc: "", Kind: "Inputs"},
				{Name: "DNTPconc", Desc: "", Kind: "Parameters"},
				{Name: "Denaturationtime", Desc: "", Kind: "Parameters"},
				{Name: "Extensiontemp", Desc: "", Kind: "Parameters"},
				{Name: "Extensiontime", Desc: "should be calculated from template length and polymerase rate\n", Kind: "Parameters"},
				{Name: "Finalextensiontime", Desc: "", Kind: "Parameters"},
				{Name: "FwdPrimer", Desc: "", Kind: "Inputs"},
				{Name: "FwdPrimerConc", Desc: "", Kind: "Parameters"},
				{Name: "InitDenaturationtime", Desc: "", Kind: "Parameters"},
				{Name: "Numberofcycles", Desc: "\t// let's be ambitious and try this as part of type polymerase Polymeraseconc Volume\n\n\t//Templatetype string  // e.g. colony, genomic, pure plasmid... will effect efficiency. We could get more sophisticated here later on...\n\t//FullTemplatesequence string // better to use Sid's type system here after proof of concept\n\t//FullTemplatelength int\t// clearly could be calculated from the sequence... Sid will have a method to do this already so check!\n\t//TargetTemplatesequence string // better to use Sid's type system here after proof of concept\n\t//TargetTemplatelengthinBP int\n\nReaction parameters: (could be a entered as thermocycle parameters type possibly?)\n", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "PCRPolymerase", Desc: "", Kind: "Inputs"},
				{Name: "ReactionVolume", Desc: "PCRprep parameters:\n", Kind: "Parameters"},
				{Name: "RevPrimer", Desc: "", Kind: "Inputs"},
				{Name: "RevPrimerConc", Desc: "", Kind: "Parameters"},
				{Name: "TargetpolymeraseConcentration", Desc: "", Kind: "Parameters"},
				{Name: "Template", Desc: "", Kind: "Inputs"},
				{Name: "Templatevolume", Desc: "", Kind: "Parameters"},
				{Name: "Reaction", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
