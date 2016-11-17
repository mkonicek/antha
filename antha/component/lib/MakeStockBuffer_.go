package lib

import (
	//"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	//"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/buffers"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

//Bufferstock		*wtype.LHComponent

//OutPlate 		*wtype.LHPlate

// Physical outputs from this protocol with types

//Buffer 			*wtype.LHComponent

// Data which is returned from this protocol, and data types

//Status string

//OriginalDiluentVolume Volume

// Input Requirement specification
func _MakeStockBufferRequirements() {

}

// Conditions to run on startup
func _MakeStockBufferSetup(_ctx context.Context, _input *MakeStockBufferInput) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _MakeStockBufferSteps(_ctx context.Context, _input *MakeStockBufferInput, _output *MakeStockBufferOutput) {
	//Bufferstockvolume := wunit.NewVolume((FinalVolume.SIValue() * FinalConcentration.SIValue()/Bufferstockconc.SIValue()),"l")
	var err error
	_output.StockConc, err = buffers.StockConcentration(_input.Moleculename, _input.MassAddedinG, _input.Diluent.CName, _input.TotalVolume)
	if err != nil {
		panic(err)
	}

	/*
		Buffer = MixInto(OutPlate,"",
		mixer.Sample(Bufferstock,BufferVolumeAdded),
		mixer.Sample(Diluent,DiluentVolume))

		Status = fmt.Sprintln( "Buffer stock volume = ", BufferVolumeAdded.ToString(), "of", Bufferstock.CName,
		"was added to ", DiluentVolume.ToString(), "of", Diluent.CName,
		"to make ", BufferVolumeAdded.SIValue() + DiluentVolume.SIValue(), "L", "of", Buffername,
		"Buffer stock conc =",FinalConcentration.ToString())

		OriginalDiluentVolume = DiluentVolume
	*/

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MakeStockBufferAnalysis(_ctx context.Context, _input *MakeStockBufferInput, _output *MakeStockBufferOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _MakeStockBufferValidation(_ctx context.Context, _input *MakeStockBufferInput, _output *MakeStockBufferOutput) {
}
func _MakeStockBufferRun(_ctx context.Context, input *MakeStockBufferInput) *MakeStockBufferOutput {
	output := &MakeStockBufferOutput{}
	_MakeStockBufferSetup(_ctx, input)
	_MakeStockBufferSteps(_ctx, input, output)
	_MakeStockBufferAnalysis(_ctx, input, output)
	_MakeStockBufferValidation(_ctx, input, output)
	return output
}

func MakeStockBufferRunSteps(_ctx context.Context, input *MakeStockBufferInput) *MakeStockBufferSOutput {
	soutput := &MakeStockBufferSOutput{}
	output := _MakeStockBufferRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MakeStockBufferNew() interface{} {
	return &MakeStockBufferElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MakeStockBufferInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MakeStockBufferRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MakeStockBufferInput{},
			Out: &MakeStockBufferOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type MakeStockBufferElement struct {
	inject.CheckedRunner
}

type MakeStockBufferInput struct {
	Diluent      *wtype.LHComponent
	MassAddedinG wunit.Mass
	Moleculename string
	TotalVolume  wunit.Volume
}

type MakeStockBufferOutput struct {
	StockConc wunit.Concentration
}

type MakeStockBufferSOutput struct {
	Data struct {
		StockConc wunit.Concentration
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MakeStockBuffer",
		Constructor: MakeStockBufferNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/MakeBuffer/MakeStockBuffer.an",
			Params: []component.ParamDesc{
				{Name: "Diluent", Desc: "Bufferstock\t\t*wtype.LHComponent\n", Kind: "Inputs"},
				{Name: "MassAddedinG", Desc: "", Kind: "Parameters"},
				{Name: "Moleculename", Desc: "", Kind: "Parameters"},
				{Name: "TotalVolume", Desc: "", Kind: "Parameters"},
				{Name: "StockConc", Desc: "Status string\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}

/*
type Mole struct {
	number float64
}*/
