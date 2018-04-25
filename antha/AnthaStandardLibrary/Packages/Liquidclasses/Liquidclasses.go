// antha/AnthaStandardLibrary/Packages/Liquidclasses/Liquidclasses.go: Part of the Antha language
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

// Example liquidclasses and liquid properties
package liquidclasses

var (
	//radius = 0.005 // radius of surface in m

	// map of liquid classes

	Liquidclass = map[string]map[string]float64{
		"water": {
			"c":     25.0,
			"d":     19.0,
			"xs":    0.1,    //0.62198 * pws / (pa - pws), // humidity ratio in saturated air at the same temperature as the water surface (kg/kg)  (kg H2O in kg Dry Air)
			"x":     0.01,   //0.62198 * pw / (pa - pw),   // equations not working
			"ro":    999.97, //density, kg 􏰀/ m􏰁3
			"mu":    0.001,  //viscosity at 20degrees, kg 􏰀/ m􏰁 /􏰀 s
			"sigma": 0.072,  //Surface tension in N/m from Wikipedia at 25 degrees
			"Lm":    334000, //(J/kg) //latent heat of melting (J/kg, Btu/lb)
		},
		"ethanol": { // dummy data
			"c":  100.0,
			"d":  19.0,
			"xs": 0.1, // humidity ratio in saturated air at the same temperature as the water surface (kg/kg)  (kg H2O in kg Dry Air)
			"x":  0.01,
		},
	}
)
