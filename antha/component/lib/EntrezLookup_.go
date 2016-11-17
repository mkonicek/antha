package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/entrez"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Valid Database list: http://www.ncbi.nlm.nih.gov/books/NBK25497/table/chapter2.T._entrez_unique_identifiers_ui/?report=objectonly

// Valid ReturnType List: http://www.ncbi.nlm.nih.gov/books/NBK25499/table/chapter4.T._valid_values_of__retmode_and/?report=objectonly

// Input parameters for this protocol

// e.g. "EF208560"
// e.g. "nucleotide", "Protein", "Gene"
// e.g. 1
// e.g. "gb", "fasta"
// e.g myproject/GFPReporter.gb. if Filename == "" no file will be generated

// Data which is returned from this protocol

// Physical inputs to this protocol

// Physical outputs from this protocol

func _EntrezLookupRequirements() {

}

// Actions to perform before protocol itself
func _EntrezLookupSetup(_ctx context.Context, _input *EntrezLookupInput) {

}

// Core process of the protocol: steps to be performed for each input
func _EntrezLookupSteps(_ctx context.Context, _input *EntrezLookupInput, _output *EntrezLookupOutput) {

	var output []byte

	_output.Filenameused, output, _output.Err = entrez.RetrieveRecords(_input.ID, _input.Database, _input.MaxReturns, _input.ReturnType, _input.Filename)

	_output.Output = string(output)
}

// Actions to perform after steps block to analyze data
func _EntrezLookupAnalysis(_ctx context.Context, _input *EntrezLookupInput, _output *EntrezLookupOutput) {

}

func _EntrezLookupValidation(_ctx context.Context, _input *EntrezLookupInput, _output *EntrezLookupOutput) {

}
func _EntrezLookupRun(_ctx context.Context, input *EntrezLookupInput) *EntrezLookupOutput {
	output := &EntrezLookupOutput{}
	_EntrezLookupSetup(_ctx, input)
	_EntrezLookupSteps(_ctx, input, output)
	_EntrezLookupAnalysis(_ctx, input, output)
	_EntrezLookupValidation(_ctx, input, output)
	return output
}

func EntrezLookupRunSteps(_ctx context.Context, input *EntrezLookupInput) *EntrezLookupSOutput {
	soutput := &EntrezLookupSOutput{}
	output := _EntrezLookupRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func EntrezLookupNew() interface{} {
	return &EntrezLookupElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &EntrezLookupInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _EntrezLookupRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &EntrezLookupInput{},
			Out: &EntrezLookupOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type EntrezLookupElement struct {
	inject.CheckedRunner
}

type EntrezLookupInput struct {
	Database   string
	Filename   string
	ID         string
	MaxReturns int
	ReturnType string
}

type EntrezLookupOutput struct {
	Err          error
	Filenameused string
	Output       string
}

type EntrezLookupSOutput struct {
	Data struct {
		Err          error
		Filenameused string
		Output       string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "EntrezLookup",
		Constructor: EntrezLookupNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/DNA/EntrezLookup/EntrezLookup.an",
			Params: []component.ParamDesc{
				{Name: "Database", Desc: "e.g. \"nucleotide\", \"Protein\", \"Gene\"\n", Kind: "Parameters"},
				{Name: "Filename", Desc: "e.g myproject/GFPReporter.gb. if Filename == \"\" no file will be generated\n", Kind: "Parameters"},
				{Name: "ID", Desc: "e.g. \"EF208560\"\n", Kind: "Parameters"},
				{Name: "MaxReturns", Desc: "e.g. 1\n", Kind: "Parameters"},
				{Name: "ReturnType", Desc: "e.g. \"gb\", \"fasta\"\n", Kind: "Parameters"},
				{Name: "Err", Desc: "", Kind: "Data"},
				{Name: "Filenameused", Desc: "", Kind: "Data"},
				{Name: "Output", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
