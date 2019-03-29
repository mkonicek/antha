package wtype

import (
	"encoding/json"
	"reflect"
	"testing"
)

func makeTipboxForTest() *LHTipbox {
	shp := NewShape(CylinderShape, "mm", 7.3, 7.3, 51.2)
	w := NewLHWell("ul", 250.0, 10.0, shp, FlatWellBottom, 7.3, 7.3, 51.2, 0.0, "mm")
	tiptype := makeTipForTest()
	tb := NewLHTipbox(8, 12, Coordinates3D{127.76, 85.48, 120.0}, "me", "mytype", tiptype, w, 9.0, 9.0, 0.5, 0.5, 0.0)
	return tb
}

func TestInflateMask(t *testing.T) {
	m := []bool{true}

	for i := 0; i < 8; i++ {
		inflated := inflateMask(m, i, 8)

		expected := make([]bool, 8)

		expected[i] = true

		if !reflect.DeepEqual(expected, inflated) {
			t.Errorf("Expected %v got %v", expected, inflated)
		}
	}
}

func TestMaskToWellCoords(t *testing.T) {
	ori := LHVChannel
	for i := 0; i < 8; i++ {
		m := make([]bool, 8)
		m[i] = true

		for j := 0; j < 12; j++ {
			expected := make([]string, 8)

			wc := WellCoords{X: j, Y: i}

			expected[i] = wc.FormatA1()

			got := maskToWellCoords(m, j, ori)

			if !reflect.DeepEqual(expected, got) {
				t.Errorf("Expected %v got %v", expected, got)
			}
		}
	}
}

// func NewLHTipbox(nrows, ncols int, height float64, manufacturer, boxtype string, tiptype *LHTip, well *LHWell, tipxoffset, tipyoffset, tipxstart, tipystart, tipzstart float64)
func TestGetTipsMasked(t *testing.T) {
	tb := makeTipboxForTest()

	mask := []bool{true}

	for i := 0; i < 96; i++ {
		wells, err := tb.GetTipsMasked(mask, LHVChannel, true)

		if err != nil {
			t.Errorf(err.Error())
		}

		if wells[0] == "" {
			t.Errorf("Ran out of tips too soon (%d)", i)
		}
	}
}

func TestGetTipsMasked2(t *testing.T) {
	tb := makeTipboxForTest()

	mask := make([]bool, 8)
	mask[2] = true

	for i := 0; i < 12; i++ {
		wells, err := tb.GetTipsMasked(mask, LHVChannel, false)
		if err != nil {
			t.Errorf(err.Error())
		}

		if wells[2] == "" {
			t.Errorf("Ran out of tips too soon (%d)", i)
		}
	}
}

func TestHasCleanTips(t *testing.T) {
	tb := makeTipboxForTest()

	m := make([]bool, 8)

	m[3] = true

	ori := LHVChannel

	for i := 0; i < 12; i++ {
		if !tb.hasCleanTips(i, m, ori) {
			t.Errorf("Offset %d mask %v has a clean tip but claims not to", i, m)
		}

		tb.Tips[i][3] = nil

		if tb.hasCleanTips(i, m, ori) {
			t.Errorf("Offset %d mask %v has no clean tip but cleams to", i, m)
		}
	}
}

func TestTrimToMask(t *testing.T) {
	wells := make([]string, 8)
	wells[1] = "B1"
	mask := []bool{true}

	trimmed := trimToMask(wells, mask)
	expected := []string{"B1"}

	if !reflect.DeepEqual(expected, trimmed) {
		t.Errorf("Expected %v, got %v", expected, trimmed)
	}

}

func TestTipboxWellCoordsToCoords(t *testing.T) {

	tb := makeTipboxForTest()

	pos, ok := tb.WellCoordsToCoords(MakeWellCoords("A1"), BottomReference)
	if !ok {
		t.Fatal("well A1 doesn't exist!")
	}

	xExpected := tb.TipXStart
	yExpected := tb.TipYStart

	if pos.X != xExpected || pos.Y != yExpected {
		t.Errorf("position was wrong: expected (%f, %f) got (%f, %f)", xExpected, yExpected, pos.X, pos.Y)
	}

}

func TestTipboxCoordsToWellCoords(t *testing.T) {

	tb := makeTipboxForTest()

	pos := Coordinates3D{
		X: tb.TipXStart + 0.75*tb.TipXOffset,
		Y: tb.TipYStart + 0.75*tb.TipYOffset,
	}

	wc, delta := tb.CoordsToWellCoords(pos)

	if e, g := "B2", wc.FormatA1(); e != g {
		t.Errorf("Wrong well coordinates: expected %s, got %s", e, g)
	}

	eDelta := -0.25 * tb.TipXOffset
	if delta.X != eDelta || delta.Y != eDelta {
		t.Errorf("Delta incorrect: expected (%f, %f), got (%f, %f)", eDelta, eDelta, delta.X, delta.Y)
	}

}

func TestTipboxGetWellBounds(t *testing.T) {

	tb := makeTipboxForTest()

	eStart := Coordinates3D{
		X: 0.5 - 0.5*7.3,
		Y: 0.5 - 0.5*7.3,
		Z: 0.0,
	}
	eSize := Coordinates3D{
		X: 9.0*11 + 7.3,
		Y: 9.0*7 + 7.3,
		Z: 51.2,
	}
	eBounds := NewBBox(eStart, eSize)
	bounds := tb.GetTipBounds()

	if e, g := eBounds.String(), bounds.String(); e != g {
		t.Errorf("GetWellBounds incorrect: expected %v, got %v", eBounds, bounds)
	}
}

func TestTipboxSerialization(t *testing.T) {

	removed := makeTipboxForTest()
	toRemove := MakeWellCoordsArray([]string{"A1", "B2", "H5"})
	for _, wc := range toRemove {
		removed.RemoveTip(wc)
	}

	if e, g := (removed.NCols()*removed.NRows() - len(toRemove)), removed.N_clean_tips(); e != g {
		t.Fatalf("LHTipbox.RemoveTip didn't work: expected %d tips remaining, got %d", e, g)
	}

	tipboxes := []*LHTipbox{
		makeTipboxForTest(),
	}

	for _, before := range tipboxes {

		var after LHTipbox
		if bs, err := json.Marshal(before); err != nil {
			t.Fatal(err)
		} else if err := json.Unmarshal(bs, &after); err != nil {
			t.Fatal(err)
		}

		for _, row := range after.Tips {
			for _, tip := range row {
				if tip.parent != &after {
					t.Fatal("parent not set correctly")
				}
			}
		}

		if !reflect.DeepEqual(before, &after) {
			t.Errorf("serialization changed the tipbox:\nbefore: %+v\n after: %+v", before, &after)
		}
	}
}
