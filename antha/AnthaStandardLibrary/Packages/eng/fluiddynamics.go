// antha/AnthaStandardLibrary/Packages/eng/fluiddynamics.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

// Package for performing engineering calculations; at present this consists of evaporation rate estimation, thawtime estimation and fluid dynamics
package eng

import (
	"fmt"

	"math"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

//Equations from Islam et al:

func CentripetalForce(mass wunit.Mass, angularfrequency float64, radius wunit.Length) (force wunit.Force) {
	forcefloat := mass.SIValue() * math.Pow(angularfrequency, 2) * radius.SIValue()
	force = wunit.NewForce(forcefloat, "N")
	return force
}

func Angularfrequency(frequency float64) (angularfrequency float64) {
	return 2 * math.Pi * frequency
}

/*
func KLa_squaremicrowell(D float64, dv float64, ai float64, RE float64, a float64, froude float64, b float64) float64 {

	log := math.Log((3.94E-4)) + math.Log((D / dv)) + math.Log(math.Pow(RE, 1.91)) + (a * (math.Pow(froude, b)))

	// fmt.Println("e ^", log)

	kla := math.Exp(log)

	return kla
} // a little unclear whether exp is e to (afr^b) from paper but assumed this is the case
*/

func KLaSquareMicrowell(D float64, wellDiameter wunit.Length, ai float64, RE float64, a float64, froude float64, b float64) (float64, error) {

	/*
		klainputs := fmt.Sprintln("D: ",D,"dv: ", dv,"ai: ", ai,"Re: ", Re,"a: ", a,"Fr: ", Fr,"b: ", b)

		var errorMessages []string = []string{
				fmt.Sprintln("Calculation inputs: ", klainputs),
				fmt.Sprintln("Derived values: ", "math.Pow(RE, 1.91): ", math.Pow(Re, 1.91), "math.Pow(froude, b): ", math.Pow(Fr, b), "(math.Pow(math.E, (a * (math.Pow(froude, b))))): ", (math.Exp(a * (math.Pow(Fr, b)))), "a * (math.Pow(froude, b)): ", a*(math.Pow(Fr, b))),
				fmt.Sprintln("e", math.E, "power", (a * (math.Pow(Fr, b)))),
		}


		if math.IsNaN(CalculatedKla) {
			Errorf("Calculated Kla is not a number. %s",strings.Join(errorMessages,";"))
		}

		if math.IsInf(CalculatedKla,0) {
			Errorf("Calculated Kla is infinite. %s",strings.Join(errorMessages,";"))
		}
	*/
	//

	klainputs := fmt.Sprintln(fmt.Sprintln("D = ", D), fmt.Sprintln("wellDiameter = ", wellDiameter), fmt.Sprintln("ai = ", ai), fmt.Sprintln("RE = ", RE), fmt.Sprintln("a = ", a), fmt.Sprintln("Froude number = ", froude), fmt.Sprintln("b = ", b))

	dv := wellDiameter.ConvertTo(wunit.ParsePrefixedUnit("m"))

	part1 := ((3.94E-4) * (D / dv) * ai * (math.Pow(RE, 1.91)))

	exponent := a * (math.Pow(froude, b))

	part2a := float64(math.Pow(math.E, exponent))

	var errorMessages []string = []string{
		fmt.Sprintln(fmt.Sprintln("Calculation inputs: "), klainputs),
		fmt.Sprintln(fmt.Sprintln("Derived values: "),
			fmt.Sprintln("math.Pow(RE, 1.91) = ", math.Pow(RE, 1.91)),
			fmt.Sprintln("math.Pow(froude, b) = ", math.Pow(froude, b)),
			fmt.Sprintln("exponent = a * (math.Pow(froude, b)) = ", a*(math.Pow(froude, b))),
			fmt.Sprintln("e", math.E, "^", (a*(math.Pow(froude, b))), "= ", math.Pow(math.E, exponent))),
	}

	if math.IsNaN(part2a) {
		err := fmt.Errorf("Calculated Kla is not a number due to sub calculation. %s = %f where exponent %s =  %f", "math.Pow(math.E, exponent) , %s = %f", "a * (math.Pow(froude, b)). Inputs = %s", a*(math.Pow(froude, b)), errorMessages)
		return 0.0, err
	}

	if math.IsInf(part2a, 0) {
		err := fmt.Errorf("Calculated Kla is infinity due to sub calculation. %s = %f. Inputs = %s", "math.Pow(math.E, exponent)", math.Pow(math.E, exponent), strings.Join(errorMessages, "\n"))
		return 0.0, err
	}

	part2 := float64(part2a)

	klaresult := part1 * part2

	return klaresult, nil
	//

	// original return
	//return ((3.94E-4) * (D / dv) * ai * (math.Pow(RE, 1.91)) * (math.Pow(math.E, (a * (math.Pow(froude, b)))))), nil

} // a little unclear whether exp is e to (afr^b) from paper but assumed this is the case

/*

func KLa_squaremicrowell(D float64, dv float64, ai float64, RE float64, a float64, froude float64, b float64) float64 {

	part1 := ((3.94E-4) * (D / dv) * ai * (math.Pow(RE, 1.91)))

	exponent := a * (math.Pow(froude, b))

	part2a := float64(math.Pow(math.E, exponent))

	part2 := float64(part2a)

	klaresult := part1 * part2

	return klaresult
} // a little unclear whether exp is e to (afr^b) from paper but assumed this is the case
*/

// RE calculates the Reynolds number for mixing fluid in a shaken microwell plate.
// is dv really the well diameter and not the shaking amplitude?
func RE(ro float64, n float64, mu float64, wellDiameter wunit.Length) float64 { // Reynolds number

	dv := wellDiameter.ConvertTo(wunit.ParsePrefixedUnit("m"))

	return (ro * n * dv * 2 / mu)
}

func ShakerSpeed(TargetRE float64, ro float64, mu float64, wellDiameter wunit.Length) (rate wunit.Rate, err error) /*float64*/ { // calulate shaker speed from target Reynolds number

	dv := wellDiameter.ConvertTo(wunit.ParsePrefixedUnit("m"))

	rps := (TargetRE * mu / (ro * dv * 2))
	rate, err = wunit.NewRate(rps, "/s")
	return rate, err
}

func Froude(shakingDiameter wunit.Length, n float64, g float64) float64 { // froude number  dt = shaken diamter in m
	dt := shakingDiameter.ConvertTo(wunit.ParsePrefixedUnit("m"))
	return (dt * (math.Pow((2 * math.Pi * n), 2)) / (2 * g))
}

const G float64 = 9.81 //acceleration due to gravity in meters per second squared

// NcritSRW calculates the minimal mixing velocity, in revolutions per second,
// to achieve turbulent flow when mixing a shaken shallow round well microtitre plate.
//
// Calculation from Micheletti 2006
// sigma = Surface tension in N/m
// dv = microwell vessel diameter in m
// ro = density, kg / m^3
// dt = shaking amplitude
func NcritSRW(sigma float64, dv wunit.Length, liquidVolume wunit.Volume, ro float64, dt float64) (rate wunit.Rate) {
	rps := math.Sqrt((sigma * dv.ConvertTo(wunit.ParsePrefixedUnit("m"))) / (4 * math.Pi * liquidVolume.SIValue() * ro * dt)) //unit = per S // established for srw with Vl = 200ul
	rate, _ = wunit.NewRate(rps, "/s")
	return rate
	//sigma = liquid surface tension N /m; dt = shaken diamter in m
}
