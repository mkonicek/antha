package wunit

import (
	"fmt"
	"github.com/pkg/errors"
	"math"
)

// Unit everything we need to know about a unit to support it
type Unit struct {
	name       string   //common name for the unit
	symbol     string   //symbol of the unit
	base       string   //the SI unit for this dimension, or derived one if there isn't one
	prefix     SIPrefix //the SI prefix which is applied to the symbol
	multiplier float64  //value to multiply by to convert "symbol"s to "base", e.g. 60 for min (SI unit = s)
	exponent   int      //the exponent for the prefix, 1 unless the prefix is inclided in a power, e.g. exponent=2 for "m^2"
}

// Name the common name for this unit
func (self *Unit) Name() string {
	return self.prefix.LongName() + self.name
}

// Symbol the symbol used to denote this unit
func (self *Unit) Symbol() string {
	return self.symbol
}

// BaseSIConversionFactor factor to multiply by in order to convert to SI value including effect of prefix
func (self *Unit) BaseSIConversionFactor() float64 {
	return math.Pow(self.prefix.Value, float64(self.exponent)) * self.multiplier
}

// BaseSIUnit Base SI or derived unit for this property, equivalent to BaseSISymbol
func (self *Unit) BaseSIUnit() string {
	return self.base
}

// ToString return the symbol and the prefix
func (self *Unit) ToString() string {
	if self.prefix.Name == " " {
		return self.symbol
	}
	return self.prefix.Name + self.symbol
}

// String a string representation of the name and symbol
func (self *Unit) String() string {
	return fmt.Sprintf("%s[%s]", self.Name(), self.ToString())
}

// Prefix get the SI prefix of this unit, or " " if none
func (self *Unit) Prefix() SIPrefix {
	return self.prefix
}

// PrefixedSymbol the symbol including any prefix
func (self *Unit) PrefixedSymbol() string {
	return self.prefix.Name + self.symbol
}

// RawSymbol symbol without prefex, equivalent to Symbol()
func (self *Unit) RawSymbol() string {
	return self.symbol
}

// BaseSISymbol Base SI or derived unit for this property, equivalent to BaseSIUnit
func (self *Unit) BaseSISymbol() string {
	return self.base
}

// ConvertTo get the conversion factor between this unit and pu
func (self *Unit) ConvertTo(pu PrefixedUnit) float64 {
	if self.base != pu.BaseSIUnit() {
		panic(errors.Errorf("cannot convert units: base units for %s and %s do not match: %s != %s", self.ToString(), pu.ToString(), self.base, pu.BaseSIUnit()))
	}

	return self.BaseSIConversionFactor() / pu.BaseSIConversionFactor()
}

// Copy return a pointer to a new Unit identical to this one
func (self *Unit) Copy() *Unit {
	ret := *self
	return &ret
}
