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

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// e.g. DMSO

// Physical outputs from this protocol with types

func _MakePCRmmxRequirements() {
}

// Conditions to run on startup
func _MakePCRmmxSetup(_ctx context.Context, _input *MakePCRmmxInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _MakePCRmmxSteps(_ctx context.Context, _input *MakePCRmmxInput, _output *MakePCRmmxOutput) {

	// rename components

	bufferVolume := (wunit.CopyVolume(_input.ReactionVolume))
	bufferVolume.DivideBy(float64(_input.BufferConcinX))

	// Make a mastermix
	samples := make([]*wtype.LHComponent, 0)
	waterSample := mixer.Sample(_input.Water, _input.WaterVolume)
	bufferSample := mixer.Sample(_input.Buffer, bufferVolume)
	samples = append(samples, waterSample, bufferSample)

	dntpSample := mixer.Sample(_input.DNTPS, _input.DNTPVol)
	samples = append(samples, dntpSample)

	if len(_input.Additives) != len(_input.AdditiveVols) {
		execute.Errorf(_ctx, "Bad things are going to happen if you have different numbers of additives and additivevolumes")
	}

	for i := range _input.Additives {
		additiveSample := mixer.Sample(_input.Additives[i], _input.AdditiveVols[i])
		samples = append(samples, additiveSample)
	}

	if _input.Hotstart == false {
		polySample := mixer.Sample(_input.PCRPolymerase, _input.PolymeraseVolume)
		samples = append(samples, polySample)
	}

	// if this is true do stuff inside {}
	if _input.AddPrimerstoMasterMix {

		FwdPrimerSample := mixer.Sample(_input.FwdPrimer, _input.FwdPrimerVol)
		samples = append(samples, FwdPrimerSample)
		RevPrimerSample := mixer.Sample(_input.RevPrimer, _input.RevPrimerVol)
		samples = append(samples, RevPrimerSample)

	}

	// pipette out to make mastermix
	mastermix := execute.MixInto(_ctx, _input.OutPlate, "", samples...)

	// rest samples to zero
	samples = make([]*wtype.LHComponent, 0)

	// if this is false do stuff inside {}
	if !_input.AddPrimerstoMasterMix {

		FwdPrimerSample := mixer.Sample(_input.FwdPrimer, _input.FwdPrimerVol)
		samples = append(samples, FwdPrimerSample)
		RevPrimerSample := mixer.Sample(_input.RevPrimer, _input.RevPrimerVol)
		samples = append(samples, RevPrimerSample)

	}

	for j := range samples {
		mastermix = execute.Mix(_ctx, mastermix, samples[j])
	}
	reaction := mastermix

	// this needs to go after an initial denaturation!
	if _input.Hotstart {
		polySample := mixer.Sample(_input.PCRPolymerase, _input.PolymeraseVolume)

		reaction = execute.Mix(_ctx, reaction, polySample)
	}

	// all done
	_output.Reaction = reaction

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MakePCRmmxAnalysis(_ctx context.Context, _input *MakePCRmmxInput, _output *MakePCRmmxOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _MakePCRmmxValidation(_ctx context.Context, _input *MakePCRmmxInput, _output *MakePCRmmxOutput) {
}
func _MakePCRmmxRun(_ctx context.Context, input *MakePCRmmxInput) *MakePCRmmxOutput {
	output := &MakePCRmmxOutput{}
	_MakePCRmmxSetup(_ctx, input)
	_MakePCRmmxSteps(_ctx, input, output)
	_MakePCRmmxAnalysis(_ctx, input, output)
	_MakePCRmmxValidation(_ctx, input, output)
	return output
}

func MakePCRmmxRunSteps(_ctx context.Context, input *MakePCRmmxInput) *MakePCRmmxSOutput {
	soutput := &MakePCRmmxSOutput{}
	output := _MakePCRmmxRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MakePCRmmxNew() interface{} {
	return &MakePCRmmxElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MakePCRmmxInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MakePCRmmxRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MakePCRmmxInput{},
			Out: &MakePCRmmxOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type MakePCRmmxElement struct {
	inject.CheckedRunner
}

type MakePCRmmxInput struct {
	AddPrimerstoMasterMix bool
	AdditiveVols          []wunit.Volume
	Additives             []*wtype.LHComponent
	Buffer                *wtype.LHComponent
	BufferConcinX         int
	DNTPS                 *wtype.LHComponent
	DNTPVol               wunit.Volume
	FwdPrimer             *wtype.LHComponent
	FwdPrimerName         string
	FwdPrimerVol          wunit.Volume
	Hotstart              bool
	OutPlate              *wtype.LHPlate
	PCRPolymerase         *wtype.LHComponent
	PolymeraseVolume      wunit.Volume
	ReactionVolume        wunit.Volume
	RevPrimer             *wtype.LHComponent
	RevPrimerName         string
	RevPrimerVol          wunit.Volume
	Templatevolume        wunit.Volume
	Water                 *wtype.LHComponent
	WaterVolume           wunit.Volume
}

type MakePCRmmxOutput struct {
	Reaction *wtype.LHComponent
}

type MakePCRmmxSOutput struct {
	Data struct {
	}
	Outputs struct {
		Reaction *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MakePCRmmx",
		Constructor: MakePCRmmxNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/PCR/MakePCRmmx.an",
			Params: []component.ParamDesc{
				{Name: "AddPrimerstoMasterMix", Desc: "", Kind: "Parameters"},
				{Name: "AdditiveVols", Desc: "", Kind: "Parameters"},
				{Name: "Additives", Desc: "e.g. DMSO\n", Kind: "Inputs"},
				{Name: "Buffer", Desc: "", Kind: "Inputs"},
				{Name: "BufferConcinX", Desc: "", Kind: "Parameters"},
				{Name: "DNTPS", Desc: "", Kind: "Inputs"},
				{Name: "DNTPVol", Desc: "", Kind: "Parameters"},
				{Name: "FwdPrimer", Desc: "", Kind: "Inputs"},
				{Name: "FwdPrimerName", Desc: "", Kind: "Parameters"},
				{Name: "FwdPrimerVol", Desc: "", Kind: "Parameters"},
				{Name: "Hotstart", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "PCRPolymerase", Desc: "", Kind: "Inputs"},
				{Name: "PolymeraseVolume", Desc: "", Kind: "Parameters"},
				{Name: "ReactionVolume", Desc: "", Kind: "Parameters"},
				{Name: "RevPrimer", Desc: "", Kind: "Inputs"},
				{Name: "RevPrimerName", Desc: "", Kind: "Parameters"},
				{Name: "RevPrimerVol", Desc: "", Kind: "Parameters"},
				{Name: "Templatevolume", Desc: "\t// let's be ambitious and try this as part of type polymerase Polymeraseconc Volume\n\n\t//Templatetype string  // e.g. colony, genomic, pure plasmid... will effect efficiency. We could get more sophisticated here later on...\n\t//FullTemplatesequence string // better to use Sid's type system here after proof of concept\n\t//FullTemplatelength int\t// clearly could be calculated from the sequence... Sid will have a method to do this already so check!\n\t//TargetTemplatesequence string // better to use Sid's type system here after proof of concept\n\t//TargetTemplatelengthinBP int\n\nReaction parameters: (could be a entered as thermocycle parameters type possibly?)\n", Kind: "Parameters"},
				{Name: "Water", Desc: "", Kind: "Inputs"},
				{Name: "WaterVolume", Desc: "PCRprep parameters:\n", Kind: "Parameters"},
				{Name: "Reaction", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
