// example of how to look up molecule properties from pubchem
package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Pubchem"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Name of compound or array of multiple compounds

// molecule type is returned consisting of name, formula, molecular weight and chemical ID (CID)

// or JSON structure if preferred

// status to be printed out in manual driver console

func _LookUpMoleculeRequirements() {
}
func _LookUpMoleculeSetup(_ctx context.Context, _input *LookUpMoleculeInput) {
}
func _LookUpMoleculeSteps(_ctx context.Context, _input *LookUpMoleculeInput, _output *LookUpMoleculeOutput) {
	var err error

	// method of making molecule from name
	_output.Compoundprops, err = pubchem.MakeMolecule(_input.Compound)
	if err != nil {
		panic(err)
	}

	// or returning properties in JSON structure
	_output.Jsonstring, err = pubchem.Compoundproperties(_input.Compound)
	if err != nil {
		panic(err)
	}

	// method of making a list of compounds from names
	_output.List, err = pubchem.MakeMolecules(_input.Compoundlist)
	if err != nil {
		panic(err)
	}

	// Print out status
	_output.Status = fmt.Sprintln("Returned data from",
		_input.Compound, "=",
		_output.Compoundprops.Moleculename,
		_output.Compoundprops.MolecularWeight,
		_output.Compoundprops.MolecularFormula,
		_output.Compoundprops.CID,
		"Data in JSON format =", _output.Jsonstring,
		"List=", _output.List)
}
func _LookUpMoleculeAnalysis(_ctx context.Context, _input *LookUpMoleculeInput, _output *LookUpMoleculeOutput) {

}

func _LookUpMoleculeValidation(_ctx context.Context, _input *LookUpMoleculeInput, _output *LookUpMoleculeOutput) {

}
func _LookUpMoleculeRun(_ctx context.Context, input *LookUpMoleculeInput) *LookUpMoleculeOutput {
	output := &LookUpMoleculeOutput{}
	_LookUpMoleculeSetup(_ctx, input)
	_LookUpMoleculeSteps(_ctx, input, output)
	_LookUpMoleculeAnalysis(_ctx, input, output)
	_LookUpMoleculeValidation(_ctx, input, output)
	return output
}

func LookUpMoleculeRunSteps(_ctx context.Context, input *LookUpMoleculeInput) *LookUpMoleculeSOutput {
	soutput := &LookUpMoleculeSOutput{}
	output := _LookUpMoleculeRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func LookUpMoleculeNew() interface{} {
	return &LookUpMoleculeElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &LookUpMoleculeInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _LookUpMoleculeRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &LookUpMoleculeInput{},
			Out: &LookUpMoleculeOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type LookUpMoleculeElement struct {
	inject.CheckedRunner
}

type LookUpMoleculeInput struct {
	Compound     string
	Compoundlist []string
}

type LookUpMoleculeOutput struct {
	Compoundprops pubchem.Molecule
	Jsonstring    string
	List          []pubchem.Molecule
	Status        string
}

type LookUpMoleculeSOutput struct {
	Data struct {
		Compoundprops pubchem.Molecule
		Jsonstring    string
		List          []pubchem.Molecule
		Status        string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "LookUpMolecule",
		Constructor: LookUpMoleculeNew,
		Desc: component.ComponentDesc{
			Desc: "example of how to look up molecule properties from pubchem\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/LookUpMolecule/LookUpMolecule.an",
			Params: []component.ParamDesc{
				{Name: "Compound", Desc: "Name of compound or array of multiple compounds\n", Kind: "Parameters"},
				{Name: "Compoundlist", Desc: "", Kind: "Parameters"},
				{Name: "Compoundprops", Desc: "molecule type is returned consisting of name, formula, molecular weight and chemical ID (CID)\n", Kind: "Data"},
				{Name: "Jsonstring", Desc: "or JSON structure if preferred\n", Kind: "Data"},
				{Name: "List", Desc: "", Kind: "Data"},
				{Name: "Status", Desc: "status to be printed out in manual driver console\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
