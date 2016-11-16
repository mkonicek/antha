// Example OD measurement protocol.
// Computes the OD and dry cell weight estimate from absorbance reading
// TODO: implement replicates from parameters
package lib

import (
	//"liquid handler"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/platereader"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

//"standard_labware"

// Input parameters for this protocol (data)

//= uL(100)
//= uL(0)
//Total_volume Volume//= ul (sample_volume+diluent_volume)
//Wavelength //= nm(600)
//Diluent_type //= (PBS)
//= (0.25)
//Replicate_count uint32 //= 1 // Note: 1 replicate means experiment is in duplicate, etc.
// calculate path length? - takes place under plate reader since this will only be necessary for plate reader protocols? labware?
// Data which is returned from this protocol, and data types
//= 0.0533
//WellCrosssectionalArea float64// should be calculated from plate and well type automatically

//Absorbance
//Absorbance
//(pathlength corrected)

//R_squared float32
//Control_absorbance [control_curve_points+1]float64//Absorbance
//Control_concentrations [control_curve_points+1]float64

// Physical Inputs to this protocol with types

//Culture

// Physical outputs from this protocol with types

// None

func _ODRequirements() {
	// sufficient sample volume available to sacrifice
}
func _ODSetup(_ctx context.Context, _input *ODInput) {
	/*control.Config(config.per_plate)
	var control_blank[total_volume]WaterSolution

	blank_absorbance = platereader.Read(ODplate,control_blank, wavelength)*/
}
func _ODSteps(_ctx context.Context, _input *ODInput, _output *ODOutput) {

	var product *wtype.LHComponent //WaterSolution
	var counter int
	if _input.MaxDilutions == 0 {
		_input.MaxDilutions = 10
	}

	for {
		product = execute.MixInto(_ctx, _input.ODplate, "", mixer.Sample(_input.Sampletotest, _input.Sample_volume), mixer.Sample(_input.Diluent, _input.Diluent_volume))
		/*Is it necessary to include platetype in Read function?
		or is the info on volume, opacity, pathlength etc implied in LHComponent?*/
		_output.Sample_absorbance = platereader.ReadAbsorbance(_input.ODplate, product, _input.Wlength)

		if _output.Sample_absorbance.Reading < 1 || counter > _input.MaxDilutions {
			break
		}
		_input.Diluent_volume.Mvalue += 1 //diluent_volume = diluent_volume + 1
		counter++
	}
} // serial dilution or could write element for finding optimum dilution or search historical data
func _ODAnalysis(_ctx context.Context, _input *ODInput, _output *ODOutput) {
	// Need to substract blank from measurement; normalise to path length of 1cm for OD value; apply conversion factor to estimate dry cell weight

	_output.Blankcorrected_absorbance = platereader.Blankcorrect(_output.Sample_absorbance, _input.Blank_absorbance)
	volumetopathlengthconversionfactor := wunit.NewLength(_input.Heightof100ulinm, "m")                               //WellCrosssectionalArea
	_output.OD = platereader.PathlengthCorrect(volumetopathlengthconversionfactor, _output.Blankcorrected_absorbance) // 0.0533 could be written as function of labware and liquid volume (or measureed height)
	_output.Estimateddrycellweight_conc = wunit.NewConcentration(_output.OD.Reading*_input.ODtoDCWconversionfactor, "g/L")
}
func _ODValidation(_ctx context.Context, _input *ODInput, _output *ODOutput) { /*
		if Sample_absorbance > 1 {
		panic("Sample likely needs further dilution")
		}
		if Sample_absorbance < 0.1 {
		warn("Low OD, sample likely needs increased volume")
		}
		}*/
	// TODO: add test of replicate variance
}
func _ODRun(_ctx context.Context, input *ODInput) *ODOutput {
	output := &ODOutput{}
	_ODSetup(_ctx, input)
	_ODSteps(_ctx, input, output)
	_ODAnalysis(_ctx, input, output)
	_ODValidation(_ctx, input, output)
	return output
}

func ODRunSteps(_ctx context.Context, input *ODInput) *ODSOutput {
	soutput := &ODSOutput{}
	output := _ODRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ODNew() interface{} {
	return &ODElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ODInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ODRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ODInput{},
			Out: &ODOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type ODElement struct {
	inject.CheckedRunner
}

type ODInput struct {
	Blank_absorbance        wtype.Absorbance
	Diluent                 *wtype.LHComponent
	Diluent_volume          wunit.Volume
	Heightof100ulinm        float64
	MaxDilutions            int
	ODplate                 *wtype.LHPlate
	ODtoDCWconversionfactor float64
	Sample_volume           wunit.Volume
	Sampletotest            *wtype.LHComponent
	Wlength                 float64
}

type ODOutput struct {
	Blankcorrected_absorbance   wtype.Absorbance
	Estimateddrycellweight_conc wunit.Concentration
	OD                          wtype.Absorbance
	Sample_absorbance           wtype.Absorbance
}

type ODSOutput struct {
	Data struct {
		Blankcorrected_absorbance   wtype.Absorbance
		Estimateddrycellweight_conc wunit.Concentration
		OD                          wtype.Absorbance
		Sample_absorbance           wtype.Absorbance
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "OD",
		Constructor: ODNew,
		Desc: component.ComponentDesc{
			Desc: "Example OD measurement protocol.\nComputes the OD and dry cell weight estimate from absorbance reading\nTODO: implement replicates from parameters\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Liquid_handling/OD/OD.an",
			Params: []component.ParamDesc{
				{Name: "Blank_absorbance", Desc: "WellCrosssectionalArea float64// should be calculated from plate and well type automatically\n", Kind: "Parameters"},
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "Diluent_volume", Desc: "= uL(0)\n", Kind: "Parameters"},
				{Name: "Heightof100ulinm", Desc: "Replicate_count uint32 //= 1 // Note: 1 replicate means experiment is in duplicate, etc.\ncalculate path length? - takes place under plate reader since this will only be necessary for plate reader protocols? labware?\nData which is returned from this protocol, and data types\n\n= 0.0533\n", Kind: "Parameters"},
				{Name: "MaxDilutions", Desc: "", Kind: "Parameters"},
				{Name: "ODplate", Desc: "", Kind: "Inputs"},
				{Name: "ODtoDCWconversionfactor", Desc: "Diluent_type //= (PBS)\n\n= (0.25)\n", Kind: "Parameters"},
				{Name: "Sample_volume", Desc: "= uL(100)\n", Kind: "Parameters"},
				{Name: "Sampletotest", Desc: "Culture\n", Kind: "Inputs"},
				{Name: "Wlength", Desc: "Total_volume Volume//= ul (sample_volume+diluent_volume)\n\nWavelength //= nm(600)\n", Kind: "Parameters"},
				{Name: "Blankcorrected_absorbance", Desc: "Absorbance\n", Kind: "Data"},
				{Name: "Estimateddrycellweight_conc", Desc: "", Kind: "Data"},
				{Name: "OD", Desc: "(pathlength corrected)\n", Kind: "Data"},
				{Name: "Sample_absorbance", Desc: "Absorbance\n", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
