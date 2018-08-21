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
	}.Run(t, func(v float64, u string) Measurement {
		return NewVolume(v, u)
	})
}
