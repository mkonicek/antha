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
	"testing"
)

type unitFromStringTest struct {
	ComponentName     string
	ContainsConc      bool
	Conc              Concentration
	ComponentNameOnly string
}

func (test *unitFromStringTest) Run(t *testing.T) {
	t.Run(test.ComponentName, func(t *testing.T) {
		a, b, c := ParseConcentration(test.ComponentName)
		if a != test.ContainsConc {
			t.Errorf("ContainsConc: expected %t, got %t", test.ContainsConc, a)
		}
		if a && !b.EqualTo(test.Conc) {
			t.Errorf("parsed concentration incorrect: expected %v, got %v", test.Conc, b)
		}
		if c != test.ComponentNameOnly {
			t.Errorf("ComponentNameOnly: expected %q, got %q", test.ComponentNameOnly, c)
		}
	})
}

type unitFromStringTests []unitFromStringTest

func (self unitFromStringTests) Run(t *testing.T) {
	for _, test := range self {
		test.Run(t)
	}
}

func TestParseConcentration(t *testing.T) {
	unitFromStringTests{
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
			ComponentName:     "1 % w/v C6",
			ContainsConc:      true,
			Conc:              NewConcentration(1.0, "% w/v"),
			ComponentNameOnly: "C6",
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
			ComponentName:     "A happy solution :-) 5 g/l",
			ContainsConc:      true,
			Conc:              NewConcentration(5.0, "g/l"),
			ComponentNameOnly: "A happy solution :-)",
		},
		{
			ComponentName:     "A sad solution :-( 5 g/l",
			ContainsConc:      true,
			Conc:              NewConcentration(5.0, "g/l"),
			ComponentNameOnly: "A sad solution :-(",
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
			ComponentName:     "1 MOPS pH 7",
			ContainsConc:      false,
			ComponentNameOnly: "1 MOPS pH 7",
		},
		{
			ComponentName:     "1 MOPS pH 7 (mM)",
			ContainsConc:      true,
			Conc:              NewConcentration(0.0, "mM"),
			ComponentNameOnly: "1 MOPS pH 7",
		},
		{
			ComponentName:     "281 mMol/l 1 MOPS pH 7",
			ContainsConc:      true,
			Conc:              NewConcentration(281.0, "mM"),
			ComponentNameOnly: "1 MOPS pH 7",
		},
		{
			ComponentName:     "(1X)",
			ContainsConc:      true,
			Conc:              NewConcentration(1, "X"),
			ComponentNameOnly: "(1X)",
		},
		{
			ComponentName:     "Magnesium Acetate",
			ContainsConc:      false,
			ComponentNameOnly: "Magnesium Acetate",
		},
		{
			ComponentName:     "0.515 v/v (S)-Styrene",
			ContainsConc:      true,
			ComponentNameOnly: "(S)-Styrene",
			Conc:              NewConcentration(0.515, "v/v"),
		},
		{
			ComponentName:     "(D)Glucose (6M)",
			ContainsConc:      true,
			ComponentNameOnly: "(D)Glucose",
			Conc:              NewConcentration(6.0, "M"),
		},

		{
			ComponentName:     "1X Antigen X",
			ContainsConc:      true,
			ComponentNameOnly: "Antigen X",
			Conc:              NewConcentration(1.0, "X"),
		},
	}.Run(t)
}

type volTest struct {
	VolString   string
	Volume      Volume
	ShouldError bool
}

func (test *volTest) unexpectedError(err error) bool {
	return (err != nil) != test.ShouldError
}

func (test *volTest) Run(t *testing.T) {
	t.Run(test.VolString, func(t *testing.T) {
		vol, err := ParseVolume(test.VolString)
		if !test.ShouldError && !vol.EqualTo(test.Volume) {
			t.Errorf("parsed volume incorrect: expected %v, got %v", test.Volume, vol)
		}
		if test.unexpectedError(err) {
			t.Errorf("error mismatched: expected %t, got error %v", test.ShouldError, err)
		}
	})
}

type volTests []volTest

func (self volTests) Run(t *testing.T) {
	for _, test := range self {
		test.Run(t)
	}
}

func TestParseVolume(t *testing.T) {
	volTests{
		{
			VolString: "10ul",
			Volume:    NewVolume(10, "ul"),
		},
		{
			VolString: "10 ul",
			Volume:    NewVolume(10, "ul"),
		},
		{
			VolString:   "10",
			Volume:      Volume{},
			ShouldError: true,
		},
	}.Run(t)
}

type valueAndUnitTest struct {
	value        float64
	unit         string
	valueandunit string
}

func (test *valueAndUnitTest) Run(t *testing.T) {
	t.Run(test.valueandunit, func(t *testing.T) {
		if val, unit := SplitValueAndUnit(test.valueandunit); val != test.value || unit != test.unit {
			t.Errorf("expected = %f, %q; got = %f, %q;", test.value, test.unit, val, unit)
		}
	})
}

type valueAndUnitTests []valueAndUnitTest

func (tests valueAndUnitTests) Run(t *testing.T) {
	for _, test := range tests {
		test.Run(t)
	}
}

func TestSplitValueAndUnit(t *testing.T) {
	valueAndUnitTests{
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
	}.Run(t)
}
