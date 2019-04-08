package wtype

import (
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/laboratory/effects/id"
)

// tests on ID arithmetic

func makeWell(idGen *id.IDGenerator) *LHWell {
	swshp := NewShape(BoxShape, "mm", 8.2, 8.2, 41.3)
	welltype := NewLHWell(idGen, "ul", 1000, 100, swshp, VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	welltype.WContents.Loc = "randomplate:A1"
	return welltype
}

func makeComponent(idGen *id.IDGenerator) *Liquid {
	rand := wutil.GetRandom()
	liq := NewLHComponent(idGen)
	liq.Type = LTWater
	liq.Smax = 9999
	liq.Vol = rand.Float64() * 10.0
	liq.Vunit = "ul"
	liq.Loc = "anotherrandomplate:A2"
	return liq
}

// mix to empty well should result in component
// in well with new ID
// well component should have old component as parent
func TestEmptyWellMix(t *testing.T) {
	idGen := id.NewIDGenerator("testing")
	c := makeComponent(idGen)
	w := makeWell(idGen)
	err := w.AddComponent(idGen, c)
	if err != nil {
		t.Fatal(err)
	}

	if w.WContents.ID == c.ID {
		t.Fatal("Well contents should have different ID to input component")
	}

	if !w.WContents.HasParent(c.ID) {
		t.Fatal("Well contents should have input component ID as parent")
	}
	if !c.HasDaughter(w.WContents.ID) {
		t.Fatal("Component mixed into well should have well contents as daughter")
	}
}

func TestFullWellMix(t *testing.T) {
	idGen := id.NewIDGenerator("testing")
	c := makeComponent(idGen)
	w := makeWell(idGen)
	idb4 := w.WContents.ID

	err := w.AddComponent(idGen, c)
	if err != nil {
		t.Fatal(err)
	}

	if w.WContents.HasParent(w.WContents.ID) {
		t.Fatal("Components should not have themselves as parents! It's just too metaphysical")
	}
	d := makeComponent(idGen)

	err = w.AddComponent(idGen, d)
	if err != nil {
		t.Fatal(err)
	}

	if w.WContents.ID == c.ID || w.WContents.ID == d.ID || w.WContents.ID == idb4 {
		t.Fatal("Well contents should have new ID after mix")
	}

	if !w.WContents.HasParent(d.ID) {
		t.Fatal("Well contents should have last parent set")
	}

	if !d.HasDaughter(w.WContents.ID) {
		t.Fatal("Component mixed into well should have well contents as daughter")
	}

	if w.WContents.HasParent(w.WContents.ID) {
		t.Fatal("Components should not have themselves as parents! It's just too metaphysical")
	}

	e := makeComponent(idGen)

	err = w.AddComponent(idGen, e)
	if err != nil {
		t.Fatal(err)
	}

	f := makeComponent(idGen)

	err = w.AddComponent(idGen, f)
	if err != nil {
		t.Fatal(err)
	}

	w2 := makeWell(idGen)

	g := makeComponent(idGen)

	err = w2.AddComponent(idGen, w.WContents)
	if err != nil {
		t.Fatal(err)
	}

	err = w2.AddComponent(idGen, g)
	if err != nil {
		t.Fatal(err)
	}
}
