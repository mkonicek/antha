package lh

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/stretchr/testify/assert"
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

func TestTipboxSerializeDeserialize(t *testing.T) {
	tb := makeTipboxForTest()

	//validateTipbox(t, tb)

	encodedTB := Encodeinterface(tb)

	decodedTB, err := DecodeGenericPlate(encodedTB.GetArg_1())
	if err != nil {
		t.Error(err)
	}

	tb2, ok := decodedTB.(*wtype.LHTipbox)
	if !ok {
		t.Errorf("Expected *wtype.LHTipbox, got %T", decodedTB)
	}

	assertTipBoxesEqual(t, tb, tb2)

}

func makeplatefortest() *wtype.LHPlate {
	swshp := wtype.NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype := wtype.NewLHWell("ul", 200, 10, swshp, wtype.VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	p := wtype.NewLHPlate("DSW96", "none", 8, 12, wtype.Coordinates{127.76, 85.48, 43.1}, welltype, 0.5, 0.5, 0.5, 0.5, 0.5)
	p.Welltype.SetWellTargets("Magick", []wtype.Coordinates{{0.0, -10.0, 0.0}, {0.0, 10.0, 0.0}})
	return p
}

func makeTipboxForTest() *wtype.LHTipbox {

	shp := wtype.NewShape("cylinder", "mm", 7.3, 7.3, 51.2)
	w := wtype.NewLHWell("ul", 20.0, 1.0, shp, 0, 7.3, 7.3, 46.0, 0.0, "mm")
	w.Extra["InnerL"] = 5.5
	w.Extra["InnerW"] = 5.5
	w.Extra["Tipeffectiveheight"] = 34.6
	tip := wtype.NewLHTip("gilson", "Gilson20", 0.5, 20.0, "ul", false, shp)
	tb := wtype.NewLHTipbox(8, 12, wtype.Coordinates{127.76, 85.48, 60.13}, "Gilson", "DL10 Tip Rack (PIPETMAX 8x20)", tip, w, 9.0, 9.0, 0.0, 0.0, 28.93)

	return tb

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
				t.Errorf("%s: no matching well found (%d != %d) for %p %s:%s", what, count, 2, w, w.ID, w.Crds.FormatA1())
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
		ws3 = append(ws3, ws...)
	}
	for _, ws := range plate.Cols {
		ws4 = append(ws4, ws...)
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

	targets := []wtype.Coordinates{{0.0, -10.0, 0.0}, {0.0, 10.0, 0.0}}
	assert.Equal(t, targets, plate.GetTargets("Magick"), "Well Targets")
	assert.Equal(t, []wtype.Coordinates{}, plate.GetTargets("Muggles"), "Well Targets")
}

func assertTipsEqual(t *testing.T, a, b *wtype.LHTip, message string) {
	if a == nil && b == nil {
		fmt.Printf("Comparing nil tips at %s\n", message)
		return
	} else if a == nil {
		t.Errorf("Expected nil tip %s", message)
	} else if b == nil {
		t.Errorf("Unexpected nil tip %s", message)
	}

	assert.Equal(t, a.ID, b.ID, "Tip(%s) ID", message)
	assert.Equal(t, a.Mnfr, b.Mnfr, "Tip(%s) Mnfr", message)
	assert.Equal(t, a.Dirty, b.Dirty, "Tip(%s) Dirty", message)
	assert.Equal(t, a.MaxVol, b.MaxVol, "Tip(%s) MaxVol", message)
	assert.Equal(t, a.MinVol, b.MinVol, "Tip(%s) MinVol", message)
	assert.Equal(t, a.Bounds, b.Bounds, "Tip(%s) Bounds", message)
	//assuming contents are the same - IDs are changed though
	//assert.Equal(t, a.Contents(), b.Contents(), "Tip(%s) Contents", message)
	assert.Equal(t, a.Shape, b.Shape, "Tip(%s) Shape", message)
}

func assertTipBoxesEqual(t *testing.T, a, b *wtype.LHTipbox) {

	assert.Equal(t, a.ID, b.ID, "Tipbox ID")
	assert.Equal(t, a.Boxname, b.Boxname, "Tipbox Boxname")
	assert.Equal(t, a.Type, b.Type, "Tipbox Type")
	assert.Equal(t, a.Mnfr, b.Mnfr, "Tipbox Mnfr")
	assert.Equal(t, a.Nrows, b.Nrows, "Tipbox Nrows")
	assert.Equal(t, a.Ncols, b.Ncols, "Tipbox Ncols")
	assert.Equal(t, a.Height, b.Height, "Tipbox Height")
	assert.Equal(t, a.NTips, b.NTips, "Tipbox NTips")
	assert.Equal(t, a.TipXOffset, b.TipXOffset, "Tipbox TipXOffset")
	assert.Equal(t, a.TipYOffset, b.TipYOffset, "Tipbox TipYOffset")
	assert.Equal(t, a.TipXStart, b.TipXStart, "Tipbox TipXStart")
	assert.Equal(t, a.TipYStart, b.TipYStart, "Tipbox TipYStart")
	assert.Equal(t, a.TipZStart, b.TipZStart, "Tipbox TipZStart")

	assertTipsEqual(t, a.Tiptype, b.Tiptype, "Tipbox.Tiptype")

	for i := 0; i < a.Nrows; i++ {
		for j := 0; j < b.Ncols; j++ {
			assertTipsEqual(t, a.Tips[j][i], b.Tips[j][i], fmt.Sprintf("Tipbox.Tips[%d][%d]", j, i))
		}
	}
}
