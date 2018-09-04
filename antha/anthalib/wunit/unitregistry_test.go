package wunit

import (
	"testing"
)

func TestUnitRegistry(t *testing.T) {
	reg := NewUnitRegistry()

	if err := reg.DeclareUnit("distance", "meters", "m", "m", SIPrefixes, 1); err != nil {
		t.Fatal(err)
	}

	//can't redeclare a unit
	if err := reg.DeclareUnit("distance", "meters", "m", "m", SIPrefixes, 1); err == nil {
		t.Fatal("redeclaration of meter got no error")
	}

	//can retrieve the correct unit
	if millimeters, err := reg.GetUnit("mm"); err != nil {
		t.Error(err)
	} else if e, g := "millimeters[mm]", millimeters.String(); e != g {
		t.Fatalf("reg.GetUnit(\"mm\") returned %q, not %q", g, e)
	}

	//can't retrieve a unit that doesn't exist
	if _, err := reg.GetUnit("miles"); err == nil {
		t.Error("got no error getting a unit that doesn't exist")
	}

	//can't declare a derived unit across measurement types
	if err := reg.DeclareDerivedUnit("volume", "meters cubed", "m^3", nil, 1, "m", 1.0); err == nil {
		t.Error("declaring derived unit across measurement types got no error")
	}

	if err := reg.DeclareDerivedUnit("distance", "yards", "yards", nil, 1, "m", 0.9144); err != nil {
		t.Error(err)
	} else if err := reg.DeclareDerivedUnit("distance", "miles", "miles", nil, 1, "km", 1609.34); err != nil {
		t.Error(err)
	}

	//can't declare an alias if it already exists
	if err := reg.DeclareAlias("distance", "yards", "miles", nil); err == nil {
		t.Error("declared an alias which shadowed a symbol without error")
	}

	if err := reg.DeclareAlias("distance", "parsec", "Parsec", nil); err == nil {
		t.Error("declared an alias for a symbol that doesn't exist without error")
	}

	if err := reg.DeclareAlias("volume", "m^3", "m", nil); err == nil {
		t.Error("declaring an alias across measurement types got no error")
	}

	//successful aliassing
	if err := reg.DeclareAlias("distance", "y", "yards", nil); err != nil {
		t.Error(err)
	} else if unit, err := reg.GetUnit("y"); err != nil {
		t.Error(err)
	} else if unit.symbol != "yards" {
		t.Errorf("alias \"y\" mapped to %v, expected yards[yards]", unit)
	}

	//cannot redeclare alias
	if err := reg.DeclareAlias("distance", "y", "yards", nil); err == nil {
		t.Error("redeclared alias with no error")
	}

	//simple conversion
	if measurement, err := reg.NewMeasurement(20.0, "km"); err != nil {
		t.Error(err)
	} else if meters, err := reg.GetUnit("m"); err != nil {
		t.Error(err)
	} else if e, g := "meters[m]", meters.String(); e != g {
		t.Fatalf("reg.GetUnit(\"m\") returned %q, not %q", g, e)
	} else if measurementInMeters, err := measurement.InUnit(meters); err != nil {
		t.Error(err)
	} else if e, g := 20000.0, measurementInMeters.RawValue(); e != g {
		t.Errorf("converting 20 km to %v: got %f", meters, g)
	}

	//simple addition
	if a, err := reg.NewMeasurement(10, "m"); err != nil {
		t.Error(err)
	} else if b, err := reg.NewMeasurement(50.0, "cm"); err != nil {
		t.Error(err)
	} else if cm, err := reg.GetUnit("cm"); err != nil {
		t.Error(err)
	} else {
		a.Add(b)
		if aInCm, err := a.InUnit(cm); err != nil {
			t.Error(err)
		} else if aInCm.RawValue() != 1050.0 {
			t.Errorf("added 50 cm to 10 m and expected 1050 cm, but got %v", aInCm)
		}
	}

	//simple subtraction
	if a, err := reg.NewMeasurement(10, "m"); err != nil {
		t.Error(err)
	} else if b, err := reg.NewMeasurement(50.0, "cm"); err != nil {
		t.Error(err)
	} else if cm, err := reg.GetUnit("cm"); err != nil {
		t.Error(err)
	} else {
		a.Subtract(b)
		if aInCm, err := a.InUnit(cm); err != nil {
			t.Error(err)
		} else if aInCm.RawValue() != 950.0 {
			t.Errorf("subtracted 50 cm from 10 m and expected 950 cm, but got %v", aInCm)
		}
	}

}
