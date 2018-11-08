package wtype

import (
	"encoding/json"
	"reflect"
	"testing"
)

func makeTipwasteForTest() *LHTipwaste {
	shp := NewShape("box", "mm", 123.0, 80.0, 92.0)
	w := NewLHWell("ul", 800000.0, 800000.0, shp, 0, 123.0, 80.0, 92.0, 0.0, "mm")
	lht := NewLHTipwaste(6000, "TipwasteForTest", "ACME Corp.", Coordinates{X: 127.76, Y: 85.48, Z: 92.0}, w, 49.5, 31.5, 0.0)
	return lht
}

func TestTipwasteWellCoordsToCoords(t *testing.T) {

	tw := makeTipwasteForTest()

	pos, ok := tw.WellCoordsToCoords(MakeWellCoords("A1"), TopReference)
	if !ok {
		t.Fatal("well A1 doesn't exist!")
	}

	xExpected := tw.WellXStart
	yExpected := tw.WellYStart

	if pos.X != xExpected || pos.Y != yExpected {
		t.Errorf("position was wrong: expected (%f, %f) got (%f, %f)", xExpected, yExpected, pos.X, pos.Y)
	}

}

func TestTipwasteCoordsToWellCoords(t *testing.T) {

	tw := makeTipwasteForTest()

	eDelta := Coordinates{X: -0.2 * tw.AsWell.GetSize().X, Y: -0.3 * tw.AsWell.GetSize().Y}

	pos := Coordinates{
		X: tw.WellXStart + eDelta.X,
		Y: tw.WellYStart + eDelta.Y,
	}

	wc, delta := tw.CoordsToWellCoords(pos)

	if e, g := "A1", wc.FormatA1(); e != g {
		t.Errorf("Wrong well coordinates: expected %s, got %s", e, g)
	}

	if delta.X != eDelta.X || delta.Y != eDelta.Y {
		t.Errorf("Delta incorrect: expected (%f, %f), got (%f, %f)", eDelta.X, eDelta.Y, delta.X, delta.Y)
	}

}

func TestTipwasteSerialisation(t *testing.T) {
	before := makeTipwasteForTest()

	var after LHTipwaste
	if data, err := json.Marshal(before); err != nil {
		t.Fatal(err)
	} else if err := json.Unmarshal(data, &after); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(before, &after) {
		t.Errorf("serialisation changed the tipwaste:\ne: %v\ng: %v", before, &after)
	}
}
