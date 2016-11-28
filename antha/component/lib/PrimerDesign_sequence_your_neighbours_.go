// This element will design outward facing primers for all .gb file sequences in a specified folder.
// Design criteria such as maximum gc content, acceptable ranges of melting temperatures and primer length may be specified by the user.
package lib

import (
	"fmt"
	//"math"
	//"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/export"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/oligos"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
	"os"
	"path/filepath"
	"strings"
)

// Input parameters for this protocol

//files     []string = []string{"STAR_0023_VECTOR_BBSI.gb", "STAR_0023_VECTOR_BBSI+Grp7+Grp14+Grp3.gb"}
//= "current" // this will check for all .gb files in the folder you select here
//      = wunit.NewTemperature(60, "C")
//      = wunit.NewTemperature(55, "C")
//      = 0.6
//     = 20
//     = 25
// number of nucleotides which primers can overlap by

// Data which is returned from this protocol

//PrimerData []oligos.Primer

// Physical inputs to this protocol

// Physical outputs from this protocol

func _PrimerDesign_sequence_your_neighboursRequirements() {

}

// Actions to perform before protocol itself
func _PrimerDesign_sequence_your_neighboursSetup(_ctx context.Context, _input *PrimerDesign_sequence_your_neighboursInput) {

}

// Core process of the protocol: steps to be performed for each input
func _PrimerDesign_sequence_your_neighboursSteps(_ctx context.Context, _input *PrimerDesign_sequence_your_neighboursInput, _output *PrimerDesign_sequence_your_neighboursOutput) {

	//var Start int
	//var End int
	var err error
	var output string
	var dirname string
	var alloutputs = make([]string, 0)
	var allprimers = make([]oligos.Primer, 0)
	var allprimerstrings = make([]string, 0)
	var primerpairs = make([]PrimerPair, 0)
	var files = make([]string, 0)

	//Search for files within current directory

	if _input.Dirname == "current" {
		dirname = "." + string(filepath.Separator)
	} else {
		dirname = _input.Dirname
	}

	d, err := os.Open(dirname)
	if err != nil {
		panic(err)
	}
	defer d.Close()

	allfiles, err := d.Readdir(-1)
	if err != nil {
		panic(err)
	}

	fmt.Println("Reading " + dirname)

	//Determine if file extension is ".gb"
	for _, file := range allfiles {
		if filepath.Ext(file.Name()) == ".gb" {
			files = append(files, file.Name())
		}

	}

	for _, file := range files {
		file = filepath.Join(dirname, file)
		sequence, _ := parser.GenbanktoAnnotatedSeq(file)

		primer1, primer2 := oligos.MakeOutwardFacingPrimers(sequence, _input.Maxgc, _input.Minlength, _input.Maxlength, _input.Mintemp, _input.Maxtemp, allprimerstrings, _input.PermittednucleotideOverlapBetweenPrimers)

		primer1.Nm = "primer1" + "_" + file

		bindingsitesinseq1 := oligos.CheckNonSpecificBinding(sequence, primer1.DNASequence)

		primer2.Nm = "primer2" + "_" + file

		bindingsitesinseq2 := oligos.CheckNonSpecificBinding(sequence, primer2.DNASequence)

		output = fmt.Sprintln(file, ",", "primer1: ", ",", primer1.Sequence(), ",", "melting temp: ", ",", primer1.MeltingTemp.ToString(), ",", "length: ", ",", primer1.Length, ",", "gc content: ", ",", primer1.GCContent, ",", "binds at", ",", bindingsitesinseq1, ",", "positions", ",", "primer2: ", ",", primer2.Sequence(), ",", "melting temp: ", ",", primer2.MeltingTemp.ToString(), ",", "length: ", ",", primer2.Length, ",", "gc content: ", ",", primer2.GCContent, ",", "binds at", ",", bindingsitesinseq2, ",", "positions", ",")
		alloutputs = append(alloutputs, output)
		allprimers = append(allprimers, primer1, primer2)
		allprimerstrings = append(allprimerstrings, primer1.Sequence(), primer2.Sequence())

		primerpairs = append(primerpairs, PrimerPair{primer1.Sequence(), primer2.Sequence()})

	}

	fmt.Println(alloutputs, allprimers)

	_output.AllOutputs = alloutputs

	if _input.ExportToFile {
		err = export.ExporttoTextFile("exported_primers.csv", _output.AllOutputs)

		if err != nil {
			panic(err.Error())
		}

	}

	_output.AllPrimers = allprimers
	_output.PrimerPairs = primerpairs

}

// Actions to perform after steps block to analyze data
func _PrimerDesign_sequence_your_neighboursAnalysis(_ctx context.Context, _input *PrimerDesign_sequence_your_neighboursInput, _output *PrimerDesign_sequence_your_neighboursOutput) {

}
func _PrimerDesign_sequence_your_neighboursValidation(_ctx context.Context, _input *PrimerDesign_sequence_your_neighboursInput, _output *PrimerDesign_sequence_your_neighboursOutput) {

	// check each sequence for binding to other sequences in folder:

	//var Start int
	//var End int
	var err error
	var output string
	var dirname string

	var files = make([]string, 0)

	//Search for files within current directory

	if _input.Dirname == "current" {
		dirname = "." + string(filepath.Separator)
	} else {
		dirname = _input.Dirname
	}

	d, err := os.Open(dirname)
	if err != nil {
		panic(err)
	}
	defer d.Close()

	allfiles, err := d.Readdir(-1)
	if err != nil {
		panic(err)
	}

	fmt.Println("Reading " + dirname)

	//Determine if file extension is ".gb"
	for _, file := range allfiles {
		if filepath.Ext(file.Name()) == ".gb" {
			files = append(files, file.Name())
		}

	}

	var nonspecificbinding = make([]string, 0)

	for _, file := range files {
		file = filepath.Join(dirname, file)
		sequence, _ := parser.GenbanktoAnnotatedSeq(file)

		for _, primer := range _output.AllPrimers {

			// only check other files
			if strings.Contains(primer.Nm, file) == false {

				bindingsites := oligos.CheckNonSpecificBinding(sequence, primer.DNASequence)

				// if binding found add to output file:

				if bindingsites > 0 {
					output = fmt.Sprintln(file, ",", "primer: ", ",", primer.Nm, ", ", primer.Sequence(), ",", "binds at", ",", bindingsites, ",", "positions")

					nonspecificbinding = append(nonspecificbinding, output)

				}
			}
		}

	}

	if _input.ExportToFile {
		err = export.ExporttoTextFile("exported_primers_bindingReport.csv", nonspecificbinding)

		if err != nil {
			panic(err.Error())
		}

	}
}

type PrimerPair struct {
	FWD string
	REV string
}

func _PrimerDesign_sequence_your_neighboursRun(_ctx context.Context, input *PrimerDesign_sequence_your_neighboursInput) *PrimerDesign_sequence_your_neighboursOutput {
	output := &PrimerDesign_sequence_your_neighboursOutput{}
	_PrimerDesign_sequence_your_neighboursSetup(_ctx, input)
	_PrimerDesign_sequence_your_neighboursSteps(_ctx, input, output)
	_PrimerDesign_sequence_your_neighboursAnalysis(_ctx, input, output)
	_PrimerDesign_sequence_your_neighboursValidation(_ctx, input, output)
	return output
}

func PrimerDesign_sequence_your_neighboursRunSteps(_ctx context.Context, input *PrimerDesign_sequence_your_neighboursInput) *PrimerDesign_sequence_your_neighboursSOutput {
	soutput := &PrimerDesign_sequence_your_neighboursSOutput{}
	output := _PrimerDesign_sequence_your_neighboursRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PrimerDesign_sequence_your_neighboursNew() interface{} {
	return &PrimerDesign_sequence_your_neighboursElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PrimerDesign_sequence_your_neighboursInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PrimerDesign_sequence_your_neighboursRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PrimerDesign_sequence_your_neighboursInput{},
			Out: &PrimerDesign_sequence_your_neighboursOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PrimerDesign_sequence_your_neighboursElement struct {
	inject.CheckedRunner
}

type PrimerDesign_sequence_your_neighboursInput struct {
	Dirname                                  string
	ExportToFile                             bool
	Maxgc                                    float64
	Maxlength                                int
	Maxtemp                                  wunit.Temperature
	Minlength                                int
	Mintemp                                  wunit.Temperature
	PermittednucleotideOverlapBetweenPrimers int
}

type PrimerDesign_sequence_your_neighboursOutput struct {
	AllOutputs  []string
	AllPrimers  []oligos.Primer
	PrimerPairs []PrimerPair
}

type PrimerDesign_sequence_your_neighboursSOutput struct {
	Data struct {
		AllOutputs  []string
		AllPrimers  []oligos.Primer
		PrimerPairs []PrimerPair
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PrimerDesign_sequence_your_neighbours",
		Constructor: PrimerDesign_sequence_your_neighboursNew,
		Desc: component.ComponentDesc{
			Desc: "This element will design outward facing primers for all .gb file sequences in a specified folder.\nDesign criteria such as maximum gc content, acceptable ranges of melting temperatures and primer length may be specified by the user.\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/DNA/PrimerDesign/PrimerDesign_sequence_your_neighbours.an",
			Params: []component.ParamDesc{
				{Name: "Dirname", Desc: "files     []string = []string{\"STAR_0023_VECTOR_BBSI.gb\", \"STAR_0023_VECTOR_BBSI+Grp7+Grp14+Grp3.gb\"}\n\n= \"current\" // this will check for all .gb files in the folder you select here\n", Kind: "Parameters"},
				{Name: "ExportToFile", Desc: "", Kind: "Parameters"},
				{Name: "Maxgc", Desc: "     = 0.6\n", Kind: "Parameters"},
				{Name: "Maxlength", Desc: "    = 25\n", Kind: "Parameters"},
				{Name: "Maxtemp", Desc: "     = wunit.NewTemperature(60, \"C\")\n", Kind: "Parameters"},
				{Name: "Minlength", Desc: "    = 20\n", Kind: "Parameters"},
				{Name: "Mintemp", Desc: "     = wunit.NewTemperature(55, \"C\")\n", Kind: "Parameters"},
				{Name: "PermittednucleotideOverlapBetweenPrimers", Desc: "number of nucleotides which primers can overlap by\n", Kind: "Parameters"},
				{Name: "AllOutputs", Desc: "", Kind: "Data"},
				{Name: "AllPrimers", Desc: "", Kind: "Data"},
				{Name: "PrimerPairs", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
