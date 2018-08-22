package wunit

import (
	"fmt"
	"testing"
)

type MeasurementConstructor func(float64, string) Measurement

type NewMeasurementTest struct {
	Value            float64
	Unit             string
	ExpectedSIValue  float64
	ExpectedBaseUnit string
	ExpectedPrefix   string
}

func (self *NewMeasurementTest) Run(t *testing.T, constructor MeasurementConstructor) {

	t.Run(fmt.Sprintf("%f_%s", self.Value, self.Unit), func(t *testing.T) {
		m := constructor(self.Value, self.Unit)

		if m.SIValue() != self.ExpectedSIValue {
			t.Errorf("wrong SIValue: expected %e, got %e", self.ExpectedSIValue, m.SIValue())
		}

		if e, g := self.ExpectedBaseUnit, m.Unit().BaseSIUnit(); e != g {
			t.Errorf("wrong base unit: expected \"%s\", got \"%s\"", e, g)
		}

		if e, g := self.ExpectedPrefix, m.Unit().Prefix().Name; e != g {
			t.Errorf("wrong base prefix: expected \"%s\", got \"%s\"", e, g)
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
		},
	}.Run(t, func(v float64, u string) Measurement {
		return NewMass(v, u)
	})
}

func TestNewMoles(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "uM",
			ExpectedSIValue:  1.0e-6,
			ExpectedBaseUnit: "M",
			ExpectedPrefix:   "u",
		},
		{
			Value:            1.0,
			Unit:             "uMol",
			ExpectedSIValue:  1.0e-6,
			ExpectedBaseUnit: "M",
			ExpectedPrefix:   "u",
		},
	}.Run(t, func(v float64, u string) Measurement {
		return NewMoles(v, u)
	})
}

func TestNewAmount(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "uM",
			ExpectedSIValue:  1.0e-6,
			ExpectedBaseUnit: "M",
			ExpectedPrefix:   "u",
		},
		{
			Value:            1.0,
			Unit:             "uMol",
			ExpectedSIValue:  1.0e-6,
			ExpectedBaseUnit: "M",
			ExpectedPrefix:   "u",
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
			ExpectedBaseUnit: "radians",
			ExpectedPrefix:   " ",
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
			ExpectedSIValue:  1.0,
			ExpectedBaseUnit: "rpm",
			ExpectedPrefix:   " ",
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
			ExpectedPrefix:   " ",
		},
		{
			Value:            1.0,
			Unit:             "ug/ml",
			ExpectedSIValue:  1.0e-6,
			ExpectedBaseUnit: "kg/l",
			ExpectedPrefix:   "m",
		},
	}.Run(t, func(v float64, u string) Measurement {
		return NewConcentration(v, u)
	})
}

func TestNewSpecificHeatCapacity(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "J/kg",
			ExpectedSIValue:  1.0,
			ExpectedBaseUnit: "J/kg",
			ExpectedPrefix:   " ",
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
			ExpectedPrefix:   " ",
		},
	}.Run(t, func(v float64, u string) Measurement {
		return NewDensity(v, u)
	})
}

func TestNewFlowRate(t *testing.T) {
	NewMeasurementTests{
		{
			Value:            1.0,
			Unit:             "ml/min",
			ExpectedSIValue:  1.0,
			ExpectedBaseUnit: "ml/min",
			ExpectedPrefix:   " ",
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
	}.Run(t, func(v float64, u string) Measurement {
		r, _ := NewRate(v, u)
		return r
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
	}.Run(t, func(v float64, u string) Measurement {
		r, _ := NewVoltage(v, u)
		return r
	})
}
