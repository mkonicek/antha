package wtype

import (
	"encoding/json"
	"reflect"
	"testing"
)

func makeTipForTest() *LHTip {
	shp := NewShape("cylinder", "mm", 7.3, 7.3, 51.2)
	return NewLHTip("me", "mytype", 0.5, 1000.0, "ul", false, shp, 44.7)
}

func TestTipSerialization(t *testing.T) {
	cmp := makeComponent()

	before := makeTipForTest()
	before.SetContents(cmp)

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
