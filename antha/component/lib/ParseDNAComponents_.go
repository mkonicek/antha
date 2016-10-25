// Example protocol demonstrating the use of the Sample function
package lib

import

// we need to import the wtype package to use the LHComponent type
// the mixer package is required to use the Sample function
(
	"fmt"
	//"github.com/antha-lang/antha/antha/anthalib/wtype"
	//"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/doe"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

//ComponentFile string

//Outputs from this element

func _ParseDNAComponentsRequirements() {
}

func _ParseDNAComponentsSetup(_ctx context.Context, _input *ParseDNAComponentsInput) {
}

func _ParseDNAComponentsSteps(_ctx context.Context, _input *ParseDNAComponentsInput, _output *ParseDNAComponentsOutput) {
	//create a string that will store the header names the element finds in the input file, this is for an internal check to make sure all 4 columns of information is found.
	headersfound := make([]string, 0)

	//create strings that will be populated with the values from each column in the input file
	_output.Partnames = make([]string, 0)
	_output.PartConcs = make([]float64, 0)
	_output.Partplusvectorlengths = make([]int, 0)
	_output.PartLengths = make([]int, 0)

	dnaparts, err := doe.RunsFromDesignPreResponses(_input.SequenceInfoFile, []string{_input.PartplusvectorlengthsHeader}, _input.SequenceInfoFileformat)

	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	// code for parsing the data from the xl file into the strings, this searches the file in direction i followed by j
	for i, partinfo := range dnaparts {

		for j := range partinfo.Setpoints {

			//First creates an array of part names
			if partinfo.Factordescriptors[j] == _input.NameHeader {

				if name, found := partinfo.Setpoints[j].(string); found {
					_output.Partnames = append(_output.Partnames, name)
				} else {
					execute.Errorf(_ctx, fmt.Sprint("wrong type", partinfo.Factordescriptors[j], partinfo.Setpoints[j]))
				}

				if i == 0 {
					headersfound = append(headersfound, _input.NameHeader)
				}

			}

			//second creats an array of plasmid concentrations
			if partinfo.Factordescriptors[j] == _input.ConcHeader {

				if conc, found := partinfo.Setpoints[j].(float64); found {
					_output.PartConcs = append(_output.PartConcs, conc)
				} else {
					execute.Errorf(_ctx, fmt.Sprint("wrong type", partinfo.Factordescriptors[j], partinfo.Setpoints[j]))
				}
				if i == 0 {
					headersfound = append(headersfound, _input.ConcHeader)
				}
			}

			//third creates an array of part lengths in bp
			if partinfo.Factordescriptors[j] == _input.PartLengthHeader {

				if partlength, found := partinfo.Setpoints[j].(int); found {
					_output.PartLengths = append(_output.PartLengths, partlength)
				} else if partlength, found := partinfo.Setpoints[j].(float64); found {
					_output.PartLengths = append(_output.PartLengths, int(partlength))
				} else {
					execute.Errorf(_ctx, fmt.Sprint("wrong type", partinfo.Factordescriptors[j], partinfo.Setpoints[j]))
				}
				if i == 0 {
					headersfound = append(headersfound, _input.PartLengthHeader)
				}
			}

			//forth creates an array of total plasmid size (part + vector) in bp
			if partinfo.Factordescriptors[j] == _input.PartplusvectorlengthsHeader {

				if partplusplasmid, found := partinfo.Setpoints[j].(int); found {
					_output.Partplusvectorlengths = append(_output.Partplusvectorlengths, partplusplasmid)
				} else if partplusplasmid, found := partinfo.Setpoints[j].(float64); found {
					_output.Partplusvectorlengths = append(_output.Partplusvectorlengths, int(partplusplasmid))
				} else {
					execute.Errorf(_ctx, fmt.Sprint("wrong type", partinfo.Factordescriptors[j], partinfo.Setpoints[j]))
				}
				if i == 0 {
					headersfound = append(headersfound, _input.PartplusvectorlengthsHeader)
				}
			}

		}

		//internal check if there are not 4 headers (as we know there should be 4) return an error telling us which ones were found and which were not
		if len(headersfound) != 4 {
			execute.Errorf(_ctx, fmt.Sprint("Only found these headers in input file: ", headersfound))
		}

	}

}

func _ParseDNAComponentsAnalysis(_ctx context.Context, _input *ParseDNAComponentsInput, _output *ParseDNAComponentsOutput) {
}

func _ParseDNAComponentsValidation(_ctx context.Context, _input *ParseDNAComponentsInput, _output *ParseDNAComponentsOutput) {
}
func _ParseDNAComponentsRun(_ctx context.Context, input *ParseDNAComponentsInput) *ParseDNAComponentsOutput {
	output := &ParseDNAComponentsOutput{}
	_ParseDNAComponentsSetup(_ctx, input)
	_ParseDNAComponentsSteps(_ctx, input, output)
	_ParseDNAComponentsAnalysis(_ctx, input, output)
	_ParseDNAComponentsValidation(_ctx, input, output)
	return output
}

func ParseDNAComponentsRunSteps(_ctx context.Context, input *ParseDNAComponentsInput) *ParseDNAComponentsSOutput {
	soutput := &ParseDNAComponentsSOutput{}
	output := _ParseDNAComponentsRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ParseDNAComponentsNew() interface{} {
	return &ParseDNAComponentsElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ParseDNAComponentsInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ParseDNAComponentsRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ParseDNAComponentsInput{},
			Out: &ParseDNAComponentsOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type ParseDNAComponentsElement struct {
	inject.CheckedRunner
}

type ParseDNAComponentsInput struct {
	ConcHeader                  string
	NameHeader                  string
	PartLengthHeader            string
	PartplusvectorlengthsHeader string
	SequenceInfoFile            string
	SequenceInfoFileformat      string
}

type ParseDNAComponentsOutput struct {
	PartConcs             []float64
	PartLengths           []int
	Partnames             []string
	Partplusvectorlengths []int
	Status                string
}

type ParseDNAComponentsSOutput struct {
	Data struct {
		PartConcs             []float64
		PartLengths           []int
		Partnames             []string
		Partplusvectorlengths []int
		Status                string
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "ParseDNAComponents",
		Constructor: ParseDNAComponentsNew,
		Desc: component.ComponentDesc{
			Desc: "Example protocol demonstrating the use of the Sample function\n",
			Path: "antha/component/an/Liquid_handling/PooledLibrary/playground/ParseDNAComponents.an",
			Params: []component.ParamDesc{
				{Name: "ConcHeader", Desc: "", Kind: "Parameters"},
				{Name: "NameHeader", Desc: "", Kind: "Parameters"},
				{Name: "PartLengthHeader", Desc: "", Kind: "Parameters"},
				{Name: "PartplusvectorlengthsHeader", Desc: "", Kind: "Parameters"},
				{Name: "SequenceInfoFile", Desc: "ComponentFile string\n", Kind: "Parameters"},
				{Name: "SequenceInfoFileformat", Desc: "", Kind: "Parameters"},
				{Name: "PartConcs", Desc: "", Kind: "Data"},
				{Name: "PartLengths", Desc: "", Kind: "Data"},
				{Name: "Partnames", Desc: "", Kind: "Data"},
				{Name: "Partplusvectorlengths", Desc: "", Kind: "Data"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
