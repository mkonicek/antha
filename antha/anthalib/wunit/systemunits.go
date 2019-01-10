package wunit

import "math"

type baseUnit struct {
	Name     string
	Symbol   string
	SISymbol string //the canonincal form for the unit which can include a prefix, defaults to Symbol
	Prefixes []SIPrefix
	Exponent int
}

type baseUnits map[string][]baseUnit

func (self baseUnits) AddTo(reg *UnitRegistry) error {
	for mType, defs := range self {
		for _, unit := range defs {
			SISymbol := unit.SISymbol
			if SISymbol == "" {
				SISymbol = unit.Symbol
			}
			if err := reg.DeclareUnit(mType, unit.Name, unit.Symbol, SISymbol, unit.Prefixes, unit.Exponent); err != nil {
				return err
			}
		}
	}
	return nil
}

func getSystemUnits() baseUnits {
	return baseUnits{
		"Concentration": {
			{
				Name:     "grams per litre",
				Symbol:   "g/l",
				SISymbol: "kg/l",
				Prefixes: SIPrefixes,
				Exponent: 1,
			},
			{
				Name:     "moles per litre",
				Symbol:   "Mol/l",
				Prefixes: SIPrefixes,
				Exponent: 1,
			},
			{
				Name:     "units per litre",
				Symbol:   "U/l",
				Prefixes: SIPrefixes,
				Exponent: 1,
			},
			{
				Name:   "relative concentration",
				Symbol: "X",
			},
			{
				Name:   "volume ratio",
				Symbol: "v/v",
			},
		},
		"Volume": {
			{
				Name:     "litre",
				Symbol:   "l",
				Prefixes: SIPrefixes,
				Exponent: 1,
			},
		},
		"Mass": {
			{
				Name:     "gram",
				Symbol:   "g",
				SISymbol: "kg",
				Prefixes: SIPrefixes,
				Exponent: 1,
			},
		},
		"Density": {
			{
				Name:     "grams per meter cubed",
				Symbol:   "g/m^3",
				SISymbol: "kg/m^3",
				Prefixes: SIPrefixes,
				Exponent: 1,
			},
		},
		"Length": {
			{
				Name:     "metre",
				Symbol:   "m",
				Prefixes: SIPrefixes,
				Exponent: 1,
			},
		},
		"Area": {
			{
				Name:     "metre squared",
				Symbol:   "m^2",
				Prefixes: SIPrefixes,
				Exponent: 2,
			},
		},
		"Temperature": {
			{
				Name:   "celsius",
				Symbol: "℃",
			},
		},
		"Time": {
			{
				Name:     "seconds",
				Symbol:   "s",
				Prefixes: []SIPrefix{Yocto, Zepto, Atto, Femto, Pico, Nano, Micro, Milli},
				Exponent: 1,
			},
		},
		"Moles": {
			{
				Name:     "moles",
				Symbol:   "Mol",
				Prefixes: SIPrefixes,
				Exponent: 1,
			},
		},
		"Angle": {
			{
				Name:   "radians",
				Symbol: "rad",
			},
		},
		"AngularVelocity": {
			{
				Name:   "radians per second",
				Symbol: "rad/s",
			},
		},
		"Energy": {
			{
				Name:     "joules",
				Symbol:   "J",
				Prefixes: SIPrefixes,
				Exponent: 1,
			},
		},
		"Force": {
			{
				Name:     "newtons",
				Symbol:   "N",
				Prefixes: SIPrefixes,
				Exponent: 1,
			},
		},
		"Pressure": {
			{
				Name:     "pascals",
				Symbol:   "Pa",
				Prefixes: SIPrefixes,
				Exponent: 1,
			},
		},
		"SpecificHeatCapacity": {
			{
				Name:     "joules per kilogram per degrees celsius",
				Symbol:   "J/kg*C",
				Prefixes: SIPrefixes,
				Exponent: 1,
			},
		},
		"Velocity": {
			{
				Name:     "meters per second",
				Symbol:   "m/s",
				Prefixes: SIPrefixes,
				Exponent: 1,
			},
		},
		"FlowRate": {
			{
				Name:     "litres per second",
				Symbol:   "l/s",
				Prefixes: SIPrefixes,
				Exponent: 1,
			},
		},
		"Rate": {
			{
				Name:   "per second",
				Symbol: "/s",
			},
		},
		"Voltage": {
			{
				Name:     "volts",
				Symbol:   "V",
				Prefixes: SIPrefixes,
				Exponent: 1,
			},
		},
	}
}

type derivedUnit struct {
	Name         string
	Symbol       string
	Prefixes     []SIPrefix
	Exponent     int
	TargetSymbol string
	TargetScale  float64 //i.e. how many target units are in 1 derived unit
}

type derivedUnits map[string][]derivedUnit

func (self derivedUnits) AddTo(reg *UnitRegistry) error {
	for mType, defs := range self {
		for _, du := range defs {
			if err := reg.DeclareDerivedUnit(mType, du.Name, du.Symbol, du.Prefixes, du.Exponent, du.TargetSymbol, du.TargetScale); err != nil {
				return err
			}
		}
	}
	return nil
}

func getSystemDerivedUnits() derivedUnits {
	return derivedUnits{
		"Concentration": {
			{
				Name:         "grams per mililitre",
				Symbol:       "g/ml",
				Prefixes:     SIPrefixes,
				Exponent:     1,
				TargetSymbol: "kg/l",
				TargetScale:  1.0,
			},
			{
				Name:         "grams per microlitre",
				Symbol:       "g/ul",
				Prefixes:     SIPrefixes,
				Exponent:     1,
				TargetSymbol: "Mg/l",
				TargetScale:  1.0,
			},
			{
				Name:         "grams per nanolitre",
				Symbol:       "g/nl",
				Prefixes:     SIPrefixes,
				Exponent:     1,
				TargetSymbol: "Gg/l",
				TargetScale:  1.0,
			},
			{
				Name:         "moles per mililitre",
				Symbol:       "Mol/ml",
				Prefixes:     SIPrefixes,
				Exponent:     1,
				TargetSymbol: "kMol/l",
				TargetScale:  1.0,
			},
			{
				Name:         "moles per microlitre",
				Symbol:       "Mol/ul",
				Prefixes:     SIPrefixes,
				Exponent:     1,
				TargetSymbol: "MMol/l",
				TargetScale:  1.0,
			},
			{
				Name:         "moles per nanolitre",
				Symbol:       "Mol/nl",
				Prefixes:     SIPrefixes,
				Exponent:     1,
				TargetSymbol: "GMol/l",
				TargetScale:  1.0,
			},
			{
				Name:         "units per mililitre",
				Symbol:       "U/ml",
				Prefixes:     SIPrefixes,
				Exponent:     1,
				TargetSymbol: "kU/l",
				TargetScale:  1.0,
			},
			{
				Name:         "units per microlitre",
				Symbol:       "U/ul",
				Prefixes:     SIPrefixes,
				Exponent:     1,
				TargetSymbol: "MU/l",
				TargetScale:  1.0,
			},
			{
				Name:         "units per nanolitre",
				Symbol:       "U/nl",
				Prefixes:     SIPrefixes,
				Exponent:     1,
				TargetSymbol: "GU/l",
				TargetScale:  1.0,
			},
			{
				Name:         "percentage weight of solution",
				Symbol:       "% w/v",
				TargetSymbol: "g/l",
				TargetScale:  10.0,
			},
			{
				Name:         "percentage volume per volume of solution",
				Symbol:       "% v/v",
				TargetSymbol: "v/v",
				TargetScale:  0.01,
			},
		},
		"Volume": {
			{
				Name:         "meters cubed",
				Symbol:       "m^3",
				Prefixes:     SIPrefixes,
				Exponent:     3,
				TargetSymbol: "l",
				TargetScale:  1000.0,
			},
		},
		"Length": {
			{
				Name:         "inches",
				Symbol:       "in",
				TargetSymbol: "mm",
				TargetScale:  25.4,
			},
		},
		"Time": {
			{
				Name:         "minutes",
				Symbol:       "min",
				TargetSymbol: "s",
				TargetScale:  60.0,
			},
			{
				Name:         "hours",
				Symbol:       "h",
				TargetSymbol: "s",
				TargetScale:  3600.0,
			},
		},
		"Angle": {
			{
				Name:         "degrees",
				Symbol:       "°",
				TargetSymbol: "rad",
				TargetScale:  (2.0 * math.Pi) / 360.0,
			},
		},
		"AngularVelocity": {
			{
				Name:         "radians per minute",
				Symbol:       "rad/min",
				TargetSymbol: "rad/s",
				TargetScale:  1.0 / 60.0,
			},
			{
				Name:         "revolutions per minute",
				Symbol:       "rpm",
				TargetSymbol: "rad/s",
				TargetScale:  2.0 * math.Pi / 60.0,
			},
		},
		"Pressure": {
			{
				Name:         "bar",
				Symbol:       "bar",
				Prefixes:     SIPrefixes,
				Exponent:     1,
				TargetSymbol: "kPa",
				TargetScale:  100.0,
			},
		},
		"FlowRate": {
			{
				Name:         "litres per minute",
				Symbol:       "l/min",
				Prefixes:     SIPrefixes,
				Exponent:     1,
				TargetSymbol: "l/s",
				TargetScale:  1 / 60.0,
			},
			{
				Name:         "litres per hour",
				Symbol:       "l/h",
				Prefixes:     SIPrefixes,
				Exponent:     1,
				TargetSymbol: "l/s",
				TargetScale:  1 / 3600.0,
			},
		},
		"Rate": {
			{
				Name:         "per minute",
				Symbol:       "/min",
				TargetSymbol: "/s",
				TargetScale:  1 / 60.0,
			},
			{
				Name:         "per hour",
				Symbol:       "/h",
				TargetSymbol: "/s",
				TargetScale:  1 / 3600.0,
			},
		},
	}
}

type unitAlias struct {
	BaseSymbol string
	BaseTarget string
	Prefixes   []SIPrefix
}

type unitAliases map[string][]unitAlias

func (self unitAliases) AddTo(reg *UnitRegistry) error {
	for mType, defs := range self {
		for _, a := range defs {
			if err := reg.DeclareAlias(mType, a.BaseSymbol, a.BaseTarget, a.Prefixes); err != nil {
				return err
			}
		}
	}
	return nil
}

func getSystemAliases() unitAliases {
	return unitAliases{
		"Concentration": {
			{
				BaseSymbol: "g/L",
				BaseTarget: "g/l",
				Prefixes:   SIPrefixes,
			},
			{
				BaseSymbol: "g/mL",
				BaseTarget: "g/ml",
				Prefixes:   SIPrefixes,
			},
			{
				BaseSymbol: "g/uL",
				BaseTarget: "g/ul",
				Prefixes:   SIPrefixes,
			},
			{
				BaseSymbol: "M",
				BaseTarget: "Mol/l",
				Prefixes:   SIPrefixes,
			},
			{
				BaseSymbol: "Molar",
				BaseTarget: "Mol/l",
				Prefixes:   SIPrefixes,
			},
			{
				BaseSymbol: "M/l",
				BaseTarget: "Mol/l",
				Prefixes:   SIPrefixes,
			},
			{
				BaseSymbol: "M/ml",
				BaseTarget: "Mol/ml",
				Prefixes:   SIPrefixes,
			},
			{
				BaseSymbol: "M/ul",
				BaseTarget: "Mol/ul",
				Prefixes:   SIPrefixes,
			},
			{
				BaseSymbol: "Mol/L",
				BaseTarget: "Mol/l",
				Prefixes:   SIPrefixes,
			},
			{
				BaseSymbol: "Mol/mL",
				BaseTarget: "Mol/ml",
				Prefixes:   SIPrefixes,
			},
			{
				BaseSymbol: "Mol/uL",
				BaseTarget: "Mol/ul",
				Prefixes:   SIPrefixes,
			},
			{
				BaseSymbol: "M/L",
				BaseTarget: "Mol/l",
				Prefixes:   SIPrefixes,
			},
			{
				BaseSymbol: "M/mL",
				BaseTarget: "Mol/ml",
				Prefixes:   SIPrefixes,
			},
			{
				BaseSymbol: "M/uL",
				BaseTarget: "Mol/ul",
				Prefixes:   SIPrefixes,
			},
			{
				BaseSymbol: "x",
				BaseTarget: "X",
			},
			{
				BaseSymbol: "w/v",
				BaseTarget: "% w/v",
			},
		},
		"Volume": {
			{
				BaseSymbol: "L",
				BaseTarget: "l",
				Prefixes:   SIPrefixes,
			},
		},
		"Length": {
			{
				BaseSymbol: `"`,
				BaseTarget: "in",
			},
		},
		"Temperature": {
			{
				BaseSymbol: "C",
				BaseTarget: "℃",
			},
			{
				BaseSymbol: "˚C",
				BaseTarget: "℃",
			},
			{
				BaseSymbol: "°C",
				BaseTarget: "℃",
			},
		},
		"Time": {
			{
				BaseSymbol: "minutes",
				BaseTarget: "min",
			},
		},
		"Angle": {
			{
				BaseSymbol: "radians",
				BaseTarget: "rad",
			},
			{
				BaseSymbol: "deg",
				BaseTarget: "°",
			},
			{
				BaseSymbol: "degrees",
				BaseTarget: "°",
			},
			{
				BaseSymbol: "˚",
				BaseTarget: "°",
			},
		},
	}
}

// makeGlobalUnitRegistry return a new registry pre-populated with system units
func makeGlobalUnitRegistry() *UnitRegistry {
	reg := NewUnitRegistry()

	if err := getSystemUnits().AddTo(reg); err != nil {
		panic(err)
	}

	if err := getSystemDerivedUnits().AddTo(reg); err != nil {
		panic(err)
	}

	if err := getSystemAliases().AddTo(reg); err != nil {
		panic(err)
	}

	return reg
}

var globalRegistry *UnitRegistry

// GetGlobalUnitRegistry gets the shared unit registry which contains system types
func GetGlobalUnitRegistry() *UnitRegistry {
	if globalRegistry == nil {
		globalRegistry = makeGlobalUnitRegistry()
	}
	return globalRegistry
}
