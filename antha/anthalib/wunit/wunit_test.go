// wunit/wunit_test.go: Part of the Antha language
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
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

type VolumeArithmetic struct {
	VolumeA    Volume
	VolumeB    Volume
	Sum        Volume
	Difference Volume
	Factor     float64
	Product    Volume
	Quotient   Volume
}

var volumearithmetictests = []VolumeArithmetic{
	{
		VolumeA:    NewVolume(1, "ul"),
		VolumeB:    NewVolume(1, "ul"),
		Sum:        NewVolume(2, "ul"),
		Difference: NewVolume(0, "ul"),
		Factor:     1.0,
		Product:    NewVolume(1, "ul"),
		Quotient:   NewVolume(1, "ul"),
	},
	{
		VolumeA:    NewVolume(100, "ul"),
		VolumeB:    NewVolume(10, "ul"),
		Sum:        NewVolume(110, "ul"),
		Difference: NewVolume(90, "ul"),
		Factor:     10.0,
		Product:    NewVolume(1000, "ul"),
		Quotient:   NewVolume(10, "ul"),
	},
	{
		VolumeA:    NewVolume(1000000, "ul"),
		VolumeB:    NewVolume(10, "ul"),
		Sum:        NewVolume(1000010, "ul"),
		Difference: NewVolume(999990, "ul"),
		Factor:     10.0,
		Product:    NewVolume(10000000, "ul"),
		Quotient:   NewVolume(100000, "ul"),
	},
	{
		VolumeA:    NewVolume(1, "l"),
		VolumeB:    NewVolume(10, "ul"),
		Sum:        NewVolume(1000010, "ul"),
		Difference: NewVolume(999990, "ul"),
		Factor:     10.0,
		Product:    NewVolume(10000000, "ul"),
		Quotient:   NewVolume(100000, "ul"),
	},
	{
		VolumeA:    NewVolume(1000, "ml"),
		VolumeB:    NewVolume(10, "ul"),
		Sum:        NewVolume(1000010, "ul"),
		Difference: NewVolume(999990, "ul"),
		Factor:     10.0,
		Product:    NewVolume(10000000, "ul"),
		Quotient:   NewVolume(100000, "ul"),
	},
	{
		VolumeA:    NewVolume(1000, "ul"),
		VolumeB:    NewVolume(-10, "ul"),
		Sum:        NewVolume(990, "ul"),
		Difference: NewVolume(1010, "ul"),
		Factor:     -10.0,
		Product:    NewVolume(-10000, "ul"),
		Quotient:   NewVolume(-100, "ul"),
	},
	{
		VolumeA:    NewVolume(-1000, "ul"),
		VolumeB:    NewVolume(10, "ul"),
		Sum:        NewVolume(-990, "ul"),
		Difference: NewVolume(-1010, "ul"),
		Factor:     -10.0,
		Product:    NewVolume(10000, "ul"),
		Quotient:   NewVolume(100, "ul"),
	},
	{
		VolumeA:    NewVolume(100, "ul"),
		VolumeB:    NewVolume(-165, "ul"),
		Sum:        NewVolume(-65, "ul"),
		Difference: NewVolume(265, "ul"),
		Factor:     10.0,
		Product:    NewVolume(1000, "ul"),
		Quotient:   NewVolume(10, "ul"),
	},
}

func TestSubstractVolumes(t *testing.T) {
	for _, testUnit := range volumearithmetictests {
		r := SubtractVolumes(testUnit.VolumeA, testUnit.VolumeB)
		rt, _ := wutil.Roundto(r.SIValue(), 4)
		tt, _ := wutil.Roundto(testUnit.Difference.SIValue(), 4)
		if rt != tt {
			t.Error(
				"For", testUnit.VolumeA, "-", testUnit.VolumeB, "\n",
				"expected", testUnit.Difference, "\n",
				"got", r, "\n",
			)
		}
	}

}

func TestAddVolumes(t *testing.T) {
	for _, testUnit := range volumearithmetictests {
		r := AddVolumes(testUnit.VolumeA, testUnit.VolumeB)
		if r.SIValue() != testUnit.Sum.SIValue() {
			t.Error(
				"For", testUnit.VolumeA, "+", testUnit.VolumeB, "\n",
				"expected", testUnit.Sum, "\n",
				"got", r, "\n",
			)
		}
	}

}

func TestMultiplyVolumes(t *testing.T) {
	for _, testUnit := range volumearithmetictests {
		r := MultiplyVolume(testUnit.VolumeA, testUnit.Factor)
		if r.SIValue() != testUnit.Product.SIValue() {
			t.Error(
				"For", testUnit.VolumeA, " x ", testUnit.Factor, "\n",
				"expected", testUnit.Product, "\n",
				"got", r, "\n",
			)
		}
	}

}

func TestDivideVolume(t *testing.T) {
	for _, testUnit := range volumearithmetictests {
		r := DivideVolume(testUnit.VolumeA, testUnit.Factor)
		rt, _ := wutil.Roundto(r.SIValue(), 4)
		tt, _ := wutil.Roundto(testUnit.Quotient.SIValue(), 4)
		if rt != tt {
			t.Error(
				"For", testUnit.VolumeA, " / ", testUnit.Factor, "\n",
				"expected", testUnit.Quotient, "\n",
				"got", r, "\n",
			)
		}
	}

}

func TestDivideVolumes(t *testing.T) {
	for _, testUnit := range volumearithmetictests {
		r, err := DivideVolumes(testUnit.Product, testUnit.VolumeA)
		if err != nil {
			t.Errorf("For DivideVolumes(%v, %v): got error: %v", testUnit.Product, testUnit.VolumeA, err)
			continue
		}
		rt, _ := wutil.Roundto(r, 4)
		tt, _ := wutil.Roundto(testUnit.Factor, 4)
		if rt != tt {
			t.Error(
				"For DivideVolumes(", testUnit.Product, ", ", testUnit.VolumeA, ")\n",
				"expected", testUnit.Factor, "\n",
				"got", r, "\n",
			)
		}
	}

	if _, err := DivideVolumes(ZeroVolume(), ZeroVolume()); err == nil {
		t.Error("divide by zero didn't cause error")
	}

}

type ConcArithmetic struct {
	ValueA     Concentration
	ValueB     Concentration
	Sum        Concentration
	Difference Concentration
	Factor     float64
	Product    Concentration
	Quotient   Concentration
}

var concarithmetictests = []ConcArithmetic{
	{
		ValueA:     NewConcentration(1, "ng/ul"),
		ValueB:     NewConcentration(1, "ng/ul"),
		Sum:        NewConcentration(2, "ng/ul"),
		Difference: NewConcentration(0, "ng/ul"),
		Factor:     1.0,
		Product:    NewConcentration(1, "ng/ul"),
		Quotient:   NewConcentration(1, "ng/ul"),
	},
	{
		ValueA:     NewConcentration(100, "ng/ul"),
		ValueB:     NewConcentration(10, "ng/ul"),
		Sum:        NewConcentration(110, "ng/ul"),
		Difference: NewConcentration(90, "ng/ul"),
		Factor:     10.0,
		Product:    NewConcentration(1000, "ng/ul"),
		Quotient:   NewConcentration(10, "ng/ul"),
	},
	{
		ValueA:     NewConcentration(1000000, "mg/l"),
		ValueB:     NewConcentration(10, "ng/ul"),
		Sum:        NewConcentration(1000010, "ng/ul"),
		Difference: NewConcentration(999990, "ng/ul"),
		Factor:     10.0,
		Product:    NewConcentration(10000000, "ng/ul"),
		Quotient:   NewConcentration(100000, "ng/ul"),
	},
	{
		ValueA:     NewConcentration(1000, "g/l"),
		ValueB:     NewConcentration(10, "ng/ul"),
		Sum:        NewConcentration(1000010, "ng/ul"),
		Difference: NewConcentration(999.99, "g/l"),
		Factor:     10.0,
		Product:    NewConcentration(10000000, "ng/ul"),
		Quotient:   NewConcentration(100, "g/l"),
	},
	{
		ValueA:     NewConcentration(1, "Mol/l"),
		ValueB:     NewConcentration(10, "mMol/l"),
		Sum:        NewConcentration(1.01, "Mol/l"),
		Difference: NewConcentration(0.99, "Mol/l"),
		Factor:     10.0,
		Product:    NewConcentration(10, "Mol/l"),
		Quotient:   NewConcentration(0.1, "Mol/l"),
	},
	{
		ValueA:     NewConcentration(2, "ng/ul"),
		ValueB:     NewConcentration(1, "ng/ul"),
		Sum:        NewConcentration(3, "ng/ul"),
		Difference: NewConcentration(1, "ng/ul"),
		Factor:     2.0,
		Product:    NewConcentration(4, "ng/ul"),
		Quotient:   NewConcentration(1, "ng/ul"),
	},
}

func TestMultiplyConcentrations(t *testing.T) {
	for _, testUnit := range concarithmetictests {
		r := MultiplyConcentration(testUnit.ValueA, testUnit.Factor)
		if r.SIValue() != testUnit.Product.SIValue() {
			t.Error(
				"For", testUnit.ValueA, "\n",
				"expected", testUnit.Product, "\n",
				"got", r, "\n",
			)
		}
	}

}

func TestDivideConcentration(t *testing.T) {
	for _, testUnit := range concarithmetictests {
		r := DivideConcentration(testUnit.ValueA, testUnit.Factor)
		if r.SIValue() != testUnit.Quotient.SIValue() {
			t.Error(
				"For", testUnit.ValueA, "\n",
				"expected", testUnit.Quotient, "\n",
				"got", r, "\n",
			)
		}
	}

}

func TestAddConcentrations(t *testing.T) {
	for _, testUnit := range concarithmetictests {
		r, err := AddConcentrations(testUnit.ValueA, testUnit.ValueB)
		if err != nil {
			t.Error(
				"Add Concentration returns error ", err.Error(), "should return nil \n",
			)
		}
		if r.SIValue() != testUnit.Sum.SIValue() {
			t.Error(
				"For addition of ", testUnit.ValueA, "and", testUnit.ValueB, "\n",
				"expected", testUnit.Sum, "\n",
				"got", r, "\n",
			)
		}
	}

	_, err := AddConcentrations(concarithmetictests[0].ValueA, concarithmetictests[4].ValueA)
	if err == nil {
		t.Error(
			"Expected Errorf but got nil. Adding of two different bases (g/l and M/l) should not be possible \n",
		)
	}

}

func TestSubtractConcentrations(t *testing.T) {
	for _, testUnit := range concarithmetictests {
		r, err := SubtractConcentrations(testUnit.ValueA, testUnit.ValueB)
		if err != nil {
			t.Error(
				"Subtract Concentration returns error ", err.Error(), "should return nil \n",
			)
		}
		if r.SIValue() != testUnit.Difference.SIValue() {
			t.Error(
				"For subtraction of ", testUnit.ValueB, " from ", testUnit.ValueA, "\n",
				"expected", testUnit.Difference, "\n",
				"got", r, "\n",
			)
		}
	}

	_, err := SubtractConcentrations(concarithmetictests[0].ValueA, concarithmetictests[4].ValueA)
	if err == nil {
		t.Error(
			"Expected Errorf but got nil. Subtracting of two different bases (g/l and M/l) should not be possible \n",
		)
	}

}

// test precision
func TestDivideConcentrationsPrecision(t *testing.T) {

	type divideTest struct {
		StockConc, TargetConc Concentration
		ExpectedFactor        float64
		ExpectedErr           error
	}

	tests := []divideTest{
		{
			StockConc:      NewConcentration(15, "X"),
			TargetConc:     NewConcentration(7.5, "X"),
			ExpectedFactor: 2.0000000000000000000,
		},
	}

	for _, test := range tests {
		r, err := DivideConcentrations(test.StockConc, test.TargetConc)

		if err != test.ExpectedErr {
			t.Error("expected: ", err, "\n",
				"got: ", test.ExpectedErr,
			)
		}

		if r != test.ExpectedFactor {
			t.Error(
				"For", fmt.Sprintf("+%v", test), "\n",
				"expected factor: ", test.ExpectedFactor, "\n",
				"got", r, "\n",
			)
		}
	}

}

// test precision
func TestDivideVolumePrecision(t *testing.T) {

	type divideTest struct {
		StockVolume, ExpectedVolume Volume
		Factor                      float64
		ExpectedErr                 error
	}

	tests := []divideTest{
		{
			StockVolume:    NewVolume(100, "ul"),
			ExpectedVolume: NewVolume(50.0, "ul"),
			Factor:         2.0000000000000000000,
		},
	}

	for _, test := range tests {
		r := DivideVolume(test.StockVolume, test.Factor)
		if !r.EqualTo(test.ExpectedVolume) {
			t.Error(
				"For", fmt.Sprintf("+%v", test), "\n",
				"expected: ", test.ExpectedVolume, "\n",
				"got", r, "\n",
			)
		}
	}

}

// test precision
func TestDivideConcentrationPrecision(t *testing.T) {

	type divideTest struct {
		StockConcentration, ExpectedConcentration Concentration
		Factor                                    float64
		ExpectedErr                               error
	}

	tests := []divideTest{
		{
			StockConcentration:    NewConcentration(0.00012207, "X"),
			ExpectedConcentration: NewConcentration(6.1035e-05, "X"),
			Factor:                2.0,
		},
		{
			StockConcentration:    NewConcentration(0.000125, "X"),
			ExpectedConcentration: NewConcentration(0.0000625, "X"),
			Factor:                2.0000000000000000000,
		},
		{
			StockConcentration:    NewConcentration(0.0625, "X"),
			ExpectedConcentration: NewConcentration(0.03125, "X"),
			Factor:                2.0000000000000000000,
		},

		{
			StockConcentration:    NewConcentration(22.0/7.0, "X"),
			ExpectedConcentration: NewConcentration(3.14285714285714, "X"),
			Factor:                1.0000000000000000000,
		},
	}

	for _, test := range tests {
		r := DivideConcentration(test.StockConcentration, test.Factor)
		if !r.EqualTo(test.ExpectedConcentration) {
			t.Error(
				"For", fmt.Sprintf("+%v", test), "\n",
				"expected: ", test.ExpectedConcentration, "\n",
				"got", r, "\n",
			)
		}
	}

}

func TestFlowRateComparison(t *testing.T) {
	a := NewFlowRate(1., "ml/min")
	b := NewFlowRate(2., "ml/min")

	if !b.GreaterThan(a) {
		t.Errorf("Got b > a (%s > %s) wrong", b, a)
	}
	if a.GreaterThan(b) {
		t.Errorf("Got a > b (%s > %s) wrong", a, b)
	}

	if b.LessThan(a) {
		t.Errorf("Got b < a (%s < %s) wrong", b, a)
	}
	if !a.LessThan(b) {
		t.Errorf("Got a < b (%s < %s) wrong", a, b)
	}
}

func TestRoundedComparisons(t *testing.T) {
	v1 := NewVolume(0.5, "ul")
	v2 := NewVolume(0.4999999, "ul")

	vrai := v1.GreaterThanRounded(v2, 7)

	if !vrai {
		t.Error(
			"For", v1.ToString(), " >_7 ", v2.ToString(), "\n",
			"expected true\n",
			"got false\n",
		)
	}

	faux := v1.LessThanRounded(v2, 7)

	if faux {
		t.Error(
			"For", v1.ToString(), " <_7 ", v2.ToString(), "\n",
			"expected false\n",
			"got true\n",
		)
	}

	faux = v1.EqualToRounded(v2, 8)

	if faux {
		t.Error(
			"For", v1.ToString(), " ==_7 ", v2.ToString(), "\n",
			"expected false\n",
			"got true\n",
		)

	}

	vrai = v1.EqualToRounded(v2, 6)

	if !vrai {
		t.Error(
			"For", v1.ToString(), " ==_6 ", v2.ToString(), "\n",
			"expected true\n",
			"got false\n",
		)

	}

	faux = v1.LessThanRounded(v2, 6)

	if faux {
		t.Error(
			"For", v1.ToString(), " <_6 ", v2.ToString(), "\n",
			"expected false\n",
			"got true\n",
		)
	}

	faux = v1.GreaterThanRounded(v2, 6)

	if faux {
		t.Error(
			"For", v1.ToString(), " >_6 ", v2.ToString(), "\n",
			"expected false\n",
			"got true\n",
		)

	}
}
