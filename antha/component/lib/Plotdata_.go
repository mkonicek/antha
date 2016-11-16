package lib

import (
	graph "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/plot"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

/*datarange*/
/*datarange*/

//	HeaderRange []string

// Data which is returned from this protocol, and data types

//	OutputData       []string

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _PlotdataRequirements() {

}

// Conditions to run on startup
func _PlotdataSetup(_ctx context.Context, _input *PlotdataInput) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _PlotdataSteps(_ctx context.Context, _input *PlotdataInput, _output *PlotdataOutput) {

	// now plot the graph

	// the data points

	plot := graph.Plot(_input.Xvalues, _input.Yvaluearray)

	graph.Export(plot, _input.Exportedfilename)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PlotdataAnalysis(_ctx context.Context, _input *PlotdataInput, _output *PlotdataOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PlotdataValidation(_ctx context.Context, _input *PlotdataInput, _output *PlotdataOutput) {

}
func _PlotdataRun(_ctx context.Context, input *PlotdataInput) *PlotdataOutput {
	output := &PlotdataOutput{}
	_PlotdataSetup(_ctx, input)
	_PlotdataSteps(_ctx, input, output)
	_PlotdataAnalysis(_ctx, input, output)
	_PlotdataValidation(_ctx, input, output)
	return output
}

func PlotdataRunSteps(_ctx context.Context, input *PlotdataInput) *PlotdataSOutput {
	soutput := &PlotdataSOutput{}
	output := _PlotdataRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PlotdataNew() interface{} {
	return &PlotdataElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PlotdataInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PlotdataRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PlotdataInput{},
			Out: &PlotdataOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PlotdataElement struct {
	inject.CheckedRunner
}

type PlotdataInput struct {
	Exportedfilename string
	Xvalues          []float64
	Yvaluearray      [][]float64
}

type PlotdataOutput struct {
}

type PlotdataSOutput struct {
	Data struct {
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Plotdata",
		Constructor: PlotdataNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Data/plotdata/Plotdata.an",
			Params: []component.ParamDesc{
				{Name: "Exportedfilename", Desc: "", Kind: "Parameters"},
				{Name: "Xvalues", Desc: "", Kind: "Parameters"},
				{Name: "Yvaluearray", Desc: "", Kind: "Parameters"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
