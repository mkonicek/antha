package wtype

import (
	"fmt"
	"strings"
	"testing"
)

// platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64

func makeplatefortest() *LHPlate {
	swshp := NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype := NewLHWell("DSW96", "", "", "ul", 1000, 100, swshp, LHWBV, 8.2, 8.2, 41.3, 4.7, "mm")
	p := NewLHPlate("testplate", "none", 8, 12, 44.1, "mm", welltype, 0.5, 0.5, 0.5, 0.5, 0.5)
	return p
}

func TestPlateCreation(t *testing.T) {
	p := makeplatefortest()
	validatePlate(t, p)
}

func TestPlateDup(t *testing.T) {
	p := makeplatefortest()
	d := p.Dup()
	validatePlate(t, d)
}

func validatePlate(t *testing.T, plate *LHPlate) {
	assertWellsEqual := func(what string, as, bs []*LHWell) {
		seen := make(map[*LHWell]int)
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

	var ws1, ws2, ws3, ws4 []*LHWell

	for _, w := range plate.HWells {
		ws1 = append(ws1, w)
	}
	for crds, w := range plate.Wellcoords {
		ws2 = append(ws2, w)

		if w.Crds != crds {
			t.Fatal(fmt.Sprintf("ERROR: Well coords not consistent -- %s != %s", w.Crds, crds))
		}

		if w.WContents.Loc == "" {
			t.Fatal(fmt.Sprintf("ERROR: Well contents do not have loc set"))
		}

		ltx := strings.Split(w.WContents.Loc, ":")

		if ltx[0] != plate.ID {
			t.Fatal(fmt.Sprintf("ERROR: Plate ID for component not consistent -- %s != %s", ltx[0], plate.ID))
		}

		if ltx[0] != w.Plateid {
			t.Fatal(fmt.Sprintf("ERROR: Plate ID for component not consistent with well -- %s != %s", ltx[0], w.Plateid))
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
	assertWellsEqual("Cols != WellCoords", ws3, ws4)

	// Check pointer-ID equality
	comp := make(map[string]*LHComponent)
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
