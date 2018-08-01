package wtype

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

func TestSampleBehaviour(t *testing.T) {
	c := NewLHComponent()
	c.CName = "water"

	c2 := NewLHComponent()
	c.CName = "cider"

	c.SetSample(true)

	if !c.IsSample() {
		t.Errorf("SetSample(true) must cause components to return true to IsSample()")
	}

	c2.SetSample(true)
	if !c2.IsSample() {
		t.Errorf("SetSample(true) must cause components to return true to IsSample()")
	}

	c.Mix(c2)

	if c.IsSample() {
		t.Errorf("Results of mixes must not be samples")
	}

	c2.SetSample(false)

	if c2.IsSample() {
		t.Errorf("SetSample(false) must cause components to return false to IsSample()")
	}

	c3 := c2.Dup()

	if c3.IsSample() {
		t.Errorf("Dup()ing a non-sample must produce a non-sample")
	}

	c3.SetSample(true)

	if !c3.IsSample() {
		t.Errorf("SetSample(true) must  cause components to return true to IsSample()... even duplicates")
	}

	if c2.IsSample() {
		t.Errorf("Duplicates must not remain linked")
	}
}

func TestDup(t *testing.T) {
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
		c.AddSubComponents(componentList)

		return c
	}
	someComponents := ComponentList{Components: map[string]wunit.Concentration{
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
		someComponents)

	newLiquid := newTestComponent("new thing",
		LTWater,
		9999,
		wunit.NewConcentration(1, "X"),
		wunit.NewVolume(2000.0, "ul"),
		ComponentList{})

	duplicated := mediaMixture.Dup()

	if err := EqualLists(mediaMixture.SubComponents, duplicated.SubComponents); err != nil {
		t.Error(err.Error())
	}

	duplicated.AddSubComponent(newLiquid, wunit.NewConcentration(0.25, "g/l"))

	if err := EqualLists(mediaMixture.SubComponents, duplicated.SubComponents); err == nil {
		t.Error("expecting lists to no longer be equal but this is not the case")
	}
}

func TestComponentSerialize(t *testing.T) {
	c := NewLHComponent()
	c.CName = "water"

	b, err := json.Marshal(c)

	if err != nil {
		t.Errorf(err.Error())
	}

	c2 := NewLHComponent()

	err = json.Unmarshal(b, &c2)

	if err != nil {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(c, c2) {
		t.Errorf("COMPONENTS NOT EQUAL AFTER MARSHAL/UNMARSHAL")
	}
}
