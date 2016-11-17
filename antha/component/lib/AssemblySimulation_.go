package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

func _AssemblySimulationRequirements() {
}

func _AssemblySimulationSetup(_ctx context.Context, _input *AssemblySimulationInput) {
}

func _AssemblySimulationSteps(_ctx context.Context, _input *AssemblySimulationInput, _output *AssemblySimulationOutput) {

	// Assembly parameters
	assembly := enzymes.Assemblyparameters{"Simulated", _input.RE, _input.VectorSeq, _input.PartsWithOverhangs}

	// Simulation
	_output.SimulationStatus, _output.NumberofSuccessfulAssemblies, _output.RestrictionSitesFound, _output.SimulatedSequence, _output.Warnings = enzymes.Assemblysimulator(assembly)

}

func _AssemblySimulationAnalysis(_ctx context.Context, _input *AssemblySimulationInput, _output *AssemblySimulationOutput) {

}

func _AssemblySimulationValidation(_ctx context.Context, _input *AssemblySimulationInput, _output *AssemblySimulationOutput) {

}
func _AssemblySimulationRun(_ctx context.Context, input *AssemblySimulationInput) *AssemblySimulationOutput {
	output := &AssemblySimulationOutput{}
	_AssemblySimulationSetup(_ctx, input)
	_AssemblySimulationSteps(_ctx, input, output)
	_AssemblySimulationAnalysis(_ctx, input, output)
	_AssemblySimulationValidation(_ctx, input, output)
	return output
}

func AssemblySimulationRunSteps(_ctx context.Context, input *AssemblySimulationInput) *AssemblySimulationSOutput {
	soutput := &AssemblySimulationSOutput{}
	output := _AssemblySimulationRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AssemblySimulationNew() interface{} {
	return &AssemblySimulationElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AssemblySimulationInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AssemblySimulationRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AssemblySimulationInput{},
			Out: &AssemblySimulationOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type AssemblySimulationElement struct {
	inject.CheckedRunner
}

type AssemblySimulationInput struct {
	PartsWithOverhangs []wtype.DNASequence
	RE                 string
	SynthesisProvider  string
	VectorSeq          wtype.DNASequence
}

type AssemblySimulationOutput struct {
	NumberofSuccessfulAssemblies int
	RestrictionSitesFound        []enzymes.Restrictionsites
	SimulatedSequence            wtype.DNASequence
	SimulationStatus             string
	Validated                    bool
	ValidationStatus             string
	Warnings                     error
}

type AssemblySimulationSOutput struct {
	Data struct {
		NumberofSuccessfulAssemblies int
		RestrictionSitesFound        []enzymes.Restrictionsites
		SimulatedSequence            wtype.DNASequence
		SimulationStatus             string
		Validated                    bool
		ValidationStatus             string
		Warnings                     error
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "AssemblySimulation",
		Constructor: AssemblySimulationNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/DNA/GeneDesign/AssemblySimulation.an",
			Params: []component.ParamDesc{
				{Name: "PartsWithOverhangs", Desc: "", Kind: "Parameters"},
				{Name: "RE", Desc: "", Kind: "Parameters"},
				{Name: "SynthesisProvider", Desc: "", Kind: "Parameters"},
				{Name: "VectorSeq", Desc: "", Kind: "Parameters"},
				{Name: "NumberofSuccessfulAssemblies", Desc: "", Kind: "Data"},
				{Name: "RestrictionSitesFound", Desc: "", Kind: "Data"},
				{Name: "SimulatedSequence", Desc: "", Kind: "Data"},
				{Name: "SimulationStatus", Desc: "", Kind: "Data"},
				{Name: "Validated", Desc: "", Kind: "Data"},
				{Name: "ValidationStatus", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
