package liquidhandling

import (
	"context"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"
	"testing"
)

func TestAssertWFContiguousNonEmpty(t *testing.T) {
	names := []string{"Empty", "OneContiguous", "TwoContiguous", "TwoDiscontiguous", "TrailingSpaces"}
	things := [][]string{{"", ""}, {"", "", "", "X"}, {"", "", "", "A", "B"}, {"A", "", "B"}, {"A", "", ""}}
	wants := []bool{false, true, true, false, true}

	for i := 0; i < len(things); i++ {
		name := names[i]
		want := wants[i]
		thing := things[i]

		testFunc := func(t *testing.T) {
			got := assertWFContiguousNonEmpty(thing)
			if got != want {
				t.Errorf("assertWFContiguousNonEmpty(%v) : expected %v got %v ", thing, want, got)
			}
		}

		t.Run(name, testFunc)
	}
}

func TestPhysicalTipCheck(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())
	robot := makeGilson()
	head := robot.HeadsLoaded[0]

	names := []string{"No", "No", "No", "No", "No", "Yes", "No"}
	things := []string{"eppendorfrack425_1.5ml", "eppendorfrack425_1.5ml", "EGEL48", "eppendorfrack425_2ml", "eppendorfrack425_2ml", "falcon6wellAgar", "falcon6wellAgar"}
	wellsFrom := [][]string{{"A1", "B1"}, {"A1", "A1"}, {"A1", "B1"}, {"A1", "A1"}, {"A1", "B1"}, {"A1", "A1"}, {"A1", "B1"}}
	wants := []bool{false, false, false, false, false, true, false}

	for i := 0; i < len(things); i++ {
		want := wants[i]
		platetype := things[i]
		name := platetype + " " + names[i]
		wf := wellsFrom[i]

		p, err := inventory.NewPlate(ctx, platetype)

		if err != nil || p == nil {
			t.Errorf("Plate type %v n'existe pas", platetype)
		}

		testFunc := func(t *testing.T) {
			// func physicalTipCheck(robot *LHProperties, head *wtype.LHHead, plt *wtype.LHPlate, wellsFrom []string) bool

			got := physicalTipCheck(robot, head, p, wf)
			if got != want {
				t.Errorf("Expected %v got %v ", want, got)
			}
		}

		t.Run(name, testFunc)
	}

}
