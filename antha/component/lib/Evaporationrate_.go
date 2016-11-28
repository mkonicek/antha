/* Evaporation calculator based on
http://www.engineeringtoolbox.com/evaporation-water-surface-d_690.html

This engineering function may need to be improved to account for vapour pressure and surface tension

gs = Θ A (xs - x) / 3600         (1)

or

gh = Θ A (xs - x)

where

gs = amount of evaporated water per second (kg/s)

gh = amount of evaporated water per hour (kg/h)

Θ = (25 + 19 v) = evaporation coefficient (kg/m2h)

v = velocity of air above the water surface (m/s)

A = water surface area (m2)

xs = humidity ratio in saturated air at the same temperature as the water surface (kg/kg)  (kg H2O in kg Dry Air)

x = humidity ratio in the air (kg/kg) (kg H2O in kg Dry Air) */
package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Liquidclasses"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/eng"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// ul

// cubesensor streams:
// in pascals atmospheric pressure of moist air (Pa) 100mBar = 1 pa. Not yet built in unit so we import it from wunit.
// input in deg C will be converted to Kelvin
// Percentage // density water vapor (kg/m3)
// // velocity of air above water in m/s ; could be calculated or measured by an anemometer

// time

// ul/h
// ul

func _EvaporationrateRequirements() {
}
func _EvaporationrateSetup(_ctx context.Context, _input *EvaporationrateInput) {
}
func _EvaporationrateSteps(_ctx context.Context, _input *EvaporationrateInput, _output *EvaporationrateOutput) {
}
func _EvaporationrateAnalysis(_ctx context.Context, _input *EvaporationrateInput, _output *EvaporationrateOutput) {

	var surfacearea wunit.Area
	if _input.Platetype.Welltype.Shape().LengthUnit == "mm" {
		wellarea, err := _input.Platetype.Welltype.CalculateMaxCrossSectionArea()
		if err != nil {
			_output.Warnings = append(_output.Warnings, fmt.Errorf(err.Error()))
		}

		// print statements like this
		fmt.Println("wellarea", wellarea.ToString())
		fmt.Println(_input.Platetype.Welltype.Xdim, _input.Platetype.Welltype.Ydim, _input.Platetype.Welltype.Zdim, _input.Platetype.Welltype.Shape())
		surfacearea = wellarea
	} else {
		_output.Warnings = append(_output.Warnings, fmt.Errorf("plate "+_input.Platetype.String()+" Wellshape "+_input.Platetype.Welltype.String()+" surface area not yet calculated due to bottom type"))
		execute.Errorf(_ctx, "plate "+_input.Platetype.String()+" Wellshape "+_input.Platetype.Welltype.String()+" surface area not yet calculated due to bottom type")
	}
	var PWS float64 = eng.Pws(_input.Temp)
	var pw float64 = eng.Pw(_input.Relativehumidity, PWS) // vapour partial pressure in Pascals

	theta, err := eng.Θ(_input.Liquid.TypeName(), _input.Airvelocity)
	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}
	var Gh = (theta *
		((surfacearea.RawValue() / 1000000) *
			((eng.Xs(PWS, _input.Pa)) - (eng.X(pw, _input.Pa))))) // Gh is rate of evaporation in kg/h
	evaporatedliquid := (Gh * (_input.Executiontime.SIValue() / 3600)) // in kg

	density, ok := liquidclasses.Liquidclass[_input.Liquid.TypeName()]["ro"]

	if !ok {
		density = liquidclasses.Liquidclass["water"]["ro"]
		_output.Warnings = append(_output.Warnings, fmt.Errorf("liquid density not found for "+_input.Liquid.TypeName()+" so used water value"))
	}

	evaporatedliquid = (evaporatedliquid * density) / 1000                         // converted to litres
	_output.Evaporatedliquid = wunit.NewVolume((evaporatedliquid * 1000000), "ul") // convert to ul

	_output.Evaporationrateestimate = Gh * 1000000 // ul/h if declared in parameters or data it doesn't need declaring again

	estimatedevaporationtime := _input.Volumeperwell.ConvertTo(wunit.ParsePrefixedUnit("ul")) / _output.Evaporationrateestimate
	_output.Estimatedevaporationtime = wunit.NewTime((estimatedevaporationtime * 3600), "s")

	_output.Status = fmt.Sprintln("Well Surface Area=",
		surfacearea.ToString(),
		"evaporation rate =", Gh*1000000, "ul/h",
		"total evaporated liquid =", _output.Evaporatedliquid.ToString(), "after", _input.Executiontime.ToString(),
		"estimated evaporation time = ", _output.Estimatedevaporationtime.ToString(),
		"Warnings =", _output.Warnings)

} // works in either analysis or steps sections

func _EvaporationrateValidation(_ctx context.Context, _input *EvaporationrateInput, _output *EvaporationrateOutput) {
	if _output.Evaporatedliquid.SIValue() > _input.Volumeperwell.SIValue() {
		execute.Errorf(_ctx, "not enough liquid, Expected that liquid volume "+_input.Volumeperwell.ToString()+" will evaporate during this time "+_input.Executiontime.ToString()+" Status:  "+_output.Status)
	}
}
func _EvaporationrateRun(_ctx context.Context, input *EvaporationrateInput) *EvaporationrateOutput {
	output := &EvaporationrateOutput{}
	_EvaporationrateSetup(_ctx, input)
	_EvaporationrateSteps(_ctx, input, output)
	_EvaporationrateAnalysis(_ctx, input, output)
	_EvaporationrateValidation(_ctx, input, output)
	return output
}

func EvaporationrateRunSteps(_ctx context.Context, input *EvaporationrateInput) *EvaporationrateSOutput {
	soutput := &EvaporationrateSOutput{}
	output := _EvaporationrateRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func EvaporationrateNew() interface{} {
	return &EvaporationrateElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &EvaporationrateInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _EvaporationrateRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &EvaporationrateInput{},
			Out: &EvaporationrateOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type EvaporationrateElement struct {
	inject.CheckedRunner
}

type EvaporationrateInput struct {
	Airvelocity      wunit.Velocity
	Executiontime    wunit.Time
	Liquid           *wtype.LHComponent
	Pa               wunit.Pressure
	Platetype        *wtype.LHPlate
	Relativehumidity float64
	Temp             wunit.Temperature
	Volumeperwell    wunit.Volume
}

type EvaporationrateOutput struct {
	Estimatedevaporationtime wunit.Time
	Evaporatedliquid         wunit.Volume
	Evaporationrateestimate  float64
	Status                   string
	Warnings                 []error
}

type EvaporationrateSOutput struct {
	Data struct {
		Estimatedevaporationtime wunit.Time
		Evaporatedliquid         wunit.Volume
		Evaporationrateestimate  float64
		Status                   string
		Warnings                 []error
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Evaporationrate",
		Constructor: EvaporationrateNew,
		Desc: component.ComponentDesc{
			Desc: " Evaporation calculator based on\nhttp://www.engineeringtoolbox.com/evaporation-water-surface-d_690.html\n\nThis engineering function may need to be improved to account for vapour pressure and surface tension\n\ngs = Θ A (xs - x) / 3600         (1)\n\nor\n\ngh = Θ A (xs - x)\n\nwhere\n\ngs = amount of evaporated water per second (kg/s)\n\ngh = amount of evaporated water per hour (kg/h)\n\nΘ = (25 + 19 v) = evaporation coefficient (kg/m2h)\n\nv = velocity of air above the water surface (m/s)\n\nA = water surface area (m2)\n\nxs = humidity ratio in saturated air at the same temperature as the water surface (kg/kg)  (kg H2O in kg Dry Air)\n\nx = humidity ratio in the air (kg/kg) (kg H2O in kg Dry Air)\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/eng/Evaporationrate/Evaporationrate.an",
			Params: []component.ParamDesc{
				{Name: "Airvelocity", Desc: "// velocity of air above water in m/s ; could be calculated or measured by an anemometer\n", Kind: "Parameters"},
				{Name: "Executiontime", Desc: "time\n", Kind: "Parameters"},
				{Name: "Liquid", Desc: "", Kind: "Inputs"},
				{Name: "Pa", Desc: "cubesensor streams:\n\nin pascals atmospheric pressure of moist air (Pa) 100mBar = 1 pa. Not yet built in unit so we import it from wunit.\n", Kind: "Parameters"},
				{Name: "Platetype", Desc: "", Kind: "Inputs"},
				{Name: "Relativehumidity", Desc: "Percentage // density water vapor (kg/m3)\n", Kind: "Parameters"},
				{Name: "Temp", Desc: "input in deg C will be converted to Kelvin\n", Kind: "Parameters"},
				{Name: "Volumeperwell", Desc: "ul\n", Kind: "Parameters"},
				{Name: "Estimatedevaporationtime", Desc: "", Kind: "Data"},
				{Name: "Evaporatedliquid", Desc: "ul\n", Kind: "Data"},
				{Name: "Evaporationrateestimate", Desc: "ul/h\n", Kind: "Data"},
				{Name: "Status", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
