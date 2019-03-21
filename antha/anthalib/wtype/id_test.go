package wtype

import (
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"testing"
)

// tests on ID arithmetic

func makeWell() *LHWell {
	swshp := NewShape(BoxShape, "mm", 8.2, 8.2, 41.3)
	welltype := NewLHWell("ul", 1000, 100, swshp, VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	welltype.WContents.Loc = "randomplate:A1"
	return welltype
}

func makeComponent() *Liquid {
	rand := wutil.GetRandom()
	A := NewLHComponent()
	A.Type = LTWater
	A.Smax = 9999
	A.Vol = rand.Float64() * 10.0
	A.Vunit = "ul"
	A.Loc = "anotherrandomplate:A2"
	return A
}

// mix to empty well should result in component
// in well with new ID
// well component should have old component as parent
func TestEmptyWellMix(t *testing.T) {
	c := makeComponent()
	w := makeWell()
	err := w.AddComponent(c)
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
	c := makeComponent()
	w := makeWell()
	idb4 := w.WContents.ID

	err := w.AddComponent(c)
	if err != nil {
		t.Fatal(err)
	}

	if w.WContents.HasParent(w.WContents.ID) {
		t.Fatal("Components should not have themselves as parents! It's just too metaphysical")
	}
	d := makeComponent()

	err = w.AddComponent(d)
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

	e := makeComponent()

	err = w.AddComponent(e)
	if err != nil {
		t.Fatal(err)
	}

	f := makeComponent()

	err = w.AddComponent(f)
	if err != nil {
		t.Fatal(err)
	}

	w2 := makeWell()

	g := makeComponent()

	err = w2.AddComponent(w.WContents)
	if err != nil {
		t.Fatal(err)
	}

	err = w2.AddComponent(g)
	if err != nil {
		t.Fatal(err)
	}
}
