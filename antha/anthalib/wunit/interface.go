// wunit/interface.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
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

package wunit

import (
	"fmt"
)

// units mapped by string
var unitMap map[string]GenericUnit

// helper function to make it easier to
// make a new unit with prefix directly
func NewPrefixedUnit(prefix string, unit string) *GenericPrefixedUnit {
	u := UnitBySymbol(unit)
	p := SIPrefixBySymbol(prefix)

	gpu := GenericPrefixedUnit{u, p}
	return &gpu
}

// get a unit from a string

func ParsePrefixedUnit(unit string) *GenericPrefixedUnit {
	parser := &SIPrefixedUnitGrammar{}
	parser.SIPrefixedUnit.Init([]byte(unit))

	if err := parser.Parse(unit); err != nil {
		e := fmt.Errorf("cannot parse %s: %s", unit, err.Error())
		panic(e)
	}

	prefix := ""
	var un string

	if len(parser.TreeTop.Children) == 1 {
		un = parser.TreeTop.Children[0].Value.(string)
	} else {
		prefix = parser.TreeTop.Children[0].Value.(string)
		un = parser.TreeTop.Children[1].Value.(string)
	}
	return NewPrefixedUnit(prefix, un)
}

// look up unit by symbol
func UnitBySymbol(sym string) GenericUnit {
	if unitMap == nil {
		unitMap = Make_units()
	}

	return unitMap[sym]
}

// generate an initial unit library
func Make_units() map[string]GenericUnit {

	units := []string{"M", "min", "l", "L", "g", "V", "J", "A", "N", "s", "radians", "degrees", "rads", "Hz", "rpm", "℃", "M/l", "g/l", "J/kg*C", "Pa", "kg/m^3", "/s", "/min", "per", `/`, "m/s", "m^2", "mm^2", "ml/min", "kg/l", "X", "U/l", "m", "v/v"}
	unitnames := []string{"mole", "minute", "litre", "litre", "Gramme", "Volt", "Joule", "Ampere", "Newton", "second", "radian", "degree", "radian", "Herz", "revolutions per minute", "Celsius", "Mol/litre", "g/litre", "Joule per kilogram per degrees celsius", "Pascal", "kg per cubic meter", "per second", "per minute", "per", "per", "metres per second", "square metres", "square metres", "millilitres/minute", "kilogram per litre", "times", "Units	 per L", "metres", "volume/volume"}
	//unitdimensions:=[]string{"amount", "time", "length^3", "length^3", "mass", "mass*length/time^2*charge", "mass*length^2/time^2", "charge/time", "charge", "mass*length/time^2", "time", "angle", "angle", "angle", "time^-1", "angle/time", "temperature", "velocity}

	unitbaseconvs := []float64{1, 0.1666666666666666667, 1, 1, 0.001, 1, 1, 1, 1, 1, 1, 0.01745329251994, 1, 1, 1, 1, 1, 0.001, 1, 1, 1, 1, 0.1666666666666666667, 1, 1, 1, 1, 0.000001, 1., 1, 1, 1, 1, 1}

	unit_map := make(map[string]GenericUnit, len(units))

	for i, u := range units {

		baseunit := u

		if u == "g" {
			baseunit = "kg"
		} else if u == "g/l" {
			baseunit = "kg/l"
		} else if u == "mm^2" {
			baseunit = "m^2"
		}
		gu := GenericUnit{unitnames[i], u, unitbaseconvs[i], baseunit}
		unit_map[u] = gu
	}

	return unit_map
}
