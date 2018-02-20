package lh

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"strings"
	"testing"
)

func TestPlateSerializeDeserialize(t *testing.T) {
	p := makeplatefortest()

	validatePlate(t, p)

	encodedP := Encodeinterface(p)

	decodedP, err := DecodeGenericPlate(string(encodedP.Arg_1))

	if err != nil {
		t.Errorf(err.Error())
	}

	plate, ok := decodedP.(*wtype.LHPlate)

	if !ok {
		t.Errorf("WANT *wtype.LHPlate got %T", decodedP)
	}

	validatePlate(t, plate)

}

func makeplatefortest() *wtype.LHPlate {
	swshp := wtype.NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype := wtype.NewLHWell("ul", 200, 10, swshp, wtype.VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	p := wtype.NewLHPlate("DSW96", "none", 8, 12, wtype.Coordinates{127.76, 85.48, 43.1}, welltype, 0.5, 0.5, 0.5, 0.5, 0.5)
	return p
}

func validatePlate(t *testing.T, plate *wtype.LHPlate) {
	assertWellsEqual := func(what string, as, bs []*wtype.LHWell) {
		seen := make(map[*wtype.LHWell]int)
		for _, w := range as {
			seen[w] += 1
		}
		for _, w := range bs {
			seen[w] += 1
		}
		for w, count := range seen {
			if count != 2 {
				t.Errorf("%s: no matching well found (%d != %d) for %p %s:%s", what, count, 2, w, w.ID, w.Crds)
			}
		}
	}

	var ws1, ws2, ws3, ws4 []*wtype.LHWell

	for _, w := range plate.HWells {
		ws1 = append(ws1, w)
	}
	for crds, w := range plate.Wellcoords {
		ws2 = append(ws2, w)

		if w.Crds.FormatA1() != crds {
			t.Fatal(fmt.Sprintf("ERROR: Well coords not consistent -- %s != %s", w.Crds.FormatA1(), crds))
		}

		if w.WContents.Loc == "" {
			t.Fatal(fmt.Sprintf("ERROR: Well contents do not have loc set"))
		}

		ltx := strings.Split(w.WContents.Loc, ":")

		if ltx[0] != plate.ID {
			t.Fatal(fmt.Sprintf("ERROR: Plate ID for component not consistent -- %s != %s", ltx[0], plate.ID))
		}

		if w.Plate != nil && ltx[0] != w.Plate.(*wtype.LHPlate).ID {
			t.Fatal(fmt.Sprintf("ERROR: Plate ID for component not consistent with well -- %s != %s", ltx[0], w.Plate.(*wtype.LHPlate).ID))
		}

		if ltx[1] != crds {
			t.Fatal(fmt.Sprintf("ERROR: Coords for component not consistent: -- %s != %s", ltx[1], crds))
		}

	}

	for _, ws := range plate.Rows {
		for _, w := range ws {
			ws3 = append(ws3, w)
		}
	}
	for _, ws := range plate.Cols {
		for _, w := range ws {
			ws4 = append(ws4, w)
		}

	}
	assertWellsEqual("HWells != Rows", ws1, ws2)
	assertWellsEqual("Rows != Cols", ws2, ws3)
	assertWellsEqual("Cols != Wellcoords", ws3, ws4)

	// Check pointer-ID equality
	comp := make(map[string]*wtype.LHComponent)
	for _, w := range append(append(ws1, ws2...), ws3...) {
		c := w.WContents
		if c == nil || c.Vol == 0.0 {
			continue
		}
		if co, seen := comp[c.ID]; seen && co != c {
			t.Errorf("component %s duplicated as %+v and %+v", c.ID, c, co)
		} else if !seen {
			comp[c.ID] = c
		}
	}
}
