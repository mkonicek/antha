package wunit

import (
	"encoding/json"
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
	exponent   int      //the exponent for the prefix. 1 unless the prefix is grouped with a unit that is raised to a power, e.g. for "m^2" exponent=2 such that 1 m^2 = 10^6 mm^2
}

// MarshalJSON marshal the unit as a JSON string
func (self *Unit) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name       string
		Symbol     string
		Base       string
		Prefix     SIPrefix
		Multiplier float64
		Exponent   int
	}{
		Name:       self.name,
		Symbol:     self.symbol,
		Base:       self.base,
		Prefix:     self.prefix,
		Multiplier: self.multiplier,
		Exponent:   self.exponent,
	})
}

// UnmarshalJSON marshal the unit as a JSON string
func (self *Unit) UnmarshalJSON(data []byte) error {
	value := struct {
		Name       string
		Symbol     string
		Base       string
		Prefix     SIPrefix
		Multiplier float64
		Exponent   int
	}{}

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	self.name = value.Name
	self.symbol = value.Symbol
	self.base = value.Base
	self.prefix = value.Prefix
	self.multiplier = value.Multiplier
	self.exponent = value.Exponent
	return nil
}

// Name get the full name of the unit
func (self *Unit) Name() string {
	return self.prefix.Name + self.name
}

// BaseSIConversionFactor factor to multiply by in order to convert to SI value including effect of prefix
func (self *Unit) BaseSIConversionFactor() float64 {
	return math.Pow(self.prefix.Value, float64(self.exponent)) * self.multiplier
}

// String a string representation of the unit name and symbol
func (self *Unit) String() string {
	return fmt.Sprintf("%s[%s]", self.Name(), self.PrefixedSymbol())
}

// Prefix get the SI prefix of this unit, or " " if none
func (self *Unit) Prefix() SIPrefix {
	return self.prefix
}

// PrefixedSymbol the symbol including any prefix
func (self *Unit) PrefixedSymbol() string {
	if self.prefix.Symbol == " " {
		return self.symbol
	}
	return self.prefix.Symbol + self.symbol
}

// RawSymbol symbol without prefex, equivalent to Symbol()
func (self *Unit) RawSymbol() string {
	return self.symbol
}

// BaseSISymbol Base SI or derived unit for this property, equivalent to BaseSISymbol
func (self *Unit) BaseSISymbol() string {
	return self.base
}

// ConvertTo get the conversion factor between this unit and pu.
// This function will call panic() if pu is not compatible with this unit, see CompatibleWith
func (self *Unit) ConvertTo(pu PrefixedUnit) (float64, error) {
	if !self.compatibleWith(pu) {
		return 0.0, errors.Errorf("cannot convert units: base units for %s and %s do not match: %s != %s", self.PrefixedSymbol(), pu.PrefixedSymbol(), self.base, pu.BaseSISymbol())
	}

	return self.BaseSIConversionFactor() / pu.BaseSIConversionFactor(), nil
}

// CompatibleWith returns true if the units can be converted to the supplied units.
// If this function returns false then calling ConvertTo with the same units will
// case a panic
func (self *Unit) compatibleWith(pu PrefixedUnit) bool {
	return self.base == pu.BaseSISymbol()
}

// Copy return a pointer to a new Unit identical to this one
func (self *Unit) Copy() *Unit {
	ret := *self
	return &ret
}
