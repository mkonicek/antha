// Protocol for resuspending freeze dried DNA with a diluent
package lib

import

// we need to import the wtype package to use the LHComponent type
// the mixer package is required to use the Sample function
(
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/doe"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/factory"
	"golang.org/x/net/context"
	"strings"
)

// Input parameters for this protocol (data)

// using the fwd primer name, returns the rev primer name

func _ParseDNAPlateFileRequirements() {
}

func _ParseDNAPlateFileSetup(_ctx context.Context, _input *ParseDNAPlateFileInput) {
}

func _ParseDNAPlateFileSteps(_ctx context.Context, _input *ParseDNAPlateFileInput, _output *ParseDNAPlateFileOutput) {

	var ReplaceMap = map[string]string{
		"01": "1",
		"02": "2",
		"03": "3",
		"04": "4",
		"05": "5",
		"06": "6",
		"07": "7",
		"08": "8",
		"09": "9",
		"_":  "",
	}

	// headers
	NameHeader := "Seq Name"
	MassHeader := "Yield_ug"
	MWHeader := "MW"
	WellHeader := "Customer Well"
	PlateNameHeader := "Plate"

	// initialise maps
	_output.PartMassMap = make(map[string]wunit.Mass)
	_output.PartMolecularWeightMap = make(map[string]float64)
	_output.PartLocationsMap = make(map[string]string)
	_output.PartPlateMap = make(map[string]string)
	headersfound := make([]string, 0)
	_output.FwdOligotoRevOligo = make(map[string]string)
	_output.PartsList = make(map[string]*wtype.LHComponent)

	dnaparts, err := doe.RunsFromDesignPreResponses(_input.SequenceInfoFile, []string{"Length", "MW", "Tm", "Yield"}, _input.SequenceInfoFileformat)

	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	// code for parsing the data from the xl file into the strings, this searches the file in direction i followed by j
	for i, partinfo := range dnaparts {

		var partname string
		var partmass float64
		var partwell string
		var partmw float64
		var platename string

		for j := range partinfo.Setpoints {

			//First creates an array of part names
			if partinfo.Factordescriptors[j] == NameHeader {

				if name, found := partinfo.Setpoints[j].(string); found {

					if name == "" || name == "BLANK" {
						fmt.Print("Skipping ", name)
					} else {
						partname = name
						_output.Partnames = append(_output.Partnames, name)
					}
				} else {
					execute.Errorf(_ctx, fmt.Sprint("wrong type", partinfo.Factordescriptors[j], partinfo.Setpoints[j]))
				}

				if i == 0 {
					headersfound = append(headersfound, NameHeader)
				}

			}

			//second create an array of plasmid masses
			if partinfo.Factordescriptors[j] == MassHeader {

				if mass, found := partinfo.Setpoints[j].(float64); found {
					partmass = mass
				} else if mass, found := partinfo.Setpoints[j].(string); found {
					if mass == "" {
						partmass = 0.0 // mw
						// empty so skip
					}
				} else {
					execute.Errorf(_ctx, fmt.Sprint("wrong type", partinfo.Factordescriptors[j], partinfo.Setpoints[j]))
				}
				if i == 0 {
					headersfound = append(headersfound, MassHeader)
				}
			}

			//third creates an array of part lengths in bp
			if partinfo.Factordescriptors[j] == MWHeader {

				if mw, found := partinfo.Setpoints[j].(int); found {
					partmw = float64(mw)
				} else if mw, found := partinfo.Setpoints[j].(float64); found {
					partmw = mw
				} else if mw, found := partinfo.Setpoints[j].(string); found {
					if mw == "" {
						partmw = 0.0 // mw
						// empty so skip
					}
				} else {
					execute.Errorf(_ctx, fmt.Sprint("wrong type", partinfo.Factordescriptors[j], partinfo.Setpoints[j]))
				}
				if i == 0 {
					headersfound = append(headersfound, MWHeader)
				}
			}

			//third creates an array of part lengths in bp
			if partinfo.Factordescriptors[j] == WellHeader {

				if well, found := partinfo.Setpoints[j].(string); found {

					for key, value := range ReplaceMap {

						if strings.Contains(well, key) {
							well = strings.Replace(well, key, value, 1)

							break
						}
					}
					partwell = well
				} else {
					execute.Errorf(_ctx, fmt.Sprint("wrong type", partinfo.Factordescriptors[j], partinfo.Setpoints[j]))
				}
				if i == 0 {
					headersfound = append(headersfound, WellHeader)
				}
			}

			//third creates an array of part lengths in bp
			if partinfo.Factordescriptors[j] == PlateNameHeader {

				if plate, found := partinfo.Setpoints[j].(string); found {

					platename = plate
				} else {
					execute.Errorf(_ctx, fmt.Sprint("wrong type", partinfo.Factordescriptors[j], partinfo.Setpoints[j]))
				}
				if i == 0 {
					headersfound = append(headersfound, PlateNameHeader)
				}
			}

			//internal check if there are not 4 headers (as we know there should be 4) return an error telling us which ones were found and which were not
			/*
				if len(headersfound)!= 4 {
					Errorf(fmt.Sprint("Only found these headers in input file: ", headersfound))
				}
			*/
		}

		if partname == "" || partname == "BLANK" {
			fmt.Print("Skipping ", partname)
		} else {
			_output.PartLocationsMap[partname] = partwell
			_output.PartMolecularWeightMap[partname] = partmw
			_output.PartMassMap[partname] = wunit.NewMass(partmass, "ug")
			_output.PartPlateMap[partname] = platename

			part := factory.GetComponentByType("dna_part")
			part.CName = partname

			_output.PartsList[partname] = part
		}

	}

	for _, partname := range _output.Partnames {

		if !strings.Contains(partname, "_Revers") {
			for _, partname2 := range _output.Partnames {
				if strings.Contains(partname2, "_Revers") && strings.Contains(partname2, partname) {
					_output.FwdOligotoRevOligo[partname] = partname2
					break
				}
			}
		}
	}
	_output.OligoPairs = len(_output.FwdOligotoRevOligo)

	_output.HeadersFound = headersfound
}
func _ParseDNAPlateFileAnalysis(_ctx context.Context, _input *ParseDNAPlateFileInput, _output *ParseDNAPlateFileOutput) {
}

func _ParseDNAPlateFileValidation(_ctx context.Context, _input *ParseDNAPlateFileInput, _output *ParseDNAPlateFileOutput) {
}
func _ParseDNAPlateFileRun(_ctx context.Context, input *ParseDNAPlateFileInput) *ParseDNAPlateFileOutput {
	output := &ParseDNAPlateFileOutput{}
	_ParseDNAPlateFileSetup(_ctx, input)
	_ParseDNAPlateFileSteps(_ctx, input, output)
	_ParseDNAPlateFileAnalysis(_ctx, input, output)
	_ParseDNAPlateFileValidation(_ctx, input, output)
	return output
}

func ParseDNAPlateFileRunSteps(_ctx context.Context, input *ParseDNAPlateFileInput) *ParseDNAPlateFileSOutput {
	soutput := &ParseDNAPlateFileSOutput{}
	output := _ParseDNAPlateFileRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ParseDNAPlateFileNew() interface{} {
	return &ParseDNAPlateFileElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ParseDNAPlateFileInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ParseDNAPlateFileRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ParseDNAPlateFileInput{},
			Out: &ParseDNAPlateFileOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type ParseDNAPlateFileElement struct {
	inject.CheckedRunner
}

type ParseDNAPlateFileInput struct {
	DNAPlate               *wtype.LHPlate
	SequenceInfoFile       string
	SequenceInfoFileformat string
}

type ParseDNAPlateFileOutput struct {
	FwdOligotoRevOligo     map[string]string
	HeadersFound           []string
	OligoPairs             int
	PartLocationsMap       map[string]string
	PartMassMap            map[string]wunit.Mass
	PartMolecularWeightMap map[string]float64
	PartPlateMap           map[string]string
	Partnames              []string
	PartsList              map[string]*wtype.LHComponent
	Platetype              string
}

type ParseDNAPlateFileSOutput struct {
	Data struct {
		FwdOligotoRevOligo     map[string]string
		HeadersFound           []string
		OligoPairs             int
		PartLocationsMap       map[string]string
		PartMassMap            map[string]wunit.Mass
		PartMolecularWeightMap map[string]float64
		PartPlateMap           map[string]string
		Partnames              []string
		Platetype              string
	}
	Outputs struct {
		PartsList map[string]*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "ParseDNAPlateFile",
		Constructor: ParseDNAPlateFileNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol for resuspending freeze dried DNA with a diluent\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/ResuspendDNA/ParseDNAInputFile.an",
			Params: []component.ParamDesc{
				{Name: "DNAPlate", Desc: "", Kind: "Inputs"},
				{Name: "SequenceInfoFile", Desc: "", Kind: "Parameters"},
				{Name: "SequenceInfoFileformat", Desc: "", Kind: "Parameters"},
				{Name: "FwdOligotoRevOligo", Desc: "using the fwd primer name, returns the rev primer name\n", Kind: "Data"},
				{Name: "HeadersFound", Desc: "", Kind: "Data"},
				{Name: "OligoPairs", Desc: "", Kind: "Data"},
				{Name: "PartLocationsMap", Desc: "", Kind: "Data"},
				{Name: "PartMassMap", Desc: "", Kind: "Data"},
				{Name: "PartMolecularWeightMap", Desc: "", Kind: "Data"},
				{Name: "PartPlateMap", Desc: "", Kind: "Data"},
				{Name: "Partnames", Desc: "", Kind: "Data"},
				{Name: "PartsList", Desc: "", Kind: "Outputs"},
				{Name: "Platetype", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
