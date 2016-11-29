//status = compiles and calculates; need to fill in correct parameters and check units
//currently using dummy values only so won't be accurate yet!
package lib

import (
	"fmt"                                                                 // we need this go library to print
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/eng" // all of our functions used here are in the Thaw.go file in the eng package which this points to
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
)

// Many of the real parameters required will be looked up via the specific labware (platetype) and liquid type which are being used.

/* e.g. the sample volume as frozen by a previous storage protocol;
could be known or measured via liquid height detection on some liquid handlers */

// These should be captured via sensors just prior to execution

// This will be monitored via the thermometer in the freezer in which the sample was stored

/* This will offer another knob to tweak (in addition to the other parameters) as a means to improve
the correlation over time as we see how accurate the calculator is in practice */

func _ThawtimeRequirements() {
}
func _ThawtimeSetup(_ctx context.Context, _input *ThawtimeInput) {
}
func _ThawtimeSteps(_ctx context.Context, _input *ThawtimeInput, _output *ThawtimeOutput) {
	/*  Step 1. we need a mass for the following equations so we calculate this by looking up
	the liquid density and multiplying by the fill volume using this function from the engineering library */

	mass := eng.Massfromvolume(_input.Fillvolume, _input.Liquid)

	/*  Step 2. Required heat energy to melt the solid is calculated using the calculated mass along with the latent heat of melting
	which we find via a liquid class look up package which is not required for import here since it's imported from the engineering library */

	q := eng.Q(_input.Liquid, mass)

	/*  Step 3. Heat will be transferred via both convection through the air and conduction through the plate walls.
	Let's first work out the heat energy transferred via convection, this uses an empirical parameter,
	the convective heat transfer coefficient of air (HC_air), this is calculated via another function in the eng library.
	In future we will make this process slightly more sophisticated by adding conditions, since this empirical equation is
	only validated between air velocities 2 - 20 m/s. It could also be adjusted to calculate heat transfer if the sample
	is agitated on a shaker to speed up thawing. */

	hc_air := eng.Hc_air(_input.Airvelocity.SIValue())

	/*  Step 4. The rate of heat transfer by convection is then calculated using this value combined with the temperature differential
	(measured by the temp sensor) and surface area dictated by the plate type (another look up called from the eng library!)*/

	convection := eng.ConvectionPowertransferred(hc_air, _input.Platetype, _input.SurfaceTemp, _input.BulkTemp)

	/*  Step 5. We now estimate the heat transfer rate via conduction. For this we need to know the thermal conductivity of the plate material
	along with the wall thickness. As before, both of these are looked up via the labware library called by this function in the eng library */

	conduction := eng.ConductionPowertransferred(_input.Platetype, _input.SurfaceTemp, _input.BulkTemp)

	/*  Step 6. We're now ready to estimate the thawtime needed by simply dividing the estimated heat required to melt/thaw (i.e. q from step 2)
	by the combined rate of heat transfer estimated to occur via both convection and conduction */
	_output.Estimatedthawtime = eng.Thawtime(convection, conduction, q)

	/* Step 7. Since there're a lot of assumptions here (liquid behaves as water, no change in temperature gradient, no heat transferred via radiation,
	imprecision in the estimates and 	empirical formaulas) we'll multiply by a fudgefactor to be safer that we've definitely thawed,
	this (and all parameters!) can be adjusted over time as we see emprically how reliable this function is as more datapoints are collected */
	_output.Thawtimeused = wunit.NewTime(_output.Estimatedthawtime.SIValue()*_input.Fudgefactor, "s")

	_output.Status = fmt.Sprintln("For", mass.ToString(), "of", _input.Liquid, "in", _input.Platetype,
		"Thawtime required =", _output.Estimatedthawtime.ToString(),
		"Thawtime used =", _output.Thawtimeused.ToString(),
		"power required =", q, "J",
		"HC_air (convective heat transfer coefficient=", hc_air,
		"Convective power=", convection, "J/s",
		"conductive power=", conduction, "J/s")

}
func _ThawtimeAnalysis(_ctx context.Context, _input *ThawtimeInput, _output *ThawtimeOutput) {

}

func _ThawtimeValidation(_ctx context.Context, _input *ThawtimeInput, _output *ThawtimeOutput) {

}
func _ThawtimeRun(_ctx context.Context, input *ThawtimeInput) *ThawtimeOutput {
	output := &ThawtimeOutput{}
	_ThawtimeSetup(_ctx, input)
	_ThawtimeSteps(_ctx, input, output)
	_ThawtimeAnalysis(_ctx, input, output)
	_ThawtimeValidation(_ctx, input, output)
	return output
}

func ThawtimeRunSteps(_ctx context.Context, input *ThawtimeInput) *ThawtimeSOutput {
	soutput := &ThawtimeSOutput{}
	output := _ThawtimeRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ThawtimeNew() interface{} {
	return &ThawtimeElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ThawtimeInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ThawtimeRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ThawtimeInput{},
			Out: &ThawtimeOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type ThawtimeElement struct {
	inject.CheckedRunner
}

type ThawtimeInput struct {
	Airvelocity wunit.Velocity
	BulkTemp    wunit.Temperature
	Fillvolume  wunit.Volume
	Fudgefactor float64
	Liquid      string
	Platetype   string
	SurfaceTemp wunit.Temperature
}

type ThawtimeOutput struct {
	Estimatedthawtime wunit.Time
	Status            string
	Thawtimeused      wunit.Time
}

type ThawtimeSOutput struct {
	Data struct {
		Estimatedthawtime wunit.Time
		Status            string
		Thawtimeused      wunit.Time
	}
	Outputs struct {
	}
}

func init() {
	if err := addComponent(component.Component{Name: "Thawtime",
		Constructor: ThawtimeNew,
		Desc: component.ComponentDesc{
			Desc: "status = compiles and calculates; need to fill in correct parameters and check units\ncurrently using dummy values only so won't be accurate yet!\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/eng/Thawtime/Thawtime.an",
			Params: []component.ParamDesc{
				{Name: "Airvelocity", Desc: "These should be captured via sensors just prior to execution\n", Kind: "Parameters"},
				{Name: "BulkTemp", Desc: "This will be monitored via the thermometer in the freezer in which the sample was stored\n", Kind: "Parameters"},
				{Name: "Fillvolume", Desc: " e.g. the sample volume as frozen by a previous storage protocol;\n\tcould be known or measured via liquid height detection on some liquid handlers\n", Kind: "Parameters"},
				{Name: "Fudgefactor", Desc: " This will offer another knob to tweak (in addition to the other parameters) as a means to improve\n\tthe correlation over time as we see how accurate the calculator is in practice\n", Kind: "Parameters"},
				{Name: "Liquid", Desc: "", Kind: "Parameters"},
				{Name: "Platetype", Desc: "Many of the real parameters required will be looked up via the specific labware (platetype) and liquid type which are being used.\n", Kind: "Parameters"},
				{Name: "SurfaceTemp", Desc: "", Kind: "Parameters"},
				{Name: "Estimatedthawtime", Desc: "", Kind: "Data"},
				{Name: "Status", Desc: "", Kind: "Data"},
				{Name: "Thawtimeused", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
