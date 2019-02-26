// antha/AnthaStandardLibrary/Packages/devices/devices.go: Part of the Antha language
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

// Look up tables stroring device properties. E.g. for use in calculating relative centrifugal force.
package devices

var (
	Shaker = map[string]map[string]float64{
		"HiGro incubator-shaker": {
			"dt": 0.008, //shaking amplitude diameter in m

		},
		"Thermomixer": {
			"dt": 0.003, //shaking amplitude diameter in m
		},
		"Kuhner": {
			"dt": 0.025, //shaking amplitude diameter in m
		},
		"3000 T-elm": {
			"dt":                 0.002,  //shaking amplitude diameter in m
			"maxrpm":             3000,   // maximum shaking speed in rpm
			"minrpm":             200,    // maximum shaking speed in rpm
			"shakeaccuracy":      25,     // +/- in rpm
			"zeroposaccuracy":    0.0001, //in m
			"heatupRate":         7,      // degrees C per minute
			"maxtemp":            99,     // in deg C
			"tempaccuracy":       0.1,    // in deg C
			"tempuniformity":     0.5,    // at 45 deg C
			"environmaxtemp":     45,     // deg C
			"environmintemp":     5,
			"enrironmaxhumidity": 80, //percent
			"weightinkg":         1.5,
			"Width":              0.142, // in m
			"Depth":              0.099,
			"Height":             0.0575,
		},
		"InhecoStaticOnDeck": {
			"Height": 0.0575, // height in m

		},
	}
)

/*
type ThermalAdaptor struct {
	Manufacturer string
	Cat-number string
	ForPlates []wtype.LHPlate
}
var (
	ShakerAdaptors = map[string]map[string]map[string]float64{
		"3000 T-elm": map[string]map[string]float64{
			"PCR_adaptor": map[string]float64{
			"Height":0.008, //height in m,
			"wells": 96.0,


		},
		},
		}
*/
