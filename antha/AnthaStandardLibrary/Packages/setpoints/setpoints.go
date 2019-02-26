// antha/AnthaStandardLibrary/Packages/setpoints/setpoints.go: Part of the Antha language
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

/* Islam, R. S., Tisi, D., Levy, M. S. & Lye, G. J. Scale-up of Escherichia coli growth and recombinant protein expression conditions from microwell to laboratory and pilot scale based on matched kLa. Biotechnol. Bioeng. 99, 1128–1139 (2008).

equation (6)

func kLa_squaremicrowell = (3.94 x 10E-4) * (D/dv)* ai * (ro * n * dv * 2/mu)^1.91 * exp ^ (a * dt(2 * math.Pi * n)^2 /(2 * g)^b) // a little unclear whether exp is e to (afr^b) from paper but assumed this is the case

kla = dimensionless
	var D = diffusion coefficient, m2 􏰀 s􏰁1
	var dv = microwell vessel diameter, m
	var ai = initial specific surface area, m􏰁1
	var RE = Reynolds number, (ro * n * dv * 2/mu), dimensionless
		var	ro	= density, kg 􏰀/ m􏰁3
		var	n 	= shaking frequency, s􏰁1
		var	mu	= viscosity, kg 􏰀/ m􏰁 /􏰀 s
	const exp = Eulers number, 2.718281828

	var Fr = Froude number = dt(2 * math.Pi * n)^2 /(2 * g), (dimensionless)
		var dt = shaking amplitude, m
		const g = acceleration due to gravity, m 􏰀/ s􏰁2
	const	a = constant
	const	b = constant
*/

// Functions used in calculating mass transfer in microwells
package setpoints //masstransfer

import (
	//"fmt"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Labware"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Liquidclasses"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/devices"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/eng"
	//"github.com/montanaflynn/stats"
)

func CalculateKlasquaremicrowell(Platetype string, liquid string, rpm float64, Shakertype string, TargetRE float64, D float64) (CalculatedKla float64) {

	// find relevant properties from labware
	dv := labware.Labwaregeometry[Platetype]["dv"] // microwell vessel diameter, m 0.017 //
	ai := labware.Labwaregeometry[Platetype]["ai"] // initial specific surface area, /m 96.0
	//var RE = Reynolds number, (ro * n * dv * 2/mu), dimensionless
	//find relevant properties from liquid type
	ro := liquidclasses.Liquidclass[liquid]["ro"] //density, kg 􏰀/ m􏰁3 999.7 // environment dependent
	mu := liquidclasses.Liquidclass[liquid]["mu"] //0.001           environment dependent                        //liquidclasses.Liquidclass[liquid]["mu"] viscosity, kg 􏰀/ m􏰁 /􏰀 s

	n := rpm / 60 //shaking frequency, s􏰁1
	//const exp = Eulers number, 2.718281828

	//Fr = Froude number = dt(2 * math.Pi * n)^2 /(2 * g), (dimensionless)

	// find relevant properties of shaker
	dt := devices.Shaker[Shakertype]["dt"] //0.008                                  //shaking amplitude, m // move to shaker package

	a := labware.Labwaregeometry[Platetype]["a"] //0.88   //
	b := labware.Labwaregeometry[Platetype]["b"] //1.24

	Fr := eng.Froude(dt, n, eng.G)
	Re := eng.RE(ro, n, mu, dv)
	//Necessaryshakerspeed := eng.Shakerspeed(TargetRE, ro, mu, dv)
	//r, _ := stats.Round(Necessaryshakerspeed*60, 3)
	//// fmt.Println("shakerspeedrequired= ", r)
	//r, _ = stats.Round(Re, 3)
	//// fmt.Println("Reynolds number = ", r)
	// if Re > 5E3 {
	// 	// fmt.Println("Turbulent flow")
	// }

	//r, _ = stats.Round(Fr, 3)
	//// fmt.Println("Froude number = ", r)
	CalculatedKla = eng.KLa_squaremicrowell(D, dv, ai, Re, a, Fr, b)

	//r, _ = stats.Round(CalculatedKla, 3)
	//// fmt.Println("kla =", r)

	// trouble shooting
	/*
		fmt.Println(D / dv)
		fmt.Println(math.Pow(Re, 1.91))
		fmt.Println(math.Pow(math.E, (a * (math.Pow(Fr, b)))))
		fmt.Println(a * (math.Pow(Fr, b)))
		fmt.Println(math.Pow(Fr, b))
		fmt.Println(math.E)
	*/
	return CalculatedKla
}
