package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Pubchem"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
	"strings"
)

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

//OriginalDiluentVolume Volume

// Input Requirement specification
func _MakeStockBuffer2Requirements() {

}

// Conditions to run on startup
func _MakeStockBuffer2Setup(_ctx context.Context, _input *MakeStockBuffer2Input) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _MakeStockBuffer2Steps(_ctx context.Context, _input *MakeStockBuffer2Input, _output *MakeStockBuffer2Output) {

	var err error

	molecule, err := pubchem.MakeMolecule(_input.Molecule.CName)

	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	_output.ConcinGperL = _input.TargetConc.GramPerL(molecule.MolecularWeight)
	_output.ConcinMperL = _input.TargetConc.MolPerL(molecule.MolecularWeight)

	_output.MassToAddinG, err = wunit.MassForTargetConcentration(_input.TargetConc, _input.TotalVolume)

	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	_output.Buffer = execute.MixInto(_ctx, _input.Vessel, "",
		mixer.SampleMass(_input.Molecule, _output.MassToAddinG, _input.MoleculeDensity),
		mixer.SampleForTotalVolume(_input.Diluent, _input.TotalVolume))

	_output.Buffer.CName = _output.ConcinMperL.ToString() + " " + _input.Molecule.CName

	_output.Status = fmt.Sprintln("Stock added = ", _output.MassToAddinG.ToString(), "of", _input.Molecule.CName,
		"was added up to ", _input.TotalVolume.SIValue(), "L with ", _input.Diluent.CName,
		"to make ", _input.TotalVolume.SIValue(), "L", "of", _output.Buffer.CName,
		"Buffer stock conc =", _input.TargetConc.ToString(), ". Extra instructions: ", strings.Join(_input.ExtraInstructions, ", "), ". Store at ", _input.StorageTemperature.ToString())

	_output.MoleculeInfo = molecule

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MakeStockBuffer2Analysis(_ctx context.Context, _input *MakeStockBuffer2Input, _output *MakeStockBuffer2Output) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _MakeStockBuffer2Validation(_ctx context.Context, _input *MakeStockBuffer2Input, _output *MakeStockBuffer2Output) {
}
func _MakeStockBuffer2Run(_ctx context.Context, input *MakeStockBuffer2Input) *MakeStockBuffer2Output {
	output := &MakeStockBuffer2Output{}
	_MakeStockBuffer2Setup(_ctx, input)
	_MakeStockBuffer2Steps(_ctx, input, output)
	_MakeStockBuffer2Analysis(_ctx, input, output)
	_MakeStockBuffer2Validation(_ctx, input, output)
	return output
}

func MakeStockBuffer2RunSteps(_ctx context.Context, input *MakeStockBuffer2Input) *MakeStockBuffer2SOutput {
	soutput := &MakeStockBuffer2SOutput{}
	output := _MakeStockBuffer2Run(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MakeStockBuffer2New() interface{} {
	return &MakeStockBuffer2Element{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MakeStockBuffer2Input{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MakeStockBuffer2Run(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MakeStockBuffer2Input{},
			Out: &MakeStockBuffer2Output{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type MakeStockBuffer2Element struct {
	inject.CheckedRunner
}

type MakeStockBuffer2Input struct {
	Diluent            *wtype.LHComponent
	ExtraInstructions  []string
	Molecule           *wtype.LHComponent
	MoleculeDensity    wunit.Density
	StorageTemperature wunit.Temperature
	TargetConc         wunit.Concentration
	TotalVolume        wunit.Volume
	Vessel             *wtype.LHPlate
}

type MakeStockBuffer2Output struct {
	Buffer       *wtype.LHComponent
	ConcinGperL  wunit.Concentration
	ConcinMperL  wunit.Concentration
	MassToAddinG wunit.Mass
	MoleculeInfo pubchem.Molecule
	Status       string
}

type MakeStockBuffer2SOutput struct {
	Data struct {
		ConcinGperL  wunit.Concentration
		ConcinMperL  wunit.Concentration
		MassToAddinG wunit.Mass
		MoleculeInfo pubchem.Molecule
		Status       string
	}
	Outputs struct {
		Buffer *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MakeStockBuffer2",
		Constructor: MakeStockBuffer2New,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/MakeBuffer/MakeStockBuffer2.an",
			Params: []component.ParamDesc{
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "ExtraInstructions", Desc: "", Kind: "Parameters"},
				{Name: "Molecule", Desc: "", Kind: "Inputs"},
				{Name: "MoleculeDensity", Desc: "", Kind: "Parameters"},
				{Name: "StorageTemperature", Desc: "", Kind: "Parameters"},
				{Name: "TargetConc", Desc: "", Kind: "Parameters"},
				{Name: "TotalVolume", Desc: "", Kind: "Parameters"},
				{Name: "Vessel", Desc: "", Kind: "Inputs"},
				{Name: "Buffer", Desc: "", Kind: "Outputs"},
				{Name: "ConcinGperL", Desc: "", Kind: "Data"},
				{Name: "ConcinMperL", Desc: "", Kind: "Data"},
				{Name: "MassToAddinG", Desc: "", Kind: "Data"},
				{Name: "MoleculeInfo", Desc: "", Kind: "Data"},
				{Name: "Status", Desc: "", Kind: "Data"},
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
