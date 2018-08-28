package wunit

import (
	"github.com/pkg/errors"
	"sort"
	"strings"
)

// UnitRegistry store all the valid units in the library
type UnitRegistry struct {
	unitByType   map[string]map[string]bool
	unitBySymbol map[string]*Unit
	aliases      map[string]string
}

// NewUnitRegistry build a new empty unit registry
func NewUnitRegistry() *UnitRegistry {
	return &UnitRegistry{
		unitByType:   make(map[string]map[string]bool),
		unitBySymbol: make(map[string]*Unit),
		aliases:      make(map[string]string),
	}
}

// DeclareUnit add a unit to the registry, as well as corresponding entries for valid prefixes
// If validPrefices is zero length, only the base symbol will be added
func (self *UnitRegistry) DeclareUnit(measurementType, name, baseSymbol, SISymbol string, validPrefixes []SIPrefix, exponent int) error {
	if len(validPrefixes) == 0 {
		validPrefixes = []SIPrefix{None} //if no prefices specified, only add default
	}
	unit := &Unit{
		name:       name,
		symbol:     baseSymbol,
		base:       SISymbol,
		multiplier: 1.0,
		exponent:   exponent,
	}

	for _, prefix := range validPrefixes {
		unit.prefix = prefix
		if err := self.declareUnit(measurementType, unit); err != nil {
			return err
		}
	}

	return nil
}

// DeclareAlias declare an alias for a target symbol such that units with the alias are converted to the target.
// This is expected to be used when there are multiple convensions for writing a unit, for example
//   reg.DeclareAlias("volume", "L", "l", SIPrefices)
// will lead to all units with "L" (e.g. "uL", "mL") being converted to "l" (e.g. "ul", "ml", etc).
// If validPrefices is zero length, only the base symbol will be added
// Note there is no value scaling, for that see DeclareDerivedUnit
func (self *UnitRegistry) DeclareAlias(measurementType, baseSymbol, baseTarget string, validPrefixes []SIPrefix) error {
	if len(validPrefixes) == 0 {
		validPrefixes = []SIPrefix{None} //if no prefices specified, only add default
	}

	for _, prefix := range validPrefixes {
		if err := self.declareAlias(measurementType, prefix.Symbol+baseSymbol, prefix.Symbol+baseTarget); err != nil {
			return err
		}
	}

	return nil
}

func (self *UnitRegistry) declareAlias(measurementType, symbol, target string) error {
	if existingTarget, ok := self.aliases[symbol]; ok {
		return errors.Errorf("cannot declare alias %s = %s: alias %s = %s already declared", symbol, target, symbol, existingTarget)
	} else if !self.ValidUnitForType(measurementType, target) {
		return errors.Errorf("cannot declare alias %s == %s: unit %q is not of type %s", symbol, target, target, measurementType)
	} else if existing, ok := self.unitBySymbol[symbol]; ok {
		return errors.Errorf("cannot declare alias %s == %s: would shadow pre-existing unit %v", symbol, target, existing)
	} else {
		self.unitByType[measurementType][symbol] = true
		self.aliases[symbol] = target
		return nil
	}
}

// GetUnit return the unit referred to by symbol
func (self *UnitRegistry) GetUnit(symbol string) (*Unit, error) {

	symbol = self.resolveAliasing(symbol)

	if unit, ok := self.unitBySymbol[symbol]; !ok {
		return nil, errors.Errorf("unknown unit symbol %q", symbol)
	} else {
		return unit.Copy(), nil
	}
}

func (self *UnitRegistry) declareUnit(measurementType string, unit *Unit) error {
	if _, ok := self.unitBySymbol[unit.ToString()]; ok {
		return errors.Errorf("cannot declare unit %q: unit already declared", unit.ToString())
	}
	if _, ok := self.unitByType[measurementType]; !ok {
		self.unitByType[measurementType] = make(map[string]bool)
	}

	self.unitByType[measurementType][unit.ToString()] = true
	self.unitBySymbol[unit.ToString()] = unit.Copy()
	return nil
}

// ValidUnitForType return true if the given symbol represents a unit that is valid
// for the given measurement type
// e.g. ValidUnitForType("Length", "m") -> true
// and  ValidUnitForType("Area", "l") -> false
func (self *UnitRegistry) ValidUnitForType(measurementType, symbol string) bool {
	return self.unitByType[measurementType][self.resolveAliasing(symbol)]
}

// AssertValidForType assert that the symbol refers to a valid unit for the given type
// the same as ValidUnitForType, except this function returns a useful
func (self *UnitRegistry) AssertValidUnitForType(measurementType, symbol string) error {
	if !self.ValidUnitForType(measurementType, symbol) {
		return errors.Errorf("invalid symbol %q for measurement type %q, valid symbols are %v", symbol, measurementType, self.unitByType[measurementType])
	}
	return nil
}

// ListValidUnitsForType returns a sorted list of all valid unit symbols for a given
// measurement type
func (self *UnitRegistry) ListValidUnitsForType(measurementType string) []string {
	if symbols, ok := self.unitByType[measurementType]; !ok {
		return nil
	} else {
		ret := make([]string, 0, len(symbols))
		for symbol := range symbols {
			ret = append(ret, symbol)
		}
		sort.Strings(ret)
		return ret
	}
}

// DeclareDerivedUnit such that references to "symbol" are converted to "target" using
// the conversion factor symbolInTargets, for each valid prefix.
// The target should already exist in the Registry.
// If validPrefixes is nil or zero length, only the base unit will be added
// e.g. DeclareDerivedUnit("pint", nil, "l", 0.568) will cause the unit "1 pint" to be
// understood as "0.568 l"
func (self *UnitRegistry) DeclareDerivedUnit(measurementType string, name, symbol string, validPrefixes []SIPrefix, exponent int, target string, symbolInTargets float64) error {
	if len(validPrefixes) == 0 {
		validPrefixes = []SIPrefix{None} //if no prefices specified, only add default
	}

	unit, err := self.GetUnit(target)
	if err != nil {
		return err
	}

	if !self.ValidUnitForType(measurementType, unit.symbol) {
		return errors.Errorf("cannot declare derived unit %s = %f %s: %s is not of type %q", symbol, symbolInTargets, target, target, measurementType)
	}

	unit.symbol = symbol
	unit.name = name
	unit.multiplier = unit.BaseSIConversionFactor() * symbolInTargets
	unit.exponent = exponent

	for _, prefix := range validPrefixes {
		unit.prefix = prefix
		if err := self.declareUnit(measurementType, unit); err != nil {
			return err
		}
	}

	return nil
}

// resolveAliasing convert the given symbol into a known symbol if it is aliased
// and also do any µ/u style conversions
func (self *UnitRegistry) resolveAliasing(symbol string) string {
	symbol = strings.Replace(symbol, "µ", "u", -1)
	symbol = strings.Trim(symbol, " ")

	if alias, ok := self.aliases[symbol]; ok {
		return alias
	}
	return symbol
}

// NewMeasurement return a new typed measurement
func (self *UnitRegistry) NewMeasurement(value float64, unitSymbol string) (*ConcreteMeasurement, error) {
	if unit, err := self.GetUnit(unitSymbol); err != nil {
		return nil, err
	} else {
		return &ConcreteMeasurement{
			Mvalue: value,
			Munit:  unit,
		}, nil
	}
}
