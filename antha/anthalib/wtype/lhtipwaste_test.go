package wtype

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/antha-lang/antha/laboratory/effects/id"
)

func makeTipwasteForTest(idGen *id.IDGenerator) *LHTipwaste {
	shp := NewShape(BoxShape, "mm", 123.0, 80.0, 92.0)
	w := NewLHWell(idGen, "ul", 800000.0, 800000.0, shp, 0, 123.0, 80.0, 92.0, 0.0, "mm")
	lht := NewLHTipwaste(idGen, 6000, "TipwasteForTest", "ACME Corp.", Coordinates3D{X: 127.76, Y: 85.48, Z: 92.0}, w, 49.5, 31.5, 0.0)
	return lht
}

func TestTipwasteWellCoordsToCoords(t *testing.T) {
	idGen := id.NewIDGenerator(t.Name())

	tw := makeTipwasteForTest(idGen)

	pos, ok := tw.WellCoordsToCoords(idGen, MakeWellCoords("A1"), TopReference)
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
	idGen := id.NewIDGenerator(t.Name())

	tw := makeTipwasteForTest(idGen)

	eDelta := Coordinates3D{X: -0.2 * tw.AsWell.GetSize().X, Y: -0.3 * tw.AsWell.GetSize().Y}

	pos := Coordinates3D{
		X: tw.WellXStart + eDelta.X,
		Y: tw.WellYStart + eDelta.Y,
	}

	wc, delta := tw.CoordsToWellCoords(idGen, pos)

	if e, g := "A1", wc.FormatA1(); e != g {
		t.Errorf("Wrong well coordinates: expected %s, got %s", e, g)
	}

	if delta.X != eDelta.X || delta.Y != eDelta.Y {
		t.Errorf("Delta incorrect: expected (%f, %f), got (%f, %f)", eDelta.X, eDelta.Y, delta.X, delta.Y)
	}

}

func TestTipwasteSerialisation(t *testing.T) {
	idGen := id.NewIDGenerator(t.Name())
	partFull := makeTipwasteForTest(idGen)
	if err := partFull.SetOffset(Coordinates3D{X: 1.0, Y: 2.0, Z: 3.0}); err != nil {
		t.Fatal(err)
	}
	partFull.Contents = 300

	tipwastes := []*LHTipwaste{
		makeTipwasteForTest(idGen),
		partFull,
	}

	for _, before := range tipwastes {
		var after LHTipwaste
		if data, err := json.Marshal(before); err != nil {
			t.Fatal(err)
		} else if err := json.Unmarshal(data, &after); err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(before.AsWell, after.AsWell) {
			t.Errorf("serialisation changed the tipwaste:\nbefore: %#v\nafter : %#v", before.AsWell, after.AsWell)
		}
	}
}
