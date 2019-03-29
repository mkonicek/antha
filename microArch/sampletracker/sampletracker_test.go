package sampletracker

import (
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func assertLocation(t *testing.T, st *SampleTracker, id string, eloc string, eok bool) {
	if loc, ok := st.GetLocationOf(id); ok != eok || eloc != loc {
		t.Errorf("GetLocationOf(%q) returned %q, %t - expected %q, %t", id, loc, ok, eloc, eok)
	}
}

func getPlateForTest() *wtype.Plate {
	cone := wtype.NewShape(wtype.CylinderShape, "mm", 5.5, 5.5, 20.4)
	welltype := wtype.NewLHWell("ul", 200, 5, cone, wtype.UWellBottom, 5.5, 5.5, 20.4, 1.4, "mm")
	return wtype.NewLHPlate("pcrplate_skirted_riser", "Unknown", 8, 12, wtype.Coordinates3D{X: 127.76, Y: 85.48, Z: 25.7}, welltype, 9, 9, 0.0, 0.0, 38.5)
}

func getLiquidForTest(name string, volume float64) *wtype.Liquid {
	ret := wtype.NewLHComponent()
	ret.CName = name
	ret.Vol = volume

	return ret
}

func TestLocations(t *testing.T) {
	st := NewSampleTracker()

	positions := map[string]string{
		"Who":          "First Base",
		"What":         "Second base",
		"I Don't Know": "Third base",
		"Why":          "Left field",
		"Because":      "Center field",
		"Tomorrow":     "Pitcher",
		"Today":        "Catcher",
		"I Don't Care": "Shortstop",
	}

	for id, pos := range positions {
		st.SetLocationOf(id, pos)
	}

	for id, epos := range positions {
		assertLocation(t, st, id, epos, true)
	}
}

func TestUpdateIDOfBeforeLocation(t *testing.T) {
	st := NewSampleTracker()

	st.UpdateIDOf("oldID", "newID")
	assertLocation(t, st, "newID", "", false)

	st.SetLocationOf("oldID", "Location")
	assertLocation(t, st, "newID", "Location", true)
}

func TestUpdateIDOfAfterLocation(t *testing.T) {
	st := NewSampleTracker()

	st.SetLocationOf("oldID", "Location")
	assertLocation(t, st, "oldID", "Location", true)

	assertLocation(t, st, "newID", "", false)
	st.UpdateIDOf("oldID", "newID")

	assertLocation(t, st, "newID", "Location", true)
	assertLocation(t, st, "oldID", "Location", true)
}

func TestInputPlates(t *testing.T) {
	p := getPlateForTest()

	for it := wtype.NewAddressIterator(p, wtype.ColumnWise, wtype.TopToBottom, wtype.LeftToRight, false); it.Valid(); it.Next() {
		well := p.GetChildByAddress(it.Curr()).(*wtype.LHWell)
		cmp := getLiquidForTest("water", 100.0)
		cmp.ID = it.Curr().FormatA1()

		if err := well.SetContents(cmp); err != nil {
			t.Fatal(err)
		}
	}

	st := NewSampleTracker()

	if got := len(st.GetInputPlates()); got != 0 {
		t.Errorf("new sample tracker had %d input plates", got)
	}

	st.SetInputPlate(p)

	if inputs := st.GetInputPlates(); len(inputs) != 1 {
		t.Errorf("expected one input plate, got %d", len(inputs))
	} else if inputs[0] != p {
		t.Errorf("input plate didn't match %p != %p", inputs[0], p)
	} else {
		for it := wtype.NewAddressIterator(p, wtype.ColumnWise, wtype.TopToBottom, wtype.LeftToRight, false); it.Valid(); it.Next() {
			well := inputs[0].GetChildByAddress(it.Curr()).(*wtype.LHWell)
			if !well.IsUserAllocated() {
				t.Errorf("Well %s wasn't marked as user allocated", well)
			}

			cmp := well.Contents()
			assertLocation(t, st, cmp.ID, cmp.Loc, true)
		}
	}

}
