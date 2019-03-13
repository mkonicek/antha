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
		if err := c.AddSubComponents(componentList); err != nil {
			t.Fatal(err)
		}

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

	if err := duplicated.AddSubComponent(newLiquid, wunit.NewConcentration(0.25, "g/l")); err != nil {
		t.Fatal(err)
	}

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

func TestDeepCopySubComponents(t *testing.T) {
	l := NewLHComponent()
	l.CName = "water"
	sc := NewLHComponent()
	sc.CName = "mush"
	if err := l.AddSubComponent(sc, wunit.NewConcentration(50.0, "g/l")); err != nil {
		t.Fatal(err)
	}

	l2 := l.Dup()

	scMapEqual := func(m1, m2 map[string]wunit.Concentration) bool {
		if len(m1) != len(m2) {
			return false
		}

		for k, v := range m1 {
			v2, ok := m2[k]

			if !ok || !v.EqualTo(v2) {
				return false
			}
		}

		return true
	}

	if !scMapEqual(l.SubComponents.Components, l2.SubComponents.Components) {
		t.Errorf("Subcomponent Maps not identical in contents after Dup")
	}

	l.SubComponents.Components["fishpaste"] = wunit.NewConcentration(100.0, "g/l")

	if scMapEqual(l.SubComponents.Components, l2.SubComponents.Components) {
		t.Errorf("Subcomponent Maps still linked after dup")
	}

}

func TestEqualTypeVolume(t *testing.T) {
	l := &Liquid{}

	if !l.EqualTypeVolume(l) {
		t.Errorf("Liquid must equal itself")
	}

	l2 := &Liquid{
		CName: "which",
	}

	if l2.EqualTypeVolume(l) {
		t.Errorf("Liquids with non-equal types must not report equal types")
	}

	l3 := &Liquid{
		Vol:   50.0,
		Vunit: "ul",
	}

	if l3.EqualTypeVolume(l) {
		t.Errorf("Liquids with non-equal volumes must not report equal volumes")
	}

	l4 := &Liquid{
		CName: "a",
		Vol:   50.0,
		Vunit: "ul",
	}

	l5 := l4.Dup()

	if !l5.EqualTypeVolume(l4) {
		t.Errorf("Liquids must be of equal type and volume after Dup()")
	}

	l6 := &Liquid{
		CName: "b",
		Vol:   50.0,
		Vunit: "ul",
		ID:    "thisismyID",
	}

	l7 := l6.Dup()

	if !l7.EqualTypeVolumeID(l6) {
		t.Errorf("Liquids must preserve IDs after Dup()")
	}

	l8 := l7.Cp()

	if l8.EqualTypeVolumeID(l7) {
		t.Errorf("Liquids must not preserve IDs after Cp()")
	}

}
