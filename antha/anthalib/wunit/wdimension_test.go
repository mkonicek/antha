package wunit

import (
	"fmt"
	"math"
	"testing"
)

type MeasurementConstructor func(float64, string) Measurement

type TestFn func(*testing.T, Measurement)

type NewMeasurementTest struct {
	Value            float64
	Unit             string
	ExpectedSIValue  float64
	ExpectedBaseUnit string
	ExpectedPrefix   string
	ExpectError      bool
	ConvertToString  map[string]float64
}

func (self *NewMeasurementTest) Run(t *testing.T, constructor MeasurementConstructor) {

	t.Run(fmt.Sprintf("%f_%s", self.Value, self.Unit), func(t *testing.T) {
		defer func() {
			if r := recover(); (r != nil) != self.ExpectError {
				t.Errorf("error mismatch: expectError: %t, got error: \"%v\"", self.ExpectError, r)
			}
		}()

		m := constructor(self.Value, self.Unit)

		if !self.ExpectError { //don't check these if we were expecting error, the defer statement will add one
			if math.Abs(m.SIValue()-self.ExpectedSIValue) > 1.0e-9 {
				t.Errorf("wrong SIValue: expected %e, got %e", self.ExpectedSIValue, m.SIValue())
			}

			if e, g := self.ExpectedBaseUnit, m.Unit().BaseSIUnit(); e != g {
				t.Errorf("wrong base unit: expected \"%s\", got \"%s\"", e, g)
			}

			if e, g := self.ExpectedPrefix, m.Unit().Prefix().Name; e != g {
				t.Errorf("wrong base prefix: expected \"%s\", got \"%s\"", e, g)
			}

			if e, g := self.Value, m.ConvertToString(self.Unit); math.Abs(e-g) > 1.0e-9 {
				t.Errorf("(\"%s\" [%T]).ConvertToString(\"%s\") = %f, expected %f", m, m, self.Unit, g, e)
			}

			for unit, e := range self.ConvertToString {
				if g := m.ConvertToString(unit); math.Abs(e-g) > 1.0e-9 {
					t.Errorf("(\"%s\" [%T]).ConvertToString(\"%s\") = %f, expected %f", m, m, self.Unit, g, e)
				}
			}

		}

	})
}

type NewMeasurementTests []*NewMeasurementTest

func (self NewMeasurementTests) Run(t *testing.T, constructor MeasurementConstructor) {
	for _, test := range self {
		test.Run(t, constructor)
	}
}

func TestNewLength(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "mm",
			ExpectedSIValue:  1e-3,
			ExpectedBaseUnit: "m",
			ExpectedPrefix:   "m",
		},
		{
			Value:            1.0,
			Unit:             "um",
			ExpectedSIValue:  1e-6,
			ExpectedBaseUnit: "m",
			ExpectedPrefix:   "u",
		},
		{
			Value:            1.0,
			Unit:             "m",
			ExpectedSIValue:  1,
			ExpectedBaseUnit: "m",
			ExpectedPrefix:   " ",
		},
		{
			Value:       1.0,
			Unit:        "parsec",
			ExpectError: true,
		},
	}.Run(t, func(v float64, u string) Measurement {
		return NewLength(v, u)
	})
}

func TestNewArea(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "m^2",
			ExpectedSIValue:  1,
			ExpectedBaseUnit: "m^2",
			ExpectedPrefix:   " ",
		},
		{
			Value:            1.0,
			Unit:             "mm^2",
			ExpectedSIValue:  1e-6,
			ExpectedBaseUnit: "m^2",
			ExpectedPrefix:   "m",
		},
		{
			Value:       1.0,
			Unit:        "hectare",
			ExpectError: true,
		},
	}.Run(t, func(v float64, u string) Measurement {
		return NewArea(v, u)
	})
}

func TestNewVolume(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "ul",
			ExpectedSIValue:  1e-6,
			ExpectedBaseUnit: "l",
			ExpectedPrefix:   "u",
		},
		{
			Value:            1.0,
			Unit:             "µl",
			ExpectedSIValue:  1e-6,
			ExpectedBaseUnit: "l",
			ExpectedPrefix:   "u",
		},
		{
			Value:            1.0,
			Unit:             "ml",
			ExpectedSIValue:  1e-3,
			ExpectedBaseUnit: "l",
			ExpectedPrefix:   "m",
		},
		{
			Value:       1.0,
			Unit:        "decibel",
			ExpectError: true,
		},
	}.Run(t, func(v float64, u string) Measurement {
		return NewVolume(v, u)
	})
}

func TestNewTemperature(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "C",
			ExpectedSIValue:  1.0,
			ExpectedBaseUnit: "℃",
			ExpectedPrefix:   " ",
		},
		{
			Value:            1.0,
			Unit:             "˚C",
			ExpectedSIValue:  1.0,
			ExpectedBaseUnit: "℃",
			ExpectedPrefix:   " ",
		},
		{
			Value:            1.0,
			Unit:             "℃",
			ExpectedSIValue:  1.0,
			ExpectedBaseUnit: "℃",
			ExpectedPrefix:   " ",
		},
		{
			Value:            1.0,
			Unit:             "°C",
			ExpectedSIValue:  1.0,
			ExpectedBaseUnit: "℃",
			ExpectedPrefix:   " ",
		},
		{
			Value:       1.0,
			Unit:        "F",
			ExpectError: true,
		},
	}.Run(t, func(v float64, u string) Measurement {
		return NewTemperature(v, u)
	})
}

func TestNewTime(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "s",
			ExpectedSIValue:  1.0,
			ExpectedBaseUnit: "s",
			ExpectedPrefix:   " ",
		},
		{
			Value:            1.0,
			Unit:             "ms",
			ExpectedSIValue:  1.0e-3,
			ExpectedBaseUnit: "s",
			ExpectedPrefix:   "m",
		},
		{
			Value:            1.0,
			Unit:             "min",
			ExpectedSIValue:  60.0,
			ExpectedBaseUnit: "s",
			ExpectedPrefix:   " ",
		},
		{
			Value:            1.0,
			Unit:             "h",
			ExpectedSIValue:  60.0 * 60.0,
			ExpectedBaseUnit: "s",
			ExpectedPrefix:   " ",
		},
		{
			Value:            1.0,
			Unit:             "days",
			ExpectedSIValue:  60.0 * 60.0 * 24.0,
			ExpectedBaseUnit: "s",
			ExpectedPrefix:   " ",
		},
		{
			Value:       1.0,
			Unit:        "aeon",
			ExpectError: true,
		},
	}.Run(t, func(v float64, u string) Measurement {
		return NewTime(v, u)
	})
}

func TestNewMass(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "kg",
			ExpectedSIValue:  1.0,
			ExpectedBaseUnit: "kg",
			ExpectedPrefix:   "k",
		},
		{
			Value:            1.0,
			Unit:             "ug",
			ExpectedSIValue:  1.0e-9,
			ExpectedBaseUnit: "kg",
			ExpectedPrefix:   "u",
			ConvertToString: map[string]float64{
				"ng": 1000.0,
			},
		},
		{
			Value:       1.0,
			Unit:        "St.Bernard",
			ExpectError: true,
		},
	}.Run(t, func(v float64, u string) Measurement {
		return NewMass(v, u)
	})
}

func TestNewMoles(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "nMol",
			ExpectedSIValue:  1.0e-9,
			ExpectedBaseUnit: "Mol",
			ExpectedPrefix:   "n",
		},
		{
			Value:            1.0,
			Unit:             "uMol",
			ExpectedSIValue:  1.0e-6,
			ExpectedBaseUnit: "Mol",
			ExpectedPrefix:   "u",
		},
		{
			Value:       1.0,
			Unit:        "CommonEuropean",
			ExpectError: true,
		},
	}.Run(t, func(v float64, u string) Measurement {
		return NewMoles(v, u)
	})
}

func TestNewAmount(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "nMol",
			ExpectedSIValue:  1.0e-9,
			ExpectedBaseUnit: "Mol",
			ExpectedPrefix:   "n",
		},
		{
			Value:            1.0,
			Unit:             "uMol",
			ExpectedSIValue:  1.0e-6,
			ExpectedBaseUnit: "Mol",
			ExpectedPrefix:   "u",
		},
		{
			Value:       1.0,
			Unit:        "shedloads",
			ExpectError: true,
		},
	}.Run(t, func(v float64, u string) Measurement {
		return NewAmount(v, u)
	})
}

func TestNewAngle(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "radians",
			ExpectedSIValue:  1.0,
			ExpectedBaseUnit: "rad",
			ExpectedPrefix:   " ",
		},
		{
			Value:            180.0,
			Unit:             "degrees",
			ExpectedSIValue:  math.Pi,
			ExpectedBaseUnit: "rad",
			ExpectedPrefix:   " ",
			ConvertToString: map[string]float64{
				"deg": 180.0,
			},
		},
		{
			Value:       1.0,
			Unit:        "gradians",
			ExpectError: true,
		},
	}.Run(t, func(v float64, u string) Measurement {
		return NewAngle(v, u)
	})
}

func TestNewAnglularVelocity(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "rpm",
			ExpectedSIValue:  math.Pi / 30.0,
			ExpectedBaseUnit: "rad/s",
			ExpectedPrefix:   " ",
			ConvertToString: map[string]float64{
				"rpm": 1.0,
			},
		},
		{
			Value:       1.0,
			Unit:        "pulsar",
			ExpectError: true,
		},
	}.Run(t, func(v float64, u string) Measurement {
		return NewAngularVelocity(v, u)
	})
}

func TestNewEnergy(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "J",
			ExpectedSIValue:  1.0,
			ExpectedBaseUnit: "J",
			ExpectedPrefix:   " ",
		},
		{
			Value:       1.0,
			Unit:        "eV",
			ExpectError: true,
		},
	}.Run(t, func(v float64, u string) Measurement {
		return NewEnergy(v, u)
	})
}

func TestNewForce(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "N",
			ExpectedSIValue:  1.0,
			ExpectedBaseUnit: "N",
			ExpectedPrefix:   " ",
		},
		{
			Value:       1.0,
			Unit:        "DarkSides",
			ExpectError: true,
		},
	}.Run(t, func(v float64, u string) Measurement {
		return NewForce(v, u)
	})
}

func TestNewPressure(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "Pa",
			ExpectedSIValue:  1.0,
			ExpectedBaseUnit: "Pa",
			ExpectedPrefix:   " ",
			ConvertToString: map[string]float64{
				"ubar": 10.0,
			},
		},
		{
			Value:            1.0,
			Unit:             "bar",
			ExpectedSIValue:  100000.0,
			ExpectedBaseUnit: "Pa",
			ExpectedPrefix:   " ",
		},
		{
			Value:       2.7,
			Unit:        "JobInterviews",
			ExpectError: true,
		},
	}.Run(t, func(v float64, u string) Measurement {
		return NewPressure(v, u)
	})
}

func TestNewConcentration(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "kg/l",
			ExpectedSIValue:  1.0,
			ExpectedBaseUnit: "kg/l",
			ExpectedPrefix:   "k",
		},
		{
			Value:            1.0,
			Unit:             "ug/ul",
			ExpectedSIValue:  1.0e-3,
			ExpectedBaseUnit: "kg/l",
			ExpectedPrefix:   "u",
		},
		{
			Value:            1.0,
			Unit:             "ug/ml",
			ExpectedSIValue:  1.0e-6,
			ExpectedBaseUnit: "kg/l",
			ExpectedPrefix:   "u",
		},
		{
			Value:       1.0,
			Unit:        "Eagles",
			ExpectError: true,
		},
	}.Run(t, func(v float64, u string) Measurement {
		return NewConcentration(v, u)
	})
}

func TestNewSpecificHeatCapacity(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "J/kg*C",
			ExpectedSIValue:  1.0,
			ExpectedBaseUnit: "J/kg*C",
			ExpectedPrefix:   " ",
		},
		{
			Value:       1.0,
			Unit:        "J/kg",
			ExpectError: true,
		},
	}.Run(t, func(v float64, u string) Measurement {
		return NewSpecificHeatCapacity(v, u)
	})
}

func TestNewDensity(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "kg/m^3",
			ExpectedSIValue:  1.0,
			ExpectedBaseUnit: "kg/m^3",
			ExpectedPrefix:   "k",
		},
		{
			Value:       1.0,
			Unit:        "ducks",
			ExpectError: true,
		},
	}.Run(t, func(v float64, u string) Measurement {
		return NewDensity(v, u)
	})
}

func TestNewFlowRate(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            60.0,
			Unit:             "ml/min",
			ExpectedSIValue:  0.001,
			ExpectedBaseUnit: "l/s",
			ExpectedPrefix:   "m",
		},
		{
			Value:       1.0,
			Unit:        "TheNile",
			ExpectError: true,
		},
	}.Run(t, func(v float64, u string) Measurement {
		return NewFlowRate(v, u)
	})
}

func TestNewVelocity(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "m/s",
			ExpectedSIValue:  1.0,
			ExpectedBaseUnit: "m/s",
			ExpectedPrefix:   " ",
		},
		{
			Value:       1.0,
			Unit:        "BatsOuttaHell",
			ExpectError: true,
		},
	}.Run(t, func(v float64, u string) Measurement {
		return NewVelocity(v, u)
	})
}

func TestNewRate(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "/s",
			ExpectedSIValue:  1.0,
			ExpectedBaseUnit: "/s",
			ExpectedPrefix:   " ",
		},
		{
			Value:            60.0,
			Unit:             "/min",
			ExpectedSIValue:  1.0,
			ExpectedBaseUnit: "/s",
			ExpectedPrefix:   " ",
		},
		{
			Value:            3600.0,
			Unit:             "/h",
			ExpectedSIValue:  1.0,
			ExpectedBaseUnit: "/s",
			ExpectedPrefix:   " ",
		},
		{
			Value:       1.0,
			Unit:        "/loc",
			ExpectError: true,
		},
	}.Run(t, func(v float64, u string) Measurement {
		if r, err := NewRate(v, u); err != nil {
			panic(err) //other NewX functions panic but this one returns error
		} else {
			return r
		}
	})
}

func TestNewVoltage(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "V",
			ExpectedSIValue:  1.0,
			ExpectedBaseUnit: "V",
			ExpectedPrefix:   " ",
		},
		{
			Value:       1.0,
			Unit:        "NylonShirts",
			ExpectError: true,
		},
	}.Run(t, func(v float64, u string) Measurement {
		if r, err := NewVoltage(v, u); err != nil {
			panic(err)
		} else {
			return r
		}
	})
}

func TestConcentration_MolPerL(t *testing.T) {
	conc := NewConcentration(1.0, "g/l")

	concInMols := conc.MolPerL(2.0)

	if concInMols.Munit.BaseSISymbol() != "Mol/l" {
		t.Errorf("concentration was converted to %s not Mol/l", concInMols.Munit.BaseSISymbol())
	}

	if concInMols.SIValue() != 0.5 {
		t.Errorf("expected concentration of 0.5 Mol/l, got %v", concInMols)
	}
}

func TestConcentration_GramPerL(t *testing.T) {
	conc := NewConcentration(1.0, "M/l")

	concInGrams := conc.GramPerL(2.0)

	if concInGrams.Munit.BaseSISymbol() != "kg/l" {
		t.Errorf("concentration was converted to %s not kg/l", concInGrams.Munit.BaseSISymbol())
	}

	if concInGrams.ConvertToString("g/l") != 2.0 {
		t.Errorf("expected concentration of 2 g/l, got %v", concInGrams)
	}
}
