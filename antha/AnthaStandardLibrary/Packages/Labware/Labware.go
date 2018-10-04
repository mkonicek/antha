// antha/AnthaStandardLibrary/Packages/Labware/Labware.go: Part of the Antha language
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

// Example labware definitions
package labware

var (
	//radius = 0.005 // radius of surface in m

	// map of some labware properties

	Labwaregeometry = map[string]map[string]float64{
		"96DSW_axygen": {
			"numberofwells": 96.0,
			//"Surfacearea":   0.000005,
			"Height": 0.01, // humidity ratio in saturated air at the same temperature as the water surface (kg/kg)  (kg H2O in kg Dry Air)
			"Radius": 0.001,
			"dv":     0.017,
		},
		"96DRW_axygen": {
			"numberofwells": 96.0,
			//"Surfacearea":   0.00004,
			"Height": 10.0, // humidity ratio in saturated air at the same temperature as the water surface (kg/kg)  (kg H2O in kg Dry Air)
			"Radius": 0.001,
		},
		"DSW24": { //24 square wells, pyramidal bases (DOT Scientific, Inc., Burton, MI)
			"numberofwells":     24.0,
			"numberofwellsides": 4.0,
			"height_m":          0.043, //m
			"width_m":           0.017, //m
			"breadth_m":         0.017, //m
			"dv":                0.017, //m
			//"Surfacearea":       0.000289, //m2 surface area in contact with air when stationary
			"Height": 43.0, //mm used in evaporation calculator ... correct to SI units
			"Radius": 17.0, //mm used in evaporation calculator ... correct to SI units
			"a":      0.88,
			"b":      1.24,
			"ai":     96.0,   //initial specific surface area, /m
			"Δx":     0.002,  // wall thickness m
			"k":      0.1,    // plastic thermal conductivity J s-1 m-1 °C-1 for polypropylene ~0.1-0.22 http://www.engineeringtoolbox.com/thermal-conductivity-d_429.html
			"A":      0.0003, // plasticsurfacearea m2 in contact with ?
		},
		"pcrplate_skirted": {
			"numberofwells":     96.0,
			"numberofwellsides": 1.0,
			"height_m":          0.02075, //m
			"width_m":           0.0055,  //m
			"breadth_m":         0.0055,  //m
			"dv":                0.0055,  //m
			//	"Surfacearea":       0.000023761375, //m2 surface area in contact with air when stationary
			"Height": 20.75,    //mm used in evaporation calculator ... correct to SI units
			"Radius": 2.75,     //mm used in evaporation calculator ... correct to SI units
			"a":      0.88,     //??
			"b":      1.24,     //??
			"ai":     96.0,     //initial specific surface area, /m
			"Δx":     0.001,    // wall thickness m estimate... check accuracy!
			"k":      0.1,      // plastic thermal conductivity J s-1 m-1 °C-1 for polypropylene ~0.1-0.22 http://www.engineeringtoolbox.com/thermal-conductivity-d_429.html
			"A":      0.000198, // plasticsurfacearea m2 in contact with ?
		},
		"greiner_384": { //
			"numberofwells":     384.0,
			"numberofwellsides": 4.0,
			"height_m":          0.01,  //m
			"width_m":           0.005, //m
			"breadth_m":         0.005, //m
			"dv":                0.005, //m
			//	"Surfacearea":       0.000025, //m2
			"Height": 10.0, //mm used in evaporation calculator ... correct to SI units
			"Radius": 5.0,  //mm used in evaporation calculator ... correct to SI units
			"a":      0.88, // ??
			"b":      1.24, //??
			"ai":     96.0, //initial specific surface area, /m ???
		},
	}
)

/* to do:
1. add func to calculate area based on radius, shape dimensions etc?

2. add labware material properties maps (probably linking to external map or db): surface tension, solvent compatibility, leachability, biocompatibility etc?
*/
