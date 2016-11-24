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

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _BasicTransferRequirements() {

}

// Conditions to run on startup
func _BasicTransferSetup(_ctx context.Context, _input *BasicTransferInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _BasicTransferSteps(_ctx context.Context, _input *BasicTransferInput, _output *BasicTransferOutput) {
	_output.Status = ""
	fmt.Println("source plate ID: ", _input.SourcePlate.ID)
	for k := 0; k < len(_input.SourceWell); k++ {

		sourceCoords := wtype.MakeWellCoordsA1(_input.SourceWell[k])
		sourceLiquid := _input.SourcePlate.WellAt(sourceCoords).WContents
		fmt.Println("SOURCE LIQUID ID: ", sourceLiquid.ID, " ", sourceLiquid.CName, " LOC: ", sourceLiquid.Loc)

		sample := mixer.Sample(sourceLiquid, _input.LiquidVolume[k])

		destCoords := wtype.MakeWellCoordsA1(_input.DestinationWell[k])
		destLiquid := _input.DestinationPlate.WellAt(destCoords).WContents

		var mixed *wtype.LHComponent
		var present bool

		if destLiquid.IsZero() {
			present = false
			mixed = execute.MixInto(_ctx, _input.DestinationPlate, _input.DestinationWell[k], sample)
		} else {
			present = true
			mixed = execute.Mix(_ctx, destLiquid, sample)

		}
		//fmt.Println("Sampling: " + sourceSolution.String() + " Resulting sample: " + sample.String())

		// this is very unsafe. We need to make sure this works properl
		mixed.Loc = _input.DestinationPlate.ID + ":" + _input.DestinationWell[k]
		_input.DestinationPlate.WellAt(wtype.MakeWellCoordsA1(_input.DestinationWell[k])).WContents = mixed
		fmt.Println(_input.DestinationPlate.Wellcoords[_input.DestinationWell[k]].WContents.ID, " ", _input.DestinationPlate.Wellcoords[_input.DestinationWell[k]].WContents.CName, " ", _input.DestinationPlate.Wellcoords[_input.DestinationWell[k]].WContents.Loc)

		_output.Status = _output.Status + _input.LiquidVolume[k].ToString() + " of " + sourceLiquid.Name() + " from " + _input.SourceWell[k] + " on " + _input.SourcePlate.Name() + " dispensed to " + _input.DestinationWell[k] + " on " + _input.DestinationPlate.Name() + " "
		if present {
			_output.Status = _output.Status + "which contained " + destLiquid.Name() + " "
		}
	}
	if !_input.DeadEnd {
		_output.DestinationPlateOut = _input.DestinationPlate
		fmt.Println("Destination plate ID passed out: ", _output.DestinationPlateOut.ID)
	}
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _BasicTransferAnalysis(_ctx context.Context, _input *BasicTransferInput, _output *BasicTransferOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _BasicTransferValidation(_ctx context.Context, _input *BasicTransferInput, _output *BasicTransferOutput) {

}
func _BasicTransferRun(_ctx context.Context, input *BasicTransferInput) *BasicTransferOutput {
	output := &BasicTransferOutput{}
	_BasicTransferSetup(_ctx, input)
	_BasicTransferSteps(_ctx, input, output)
	_BasicTransferAnalysis(_ctx, input, output)
	_BasicTransferValidation(_ctx, input, output)
	return output
}

func BasicTransferRunSteps(_ctx context.Context, input *BasicTransferInput) *BasicTransferSOutput {
	soutput := &BasicTransferSOutput{}
	output := _BasicTransferRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func BasicTransferNew() interface{} {
	return &BasicTransferElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &BasicTransferInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _BasicTransferRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &BasicTransferInput{},
			Out: &BasicTransferOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type BasicTransferElement struct {
	inject.CheckedRunner
}

type BasicTransferInput struct {
	DeadEnd          bool
	DestinationPlate *wtype.LHPlate
	DestinationWell  []string
	LiquidVolume     []wunit.Volume
	SourcePlate      *wtype.LHPlate
	SourceWell       []string
}

type BasicTransferOutput struct {
	DestinationPlateOut *wtype.LHPlate
	Status              string
}

type BasicTransferSOutput struct {
	Data struct {
		Status string
	}
	Outputs struct {
		DestinationPlateOut *wtype.LHPlate
	}
}

func init() {
	if err := addComponent(component.Component{Name: "BasicTransfer",
		Constructor: BasicTransferNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/BS/BasicTransfer.an",
			Params: []component.ParamDesc{
				{Name: "DeadEnd", Desc: "", Kind: "Parameters"},
				{Name: "DestinationPlate", Desc: "", Kind: "Inputs"},
				{Name: "DestinationWell", Desc: "", Kind: "Parameters"},
				{Name: "LiquidVolume", Desc: "", Kind: "Parameters"},
				{Name: "SourcePlate", Desc: "", Kind: "Inputs"},
				{Name: "SourceWell", Desc: "", Kind: "Parameters"},
				{Name: "DestinationPlateOut", Desc: "", Kind: "Outputs"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
