package wtype

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

func makeTipForTest() *LHTip {
	shp := NewShape(CylinderShape, "mm", 7.3, 7.3, 51.2)
	return NewLHTip("me", "mytype", 0.5, 1000.0, "ul", false, shp, 44.7)
}

func TestTipSerialization(t *testing.T) {

	tips := []*LHTip{
		{
			ID:     "tipID",
			Type:   "tiptype",
			Mnfr:   "manufacturer",
			Dirty:  true,
			MaxVol: wunit.NewVolume(50, "ul"),
			MinVol: wunit.NewVolume(5, "ul"),
			Shape:  &Shape{},
			Bounds: BBox{
				Position: Coordinates3D{X: 12.0, Y: 24.0, Z: 48.0},
				Size:     Coordinates3D{X: 10.0, Y: 20.0, Z: 40.0},
			},
			EffectiveHeight: 42.0,
			Filtered:        true,
			contents:        makeComponent(),
		},
		makeTipForTest(),
	}

	for _, before := range tips {
		var after LHTip
		if bs, err := json.Marshal(before); err != nil {
			t.Fatal(err)
		} else if err := json.Unmarshal(bs, &after); err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(before, &after) {
			t.Errorf("serialization changed the tip:\nbefore: %+v\n after: %+v", before, &after)
		}
	}
}
