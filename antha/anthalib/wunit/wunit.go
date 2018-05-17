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
	// add to this measurement - synonym for less typing
	A(m Measurement)
	// subtract from this measurement
	Subtract(m Measurement)
	// subtract - synonym for less typing
	S(m Measurement)
	// multiply measurement by a factor
	MultiplyBy(factor float64)
	// multiply by a factor - synonym for less typing
	M(factor float64)
	// divide measurement by a factor
	DivideBy(factor float64)
	// divide - synonym for less typing
	D(factor float64)
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
	Munit *GenericPrefixedUnit
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

// value when converted to SI units
func (cm *ConcreteMeasurement) SIValue() float64 {
	if cm == nil {
		return 0.0
	}
	return cm.Mvalue * cm.Munit.BaseSIConversionFactor()
}

// value without conversion
func (cm *ConcreteMeasurement) RawValue() float64 {
	if cm == nil {
		return 0.0
	}
	return cm.Mvalue
}

// get unit with prefix
func (cm *ConcreteMeasurement) Unit() PrefixedUnit {
	if cm == nil {
		return NewGenericPrefixedUnit()
	}
	return cm.Munit
}

// set the value of this measurement
func (cm *ConcreteMeasurement) SetValue(v float64) float64 {
	if cm == nil {
		return 0.0
	}
	cm.Mvalue = v
	return v
}

// convert to a different unit
// nb this is NOT destructive
func (cm *ConcreteMeasurement) ConvertTo(p PrefixedUnit) float64 {
	if cm == nil {
		return 0.0
	}
	return cm.Unit().ConvertTo(p) * cm.RawValue()
}

func (cm *ConcreteMeasurement) ConvertToString(s string) float64 {
	if cm == nil {
		return 0.0
	}
	ppu := ParsePrefixedUnit(s)
	return cm.ConvertTo(ppu)
}

func (cm *ConcreteMeasurement) String() string {
	return cm.ToString()
}

// add to this

func (cm *ConcreteMeasurement) Add(m Measurement) {
	if cm == nil {
		return
	}
	// ideally should check these have the same Dimension
	// need to improve this

	cm.SetValue(m.ConvertTo(cm.Unit()) + cm.RawValue())

}

func (cm *ConcreteMeasurement) A(m Measurement) {
	cm.Add(m)
}

// subtract

func (cm *ConcreteMeasurement) Subtract(m Measurement) {
	if cm == nil {
		return
	}
	// ideally should check these have the same Dimension
	// need to improve this

	cm.SetValue(cm.RawValue() - m.ConvertTo(cm.Unit()))

}
func (cm *ConcreteMeasurement) S(m Measurement) {
	cm.Subtract(m)
}

// multiply
func (cm *ConcreteMeasurement) MultiplyBy(factor float64) {
	if cm == nil {
		return
	}
	// ideally should check these have the same Dimension
	// need to improve this

	cm.SetValue(cm.RawValue() * float64(factor))

}

func (cm *ConcreteMeasurement) M(factor float64) {
	cm.MultiplyBy(factor)
}

func (cm *ConcreteMeasurement) DivideBy(factor float64) {

	if cm == nil {
		return
	}
	// ideally should check these have the same Dimension
	// need to improve this

	cm.SetValue(cm.RawValue() / float64(factor))

}

func (cm *ConcreteMeasurement) D(factor float64) {
	cm.DivideBy(factor)
}

// define a zero

func (cm *ConcreteMeasurement) IsNil() bool {
	if cm == nil || cm.Munit == nil {
		return true
	}

	return false
}

func (cm *ConcreteMeasurement) IsZero() bool {
	if cm == nil || cm.Munit == nil || cm.Mvalue < 0.00000000001 {
		return true
	}
	return false
}

// less sensitive comparison operators

func (cm *ConcreteMeasurement) LessThanRounded(m Measurement, p int) bool {
	// nil means less than everything
	if cm == nil {
		return true
	}
	// returns true if this is less than m
	v := wutil.RoundIgnoreNan(m.ConvertTo(cm.Unit()), p)
	v2 := wutil.RoundIgnoreNan(cm.RawValue(), p)

	return v > v2
}

func (cm *ConcreteMeasurement) GreaterThanRounded(m Measurement, p int) bool {
	if cm == nil {
		return false
	}
	// returns true if this is greater than m
	v := wutil.RoundIgnoreNan(m.ConvertTo(cm.Unit()), p)
	v2 := wutil.RoundIgnoreNan(cm.RawValue(), p)
	return v < v2

}

func (cm *ConcreteMeasurement) EqualToRounded(m Measurement, p int) bool {
	// this is not equal to anything
	if cm == nil {
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
	if cm == nil {
		return true
	}
	// returns true if this is less than m
	v := m.ConvertTo(cm.Unit())

	return v > cm.RawValue()
}

func (cm *ConcreteMeasurement) LessThanFloat(f float64) bool {
	if cm == nil {
		return true
	}
	// assumes the units work out

	return cm.RawValue() < f
}

func (cm *ConcreteMeasurement) GreaterThan(m Measurement) bool {
	if cm == nil {
		return false
	}
	// returns true if this is greater than m
	v := m.ConvertTo(cm.Unit())
	return v < cm.RawValue()
}

func (cm *ConcreteMeasurement) GreaterThanFloat(f float64) bool {
	if cm == nil {
		return false
	}
	if cm.RawValue() > f {
		return true
	}

	return false
}

// XXX This should be made more literal and rounded behaviour explicitly called for by user
func (cm *ConcreteMeasurement) EqualTo(m Measurement) bool {
	// this is not equal to anything

	if cm == nil {
		return false
	}
	// returns true if this is equal to m
	v := m.ConvertTo(cm.Unit())

	dif := math.Abs(v - cm.RawValue())

	epsilon := math.Nextafter(1, 2) - 1
	return dif < (epsilon * 10000)
}

func (cm *ConcreteMeasurement) EqualToFloat(f float64) bool {
	if cm == nil {
		return false
	}

	return f == cm.RawValue()
}

func (cm *ConcreteMeasurement) ToString() string {
	return fmt.Sprintf("%g %s", cm.RawValue(), cm.Unit().PrefixedSymbol())
}

/**********/

func NewPMeasurement(v float64, pu string) *ConcreteMeasurement {
	cm := ConcreteMeasurement{v, ParsePrefixedUnit(pu)}
	return &cm
}

// helper function for creating a new measurement
func NewMeasurement(v float64, prefix string, unit string) *ConcreteMeasurement {
	gpu := NewPrefixedUnit(prefix, unit)
	cm := ConcreteMeasurement{v, gpu}
	return &cm
}
