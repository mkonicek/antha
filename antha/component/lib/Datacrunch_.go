//Some examples functions
// Calculate rate of reaction, V, of enzyme displaying Micahelis-Menten kinetics with Vmax, Km and [S] declared
// Calculating [S] and V from g/l concentration and looking up molecular weight of named substrate
// Calculating [S] and V from g/l concentration of DNA of known sequence
// Calculating [S] and V from g/l concentration of Protein product of DNA of known sequence

package lib

import (
	"fmt"
	//"math"
	//"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Pubchem"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol

//Amount
// i.e. Moles, M

//Amount

// Data which is returned from this protocol

// Physical inputs to this protocol

// Physical outputs from this protocol

func _DatacrunchRequirements() {

}

// Actions to perform before protocol itself
func _DatacrunchSetup(_ctx context.Context, _input *DatacrunchInput) {

}

// Core process of the protocol: steps to be performed for each input
func _DatacrunchSteps(_ctx context.Context, _input *DatacrunchInput, _output *DatacrunchOutput) {
	// Work out rate of reaction, V of enzyme with Michaelis-Menten kinetics and [S], Km and Vmax declared
	//Using declared values for S and unit of S
	km := wunit.NewAmount(_input.Km, _input.Kmunit) //.SIValue()
	s := wunit.NewAmount(_input.S, _input.Sunit)    //.SIValue()

	_output.V = ((s.SIValue() * _input.Vmax) / (s.SIValue() + km.SIValue()))

	// Now working out Molarity of Substrate based on conc and looking up molecular weight in pubchem

	// Look up properties
	substrate_mw, err := pubchem.MakeMolecule(_input.Substrate_name)
	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	// calculate moles
	submoles := sequences.Moles(_input.SubstrateConc, substrate_mw.MolecularWeight, _input.SubstrateVol)
	// calculate molar concentration
	submolarconc := sequences.GtoMolarConc(_input.SubstrateConc, substrate_mw.MolecularWeight)

	// make a new amount
	s = wunit.NewAmount(submolarconc, "M")

	// use michaelis menton equation
	v_substrate_name := ((s.SIValue() * _input.Vmax) / (s.SIValue() + km.SIValue()))

	// Now working out Molarity of Substrate from DNA Sequence
	// calculate molar concentration
	dna_mw := sequences.MassDNA(_input.DNA_seq, false, false)
	dnamolarconc := sequences.GtoMolarConc(_input.DNAConc, dna_mw)

	// make a new amount
	s = wunit.NewAmount(dnamolarconc, "M")

	// use michaelis menton equation
	v_dna := ((s.SIValue() * _input.Vmax) / (s.SIValue() + km.SIValue()))

	// Now working out Molarity of Substrate from Protein product of dna Sequence

	// translate
	orf, orftrue := sequences.FindORF(_input.DNA_seq)
	var protein_mw float64
	if orftrue == true {
		protein_mw_kDA := sequences.Molecularweight(orf)
		protein_mw = protein_mw_kDA * 1000
		_output.Orftrue = orftrue
	}

	// calculate molar concentration
	proteinmolarconc := sequences.GtoMolarConc(_input.ProteinConc, protein_mw)

	// make a new amount
	s = wunit.NewAmount(submolarconc, "M")

	// use michaelis menton equation
	v_protein := ((s.SIValue() * _input.Vmax) / (s.SIValue() + km.SIValue()))

	// print report
	_output.Status = fmt.Sprintln(
		"Rate, V of enzyme at substrate conc", _input.S, _input.Sunit,
		"of enzyme with Km", km.ToString(),
		"and Vmax", _input.Vmax, _input.Vmaxunit,
		"=", _output.V, _input.Vunit, ".",
		"Substrate =", _input.Substrate_name, ". We have", _input.SubstrateVol.ToString(), "of", _input.Substrate_name, "at concentration of", _input.SubstrateConc.ToString(),
		"Therefore... Moles of", _input.Substrate_name, "=", submoles, "Moles.",
		"Molar Concentration of", _input.Substrate_name, "=", submolarconc, "Mol/L.",
		"Rate, V = ", v_substrate_name, _input.Vmaxunit,
		"Substrate =", "DNA Sequence of", _input.Gene_name, "We have", "concentration of", _input.DNAConc.ToString(),
		"Therefore... Molar conc", "=", dnamolarconc, "Mol/L",
		"Rate, V = ", v_dna, _input.Vmaxunit,
		"Substrate =", "protein from DNA sequence", _input.Gene_name, ".",
		"We have", "concentration of", _input.ProteinConc.ToString(),
		"Therefore... Molar conc", "=", proteinmolarconc, "Mol/L",
		"Rate, V = ", v_protein, _input.Vmaxunit)
}

// Actions to perform after steps block to analyze data
func _DatacrunchAnalysis(_ctx context.Context, _input *DatacrunchInput, _output *DatacrunchOutput) {

}

func _DatacrunchValidation(_ctx context.Context, _input *DatacrunchInput, _output *DatacrunchOutput) {

}
func _DatacrunchRun(_ctx context.Context, input *DatacrunchInput) *DatacrunchOutput {
	output := &DatacrunchOutput{}
	_DatacrunchSetup(_ctx, input)
	_DatacrunchSteps(_ctx, input, output)
	_DatacrunchAnalysis(_ctx, input, output)
	_DatacrunchValidation(_ctx, input, output)
	return output
}

func DatacrunchRunSteps(_ctx context.Context, input *DatacrunchInput) *DatacrunchSOutput {
	soutput := &DatacrunchSOutput{}
	output := _DatacrunchRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func DatacrunchNew() interface{} {
	return &DatacrunchElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &DatacrunchInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _DatacrunchRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &DatacrunchInput{},
			Out: &DatacrunchOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type DatacrunchElement struct {
	inject.CheckedRunner
}

type DatacrunchInput struct {
	DNAConc        wunit.Concentration
	DNA_seq        string
	Gene_name      string
	Km             float64
	Kmunit         string
	ProteinConc    wunit.Concentration
	S              float64
	SubstrateConc  wunit.Concentration
	SubstrateVol   wunit.Volume
	Substrate_name string
	Sunit          string
	Vmax           float64
	Vmaxunit       string
	Vunit          string
}

type DatacrunchOutput struct {
	Orftrue bool
	Status  string
	V       float64
}

type DatacrunchSOutput struct {
	Data struct {
		Orftrue bool
		Status  string
		V       float64
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Datacrunch",
		Constructor: DatacrunchNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/AnthaAcademy/Lesson5_Units2/A_Datacrunch.an",
			Params: []component.ParamDesc{
				{Name: "DNAConc", Desc: "", Kind: "Parameters"},
				{Name: "DNA_seq", Desc: "", Kind: "Parameters"},
				{Name: "Gene_name", Desc: "", Kind: "Parameters"},
				{Name: "Km", Desc: "Amount\n", Kind: "Parameters"},
				{Name: "Kmunit", Desc: "i.e. Moles, M\n", Kind: "Parameters"},
				{Name: "ProteinConc", Desc: "", Kind: "Parameters"},
				{Name: "S", Desc: "Amount\n", Kind: "Parameters"},
				{Name: "SubstrateConc", Desc: "", Kind: "Parameters"},
				{Name: "SubstrateVol", Desc: "", Kind: "Parameters"},
				{Name: "Substrate_name", Desc: "", Kind: "Parameters"},
				{Name: "Sunit", Desc: "", Kind: "Parameters"},
				{Name: "Vmax", Desc: "", Kind: "Parameters"},
				{Name: "Vmaxunit", Desc: "", Kind: "Parameters"},
				{Name: "Vunit", Desc: "", Kind: "Parameters"},
				{Name: "Orftrue", Desc: "", Kind: "Data"},
				{Name: "Status", Desc: "", Kind: "Data"},
				{Name: "V", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
