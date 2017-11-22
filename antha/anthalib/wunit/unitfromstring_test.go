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
	unitFromStringTest{
		ComponentName:     "Glucose (M)",
		ContainsConc:      true,
		Conc:              NewConcentration(0.0, "M"),
		ComponentNameOnly: "Glucose",
	},
	unitFromStringTest{
		ComponentName:     "Glucose (mM)",
		ContainsConc:      true,
		Conc:              NewConcentration(0.0, "mM"),
		ComponentNameOnly: "Glucose",
	},
	unitFromStringTest{
		ComponentName:     "5g/L Glucose",
		ContainsConc:      true,
		Conc:              NewConcentration(5.0, "g/L"),
		ComponentNameOnly: "Glucose",
	},
	unitFromStringTest{
		ComponentName:     "5 g/L Glucose",
		ContainsConc:      true,
		Conc:              NewConcentration(5.0, "g/L"),
		ComponentNameOnly: "Glucose",
	},
	unitFromStringTest{
		ComponentName:     "Glucose 5g/L",
		ContainsConc:      true,
		Conc:              NewConcentration(5.0, "g/L"),
		ComponentNameOnly: "Glucose",
	},
	unitFromStringTest{
		ComponentName:     "1mM/l C6",
		ContainsConc:      true,
		Conc:              NewConcentration(1.0, "mM/l"),
		ComponentNameOnly: "C6",
	},
	unitFromStringTest{
		ComponentName:     "C6",
		ContainsConc:      false,
		Conc:              NewConcentration(0.0, "mM"),
		ComponentNameOnly: "C6",
	},
	unitFromStringTest{
		ComponentName:     "1 mM Ammonium Sulphate",
		ContainsConc:      true,
		Conc:              NewConcentration(1.0, "mM"),
		ComponentNameOnly: "Ammonium Sulphate",
	},
	unitFromStringTest{
		ComponentName:     "Ammonium Sulphate (mM)",
		ContainsConc:      true,
		Conc:              NewConcentration(0.0, "mM"),
		ComponentNameOnly: "Ammonium Sulphate",
	},
	unitFromStringTest{
		ComponentName:     "E.coli SuperFolder GFP (g/L)",
		ContainsConc:      true,
		Conc:              NewConcentration(0.0, "g/L"),
		ComponentNameOnly: "E.coli SuperFolder GFP",
	},
}

type volTest struct {
	VolString    string
	Volume       Volume
	ErrorMessage string
}

var volTests = []volTest{
	volTest{
		VolString:    "10ul",
		Volume:       NewVolume(10, "ul"),
		ErrorMessage: "",
	},
	volTest{
		VolString:    "10 ul",
		Volume:       NewVolume(10, "ul"),
		ErrorMessage: "",
	},
	volTest{
		VolString:    "10",
		Volume:       Volume{},
		ErrorMessage: "no valid unit found for 10: valid units are: [L l ml nl ul]",
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
