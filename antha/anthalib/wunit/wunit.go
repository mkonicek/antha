// wunit/wunit.go: Part of the Antha language
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
	"math"

	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

// structure defining a base unit
type BaseUnit interface {
	// unit name
	Name() string
	// unit symbol
	Symbol() string
	// multiply by this to get SI value
	// nb this should be a function since we actually need
	// an affine transformation
	BaseSIConversionFactor() float64 // this can be calculated in many cases
	// if we convert to the SI units what is the appropriate unit symbol
	BaseSIUnit() string // if we use the above, what unit do we get?
	// print this
	ToString() string
}

// a unit with an SI prefix
type PrefixedUnit interface {
	BaseUnit
	// the prefix of the unit
	Prefix() SIPrefix
	// the symbol including prefix
	PrefixedSymbol() string
	// the symbol excluding prefix
	RawSymbol() string
	// appropriate unit if we ask for SI values
	BaseSISymbol() string
	// returns conversion factor from *this* unit to the other
	ConvertTo(pu PrefixedUnit) float64
	// CompatibleWith returns true iff the PrefixedUnit can be converted to rhs
	// if false then calling ConvertTo will cause a panic()
	CompatibleWith(PrefixedUnit) bool
}

// fundamental representation of a value in the system
type Measurement interface {
	// the value in base SI units
	SIValue() float64
	// the value in the current units
	RawValue() float64
	// unit plus prefix
	Unit() PrefixedUnit
	// set the value, this must be thread-safe
	// returns old value
	SetValue(v float64) float64
	// convert units
	ConvertTo(p PrefixedUnit) float64
	// wrapper for above
	ConvertToString(s string) float64
	// add to this measurement
	Add(m Measurement)
	// subtract from this measurement
	Subtract(m Measurement)
	// multiply measurement by a factor
	MultiplyBy(factor float64)
	// divide measurement by a factor
	DivideBy(factor float64)
	// comparison operators
	LessThan(m Measurement) bool
	GreaterThan(m Measurement) bool
	EqualTo(m Measurement) bool

	// A nice string representation
	ToString() string
}

// structure implementing the Measurement interface
type ConcreteMeasurement struct {
	// the raw value
	Mvalue float64
	// the relevant units
	Munit *Unit
}

/*
func AddMeasurements(a Measurement, b Measurement) (c Measurement) {
	if a.Unit().BaseSIUnit() == b.Unit().BaseSIUnit() {

		apointer := *a

		c = &apointer
		&c.Add(&b)
		/* *(CopyVolume(&A))
		(&C).Add(&B)
	}
	return c
}*/

func isNil(cm *ConcreteMeasurement) bool {
	if cm == nil {
		return true
	}
	if cm.Munit == nil {
		return true
	}
	return false
}

// value when converted to SI units
func (cm *ConcreteMeasurement) SIValue() float64 {
	if isNil(cm) {
		return 0.0
	}
	return cm.ConvertToString(cm.Unit().BaseSIUnit())
}

// value without conversion
func (cm *ConcreteMeasurement) RawValue() float64 {
	if isNil(cm) {
		return 0.0
	}
	return cm.Mvalue
}

// get unit with prefix
func (cm *ConcreteMeasurement) Unit() PrefixedUnit {
	if isNil(cm) {
		var ret *Unit
		return ret
	}
	return cm.Munit
}

// set the value of this measurement
func (cm *ConcreteMeasurement) SetValue(v float64) float64 {
	if isNil(cm) {
		return 0.0
	}
	cm.Mvalue = v
	return v
}

// convert to a different unit
// nb this is NOT destructive
func (cm *ConcreteMeasurement) ConvertTo(p PrefixedUnit) float64 {
	if isNil(cm) {
		return 0.0
	}
	return cm.Unit().ConvertTo(p) * cm.RawValue()
}

// ConvertToString return the measurement to the units specified by symbol.
// panic()s if symbol is not a known symbol
func (cm *ConcreteMeasurement) ConvertToString(s string) float64 {
	if isNil(cm) {
		return 0.0
	}
	reg := GetGlobalUnitRegistry()

	if unit, err := reg.GetUnit(s); err != nil {
		panic(err)
	} else {
		return cm.ConvertTo(unit)
	}
}

// String will return a summary of the ConcreteMeasurement Value and prefixed unit as a string.
// The value will be formatted in scientific notation for large exponents and the value unbounded.
// The Summary() method should be used to return a rounded string.
func (cm *ConcreteMeasurement) String() string {
	return fmt.Sprintf("%g %s", cm.RawValue(), cm.Unit().PrefixedSymbol())
}

// add to this

func (cm *ConcreteMeasurement) Add(m Measurement) {
	if isNil(cm) {
		return
	}
	// ideally should check these have the same Dimension
	// need to improve this

	cm.SetValue(m.ConvertTo(cm.Unit()) + cm.RawValue())

}

// subtract

func (cm *ConcreteMeasurement) Subtract(m Measurement) {
	if isNil(cm) {
		return
	}
	// ideally should check these have the same Dimension
	// need to improve this

	cm.SetValue(cm.RawValue() - m.ConvertTo(cm.Unit()))

}

// multiply
func (cm *ConcreteMeasurement) MultiplyBy(factor float64) {
	if isNil(cm) {
		return
	}
	// ideally should check these have the same Dimension
	// need to improve this

	cm.SetValue(cm.RawValue() * float64(factor))

}

func (cm *ConcreteMeasurement) DivideBy(factor float64) {

	if isNil(cm) {
		return
	}
	// ideally should check these have the same Dimension
	// need to improve this

	cm.SetValue(cm.RawValue() / float64(factor))

}

// define a zero

func (cm *ConcreteMeasurement) IsNil() bool {
	return isNil(cm)
}

func (cm *ConcreteMeasurement) IsZero() bool {
	if isNil(cm) || cm.Mvalue < 0.00000000001 {
		return true
	}
	return false
}

// less sensitive comparison operators

func (cm *ConcreteMeasurement) LessThanRounded(m Measurement, p int) bool {
	// nil means less than everything
	if isNil(cm) {
		return true
	}
	// returns true if this is less than m
	v := wutil.RoundIgnoreNan(m.ConvertTo(cm.Unit()), p)
	v2 := wutil.RoundIgnoreNan(cm.RawValue(), p)

	return v > v2
}

func (cm *ConcreteMeasurement) GreaterThanRounded(m Measurement, p int) bool {
	if isNil(cm) {
		return false
	}
	// returns true if this is greater than m
	v := wutil.RoundIgnoreNan(m.ConvertTo(cm.Unit()), p)
	v2 := wutil.RoundIgnoreNan(cm.RawValue(), p)
	return v < v2

}

func (cm *ConcreteMeasurement) EqualToRounded(m Measurement, p int) bool {
	// this is not equal to anything
	if isNil(cm) {
		return false
	}

	// returns true if this is equal to m
	v := wutil.RoundIgnoreNan(m.ConvertTo(cm.Unit()), p)
	v2 := wutil.RoundIgnoreNan(cm.RawValue(), p)

	return v == v2
}

// comparison operators

func (cm *ConcreteMeasurement) LessThan(m Measurement) bool {
	// nil means less than everything
	if isNil(cm) {
		return true
	}
	// returns true if this is less than m
	v := m.ConvertTo(cm.Unit())

	return v > cm.RawValue()
}

func (cm *ConcreteMeasurement) GreaterThan(m Measurement) bool {
	if isNil(cm) {
		return false
	}
	// returns true if this is greater than m
	v := m.ConvertTo(cm.Unit())
	return v < cm.RawValue()
}

// XXX This should be made more literal and rounded behaviour explicitly called for by user
func (cm *ConcreteMeasurement) EqualTo(m Measurement) bool {
	// this is not equal to anything

	if isNil(cm) {
		return false
	}
	// returns true if this is equal to m
	v := m.ConvertTo(cm.Unit())

	dif := math.Abs(v - cm.RawValue())

	epsilon := math.Nextafter(1, 2) - 1
	return dif < (epsilon * 10000)
}

// EqualToTolerance return true if m is within a small tolerace, tol, of the measurement
// where tol is expressed in the same units are the receiver
func (cm *ConcreteMeasurement) EqualToTolerance(m Measurement, tol float64) bool {
	if isNil(cm) {
		return false
	}

	return math.Abs(m.ConvertTo(cm.Unit())-cm.RawValue()) < tol
}

// ToString will return a summary of the ConcreteMeasurement Value and prefixed unit as a string.
// The value will be formatted in scientific notation for large exponents and will be bounded to 3 decimal places.
// The String() method should be used to use the unbounded value.
func (cm *ConcreteMeasurement) ToString() string {
	return fmt.Sprintf("%.3g %s", cm.RawValue(), cm.Unit().PrefixedSymbol())
}

/**********/

func NewMeasurement(v float64, pu string) *ConcreteMeasurement {
	if value, err := GetGlobalUnitRegistry().NewMeasurement(v, pu); err != nil {
		panic(err)
	} else {
		return value
	}
}
