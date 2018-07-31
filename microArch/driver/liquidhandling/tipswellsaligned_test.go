package liquidhandling

import (
	"context"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"
	"testing"
)

//makeAlignmentTestPlate make a plate setting only the important things
func makeTestPlate(wellsX, wellsY int, offsetX, offsetY float64) *wtype.LHPlate {
	plateSize := wtype.Coordinates{X: 127.76, Y: 85.48, Z: 15.0}
	wellSize := wtype.Coordinates{X: plateSize.X / float64(wellsX), Y: plateSize.Y / float64(wellsY), Z: 15.0}

	shape := wtype.NewShape("box", "mm", wellSize.X, wellSize.Y, wellSize.Z)
	well := wtype.NewLHWell("ul", 100.0, 10.0, shape, wtype.FlatWellBottom, wellSize.X, wellSize.Y, wellSize.Z, 0.0, "mm")
	return wtype.NewLHPlate("testplate", "", wellsX, wellsY, plateSize, well, offsetX, offsetY, 0.0, 0.0, 0.0)
}

type canHeadReachTest struct {
	Name          string                   //to identify the test
	Independent   bool                     //is the head capable of independent multi channel
	Orientation   wtype.ChannelOrientation //what orientation is the channel
	Multi         int                      //number of channels
	Plate         *wtype.LHPlate           //the plate to use for the test
	WellAddresses []string                 //well addresses that we want to move to
	Expected      bool
}

func (self *canHeadReachTest) Run(t *testing.T) {
	t.Run(self.Name, self.run)
}

func (self *canHeadReachTest) run(t *testing.T) {

	head := &wtype.LHHead{
		Adaptor: &wtype.LHAdaptor{
			Params: &wtype.LHChannelParameter{
				Independent: self.Independent,
				Orientation: self.Orientation,
				Multi:       self.Multi,
			},
		},
	}

	wc := wtype.WCArrayFromStrings(self.WellAddresses)
	if g := CanHeadReach(head, self.Plate, wc); g != self.Expected {
		t.Errorf("got %t, expected %t", g, self.Expected)
	}
}

func TestCanHeadReachVChannel96Plate(t *testing.T) {

	plate := makeTestPlate(8, 12, 9.0, 9.0)

	tests := []*canHeadReachTest{
		{
			Name:          "non-independent 8-well in A1",
			Independent:   false,
			Orientation:   wtype.LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1"},
			Expected:      true,
		},
		{
			Name:          "non-independent 8-well in A1-H1",
			Independent:   false,
			Orientation:   wtype.LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
			Expected:      true,
		},
		{
			Name:          "wrong rows",
			Independent:   false,
			Orientation:   wtype.LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1", "C1"},
			Expected:      false,
		},
		{
			Name:          "wrong columns",
			Independent:   false,
			Orientation:   wtype.LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1", "B2"},
			Expected:      false,
		},
		{
			Name:          "wrong order",
			Independent:   false,
			Orientation:   wtype.LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"B1", "A1"},
			Expected:      false,
		},
		{
			Name:          "independent rows",
			Independent:   true,
			Orientation:   wtype.LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1", "", "C1"},
			Expected:      true,
		},
		{
			Name:          "independent columns",
			Independent:   true,
			Orientation:   wtype.LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1", "B2"},
			Expected:      false,
		},
		{
			Name:          "independent wrong order",
			Independent:   true,
			Orientation:   wtype.LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"B1", "A1"},
			Expected:      false,
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func TestCanHeadReachHChannelPCRPlate(t *testing.T) {
	plate := makeTestPlate(8, 12, 9.0, 9.0)

	tests := []*canHeadReachTest{
		{
			Name:          "non-independent 8-well in A1",
			Independent:   false,
			Orientation:   wtype.LHHChannel,
			Multi:         12,
			Plate:         plate,
			WellAddresses: []string{"A1"},
			Expected:      true,
		},
		{
			Name:          "non-independent 8-well in A1-H1",
			Independent:   false,
			Orientation:   wtype.LHHChannel,
			Multi:         12,
			Plate:         plate,
			WellAddresses: []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9", "A10", "A11", "A12"},
			Expected:      true,
		},
		{
			Name:          "wrong rows",
			Independent:   false,
			Orientation:   wtype.LHHChannel,
			Multi:         12,
			Plate:         plate,
			WellAddresses: []string{"A1", "B2"},
			Expected:      false,
		},
		{
			Name:          "wrong columns",
			Independent:   false,
			Orientation:   wtype.LHHChannel,
			Multi:         12,
			Plate:         plate,
			WellAddresses: []string{"A1", "B1"},
			Expected:      false,
		},
		{
			Name:          "wrong order",
			Independent:   false,
			Orientation:   wtype.LHHChannel,
			Multi:         12,
			Plate:         plate,
			WellAddresses: []string{"A2", "A1"},
			Expected:      false,
		},
		{
			Name:          "independent rows",
			Independent:   true,
			Orientation:   wtype.LHHChannel,
			Multi:         12,
			Plate:         plate,
			WellAddresses: []string{"A1", "", "A3"},
			Expected:      true,
		},
		{
			Name:          "independent columns",
			Independent:   true,
			Orientation:   wtype.LHHChannel,
			Multi:         12,
			Plate:         plate,
			WellAddresses: []string{"A1", "B2"},
			Expected:      false,
		},
		{
			Name:          "independent wrong order",
			Independent:   true,
			Orientation:   wtype.LHHChannel,
			Multi:         12,
			Plate:         plate,
			WellAddresses: []string{"A2", "A1"},
			Expected:      false,
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func TestCanHeadReach384Plate(t *testing.T) {

	plate := makeTestPlate(16, 24, 4.5, 4.5)

	tests := []*canHeadReachTest{
		{
			Name:          "non-independent 8-well in A1",
			Independent:   false,
			Orientation:   wtype.LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1"},
			Expected:      true,
		},
		{
			Name:          "non-independent every other well",
			Independent:   false,
			Orientation:   wtype.LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1", "C1", "E1", "G1", "I1", "K1", "M1", "O1"},
			Expected:      true,
		},
		{
			Name:          "non-independent every other well offset",
			Independent:   false,
			Orientation:   wtype.LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"B1", "D1", "F1", "H1", "J1", "L1", "N1", "P1"},
			Expected:      true,
		},
		{
			Name:          "non-independent can't skip wells",
			Independent:   false,
			Orientation:   wtype.LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1", "C1", "E1", "I1", "K1", "M1", "O1"}, //missing G1
			Expected:      false,
		},
		{
			Name:          "non-independent can't do adjacent wells",
			Independent:   false,
			Orientation:   wtype.LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1", "B1"},
			Expected:      false,
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func TestCanHeadReachWeirdPlate(t *testing.T) {
	plate := makeTestPlate(16, 24, 4, 4)

	tests := []*canHeadReachTest{
		{
			Name:          "non-independent 8-well in A1",
			Independent:   false,
			Orientation:   wtype.LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1"},
			Expected:      true,
		},
		{
			Name:          "non-independent can't spread adaptors",
			Independent:   false,
			Orientation:   wtype.LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1", "B1"}, //the wells are 4 mm apart so you can't actually do this
			Expected:      false,
		},
		{
			Name:          "non-independent can't spread adaptors",
			Independent:   false,
			Orientation:   wtype.LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1", "C1"}, //the wells are 8 mm apart so you can't actually do this
			Expected:      false,
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}

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
	robot := MakeGilsonForTest(defaultTipList())
	head := robot.GetLoadedHead(0)

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
			got := physicalTipCheck(head, p, wf)
			if got != want {
				t.Errorf("Expected %v got %v ", want, got)
			}
		}

		t.Run(name, testFunc)
	}

}
