package wunit

import (
	"fmt"
	"testing"
)

type concConversionTest struct {
	StockConc    Concentration
	TargetConc   Concentration
	TotalVolume  Volume
	VolumeNeeded Volume
	ShouldError  bool
}

func (test *concConversionTest) unexpectedError(err error) bool {
	return (err != nil) != test.ShouldError
}

func (test *concConversionTest) Run(t *testing.T) {
	t.Run(fmt.Sprintf("(%v, %v, %v)", test.TargetConc, test.StockConc, test.TotalVolume), func(t *testing.T) {
		vol, err := VolumeForTargetConcentration(test.TargetConc, test.StockConc, test.TotalVolume)

		if test.unexpectedError(err) {
			t.Errorf("expecting error %t, got error %v", test.ShouldError, err)
		}

		if !test.ShouldError && !vol.EqualTo(test.VolumeNeeded) {
			t.Errorf("expected volume %v, got volume %v", test.VolumeNeeded, vol)
		}
	})
}

type concConversionTests []concConversionTest

func (self concConversionTests) Run(t *testing.T) {
	for _, test := range self {
		test.Run(t)
	}
}

type massConversionTest struct {
	Conc      Concentration
	Vol       Volume
	Mass      Mass
	MassError bool
	VolError  bool
}

func (test *massConversionTest) TestMassForTargetConcentration(t *testing.T) {
	t.Run(fmt.Sprintf("MassForTargetConcentration(%v, %v)", test.Conc, test.Vol), func(t *testing.T) {
		if mass, err := MassForTargetConcentration(test.Conc, test.Vol); (err != nil) != test.MassError {
			t.Errorf("expecting error %t, got error %v", test.MassError, err)
		} else if !test.MassError && !mass.EqualTo(test.Mass) {
			t.Errorf("wrong result: expected %v, got %v", test.Mass, mass)
		}
	})
}

func (test *massConversionTest) TestVolumeForTargetMass(t *testing.T) {
	t.Run(fmt.Sprintf("VolumeForTargetMass(%v, %v)", test.Conc, test.Vol), func(t *testing.T) {
		t.Run(fmt.Sprintf("(%v, %v)", test.Mass, test.Conc), func(t *testing.T) {
			if vol, err := VolumeForTargetMass(test.Mass, test.Conc); (err != nil) != test.VolError {
				t.Errorf("expecting error: %t, got error: %v", test.VolError, err)
			} else if !test.VolError && !test.Vol.EqualToTolerance(vol, 1.0e-9) {
				t.Errorf("expected vol: %v, got vol: %v", test.Vol, vol)
			}
		})
	})
}

type massConversionTests []massConversionTest

func (self massConversionTests) Run(t *testing.T) {
	for _, test := range self {
		test.TestMassForTargetConcentration(t)
	}
	for _, test := range self {
		test.TestVolumeForTargetMass(t)
	}
}

type densityConversionTest struct {
	Density     Density
	Vol         Volume
	Mass        Mass
	ShouldError bool
}

func (test *densityConversionTest) unexpectedError(err error) bool {
	return (err != nil) != test.ShouldError
}

func (test *densityConversionTest) TestMassToVolume(t *testing.T) {
	t.Run(fmt.Sprintf("(%v, %v)", test.Density, test.Vol), func(t *testing.T) {
		if vol, err := MassToVolume(test.Mass, test.Density); test.unexpectedError(err) {
			t.Errorf("expecting error %t, got error %v", test.ShouldError, err)
		} else if !vol.EqualTo(test.Vol) {
			t.Error(
				"Expected vol:", test.Vol.ToString(), "\n",
				"Got vol:", vol.ToString(), "\n",
			)
		}
	})
}

func (test *densityConversionTest) TestVolumeToMass(t *testing.T) {
	t.Run(fmt.Sprintf("(%v, %v)", test.Density, test.Mass), func(t *testing.T) {
		if mass, err := VolumeToMass(test.Vol, test.Density); test.unexpectedError(err) {
			t.Errorf("expecting error %t, got error %v", test.ShouldError, err)
		} else if !mass.EqualTo(test.Mass) {
			t.Error(
				"for", fmt.Sprintf("%+v", test), "\n",
				"Expected mass:", test.Mass.ToString(), "\n",
				"Got mass:", mass.ToString(), "\n",
			)
		}

	})
}

type densityConversionTests []densityConversionTest

func (self densityConversionTests) Run(t *testing.T) {
	for _, test := range self {
		test.TestMassToVolume(t)
	}
	for _, test := range self {
		test.TestVolumeToMass(t)
	}
}

var (
	// Some Concentrations

	x2   = NewConcentration(2, "X")
	x100 = NewConcentration(100, "X")

	kgPerL01 = NewConcentration(0.1, "kg/L")
	gPerL1   = NewConcentration(1, "g/L")
	gPerL100 = NewConcentration(100, "g/L")

	mPerL01  = NewConcentration(0.1, "M/L")
	mMPerL1  = NewConcentration(1, "mM/L")
	uMPerL10 = NewConcentration(10, "uM")

	// Some volumes
	ul1   = NewVolume(1, "ul")
	ul100 = NewVolume(100, "ul")
)

func TestVolumeForTargetConcentration(t *testing.T) {

	concConversionTests{
		{
			StockConc:    kgPerL01,
			TargetConc:   gPerL1,
			TotalVolume:  ul100,
			VolumeNeeded: ul1,
		},
		{
			StockConc:    gPerL100,
			TargetConc:   gPerL1,
			TotalVolume:  ul100,
			VolumeNeeded: ul1,
		},
		{
			StockConc:    x100,
			TargetConc:   x2,
			TotalVolume:  ul100,
			VolumeNeeded: NewVolume(2.0, "ul"),
		},
		{
			StockConc:    mMPerL1,
			TargetConc:   uMPerL10,
			TotalVolume:  ul100,
			VolumeNeeded: NewVolume(1.0, "ul"),
		},
		{
			StockConc:    mPerL01,
			TargetConc:   mMPerL1,
			TotalVolume:  ul100,
			VolumeNeeded: NewVolume(1.0, "ul"),
		},
		{
			StockConc:   NewConcentration(1.0, "g/l"),
			TargetConc:  NewConcentration(2.0, "g/l"),
			TotalVolume: NewVolume(100.0, "ul"),
			ShouldError: true,
		},
		{
			StockConc:    NewConcentration(1.0, "g/l"),
			TargetConc:   NewConcentration(2.0, "g/l"),
			TotalVolume:  NewVolume(0.0, "ul"),
			VolumeNeeded: NewVolume(0.0, "ul"),
		},
		{
			StockConc:   NewConcentration(1.0, "g/l"),
			TargetConc:  NewConcentration(0.5, "X"),
			TotalVolume: NewVolume(100.0, "ul"),
			ShouldError: true,
		},
	}.Run(t)
}

func TestMassConversion(t *testing.T) {

	massConversionTests{
		{
			Conc: NewConcentration(1.0, "g/L"),
			Mass: NewMass(1000.0, "mg"),
			Vol:  NewVolume(1.0, "l"),
		},
		{
			Conc: NewConcentration(1.0, "kg/L"),
			Mass: NewMass(1.0, "kg"),
			Vol:  NewVolume(1.0, "l"),
		},
		{
			Conc: NewConcentration(0.1, "mg/L"),
			Mass: NewMass(0.1, "mg"),
			Vol:  NewVolume(1.0, "l"),
		},
		{
			Conc: NewConcentration(100, "ng/ul"),
			Mass: NewMass(100, "ng"),
			Vol:  NewVolume(1.0, "ul"),
		},
		{
			Conc:     NewConcentration(0, "g/l"),
			Mass:     NewMass(0, "g"),
			Vol:      NewVolume(1.0, "l"),
			VolError: true,
		},
		{
			Conc: NewConcentration(100, "ng/ul"),
			Mass: NewMass(0, "ng"),
			Vol:  NewVolume(0, "ul"),
		},
		{
			Conc:      NewConcentration(100, "X"),
			Mass:      NewMass(100, "ng"),
			Vol:       NewVolume(1.0, "ul"),
			VolError:  true,
			MassError: true,
		},
	}.Run(t)
}

func TestDensityConversion(t *testing.T) {
	densityConversionTests{
		{
			Density: NewDensity(1.0, "kg/m^3"),
			Mass:    NewMass(1.0, "kg"),
			Vol:     NewVolume(1000, "l"),
		},
		{
			Density: NewDensity(1000.0, "kg/m^3"),
			Mass:    NewMass(1.0, "g"),
			Vol:     NewVolume(1, "ml"),
		},
		{
			Density: NewDensity(1000.0, "kg/m^3"),
			Mass:    NewMass(0.0, "g"),
			Vol:     NewVolume(0, "l"),
		},
	}.Run(t)
}
