package lib

import (
	"fmt"
	//"math/rand"
	//"github.com/montanaflynn/stats"
	graph "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/plot"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/spreadsheet"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

//                                                                         = "plotinumdata.xlsx"
//                                                                        = 0
/*datarange*/ //  = []string{"a4", "a16"}                                                           // row in A1 format i.e string{A,E} would indicate all data between those points
/*datarange*/ //= [][]string{[]string{"b4", "b16"}, []string{"c4", "c16"}, []string{"d4", "d16"}} // column in A1 format i.e string{1,12} would indicate all data between those points
//= "Excelfile.jpg"
//	HeaderRange []string // if Bycolumn == true, format would be e.g. string{A1,E1} else e.g. string{A1,A12}

// Data which is returned from this protocol, and data types

//	OutputData       []string

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _Plotdata_spreadsheetRequirements() {

}

// Conditions to run on startup
func _Plotdata_spreadsheetSetup(_ctx context.Context, _input *Plotdata_spreadsheetInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _Plotdata_spreadsheetSteps(_ctx context.Context, _input *Plotdata_spreadsheetInput, _output *Plotdata_spreadsheetOutput) {

	// Get some data.

	file, err := spreadsheet.OpenFile(_input.Filename)

	sheet := file.Sheets[_input.Sheet]

	Xdatarange, err := spreadsheet.ConvertMinMaxtoArray(_input.Xminmax)
	if err != nil {
		fmt.Println(_input.Xminmax, Xdatarange)
		panic(err)
	}
	//	fmt.Println(Xdatarange)

	Ydatarangearray := make([][]string, 0)
	for i, Yminmax := range _input.Yminmaxarray {
		Ydatarange, err := spreadsheet.ConvertMinMaxtoArray(Yminmax)
		if err != nil {
			panic(err)
		}
		if len(Xdatarange) != len(Ydatarange) {
			panicmessage := fmt.Errorf("for index", i, "of array", "len(Xdatarange) != len(Ydatarange)")
			panic(panicmessage.Error())
		}
		Ydatarangearray = append(Ydatarangearray, Ydatarange)
		//	fmt.Println(Ydatarange)
	}

	// now plot the graph

	// the data points

	graph.PlotfromMinMaxpairs(sheet, _input.Xminmax, _input.Yminmaxarray, _input.Exportedfilename)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _Plotdata_spreadsheetAnalysis(_ctx context.Context, _input *Plotdata_spreadsheetInput, _output *Plotdata_spreadsheetOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _Plotdata_spreadsheetValidation(_ctx context.Context, _input *Plotdata_spreadsheetInput, _output *Plotdata_spreadsheetOutput) {

}
func _Plotdata_spreadsheetRun(_ctx context.Context, input *Plotdata_spreadsheetInput) *Plotdata_spreadsheetOutput {
	output := &Plotdata_spreadsheetOutput{}
	_Plotdata_spreadsheetSetup(_ctx, input)
	_Plotdata_spreadsheetSteps(_ctx, input, output)
	_Plotdata_spreadsheetAnalysis(_ctx, input, output)
	_Plotdata_spreadsheetValidation(_ctx, input, output)
	return output
}

func Plotdata_spreadsheetRunSteps(_ctx context.Context, input *Plotdata_spreadsheetInput) *Plotdata_spreadsheetSOutput {
	soutput := &Plotdata_spreadsheetSOutput{}
	output := _Plotdata_spreadsheetRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Plotdata_spreadsheetNew() interface{} {
	return &Plotdata_spreadsheetElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Plotdata_spreadsheetInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Plotdata_spreadsheetRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Plotdata_spreadsheetInput{},
			Out: &Plotdata_spreadsheetOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Plotdata_spreadsheetElement struct {
	inject.CheckedRunner
}

type Plotdata_spreadsheetInput struct {
	Exportedfilename string
	Filename         string
	Sheet            int
	Xminmax          []string
	Yminmaxarray     [][]string
}

type Plotdata_spreadsheetOutput struct {
}

type Plotdata_spreadsheetSOutput struct {
	Data struct {
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Plotdata_spreadsheet",
		Constructor: Plotdata_spreadsheetNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/plotdata/Plotdata_fromxlsx.an",
			Params: []component.ParamDesc{
				{Name: "Exportedfilename", Desc: "= \"Excelfile.jpg\"\n", Kind: "Parameters"},
				{Name: "Filename", Desc: "                                                                        = \"plotinumdata.xlsx\"\n", Kind: "Parameters"},
				{Name: "Sheet", Desc: "                                                                       = 0\n", Kind: "Parameters"},
				{Name: "Xminmax", Desc: " = []string{\"a4\", \"a16\"}                                                           // row in A1 format i.e string{A,E} would indicate all data between those points\n", Kind: "Parameters"},
				{Name: "Yminmaxarray", Desc: "= [][]string{[]string{\"b4\", \"b16\"}, []string{\"c4\", \"c16\"}, []string{\"d4\", \"d16\"}} // column in A1 format i.e string{1,12} would indicate all data between those points\n", Kind: "Parameters"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
