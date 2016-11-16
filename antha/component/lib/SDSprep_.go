package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

//Input parameters for this protocol. Single instance of an SDS-PAGE sample preperation step.
//Mix 10ul of 4x stock buffer with 30ul of proteinX sample to create 40ul sample for loading.

//ProteinX
//30uL

//SDSBuffer
//10ul
//100g/L

//25g/L
//40uL

//5min
//95oC

//Biologicals

//Purified protein or cell lysate...

//Chemicals

//Consumables

//Contains protein and buffer
//Final plate with mixed components

//Biologicals

func _SDSprepSetup(_ctx context.Context, _input *SDSprepInput) {
}

func _SDSprepSteps(_ctx context.Context, _input *SDSprepInput, _output *SDSprepOutput) {

	//Method 1. Mix two things. DOES NOT WORK as recognises protein to be 1 single entity and won't handle as seperate components. ie end result is 5 things created all
	//from the same well. Check typeIIs workflow for hints.
	//
	//	Step1a
	//	LoadSample = MixInto(OutPlate,
	//	mixer.Sample(Protein, SampleVolume),
	//	mixer.Sample(Buffer, BufferVolume))
	//Try something else. Outputs are an array taking in a single (not array) of protein and buffer. Do this 12 times.

	samples := make([]*wtype.LHComponent, 0)
	bufferSample := mixer.Sample(_input.Buffer, _input.BufferVolume)
	bufferSample.CName = _input.BufferName
	samples = append(samples, bufferSample)

	proteinSample := mixer.Sample(_input.Protein, _input.SampleVolume)
	proteinSample.CName = _input.SampleName
	samples = append(samples, proteinSample)
	fmt.Println("This is a sample list ", samples)
	_output.LoadSample = execute.MixInto(_ctx, _input.OutPlate, "", samples...)

	//Methods 2.Make a sample of two things creating a list
	//	Step 1b

	//	sample	    := make([]wtype.LHComponent, 0)

	//	bufferPart  := mixer.Sample(Buffer, BufferVolume)
	//	sample	     = append([]samples, bufferSample)

	//	proteinPart := mixer.Sample(Protein, SampleVolume)
	//	sample      = append([]samples, proteinSample)

	//	LoadSample   = MixInto(OutPlate, sample...)

	//Denature the load mixture at specified temperature and time ie 95oC for 5min
	//	Step2
	_output.LoadSample = execute.Incubate(_ctx, _output.LoadSample, _input.DenatureTemp, _input.DenatureTime, false)

	//Load the water in EPAGE gel wells
	//	Step3

	//	var water water volume
	//	waterLoad := mixer.Sample(Water, WaterLoadVolume)
	//
	//Load the LoadSample into EPAGE gel
	//
	//	Loader = MixInto(EPAGE48, LoadSample)
	//
	//
	//

	//	Status = fmtSprintln(BufferVolume.ToString() "uL of", BufferName,"mixed with", SampleVolume.ToString(), "uL of", SampleName, "Total load sample available is", ReactionVolume.ToString())
}

func _SDSprepAnalysis(_ctx context.Context, _input *SDSprepInput, _output *SDSprepOutput) {
}

func _SDSprepValidation(_ctx context.Context, _input *SDSprepInput, _output *SDSprepOutput) {
}
func _SDSprepRun(_ctx context.Context, input *SDSprepInput) *SDSprepOutput {
	output := &SDSprepOutput{}
	_SDSprepSetup(_ctx, input)
	_SDSprepSteps(_ctx, input, output)
	_SDSprepAnalysis(_ctx, input, output)
	_SDSprepValidation(_ctx, input, output)
	return output
}

func SDSprepRunSteps(_ctx context.Context, input *SDSprepInput) *SDSprepSOutput {
	soutput := &SDSprepSOutput{}
	output := _SDSprepRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func SDSprepNew() interface{} {
	return &SDSprepElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &SDSprepInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _SDSprepRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &SDSprepInput{},
			Out: &SDSprepOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type SDSprepElement struct {
	inject.CheckedRunner
}

type SDSprepInput struct {
	Buffer             *wtype.LHComponent
	BufferName         string
	BufferStockConc    wunit.Concentration
	BufferVolume       wunit.Volume
	DenatureTemp       wunit.Temperature
	DenatureTime       wunit.Time
	FinalConcentration wunit.Concentration
	InPlate            *wtype.LHPlate
	OutPlate           *wtype.LHPlate
	Protein            *wtype.LHComponent
	ReactionVolume     wunit.Volume
	SampleName         string
	SampleVolume       wunit.Volume
}

type SDSprepOutput struct {
	LoadSample *wtype.LHComponent
	Status     string
}

type SDSprepSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		LoadSample *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "SDSprep",
		Constructor: SDSprepNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/SDSprep/SDSprep.an",
			Params: []component.ParamDesc{
				{Name: "Buffer", Desc: "", Kind: "Inputs"},
				{Name: "BufferName", Desc: "SDSBuffer\n", Kind: "Parameters"},
				{Name: "BufferStockConc", Desc: "100g/L\n", Kind: "Parameters"},
				{Name: "BufferVolume", Desc: "10ul\n", Kind: "Parameters"},
				{Name: "DenatureTemp", Desc: "95oC\n", Kind: "Parameters"},
				{Name: "DenatureTime", Desc: "5min\n", Kind: "Parameters"},
				{Name: "FinalConcentration", Desc: "25g/L\n", Kind: "Parameters"},
				{Name: "InPlate", Desc: "Contains protein and buffer\n", Kind: "Inputs"},
				{Name: "OutPlate", Desc: "Final plate with mixed components\n", Kind: "Inputs"},
				{Name: "Protein", Desc: "Purified protein or cell lysate...\n", Kind: "Inputs"},
				{Name: "ReactionVolume", Desc: "40uL\n", Kind: "Parameters"},
				{Name: "SampleName", Desc: "ProteinX\n", Kind: "Parameters"},
				{Name: "SampleVolume", Desc: "30uL\n", Kind: "Parameters"},
				{Name: "LoadSample", Desc: "Biologicals\n", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
