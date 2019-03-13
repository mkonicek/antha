package wtype

import (
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

func TestAddandGetComponent(t *testing.T) {
	newTestComponent := func(
		name string,
		typ LiquidType,
		smax float64,
		conc wunit.Concentration,
		vol wunit.Volume,
		componentList ComponentList,
	) *Liquid {
		c := NewLHComponent()
		c.SetName(name)
		c.Type = typ
		c.Smax = smax
		c.SetConcentration(conc)
		if err := c.AddSubComponents(componentList); err != nil {
			t.Fatal(err)
		}

		return c
	}
	someComponents := ComponentList{
		Components: map[string]wunit.Concentration{
			"glycerol": wunit.NewConcentration(0.25, "g/l"),
			"IPTG":     wunit.NewConcentration(0.25, "mM/l"),
			"water":    wunit.NewConcentration(0.25, "v/v"),
			"LB":       wunit.NewConcentration(0.25, "X"),
		},
	}
	mediaMixture := newTestComponent("LB",
		LTWater,
		9999,
		wunit.NewConcentration(1, "X"),
		wunit.NewVolume(2000.0, "ul"),
		ComponentList{})

	if err := mediaMixture.AddSubComponents(someComponents); err != nil {
		t.Error(err)
	}

	tests := ComponentList{
		Components: map[string]wunit.Concentration{
			"Glycerol":  wunit.NewConcentration(0.25, "g/l"),
			"GLYCEROL ": wunit.NewConcentration(0.25, "g/l"),
		},
	}

	if err := mediaMixture.AddSubComponents(tests); err == nil {
		t.Errorf("expected error adding equalfold sub components to liquid but no error reported. New Sub components: %v", mediaMixture.SubComponents.AllComponents())
	}

	for test := range tests.Components {
		if !mediaMixture.HasSubComponent(test) {
			t.Errorf(
				"Expected sub component %s to be found but only found these: %+v",
				test,
				mediaMixture.SubComponents.AllComponents(),
			)
		}
	}
}
