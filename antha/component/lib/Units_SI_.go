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

func _Units_SIRequirements() {

}

// Actions to perform before protocol itself
func _Units_SISetup(_ctx context.Context, _input *Units_SIInput) {

}

// Core process of the protocol: steps to be performed for each input
func _Units_SISteps(_ctx context.Context, _input *Units_SIInput, _output *Units_SIOutput) {

	_output.SIMass = _input.MyMass.SIValue()
	_output.SIMassUnit = _input.MyMass.Unit().BaseSISymbol()

	_output.RawMass = _input.MyMass.RawValue()
	_output.RawMassUnit = _input.MyMass.Unit().PrefixedSymbol()

}

// Actions to perform after steps block to analyze data
func _Units_SIAnalysis(_ctx context.Context, _input *Units_SIInput, _output *Units_SIOutput) {

}

func _Units_SIValidation(_ctx context.Context, _input *Units_SIInput, _output *Units_SIOutput) {

}
func _Units_SIRun(_ctx context.Context, input *Units_SIInput) *Units_SIOutput {
	output := &Units_SIOutput{}
	_Units_SISetup(_ctx, input)
	_Units_SISteps(_ctx, input, output)
	_Units_SIAnalysis(_ctx, input, output)
	_Units_SIValidation(_ctx, input, output)
	return output
}

func Units_SIRunSteps(_ctx context.Context, input *Units_SIInput) *Units_SISOutput {
	soutput := &Units_SISOutput{}
	output := _Units_SIRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Units_SINew() interface{} {
	return &Units_SIElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Units_SIInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Units_SIRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Units_SIInput{},
			Out: &Units_SIOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Units_SIElement struct {
	inject.CheckedRunner
}

type Units_SIInput struct {
	MyMass wunit.Mass
}

type Units_SIOutput struct {
	RawMass     float64
	RawMassUnit string
	SIMass      float64
	SIMassUnit  string
}

type Units_SISOutput struct {
	Data struct {
		RawMass     float64
		RawMassUnit string
		SIMass      float64
		SIMassUnit  string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Units_SI",
		Constructor: Units_SINew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/AnthaAcademy/Lesson0_Units/C_units_SI.an",
			Params: []component.ParamDesc{
				{Name: "MyMass", Desc: "", Kind: "Parameters"},
				{Name: "RawMass", Desc: "", Kind: "Data"},
				{Name: "RawMassUnit", Desc: "", Kind: "Data"},
				{Name: "SIMass", Desc: "", Kind: "Data"},
				{Name: "SIMassUnit", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
