// Protocol for resuspending freeze dried DNA with a diluent
package lib

import

// we need to import the wtype package to use the LHComponent type
// the mixer package is required to use the Sample function
(
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Input parameters for this protocol (data)

func _ParseDNAPlateFileRequirements() {
}

func _ParseDNAPlateFileSetup(_ctx context.Context, _input *ParseDNAPlateFileInput) {
}

func _ParseDNAPlateFileSteps(_ctx context.Context, _input *ParseDNAPlateFileInput, _output *ParseDNAPlateFileOutput) {

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
	DNAPlate  *wtype.LHPlate
	Diluent   *wtype.LHComponent
	Platefile string
	Vendor    string
}

type ParseDNAPlateFileOutput struct {
	PartLocationsMap       map[string]string
	PartMassMap            map[string]wunit.Mass
	PartMolecularWeightMap map[string]float64
	Parts                  []string
	Platetype              string
	ResuspendedDNAMap      map[string]*wtype.LHComponent
}

type ParseDNAPlateFileSOutput struct {
	Data struct {
		PartLocationsMap       map[string]string
		PartMassMap            map[string]wunit.Mass
		PartMolecularWeightMap map[string]float64
		Parts                  []string
		Platetype              string
	}
	Outputs struct {
		ResuspendedDNAMap map[string]*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "ParseDNAPlateFile",
		Constructor: ParseDNAPlateFileNew,
		Desc: component.ComponentDesc{
			Desc: "Protocol for resuspending freeze dried DNA with a diluent\n",
			Path: "antha/component/an/Liquid_handling/ResuspendDNA/ParseDNAInputFile.an",
			Params: []component.ParamDesc{
				{Name: "DNAPlate", Desc: "", Kind: "Inputs"},
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "Platefile", Desc: "", Kind: "Parameters"},
				{Name: "Vendor", Desc: "", Kind: "Parameters"},
				{Name: "PartLocationsMap", Desc: "", Kind: "Data"},
				{Name: "PartMassMap", Desc: "", Kind: "Data"},
				{Name: "PartMolecularWeightMap", Desc: "", Kind: "Data"},
				{Name: "Parts", Desc: "", Kind: "Data"},
				{Name: "Platetype", Desc: "", Kind: "Data"},
				{Name: "ResuspendedDNAMap", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
