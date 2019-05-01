// /anthalib/simulator/liquidhandling/simulator_test.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package liquidhandling

import (
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

func TestUnknownLocations(t *testing.T) {
	assertPropsInvalid := func(t *testing.T, props *liquidhandling.LHProperties, message string) {
		if _, err := NewVirtualLiquidHandler(props, nil); err == nil {
			t.Errorf("missing error in %s", message)
		}
	}

	lhp := defaultLHProperties()
	lhp.Preferences.Tipboxes = append(lhp.Preferences.Tipboxes, "undefined_pref")
	assertPropsInvalid(t, lhp, "tipboxes")
}

func TestNewVirtualLiquidHandler_ValidProps(t *testing.T) {
	(&SimulatorTest{"Create Valid VLH", nil, nil, nil, nil, nil}).Run(t)
}

func TestVLH_AddPlateTo(t *testing.T) {
	// create a wide deckposition and enable input plates to go there
	propsWithWideSlot := defaultLHProperties()
	propsWithWideSlot.Positions["widePosition"] = wtype.NewLHPosition("widePosition", wtype.Coordinates3D{X: -300}, wtype.Coordinates2D{X: 300.0, Y: 85.48})
	propsWithWideSlot.Preferences.Inputs = append(propsWithWideSlot.Preferences.Inputs, "widePosition")

	SimulatorTests{
		{
			Name: "OK",
			Instructions: []TestRobotInstruction{
				&Initialize{},
				&AddPlateTo{"tipbox_1", defaultLHTipbox("tipbox1"), "tipbox1"},
				&AddPlateTo{"tipbox_2", defaultLHTipbox("tipbox2"), "tipbox2"},
				&AddPlateTo{"input_1", defaultLHPlate("input1"), "input1"},
				&AddPlateTo{"input_2", defaultLHPlate("input2"), "input2"},
				&AddPlateTo{"output_1", defaultLHPlate("output1"), "output1"},
				&AddPlateTo{"output_2", defaultLHPlate("output2"), "output2"},
				&AddPlateTo{"tipwaste", defaultLHTipwaste("tipwaste"), "tipwaste"},
			},
		},
		{
			Name: "non plate type",
			Instructions: []TestRobotInstruction{
				&Initialize{},
				&AddPlateTo{"tipbox_1", "my plate's gone stringy", "not_a_plate"},
			},
			ExpectedErrors: []string{"(err) AddPlateTo[1]: Couldn't add object of type string to tipbox_1"},
		},
		{
			Name: "location full",
			Instructions: []TestRobotInstruction{
				&Initialize{},
				&AddPlateTo{"tipbox_1", defaultLHTipbox("p0"), "p0"},
				&AddPlateTo{"tipbox_1", defaultLHTipbox("p1"), "p1"},
			},
			ExpectedErrors: []string{"(err) AddPlateTo[2]: Couldn't add tipbox \"p1\" to location \"tipbox_1\" which already contains tipbox \"p0\""},
		},
		{
			Name: "tipbox on tipwaste location",
			Instructions: []TestRobotInstruction{
				&Initialize{},
				&AddPlateTo{"tipwaste", defaultLHTipbox("tipbox"), "tipbox"},
			},
			ExpectedErrors: []string{"(err) AddPlateTo[1]: Slot \"tipwaste\" can't accept tipbox \"tipbox\", only tipwaste allowed"},
		},
		{
			Name: "tipwaste on tipbox location",
			Instructions: []TestRobotInstruction{
				&Initialize{},
				&AddPlateTo{"tipbox_1", defaultLHTipwaste("tipwaste"), "tipwaste"},
			},
			ExpectedErrors: []string{"(err) AddPlateTo[1]: Slot \"tipbox_1\" can't accept tipwaste \"tipwaste\", only tipbox allowed"},
		},
		{
			Name: "unknown location",
			Instructions: []TestRobotInstruction{
				&Initialize{},
				&AddPlateTo{"ruritania", defaultLHTipbox("aTipbox"), "aTipbox"},
			},
			ExpectedErrors: []string{"(err) AddPlateTo[1]: Cannot put tipbox \"aTipbox\" at unknown slot \"ruritania\""},
		},
		{
			Name: "too big",
			Instructions: []TestRobotInstruction{
				&Initialize{},
				&AddPlateTo{"output_1", wideLHPlate("plate1"), "plate1"},
			},
			ExpectedErrors: []string{
				"(err) AddPlateTo[1]: Footprint of plate \"plate1\"[300mm x 85.48mm] doesn't fit slot \"output_1\"[127.76mm x 85.48mm]",
			},
		},
		{
			Name:  "non-SBS format plate",
			Props: propsWithWideSlot,
			Instructions: []TestRobotInstruction{
				&Initialize{},
				&AddPlateTo{"widePosition", wideLHPlate("plate1"), "plate1"},
			},
		},
	}.Run(t)
}

func Test_SetPippetteSpeed(t *testing.T) {
	SimulatorTests{
		{
			Name: "OK",
			Instructions: []TestRobotInstruction{
				&Initialize{},
				&SetPipetteSpeed{0, -1, 5.},
			},
		},
		{
			Name: "too low",
			Instructions: []TestRobotInstruction{
				&Initialize{},
				&SetPipetteSpeed{0, -1, 0.001},
			},
			ExpectedErrors: []string{
				"(warn) SetPipetteSpeed[1]: Setting Head 0 channels 0-7 speed to 0.001 ml/min is outside allowable range [0.1 ml/min:10 ml/min]",
			},
		},
		{
			Name: "too high",
			Instructions: []TestRobotInstruction{
				&Initialize{},
				&SetPipetteSpeed{0, -1, 15.},
			},
			ExpectedErrors: []string{
				"(warn) SetPipetteSpeed[1]: Setting Head 0 channels 0-7 speed to 15 ml/min is outside allowable range [0.1 ml/min:10 ml/min]",
			},
		},
		{
			Name: "Independent",
			Instructions: []TestRobotInstruction{
				&Initialize{},
				&SetPipetteSpeed{0, 3, 5.},
			},
			ExpectedErrors: []string{
				"(warn) SetPipetteSpeed[1]: Head 0 is not independent, setting pipette speed for channel 3 sets all other channels as well",
			},
		},
	}.Run(t)
}

func TestSetDriveSpeed(t *testing.T) {
	props := defaultLHProperties()

	//set max and minimum drive speeds
	props.HeadAssemblies[0].VelocityLimits = &wtype.VelocityRange{
		Min: &wunit.Velocity3D{
			X: wunit.NewVelocity(0.5, "mm/s"),
			Y: wunit.NewVelocity(0.5, "mm/s"),
			Z: wunit.NewVelocity(0.1, "mm/s"),
		},
		Max: &wunit.Velocity3D{
			X: wunit.NewVelocity(50, "mm/s"),
			Y: wunit.NewVelocity(50, "mm/s"),
			Z: wunit.NewVelocity(10, "mm/s"),
		},
	}

	SimulatorTests{
		{
			Name: "OK",
			Instructions: []TestRobotInstruction{
				&Initialize{},
				&SetDriveSpeed{drive: "X", speed: 5.},
				&SetDriveSpeed{drive: "Y", speed: 5.},
				&SetDriveSpeed{drive: "Z", speed: 5.},
			},
		},
		{
			Name: "invalid drive",
			Instructions: []TestRobotInstruction{
				&Initialize{},
				&SetDriveSpeed{drive: "Q", speed: 5.},
			},
			ExpectedErrors: []string{
				"(err) SetDriveSpeed[1]: while setting head group 0 drive Q speed to 5 mm/s: unknown axis \"Q\"",
			},
		},
		{
			Name: "negative value",
			Instructions: []TestRobotInstruction{
				&Initialize{},
				&SetDriveSpeed{drive: "Z", speed: -5.},
			},
			ExpectedErrors: []string{
				"(err) SetDriveSpeed[1]: while setting head group 0 drive Z speed to -5 mm/s: speed must be positive",
			},
		},
		{
			Name: "OK - with speed limits",
			Instructions: []TestRobotInstruction{
				&Initialize{},
				&SetDriveSpeed{drive: "X", speed: 5.},
				&SetDriveSpeed{drive: "Y", speed: 5.},
				&SetDriveSpeed{drive: "Z", speed: 5.},
			},
		},
		{
			Name:  "Outside range",
			Props: props,
			Instructions: []TestRobotInstruction{
				&Initialize{},
				&SetDriveSpeed{drive: "X", speed: 0.2},
				&SetDriveSpeed{drive: "Y", speed: 0.2},
				&SetDriveSpeed{drive: "Z", speed: 0.01},
			},
			ExpectedErrors: []string{
				"(err) SetDriveSpeed[1]: while setting head group 0 drive X speed to 0.2 mm/s: 0.2 mm/s is outside allowable range [0.5 mm/s - 50 mm/s]",
				"(err) SetDriveSpeed[2]: while setting head group 0 drive Y speed to 0.2 mm/s: 0.2 mm/s is outside allowable range [0.5 mm/s - 50 mm/s]",
				"(err) SetDriveSpeed[3]: while setting head group 0 drive Z speed to 0.01 mm/s: 0.01 mm/s is outside allowable range [0.1 mm/s - 10 mm/s]",
			},
		},
	}.Run(t)
}

// ########################################################################################################################
// ########################################################## Move
// ########################################################################################################################
func testLayout() *SetupFn {
	var ret SetupFn = func(vlh *VirtualLiquidHandler) {
		vlh.Initialize()
		vlh.AddPlateTo("tipbox_1", defaultLHTipbox("tipbox1"), "tipbox1")
		vlh.AddPlateTo("tipbox_2", defaultLHTipbox("tipbox2"), "tipbox2")
		vlh.AddPlateTo("input_1", defaultLHPlate("plate1"), "plate1")
		vlh.AddPlateTo("input_2", defaultLHPlate("plate2"), "plate2")
		vlh.AddPlateTo("output_1", defaultLHPlate("plate3"), "plate3")
		vlh.AddPlateTo("waste", defaultLHPlate("wasteplate"), "wasteplate")
		vlh.AddPlateTo("tipwaste", defaultLHTipwaste("tipwaste"), "tipwaste")
	}
	return &ret
}

func testLayoutTransposed() *SetupFn {
	var ret SetupFn = func(vlh *VirtualLiquidHandler) {
		vlh.Initialize()
		vlh.AddPlateTo("tipbox_1", defaultLHTipbox("tipbox1"), "tipbox1")
		vlh.AddPlateTo("input_2", defaultLHTipbox("tipbox2"), "tipbox2")
		vlh.AddPlateTo("tipwaste", defaultLHPlate("plate1"), "plate1")
		vlh.AddPlateTo("tipbox_2", defaultLHPlate("plate2"), "plate2")
		vlh.AddPlateTo("output_1", defaultLHPlate("plate3"), "plate3")
		vlh.AddPlateTo("input_1", defaultLHTipwaste("tipwaste"), "tipwaste")
	}
	return &ret
}

func testTroughLayout() *SetupFn {
	var ret SetupFn = func(vlh *VirtualLiquidHandler) {
		vlh.Initialize()
		vlh.AddPlateTo("tipbox_1", defaultLHTipbox("tipbox1"), "tipbox1")
		vlh.AddPlateTo("tipbox_2", defaultLHTipbox("tipbox2"), "tipbox2")
		vlh.AddPlateTo("input_1", troughLHPlate("trough1"), "trough1")
		vlh.AddPlateTo("input_2", defaultLHPlate("plate2"), "plate2")
		vlh.AddPlateTo("output_1", defaultLHPlate("plate3"), "plate3")
		vlh.AddPlateTo("tipwaste", defaultLHTipwaste("tipwaste"), "tipwaste")
	}
	return &ret
}

func Test_Move(t *testing.T) {
	SimulatorTests{
		{
			Name: "OK_1",
			Setup: []*SetupFn{
				testLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2"},
					wellcoords:   []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					head:         0,
				},
			},
			Assertions: []*AssertionFn{
				positionAssertion(0, wtype.Coordinates3D{X: 200.0, Y: 0.0, Z: 62.2}),
			},
		},
		{
			Name: "OK_2",
			Setup: []*SetupFn{
				testLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste"},
					wellcoords:   []string{"A1", "A1", "A1", "A1", "A1", "A1", "A1", "A1"},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{-31.5, -22.5, -13.5, -4.5, 4.5, 13.5, 22.5, 31.5},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste"},
					head:         0,
				},
			},
			Assertions: []*AssertionFn{
				positionAssertion(0, wtype.Coordinates3D{X: 49.5, Y: 400., Z: 93.}),
			},
		},
		{
			Name: "OK_2.5",
			Setup: []*SetupFn{
				testLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"tipwaste", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{-31.5, 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"tipwaste", "", "", "", "", "", "", ""},
					head:         0,
				},
			},
			Assertions: []*AssertionFn{
				positionAssertion(0, wtype.Coordinates3D{X: 49.5, Y: 400., Z: 93}),
			},
		},
		{
			Name: "OK_3",
			Setup: []*SetupFn{
				testLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"", "", "", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					wellcoords:   []string{"", "", "", "A1", "B1", "C1", "D1", "E1"},
					reference:    []int{0, 0, 0, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"", "", "", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					head:         0,
				},
			},
			Assertions: []*AssertionFn{
				positionAssertion(0, wtype.Coordinates3D{X: 0.0, Y: -27.0, Z: 62.2}),
			},
		},
		{
			Name: "OK_4",
			Setup: []*SetupFn{
				testLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"H1", "", "", "", "", "", "", ""},
					reference:    []int{1, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
			},
			Assertions: []*AssertionFn{
				positionAssertion(0, wtype.Coordinates3D{X: 400.0, Y: 63.0, Z: 25.7}),
			},
		},
		{
			Name: "OK_5",
			Setup: []*SetupFn{
				testLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"", "", "", "", "input_1", "", "", ""},
					wellcoords:   []string{"", "", "", "", "H1", "", "", ""},
					reference:    []int{0, 0, 0, 0, 1, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					plate_type:   []string{"", "", "", "", "plate", "", "", ""},
					head:         0,
				},
			},
			Assertions: []*AssertionFn{
				positionAssertion(0, wtype.Coordinates3D{X: 400.0, Y: 27.0, Z: 25.7}),
			},
		},
		{
			Name: "OK_trough",
			Setup: []*SetupFn{
				testTroughLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"},
					wellcoords:   []string{"A1", "A1", "A1", "A1", "A1", "A1", "A1", "A1"},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"trough", "trough", "trough", "trough", "trough", "trough", "trough", "trough"},
					head:         0,
				},
			},
			Assertions: []*AssertionFn{
				positionAssertion(0, wtype.Coordinates3D{X: 400.0, Y: -31.5, Z: 46.8}),
			},
		},
		{
			Name: "OK allow cones in trough",
			Setup: []*SetupFn{
				testTroughLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"},
					wellcoords:   []string{"A1", "A1", "A1", "A1", "A1", "A1", "A1", "A1"},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{-1., -1., -1., -1., -1., -1., -1., -1.},
					plate_type:   []string{"trough", "trough", "trough", "trough", "trough", "trough", "trough", "trough"},
					head:         0,
				},
			},
			Assertions: []*AssertionFn{
				positionAssertion(0, wtype.Coordinates3D{X: 400.0, Y: -31.5, Z: 44.8}),
			},
		},
		{
			Name: "cones collide with plate",
			Setup: []*SetupFn{
				testTroughLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"},
					wellcoords:   []string{"A1", "A1", "A1", "A1", "A1", "A1", "A1", "A1"},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{4.0, 4.0, 4.0, 4.0, 4.0, 4.0, 4.0, 4.0},
					offsetZ:      []float64{-1., -1., -1., -1., -1., -1., -1., -1.},
					plate_type:   []string{"trough", "trough", "trough", "trough", "trough", "trough", "trough", "trough"},
					head:         0,
				},
			},
			ExpectedErrors: []string{
				"(err) Move[0]: head 0 channels 0-7 to 1 mm below well_top of A1,A1,A1,A1,A1,A1,A1,A1@trough1 at position input_1: collision detected: head 0 channel 7 and plate \"trough1\" of type trough at position input_1",
			},
			Assertions: []*AssertionFn{
				positionAssertion(0, wtype.Coordinates3D{X: 400.0, Y: -27.5, Z: 44.8}),
			},
		},
		{
			Name: "unknown location",
			Setup: []*SetupFn{
				testLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"tipbox7", "tipbox7", "tipbox7", "tipbox7", "tipbox7", "tipbox7", "tipbox7", "tipbox7"},
					wellcoords:   []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					plate_type:   []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					head:         0,
				},
				&Move{
					deckposition: []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					wellcoords:   []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					plate_type:   []string{"tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste"},
					head:         0,
				},
			},
			ExpectedErrors: []string{
				"(err) Move[0]: Unknown location \"tipbox7\"",
				"(warn) Move[1]: Object found at tipbox_1 was type \"tipbox\", named \"tipbox1\", not \"tipwaste\" as expected",
				"(warn) Move[1]: Object found at tipbox_1 was type \"tipbox\", named \"tipbox1\", not \"tipwaste\" as expected",
				"(warn) Move[1]: Object found at tipbox_1 was type \"tipbox\", named \"tipbox1\", not \"tipwaste\" as expected",
				"(warn) Move[1]: Object found at tipbox_1 was type \"tipbox\", named \"tipbox1\", not \"tipwaste\" as expected",
				"(warn) Move[1]: Object found at tipbox_1 was type \"tipbox\", named \"tipbox1\", not \"tipwaste\" as expected",
				"(warn) Move[1]: Object found at tipbox_1 was type \"tipbox\", named \"tipbox1\", not \"tipwaste\" as expected",
				"(warn) Move[1]: Object found at tipbox_1 was type \"tipbox\", named \"tipbox1\", not \"tipwaste\" as expected",
				"(warn) Move[1]: Object found at tipbox_1 was type \"tipbox\", named \"tipbox1\", not \"tipwaste\" as expected",
			},
		},
		{
			Name: "unknown head",
			Setup: []*SetupFn{
				testLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					wellcoords:   []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					plate_type:   []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					head:         1,
				},
				&Move{
					deckposition: []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					wellcoords:   []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					plate_type:   []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					head:         -1,
				},
			},
			ExpectedErrors: []string{
				"(err) Move[0]: head assembly 0: unknown head 1",
				"(err) Move[1]: head assembly 0: unknown head -1",
			},
		},
		{
			Name: "invalid wellcoords",
			Setup: []*SetupFn{
				testLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					wellcoords:   []string{"B1", "C1", "D1", "E1", "F1", "G1", "H1", "I1"},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					plate_type:   []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					head:         0,
				},
				&Move{
					deckposition: []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					wellcoords:   []string{"B1", "C1", "D1", "E1", "F1", "G1", "H1", "not_a_well"},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					plate_type:   []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					head:         0,
				},
			},
			ExpectedErrors: []string{
				"(err) Move[0]: Request for well I1 in object \"tipbox1\" at \"tipbox_1\" which is of size [8x12]",
				"(err) Move[1]: invalid argument wellcoords: couldn't parse \"not_a_well\"",
			},
		},
		{
			Name: "Invalid reference",
			Setup: []*SetupFn{
				testLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					wellcoords:   []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
					reference:    []int{-1, -1, -1, -1, -1, -1, -1, -1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					plate_type:   []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					head:         0,
				},
				&Move{
					deckposition: []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					wellcoords:   []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
					reference:    []int{3, 3, 3, 3, 3, 3, 3, 3},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					plate_type:   []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					head:         0,
				},
			},
			ExpectedErrors: []string{
				"(err) Move[0]: invalid argument reference: unknown value -1",
				"(err) Move[1]: invalid argument reference: unknown value 3",
			},
		},
		{
			Name: "Inconsistent references",
			Setup: []*SetupFn{
				testLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"},
					wellcoords:   []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
					reference:    []int{0, 0, 0, 0, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{-5., -5., -5., -5., -5., -5., -5., -5.},
					plate_type:   []string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"},
					head:         0,
				},
			},
			ExpectedErrors: []string{
				"(err) Move[0]: head 0 channels 0-7 to 5 mm below {well_bottom,well_top} of A1-H1@plate1 at position input_1: requires moving channels 4-7 relative to non-independent head",
			},
		},
		{
			Name: "offsets differ",
			Setup: []*SetupFn{
				testLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					wellcoords:   []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 3., 0., 0.},
					offsetY:      []float64{0., 0., 0., 1., 0., 0., 0., 0.},
					offsetZ:      []float64{0., 0., 0., 0., 1., 0., 0., 0.},
					plate_type:   []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					head:         0,
				},
			},
			ExpectedErrors: []string{
				"(err) Move[0]: head 0 channels 0-7 to 1 mm above well_top of A1-H1@tipbox1 at position tipbox_1: requires moving channels 3-5 relative to non-independent head",
			},
		},
		{
			Name: "layout mismatch",
			Setup: []*SetupFn{
				testLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					wellcoords:   []string{"A1", "B2", "C1", "D2", "E1", "F2", "G1", "H2"},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					plate_type:   []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					head:         0,
				},
			},
			ExpectedErrors: []string{
				"(err) Move[0]: head 0 channels 0-7 to well_top of A1,B2,C1,D2,E1,F2,G1,H2@tipbox1 at position tipbox_1: requires moving channels 1,3,5,7 relative to non-independent head",
			},
		},
	}.Run(t)
}

func TestCrashes(t *testing.T) {
	SimulatorTests{
		{
			Name:  "crash into tipbox",
			Props: multiheadLHProperties(),
			Setup: []*SetupFn{
				testLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2"},
					wellcoords:   []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{-1., -1., -1., -1., -1., -1., -1., -1.},
					plate_type:   []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					head:         0,
				},
			},
			ExpectedErrors: []string{
				"(err) Move[0]: head 0 channels 0-7 to 1 mm below well_top of A1-H1@tipbox2 at position tipbox_2: collision detected: head 0 channels 0-7 and head 1 channels 0-7 and tips A1-H1,A3-H3@tipbox2 at position tipbox_2",
			},
			Assertions: []*AssertionFn{
				positionAssertion(0, wtype.Coordinates3D{X: 128.0, Y: 0.0, Z: 60.2}),
			},
		},
		{
			Name:  "collides with tipbox in front",
			Props: multiheadLHProperties(),
			Setup: []*SetupFn{
				testLayoutTransposed(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					wellcoords:   []string{"E12", "F12", "G12", "H12"},
					reference:    []int{1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1.},
					plate_type:   []string{"tipbox", "tipbox", "tipbox", "tipbox"},
					head:         0,
				},
				&LoadTips{
					channels:  []int{0, 1, 2, 3},
					head:      0,
					multi:     4,
					platetype: []string{"tipbox", "tipbox", "tipbox", "tipbox"},
					position:  []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					well:      []string{"E12", "F12", "G12", "H12"},
				},
			},
			ExpectedErrors: []string{
				"(err) LoadTips[1]: from E12-H12@tipbox1 at position \"tipbox_1\" to head 0 channels 0-3: collision detected: head 0 channels 5-7 and tips A12-C12@tipbox2 at position input_2",
			},
			Assertions: []*AssertionFn{},
		},
		{
			Name: "trying to move channel cones into a well",
			Setup: []*SetupFn{
				testLayout(),
				//preloadFilledTips(0, "tipbox_1", []int{0}, "water", 100.),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
			},
			ExpectedErrors: []string{
				"(err) Move[0]: head 0 channel 0 to 1 mm above well_bottom of A1@plate1 at position input_1: collision detected: head 0 channels 0-7 and plate \"plate1\" of type plate at position input_1",
			},
		},
	}.Run(t)
}

func Test_Multihead(t *testing.T) {
	SimulatorTests{
		{
			Name:  "constrained heads",
			Props: multiheadLHProperties(),
			Setup: []*SetupFn{
				testLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2"},
					wellcoords:   []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					head:         0,
				},
			},
			Assertions: []*AssertionFn{
				positionAssertion(0, wtype.Coordinates3D{X: 128.0, Y: 0.0, Z: 62.2}),
				positionAssertion(1, wtype.Coordinates3D{X: 146.0, Y: 0.0, Z: 62.2}),
			},
		},
		{
			Name:  "can't move while a tip is loaded on another head in the same assembly",
			Props: multiheadLHProperties(),
			Setup: []*SetupFn{
				testLayout(),
				preloadAdaptorTips(1, "tipbox_1", []int{0}),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					wellcoords:   []string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					head:         0,
				},
			},
			ExpectedErrors: []string{
				"(err) Move[0]: head 0 channels 0-7 to 1 mm above well_top of A12-H12@tipbox1 at position tipbox_1: cannot move head 0 while tip loaded on head 1 channel 0",
			},
			Assertions: []*AssertionFn{},
		},
	}.Run(t)
}

func TestMotionLimits(t *testing.T) {
	SimulatorTests{
		{
			Name:  "outside limits left",
			Props: multiheadLHProperties(),
			Setup: []*SetupFn{
				testLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					wellcoords:   []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					head:         1,
				},
			},
			ExpectedErrors: []string{
				"(err) Move[0]: head 1 channels 0-7 to 1 mm above well_top of A1-H1@tipbox1 at position tipbox_1: head cannot reach position: position is 9mm too far left, please try rearranging the deck",
			},
			Assertions: []*AssertionFn{
				positionAssertion(0, wtype.Coordinates3D{X: -18.0, Y: 0.0, Z: 62.2}),
				positionAssertion(1, wtype.Coordinates3D{X: 0.0, Y: 0.0, Z: 62.2}),
			},
		},
		{
			Name:  "outside limits right",
			Props: multiheadLHProperties(),
			Setup: []*SetupFn{
				testLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"},
					wellcoords:   []string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{50., 50., 50., 50., 50., 50., 50., 50.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"},
					head:         0,
				},
			},
			ExpectedErrors: []string{
				"(err) Move[0]: head 0 channels 0-7 to 1 mm above well_top of A12-H12@plate1 at position input_1: head cannot reach position: position is 30mm too far right, please try rearranging the deck",
			},
			Assertions: []*AssertionFn{
				positionAssertion(0, wtype.Coordinates3D{X: 405.0, Y: 0.0, Z: 26.7}),
				positionAssertion(1, wtype.Coordinates3D{X: 423.0, Y: 0.0, Z: 26.7}),
			},
		},
		{
			Name:  "outside limits forward",
			Props: multiheadLHProperties(),
			Setup: []*SetupFn{
				testLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"waste"},
					wellcoords:   []string{"H12"},
					reference:    []int{1},
					offsetX:      []float64{0.},
					offsetY:      []float64{30.},
					offsetZ:      []float64{1.},
					plate_type:   []string{"plate"},
					head:         0,
				},
			},
			ExpectedErrors: []string{
				"(err) Move[0]: head 0 channel 0 to 1 mm above well_top of H12@wasteplate at position waste: head cannot reach position: position is 7mm too far forwards, please try rearranging the deck",
			},
			Assertions: []*AssertionFn{
				positionAssertion(0, wtype.Coordinates3D{X: 355.0, Y: 265.0, Z: 26.7}),
				positionAssertion(1, wtype.Coordinates3D{X: 373.0, Y: 265.0, Z: 26.7}),
			},
		},
		{
			Name:  "outside limits backwards",
			Props: multiheadLHProperties(),
			Setup: []*SetupFn{
				testLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"", "", "", "", "", "", "", "input_1"},
					wellcoords:   []string{"", "", "", "", "", "", "", "A12"},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"", "", "", "", "", "", "", "plate"},
					head:         0,
				},
			},
			ExpectedErrors: []string{
				"(err) Move[0]: head 0 channel 7 to 1 mm above well_top of A12@plate1 at position input_1: head cannot reach position: position is 63mm too far backwards, please try rearranging the deck",
			},
			Assertions: []*AssertionFn{
				positionAssertion(0, wtype.Coordinates3D{X: 355.0, Y: -63.0, Z: 26.7}),
				positionAssertion(1, wtype.Coordinates3D{X: 373.0, Y: -63.0, Z: 26.7}),
			},
		},
		{
			Name:  "outside limits too high",
			Props: multiheadLHProperties(),
			Setup: []*SetupFn{
				testLayout(),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1"},
					wellcoords:   []string{"A1"},
					reference:    []int{1},
					offsetX:      []float64{0.},
					offsetY:      []float64{0.},
					offsetZ:      []float64{600.},
					plate_type:   []string{"plate"},
					head:         0,
				},
			},
			ExpectedErrors: []string{
				"(err) Move[0]: head 0 channel 0 to 600 mm above well_top of A1@plate1 at position input_1: head cannot reach position: position is 25.7mm too high, please try lowering the object on the deck",
			},
			Assertions: []*AssertionFn{
				positionAssertion(0, wtype.Coordinates3D{X: 256.0, Y: 0.0, Z: 625.7}),
				positionAssertion(1, wtype.Coordinates3D{X: 274.0, Y: 0.0, Z: 625.7}),
			},
		},
		{
			Name:  "outside limits too low",
			Props: multiheadConstrainedLHProperties(),
			Setup: []*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1"},
					wellcoords:   []string{"A4"},
					reference:    []int{0},
					offsetX:      []float64{0.},
					offsetY:      []float64{0.},
					offsetZ:      []float64{0.5},
					plate_type:   []string{"plate"},
					head:         0,
				},
			},
			ExpectedErrors: []string{
				"(err) Move[0]: head 0 channel 0 to 0.5 mm above well_bottom of A4@plate1 at position input_1: head cannot reach position: position is 8.1mm too low, please try adding a riser to the object on the deck",
			},
			Assertions: []*AssertionFn{
				positionAssertion(0, wtype.Coordinates3D{X: 283.0, Y: 0.0, Z: 51.9}),
				positionAssertion(1, wtype.Coordinates3D{X: 301.0, Y: 0.0, Z: 51.9}),
			},
		},
		{
			Name:  "outside limits too low and far back",
			Props: multiheadConstrainedLHProperties(),
			Setup: []*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{7}),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"", "", "", "", "", "", "", "input_1"},
					wellcoords:   []string{"", "", "", "", "", "", "", "A4"},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0, 0, 0, 0, 0, 0, 0, 0.},
					offsetY:      []float64{0, 0, 0, 0, 0, 0, 0, 0.},
					offsetZ:      []float64{0, 0, 0, 0, 0, 0, 0, 0.5},
					plate_type:   []string{"", "", "", "", "", "", "", "plate"},
					head:         0,
				},
			},
			ExpectedErrors: []string{
				"(err) Move[0]: head 0 channel 7 to 0.5 mm above well_bottom of A4@plate1 at position input_1: head cannot reach position: position is 63mm too far backwards and 8.1mm too low, please try rearranging the deck and adding a riser to the object on the deck",
			},
			Assertions: []*AssertionFn{
				positionAssertion(0, wtype.Coordinates3D{X: 283.0, Y: -63.0, Z: 51.9}),
				positionAssertion(1, wtype.Coordinates3D{X: 301.0, Y: -63.0, Z: 51.9}),
			},
		},
	}.Run(t)
}

// ########################################################################################################################
// ########################################################## Tip Loading/Unloading
// ########################################################################################################################

func TestLoadTipsNoOverride(t *testing.T) {

	mtp := moveToParams{
		Multi:        8,
		Head:         0,
		Reference:    1,
		Deckposition: "tipbox_1",
		Platetype:    "tipbox",
		Offset:       []float64{0., 0., 5.},
		Cols:         12,
		Rows:         8,
	}
	misaligned_mtp := moveToParams{
		Multi:        8,
		Head:         0,
		Reference:    1,
		Deckposition: "tipbox_1",
		Platetype:    "tipbox",
		Offset:       []float64{0., 2., 5.},
		Cols:         12,
		Rows:         8,
	}

	SimulatorTests{
		{
			Name: "OK - single tip",
			Setup: []*SetupFn{
				testLayout(),
				moveTo(7, 11, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0},
					head:      0,
					multi:     1,
					platetype: []string{"tipbox", "", "", "", "", "", "", ""},
					position:  []string{"tipbox_1", "", "", "", "", "", "", ""},
					well:      []string{"H12", "", "", "", "", "", "", ""},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{"H12"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name: "OK - single tip (alt)",
			Setup: []*SetupFn{
				testLayout(),
				moveTo(-7, 0, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{7},
					head:      0,
					multi:     1,
					platetype: []string{"", "", "", "", "", "", "", "tipbox"},
					position:  []string{"", "", "", "", "", "", "", "tipbox_1"},
					well:      []string{"", "", "", "", "", "", "", "A1"},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{"A1"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{7, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name: "OK - single tip above space",
			Setup: []*SetupFn{
				testLayout(),
				removeTipboxTips("tipbox_1", []string{"H12"}),
				moveTo(6, 11, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0},
					head:      0,
					multi:     1,
					platetype: []string{"tipbox", "", "", "", "", "", "", ""},
					position:  []string{"tipbox_1", "", "", "", "", "", "", ""},
					well:      []string{"G12", "", "", "", "", "", "", ""},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{"H12", "G12"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name: "OK - single tip below space (alt)",
			Setup: []*SetupFn{
				testLayout(),
				removeTipboxTips("tipbox_1", []string{"A1"}),
				moveTo(-6, 0, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{7},
					head:      0,
					multi:     1,
					platetype: []string{"", "", "", "", "", "", "", "tipbox"},
					position:  []string{"", "", "", "", "", "", "", "tipbox_1"},
					well:      []string{"", "", "", "", "", "", "", "B1"},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{"A1", "B1"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{7, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name: "OK - 3 tips at once",
			Setup: []*SetupFn{
				testLayout(),
				moveTo(5, 11, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0, 1, 2},
					head:      0,
					multi:     3,
					platetype: []string{"tipbox", "tipbox", "tipbox", "", "", "", "", ""},
					position:  []string{"tipbox_1", "tipbox_1", "tipbox_1", "", "", "", "", ""},
					well:      []string{"F12", "G12", "H12", "", "", "", "", ""},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{"F12", "G12", "H12"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{
					{0, "", 0},
					{1, "", 0},
					{2, "", 0},
				}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name: "OK - 3 tips at once (alt)",
			Setup: []*SetupFn{
				testLayout(),
				moveTo(-5, 0, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{5, 6, 7},
					head:      0,
					multi:     3,
					platetype: []string{"", "", "", "", "", "tipbox", "tipbox", "tipbox"},
					position:  []string{"", "", "", "", "", "tipbox_1", "tipbox_1", "tipbox_1"},
					well:      []string{"", "", "", "", "", "A1", "B1", "C1"},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{"A1", "B1", "C1"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{
					{5, "", 0},
					{6, "", 0},
					{7, "", 0},
				}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name:  "OK - 3 tips (independent)",
			Props: IndependentLHProperties(),
			Setup: []*SetupFn{
				testLayout(),
				moveTo(0, 0, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0, 4, 7},
					head:      0,
					multi:     3,
					platetype: []string{"tipbox", "", "", "", "tipbox", "", "", "tipbox"},
					position:  []string{"tipbox_1", "", "", "", "tipbox_1", "", "", "tipbox_1"},
					well:      []string{"A1", "", "", "", "E1", "", "", "H1"},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{"A1", "E1", "H1"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{
					{0, "", 0},
					{4, "", 0},
					{7, "", 0},
				}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name: "OK - 8 tips at once",
			Setup: []*SetupFn{
				testLayout(),
				moveTo(0, 11, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0, 1, 2, 3, 4, 5, 6, 7},
					head:      0,
					multi:     8,
					platetype: []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					position:  []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					well:      []string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{
					{0, "", 0},
					{1, "", 0},
					{2, "", 0},
					{3, "", 0},
					{4, "", 0},
					{5, "", 0},
					{6, "", 0},
					{7, "", 0},
				}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name: "OK - 2 groups of 4",
			Setup: []*SetupFn{
				testLayout(),
				moveTo(4, 0, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0, 1, 2, 3},
					head:      0,
					multi:     4,
					platetype: []string{"tipbox", "tipbox", "tipbox", "tipbox", "", "", "", ""},
					position:  []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "", "", "", ""},
					well:      []string{"E1", "F1", "G1", "H1", "", "", "", ""},
				},
				&Move{
					deckposition: []string{"", "", "", "", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					wellcoords:   []string{"", "", "", "", "A1", "B1", "C1", "D1"},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"", "", "", "", "tipbox", "tipbox", "tipbox", "tipbox"},
					head:         0,
				},
				&LoadTips{
					channels:  []int{4, 5, 6, 7},
					head:      0,
					multi:     4,
					platetype: []string{"", "", "", "", "tipbox", "tipbox", "tipbox", "tipbox"},
					position:  []string{"", "", "", "", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					well:      []string{"", "", "", "", "A1", "B1", "C1", "D1"},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{
					{0, "", 0},
					{1, "", 0},
					{2, "", 0},
					{3, "", 0},
					{4, "", 0},
					{5, "", 0},
					{6, "", 0},
					{7, "", 0},
				}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name: "unknown channel 8",
			Setup: []*SetupFn{
				testLayout(),
				moveTo(0, 0, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{8},
					head:      0,
					multi:     1,
					platetype: []string{"", "", "", "", "", "", "", "tipbox"},
					position:  []string{"", "", "", "", "", "", "", "tipbox_1"},
					well:      []string{"", "", "", "", "", "", "", "H12"},
				},
			},
			ExpectedErrors: []string{
				"(err) LoadTips[0]: Unknown channel \"8\"",
			},
		},
		{
			Name: "unknown channel -1",
			Setup: []*SetupFn{
				testLayout(),
				moveTo(0, 0, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{-1},
					head:      0,
					multi:     1,
					platetype: []string{"", "", "", "", "", "", "", "tipbox"},
					position:  []string{"", "", "", "", "", "", "", "tipbox_1"},
					well:      []string{"", "", "", "", "", "", "", "H12"},
				},
			},
			ExpectedErrors: []string{
				"(err) LoadTips[0]: Unknown channel \"-1\"",
			},
		},
		{
			Name: "duplicate channels",
			Setup: []*SetupFn{
				testLayout(),
				moveTo(0, 11, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0, 1, 2, 3, 4, 5, 6, 3},
					head:      0,
					multi:     8,
					platetype: []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					position:  []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					well:      []string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"},
				},
			},
			ExpectedErrors: []string{
				"(err) LoadTips[0]: Channel3 appears more than once",
			},
		},
		{
			Name: "unknown head",
			Setup: []*SetupFn{
				testLayout(),
				moveTo(7, 11, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0},
					head:      1,
					multi:     1,
					platetype: []string{"tipbox", "", "", "", "", "", "", ""},
					position:  []string{"tipbox_1", "", "", "", "", "", "", ""},
					well:      []string{"H12", "", "", "", "", "", "", ""},
				},
			},
			ExpectedErrors: []string{
				"(err) LoadTips[0]: head assembly 0: unknown head 1",
			},
		},
		{
			Name: "unknown head -1",
			Setup: []*SetupFn{
				testLayout(),
				moveTo(7, 11, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0},
					head:      -1,
					multi:     1,
					platetype: []string{"tipbox", "", "", "", "", "", "", ""},
					position:  []string{"tipbox_1", "", "", "", "", "", "", ""},
					well:      []string{"H12", "", "", "", "", "", "", ""},
				},
			},
			ExpectedErrors: []string{
				"(err) LoadTips[0]: head assembly 0: unknown head -1",
			},
		},
		{
			Name: "OK - argument expansion",
			Setup: []*SetupFn{
				testLayout(),
				moveTo(7, 11, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0},
					head:      0,
					multi:     1,
					platetype: []string{"tipbox"},
					position:  []string{"tipbox_1"},
					well:      []string{"H12"},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{"H12"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name: "mismatching multi",
			Setup: []*SetupFn{
				testLayout(),
				moveTo(7, 11, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0},
					head:      0,
					multi:     4,
					platetype: []string{"tipbox", "", "", ""},
					position:  []string{"tipbox_1", "", "", ""},
					well:      []string{"H12", "", "", ""},
				},
			},
			ExpectedErrors: []string{
				"(err) LoadTips[0]: from H12@tipbox1 at position \"tipbox_1\" to head 0 channel 0: multi should equal 1, not 4",
			},
		},
		{
			Name: "tip missing",
			Setup: []*SetupFn{
				testLayout(),
				removeTipboxTips("tipbox_1", []string{"H12"}),
				moveTo(7, 11, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0},
					head:      0,
					multi:     1,
					platetype: []string{"tipbox", "", "", "", "", "", "", ""},
					position:  []string{"tipbox_1", "", "", "", "", "", "", ""},
					well:      []string{"H12", "", "", "", "", "", "", ""},
				},
			},
			ExpectedErrors: []string{
				"(err) LoadTips[0]: from H12@tipbox1 at position \"tipbox_1\" to head 0 channel 0: no tip at H12",
			},
		},
		{
			Name: "8 tips missing",
			Setup: []*SetupFn{
				testLayout(),
				removeTipboxTips("tipbox_1", []string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"}),
				moveTo(0, 0, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0, 1, 2, 3, 4, 5, 6, 7},
					head:      0,
					multi:     8,
					platetype: []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					position:  []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					well:      []string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"},
				},
			},
			ExpectedErrors: []string{
				"(err) LoadTips[0]: from A12-H12@tipbox1 at position \"tipbox_1\" to head 0 channels 0-7: no tips at A12-H12",
			},
		},
		{
			Name: "tip already loaded",
			Setup: []*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
				moveTo(7, 11, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0},
					head:      0,
					multi:     1,
					platetype: []string{"tipbox", "", "", "", "", "", "", ""},
					position:  []string{"tipbox_1", "", "", "", "", "", "", ""},
					well:      []string{"H12", "", "", "", "", "", "", ""},
				},
			},
			ExpectedErrors: []string{
				"(err) LoadTips[0]: from H12@tipbox1 at position \"tipbox_1\" to head 0 channel 0: tip already loaded on head 0 channel 0",
			},
		},
		{
			Name: "tips already loaded",
			Setup: []*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}),
				moveTo(0, 11, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0, 1, 2, 3, 4, 5, 6, 7},
					head:      0,
					multi:     8,
					platetype: []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					position:  []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					well:      []string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"},
				},
			},
			ExpectedErrors: []string{
				"(err) LoadTips[0]: from A12-H12@tipbox1 at position \"tipbox_1\" to head 0 channels 0-7: tips already loaded on head 0 channels 0-7",
			},
		},
		{
			Name: "extra tip in the way",
			Setup: []*SetupFn{
				testLayout(),
				moveTo(6, 11, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0},
					head:      0,
					multi:     1,
					platetype: []string{"tipbox", "", "", "", "", "", "", ""},
					position:  []string{"tipbox_1", "", "", "", "", "", "", ""},
					well:      []string{"G12", "", "", "", "", "", "", ""},
				},
			},
			ExpectedErrors: []string{
				"(err) LoadTips[0]: from G12@tipbox1 at position \"tipbox_1\" to head 0 channel 0: collision detected: head 0 channel 1 and tip H12@tipbox1 at position tipbox_1",
			},
		},
		{
			Name: "not aligned to move",
			Setup: []*SetupFn{
				testLayout(),
				moveTo(5, 11, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0, 1, 2},
					head:      0,
					multi:     3,
					platetype: []string{"tipbox", "tipbox", "tipbox", "", "", "", "", ""},
					position:  []string{"tipbox_1", "tipbox_1", "tipbox_1", "", "", "", "", ""},
					well:      []string{"E12", "G12", "H12", "", "", "", "", ""},
				},
			},
			ExpectedErrors: []string{
				"(err) LoadTips[0]: from E12,G12-H12@tipbox1 at position \"tipbox_1\" to head 0 channels 0-2: channel 0 is misaligned with tip at E12 by 9mm",
			},
		},
		{
			Name: "multiple not aligned to move",
			Setup: []*SetupFn{
				testLayout(),
				moveTo(5, 11, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0, 1, 2},
					head:      0,
					multi:     3,
					platetype: []string{"tipbox", "tipbox", "tipbox", "", "", "", "", ""},
					position:  []string{"tipbox_1", "tipbox_1", "tipbox_1", "", "", "", "", ""},
					well:      []string{"G12", "F12", "H12", "", "", "", "", ""},
				},
			},
			ExpectedErrors: []string{
				"(err) LoadTips[0]: from G12,F12,H12@tipbox1 at position \"tipbox_1\" to head 0 channels 0-2: channels 0-1 are misaligned with tips at G12,F12 by 9,9mm respectively",
			},
		},
		{
			Name: "misalignment single",
			Setup: []*SetupFn{
				testLayout(),
				moveTo(7, 11, misaligned_mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0},
					head:      0,
					multi:     1,
					platetype: []string{"tipbox", "", "", "", "", "", "", ""},
					position:  []string{"tipbox_1", "", "", "", "", "", "", ""},
					well:      []string{"H12", "", "", "", "", "", "", ""},
				},
			},
			ExpectedErrors: []string{
				"(err) LoadTips[0]: from H12@tipbox1 at position \"tipbox_1\" to head 0 channel 0: channel 0 is misaligned with tip at H12 by 2mm",
			},
		},
		{
			Name: "misalignment multi",
			Setup: []*SetupFn{
				testLayout(),
				moveTo(5, 11, misaligned_mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0, 1, 2},
					head:      0,
					multi:     3,
					platetype: []string{"tipbox", "tipbox", "tipbox", "", "", "", "", ""},
					position:  []string{"tipbox_1", "tipbox_1", "tipbox_1", "", "", "", "", ""},
					well:      []string{"F12", "G12", "H12", "", "", "", "", ""},
				},
			},
			ExpectedErrors: []string{
				"(err) LoadTips[0]: from F12-H12@tipbox1 at position \"tipbox_1\" to head 0 channels 0-2: channels 0-2 are misaligned with tips at F12-H12 by 2,2,2mm respectively",
			},
		},
	}.Run(t)
}

func TestLoadTipsOverride(t *testing.T) {

	mtp := moveToParams{
		Multi:        8,
		Head:         0,
		Reference:    1,
		Deckposition: "tipbox_1",
		Platetype:    "tipbox",
		Offset:       []float64{0., 0., 5.},
		Cols:         12,
		Rows:         8,
	}

	propsLTR := defaultLHProperties()
	propsLTR.Heads[0].TipLoading = wtype.TipLoadingBehaviour{
		OverrideLoadTipsCommand:    true,
		AutoRefillTipboxes:         true,
		LoadingOrder:               wtype.ColumnWise,
		VerticalLoadingDirection:   wtype.BottomToTop,
		HorizontalLoadingDirection: wtype.LeftToRight,
		ChunkingBehaviour:          wtype.ReverseSequentialTipLoading,
	}

	propsRTL := defaultLHProperties()
	propsRTL.Heads[0].TipLoading = wtype.TipLoadingBehaviour{
		OverrideLoadTipsCommand:    true,
		AutoRefillTipboxes:         true,
		LoadingOrder:               wtype.ColumnWise,
		VerticalLoadingDirection:   wtype.BottomToTop,
		HorizontalLoadingDirection: wtype.RightToLeft,
		ChunkingBehaviour:          wtype.ReverseSequentialTipLoading,
	}

	SimulatorTests{
		{
			Name:  "OK - single tip LTR override (A1 -> H1)",
			Props: propsLTR,
			Setup: []*SetupFn{
				testLayout(),
				moveTo(0, 0, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0},
					head:      0,
					multi:     1,
					platetype: []string{"tipbox", "", "", "", "", "", "", ""},
					position:  []string{"tipbox_1", "", "", "", "", "", "", ""},
					well:      []string{"A1", "", "", "", "", "", "", ""},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{"H1"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name:  "OK - single tip RTL override (A12 -> H12)",
			Props: propsRTL,
			Setup: []*SetupFn{
				testLayout(),
				moveTo(0, 11, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0},
					head:      0,
					multi:     1,
					platetype: []string{"tipbox", "", "", "", "", "", "", ""},
					position:  []string{"tipbox_1", "", "", "", "", "", "", ""},
					well:      []string{"A12", "", "", "", "", "", "", ""},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{"H12"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name:  "OK - single tip LTR override (A1 -> D1)",
			Props: propsLTR,
			Setup: []*SetupFn{
				testLayout(),
				removeTipboxTips("tipbox_1", []string{"E1", "F1", "G1", "H1"}),
				moveTo(0, 0, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0},
					head:      0,
					multi:     1,
					platetype: []string{"tipbox", "", "", "", "", "", "", ""},
					position:  []string{"tipbox_1", "", "", "", "", "", "", ""},
					well:      []string{"A1", "", "", "", "", "", "", ""},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{"D1", "E1", "F1", "G1", "H1"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name:  "OK - single tip RTL override (A12 -> D12)",
			Props: propsRTL,
			Setup: []*SetupFn{
				testLayout(),
				removeTipboxTips("tipbox_1", []string{"E12", "F12", "G12", "H12"}),
				moveTo(0, 11, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0},
					head:      0,
					multi:     1,
					platetype: []string{"tipbox", "", "", "", "", "", "", ""},
					position:  []string{"tipbox_1", "", "", "", "", "", "", ""},
					well:      []string{"A12", "", "", "", "", "", "", ""},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{"D12", "E12", "F12", "G12", "H12"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name:  "8 tips LTR",
			Props: propsLTR,
			Setup: []*SetupFn{
				testLayout(),
				moveTo(0, 11, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0, 1, 2, 3, 4, 5, 6, 7},
					head:      0,
					multi:     8,
					platetype: []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					position:  []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					well:      []string{"A11", "B11", "C11", "D11", "E11", "F11", "G11", "H11"},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}, {1, "", 0}, {2, "", 0}, {3, "", 0}, {4, "", 0}, {5, "", 0}, {6, "", 0}, {7, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name:  "8 tips RTL",
			Props: propsRTL,
			Setup: []*SetupFn{
				testLayout(),
				moveTo(0, 11, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0, 1, 2, 3, 4, 5, 6, 7},
					head:      0,
					multi:     8,
					platetype: []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					position:  []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					well:      []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}, {1, "", 0}, {2, "", 0}, {3, "", 0}, {4, "", 0}, {5, "", 0}, {6, "", 0}, {7, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name:  "8 tips LTR override with split",
			Props: propsLTR,
			Setup: []*SetupFn{
				testLayout(),
				removeTipboxTips("tipbox_1", []string{"E1", "F1", "G1", "H1"}),
				moveTo(0, 11, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0, 1, 2, 3, 4, 5, 6, 7},
					head:      0,
					multi:     8,
					platetype: []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					position:  []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					well:      []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{"A1", "B1", "C1", "D1", "E2", "F2", "G2", "H2", "E1", "F1", "G1", "H1"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}, {1, "", 0}, {2, "", 0}, {3, "", 0}, {4, "", 0}, {5, "", 0}, {6, "", 0}, {7, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name:  "8 tips RTL override with split",
			Props: propsRTL,
			Setup: []*SetupFn{
				testLayout(),
				removeTipboxTips("tipbox_1", []string{"E12", "F12", "G12", "H12"}),
				moveTo(0, 0, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0, 1, 2, 3, 4, 5, 6, 7},
					head:      0,
					multi:     8,
					platetype: []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					position:  []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					well:      []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{"A12", "B12", "C12", "D12", "E11", "F11", "G11", "H11", "E12", "F12", "G12", "H12"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}, {1, "", 0}, {2, "", 0}, {3, "", 0}, {4, "", 0}, {5, "", 0}, {6, "", 0}, {7, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name:  "8 tips LTR override with boxchange",
			Props: propsLTR,
			Setup: []*SetupFn{
				testLayout(),
				removeTipboxTips("tipbox_1", []string{
					"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1",
					"A2", "B2", "C2", "D2", "E2", "F2", "G2", "H2",
					"A3", "B3", "C3", "D3", "E3", "F3", "G3", "H3",
					"A4", "B4", "C4", "D4", "E4", "F4", "G4", "H4",
					"A5", "B5", "C5", "D5", "E5", "F5", "G5", "H5",
					"A6", "B6", "C6", "D6", "E6", "F6", "G6", "H6",
					"A7", "B7", "C7", "D7", "E7", "F7", "G7", "H7",
					"A8", "B8", "C8", "D8", "E8", "F8", "G8", "H8",
					"A9", "B9", "C9", "D9", "E9", "F9", "G9", "H9",
					"A10", "B10", "C10", "D10", "E10", "F10", "G10", "H10",
					"A11", "B11", "C11", "D11", "E11", "F11", "G11", "H11",
					"A12", "B12", "C12", "D12",
				}),
				moveTo(0, 0, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0, 1, 2, 3, 4, 5, 6, 7},
					head:      0,
					multi:     8,
					platetype: []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					position:  []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					well:      []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}, {1, "", 0}, {2, "", 0}, {3, "", 0}, {4, "", 0}, {5, "", 0}, {6, "", 0}, {7, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name:  "8 tips LTR override without boxchange",
			Props: propsLTR,
			Setup: []*SetupFn{
				testLayout(),
				removeTipboxTips("tipbox_1", []string{
					"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1",
					"A2", "B2", "C2", "D2", "E2", "F2", "G2", "H2",
					"A3", "B3", "C3", "D3", "E3", "F3", "G3", "H3",
					"A4", "B4", "C4", "D4", "E4", "F4", "G4", "H4",
					"A5", "B5", "C5", "D5", "E5", "F5", "G5", "H5",
					"A6", "B6", "C6", "D6", "E6", "F6", "G6", "H6",
					"A7", "B7", "C7", "D7", "E7", "F7", "G7", "H7",
					"A8", "B8", "C8", "D8", "E8", "F8", "G8", "H8",
					"A9", "B9", "C9", "D9", "E9", "F9", "G9", "H9",
					"A10", "B10", "C10", "D10", "E10", "F10", "G10", "H10",
					"E11", "F11", "G11", "H11",
					"A12", "B12", "C12", "D12",
				}),
				moveTo(0, 0, mtp),
			},
			Instructions: []TestRobotInstruction{
				&LoadTips{
					channels:  []int{0, 1, 2, 3, 4, 5, 6, 7},
					head:      0,
					multi:     8,
					platetype: []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
					position:  []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
					well:      []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{
					"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1",
					"A2", "B2", "C2", "D2", "E2", "F2", "G2", "H2",
					"A3", "B3", "C3", "D3", "E3", "F3", "G3", "H3",
					"A4", "B4", "C4", "D4", "E4", "F4", "G4", "H4",
					"A5", "B5", "C5", "D5", "E5", "F5", "G5", "H5",
					"A6", "B6", "C6", "D6", "E6", "F6", "G6", "H6",
					"A7", "B7", "C7", "D7", "E7", "F7", "G7", "H7",
					"A8", "B8", "C8", "D8", "E8", "F8", "G8", "H8",
					"A9", "B9", "C9", "D9", "E9", "F9", "G9", "H9",
					"A10", "B10", "C10", "D10", "E10", "F10", "G10", "H10",
					"A11", "B11", "C11", "D11", "E11", "F11", "G11", "H11",
					"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12",
				}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}, {1, "", 0}, {2, "", 0}, {3, "", 0}, {4, "", 0}, {5, "", 0}, {6, "", 0}, {7, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
	}.Run(t)
}

func Test_UnloadTips(t *testing.T) {
	SimulatorTests{
		{
			Name: "OK - single tip",
			Setup: []*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"tipwaste", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"tipwaste", "", "", "", "", "", "", ""},
					head:         0,
				},
				&UnloadTips{
					channels:  []int{0},
					head:      0,
					multi:     1,
					platetype: []string{"tipwaste", "", "", "", "", "", "", ""},
					position:  []string{"tipwaste", "", "", "", "", "", "", ""},
					well:      []string{"A1", "", "", "", "", "", "", ""},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{}),
				tipwasteAssertion("tipwaste", 1),
			},
		},
		{
			Name: "OK - 8 tips",
			Setup: []*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste"},
					wellcoords:   []string{"A1", "A1", "A1", "A1", "A1", "A1", "A1", "A1"},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{-31.5, -22.5, -13.5, -4.5, 4.5, 13.5, 22.5, 31.5},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste"},
					head:         0,
				},
				&UnloadTips{
					channels:  []int{0, 1, 2, 3, 4, 5, 6, 7},
					head:      0,
					multi:     8,
					platetype: []string{"tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste"},
					position:  []string{"tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste"},
					well:      []string{"A1", "A1", "A1", "A1", "A1", "A1", "A1", "A1"},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{}),
				tipwasteAssertion("tipwaste", 8),
			},
		},
		//		commented out due to tips colliding with tipbox
		/*		{
					Name:  "OK - 8 tips back to a tipbox",
					Props: nil,
					Setup: []*SetupFn{
						testLayout(),
						removeTipboxTips("tipbox_1", []string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"}),
						preloadAdaptorTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}),
					},
					Instructions: []TestRobotInstruction{
						&Move{
							deckposition: []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
							wellcoords:   []string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"},
							reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
							offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
							offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
							offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
							plate_type:   []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
							head:         0,
						},
						&UnloadTips{
							channels:  []int{0, 1, 2, 3, 4, 5, 6, 7},
							head:      0,
							multi:     8,
							platetype: []string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
							position:  []string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
							well:      []string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"},
						},
					},
					Assertions: []*AssertionFn{
						tipboxAssertion("tipbox_1", []string{}),
						tipboxAssertion("tipbox_2", []string{}),
						adaptorAssertion(0, []tipDesc{}),
						tipwasteAssertion("tipwaste", 0),
					},
				},
		*/
		{
			Name:  "OK - independent tips",
			Props: IndependentLHProperties(),
			Setup: []*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}),
			},
			Instructions: []TestRobotInstruction{
				&UnloadTips{
					channels:  []int{0, 2, 4, 6},
					head:      0,
					multi:     4,
					platetype: []string{"tipwaste", "", "tipwaste", "", "tipwaste", "", "tipwaste", ""},
					position:  []string{"tipwaste", "", "tipwaste", "", "tipwaste", "", "tipwaste", ""},
					well:      []string{"A1", "", "A1", "", "A1", "", "A1", ""},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{
					{1, "", 0},
					{3, "", 0},
					{5, "", 0},
					{7, "", 0},
				}),
				tipwasteAssertion("tipwaste", 4),
			},
		},
		{
			Name: "can only unload all tips",
			Setup: []*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste"},
					wellcoords:   []string{"A1", "A1", "A1", "A1", "A1", "A1", "A1", "A1"},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{-31.5, -22.5, -13.5, -4.5, 4.5, 13.5, 22.5, 31.5},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste"},
					head:         0,
				},
				&UnloadTips{
					channels:  []int{0, 2, 4, 6},
					head:      0,
					multi:     4,
					platetype: []string{"tipwaste", "", "tipwaste", "", "tipwaste", "", "tipwaste", ""},
					position:  []string{"tipwaste", "", "tipwaste", "", "tipwaste", "", "tipwaste", ""},
					well:      []string{"A1", "", "A1", "", "A1", "", "A1", ""},
				},
			},
			ExpectedErrors: []string{
				"(err) UnloadTips[1]: Cannot unload tips from head0 channels 0,2,4,6 without unloading tips from channels 1,3,5,7 (head isn't independent)",
			},
		},
		{
			Name: "can't unload to a plate",
			Setup: []*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A12", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&UnloadTips{
					channels:  []int{0},
					head:      0,
					multi:     1,
					platetype: []string{"plate", "", "", "", "", "", "", ""},
					position:  []string{"input_1", "", "", "", "", "", "", ""},
					well:      []string{"A1", "", "", "", "", "", "", ""},
				},
			},
			ExpectedErrors: []string{
				"(err) UnloadTips[1]: Cannot unload tips to plate \"plate1\" at location input_1",
			},
		},
		{
			Name: "wrong well",
			Setup: []*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"tipwaste", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{-31.5, 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 0., 0., 0., 0., 0., 0., 0.},
					plate_type:   []string{"tipwaste", "", "", "", "", "", "", ""},
					head:         0,
				},
				&UnloadTips{
					channels:  []int{0},
					head:      0,
					multi:     1,
					platetype: []string{"tipwaste", "", "", "", "", "", "", ""},
					position:  []string{"tipwaste", "", "", "", "", "", "", ""},
					well:      []string{"B1", "", "", "", "", "", "", ""},
				},
			},
			ExpectedErrors: []string{
				"(err) UnloadTips[1]: Cannot unload to address B1 in tipwaste \"tipwaste\" size [1x1]",
			},
		},
	}.Run(t)
}

func Test_Aspirate(t *testing.T) {
	SimulatorTests{
		{
			Name: "OK - single channel",
			Setup: []*SetupFn{
				testLayout(),
				prefillWells("input_1", []string{"A1"}, "water", 200.),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Aspirate{
					volume:     []float64{100., 0., 0., 0., 0., 0., 0., 0.},
					overstroke: false,
					head:       0,
					multi:      1,
					platetype:  []string{"plate", "", "", "", "", "", "", ""},
					what:       []string{"water", "", "", "", "", "", "", ""},
					llf:        []bool{false, false, false, false, false, false, false, false},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "water", 100}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name: "OK - 8 channel",
			Setup: []*SetupFn{
				testLayout(),
				prefillWells("input_1", []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}, "water", 200.),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"},
					wellcoords:   []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"},
					head:         0,
				},
				&Aspirate{
					volume:     []float64{100., 100., 100., 100., 100., 100., 100., 100.},
					overstroke: false,
					head:       0,
					multi:      8,
					platetype:  []string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"},
					what:       []string{"water", "water", "water", "water", "water", "water", "water", "water"},
					llf:        []bool{false, false, false, false, false, false, false, false},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{
					{0, "water", 100},
					{1, "water", 100},
					{2, "water", 100},
					{3, "water", 100},
					{4, "water", 100},
					{5, "water", 100},
					{6, "water", 100},
					{7, "water", 100},
				}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name: "OK - 8 channel trough",
			Setup: []*SetupFn{
				testTroughLayout(),
				prefillWells("input_1", []string{"A1"}, "water", 10000.),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"},
					wellcoords:   []string{"A1", "A1", "A1", "A1", "A1", "A1", "A1", "A1"},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"trough", "trough", "trough", "trough", "trough", "trough", "trough", "trough"},
					head:         0,
				},
				&Aspirate{
					volume:     []float64{100., 100., 100., 100., 100., 100., 100., 100.},
					overstroke: false,
					head:       0,
					multi:      8,
					platetype:  []string{"trough", "trough", "trough", "trough", "trough", "trough", "trough", "trough"},
					what:       []string{"water", "water", "water", "water", "water", "water", "water", "water"},
					llf:        []bool{false, false, false, false, false, false, false, false},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{
					{0, "water", 100},
					{1, "water", 100},
					{2, "water", 100},
					{3, "water", 100},
					{4, "water", 100},
					{5, "water", 100},
					{6, "water", 100},
					{7, "water", 100},
				}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name: "Fail - take too much from trough",
			Setup: []*SetupFn{
				testTroughLayout(),
				prefillWells("input_1", []string{"A1"}, "water", 5400.),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"},
					wellcoords:   []string{"A1", "A1", "A1", "A1", "A1", "A1", "A1", "A1"},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"trough", "trough", "trough", "trough", "trough", "trough", "trough", "trough"},
					head:         0,
				},
				&Aspirate{
					volume:     []float64{100., 100., 100., 100., 100., 100., 100., 100.},
					overstroke: false,
					head:       0,
					multi:      8,
					platetype:  []string{"trough", "trough", "trough", "trough", "trough", "trough", "trough", "trough"},
					what:       []string{"water", "water", "water", "water", "water", "water", "water", "water"},
					llf:        []bool{false, false, false, false, false, false, false, false},
				},
			},
			ExpectedErrors: []string{
				"(warn) Aspirate[1]: 100 ul of water to head 0 channels 0-7: taking 800 ul from A1@trough1 which contains only 400 ul working volume, possible aspiration of residual (see ANTHA-2704)",
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{
					{0, "water", 100.0},
					{1, "water", 100.0},
					{2, "water", 100.0},
					{3, "water", 100.0},
					{4, "water", 100.0},
					{5, "water", 100.0},
					{6, "water", 100.0},
					{7, "water", 100.0},
				}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name: "Fail - Aspirate with no tip",
			Setup: []*SetupFn{
				testLayout(),
				prefillWells("input_1", []string{"A1", "B1"}, "water", 200.),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Aspirate{
					volume:     []float64{100., 100., 0., 0., 0., 0., 0., 0.},
					overstroke: false,
					head:       0,
					multi:      2,
					platetype:  []string{"plate", "plate", "", "", "", "", "", ""},
					what:       []string{"water", "water", "", "", "", "", "", ""},
					llf:        []bool{false, false, false, false, false, false, false, false},
				},
			},
			ExpectedErrors: []string{
				"(err) Aspirate[1]: 100 ul of water to head 0 channels 0-1: missing tip on channel 1",
			},
		},
		{
			Name: "Fail - Underfull tip",
			Setup: []*SetupFn{
				testLayout(),
				prefillWells("input_1", []string{"A1"}, "water", 200.),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Aspirate{
					volume:     []float64{20., 0., 0., 0., 0., 0., 0., 0.},
					overstroke: false,
					head:       0,
					multi:      1,
					platetype:  []string{"plate", "", "", "", "", "", "", ""},
					what:       []string{"water", "", "", "", "", "", "", ""},
					llf:        []bool{false, false, false, false, false, false, false, false},
				},
			},
			ExpectedErrors: []string{
				"(warn) Aspirate[1]: 20 ul of water to head 0 channel 0: minimum tip volume is 50 ul",
			},
		},
		{
			Name: "Fail - Overfull tip",
			Setup: []*SetupFn{
				testLayout(),
				prefillWells("input_1", []string{"A1", "B1", "C1", "D1", "E1", "F1"}, "water", 200.),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Aspirate{
					volume:     []float64{175., 0., 0., 0., 0., 0., 0., 0.},
					overstroke: false,
					head:       0,
					multi:      1,
					platetype:  []string{"plate", "", "", "", "", "", "", ""},
					what:       []string{"water", "", "", "", "", "", "", ""},
					llf:        []bool{false, false, false, false, false, false, false, false},
				},
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"B1", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Aspirate{
					volume:     []float64{175., 0., 0., 0., 0., 0., 0., 0.},
					overstroke: false,
					head:       0,
					multi:      1,
					platetype:  []string{"plate", "", "", "", "", "", "", ""},
					what:       []string{"water", "", "", "", "", "", "", ""},
					llf:        []bool{false, false, false, false, false, false, false, false},
				},
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"C1", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Aspirate{
					volume:     []float64{175., 0., 0., 0., 0., 0., 0., 0.},
					overstroke: false,
					head:       0,
					multi:      1,
					platetype:  []string{"plate", "", "", "", "", "", "", ""},
					what:       []string{"water", "", "", "", "", "", "", ""},
					llf:        []bool{false, false, false, false, false, false, false, false},
				},
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"D1", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Aspirate{
					volume:     []float64{175., 0., 0., 0., 0., 0., 0., 0.},
					overstroke: false,
					head:       0,
					multi:      1,
					platetype:  []string{"plate", "", "", "", "", "", "", ""},
					what:       []string{"water", "", "", "", "", "", "", ""},
					llf:        []bool{false, false, false, false, false, false, false, false},
				},
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"E1", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Aspirate{
					volume:     []float64{175., 0., 0., 0., 0., 0., 0., 0.},
					overstroke: false,
					head:       0,
					multi:      1,
					platetype:  []string{"plate", "", "", "", "", "", "", ""},
					what:       []string{"water", "", "", "", "", "", "", ""},
					llf:        []bool{false, false, false, false, false, false, false, false},
				},
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"F1", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Aspirate{
					volume:     []float64{175., 0., 0., 0., 0., 0., 0., 0.},
					overstroke: false,
					head:       0,
					multi:      1,
					platetype:  []string{"plate", "", "", "", "", "", "", ""},
					what:       []string{"water", "", "", "", "", "", "", ""},
					llf:        []bool{false, false, false, false, false, false, false, false},
				},
			},
			ExpectedErrors: []string{
				"(err) Aspirate[11]: 175 ul of water to head 0 channel 0: channel 0 contains 875 ul, command exceeds maximum volume 1000 ul",
			},
		},
		{
			Name: "Fail - non-independent head can only aspirate equal volumes",
			Setup: []*SetupFn{
				testLayout(),
				prefillWells("input_1", []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}, "water", 200.),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"},
					wellcoords:   []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"},
					head:         0,
				},
				&Aspirate{
					volume:     []float64{50., 60., 70., 80., 90., 100., 110., 120.},
					overstroke: false,
					head:       0,
					multi:      8,
					platetype:  []string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"},
					what:       []string{"water", "water", "water", "water", "water", "water", "water", "water"},
					llf:        []bool{false, false, false, false, false, false, false, false},
				},
			},
			ExpectedErrors: []string{
				"(err) Aspirate[1]: {50,60,70,80,90,100,110,120} ul of water to head 0 channels 0-7: channels cannot aspirate different volumes in non-independent head",
			},
		},
		{
			Name: "Fail - tip not in well",
			Setup: []*SetupFn{
				testLayout(),
				prefillWells("input_1", []string{"A1"}, "water", 200.),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{50., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Aspirate{
					volume:     []float64{100., 0., 0., 0., 0., 0., 0., 0.},
					overstroke: false,
					head:       0,
					multi:      1,
					platetype:  []string{"plate", "", "", "", "", "", "", ""},
					what:       []string{"water", "", "", "", "", "", "", ""},
					llf:        []bool{false, false, false, false, false, false, false, false},
				},
			},
			ExpectedErrors: []string{
				"(err) Aspirate[1]: 100 ul of water to head 0 channel 0: tip on channel 0 not in a well",
			},
		},
		{
			Name: "Fail - Well doesn't contain enough",
			Setup: []*SetupFn{
				testLayout(),
				prefillWells("input_1", []string{"A1"}, "water", 200.),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Aspirate{
					volume:     []float64{535, 0., 0., 0., 0., 0., 0., 0.},
					overstroke: false,
					head:       0,
					multi:      1,
					platetype:  []string{"plate", "", "", "", "", "", "", ""},
					what:       []string{"water", "", "", "", "", "", "", ""},
					llf:        []bool{false, false, false, false, false, false, false, false},
				},
			},
			ExpectedErrors: []string{
				"(err) Aspirate[1]: 535 ul of water to head 0 channel 0: taking 535 ul from A1@plate1 which contains only 200 ul working plus 5 ul residual volume (possibly a side effect of ANTHA-2704, try using a plate with a larger residual volume)",
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{
					{0, "water", 535.0},
				}),
				plateAssertion("input_1", []wellDesc{{"A1", "water", 0.0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name: "Fail - inadvertant aspiration",
			Setup: []*SetupFn{
				testLayout(),
				prefillWells("input_1", []string{"A1", "B1"}, "water", 200.),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1}),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Aspirate{
					volume:     []float64{98.6, 0., 0., 0., 0., 0., 0., 0.},
					overstroke: false,
					head:       0,
					multi:      1,
					platetype:  []string{"plate", "", "", "", "", "", "", ""},
					what:       []string{"water", "", "", "", "", "", "", ""},
					llf:        []bool{false, false, false, false, false, false, false, false},
				},
			},
			ExpectedErrors: []string{
				"(err) Aspirate[1]: 98.6 ul of water to head 0 channel 0: channel 1 will inadvertantly aspirate water from well B1@plate1 as head is not independent",
			},
		},
	}.Run(t)
}

func Test_Dispense(t *testing.T) {
	SimulatorTests{
		{
			Name: "OK - single channel",
			Setup: []*SetupFn{
				testLayout(),
				preloadFilledTips(0, "tipbox_1", []int{0}, "water", 100.),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Dispense{
					volume:    []float64{50., 0., 0., 0., 0., 0., 0., 0.},
					blowout:   []bool{false, false, false, false, false, false, false, false},
					head:      0,
					multi:     1,
					platetype: []string{"plate", "", "", "", "", "", "", ""},
					what:      []string{"water", "", "", "", "", "", "", ""},
					llf:       []bool{false, false, false, false, false, false, false, false},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				plateAssertion("input_1", []wellDesc{{"A1", "water", 50.}}),
				adaptorAssertion(0, []tipDesc{{0, "water", 50.}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name: "OK - mixing",
			Setup: []*SetupFn{
				testLayout(),
				preloadFilledTips(0, "tipbox_1", []int{0}, "water", 100.),
				prefillWells("input_1", []string{"A1"}, "green", 50.),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Dispense{
					volume:    []float64{50., 0., 0., 0., 0., 0., 0., 0.},
					blowout:   []bool{false, false, false, false, false, false, false, false},
					head:      0,
					multi:     1,
					platetype: []string{"plate", "", "", "", "", "", "", ""},
					what:      []string{"water", "", "", "", "", "", "", ""},
					llf:       []bool{false, false, false, false, false, false, false, false},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				plateAssertion("input_1", []wellDesc{{"A1", "0.5 v/v green+0.5 v/v water", 100.}}),
				adaptorAssertion(0, []tipDesc{{0, "water", 50.}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name: "OK - single channel slightly above well",
			Setup: []*SetupFn{
				testLayout(),
				preloadFilledTips(0, "tipbox_1", []int{0}, "water", 100.),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{1, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{3., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Dispense{
					volume:    []float64{50., 0., 0., 0., 0., 0., 0., 0.},
					blowout:   []bool{false, false, false, false, false, false, false, false},
					head:      0,
					multi:     1,
					platetype: []string{"plate", "", "", "", "", "", "", ""},
					what:      []string{"water", "", "", "", "", "", "", ""},
					llf:       []bool{false, false, false, false, false, false, false, false},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				plateAssertion("input_1", []wellDesc{{"A1", "water", 50.}}),
				adaptorAssertion(0, []tipDesc{{0, "water", 50.}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name: "OK - 8 channel",
			Setup: []*SetupFn{
				testLayout(),
				preloadFilledTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}, "water", 100.),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"},
					wellcoords:   []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"},
					head:         0,
				},
				&Dispense{
					volume:    []float64{50., 50., 50., 50., 50., 50., 50., 50.},
					blowout:   []bool{false, false, false, false, false, false, false, false},
					head:      0,
					multi:     8,
					platetype: []string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"},
					what:      []string{"water", "water", "water", "water", "water", "water", "water", "water"},
					llf:       []bool{false, false, false, false, false, false, false, false},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{
					{0, "water", 50.},
					{1, "water", 50.},
					{2, "water", 50.},
					{3, "water", 50.},
					{4, "water", 50.},
					{5, "water", 50.},
					{6, "water", 50.},
					{7, "water", 50.},
				}),
				plateAssertion("input_1", []wellDesc{
					{"A1", "water", 50.},
					{"B1", "water", 50.},
					{"C1", "water", 50.},
					{"D1", "water", 50.},
					{"E1", "water", 50.},
					{"F1", "water", 50.},
					{"G1", "water", 50.},
					{"H1", "water", 50.},
				}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name: "Fail - no tips",
			Setup: []*SetupFn{
				testLayout(),
				//preloadFilledTips(0, "tipbox_1", []int{0}, "water", 100.),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{1, 1, 1, 1, 1, 1, 1, 1},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Dispense{
					volume:    []float64{50., 0., 0., 0., 0., 0., 0., 0.},
					blowout:   []bool{false, false, false, false, false, false, false, false},
					head:      0,
					multi:     1,
					platetype: []string{"plate", "", "", "", "", "", "", ""},
					what:      []string{"water", "", "", "", "", "", "", ""},
					llf:       []bool{false, false, false, false, false, false, false, false},
				},
			},
			ExpectedErrors: []string{
				"(err) Dispense[1]: 50 ul of water from head 0 channel 0 to A1@plate1: no tip loaded on channel 0",
			},
		},
		{
			Name: "Fail - not enough in tip",
			Setup: []*SetupFn{
				testLayout(),
				preloadFilledTips(0, "tipbox_1", []int{0}, "water", 100.),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Dispense{
					volume:    []float64{150., 0., 0., 0., 0., 0., 0., 0.},
					blowout:   []bool{false, false, false, false, false, false, false, false},
					head:      0,
					multi:     1,
					platetype: []string{"plate", "", "", "", "", "", "", ""},
					what:      []string{"water", "", "", "", "", "", "", ""},
					llf:       []bool{false, false, false, false, false, false, false, false},
				},
			},
			ExpectedErrors: []string{
				"(warn) Dispense[1]: 150 ul of water from head 0 channel 0 to A1@plate1: tip on channel 0 contains only 100 ul, but blowout flag is false",
			},
		},
		{
			Name: "Fail - well over-full",
			Setup: []*SetupFn{
				testLayout(),
				preloadFilledTips(0, "tipbox_1", []int{0}, "water", 1000.),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Dispense{
					volume:    []float64{500., 0., 0., 0., 0., 0., 0., 0.},
					blowout:   []bool{false, false, false, false, false, false, false, false},
					head:      0,
					multi:     1,
					platetype: []string{"plate", "", "", "", "", "", "", ""},
					what:      []string{"water", "", "", "", "", "", "", ""},
					llf:       []bool{false, false, false, false, false, false, false, false},
				},
			},
			ExpectedErrors: []string{
				"(warn) Dispense[1]: 500 ul of water from head 0 channel 0 to A1@plate1: overfilling well A1@plate1 to 500 ul of 200 ul max volume",
			},
		},
		{
			Name: "Fail - not in a well",
			Setup: []*SetupFn{
				testLayout(),
				preloadFilledTips(0, "tipbox_1", []int{0}, "water", 100.),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{200., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Dispense{
					volume:    []float64{50., 0., 0., 0., 0., 0., 0., 0.},
					blowout:   []bool{false, false, false, false, false, false, false, false},
					head:      0,
					multi:     1,
					platetype: []string{"plate", "", "", "", "", "", "", ""},
					what:      []string{"water", "", "", "", "", "", "", ""},
					llf:       []bool{false, false, false, false, false, false, false, false},
				},
			},
			ExpectedErrors: []string{
				"(err) Dispense[1]: 50 ul of water from head 0 channel 0 to @<unnamed>: no well within 5 mm below tip on channel 0",
			},
		},
		{
			Name: "Fail - dispensing to tipwaste",
			Setup: []*SetupFn{
				testLayout(),
				preloadFilledTips(0, "tipbox_1", []int{0}, "water", 100.),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"tipwaste", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{1, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"tipwaste", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Dispense{
					volume:    []float64{50., 0., 0., 0., 0., 0., 0., 0.},
					blowout:   []bool{false, false, false, false, false, false, false, false},
					head:      0,
					multi:     1,
					platetype: []string{"tipwaste", "", "", "", "", "", "", ""},
					what:      []string{"water", "", "", "", "", "", "", ""},
					llf:       []bool{false, false, false, false, false, false, false, false},
				},
			},
			ExpectedErrors: []string{
				"(warn) Dispense[1]: 50 ul of water from head 0 channel 0 to A1@tipwaste: dispensing to tipwaste",
			},
		},
		{
			Name: "fail - independence other tips in wells",
			Setup: []*SetupFn{
				testLayout(),
				preloadFilledTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}, "water", 100.),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Dispense{
					volume:    []float64{50., 0, 0, 0, 0, 0, 0, 0},
					blowout:   []bool{false, false, false, false, false, false, false, false},
					head:      0,
					multi:     1,
					platetype: []string{"plate", "", "", "", "", "", "", ""},
					what:      []string{"water", "", "", "", "", "", "", ""},
					llf:       []bool{false, false, false, false, false, false, false, false},
				},
			},
			ExpectedErrors: []string{
				"(err) Dispense[1]: 50 ul of water from head 0 channel 0 to A1@plate1: must also dispense 50 ul from channels 1-7 as head is not independent",
			},
		},
		{
			Name: "fail - independence other tip not in a well",
			Setup: []*SetupFn{
				testLayout(),
				preloadFilledTips(0, "tipbox_1", []int{0, 1}, "water", 100.),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"H1", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Dispense{
					volume:    []float64{50., 0, 0, 0, 0, 0, 0, 0},
					blowout:   []bool{false, false, false, false, false, false, false, false},
					head:      0,
					multi:     1,
					platetype: []string{"plate", "", "", "", "", "", "", ""},
					what:      []string{"water", "", "", "", "", "", "", ""},
					llf:       []bool{false, false, false, false, false, false, false, false},
				},
			},
			ExpectedErrors: []string{
				"(err) Dispense[1]: 50 ul of water from head 0 channel 0 to H1@plate1: must also dispense 50 ul from channel 1 as head is not independent",
			},
		},
		{
			Name: "Fail - independence, different volumes",
			Setup: []*SetupFn{
				testLayout(),
				preloadFilledTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}, "water", 100.),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"},
					wellcoords:   []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"},
					head:         0,
				},
				&Dispense{
					volume:    []float64{50., 60., 50., 50., 50., 50., 50., 50.},
					blowout:   []bool{false, false, false, false, false, false, false, false},
					head:      0,
					multi:     8,
					platetype: []string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"},
					what:      []string{"water", "water", "water", "water", "water", "water", "water", "water"},
					llf:       []bool{false, false, false, false, false, false, false, false},
				},
			},
			ExpectedErrors: []string{
				"(err) Dispense[1]: {50,60,50,50,50,50,50,50} ul of water from head 0 channels 0-7 to A1-H1@plate1: channels cannot dispense different volumes in non-independent head",
			},
		},
	}.Run(t)
}

func Test_Mix(t *testing.T) {
	SimulatorTests{
		{
			Name: "OK - single channel",
			Setup: []*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
				prefillWells("input_1", []string{"A1"}, "water", 200.),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Mix{
					head:      0,
					volume:    []float64{50., 0., 0., 0., 0., 0., 0., 0.},
					platetype: []string{"plate", "", "", "", "", "", "", ""},
					cycles:    []int{5, 0, 0, 0, 0, 0, 0, 0},
					multi:     1,
					what:      []string{"water", "", "", "", "", "", "", ""},
					blowout:   []bool{false, false, false, false, false, false, false, false},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				plateAssertion("input_1", []wellDesc{{"A1", "water", 200.}}),
				adaptorAssertion(0, []tipDesc{{0, "", 0.}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name: "OK - 8 channel",
			Setup: []*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}),
				prefillWells("input_1", []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}, "water", 200.),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"},
					wellcoords:   []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"},
					head:         0,
				},
				&Mix{
					head:      0,
					volume:    []float64{50., 50., 50., 50., 50., 50., 50., 50.},
					platetype: []string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"},
					cycles:    []int{5, 5, 5, 5, 5, 5, 5, 5},
					multi:     8,
					what:      []string{"water", "water", "water", "water", "water", "water", "water", "water"},
					blowout:   []bool{false, false, false, false, false, false, false, false},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				plateAssertion("input_1", []wellDesc{
					{"A1", "water", 200.},
					{"B1", "water", 200.},
					{"C1", "water", 200.},
					{"D1", "water", 200.},
					{"E1", "water", 200.},
					{"F1", "water", 200.},
					{"G1", "water", 200.},
					{"H1", "water", 200.},
				}),
				adaptorAssertion(0, []tipDesc{
					{0, "", 0.},
					{1, "", 0.},
					{2, "", 0.},
					{3, "", 0.},
					{4, "", 0.},
					{5, "", 0.},
					{6, "", 0.},
					{7, "", 0.},
				}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			Name: "Fail - independece problems",
			Setup: []*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}),
				prefillWells("input_1", []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}, "water", 200.),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"},
					wellcoords:   []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"},
					head:         0,
				},
				&Mix{
					head:      0,
					volume:    []float64{50., 60., 50., 50., 50., 50., 50., 50.},
					platetype: []string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"},
					cycles:    []int{5, 5, 5, 5, 5, 2, 2, 2},
					multi:     8,
					what:      []string{"water", "water", "water", "water", "water", "water", "water", "water"},
					blowout:   []bool{false, false, false, false, false, false, false, false},
				},
			},
			ExpectedErrors: []string{
				"(err) Mix[1]: {50,60,50,50,50,50,50,50} ul {5,5,5,5,5,2,2,2} times in wells A1,B1,C1,D1,E1,F1,G1,H1 of plate \"plate1\": cannot manipulate different volumes with non-independent head",
				"(err) Mix[1]: {50,60,50,50,50,50,50,50} ul {5,5,5,5,5,2,2,2} times in wells A1,B1,C1,D1,E1,F1,G1,H1 of plate \"plate1\": cannot vary number of mix cycles with non-independent head",
			},
		},
		{
			Name: "Fail - wrong platetype",
			Setup: []*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
				prefillWells("input_1", []string{"A1"}, "water", 200.),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{0, 0, 0, 0, 0, 0, 0, 0},
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Mix{
					head:      0,
					volume:    []float64{50., 0., 0., 0., 0., 0., 0., 0.},
					platetype: []string{"notaplate", "", "", "", "", "", "", ""},
					cycles:    []int{5, 0, 0, 0, 0, 0, 0, 0},
					multi:     1,
					what:      []string{"water", "", "", "", "", "", "", ""},
					blowout:   []bool{false, false, false, false, false, false, false, false},
				},
			},
			ExpectedErrors: []string{
				"(warn) Mix[1]: 50 ul 5 times in well A1 of plate \"plate1\": plate \"plate1\" is of type \"plate\", not \"notaplate\"",
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				plateAssertion("input_1", []wellDesc{{"A1", "water", 200.}}),
				adaptorAssertion(0, []tipDesc{{0, "", 0.}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
	}.Run(t)
}

func Test_LiquidLevelFollow(t *testing.T) {
	SimulatorTests{
		{
			Name: "OK - single channel",
			Setup: []*SetupFn{
				testLayout(),
				prefillWells("input_1", []string{"A1"}, "water", 200.),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
			},
			Instructions: []TestRobotInstruction{
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A1", "", "", "", "", "", "", ""},
					reference:    []int{2, 2, 2, 2, 2, 2, 2, 2}, //2 == liquidlevel
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Aspirate{
					volume:     []float64{100., 0., 0., 0., 0., 0., 0., 0.},
					overstroke: false,
					head:       0,
					multi:      1,
					platetype:  []string{"plate", "", "", "", "", "", "", ""},
					what:       []string{"water", "", "", "", "", "", "", ""},
					llf:        []bool{true, true, true, true, true, true, true, true},
				},
				&Move{
					deckposition: []string{"input_1", "", "", "", "", "", "", ""},
					wellcoords:   []string{"A2", "", "", "", "", "", "", ""},
					reference:    []int{2, 2, 2, 2, 2, 2, 2, 2}, //2 == liquidlevel
					offsetX:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetY:      []float64{0., 0., 0., 0., 0., 0., 0., 0.},
					offsetZ:      []float64{1., 1., 1., 1., 1., 1., 1., 1.},
					plate_type:   []string{"plate", "", "", "", "", "", "", ""},
					head:         0,
				},
				&Dispense{
					volume:    []float64{100., 0., 0., 0., 0., 0., 0., 0.},
					blowout:   []bool{false, false, false, false, false, false, false, false},
					head:      0,
					multi:     1,
					platetype: []string{"plate", "", "", "", "", "", "", ""},
					what:      []string{"water", "", "", "", "", "", "", ""},
					llf:       []bool{true, true, true, true, true, true, true, true},
				},
			},
			Assertions: []*AssertionFn{
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "water", 0}}),
				plateAssertion("input_1", []wellDesc{{"A1", "water", 99.5}, {"A2", "water", 100}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
	}.Run(t)
}

func component(name string) *wtype.Liquid {
	A := wtype.NewLHComponent()
	A.CName = name
	A.Type = wtype.LTWater
	A.Smax = 9999
	return A
}

func Test_Workflow(t *testing.T) {

	get8tips := func(column int) []TestRobotInstruction {
		wc := make([]string, 8)
		for i := range wc {
			c := wtype.WellCoords{X: column, Y: i}
			wc[i] = c.FormatA1()
		}
		return []TestRobotInstruction{
			&Move{
				[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
				wc,
				[]int{1, 1, 1, 1, 1, 1, 1, 1},
				[]float64{0, 0, 0, 0, 0, 0, 0, 0},
				[]float64{0, 0, 0, 0, 0, 0, 0, 0},
				[]float64{5, 5, 5, 5, 5, 5, 5, 5},
				[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
				0,
			},
			&LoadTips{
				[]int{0, 1, 2, 3, 4, 5, 6, 7},
				0,
				8,
				[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},
				[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"},
				wc,
			},
		}
	}

	get1tip := func(wc string) []TestRobotInstruction {
		return []TestRobotInstruction{
			&Move{
				[]string{"tipbox_1", "", "", "", "", "", "", ""},
				[]string{wc, "", "", "", "", "", "", ""},
				[]int{1, 1, 1, 1, 1, 1, 1, 1},
				[]float64{0, 0, 0, 0, 0, 0, 0, 0},
				[]float64{0, 0, 0, 0, 0, 0, 0, 0},
				[]float64{5, 5, 5, 5, 5, 5, 5, 5},
				[]string{"tipbox", "", "", "", "", "", "", ""},
				0,
			},
			&LoadTips{
				[]int{0},
				0,
				1,
				[]string{"tipbox", "", "", "", "", "", "", ""},
				[]string{"tipbox_1", "", "", "", "", "", "", ""},
				[]string{wc, "", "", "", "", "", "", ""},
			},
		}
	}

	dropTips := func(channels []int) []TestRobotInstruction {
		pl := make([]string, 8)
		pt := make([]string, 8)
		wl := make([]string, 8)
		for _, ch := range channels {
			pl[ch] = "tipwaste"
			pt[ch] = "tipwaste"
			wl[ch] = "A1"
		}
		return []TestRobotInstruction{
			&Move{
				pl,
				wl,
				[]int{1, 1, 1, 1, 1, 1, 1, 1},
				[]float64{0, 0, 0, 0, 0, 0, 0, 0},
				[]float64{-31.5, -22.5, -13.5, -4.5, 4.5, 13.5, 22.5, 31.5},
				[]float64{5, 5, 5, 5, 5, 5, 5, 5},
				pt,
				0,
			},
			&UnloadTips{channels, 0, len(channels), pl, pt, wl},
		}
	}

	suck := func(plateloc string, wells []string, what []string, volume float64) []TestRobotInstruction {
		multi := 0
		for i := range wells {
			if wells[i] != "" || what[i] != "" {
				multi += 1
			}
		}
		pl := make([]string, 8)
		pt := make([]string, 8)
		v := make([]float64, 8)
		for i, w := range wells {
			if w != "" {
				pl[i] = plateloc
				pt[i] = "plate"
				v[i] = volume
			}
		}
		return []TestRobotInstruction{
			&Move{
				pl,
				wells,
				[]int{0, 0, 0, 0, 0, 0, 0, 0},
				[]float64{0, 0, 0, 0, 0, 0, 0, 0},
				[]float64{0, 0, 0, 0, 0, 0, 0, 0},
				[]float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5},
				pt,
				0,
			},
			&Aspirate{
				v,
				false,
				0,
				multi,
				pt,
				what,
				[]bool{false, false, false, false, false, false, false, false},
			},
		}
	}
	blow := func(plateloc string, wells []string, what []string, volume float64) []TestRobotInstruction {
		multi := 0
		for i := range wells {
			if wells[i] != "" || what[i] != "" {
				multi += 1
			}
		}
		pl := make([]string, 8)
		pt := make([]string, 8)
		v := make([]float64, 8)
		for i, w := range wells {
			if w != "" {
				pl[i] = plateloc
				pt[i] = "plate"
				v[i] = volume
			}
		}
		return []TestRobotInstruction{
			&Move{
				pl,
				wells,
				[]int{1, 1, 1, 1, 1, 1, 1, 1},
				[]float64{0, 0, 0, 0, 0, 0, 0, 0},
				[]float64{0, 0, 0, 0, 0, 0, 0, 0},
				[]float64{1, 1, 1, 1, 1, 1, 1, 1},
				pt,
				0,
			},
			&Dispense{
				v,
				[]bool{false, false, false, false, false, false, false, false},
				0,
				multi,
				pt,
				what,
				[]bool{false, false, false, false, false, false, false, false},
			},
		}
	}

	//plates
	input_plate := defaultLHPlate("input")
	output_plate := defaultLHPlate("output")

	//tips - using small tipbox so I don't have to worry about using different tips
	tipbox := smallLHTipbox("tipbox")

	//tipwaste
	tipwaste := defaultLHTipwaste("tipwaste")

	//setup the input plate
	wc := wtype.MakeWellCoords("A1")
	comp := []*wtype.Liquid{
		component("water"),
		component("red"),
		component("green"),
		component("water"),
	}
	for x := 0; x < len(comp); x++ {
		wc.X = x
		comp[x].Vol = 200.
		for y := 0; y < 8; y++ {
			wc.Y = y
			well := input_plate.GetChildByAddress(wc).(*wtype.LHWell)
			err := well.AddComponent(comp[x].Dup())
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	//add all the plates to the robot
	inst := []TestRobotInstruction{
		&Initialize{},
		&AddPlateTo{"input_1", input_plate, "input"},
		&AddPlateTo{"output_1", output_plate, "output"},
		&AddPlateTo{"tipbox_1", tipbox, "tipbox"},
		&AddPlateTo{"tipwaste", tipwaste, "tipwaste"},
	}

	//make the green gradient
	green := make([]float64, 8)
	for i := range green {
		green[i] = float64(7-i) / 7
	}
	inst = append(inst, get1tip("H12")...)
	for y := 0; y < 8; y++ {
		if green[y] == 0. {
			continue
		}
		wc.X = 2
		wc.Y = y
		inst = append(inst, suck("input_1",
			[]string{wc.FormatA1(), "", "", "", "", "", "", ""},
			[]string{"water", "", "", "", "", "", "", ""},
			195.*green[y])...)
		wc.X = 4
		inst = append(inst, blow("input_1",
			[]string{wc.FormatA1(), "", "", "", "", "", "", ""},
			[]string{"water", "", "", "", "", "", "", ""},
			195.*green[y])...)
	}
	inst = append(inst, dropTips([]int{0})...)
	inst = append(inst, get1tip("G12")...)
	for y := 0; y < 8; y++ {
		if (1 - green[y]) == 0. {
			continue
		}
		wc.X = 3
		wc.Y = y
		inst = append(inst, suck("input_1",
			[]string{wc.FormatA1(), "", "", "", "", "", "", ""},
			[]string{"water", "", "", "", "", "", "", ""},
			195.*(1-green[y]))...)
		wc.X = 4
		inst = append(inst, blow("input_1",
			[]string{wc.FormatA1(), "", "", "", "", "", "", ""},
			[]string{"water", "", "", "", "", "", "", ""},
			195.*(1-green[y]))...)
	}
	inst = append(inst, dropTips([]int{0})...)

	//make the red gradient
	red := make([]float64, 12)
	for i := range red {
		red[i] = float64(11-i) / 11
	}
	inst = append(inst, get8tips(10)...)

	from_wells := make([]string, 8)
	for i := range from_wells {
		wc.X = 1
		wc.Y = i
		from_wells[i] = wc.FormatA1()
	}
	for x := 0; x < 12; x++ {
		if red[x] == 0. {
			continue
		}
		to_wells := make([]string, 8)
		for i := range to_wells {
			wc.X = x
			wc.Y = i
			to_wells[i] = wc.FormatA1()
		}

		inst = append(inst, suck("input_1",
			from_wells,
			[]string{"water", "water", "water", "water", "water", "water", "water", "water"},
			5.*red[x])...)
		inst = append(inst, blow("output_1",
			to_wells,
			[]string{"water", "water", "water", "water", "water", "water", "water", "water"},
			5.*red[x])...)
	}
	inst = append(inst, dropTips([]int{0, 1, 2, 3, 4, 5, 6, 7})...)

	//transfer the green gradient
	inst = append(inst, get8tips(9)...)
	from_wells = make([]string, 8)
	for i := range from_wells {
		wc.X = 4
		wc.Y = i
		from_wells[i] = wc.FormatA1()
	}
	for x := 0; x < 12; x++ {
		to_wells := make([]string, 8)
		for i := range to_wells {
			wc.X = x
			wc.Y = i
			to_wells[i] = wc.FormatA1()
		}

		inst = append(inst, suck("input_1",
			from_wells,
			[]string{"water", "water", "water", "water", "water", "water", "water", "water"},
			5.)...)
		inst = append(inst, blow("output_1",
			to_wells,
			[]string{"water", "water", "water", "water", "water", "water", "water", "water"},
			5.)...)
	}
	inst = append(inst, dropTips([]int{0, 1, 2, 3, 4, 5, 6, 7})...)

	//make up to 20ul
	inst = append(inst, get8tips(8)...)
	from_wells = make([]string, 8)
	for i := range from_wells {
		wc.X = 0
		wc.Y = i
		from_wells[i] = wc.FormatA1()
	}
	for x := 0; x < 12; x++ {
		to_wells := make([]string, 8)
		for i := range to_wells {
			wc.X = x
			wc.Y = i
			to_wells[i] = wc.FormatA1()
		}

		inst = append(inst, suck("input_1",
			from_wells,
			[]string{"water", "water", "water", "water", "water", "water", "water", "water"},
			10.+5.*(1.-red[x]))...)
		inst = append(inst, blow("output_1",
			to_wells,
			[]string{"water", "water", "water", "water", "water", "water", "water", "water"},
			10.+5.*(1.-red[x]))...)
	}
	inst = append(inst, dropTips([]int{0, 1, 2, 3, 4, 5, 6, 7})...)

	//and finally
	inst = append(inst, &Finalize{})

	(&SimulatorTest{
		Name:         "Run Workflow",
		Setup:        []*SetupFn{},
		Instructions: inst,
		Assertions: []*AssertionFn{
			tipboxAssertion("tipbox_1", []string{
				"H12", "G12",
				"A11", "B11", "C11", "D11", "E11", "F11", "G11", "H11",
				"A10", "B10", "C10", "D10", "E10", "F10", "G10", "H10",
				"A9", "B9", "C9", "D9", "E9", "F9", "G9", "H9",
			}),
			plateAssertion("input_1", []wellDesc{
				{"A1", "water", 50.0 - 12.0*0.5},
				{"B1", "water", 50.0 - 12.0*0.5},
				{"C1", "water", 50.0 - 12.0*0.5},
				{"D1", "water", 50.0 - 12.0*0.5},
				{"E1", "water", 50.0 - 12.0*0.5},
				{"F1", "water", 50.0 - 12.0*0.5},
				{"G1", "water", 50.0 - 12.0*0.5},
				{"H1", "water", 50.0 - 12.0*0.5},
				{"A2", "red", 170.0 - 11.0*0.5},
				{"B2", "red", 170.0 - 11.0*0.5},
				{"C2", "red", 170.0 - 11.0*0.5},
				{"D2", "red", 170.0 - 11.0*0.5},
				{"E2", "red", 170.0 - 11.0*0.5},
				{"F2", "red", 170.0 - 11.0*0.5},
				{"G2", "red", 170.0 - 11.0*0.5},
				{"H2", "red", 170.0 - 11.0*0.5},
				{"A3", "green", 5.0 - 0.5},
				{"B3", "green", 32.86 - 0.5},
				{"C3", "green", 60.71 - 0.5},
				{"D3", "green", 88.57 - 0.5},
				{"E3", "green", 116.43 - 0.5},
				{"F3", "green", 144.29 - 0.5},
				{"G3", "green", 172.14 - 0.5},
				{"H3", "green", 200.},
				{"A4", "water", 200.},
				{"B4", "water", 172.14 - 0.5},
				{"C4", "water", 144.29 - 0.5},
				{"D4", "water", 116.43 - 0.5},
				{"E4", "water", 88.57 - 0.5},
				{"F4", "water", 60.71 - 0.5},
				{"G4", "water", 32.86 - 0.5},
				{"H4", "water", 5. - 0.5},
				{"A5", "green", 135.0 - 12.0*0.5},
				{"B5", "0.857 v/v green+0.143 v/v water", 135.0 - 12.0*0.5},
				{"C5", "0.714 v/v green+0.286 v/v water", 135.0 - 12.0*0.5},
				{"D5", "0.571 v/v green+0.429 v/v water", 135.0 - 12.0*0.5},
				{"E5", "0.429 v/v green+0.571 v/v water", 135.0 - 12.0*0.5},
				{"F5", "0.286 v/v green+0.714 v/v water", 135.0 - 12.0*0.5},
				{"G5", "0.143 v/v green+0.857 v/v water", 135.0 - 12.0*0.5},
				{"H5", "water", 135.0 - 12.0*0.5},
			}),
			plateAssertion("output_1", []wellDesc{
				{"A1", "0.25 v/v green+0.25 v/v red+0.5 v/v water", 20.},
				{"B1", "0.214 v/v green+0.25 v/v red+0.536 v/v water", 20.},
				{"C1", "0.179 v/v green+0.25 v/v red+0.571 v/v water", 20.},
				{"D1", "0.143 v/v green+0.25 v/v red+0.607 v/v water", 20.},
				{"E1", "0.107 v/v green+0.25 v/v red+0.643 v/v water", 20.},
				{"F1", "0.0714 v/v green+0.25 v/v red+0.679 v/v water", 20.},
				{"G1", "0.0357 v/v green+0.25 v/v red+0.714 v/v water", 20.},
				{"H1", "0.25 v/v red+0.75 v/v water", 20.},
				{"A2", "0.25 v/v green+0.227 v/v red+0.523 v/v water", 20.},
				{"B2", "0.214 v/v green+0.227 v/v red+0.558 v/v water", 20.},
				{"C2", "0.179 v/v green+0.227 v/v red+0.594 v/v water", 20.},
				{"D2", "0.143 v/v green+0.227 v/v red+0.63 v/v water", 20.},
				{"E2", "0.107 v/v green+0.227 v/v red+0.666 v/v water", 20.},
				{"F2", "0.0714 v/v green+0.227 v/v red+0.701 v/v water", 20.},
				{"G2", "0.0357 v/v green+0.227 v/v red+0.737 v/v water", 20.},
				{"H2", "0.227 v/v red+0.773 v/v water", 20.},
				{"A3", "0.25 v/v green+0.205 v/v red+0.545 v/v water", 20.},
				{"B3", "0.214 v/v green+0.205 v/v red+0.581 v/v water", 20.},
				{"C3", "0.179 v/v green+0.205 v/v red+0.617 v/v water", 20.},
				{"D3", "0.143 v/v green+0.205 v/v red+0.653 v/v water", 20.},
				{"E3", "0.107 v/v green+0.205 v/v red+0.688 v/v water", 20.},
				{"F3", "0.0714 v/v green+0.205 v/v red+0.724 v/v water", 20.},
				{"G3", "0.0357 v/v green+0.205 v/v red+0.76 v/v water", 20.},
				{"H3", "0.205 v/v red+0.795 v/v water", 20.},
				{"A4", "0.25 v/v green+0.182 v/v red+0.568 v/v water", 20.},
				{"B4", "0.214 v/v green+0.182 v/v red+0.604 v/v water", 20.},
				{"C4", "0.179 v/v green+0.182 v/v red+0.64 v/v water", 20.},
				{"D4", "0.143 v/v green+0.182 v/v red+0.675 v/v water", 20.},
				{"E4", "0.107 v/v green+0.182 v/v red+0.711 v/v water", 20.},
				{"F4", "0.0714 v/v green+0.182 v/v red+0.747 v/v water", 20.},
				{"G4", "0.0357 v/v green+0.182 v/v red+0.782 v/v water", 20.},
				{"H4", "0.182 v/v red+0.818 v/v water", 20.},
				{"A5", "0.25 v/v green+0.159 v/v red+0.591 v/v water", 20.},
				{"B5", "0.214 v/v green+0.159 v/v red+0.627 v/v water", 20.},
				{"C5", "0.179 v/v green+0.159 v/v red+0.662 v/v water", 20.},
				{"D5", "0.143 v/v green+0.159 v/v red+0.698 v/v water", 20.},
				{"E5", "0.107 v/v green+0.159 v/v red+0.734 v/v water", 20.},
				{"F5", "0.0714 v/v green+0.159 v/v red+0.769 v/v water", 20.},
				{"G5", "0.0357 v/v green+0.159 v/v red+0.805 v/v water", 20.},
				{"H5", "0.159 v/v red+0.841 v/v water", 20.},
				{"A6", "0.25 v/v green+0.136 v/v red+0.614 v/v water", 20.},
				{"B6", "0.214 v/v green+0.136 v/v red+0.649 v/v water", 20.},
				{"C6", "0.179 v/v green+0.136 v/v red+0.685 v/v water", 20.},
				{"D6", "0.143 v/v green+0.136 v/v red+0.721 v/v water", 20.},
				{"E6", "0.107 v/v green+0.136 v/v red+0.756 v/v water", 20.},
				{"F6", "0.0714 v/v green+0.136 v/v red+0.792 v/v water", 20.},
				{"G6", "0.0357 v/v green+0.136 v/v red+0.828 v/v water", 20.},
				{"H6", "0.136 v/v red+0.864 v/v water", 20.},
				{"A7", "0.25 v/v green+0.114 v/v red+0.636 v/v water", 20.},
				{"B7", "0.214 v/v green+0.114 v/v red+0.672 v/v water", 20.},
				{"C7", "0.179 v/v green+0.114 v/v red+0.708 v/v water", 20.},
				{"D7", "0.143 v/v green+0.114 v/v red+0.744 v/v water", 20.},
				{"E7", "0.107 v/v green+0.114 v/v red+0.779 v/v water", 20.},
				{"F7", "0.0714 v/v green+0.114 v/v red+0.815 v/v water", 20.},
				{"G7", "0.0357 v/v green+0.114 v/v red+0.851 v/v water", 20.},
				{"H7", "0.114 v/v red+0.886 v/v water", 20.},
				{"A8", "0.25 v/v green+0.0909 v/v red+0.659 v/v water", 20.},
				{"B8", "0.214 v/v green+0.0909 v/v red+0.695 v/v water", 20.},
				{"C8", "0.179 v/v green+0.0909 v/v red+0.731 v/v water", 20.},
				{"D8", "0.143 v/v green+0.0909 v/v red+0.766 v/v water", 20.},
				{"E8", "0.107 v/v green+0.0909 v/v red+0.802 v/v water", 20.},
				{"F8", "0.0714 v/v green+0.0909 v/v red+0.838 v/v water", 20.},
				{"G8", "0.0357 v/v green+0.0909 v/v red+0.873 v/v water", 20.},
				{"H8", "0.0909 v/v red+0.909 v/v water", 20.},
				{"A9", "0.25 v/v green+0.0682 v/v red+0.682 v/v water", 20.},
				{"B9", "0.214 v/v green+0.0682 v/v red+0.718 v/v water", 20.},
				{"C9", "0.179 v/v green+0.0682 v/v red+0.753 v/v water", 20.},
				{"D9", "0.143 v/v green+0.0682 v/v red+0.789 v/v water", 20.},
				{"E9", "0.107 v/v green+0.0682 v/v red+0.825 v/v water", 20.},
				{"F9", "0.0714 v/v green+0.0682 v/v red+0.86 v/v water", 20.},
				{"G9", "0.0357 v/v green+0.0682 v/v red+0.896 v/v water", 20.},
				{"H9", "0.0682 v/v red+0.932 v/v water", 20.},
				{"A10", "0.25 v/v green+0.0455 v/v red+0.705 v/v water", 20.},
				{"B10", "0.214 v/v green+0.0455 v/v red+0.74 v/v water", 20.},
				{"C10", "0.179 v/v green+0.0455 v/v red+0.776 v/v water", 20.},
				{"D10", "0.143 v/v green+0.0455 v/v red+0.812 v/v water", 20.},
				{"E10", "0.107 v/v green+0.0455 v/v red+0.847 v/v water", 20.},
				{"F10", "0.0714 v/v green+0.0455 v/v red+0.883 v/v water", 20.},
				{"G10", "0.0357 v/v green+0.0455 v/v red+0.919 v/v water", 20.},
				{"H10", "0.0455 v/v red+0.955 v/v water", 20.},
				{"A11", "0.25 v/v green+0.0227 v/v red+0.727 v/v water", 20.},
				{"B11", "0.214 v/v green+0.0227 v/v red+0.763 v/v water", 20.},
				{"C11", "0.179 v/v green+0.0227 v/v red+0.799 v/v water", 20.},
				{"D11", "0.143 v/v green+0.0227 v/v red+0.834 v/v water", 20.},
				{"E11", "0.107 v/v green+0.0227 v/v red+0.87 v/v water", 20.},
				{"F11", "0.0714 v/v green+0.0227 v/v red+0.906 v/v water", 20.},
				{"G11", "0.0357 v/v green+0.0227 v/v red+0.942 v/v water", 20.},
				{"H11", "0.0227 v/v red+0.977 v/v water", 20.},
				{"A12", "0.25 v/v green+0.75 v/v water", 20.},
				{"B12", "0.214 v/v green+0.786 v/v water", 20.},
				{"C12", "0.179 v/v green+0.821 v/v water", 20.},
				{"D12", "0.143 v/v green+0.857 v/v water", 20.},
				{"E12", "0.107 v/v green+0.893 v/v water", 20.},
				{"F12", "0.0714 v/v green+0.929 v/v water", 20.},
				{"G12", "0.0357 v/v green+0.964 v/v water", 20.},
				{"H12", "1 v/v water", 20.},
			}),
			adaptorAssertion(0, []tipDesc{}),
			tipwasteAssertion("tipwaste", 26),
		},
	}).Run(t)

}
