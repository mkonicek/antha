// demo protocol of how to convert units to string
package lib

import (
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol

// Data which is returned from this protocol

// Physical inputs to this protocol

// Physical outputs from this protocol

func _Units_ToStringRequirements() {

}

// Actions to perform before protocol itself
func _Units_ToStringSetup(_ctx context.Context, _input *Units_ToStringInput) {

}

// Core process of the protocol: steps to be performed for each input
func _Units_ToStringSteps(_ctx context.Context, _input *Units_ToStringInput, _output *Units_ToStringOutput) {
	_output.VolumeasString = _input.MyVolume.ToString()
	_output.TempasString = _input.MyTemperature.ToString()
	_output.ConcasString = _input.MyConc.ToString()
	_output.MolesasString = _input.MyMoles.ToString()
	_output.MassasString = _input.MyMass.ToString()
	_output.FlowrateString = _input.MyFlowrate.ToString()

	// Exercise: Add an equivalent process for a FlowRate
}

// Actions to perform after steps block to analyze data
func _Units_ToStringAnalysis(_ctx context.Context, _input *Units_ToStringInput, _output *Units_ToStringOutput) {

}

func _Units_ToStringValidation(_ctx context.Context, _input *Units_ToStringInput, _output *Units_ToStringOutput) {

}
func _Units_ToStringRun(_ctx context.Context, input *Units_ToStringInput) *Units_ToStringOutput {
	output := &Units_ToStringOutput{}
	_Units_ToStringSetup(_ctx, input)
	_Units_ToStringSteps(_ctx, input, output)
	_Units_ToStringAnalysis(_ctx, input, output)
	_Units_ToStringValidation(_ctx, input, output)
	return output
}

func Units_ToStringRunSteps(_ctx context.Context, input *Units_ToStringInput) *Units_ToStringSOutput {
	soutput := &Units_ToStringSOutput{}
	output := _Units_ToStringRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Units_ToStringNew() interface{} {
	return &Units_ToStringElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Units_ToStringInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Units_ToStringRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Units_ToStringInput{},
			Out: &Units_ToStringOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Units_ToStringElement struct {
	inject.CheckedRunner
}

type Units_ToStringInput struct {
	MyConc        wunit.Concentration
	MyFlowrate    wunit.FlowRate
	MyMass        wunit.Mass
	MyMoles       wunit.Moles
	MyTemperature wunit.Temperature
	MyTime        wunit.Time
	MyVolume      wunit.Volume
}

type Units_ToStringOutput struct {
	ConcasString   string
	FlowrateString string
	MassasString   string
	MolesasString  string
	TempasString   string
	VolumeasString string
}

type Units_ToStringSOutput struct {
	Data struct {
		ConcasString   string
		FlowrateString string
		MassasString   string
		MolesasString  string
		TempasString   string
		VolumeasString string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Units_ToString",
		Constructor: Units_ToStringNew,
		Desc: component.ComponentDesc{
			Desc: "demo protocol of how to convert units to string\n",
			Path: "antha/component/an/AnthaAcademy/Lesson4_Units/B_units_toString.an",
			Params: []component.ParamDesc{
				{Name: "MyConc", Desc: "", Kind: "Parameters"},
				{Name: "MyFlowrate", Desc: "", Kind: "Parameters"},
				{Name: "MyMass", Desc: "", Kind: "Parameters"},
				{Name: "MyMoles", Desc: "", Kind: "Parameters"},
				{Name: "MyTemperature", Desc: "", Kind: "Parameters"},
				{Name: "MyTime", Desc: "", Kind: "Parameters"},
				{Name: "MyVolume", Desc: "", Kind: "Parameters"},
				{Name: "ConcasString", Desc: "", Kind: "Data"},
				{Name: "FlowrateString", Desc: "", Kind: "Data"},
				{Name: "MassasString", Desc: "", Kind: "Data"},
				{Name: "MolesasString", Desc: "", Kind: "Data"},
				{Name: "TempasString", Desc: "", Kind: "Data"},
				{Name: "VolumeasString", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
