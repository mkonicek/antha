package wunit

import (
	"github.com/pkg/errors"
	"strings"
	"sync"
)

// UnitRegistry store all the valid units in the library
type UnitRegistry struct {
	unitByType        map[string]map[string]bool
	orderedUnitByType map[string][]string
	unitBySymbol      map[string]*Unit
	aliases           map[string]string
	mutex             *sync.Mutex
}

// NewUnitRegistry build a new empty unit registry
func NewUnitRegistry() *UnitRegistry {
	return &UnitRegistry{
		unitByType:        make(map[string]map[string]bool),
		orderedUnitByType: make(map[string][]string),
		unitBySymbol:      make(map[string]*Unit),
		aliases:           make(map[string]string),
		mutex:             &sync.Mutex{},
	}
}

// DeclareUnit add a unit to the registry, as well as corresponding entries for valid prefixes
// If validPrefixes is zero length, only the base symbol will be added
func (self *UnitRegistry) DeclareUnit(measurementType, name, baseSymbol, SISymbol string, validPrefixes []SIPrefix, exponent int) error {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	unit := &Unit{
		name:       name,
		symbol:     baseSymbol,
		siSymbol:   SISymbol,
		multiplier: 1.0,
		exponent:   exponent,
		prefix:     None,
	}

	//add the prefix-less unit
	if err := self.declareUnit(measurementType, unit); err != nil {
		return err
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
//   reg.DeclareAlias("volume", "L", "l", SIPrefixes)
// will lead to all units with "L" (e.g. "uL", "mL") being converted to "l" (e.g. "ul", "ml", etc).
// If validPrefixes is zero length, only the base symbol will be added
// Note there is no value scaling, for that see DeclareDerivedUnit
func (self *UnitRegistry) DeclareAlias(measurementType, baseSymbol, baseTarget string, validPrefixes []SIPrefix) error {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	//add the prefix-less alias
	if err := self.declareAlias(measurementType, baseSymbol, baseTarget); err != nil {
		return err
	}

	for _, prefix := range validPrefixes {
		if err := self.declareAlias(measurementType, prefix.Symbol+baseSymbol, prefix.Symbol+baseTarget); err != nil {
			return err
		}
	}

	return nil
}

// declareAlias declare an alias, should have the lock when calling
func (self *UnitRegistry) declareAlias(measurementType, symbol, target string) error {
	if existingTarget, ok := self.aliases[symbol]; ok {
		return errors.Errorf("cannot declare alias %s = %s: alias %s = %s already declared", symbol, target, symbol, existingTarget)
	} else if !self.validUnitForType(measurementType, target) {
		return errors.Errorf("cannot declare alias %s == %s: unit %q is not of type %s", symbol, target, target, measurementType)
	} else if existing, ok := self.unitBySymbol[symbol]; ok {
		return errors.Errorf("cannot declare alias %s == %s: would shadow pre-existing unit %v", symbol, target, existing)
	} else {
		self.addUnitByType(measurementType, symbol)
		self.aliases[symbol] = target
		return nil
	}
}

// GetUnit return the unit referred to by symbol
func (self *UnitRegistry) GetUnit(symbol string) (*Unit, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return self.getUnit(symbol)
}

// getUnit should be called with the lock
func (self *UnitRegistry) getUnit(symbol string) (*Unit, error) {
	symbol = self.resolveAliasing(symbol)

	if unit, ok := self.unitBySymbol[symbol]; !ok {
		return nil, errors.Errorf("unknown unit symbol %q", symbol)
	} else {
		return unit.Copy(), nil
	}
}

// declareUnit declare a new unit, should have the lock when calling
func (self *UnitRegistry) declareUnit(measurementType string, unit *Unit) error {
	if _, ok := self.unitBySymbol[unit.PrefixedSymbol()]; ok {
		return errors.Errorf("cannot declare unit %q: unit already declared", unit.PrefixedSymbol())
	}
	self.addUnitByType(measurementType, unit.PrefixedSymbol())
	self.unitBySymbol[unit.PrefixedSymbol()] = unit.Copy()
	return nil
}

func (self *UnitRegistry) addUnitByType(measurementType, symbol string) {
	if _, ok := self.unitByType[measurementType]; !ok {
		self.unitByType[measurementType] = make(map[string]bool)
	}

	self.unitByType[measurementType][symbol] = true

	ordered := self.orderedUnitByType[measurementType]

	if ordered == nil {
		// 32 is a magic number, hopefully more than the number of symbols for any
		// particular measurement type
		ordered = make([]string, 0, 32)
	}

	self.orderedUnitByType[measurementType] = insertionSort(ordered, symbol)
}

// insertionSort inserts value into slice such that ordering is maintained
// slice can be empty or nil
func insertionSort(slice []string, value string) []string {
	//find out where to add the value
	pos := 0
	for ; pos < len(slice); pos++ {
		if value <= slice[pos] {
			break
		}
	}

	//make sure there's space for the value
	slice = append(slice, "")
	slice = append(slice[:pos+1], slice[pos:len(slice)-1]...)
	slice[pos] = value

	return slice
}

// ValidUnitForType return true if the given symbol represents a unit that is valid
// for the given measurement type
// e.g. ValidUnitForType("Length", "m") -> true
// and  ValidUnitForType("Area", "l") -> false
func (self *UnitRegistry) ValidUnitForType(measurementType, symbol string) bool {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return self.validUnitForType(measurementType, symbol)
}

// validUnitForType should have the lock when calling
func (self *UnitRegistry) validUnitForType(measurementType, symbol string) bool {
	return self.unitByType[measurementType][self.resolveAliasing(symbol)]
}

// AssertValidForType assert that the symbol refers to a valid unit for the given type
// the same as ValidUnitForType, except this function returns a useful
func (self *UnitRegistry) AssertValidUnitForType(measurementType, symbol string) error {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	if !self.validUnitForType(measurementType, symbol) {
		return errors.Errorf("invalid symbol %q for measurement type %q, valid symbols are %v", symbol, measurementType, self.orderedUnitByType[measurementType])
	}
	return nil
}

// ListValidUnitsForType returns a sorted list of all valid unit symbols for a given
// measurement type
func (self *UnitRegistry) ListValidUnitsForType(measurementType string) []string {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return self.orderedUnitByType[measurementType]
}

// DeclareDerivedUnit such that references to "symbol" are converted to "target" using
// the conversion factor symbolInTargets, for each valid prefix.
// The target should already exist in the Registry.
// If validPrefixes is nil or zero length, only the base unit will be added
// e.g. DeclareDerivedUnit("pint", nil, "l", 0.568) will cause the unit "1 pint" to be
// understood as "0.568 l"
func (self *UnitRegistry) DeclareDerivedUnit(measurementType string, name, symbol string, validPrefixes []SIPrefix, exponent int, target string, symbolInTargets float64) error {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	unit, err := self.getUnit(target)
	if err != nil {
		return err
	}

	if !self.validUnitForType(measurementType, unit.symbol) {
		return errors.Errorf("cannot declare derived unit %s = %f %s: %s is not of type %q", symbol, symbolInTargets, target, target, measurementType)
	}

	unit.symbol = symbol
	unit.name = name
	unit.multiplier = unit.getBaseSIConversionFactor() * symbolInTargets
	unit.exponent = exponent
	unit.prefix = None

	//declare the prefix-less unit
	if err := self.declareUnit(measurementType, unit); err != nil {
		return err
	}

	for _, prefix := range validPrefixes {
		unit.prefix = prefix
		if err := self.declareUnit(measurementType, unit); err != nil {
			return err
		}
	}

	return nil
}

// resolveAliasing convert the given symbol into a known symbol if it is aliased
// and also do any µ/u style conversions, should have the mutex lock when calling
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
	self.mutex.Lock()
	defer self.mutex.Unlock()
	if unit, err := self.getUnit(unitSymbol); err != nil {
		return nil, err
	} else {
		return &ConcreteMeasurement{
			Mvalue: value,
			Munit:  unit,
		}, nil
	}
}
