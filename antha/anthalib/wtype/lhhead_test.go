package wtype

import (
	"testing"
)

//makeAlignmentTestPlate make a plate setting only the important things
func makeTestPlate(wellsX, wellsY int, offsetX, offsetY float64) *LHPlate {
	plateSize := Coordinates3D{X: 127.76, Y: 85.48, Z: 15.0}
	wellSize := Coordinates3D{X: plateSize.X / float64(wellsX), Y: plateSize.Y / float64(wellsY), Z: 15.0}

	shape := NewShape(BoxShape, "mm", wellSize.X, wellSize.Y, wellSize.Z)
	well := NewLHWell("ul", 100.0, 10.0, shape, FlatWellBottom, wellSize.X, wellSize.Y, wellSize.Z, 0.0, "mm")
	return NewLHPlate("testplate", "", wellsX, wellsY, plateSize, well, offsetX, offsetY, 0.0, 0.0, 0.0)
}

func TestHeadDup(t *testing.T) {
	head := &LHHead{
		Name:         "headName",
		Manufacturer: "headMfg",
		ID:           "originalID",
		Adaptor: &LHAdaptor{
			ID:     "originalID",
			Params: &LHChannelParameter{},
		},
		Params: &LHChannelParameter{
			ID: "originalID",
		},
	}

	newID := head.Dup()
	oldID := head.DupKeepIDs()

	if head.ID != oldID.ID {
		t.Error("head.ID was changed by DupKeepIDs")
	}
	if head.Adaptor.ID != oldID.Adaptor.ID {
		t.Error("head.Adaptor.ID was changed by DupKeepIDs")
	}
	if head.Params.ID != oldID.Params.ID {
		t.Error("head.Params.ID was changed by DupKeepIDs")
	}

	if head.ID == newID.ID {
		t.Error("head.ID was changed by Dup")
	}
	if head.Adaptor.ID == newID.Adaptor.ID {
		t.Error("head.Adaptor.ID was changed by Dup")
	}
	if head.Params.ID == newID.Params.ID {
		t.Error("head.Params.ID was changed by Dup")
	}

}

type headCanReachTest struct {
	Name          string             //to identify the test
	Independent   bool               //is the head capable of independent multi channel
	Orientation   ChannelOrientation //what orientation is the channel
	Multi         int                //number of channels
	Plate         *LHPlate           //the plate to use for the test
	WellAddresses []string           //well addresses that we want to move to
	Expected      bool
}

func (self *headCanReachTest) Run(t *testing.T) {
	t.Run(self.Name, self.run)
}

func (self *headCanReachTest) run(t *testing.T) {

	head := &LHHead{
		Adaptor: &LHAdaptor{
			Params: &LHChannelParameter{
				Independent: self.Independent,
				Orientation: self.Orientation,
				Multi:       self.Multi,
			},
		},
	}

	wc := WCArrayFromStrings(self.WellAddresses)
	if g := head.CanReach(self.Plate, wc); g != self.Expected {
		t.Errorf("got %t, expected %t", g, self.Expected)
	}
}

func TestHeadCanReachVChannel96Plate(t *testing.T) {

	plate := makeTestPlate(8, 12, 9.0, 9.0)

	tests := []*headCanReachTest{
		{
			Name:          "non-independent 8-well in A1",
			Independent:   false,
			Orientation:   LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1"},
			Expected:      true,
		},
		{
			Name:          "non-independent 8-well in A1-H1",
			Independent:   false,
			Orientation:   LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
			Expected:      true,
		},
		{
			Name:          "independent skipping a well",
			Independent:   true,
			Orientation:   LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1", "B1", "D1", "E1", "F1", "G1", "H1"}, //double the gap between channels 1 and 2
			Expected:      true,
		},
		{
			Name:          "wrong rows",
			Independent:   false,
			Orientation:   LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1", "C1"},
			Expected:      false,
		},
		{
			Name:          "wrong columns",
			Independent:   false,
			Orientation:   LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1", "B2"},
			Expected:      false,
		},
		{
			Name:          "wrong order",
			Independent:   false,
			Orientation:   LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"B1", "A1"},
			Expected:      false,
		},
		{
			Name:          "independent rows",
			Independent:   true,
			Orientation:   LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1", "", "C1"},
			Expected:      true,
		},
		{
			Name:          "independent columns",
			Independent:   true,
			Orientation:   LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1", "B2"},
			Expected:      false,
		},
		{
			Name:          "independent wrong order",
			Independent:   true,
			Orientation:   LHVChannel,
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

func TestHeadCanReachHChannelPCRPlate(t *testing.T) {
	plate := makeTestPlate(8, 12, 9.0, 9.0)

	tests := []*headCanReachTest{
		{
			Name:          "non-independent 8-well in A1",
			Independent:   false,
			Orientation:   LHHChannel,
			Multi:         12,
			Plate:         plate,
			WellAddresses: []string{"A1"},
			Expected:      true,
		},
		{
			Name:          "non-independent 8-well in A1-H1",
			Independent:   false,
			Orientation:   LHHChannel,
			Multi:         12,
			Plate:         plate,
			WellAddresses: []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9", "A10", "A11", "A12"},
			Expected:      true,
		},
		{
			Name:          "wrong rows",
			Independent:   false,
			Orientation:   LHHChannel,
			Multi:         12,
			Plate:         plate,
			WellAddresses: []string{"A1", "B2"},
			Expected:      false,
		},
		{
			Name:          "wrong columns",
			Independent:   false,
			Orientation:   LHHChannel,
			Multi:         12,
			Plate:         plate,
			WellAddresses: []string{"A1", "B1"},
			Expected:      false,
		},
		{
			Name:          "wrong order",
			Independent:   false,
			Orientation:   LHHChannel,
			Multi:         12,
			Plate:         plate,
			WellAddresses: []string{"A2", "A1"},
			Expected:      false,
		},
		{
			Name:          "independent rows",
			Independent:   true,
			Orientation:   LHHChannel,
			Multi:         12,
			Plate:         plate,
			WellAddresses: []string{"A1", "", "A3"},
			Expected:      true,
		},
		{
			Name:          "independent columns",
			Independent:   true,
			Orientation:   LHHChannel,
			Multi:         12,
			Plate:         plate,
			WellAddresses: []string{"A1", "B2"},
			Expected:      false,
		},
		{
			Name:          "independent wrong order",
			Independent:   true,
			Orientation:   LHHChannel,
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

func TestHeadCanReach384Plate(t *testing.T) {

	plate := makeTestPlate(16, 24, 4.5, 4.5)

	tests := []*headCanReachTest{
		{
			Name:          "non-independent 8-well in A1",
			Independent:   false,
			Orientation:   LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1"},
			Expected:      true,
		},
		{
			Name:          "non-independent every other well",
			Independent:   false,
			Orientation:   LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1", "C1", "E1", "G1", "I1", "K1", "M1", "O1"},
			Expected:      true,
		},
		{
			Name:          "non-independent every other well offset",
			Independent:   false,
			Orientation:   LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"B1", "D1", "F1", "H1", "J1", "L1", "N1", "P1"},
			Expected:      true,
		},
		{
			Name:          "non-independent can't skip wells",
			Independent:   false,
			Orientation:   LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1", "C1", "E1", "I1", "K1", "M1", "O1"}, //missing G1
			Expected:      false,
		},
		{
			Name:          "non-independent can't do adjacent wells",
			Independent:   false,
			Orientation:   LHVChannel,
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

func TestHeadCanReachTrough(t *testing.T) {
	troughY := makeTestPlate(8, 1, 9.0, 0.0)
	troughX := makeTestPlate(1, 12, 0.0, 9.0)

	tests := []*headCanReachTest{
		{
			Name:          "non-independent in A1",
			Independent:   false,
			Orientation:   LHVChannel,
			Multi:         8,
			Plate:         troughY,
			WellAddresses: []string{"A1"},
			Expected:      true,
		},
		{
			Name:          "non-independent all channels in A1",
			Independent:   false,
			Orientation:   LHVChannel,
			Multi:         8,
			Plate:         troughY,
			WellAddresses: []string{"A1", "A1", "A1", "A1", "A1", "A1", "A1", "A1"},
			Expected:      true,
		},
		{
			Name:          "non-independent in A1",
			Independent:   false,
			Orientation:   LHHChannel,
			Multi:         8,
			Plate:         troughX,
			WellAddresses: []string{"A1"},
			Expected:      true,
		},
		{
			Name:          "non-independent all channels in A1",
			Independent:   false,
			Orientation:   LHHChannel,
			Multi:         8,
			Plate:         troughX,
			WellAddresses: []string{"A1", "A1", "A1", "A1", "A1", "A1", "A1", "A1"},
			Expected:      true,
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func TestHeadCanReachTwoRowTrough(t *testing.T) {
	trough := makeTestPlate(2, 12, 36.0, 9.0)

	tests := []*headCanReachTest{
		{
			Name:          "non-independent in A1",
			Independent:   false,
			Orientation:   LHHChannel,
			Multi:         8,
			Plate:         trough,
			WellAddresses: []string{"A1"},
			Expected:      true,
		},
		{
			Name:          "non-independent all channels in A1 and B1",
			Independent:   false,
			Orientation:   LHHChannel,
			Multi:         8,
			Plate:         trough,
			WellAddresses: []string{"A1", "A1", "A1", "A1", "B1", "B1", "B1", "B1"},
			Expected:      false, //we don't support this functionality currently
		},
	}

	for _, test := range tests {
		test.Run(t)
	}
}

func TestHeadCanReachWeirdPlate(t *testing.T) {
	plate := makeTestPlate(16, 24, 4, 4)

	tests := []*headCanReachTest{
		{
			Name:          "non-independent 8-well in A1",
			Independent:   false,
			Orientation:   LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1"},
			Expected:      true,
		},
		{
			Name:          "non-independent can't spread adaptors",
			Independent:   false,
			Orientation:   LHVChannel,
			Multi:         8,
			Plate:         plate,
			WellAddresses: []string{"A1", "B1"}, //the wells are 4 mm apart so you can't actually do this
			Expected:      false,
		},
		{
			Name:          "non-independent can't spread adaptors",
			Independent:   false,
			Orientation:   LHVChannel,
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
