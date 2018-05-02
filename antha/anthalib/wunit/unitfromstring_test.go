// wunit/unitfromstring.go: Part of the Antha language
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
	"reflect"
	"testing"
)

type unitFromStringTest struct {
	ComponentName     string
	ContainsConc      bool
	Conc              Concentration
	ComponentNameOnly string
}

var componentWithConcstests = []unitFromStringTest{
	{
		ComponentName:     "Glucose (M)",
		ContainsConc:      true,
		Conc:              NewConcentration(0.0, "M"),
		ComponentNameOnly: "Glucose",
	},
	{
		ComponentName:     "Glucose (mM)",
		ContainsConc:      true,
		Conc:              NewConcentration(0.0, "mM"),
		ComponentNameOnly: "Glucose",
	},
	{
		ComponentName:     "5g/L Glucose",
		ContainsConc:      true,
		Conc:              NewConcentration(5.0, "g/L"),
		ComponentNameOnly: "Glucose",
	},
	{
		ComponentName:     "5 g/L Glucose",
		ContainsConc:      true,
		Conc:              NewConcentration(5.0, "g/L"),
		ComponentNameOnly: "Glucose",
	},
	{
		ComponentName:     "Glucose 5g/L",
		ContainsConc:      true,
		Conc:              NewConcentration(5.0, "g/L"),
		ComponentNameOnly: "Glucose",
	},
	{
		ComponentName:     "1mM/l C6",
		ContainsConc:      true,
		Conc:              NewConcentration(1.0, "mM/l"),
		ComponentNameOnly: "C6",
	},
	{
		ComponentName:     "C6",
		ContainsConc:      false,
		Conc:              NewConcentration(0.0, "mM"),
		ComponentNameOnly: "C6",
	},
	{
		ComponentName:     "1 mM Ammonium Sulphate",
		ContainsConc:      true,
		Conc:              NewConcentration(1.0, "mM"),
		ComponentNameOnly: "Ammonium Sulphate",
	},
	{
		ComponentName:     "Ammonium Sulphate (mM)",
		ContainsConc:      true,
		Conc:              NewConcentration(0.0, "mM"),
		ComponentNameOnly: "Ammonium Sulphate",
	},
	{
		ComponentName:     "E.coli SuperFolder GFP (g/L)",
		ContainsConc:      true,
		Conc:              NewConcentration(0.0, "g/L"),
		ComponentNameOnly: "E.coli SuperFolder GFP",
	},
	{
		ComponentName:     "solutionX",
		ContainsConc:      false,
		Conc:              NewConcentration(0.0, "g/L"),
		ComponentNameOnly: "solutionX",
	},
	{
		ComponentName:     "1X solutionX",
		ContainsConc:      true,
		Conc:              NewConcentration(1.0, "X"),
		ComponentNameOnly: "solutionX",
	},
	{
		ComponentName:     "solutionX (X)",
		ContainsConc:      true,
		Conc:              NewConcentration(0.0, "X"),
		ComponentNameOnly: "solutionX",
	},
	{
		ComponentName:     "solutionX X",
		ContainsConc:      true,
		Conc:              NewConcentration(0.0, "X"),
		ComponentNameOnly: "solutionX",
	},
	{
		ComponentName:     "1mM rumm",
		ContainsConc:      true,
		Conc:              NewConcentration(1.0, "mM"),
		ComponentNameOnly: "rumm",
	},
	{
		ComponentName:     "rumm (mM)",
		ContainsConc:      true,
		Conc:              NewConcentration(0.0, "mM"),
		ComponentNameOnly: "rumm",
	},
	{
		ComponentName:     "X",
		ContainsConc:      false,
		Conc:              NewConcentration(0.0, "g/L"),
		ComponentNameOnly: "X",
	},
	{
		ComponentName:     "1X",
		ContainsConc:      true,
		Conc:              NewConcentration(1, "X"),
		ComponentNameOnly: "1X",
	},
	{
		ComponentName:     "(1X)",
		ContainsConc:      true,
		Conc:              NewConcentration(1, "X"),
		ComponentNameOnly: "(1X)",
	},
}

type volTest struct {
	VolString    string
	Volume       Volume
	ErrorMessage string
}

var volTests = []volTest{
	{
		VolString:    "10ul",
		Volume:       NewVolume(10, "ul"),
		ErrorMessage: "",
	},
	{
		VolString:    "10 ul",
		Volume:       NewVolume(10, "ul"),
		ErrorMessage: "",
	},
	{
		VolString:    "10",
		Volume:       Volume{},
		ErrorMessage: "no valid unit found for 10: valid units are: [L l mL ml nL nl pL pl uL ul]",
	},
}

func TestParseVolume(t *testing.T) {

	for _, test := range volTests {
		vol, err := ParseVolume(test.VolString)
		if !reflect.DeepEqual(vol, test.Volume) {
			t.Error(
				"for", fmt.Sprintf("%+v", test), "\n",
				"Expected:", test.Volume, "\n",
				"Got:", vol, "\n",
			)
		}
		if err != nil {
			if err.Error() != test.ErrorMessage {
				t.Error(
					"for", fmt.Sprintf("%+v", test), "\n",
					"Expected error:", test.ErrorMessage, "\n",
					"Got:", err.Error(), "\n",
				)
			}
		}

		if err == nil && test.ErrorMessage != "" {
			t.Error(
				"for", fmt.Sprintf("%+v", test), "\n",
				"Expected error:", test.ErrorMessage, "\n",
				"Got:", "nil", "\n",
			)
		}
	}
}

func TestParseConcentration(t *testing.T) {

	for _, test := range componentWithConcstests {
		a, b, c := ParseConcentration(test.ComponentName)
		if a != test.ContainsConc {
			t.Error(
				"for", fmt.Sprintf("%+v", test), "\n",
				"Expected:", test.ContainsConc, "\n",
				"Got:", a, "\n",
			)
		}
		if a && !b.EqualTo(test.Conc) {
			t.Error(
				"for", fmt.Sprintf("%+v", test), "\n",
				"Expected:", test.Conc, "\n",
				"Got:", b, "\n",
			)
		}
		if a && b.Unit().PrefixedSymbol() != test.Conc.Unit().PrefixedSymbol() {
			t.Error(
				"for", fmt.Sprintf("%+v", test), "\n",
				"Expected:", test.Conc, "\n",
				"Got:", b, "\n",
			)
		}
		if c != test.ComponentNameOnly {
			t.Error(
				"for", fmt.Sprintf("%+v", test), "\n",
				"Expected:", test.ComponentNameOnly, "\n",
				"Got:", c, "\n",
			)
		}
	}
}

type valueAndUnitTest struct {
	value        float64
	unit         string
	valueandunit string
}

var volandUnitTests = []valueAndUnitTest{
	{
		value:        10,
		unit:         "s",
		valueandunit: "10s",
	},
	{
		value:        10,
		unit:         "s",
		valueandunit: "10 s",
	},
	{
		value:        10,
		unit:         "",
		valueandunit: "10",
	},

	{
		value:        0,
		unit:         "s",
		valueandunit: "s",
	},
	{
		value:        10.9090,
		unit:         "ms",
		valueandunit: "10.9090ms",
	},
	{
		value:        2.16e+04,
		unit:         "s",
		valueandunit: "2.16e+04 s",
	},

	{
		value:        2.16e+04,
		unit:         "/s",
		valueandunit: "2.16e+04 /s",
	},
}

func TestSplitValueAndUnit(t *testing.T) {
	for _, test := range volandUnitTests {
		val, unit := SplitValueAndUnit(test.valueandunit)
		if val != test.value {
			t.Error(
				"for", fmt.Sprintf("%+v", test), "\n",
				"Expected:", test.value, "\n",
				"Got:", val, "\n",
			)
		}
		if unit != test.unit {
			t.Error(
				"for", fmt.Sprintf("%+v", test), "\n",
				"Expected:", test.unit, "\n",
				"Got:", unit, "\n",
			)
		}
	}
}
