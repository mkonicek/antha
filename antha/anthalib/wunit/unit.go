package wunit

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"math"
)

// Unit everything we need to know about a unit to support it
type Unit struct {
	name       string   //common name for the unit
	symbol     string   //symbol of the unit
	siSymbol   string   //the symbol of the SI unit for this unit
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
		Base:       self.siSymbol,
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
	self.siSymbol = value.Base
	self.prefix = value.Prefix
	self.multiplier = value.Multiplier
	self.exponent = value.Exponent
	return nil
}

// GobEncode encode the unit as gob
func (self *Unit) GobEncode() ([]byte, error) {
	var out bytes.Buffer
	if err := gob.NewEncoder(&out).Encode(struct {
		Name       string
		Symbol     string
		Base       string
		Prefix     SIPrefix
		Multiplier float64
		Exponent   int
	}{
		Name:       self.name,
		Symbol:     self.symbol,
		Base:       self.siSymbol,
		Prefix:     self.prefix,
		Multiplier: self.multiplier,
		Exponent:   self.exponent,
	}); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// GobDecode deserialise gob
func (self *Unit) GobDecode(b []byte) error {
	value := struct {
		Name       string
		Symbol     string
		Base       string
		Prefix     SIPrefix
		Multiplier float64
		Exponent   int
	}{}

	if err := gob.NewDecoder(bytes.NewBuffer(b)).Decode(&value); err != nil {
		return err
	}
	self.name = value.Name
	self.symbol = value.Symbol
	self.siSymbol = value.Base
	self.prefix = value.Prefix
	self.multiplier = value.Multiplier
	self.exponent = value.Exponent
	return nil
}

// Name get the full name of the unit
func (self *Unit) Name() string {
	if self == nil {
		return ""
	}
	return self.prefix.Name + self.name
}

// getBaseSIConversionFactor factor to multiply by in order to convert to SI value including effect of prefix
func (self *Unit) getBaseSIConversionFactor() float64 {
	if self == nil {
		return 0.0
	}
	return math.Pow(self.prefix.Value, float64(self.exponent)) * self.multiplier
}

// String a string representation of the unit name and symbol
func (self *Unit) String() string {
	if self == nil {
		return ""
	}
	return fmt.Sprintf("%s[%s]", self.Name(), self.PrefixedSymbol())
}

// Prefix get the SI prefix of this unit, or " " if none
func (self *Unit) Prefix() SIPrefix {
	if self == nil {
		return SIPrefix{}
	}
	return self.prefix
}

// PrefixedSymbol the symbol including any prefix
func (self *Unit) PrefixedSymbol() string {
	if self == nil {
		return ""
	}
	if self.prefix.Symbol == " " {
		return self.symbol
	}
	return self.prefix.Symbol + self.symbol
}

// RawSymbol symbol without prefex, equivalent to Symbol()
func (self *Unit) RawSymbol() string {
	if self == nil {
		return ""
	}
	return self.symbol
}

// BaseSISymbol Base SI or derived unit for this property, equivalent to BaseSISymbol
func (self *Unit) BaseSISymbol() string {
	if self == nil {
		return ""
	}
	return self.siSymbol
}

// getConversionFactor get the conversion factor between this unit and rhs.
func (self *Unit) getConversionFactor(rhs *Unit) (float64, error) {
	if self == nil || rhs == nil {
		return 0.0, errors.New("cannot convert units: nil units provided")
	}
	if !self.compatibleWith(rhs) {
		return 0.0, errors.Errorf("cannot convert units: base units for %s and %s do not match: %s != %s", self.PrefixedSymbol(), rhs.PrefixedSymbol(), self.siSymbol, rhs.BaseSISymbol())
	}

	return self.getBaseSIConversionFactor() / rhs.getBaseSIConversionFactor(), nil
}

// compatibleWith returns true if the units can be converted to the supplied units.
func (self *Unit) compatibleWith(pu PrefixedUnit) bool {
	return self.siSymbol == pu.BaseSISymbol()
}

// Copy return a pointer to a new Unit identical to this one
func (self *Unit) Copy() *Unit {
	ret := *self
	return &ret
}
