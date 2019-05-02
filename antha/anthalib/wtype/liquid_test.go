package wtype

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

func TestSampleBehaviour(t *testing.T) {
	c := NewLiquid("water", LTWater, wunit.NewVolume(20.0, "ul"))
	c2 := NewLiquid("cider", LTWater, wunit.NewVolume(20.0, "ul"))

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

func TestLiquidSources(t *testing.T) {
	// let's make some delicious squash

	// first we have some concentrates
	appleConc := NewLiquid("Apple Concentrate", LTWater, wunit.NewVolume(1, "l"))
	berryConc := NewLiquid("Blackberry Concentrate", LTWater, wunit.NewVolume(1, "l"))

	// let's mix together some of each
	appleBerryConc, err := appleConc.Sample(wunit.NewVolume(200.0, "ml"))
	if err != nil {
		t.Fatal(err)
	}

	if bSample, err := berryConc.Sample(wunit.NewVolume(300.0, "ml")); err != nil {
		t.Fatal(err)
	} else {
		appleBerryConc.Mix(bSample)
	}

	// this pre-cursor is meaningful, so give it a name we can refer to later
	appleBerryConc.SetName("Apple and Blackberry Concentrate")

	// now let's get some water
	water := NewLHComponent()
	water.Type = LTWater
	water.SetName("water")
	water.SetVolume(wunit.NewVolume(450.0, "ml"))

	appleBerry, err := appleBerryConc.Sample(wunit.NewVolume(50.0, "ml"))
	if err != nil {
		t.Fatal(err)
	}

	// and use it to dilute the concentrate
	appleBerry.Mix(water)

	// now check that we ended up with delicious squash
	sourceNames := appleBerry.Sources.Names()
	expectedNames := []string{"Apple and Blackberry Concentrate", "water"}
	if !reflect.DeepEqual(sourceNames, expectedNames) {
		t.Fatalf("source name mismatch:\ne: %q\ng: %q", expectedNames, sourceNames)
	}

	// name the squash itself
	appleBerry.SetName("Apple and Blackberry Squash")

	expectedVolumes := map[string]wunit.Volume{
		"Apple and Blackberry Concentrate": wunit.NewVolume(50.0, "ml"),
		"Apple Concentrate":                wunit.NewVolume(20.0, "ml"),
		"Blackberry Concentrate":           wunit.NewVolume(30.0, "ml"),
		"water":                            wunit.NewVolume(450.0, "ml"),
		"Bose-Einstein Condensate":         wunit.NewVolume(0.0, "ml"),
	}
	for name, eVol := range expectedVolumes {
		if gVol := appleBerry.Sources.VolumeOf(name); !eVol.EqualTo(gVol) {
			t.Errorf("wrong volume for %q: expected %s, got %s", name, eVol, gVol)
		}
	}
}

// Mix helper to make the examples look more like element code
func Mix(liquids ...*Liquid) *Liquid {
	ret := NewLHComponent()
	for _, l := range liquids {
		ret.Mix(l)
	}
	return ret
}

// Sample helper to make the examples look more like element code
func Sample(liquid *Liquid, volume wunit.Volume) *Liquid {
	ret, err := liquid.Sample(volume)
	if err != nil {
		panic(err)
	}
	return ret
}

func ExampleLiquid_Sources() {
	// make two liquids
	apple := NewLiquid("Apple Juice", LTWater, wunit.NewVolume(1, "l"))
	berry := NewLiquid("Blackberry Juice", LTWater, wunit.NewVolume(1, "l"))

	// now combine them
	mixture := Mix(apple, berry)

	// sources contains a description of what went in
	fmt.Println(mixture.Sources.Summarize())

	// we can also give the output a specific name
	mixture.SetName("Apple & Blackberry Juice")

	// which will then exist in sources later on
	mixture = Mix(mixture, NewLiquid("Lemonade", LTWater, wunit.NewVolume(0.5, "l")))
	fmt.Println("")
	fmt.Println(mixture.Sources.Summarize())

	// Output:
	// "Apple Juice": 1 l
	// "Blackberry Juice": 1 l
	//
	// "Apple & Blackberry Juice": 2 l
	//   "Apple Juice": 1 l
	//   "Blackberry Juice": 1 l
	// "Lemonade": 0.5 l

}

func ExampleLiquid_Sources_sampling() {

	// make two liquids
	apple := NewLiquid("Apple Juice", LTWater, wunit.NewVolume(1000.0, "ml"))
	berry := NewLiquid("Blackberry Juice", LTWater, wunit.NewVolume(1000.0, "ml"))

	// you can also sample liquids before mixing them
	aSample := Sample(apple, wunit.NewVolume(50.0, "ml"))
	bSample := Sample(berry, wunit.NewVolume(45.0, "ml"))

	// combine as before
	mixture := Mix(aSample, bSample)
	mixture.SetName("Apple & Blackberry Juice")
	mixture = Mix(mixture, NewLiquid("Water", LTWater, wunit.NewVolume(100, "ml")))
	fmt.Println(mixture.Sources.Summarize())

	// Output:
	// "Apple & Blackberry Juice": 95 ml
	//   "Apple Juice": 50 ml
	//   "Blackberry Juice": 45 ml
	// "Water": 100 ml

}

func ExampleLiquid_Sources_transfers() {
	// multiple mixes of the same thing are combined unless a new name is set
	apple := NewLiquid("Apple Juice", LTWater, wunit.NewVolume(1, "l"))
	berry := NewLiquid("Blackberry Juice", LTWater, wunit.NewVolume(1, "l"))

	samples := []*Liquid{Sample(apple, wunit.NewVolume(10, "ml"))}
	for i := 0; i < 10; i++ {
		samples = append(samples, Sample(berry, wunit.NewVolume(1, "ml")))
	}

	mixture := Mix(samples...)
	fmt.Println(mixture.Sources.Summarize())

	// Output:
	// "Apple Juice": 10 ml
	// "Blackberry Juice": 10 ml
}

func ExampleLiquid_Sources_transfers2() {
	// multiple mixes of the same thing are ignored unless a new name is set
	apple := NewLiquid("Apple Juice", LTWater, wunit.NewVolume(1, "l"))
	berry := NewLiquid("Blackberry Juice", LTWater, wunit.NewVolume(1, "l"))

	mixture := Sample(apple, wunit.NewVolume(10, "ml"))
	for i := 0; i < 10; i++ {
		mixture = Mix(mixture, Sample(berry, wunit.NewVolume(1, "ml")))
	}

	fmt.Println(mixture.Sources.Summarize())

	// Output:
	// "Apple Juice": 10 ml
	// "Blackberry Juice": 10 ml
}

func ExampleLiquid_Sources_naming() {
	// Liquids with the same name are assumed to be the same for the purposes of tracking history.
	//
	// For example if we make two liquids containing different things:
	mixture1 := Mix(NewLiquid("Apple Juice", LTWater, wunit.NewVolume(10.0, "ml")), NewLiquid("Blackberry Juice", LTWater, wunit.NewVolume(10.0, "ml")))
	mixture2 := Mix(NewLiquid("Pear Juice", LTWater, wunit.NewVolume(10.0, "ml")), NewLiquid("Guava Juice", LTWater, wunit.NewVolume(10.0, "ml")))

	// but give them the same name
	mixture1.SetName("Fruit Juice")
	mixture2.SetName("Fruit Juice")

	// then combining them combines their sources
	juice := Mix(mixture1, mixture2)

	fmt.Println(juice.Sources.Summarize())

	// Output:
	// "Fruit Juice": 40 ml
	//   "Apple Juice": 10 ml
	//   "Blackberry Juice": 10 ml
	//   "Guava Juice": 10 ml
	//   "Pear Juice": 10 ml

}
