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
	"github.com/pkg/errors"
	"math"

	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

// PrefixedUnit a unit with an SI prefix
type PrefixedUnit interface {
	// Name get the full name of the unit
	Name() string
	// String return a string including the long name of the unit, its prefix, and symbol, e.g. "miligrams[mg]"
	String() string
	// Prefix get the SI prefix associated with the unit, or None if none, e.g. Mili
	Prefix() SIPrefix
	// PrefixedSymbol get the symbol including any prefix, e.g. "mg"
	PrefixedSymbol() string
	// RawSymbol the unit symbol excluding any prefix, e.g. "g"
	RawSymbol() string
	// BaseSISymbol returns the symbol of the appropriate unit if we ask for SI values, e.g. "kg"
	BaseSISymbol() string
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
	// InUnit get a new Measurement with the new units, returns error if units are not compatible
	InUnit(p PrefixedUnit) (Measurement, error)
	// InStringUnit wrapper for InUnit which fetches the unit from the global UnitRegistry
	InStringUnit(symbol string) (Measurement, error)
	// MustInUnit get a new Measurement with the new units, equivalent to InUnit except calls panic() if units are not compatible
	MustInUnit(p PrefixedUnit) Measurement
	// MustInStringUnit wrapper for InUnit which fetches the unit from the global UnitRegistry,
	// equivalent to InStringUnit but calls panic() if units are incompatible
	MustInStringUnit(symbol string) Measurement
	// ConvertToString deprecated, please use ConvertTo or InStringUnit
	ConvertToString(s string) float64
	// IncrBy add to this measurement
	IncrBy(m Measurement) error
	// DecrBy subtract from this measurement
	DecrBy(m Measurement) error
	// Add deprecated, please use IncrBy
	Add(m Measurement)
	// Subtract deprecated, please use DecrBy
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
	return cm.ConvertToString(cm.Unit().BaseSISymbol())
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

// ConvertTo return a new measurement in the new units
func (cm *ConcreteMeasurement) InUnit(p PrefixedUnit) (Measurement, error) {
	if isNil(cm) {
		return &ConcreteMeasurement{}, nil
	} else if rhs, ok := p.(*Unit); !ok { //since we currently don't have any methods in PrefixedUnit for unit conversion
		return nil, errors.Errorf("unsupported PrefixedUnit type %T", p)
	} else if factor, err := cm.Munit.getConversionFactor(rhs); err != nil {
		return nil, err
	} else if unit, ok := p.(*Unit); !ok {
		return nil, errors.Errorf("cannot convert unit type %T to *Unit", unit)
	} else {
		return &ConcreteMeasurement{Mvalue: factor * cm.RawValue(), Munit: unit}, nil
	}
}

// InStringUnit return a new measurement in the new units
func (cm *ConcreteMeasurement) InStringUnit(symbol string) (Measurement, error) {
	if unit, err := GetGlobalUnitRegistry().GetUnit(symbol); err != nil {
		return nil, err
	} else {
		return cm.InUnit(unit)
	}
}

// MustInUnit convert to the given unit, calls panic() if the units are not compatible
func (cm *ConcreteMeasurement) MustInUnit(p PrefixedUnit) Measurement {
	if ret, err := cm.InUnit(p); err != nil {
		panic(err)
	} else {
		return ret
	}
}

// MustInStringUnit return a new measurement in the new units
func (cm *ConcreteMeasurement) MustInStringUnit(symbol string) Measurement {
	if ret, err := cm.InStringUnit(symbol); err != nil {
		panic(err)
	} else {
		return ret
	}
}

// ConvertToString deprecated, please use InStringUnit
func (cm *ConcreteMeasurement) ConvertToString(s string) float64 {
	if unit, err := cm.InStringUnit(s); err != nil {
		panic(err)
	} else {
		return unit.RawValue()
	}
}

// String will return a summary of the ConcreteMeasurement Value and prefixed unit as a string.
// The value will be formatted in scientific notation for large exponents and the value unbounded.
// The Summary() method should be used to return a rounded string.
func (cm *ConcreteMeasurement) String() string {
	if cm.IsNil() {
		return ""
	}
	return fmt.Sprintf("%g %s", cm.RawValue(), cm.Unit().PrefixedSymbol())
}

// add to this

// IncrBy add the measurement m to the receiver
func (cm *ConcreteMeasurement) IncrBy(m Measurement) error {
	if isNil(cm) {
		return nil
	} else if rhs, err := m.InUnit(cm.Unit()); err != nil {
		return err
	} else {
		cm.SetValue(rhs.RawValue() + cm.RawValue())
	}
	return nil
}

// DecrBy subtract m from the receiver
func (cm *ConcreteMeasurement) DecrBy(m Measurement) error {
	if isNil(cm) {
		return nil
	} else if rhs, err := m.InUnit(cm.Unit()); err != nil {
		return err
	} else {
		cm.SetValue(cm.RawValue() - rhs.RawValue())
	}
	return nil
}

// Add deprecated, please use IncrBy
func (cm *ConcreteMeasurement) Add(m Measurement) {
	if err := cm.IncrBy(m); err != nil {
		panic(err)
	}
}

// Subtract deprecated, please use DecrBy
func (cm *ConcreteMeasurement) Subtract(m Measurement) {
	if err := cm.DecrBy(m); err != nil {
		panic(err)
	}
}

// multiply
func (cm *ConcreteMeasurement) MultiplyBy(factor float64) {
	if isNil(cm) {
		return
	}
	cm.SetValue(cm.RawValue() * float64(factor))
}

func (cm *ConcreteMeasurement) DivideBy(factor float64) {
	if isNil(cm) {
		return
	}
	cm.SetValue(cm.RawValue() / float64(factor))
}

// define a zero

func (cm *ConcreteMeasurement) IsNil() bool {
	return isNil(cm)
}

func (cm *ConcreteMeasurement) IsZero() bool {
	if isNil(cm) || math.Abs(cm.Mvalue) < 0.00000000001 {
		return true
	}
	return false
}

// IsPositive true if the measurement is positive by more than a very small delta
func (cm *ConcreteMeasurement) IsPositive() bool {
	return cm.Mvalue > 0.0 && !cm.IsZero()
}

// IsNegative true if the measurement is negative by more than a very small delta
func (cm *ConcreteMeasurement) IsNegative() bool {
	return cm.Mvalue < 0.0 && !cm.IsZero()
}

// less sensitive comparison operators

func (cm *ConcreteMeasurement) LessThanRounded(m Measurement, p int) bool {
	// nil means less than everything
	if isNil(cm) {
		return true
	} else if rhs, err := m.InUnit(cm.Unit()); err != nil {
		panic(err)
	} else {
		return wutil.RoundIgnoreNan(rhs.RawValue(), p) > wutil.RoundIgnoreNan(cm.RawValue(), p)
	}
}

func (cm *ConcreteMeasurement) GreaterThanRounded(m Measurement, p int) bool {
	if isNil(cm) {
		return false
	} else if rhs, err := m.InUnit(cm.Unit()); err != nil {
		panic(err)
	} else {
		return wutil.RoundIgnoreNan(rhs.RawValue(), p) < wutil.RoundIgnoreNan(cm.RawValue(), p)
	}
}

func (cm *ConcreteMeasurement) EqualToRounded(m Measurement, p int) bool {
	// this is not equal to anything
	if isNil(cm) {
		return false
	} else if rhs, err := m.InUnit(cm.Unit()); err != nil {
		panic(err)
	} else {
		return wutil.RoundIgnoreNan(rhs.RawValue(), p) == wutil.RoundIgnoreNan(cm.RawValue(), p)
	}
}

// comparison operators

func (cm *ConcreteMeasurement) LessThan(m Measurement) bool {
	// nil means less than everything
	if isNil(cm) {
		return true
	} else if rhs, err := m.InUnit(cm.Unit()); err != nil {
		panic(err)
	} else {
		return rhs.RawValue() > cm.RawValue()
	}
}

func (cm *ConcreteMeasurement) GreaterThan(m Measurement) bool {
	if isNil(cm) {
		return false
	} else if rhs, err := m.InUnit(cm.Unit()); err != nil {
		panic(err)
	} else {
		return rhs.RawValue() < cm.RawValue()
	}
}

// XXX This should be made more literal and rounded behaviour explicitly called for by user
func (cm *ConcreteMeasurement) EqualTo(m Measurement) bool {
	// this is not equal to anything

	if isNil(cm) {
		return false
	} else if rhs, err := m.InUnit(cm.Unit()); err != nil {
		panic(err)
	} else {
		dif := math.Abs(rhs.RawValue() - cm.RawValue())
		epsilon := math.Nextafter(1, 2) - 1
		return dif < (epsilon * 10000)
	}
}

// EqualToTolerance return true if the two measurements are within a small tolerace, tol, of each other
// where tol is expressed in the same units as the receiver
func (cm *ConcreteMeasurement) EqualToTolerance(m Measurement, tol float64) bool {
	if isNil(cm) {
		return false
	} else if rhs, err := m.InUnit(cm.Unit()); err != nil {
		panic(err)
	} else {
		return math.Abs(rhs.RawValue()-cm.RawValue()) < tol
	}
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
