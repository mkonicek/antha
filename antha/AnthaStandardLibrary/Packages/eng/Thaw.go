// antha/AnthaStandardLibrary/Packages/eng/Thaw.go: Part of the Antha language
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
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Labware"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Liquidclasses"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"math"
)

/*Heat Required to Melt a Solid
The heat required to melt a solid can be calculated as

q = Lm m  (1)

where

q = required heat (J, Btu)

Lm = latent heat of melting (J/kg, Btu/lb)

m = mass of subsance (kg, lb)

Example - Required Heat to melt Ice (Water)

The heat required to melt 10 kg of water can be calculated as

q = Lm m

  = 334 103 (J/kg) 10 (kg)

  = 3340000 (J)

  = 3340 (kJ)
*/

func Massfromvolume(Volume wunit.Volume, liquid string) (Mass wunit.Mass) {
	volume := Volume.SIValue()
	mass := ((liquidclasses.Liquidclass[liquid]["ro"]) * volume)
	Mass = wunit.NewMass(mass, "g")
	return Mass
}

// Q = required heat (J, Btu)

func Q(liquid string, Mass wunit.Mass) (q float64) {
	return ((liquidclasses.Liquidclass[liquid]["Lm"]) * Mass.SIValue())
}

/* The convective heat transfer coefficient of air is approximately equal to

hc = 10.45 - v + 10 v1/2    (2)

where

v = the relative speed of the object through the air (m/s)

*/

/*
q = hc A dT         (1)

where

q = heat transferred per unit time (W)

A = heat transfer area of the surface (m2)

hc= convective heat transfer coefficient of the process (W/(m2K) or W/(m2oC))

dT = temperature difference between the surface and the bulk fluid (K or oC)

*/
//heat transferred by surface via air
func Hc_air(v float64) (hc_air float64) {
	vhalf := math.Pow(v, 0.5)
	Hc_airpow := (10 * vhalf) // v in m/s // Note! - this is an empirical equation and can be used for velocities - v - from 2 to 20 m/s.
	hc_air = (10.45 - v + Hc_airpow)
	return hc_air
}

func ConvectionPowertransferred(hc_air float64, Platetype string, SurfaceTemp wunit.Temperature, BulkTemp wunit.Temperature) (convectionpowertransferred float64) {
	surfaceTemp := SurfaceTemp.SIValue()
	bulkTemp := BulkTemp.SIValue()
	return (hc_air * labware.Labwaregeometry[Platetype]["A"] * (surfaceTemp - bulkTemp))
}

/*
P=kAΔT/Δx where P is the power transferred,
k is thermal conductivity,
A is the area of the surface through which energy will flow,
ΔTΔx is temperature gradient.
ΔT is the temperature difference between inner and outer surface,
Δx is the thickness of container.

Reference https://www.physicsforums.com/threads/calculate-how-long-it-will-take-for-the-ice-to-melt.531908/

*/

func ConductionPowertransferred(Platetype string, SurfaceTemp wunit.Temperature, BulkTemp wunit.Temperature) (conductionpowertransferred float64) { // W or J/s
	surfaceTemp := SurfaceTemp.SIValue()
	bulkTemp := BulkTemp.SIValue()
	return (labware.Labwaregeometry[Platetype]["k"] * labware.Labwaregeometry[Platetype]["A"] * ((surfaceTemp - bulkTemp) / labware.Labwaregeometry[Platetype]["Δx"]))
}

func Thawtime(convectionpowertransferred float64, conductionpowertransferred float64, q float64) (Thawtimerequired wunit.Time) { //(liquid string, airvelocity float64) float64 {

	thawtimerequired := q / (convectionpowertransferred + conductionpowertransferred)

	Thawtimerequired = wunit.NewTime(thawtimerequired, "s")

	return Thawtimerequired //(liquidclasses.Liquidclass[liquid]["c"]) + ((liquidclasses.Liquidclass[liquid]["d"]) * airvelocity)
}

func Thawfloat(convectionpowertransferred float64, conductionpowertransferred float64, q float64) (Thawtimerequired float64) { //(liquid string, airvelocity float64) float64 {

	Thawtimerequired = q / (convectionpowertransferred + conductionpowertransferred)

	//Thawtimerequired = wunit.NewTime(thawtimerequired, "s")

	return Thawtimerequired //(liquidclasses.Liquidclass[liquid]["c"]) + ((liquidclasses.Liquidclass[liquid]["d"]) * airvelocity)
}
