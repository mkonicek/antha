package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
	"strconv"
)

//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Inventory"

// Input parameters for this protocol (data)

//Volume //Mass // Should be Mass

//  +/- x  e.g. 7.0 +/- 0.2

//LiqComponentkeys	[]string
//Solidcomponentkeys	[]string // name or barcode id
//Acidkey string
//Basekey string

// Physical Inputs to this protocol with types

// should be new type or field indicating solid and mass
/*Acid					*wtype.LHComponent
Base 					*wtype.LHComponent
*/

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

// Input Requirement specification
func _MakeMediaRequirements() {

}

// Conditions to run on startup
func _MakeMediaSetup(_ctx context.Context, _input *MakeMediaInput) {}

// The core process for this protocol, with the steps to be performed
// for every input
func _MakeMediaSteps(_ctx context.Context, _input *MakeMediaInput, _output *MakeMediaOutput) {
	recipestring := make([]string, 0)
	var step string
	stepcounter := 1 // counting from 1 is what makes us human
	liquids := make([]*wtype.LHComponent, 0)
	step = text.Print("Recipe for: ", _input.Name)
	recipestring = append(recipestring, step)

	step = text.Print("Step"+strconv.Itoa(stepcounter)+": ", "Prepare clean "+_input.Vessel.Type)
	recipestring = append(recipestring, step)
	stepcounter++

	for _, liq := range _input.LiqComponents {

		volrequired, err := wunit.VolumeForTargetConcentration(_input.LiqComponentFinalConcentrations[liq.CName], _input.LiqComponentStockConcentrations[liq.CName], _input.TotalVolume)
		if err != nil {
			execute.Errorf(_ctx, fmt.Sprint("Liquid: ", liq.CName, " ", err.Error()))
		}
		liqsamp := mixer.Sample(liq, volrequired)
		liquids = append(liquids, liqsamp)
		step = text.Print("Step"+strconv.Itoa(stepcounter)+": ", "add "+volrequired.ToString()+" of "+_input.LiqComponentStockConcentrations[liq.CName].ToString()+" "+liq.CName)
		recipestring = append(recipestring, step)
		stepcounter++
	}

	//solids := make([]*wtype.LHComponent,0)

	for _, sol := range _input.SolidComponents {

		targetMass, err := wunit.MassForTargetConcentration(_input.SolidComponentFinalConcentrations[sol.CName], _input.TotalVolume)
		if err != nil {
			execute.Errorf(_ctx, fmt.Sprint("Solid: ", sol.CName, " ", err.Error()))
		}
		solsamp := mixer.SampleMass(sol, targetMass, _input.SolidComponentDensities[sol.CName])
		liquids = append(liquids, solsamp)
		step = text.Print("Step"+strconv.Itoa(stepcounter)+": ", "add "+targetMass.ToString()+" of "+sol.CName)
		recipestring = append(recipestring, step)
		stepcounter++
		//stepcounter = stepcounter + k
	}

	watersample := mixer.SampleForTotalVolume(_input.Diluent, _input.TotalVolume)
	liquids = append(liquids, watersample)
	step = text.Print("Step"+strconv.Itoa(stepcounter)+": ", "add up to "+_input.TotalVolume.ToString()+" of "+_input.Diluent.CName)
	recipestring = append(recipestring, step)

	stepcounter++

	// Add pH handling functions and driver calls etc...

	description := fmt.Sprint("adjust pH to ", _input.PH_setPoint, " +/-", _input.PH_tolerance, " for temp ", _input.PH_setPointTemp.ToString(), "C")
	step = text.Print("Step"+strconv.Itoa(stepcounter)+": ", description)
	recipestring = append(recipestring, step)
	stepcounter++

	if _input.Sterilise {
		description := fmt.Sprint("Now Sterilise by ", _input.SterilisationMethod)
		step = text.Print("Step"+strconv.Itoa(stepcounter)+": ", description)
		recipestring = append(recipestring, step)
		stepcounter++
	}

	for _, extrastep := range _input.AdditionalInstructions {
		description := fmt.Sprint("Now, ", extrastep)
		step = text.Print("Step"+strconv.Itoa(stepcounter)+": ", description)
		recipestring = append(recipestring, step)
		stepcounter++
	}

	/*
		prepH := MixInto(Vessel,liquids...)

		pHactual := prepH.Measure("pH")

		step = text.Print("pH measured = ", pHactual)
		recipestring = append(recipestring,step)

		//pHactual = wutil.Roundto(pHactual,PH_tolerance)

		pHmax := PH_setpoint + PH_tolerance
		pHmin := PH_setpoint - PH_tolerance

		if pHactual < pHmax || pHactual < pHmin {
			// basically just a series of sample, stir, wait and recheck pH
		Media, newph, componentadded = prepH.AdjustpH(PH_setPoint, pHactual, PH_setPointTemp,Acid,Base)

		step = text.Print("Adjusted pH = ", newpH)
		recipestring = append(recipestring,step)

		step = text.Print("Component added = ", componentadded.Vol + componentadded.Vunit + " of " + componentadded.Conc + componentadded.Cunit + " " + componentadded.CName + )
		recipestring = append(recipestring,step)
		}
	*/
	_output.Media = execute.MixInto(_ctx, _input.Vessel, "", liquids...)
	_output.RecipeSteps = recipestring

	fmt.Println(recipestring)
	_output.Status = fmt.Sprintln(recipestring)

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _MakeMediaAnalysis(_ctx context.Context, _input *MakeMediaInput, _output *MakeMediaOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _MakeMediaValidation(_ctx context.Context, _input *MakeMediaInput, _output *MakeMediaOutput) {
}
func _MakeMediaRun(_ctx context.Context, input *MakeMediaInput) *MakeMediaOutput {
	output := &MakeMediaOutput{}
	_MakeMediaSetup(_ctx, input)
	_MakeMediaSteps(_ctx, input, output)
	_MakeMediaAnalysis(_ctx, input, output)
	_MakeMediaValidation(_ctx, input, output)
	return output
}

func MakeMediaRunSteps(_ctx context.Context, input *MakeMediaInput) *MakeMediaSOutput {
	soutput := &MakeMediaSOutput{}
	output := _MakeMediaRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func MakeMediaNew() interface{} {
	return &MakeMediaElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &MakeMediaInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _MakeMediaRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &MakeMediaInput{},
			Out: &MakeMediaOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type MakeMediaElement struct {
	inject.CheckedRunner
}

type MakeMediaInput struct {
	AdditionalInstructions            []string
	Diluent                           *wtype.LHComponent
	LiqComponentFinalConcentrations   map[string]wunit.Concentration
	LiqComponentStockConcentrations   map[string]wunit.Concentration
	LiqComponents                     []*wtype.LHComponent
	Name                              string
	PH_setPoint                       float64
	PH_setPointTemp                   wunit.Temperature
	PH_tolerance                      float64
	SolidComponentDensities           map[string]wunit.Density
	SolidComponentFinalConcentrations map[string]wunit.Concentration
	SolidComponents                   []*wtype.LHComponent
	SterilisationMethod               string
	Sterilise                         bool
	TotalVolume                       wunit.Volume
	Vessel                            *wtype.LHPlate
}

type MakeMediaOutput struct {
	Media       *wtype.LHComponent
	RecipeSteps []string
	Status      string
}

type MakeMediaSOutput struct {
	Data struct {
		RecipeSteps []string
		Status      string
	}
	Outputs struct {
		Media *wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "MakeMedia",
		Constructor: MakeMediaNew,
		Desc: component.ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/MakeMedia/MakeMedia.an",
			Params: []component.ParamDesc{
				{Name: "AdditionalInstructions", Desc: "", Kind: "Inputs"},
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "LiqComponentFinalConcentrations", Desc: "", Kind: "Parameters"},
				{Name: "LiqComponentStockConcentrations", Desc: "", Kind: "Parameters"},
				{Name: "LiqComponents", Desc: "", Kind: "Inputs"},
				{Name: "Name", Desc: "", Kind: "Parameters"},
				{Name: "PH_setPoint", Desc: "", Kind: "Parameters"},
				{Name: "PH_setPointTemp", Desc: "", Kind: "Parameters"},
				{Name: "PH_tolerance", Desc: " +/- x  e.g. 7.0 +/- 0.2\n", Kind: "Parameters"},
				{Name: "SolidComponentDensities", Desc: "", Kind: "Parameters"},
				{Name: "SolidComponentFinalConcentrations", Desc: "Volume //Mass // Should be Mass\n", Kind: "Parameters"},
				{Name: "SolidComponents", Desc: "should be new type or field indicating solid and mass\n", Kind: "Inputs"},
				{Name: "SterilisationMethod", Desc: "", Kind: "Inputs"},
				{Name: "Sterilise", Desc: "Acid\t\t\t\t\t*wtype.LHComponent\n\tBase \t\t\t\t\t*wtype.LHComponent\n", Kind: "Inputs"},
				{Name: "TotalVolume", Desc: "", Kind: "Parameters"},
				{Name: "Vessel", Desc: "", Kind: "Inputs"},
				{Name: "Media", Desc: "", Kind: "Outputs"},
				{Name: "RecipeSteps", Desc: "", Kind: "Data"},
				{Name: "Status", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}

/*
type Mole struct {
	number float64
}*/
