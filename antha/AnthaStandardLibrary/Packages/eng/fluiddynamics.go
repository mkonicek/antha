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
	"math"

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

func KLa_squaremicrowell(D float64, dv float64, ai float64, RE float64, a float64, froude float64, b float64) float64 {
	return ((3.94E-4) * (D / dv) * ai * (math.Pow(RE, 1.91)) * (math.Pow(math.E, (a * (math.Pow(froude, b))))))
} // a little unclear whether exp is e to (afr^b) from paper but assumed this is the case

func RE(ro float64, n float64, mu float64, dv float64) float64 { // Reynolds number

	return (ro * n * dv * 2 / mu)
}

func Shakerspeed(TargetRE float64, ro float64, mu float64, dv float64) wunit.Rate /*float64*/ { // calulate shaker speed from target Reynolds number
	rps := (TargetRE * mu / (ro * dv * 2))
	rate, err := wunit.NewRate(rps, "/s")
	if err != nil {
		panic(err)
	}

	return rate
}

func Froude(dt float64, n float64, g float64) float64 { // froude number  dt = shaken diamter in m
	return (dt * (math.Pow((2 * math.Pi * n), 2)) / (2 * g))
}

const G float64 = 9.81 //acceleration due to gravity in meters per second squared

//Micheletti 2006:

func Ncrit_srw(sigma float64, dv float64, Vl float64, ro float64, dt float64) wunit.Rate {
	rps := math.Sqrt((sigma * dv) / (4 * math.Pi * Vl * ro * dt)) //unit = per S // established for srw with Vl = 200ul
	rate, err := wunit.NewRate(rps, "/s")
	if err != nil {
		panic(err)
	}
	return rate
}
