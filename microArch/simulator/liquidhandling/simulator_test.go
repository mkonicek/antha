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

package liquidhandling_test

import (
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	lh "github.com/antha-lang/antha/microArch/simulator/liquidhandling"
)

func TestUnknownLocations(t *testing.T) {
	tests := make([]SimulatorTest, 0)

	lhp := default_lhproperties()
	lhp.Tip_preferences = append(lhp.Tip_preferences, "undefined_tip_pref")
	tests = append(tests, SimulatorTest{"Undefined Tip_preference", lhp, nil, nil,
		[]string{"(warn) NewVirtualLiquidHandler: Undefined location \"undefined_tip_pref\" referenced in tip preferences"},
		nil})

	lhp = default_lhproperties()
	lhp.Input_preferences = append(lhp.Tip_preferences, "undefined_input_pref")
	tests = append(tests, SimulatorTest{"passing undefined Input_preference", lhp, nil, nil,
		[]string{"(warn) NewVirtualLiquidHandler: Undefined location \"undefined_input_pref\" referenced in input preferences"},
		nil})

	lhp = default_lhproperties()
	lhp.Output_preferences = append(lhp.Tip_preferences, "undefined_output_pref")
	tests = append(tests, SimulatorTest{"passing undefined Output_preference", lhp, nil, nil,
		[]string{"(warn) NewVirtualLiquidHandler: Undefined location \"undefined_output_pref\" referenced in output preferences"},
		nil})

	lhp = default_lhproperties()
	lhp.Tipwaste_preferences = append(lhp.Tip_preferences, "undefined_tipwaste_pref")
	tests = append(tests, SimulatorTest{"passing undefined Tipwaste_preference", lhp, nil, nil,
		[]string{"(warn) NewVirtualLiquidHandler: Undefined location \"undefined_tipwaste_pref\" referenced in tipwaste preferences"},
		nil})

	lhp = default_lhproperties()
	lhp.Wash_preferences = append(lhp.Tip_preferences, "undefined_wash_pref")
	tests = append(tests, SimulatorTest{"passing undefined Wash_preference", lhp, nil, nil,
		[]string{"(warn) NewVirtualLiquidHandler: Undefined location \"undefined_wash_pref\" referenced in wash preferences"},
		nil})

	lhp = default_lhproperties()
	lhp.Waste_preferences = append(lhp.Tip_preferences, "undefined_waste_pref")
	tests = append(tests, SimulatorTest{"passing undefined Waste_preference", lhp, nil, nil,
		[]string{"(warn) NewVirtualLiquidHandler: Undefined location \"undefined_waste_pref\" referenced in waste preferences"},
		nil})

	for _, test := range tests {
		test.run(t)
	}
}

func TestMissingPrefs(t *testing.T) {
	tests := make([]SimulatorTest, 0)

	lhp := default_lhproperties()
	lhp.Tip_preferences = make([]string, 0)
	tests = append(tests, SimulatorTest{"passing missing Tip_preferences", lhp, nil, nil,
		[]string{"(warn) NewVirtualLiquidHandler: No tip preferences specified"},
		nil})

	lhp = default_lhproperties()
	lhp.Input_preferences = make([]string, 0)
	tests = append(tests, SimulatorTest{"passing missing Input_preferences", lhp, nil, nil,
		[]string{"(warn) NewVirtualLiquidHandler: No input preferences specified"},
		nil})

	lhp = default_lhproperties()
	lhp.Output_preferences = make([]string, 0)
	tests = append(tests, SimulatorTest{"passing missing Output_preferences", lhp, nil, nil,
		[]string{"(warn) NewVirtualLiquidHandler: No output preferences specified"},
		nil})

	lhp = default_lhproperties()
	lhp.Tipwaste_preferences = make([]string, 0)
	tests = append(tests, SimulatorTest{"passing missing TipWaste_preferences", lhp, nil, nil,
		[]string{"(warn) NewVirtualLiquidHandler: No tipwaste preferences specified"},
		nil})

	lhp = default_lhproperties()
	lhp.Wash_preferences = make([]string, 0)
	tests = append(tests, SimulatorTest{"passing missing Wash_preferences", lhp, nil, nil,
		[]string{"(warn) NewVirtualLiquidHandler: No wash preferences specified"},
		nil})

	lhp = default_lhproperties()
	lhp.Waste_preferences = make([]string, 0)
	tests = append(tests, SimulatorTest{"passing missing Waste_preferences", lhp, nil, nil,
		[]string{"(warn) NewVirtualLiquidHandler: No waste preferences specified"},
		nil})

	for _, test := range tests {
		test.run(t)
	}
}

func TestNewVirtualLiquidHandler_ValidProps(t *testing.T) {
	test := SimulatorTest{"Create Valid VLH", nil, nil, nil, nil, nil}
	test.run(t)
}

func TestVLH_AddPlateTo(t *testing.T) {
	tests := []SimulatorTest{
		{
			"OK", //name
			nil,  //default params
			nil,  //no setup
			[]TestRobotInstruction{
				&Initialize{},
				&AddPlateTo{"tipbox_1", default_lhtipbox("tipbox1"), "tipbox1"},
				&AddPlateTo{"tipbox_2", default_lhtipbox("tipbox2"), "tipbox2"},
				&AddPlateTo{"input_1", default_lhplate("input1"), "input1"},
				&AddPlateTo{"input_2", default_lhplate("input2"), "input2"},
				&AddPlateTo{"output_1", default_lhplate("output1"), "output1"},
				&AddPlateTo{"output_2", default_lhplate("output2"), "output2"},
				&AddPlateTo{"tipwaste", default_lhtipwaste("tipwaste"), "tipwaste"},
			},
			nil, //no errors
			nil, //no assertions
		},
		{
			"non plate type", //name
			nil,              //default params
			nil,              //no setup
			[]TestRobotInstruction{
				&Initialize{},
				&AddPlateTo{"tipbox_1", "my plate's gone stringy", "not_a_plate"},
			},
			[]string{"(err) AddPlateTo: Couldn't add object of type string to tipbox_1"},
			nil, //no assertions
		},
		{
			"location full", //name
			nil,             //default params
			nil,             //no setup
			[]TestRobotInstruction{
				&Initialize{},
				&AddPlateTo{"tipbox_1", default_lhtipbox("p0"), "p0"},
				&AddPlateTo{"tipbox_1", default_lhtipbox("p1"), "p1"},
			},
			[]string{"(err) AddPlateTo: Couldn't add tipbox \"p1\" to location \"tipbox_1\" which already contains tipbox \"p0\""},
			nil, //no assertions
		},
		{
			"tipbox on tipwaste location", //name
			nil, //default params
			nil, //no setup
			[]TestRobotInstruction{
				&Initialize{},
				&AddPlateTo{"tipwaste", default_lhtipbox("tipbox"), "tipbox"},
			},
			[]string{"(err) AddPlateTo: Slot \"tipwaste\" can't accept tipbox \"tipbox\", only tipwaste allowed"},
			nil, //no assertions
		},
		{
			"tipwaste on tipbox location", //name
			nil, //default params
			nil, //no setup
			[]TestRobotInstruction{
				&Initialize{},
				&AddPlateTo{"tipbox_1", default_lhtipwaste("tipwaste"), "tipwaste"},
			},
			[]string{"(err) AddPlateTo: Slot \"tipbox_1\" can't accept tipwaste \"tipwaste\", only tipbox allowed"},
			nil, //no assertions
		},
		{
			"unknown location", //name
			nil,                //default params
			nil,                //no setup
			[]TestRobotInstruction{
				&Initialize{},
				&AddPlateTo{"ruritania", default_lhtipbox("aTipbox"), "aTipbox"},
			},
			[]string{"(err) AddPlateTo: Cannot put tipbox \"aTipbox\" at unknown slot \"ruritania\""},
			nil, //no assertions
		},
		{
			"too big", //name
			nil,       //default params
			nil,       //no setup
			[]TestRobotInstruction{
				&Initialize{},
				&AddPlateTo{"output_1", wide_lhplate("plate1"), "plate1"},
			},
			[]string{ //errors
				"(err) AddPlateTo: Footprint of plate \"plate1\"[300mm x 85.48mm] doesn't fit slot \"output_1\"[127.76mm x 85.48mm]",
			},
			nil, //no assertions
		},
	}

	for _, test := range tests {
		test.run(t)
	}
}

func Test_SetPippetteSpeed(t *testing.T) {
	tests := []SimulatorTest{
		{
			"OK", //name
			nil,  //default params
			nil,  //no setup
			[]TestRobotInstruction{
				&Initialize{},
				&SetPipetteSpeed{0, -1, 5.},
			},
			nil, //no errors
			nil, //no assertions
		},
		{
			"too low", //name
			nil,       //default params
			nil,       //no setup
			[]TestRobotInstruction{
				&Initialize{},
				&SetPipetteSpeed{0, -1, 0.001},
			},
			[]string{
				"(warn) SetPipetteSpeed: Setting Head 0 channels 0-7 speed to 0.001 ml/min is outside allowable range [0.1 ml/min:10 ml/min]",
			},
			nil, //no assertions
		},
		{
			"too high", //name
			nil,        //default params
			nil,        //no setup
			[]TestRobotInstruction{
				&Initialize{},
				&SetPipetteSpeed{0, -1, 15.},
			},
			[]string{
				"(warn) SetPipetteSpeed: Setting Head 0 channels 0-7 speed to 15 ml/min is outside allowable range [0.1 ml/min:10 ml/min]",
			},
			nil, //no assertions
		},
		{
			"Independent", //name
			nil,           //default params
			nil,           //no setup
			[]TestRobotInstruction{
				&Initialize{},
				&SetPipetteSpeed{0, 3, 5.},
			},
			[]string{
				"(warn) SetPipetteSpeed: Head 0 is not independent, setting pipette speed for channel 3 sets all other channels as well",
			},
			nil, //no assertions
		},
	}

	for _, test := range tests {
		test.run(t)
	}
}

// ########################################################################################################################
// ########################################################## Move
// ########################################################################################################################
func testLayout() *SetupFn {
	var ret SetupFn = func(vlh *lh.VirtualLiquidHandler) {
		vlh.Initialize()
		vlh.AddPlateTo("tipbox_1", default_lhtipbox("tipbox1"), "tipbox1")
		vlh.AddPlateTo("tipbox_2", default_lhtipbox("tipbox2"), "tipbox2")
		vlh.AddPlateTo("input_1", default_lhplate("plate1"), "plate1")
		vlh.AddPlateTo("input_2", default_lhplate("plate2"), "plate2")
		vlh.AddPlateTo("output_1", default_lhplate("plate3"), "plate3")
		vlh.AddPlateTo("tipwaste", default_lhtipwaste("tipwaste"), "tipwaste")
	}
	return &ret
}

func testTroughLayout() *SetupFn {
	var ret SetupFn = func(vlh *lh.VirtualLiquidHandler) {
		vlh.Initialize()
		vlh.AddPlateTo("tipbox_1", default_lhtipbox("tipbox1"), "tipbox1")
		vlh.AddPlateTo("tipbox_2", default_lhtipbox("tipbox2"), "tipbox2")
		vlh.AddPlateTo("input_1", lhplate_trough12("trough1"), "trough1")
		vlh.AddPlateTo("input_2", default_lhplate("plate2"), "plate2")
		vlh.AddPlateTo("output_1", default_lhplate("plate3"), "plate3")
		vlh.AddPlateTo("tipwaste", default_lhtipwaste("tipwaste"), "tipwaste")
	}
	return &ret
}

func Test_Move(t *testing.T) {

	tests := []SimulatorTest{
		{
			"OK_1",
			nil,
			[]*SetupFn{
				testLayout(),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2"}, //deckposition
					[]string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},                                                 //wellcoords
					[]int{1, 1, 1, 1, 1, 1, 1, 1},                                                                            //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},                                                                //offsetZ
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //plate_type
					0, //head
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				positionAssertion(0, wtype.Coordinates{X: 204.5, Y: 4.5, Z: 62.2}),
			},
		},
		{
			"OK_2",
			nil,
			[]*SetupFn{
				testLayout(),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste"}, //deckposition
					[]string{"A1", "A1", "A1", "A1", "A1", "A1", "A1", "A1"},                                                 //wellcoords
					[]int{1, 1, 1, 1, 1, 1, 1, 1},                                                                            //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetX
					[]float64{-31.5, -22.5, -13.5, -4.5, 4.5, 13.5, 22.5, 31.5},                                              //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},                                                                //offsetZ
					[]string{"tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste"}, //plate_type
					0, //head
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				positionAssertion(0, wtype.Coordinates{X: 111., Y: 440., Z: 93.}),
			},
		},
		{
			"OK_2.5",
			nil,
			[]*SetupFn{
				testLayout(),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"tipwaste", "", "", "", "", "", "", ""}, //deckposition
					[]string{"A1", "", "", "", "", "", "", ""},       //wellcoords
					[]int{1, 1, 1, 1, 1, 1, 1, 1},                    //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},        //offsetX
					[]float64{-31.5, 0., 0., 0., 0., 0., 0., 0.},     //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},        //offsetZ
					[]string{"tipwaste", "", "", "", "", "", "", ""}, //plate_type
					0, //head
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				positionAssertion(0, wtype.Coordinates{X: 111., Y: 440., Z: 93}),
			},
		},
		{
			"OK_3",
			nil,
			[]*SetupFn{
				testLayout(),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"", "", "", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //deckposition
					[]string{"", "", "", "A1", "B1", "C1", "D1", "E1"},                               //wellcoords
					[]int{0, 0, 0, 1, 1, 1, 1, 1},                                                    //reference (first 3 should be ignored)
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                        //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                        //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},                                        //offsetZ
					[]string{"", "", "", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},           //plate_type
					0, //head
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				positionAssertion(0, wtype.Coordinates{X: 4.5, Y: -22.5, Z: 62.2}),
			},
		},
		{
			"OK_4",
			nil,
			[]*SetupFn{
				testLayout(),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"H1", "", "", "", "", "", "", ""},      //wellcoords
					[]int{1, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				positionAssertion(0, wtype.Coordinates{X: 417.1, Y: 77.0, Z: 25.7}),
			},
		},
		{
			"OK_5",
			nil,
			[]*SetupFn{
				testLayout(),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"", "", "", "", "input_1", "", "", ""}, //deckposition
					[]string{"", "", "", "", "H1", "", "", ""},      //wellcoords
					[]int{0, 0, 0, 0, 1, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetZ
					[]string{"", "", "", "", "plate", "", "", ""},   //plate_type
					0, //head
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				positionAssertion(0, wtype.Coordinates{X: 417.1, Y: 41.0, Z: 25.7}),
			},
		},
		{
			"OK_trough",
			nil,
			[]*SetupFn{
				testTroughLayout(),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"}, //deckposition
					[]string{"A1", "A1", "A1", "A1", "A1", "A1", "A1", "A1"},                                         //wellcoords
					[]int{1, 1, 1, 1, 1, 1, 1, 1},                                                                    //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},                                                        //offsetZ
					[]string{"trough", "trough", "trough", "trough", "trough", "trough", "trough", "trough"},         //plate_type
					0, //head
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				positionAssertion(0, wtype.Coordinates{X: 418.5, Y: 15.7, Z: 46.8}),
			},
		},
		{
			"unknown location",
			nil,
			[]*SetupFn{
				testLayout(),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"tipbox7", "tipbox7", "tipbox7", "tipbox7", "tipbox7", "tipbox7", "tipbox7", "tipbox7"}, //deckposition
					[]string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},                                         //wellcoords
					[]int{1, 1, 1, 1, 1, 1, 1, 1},                                                                    //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetY
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetZ
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},         //plate_type
					0, //head
				},
				&Move{
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //deckposition
					[]string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},                                                 //wellcoords
					[]int{1, 1, 1, 1, 1, 1, 1, 1},                                                                            //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetY
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetZ
					[]string{"tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste"}, //plate_type
					0, //head
				},
			},
			[]string{ //errors
				"(err) Move: Unknown location \"tipbox7\"",
				"(warn) Move: Object found at tipbox_1 was type \"tipbox\", named \"tipbox1\", not \"tipwaste\" as expected",
				"(warn) Move: Object found at tipbox_1 was type \"tipbox\", named \"tipbox1\", not \"tipwaste\" as expected",
				"(warn) Move: Object found at tipbox_1 was type \"tipbox\", named \"tipbox1\", not \"tipwaste\" as expected",
				"(warn) Move: Object found at tipbox_1 was type \"tipbox\", named \"tipbox1\", not \"tipwaste\" as expected",
				"(warn) Move: Object found at tipbox_1 was type \"tipbox\", named \"tipbox1\", not \"tipwaste\" as expected",
				"(warn) Move: Object found at tipbox_1 was type \"tipbox\", named \"tipbox1\", not \"tipwaste\" as expected",
				"(warn) Move: Object found at tipbox_1 was type \"tipbox\", named \"tipbox1\", not \"tipwaste\" as expected",
				"(warn) Move: Object found at tipbox_1 was type \"tipbox\", named \"tipbox1\", not \"tipwaste\" as expected",
			},
			nil, //assertions
		},
		{
			"unknown head",
			nil,
			[]*SetupFn{
				testLayout(),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //deckposition
					[]string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},                                                 //wellcoords
					[]int{1, 1, 1, 1, 1, 1, 1, 1},                                                                            //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetY
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetZ
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //plate_type
					1, //head
				},
				&Move{
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //deckposition
					[]string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},                                                 //wellcoords
					[]int{1, 1, 1, 1, 1, 1, 1, 1},                                                                            //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetY
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetZ
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //plate_type
					-1, //head
				},
			},
			[]string{ //errors
				"(err) Move: head assembly 0: unknown head 1",
				"(err) Move: head assembly 0: unknown head -1",
			},
			nil, //assertions
		},
		{
			"invalid wellcoords",
			nil,
			[]*SetupFn{
				testLayout(),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //deckposition
					[]string{"B1", "C1", "D1", "E1", "F1", "G1", "H1", "I1"},                                                 //wellcoords
					[]int{1, 1, 1, 1, 1, 1, 1, 1},                                                                            //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetY
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetZ
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //plate_type
					0, //head
				},
				&Move{
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //deckposition
					[]string{"B1", "C1", "D1", "E1", "F1", "G1", "H1", "not_a_well"},                                         //wellcoords
					[]int{1, 1, 1, 1, 1, 1, 1, 1},                                                                            //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetY
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetZ
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //plate_type
					0, //head
				},
			},
			[]string{ //errors
				"(err) Move: Request for well I1 in object \"tipbox1\" at \"tipbox_1\" which is of size [8x12]",
				"(err) Move: couldn't parse well coordinates \"not_a_well\"",
			},
			nil, //assertions
		},
		{
			"Invalid reference",
			nil,
			[]*SetupFn{
				testLayout(),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //deckposition
					[]string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},                                                 //wellcoords
					[]int{-1, -1, -1, -1, -1, -1, -1, -1},                                                                    //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetY
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetZ
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //plate_type
					0, //head
				},
				&Move{
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //deckposition
					[]string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},                                                 //wellcoords
					[]int{3, 3, 3, 3, 3, 3, 3, 3},                                                                            //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetY
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetZ
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //plate_type
					0, //head
				},
			},
			[]string{ //errors
				"(err) Move: Invalid reference -1",
				"(err) Move: Invalid reference 3",
			},
			nil, //assertions
		},
		{
			"Inconsistent references",
			nil,
			[]*SetupFn{
				testLayout(),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"}, //deckposition
					[]string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},                                         //wellcoords
					[]int{0, 0, 0, 0, 1, 1, 1, 1},                                                                    //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetY
					[]float64{-5., -5., -5., -5., -5., -5., -5., -5.},                                                //offsetZ
					[]string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"},                 //plate_type
					0, //head
				},
			},
			[]string{ //errors
				"(err) Move: head 0 channels 0-7 to A1-H1@plate at position input_1: requires moving channels 4-7 relative to non-independent head",
			},
			nil, //assertions
		},
		{
			"offsets differ",
			nil,
			[]*SetupFn{
				testLayout(),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //deckposition
					[]string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},                                                 //wellcoords
					[]int{1, 1, 1, 1, 1, 1, 1, 1},                                                                            //reference
					[]float64{0., 0., 0., 0., 0., 3., 0., 0.},                                                                //offsetX
					[]float64{0., 0., 0., 1., 0., 0., 0., 0.},                                                                //offsetY
					[]float64{0., 0., 0., 0., 1., 0., 0., 0.},                                                                //offsetZ
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //plate_type
					0, //head
				},
			},
			[]string{ //errors
				"(err) Move: head 0 channels 0-7 to A1-H1@tipbox at position tipbox_1: requires moving channels 3-5 relative to non-independent head",
			},
			nil, //assertions
		},
		{
			"layout mismatch",
			nil,
			[]*SetupFn{
				testLayout(),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //deckposition
					[]string{"A1", "B2", "C1", "D2", "E1", "F2", "G1", "H2"},                                                 //wellcoords
					[]int{1, 1, 1, 1, 1, 1, 1, 1},                                                                            //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetY
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetZ
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //plate_type
					0, //head
				},
			},
			[]string{ //errors
				"(err) Move: head 0 channels 0-7 to A1,B2,C1,D2,E1,F2,G1,H2@tipbox at position tipbox_1: requires moving channels 1,3,5,7 relative to non-independent head",
			},
			nil, //assertions
		},
	}

	for _, test := range tests {
		test.run(t)
	}
}

func Test_MoveConstraints(t *testing.T) {

	tests := []SimulatorTest{
		{
			"constrained heads",
			multihead_lhproperties(),
			[]*SetupFn{
				testLayout(),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2"}, //deckposition
					[]string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},                                                 //wellcoords
					[]int{1, 1, 1, 1, 1, 1, 1, 1},                                                                            //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},                                                                //offsetZ
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //plate_type
					0, //head
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				positionAssertion(0, wtype.Coordinates{X: 132.5, Y: 4.5, Z: 62.2}),
				positionAssertion(1, wtype.Coordinates{X: 132.5, Y: 22.5, Z: 62.2}),
			},
		},
		{
			"outside limits",
			multihead_lhproperties(),
			[]*SetupFn{
				testLayout(),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2", "tipbox_2"}, //deckposition
					[]string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},                                                 //wellcoords
					[]int{1, 1, 1, 1, 1, 1, 1, 1},                                                                            //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},                                                                //offsetZ
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //plate_type
					1, //head
				},
			},
			[]string{ //errors
				"(err) Move: head 1 channels 0-7 to A1-H1@tipbox at position tipbox_2: movement limits prevent moving into position",
			},
			[]*AssertionFn{ //assertions
				positionAssertion(0, wtype.Coordinates{X: 132.5, Y: -13.5, Z: 62.2}),
				positionAssertion(1, wtype.Coordinates{X: 132.5, Y: 4.5, Z: 62.2}),
			},
		},
		{
			"tips on other head",
			multihead_lhproperties(),
			[]*SetupFn{
				testLayout(),
				preloadAdaptorTips(1, "tipbox_1", []int{0}),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //deckposition
					[]string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"},                                         //wellcoords
					[]int{1, 1, 1, 1, 1, 1, 1, 1},                                                                            //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},                                                                //offsetZ
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //plate_type
					0, //head
				},
				&LoadTips{
					[]int{0, 1, 2, 3, 4, 5, 6, 7}, //channels
					0, //head
					8, //multi
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //tipbox
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //location
					[]string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"},                                         //well
				},
			},
			[]string{ //errors
				"(err) LoadTips: from A12-H12@tipbox1 at position \"tipbox_1\" to head 0 channels 0-7: tip already loaded on head 1 channel 0",
			},
			[]*AssertionFn{ //assertions
			},
		},
	}

	for _, test := range tests {
		test.run(t)
	}
}

// ########################################################################################################################
// ########################################################## Tip Loading/Unloading
// ########################################################################################################################

func TestLoadTips(t *testing.T) {

	mtp := moveToParams{
		8,                     //Multi           int
		0,                     //Head            int
		1,                     //Reference       int
		"tipbox_1",            //Deckposition    string
		"tipbox",              //Platetype       string
		[]float64{0., 0., 5.}, //Offset          wtype.Coords
		12, //Cols            int
		8,  //Rows            int
	}
	misaligned_mtp := moveToParams{
		8,                     //Multi           int
		0,                     //Head            int
		1,                     //Reference       int
		"tipbox_1",            //Deckposition    string
		"tipbox",              //Platetype       string
		[]float64{0., 2., 5.}, //Offset          wtype.Coords
		12, //Cols            int
		8,  //Rows            int
	}

	tests := []SimulatorTest{
		{
			"OK - single tip",
			nil,
			[]*SetupFn{
				testLayout(),
				moveTo(7, 11, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0}, //channels
					0,        //head
					1,        //multi
					[]string{"tipbox", "", "", "", "", "", "", ""},   //tipbox
					[]string{"tipbox_1", "", "", "", "", "", "", ""}, //location
					[]string{"H12", "", "", "", "", "", "", ""},      //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{"H12"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			"OK - single tip (alt)",
			nil,
			[]*SetupFn{
				testLayout(),
				moveTo(-7, 0, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{7}, //channels
					0,        //head
					1,        //multi
					[]string{"", "", "", "", "", "", "", "tipbox"},   //tipbox
					[]string{"", "", "", "", "", "", "", "tipbox_1"}, //location
					[]string{"", "", "", "", "", "", "", "A1"},       //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{"A1"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{7, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			"OK - single tip above space",
			nil,
			[]*SetupFn{
				testLayout(),
				removeTipboxTips("tipbox_1", []string{"H12"}),
				moveTo(6, 11, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0}, //channels
					0,        //head
					1,        //multi
					[]string{"tipbox", "", "", "", "", "", "", ""},   //tipbox
					[]string{"tipbox_1", "", "", "", "", "", "", ""}, //location
					[]string{"G12", "", "", "", "", "", "", ""},      //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{"H12", "G12"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			"OK - single tip below space (alt)",
			nil,
			[]*SetupFn{
				testLayout(),
				removeTipboxTips("tipbox_1", []string{"A1"}),
				moveTo(-6, 0, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{7}, //channels
					0,        //head
					1,        //multi
					[]string{"", "", "", "", "", "", "", "tipbox"},   //tipbox
					[]string{"", "", "", "", "", "", "", "tipbox_1"}, //location
					[]string{"", "", "", "", "", "", "", "B1"},       //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{"A1", "B1"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{7, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			"OK - 3 tips at once",
			nil,
			[]*SetupFn{
				testLayout(),
				moveTo(5, 11, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0, 1, 2}, //channels
					0,              //head
					3,              //multi
					[]string{"tipbox", "tipbox", "tipbox", "", "", "", "", ""},       //tipbox
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "", "", "", "", ""}, //location
					[]string{"F12", "G12", "H12", "", "", "", "", ""},                //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
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
			"OK - 3 tips at once (alt)",
			nil,
			[]*SetupFn{
				testLayout(),
				moveTo(-5, 0, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{5, 6, 7}, //channels
					0,              //head
					3,              //multi
					[]string{"", "", "", "", "", "tipbox", "tipbox", "tipbox"},       //tipbox
					[]string{"", "", "", "", "", "tipbox_1", "tipbox_1", "tipbox_1"}, //location
					[]string{"", "", "", "", "", "A1", "B1", "C1"},                   //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
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
			"OK - 3 tips (independent)",
			independent_lhproperties(),
			[]*SetupFn{
				testLayout(),
				moveTo(0, 0, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0, 4, 7}, //channels
					0,              //head
					3,              //multi
					[]string{"tipbox", "", "", "", "tipbox", "", "", "tipbox"},       //tipbox
					[]string{"tipbox_1", "", "", "", "tipbox_1", "", "", "tipbox_1"}, //location
					[]string{"A1", "", "", "", "E1", "", "", "H1"},                   //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
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
			"OK - 8 tips at once",
			nil,
			[]*SetupFn{
				testLayout(),
				moveTo(0, 11, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0, 1, 2, 3, 4, 5, 6, 7}, //channels
					0, //head
					8, //multi
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //tipbox
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //location
					[]string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"},                                         //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
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
			"OK - 2 groups of 4",
			nil,
			[]*SetupFn{
				testLayout(),
				moveTo(4, 0, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0, 1, 2, 3}, //channels
					0,                 //head
					4,                 //multi
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "", "", "", ""},         //tipbox
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "", "", "", ""}, //location
					[]string{"E1", "F1", "G1", "H1", "", "", "", ""},                         //well
				},
				&Move{
					[]string{"", "", "", "", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //deckposition
					[]string{"", "", "", "", "A1", "B1", "C1", "D1"},                         //wellcoords
					[]int{1, 1, 1, 1, 1, 1, 1, 1},                                            //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},                                //offsetZ
					[]string{"", "", "", "", "tipbox", "tipbox", "tipbox", "tipbox"},         //plate_type
					0, //head
				},
				&LoadTips{
					[]int{4, 5, 6, 7}, //channels
					0,                 //head
					4,                 //multi
					[]string{"", "", "", "", "tipbox", "tipbox", "tipbox", "tipbox"},         //tipbox
					[]string{"", "", "", "", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //location
					[]string{"", "", "", "", "A1", "B1", "C1", "D1"},                         //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
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
			"unknown channel 8",
			nil,
			[]*SetupFn{
				testLayout(),
				moveTo(0, 0, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{8}, //channels
					0,        //head
					1,        //multi
					[]string{"", "", "", "", "", "", "", "tipbox"},   //tipbox
					[]string{"", "", "", "", "", "", "", "tipbox_1"}, //location
					[]string{"", "", "", "", "", "", "", "H12"},      //well
				},
			},
			[]string{ //errors
				"(err) LoadTips: Unknown channel \"8\"",
			},
			nil, //assertions
		},
		{
			"unknown channel -1",
			nil,
			[]*SetupFn{
				testLayout(),
				moveTo(0, 0, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{-1}, //channels
					0,         //head
					1,         //multi
					[]string{"", "", "", "", "", "", "", "tipbox"},   //tipbox
					[]string{"", "", "", "", "", "", "", "tipbox_1"}, //location
					[]string{"", "", "", "", "", "", "", "H12"},      //well
				},
			},
			[]string{ //errors
				"(err) LoadTips: Unknown channel \"-1\"",
			},
			nil, //assertions
		},
		{
			"duplicate channels",
			nil,
			[]*SetupFn{
				testLayout(),
				moveTo(0, 11, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0, 1, 2, 3, 4, 5, 6, 3}, //channels
					0, //head
					8, //multi
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //tipbox
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //location
					[]string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"},                                         //well
				},
			},
			[]string{ //errors
				"(err) LoadTips: Channel3 appears more than once",
			},
			nil, //assertions
		},
		{
			"unknown head",
			nil,
			[]*SetupFn{
				testLayout(),
				moveTo(7, 11, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0}, //channels
					1,        //head
					1,        //multi
					[]string{"tipbox", "", "", "", "", "", "", ""},   //tipbox
					[]string{"tipbox_1", "", "", "", "", "", "", ""}, //location
					[]string{"H12", "", "", "", "", "", "", ""},      //well
				},
			},
			[]string{ //errors
				"(err) LoadTips: head assembly 0: unknown head 1",
			},
			nil, //assertions
		},
		{
			"unknown head -1",
			nil,
			[]*SetupFn{
				testLayout(),
				moveTo(7, 11, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0}, //channels
					-1,       //head
					1,        //multi
					[]string{"tipbox", "", "", "", "", "", "", ""},   //tipbox
					[]string{"tipbox_1", "", "", "", "", "", "", ""}, //location
					[]string{"H12", "", "", "", "", "", "", ""},      //well
				},
			},
			[]string{ //errors
				"(err) LoadTips: head assembly 0: unknown head -1",
			},
			nil, //assertions
		},
		{
			"OK - argument expansion",
			nil,
			[]*SetupFn{
				testLayout(),
				moveTo(7, 11, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0},             //channels
					0,                    //head
					1,                    //multi
					[]string{"tipbox"},   //tipbox
					[]string{"tipbox_1"}, //location
					[]string{"H12"},      //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{"H12"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			"mismatching multi",
			nil,
			[]*SetupFn{
				testLayout(),
				moveTo(7, 11, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0}, //channels
					0,        //head
					4,        //multi
					[]string{"tipbox", "", "", ""},   //tipbox
					[]string{"tipbox_1", "", "", ""}, //location
					[]string{"H12", "", "", ""},      //well
				},
			},
			[]string{ //errors
				"(err) LoadTips: from H12@tipbox1 at position \"tipbox_1\" to head 0 channel 0 : multi should equal 1, not 4",
			},
			nil, //assertions
		},
		{
			"tip missing",
			nil,
			[]*SetupFn{
				testLayout(),
				removeTipboxTips("tipbox_1", []string{"H12"}),
				moveTo(7, 11, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0}, //channels
					0,        //head
					1,        //multi
					[]string{"tipbox", "", "", "", "", "", "", ""},   //tipbox
					[]string{"tipbox_1", "", "", "", "", "", "", ""}, //location
					[]string{"H12", "", "", "", "", "", "", ""},      //well
				},
			},
			[]string{ //errors
				"(err) LoadTips: from H12@tipbox1 at position \"tipbox_1\" to head 0 channel 0 : no tip at H12",
			},
			nil, //assertions
		},
		{
			"8 tips missing",
			nil,
			[]*SetupFn{
				testLayout(),
				removeTipboxTips("tipbox_1", []string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"}),
				moveTo(0, 0, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0, 1, 2, 3, 4, 5, 6, 7}, //channels
					0, //head
					8, //multi
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //tipbox
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //location
					[]string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"},                                         //well
				},
			},
			[]string{ //errors
				"(err) LoadTips: from A12-H12@tipbox1 at position \"tipbox_1\" to head 0 channels 0-7 : no tips at A12-H12",
			},
			nil, //assertions
		},
		{
			"tip already loaded",
			nil,
			[]*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
				moveTo(7, 11, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0}, //channels
					0,        //head
					1,        //multi
					[]string{"tipbox", "", "", "", "", "", "", ""},   //tipbox
					[]string{"tipbox_1", "", "", "", "", "", "", ""}, //location
					[]string{"H12", "", "", "", "", "", "", ""},      //well
				},
			},
			[]string{ //errors
				"(err) LoadTips: from H12@tipbox1 at position \"tipbox_1\" to head 0 channel 0: tip already loaded on head 0 channel 0",
			},
			nil, //assertions
		},
		{
			"tips already loaded",
			nil,
			[]*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}),
				moveTo(0, 11, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0, 1, 2, 3, 4, 5, 6, 7}, //channels
					0, //head
					8, //multi
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //tipbox
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //location
					[]string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"},                                         //well
				},
			},
			[]string{ //errors
				"(err) LoadTips: from A12-H12@tipbox1 at position \"tipbox_1\" to head 0 channels 0-7: tips already loaded on head 0 channels 0-7",
			},
			nil, //assertions
		},
		{
			"extra tip in the way",
			nil,
			[]*SetupFn{
				testLayout(),
				moveTo(6, 11, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0}, //channels
					0,        //head
					1,        //multi
					[]string{"tipbox", "", "", "", "", "", "", ""},   //tipbox
					[]string{"tipbox_1", "", "", "", "", "", "", ""}, //location
					[]string{"G12", "", "", "", "", "", "", ""},      //well
				},
			},
			[]string{ //errors
				"(err) LoadTips: from G12@tipbox1 at position \"tipbox_1\" to head 0 channel 0 : channel 1 collides with tip \"H12@tipbox1\" (head not independent)",
			},
			nil, //assertions
		},
		{
			"not aligned to move",
			nil,
			[]*SetupFn{
				testLayout(),
				moveTo(5, 11, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0, 1, 2}, //channels
					0,              //head
					3,              //multi
					[]string{"tipbox", "tipbox", "tipbox", "", "", "", "", ""},       //tipbox
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "", "", "", "", ""}, //location
					[]string{"E12", "G12", "H12", "", "", "", "", ""},                //well
				},
			},
			[]string{ //errors
				"(err) LoadTips: from E12,G12-H12@tipbox1 at position \"tipbox_1\" to head 0 channels 0-2 : channel 0 is misaligned with tip at E12 by 9mm",
			},
			nil, //assertions
		},
		{
			"multiple not aligned to move",
			nil,
			[]*SetupFn{
				testLayout(),
				moveTo(5, 11, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0, 1, 2}, //channels
					0,              //head
					3,              //multi
					[]string{"tipbox", "tipbox", "tipbox", "", "", "", "", ""},       //tipbox
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "", "", "", "", ""}, //location
					[]string{"G12", "F12", "H12", "", "", "", "", ""},                //well
				},
			},
			[]string{ //errors
				"(err) LoadTips: from G12,F12,H12@tipbox1 at position \"tipbox_1\" to head 0 channels 0-2 : channels 0-1 are misaligned with tips at G12,F12 by 9,9mm respectively",
			},
			nil, //assertions
		},
		{
			"misalignment single",
			nil,
			[]*SetupFn{
				testLayout(),
				moveTo(7, 11, misaligned_mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0}, //channels
					0,        //head
					1,        //multi
					[]string{"tipbox", "", "", "", "", "", "", ""},   //tipbox
					[]string{"tipbox_1", "", "", "", "", "", "", ""}, //location
					[]string{"H12", "", "", "", "", "", "", ""},      //well
				},
			},
			[]string{ //errors
				"(err) LoadTips: from H12@tipbox1 at position \"tipbox_1\" to head 0 channel 0 : channel 0 is misaligned with tip at H12 by 2mm",
			},
			nil, //assertions
		},
		{
			"misalignment multi",
			nil,
			[]*SetupFn{
				testLayout(),
				moveTo(5, 11, misaligned_mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0, 1, 2}, //channels
					0,              //head
					3,              //multi
					[]string{"tipbox", "tipbox", "tipbox", "", "", "", "", ""},       //tipbox
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "", "", "", "", ""}, //location
					[]string{"F12", "G12", "H12", "", "", "", "", ""},                //well
				},
			},
			[]string{ //errors
				"(err) LoadTips: from F12-H12@tipbox1 at position \"tipbox_1\" to head 0 channels 0-2 : channels 0-2 are misaligned with tips at F12-H12 by 2,2,2mm respectively",
			},
			nil, //assertions
		},
	}

	for _, test := range tests {
		test.run(t)
	}
}

func TestLoadTipsOverride(t *testing.T) {

	mtp := moveToParams{
		8,                     //Multi           int
		0,                     //Head            int
		1,                     //Reference       int
		"tipbox_1",            //Deckposition    string
		"tipbox",              //Platetype       string
		[]float64{0., 0., 5.}, //Offset          wtype.Coords
		12, //Cols            int
		8,  //Rows            int
	}

	propsLTR := default_lhproperties()
	propsLTR.Heads[0].TipLoading = wtype.TipLoadingBehaviour{
		OverrideLoadTipsCommand:    true,
		AutoRefillTipboxes:         true,
		LoadingOrder:               wtype.ColumnWise,
		VerticalLoadingDirection:   wtype.BottomToTop,
		HorizontalLoadingDirection: wtype.LeftToRight,
		ChunkingBehaviour:          wtype.ReverseSequentialTipLoading,
	}

	propsRTL := default_lhproperties()
	propsRTL.Heads[0].TipLoading = wtype.TipLoadingBehaviour{
		OverrideLoadTipsCommand:    true,
		AutoRefillTipboxes:         true,
		LoadingOrder:               wtype.ColumnWise,
		VerticalLoadingDirection:   wtype.BottomToTop,
		HorizontalLoadingDirection: wtype.RightToLeft,
		ChunkingBehaviour:          wtype.ReverseSequentialTipLoading,
	}

	tests := []SimulatorTest{
		{
			"OK - single tip LTR override (A1 -> H1)",
			propsLTR,
			[]*SetupFn{
				testLayout(),
				moveTo(0, 0, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0}, //channels
					0,        //head
					1,        //multi
					[]string{"tipbox", "", "", "", "", "", "", ""},   //tipbox
					[]string{"tipbox_1", "", "", "", "", "", "", ""}, //location
					[]string{"A1", "", "", "", "", "", "", ""},       //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{"H1"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			"OK - single tip RTL override (A12 -> H12)",
			propsRTL,
			[]*SetupFn{
				testLayout(),
				moveTo(0, 11, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0}, //channels
					0,        //head
					1,        //multi
					[]string{"tipbox", "", "", "", "", "", "", ""},   //tipbox
					[]string{"tipbox_1", "", "", "", "", "", "", ""}, //location
					[]string{"A12", "", "", "", "", "", "", ""},      //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{"H12"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			"OK - single tip LTR override (A1 -> D1)",
			propsLTR,
			[]*SetupFn{
				testLayout(),
				removeTipboxTips("tipbox_1", []string{"E1", "F1", "G1", "H1"}),
				moveTo(0, 0, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0}, //channels
					0,        //head
					1,        //multi
					[]string{"tipbox", "", "", "", "", "", "", ""},   //tipbox
					[]string{"tipbox_1", "", "", "", "", "", "", ""}, //location
					[]string{"A1", "", "", "", "", "", "", ""},       //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{"D1", "E1", "F1", "G1", "H1"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			"OK - single tip RTL override (A12 -> D12)",
			propsRTL,
			[]*SetupFn{
				testLayout(),
				removeTipboxTips("tipbox_1", []string{"E12", "F12", "G12", "H12"}),
				moveTo(0, 11, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0}, //channels
					0,        //head
					1,        //multi
					[]string{"tipbox", "", "", "", "", "", "", ""},   //tipbox
					[]string{"tipbox_1", "", "", "", "", "", "", ""}, //location
					[]string{"A12", "", "", "", "", "", "", ""},      //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{"D12", "E12", "F12", "G12", "H12"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			"8 tips LTR",
			propsLTR,
			[]*SetupFn{
				testLayout(),
				moveTo(0, 11, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0, 1, 2, 3, 4, 5, 6, 7}, //channels
					0, //head
					8, //multi
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //tipbox
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //location
					[]string{"A11", "B11", "C11", "D11", "E11", "F11", "G11", "H11"},                                         //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}, {1, "", 0}, {2, "", 0}, {3, "", 0}, {4, "", 0}, {5, "", 0}, {6, "", 0}, {7, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			"8 tips RTL",
			propsRTL,
			[]*SetupFn{
				testLayout(),
				moveTo(0, 11, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0, 1, 2, 3, 4, 5, 6, 7}, //channels
					0, //head
					8, //multi
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //tipbox
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //location
					[]string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},                                                 //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}, {1, "", 0}, {2, "", 0}, {3, "", 0}, {4, "", 0}, {5, "", 0}, {6, "", 0}, {7, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			"8 tips LTR override with split",
			propsLTR,
			[]*SetupFn{
				testLayout(),
				removeTipboxTips("tipbox_1", []string{"E1", "F1", "G1", "H1"}),
				moveTo(0, 11, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0, 1, 2, 3, 4, 5, 6, 7}, //channels
					0, //head
					8, //multi
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //tipbox
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //location
					[]string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},                                                 //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{"A1", "B1", "C1", "D1", "E2", "F2", "G2", "H2", "E1", "F1", "G1", "H1"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}, {1, "", 0}, {2, "", 0}, {3, "", 0}, {4, "", 0}, {5, "", 0}, {6, "", 0}, {7, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			"8 tips RTL override with split",
			propsRTL,
			[]*SetupFn{
				testLayout(),
				removeTipboxTips("tipbox_1", []string{"E12", "F12", "G12", "H12"}),
				moveTo(0, 0, mtp),
			},
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0, 1, 2, 3, 4, 5, 6, 7}, //channels
					0, //head
					8, //multi
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //tipbox
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //location
					[]string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},                                                 //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{"A12", "B12", "C12", "D12", "E11", "F11", "G11", "H11", "E12", "F12", "G12", "H12"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}, {1, "", 0}, {2, "", 0}, {3, "", 0}, {4, "", 0}, {5, "", 0}, {6, "", 0}, {7, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			"8 tips LTR override with boxchange",
			propsLTR,
			[]*SetupFn{
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
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0, 1, 2, 3, 4, 5, 6, 7}, //channels
					0, //head
					8, //multi
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //tipbox
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //location
					[]string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},                                                 //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "", 0}, {1, "", 0}, {2, "", 0}, {3, "", 0}, {4, "", 0}, {5, "", 0}, {6, "", 0}, {7, "", 0}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			"8 tips LTR override without boxchange",
			propsLTR,
			[]*SetupFn{
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
			[]TestRobotInstruction{
				&LoadTips{
					[]int{0, 1, 2, 3, 4, 5, 6, 7}, //channels
					0, //head
					8, //multi
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //tipbox
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //location
					[]string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},                                                 //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
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
	}

	for _, test := range tests {
		test.run(t)
	}
}

func Test_UnloadTips(t *testing.T) {

	tests := []SimulatorTest{
		{
			"OK - single tip",
			nil,
			[]*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"tipwaste", "", "", "", "", "", "", ""}, //deckposition
					[]string{"A1", "", "", "", "", "", "", ""},       //wellcoords
					[]int{1, 1, 1, 1, 1, 1, 1, 1},                    //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},        //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},        //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},        //offsetZ
					[]string{"tipwaste", "", "", "", "", "", "", ""}, //plate_type
					0, //head
				},
				&UnloadTips{
					[]int{0}, //channels
					0,        //head
					1,        //multi
					[]string{"tipwaste", "", "", "", "", "", "", ""}, //tipbox
					[]string{"tipwaste", "", "", "", "", "", "", ""}, //location
					[]string{"A1", "", "", "", "", "", "", ""},       //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{}),
				tipwasteAssertion("tipwaste", 1),
			},
		},
		{
			"OK - 8 tips",
			nil,
			[]*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste"}, //deckposition
					[]string{"A1", "A1", "A1", "A1", "A1", "A1", "A1", "A1"},                                                 //wellcoords
					[]int{1, 1, 1, 1, 1, 1, 1, 1},                                                                            //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetX
					[]float64{-31.5, -22.5, -13.5, -4.5, 4.5, 13.5, 22.5, 31.5},                                              //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},                                                                //offsetZ
					[]string{"tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste"}, //plate_type
					0, //head
				},
				&UnloadTips{
					[]int{0, 1, 2, 3, 4, 5, 6, 7}, //channels
					0, //head
					8, //multi
					[]string{"tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste"}, //tipbox
					[]string{"tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste"}, //location
					[]string{"A1", "A1", "A1", "A1", "A1", "A1", "A1", "A1"},                                                 //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{}),
				tipwasteAssertion("tipwaste", 8),
			},
		},
		{
			"OK - 8 tips back to a tipbox",
			nil,
			[]*SetupFn{
				testLayout(),
				removeTipboxTips("tipbox_1", []string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"}),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //deckposition
					[]string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"},                                         //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                                                                            //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},                                                                //offsetZ
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //plate_type
					0, //head
				},
				&UnloadTips{
					[]int{0, 1, 2, 3, 4, 5, 6, 7}, //channels
					0, //head
					8, //multi
					[]string{"tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox", "tipbox"},                 //tipbox
					[]string{"tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1", "tipbox_1"}, //location
					[]string{"A12", "B12", "C12", "D12", "E12", "F12", "G12", "H12"},                                         //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			"OK - independent tips",
			independent_lhproperties(),
			[]*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}),
			},
			[]TestRobotInstruction{
				&UnloadTips{
					[]int{0, 2, 4, 6}, //channels
					0,                 //head
					4,                 //multi
					[]string{"tipwaste", "", "tipwaste", "", "tipwaste", "", "tipwaste", ""}, //tipbox
					[]string{"tipwaste", "", "tipwaste", "", "tipwaste", "", "tipwaste", ""}, //location
					[]string{"A1", "", "A1", "", "A1", "", "A1", ""},                         //well
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
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
			"can only unload all tips",
			nil,
			[]*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste"}, //deckposition
					[]string{"A1", "A1", "A1", "A1", "A1", "A1", "A1", "A1"},                                                 //wellcoords
					[]int{1, 1, 1, 1, 1, 1, 1, 1},                                                                            //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                                //offsetX
					[]float64{-31.5, -22.5, -13.5, -4.5, 4.5, 13.5, 22.5, 31.5},                                              //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},                                                                //offsetZ
					[]string{"tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste", "tipwaste"}, //plate_type
					0, //head
				},
				&UnloadTips{
					[]int{0, 2, 4, 6}, //channels
					0,                 //head
					4,                 //multi
					[]string{"tipwaste", "", "tipwaste", "", "tipwaste", "", "tipwaste", ""}, //tipbox
					[]string{"tipwaste", "", "tipwaste", "", "tipwaste", "", "tipwaste", ""}, //location
					[]string{"A1", "", "A1", "", "A1", "", "A1", ""},                         //well
				},
			},
			[]string{ //errors
				"(err) UnloadTips: Cannot unload tips from head0 channels 0,2,4,6 without unloading tips from channels 1,3,5,7 (head isn't independent)",
			},
			nil, //assertions
		},
		{
			"can't unload to a plate",
			nil,
			[]*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"A12", "", "", "", "", "", "", ""},     //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},       //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
				&UnloadTips{
					[]int{0}, //channels
					0,        //head
					1,        //multi
					[]string{"plate", "", "", "", "", "", "", ""},   //tipbox
					[]string{"input_1", "", "", "", "", "", "", ""}, //location
					[]string{"A1", "", "", "", "", "", "", ""},      //well
				},
			},
			[]string{ //errors
				"(err) UnloadTips: Cannot unload tips to plate \"plate1\" at location input_1",
			},
			nil,
		},
		{
			"wrong well",
			nil,
			[]*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"tipwaste", "", "", "", "", "", "", ""}, //deckposition
					[]string{"A1", "", "", "", "", "", "", ""},       //wellcoords
					[]int{1, 1, 1, 1, 1, 1, 1, 1},                    //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},        //offsetX
					[]float64{-31.5, 0., 0., 0., 0., 0., 0., 0.},     //offsetY
					[]float64{1., 0., 0., 0., 0., 0., 0., 0.},        //offsetZ
					[]string{"tipwaste", "", "", "", "", "", "", ""}, //plate_type
					0, //head
				},
				&UnloadTips{
					[]int{0}, //channels
					0,        //head
					1,        //multi
					[]string{"tipwaste", "", "", "", "", "", "", ""}, //tipbox
					[]string{"tipwaste", "", "", "", "", "", "", ""}, //location
					[]string{"B1", "", "", "", "", "", "", ""},       //well
				},
			},
			[]string{ //errors
				"(err) UnloadTips: Cannot unload to address B1 in tipwaste \"tipwaste\" size [1x1]",
			},
			nil,
		},
	}

	for _, test := range tests {
		test.run(t)
	}
}

func Test_Aspirate(t *testing.T) {

	tests := []SimulatorTest{
		{
			"OK - single channel",
			nil,
			[]*SetupFn{
				testLayout(),
				prefillWells("input_1", []string{"A1"}, "water", 200.),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"A1", "", "", "", "", "", "", ""},      //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},       //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
				&Aspirate{
					[]float64{100., 0., 0., 0., 0., 0., 0., 0.},                    //volume     []float64
					[]bool{false, false, false, false, false, false, false, false}, //overstroke []bool
					0, //head       int
					1, //multi      int
					[]string{"plate", "", "", "", "", "", "", ""},                  //platetype  []string
					[]string{"water", "", "", "", "", "", "", ""},                  //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //llf        []bool
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{{0, "water", 100}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			"OK - 8 channel",
			nil,
			[]*SetupFn{
				testLayout(),
				prefillWells("input_1", []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}, "water", 200.),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"}, //deckposition
					[]string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},                                         //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                                                                    //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},                                                        //offsetZ
					[]string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"},                 //plate_type
					0, //head
				},
				&Aspirate{
					[]float64{100., 100., 100., 100., 100., 100., 100., 100.},      //volume     []float64
					[]bool{false, false, false, false, false, false, false, false}, //overstroke []bool
					0, //head       int
					8, //multi      int
					[]string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"}, //platetype  []string
					[]string{"water", "water", "water", "water", "water", "water", "water", "water"}, //what       []string
					[]bool{false, false, false, false, false, false, false, false},                   //llf        []bool
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
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
			"OK - 8 channel trough",
			nil,
			[]*SetupFn{
				testTroughLayout(),
				prefillWells("input_1", []string{"A1"}, "water", 10000.),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"}, //deckposition
					[]string{"A1", "A1", "A1", "A1", "A1", "A1", "A1", "A1"},                                         //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                                                                    //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},                                                        //offsetZ
					[]string{"trough", "trough", "trough", "trough", "trough", "trough", "trough", "trough"},         //plate_type
					0, //head
				},
				&Aspirate{
					[]float64{100., 100., 100., 100., 100., 100., 100., 100.},      //volume     []float64
					[]bool{false, false, false, false, false, false, false, false}, //overstroke []bool
					0, //head       int
					8, //multi      int
					[]string{"trough", "trough", "trough", "trough", "trough", "trough", "trough", "trough"}, //platetype  []string
					[]string{"water", "water", "water", "water", "water", "water", "water", "water"},         //what       []string
					[]bool{false, false, false, false, false, false, false, false},                           //llf        []bool
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
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
			"Fail - take too much from trough",
			nil,
			[]*SetupFn{
				testTroughLayout(),
				prefillWells("input_1", []string{"A1"}, "water", 5400.),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"}, //deckposition
					[]string{"A1", "A1", "A1", "A1", "A1", "A1", "A1", "A1"},                                         //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                                                                    //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},                                                        //offsetZ
					[]string{"trough", "trough", "trough", "trough", "trough", "trough", "trough", "trough"},         //plate_type
					0, //head
				},
				&Aspirate{
					[]float64{100., 100., 100., 100., 100., 100., 100., 100.},      //volume     []float64
					[]bool{false, false, false, false, false, false, false, false}, //overstroke []bool
					0, //head       int
					8, //multi      int
					[]string{"trough", "trough", "trough", "trough", "trough", "trough", "trough", "trough"}, //platetype  []string
					[]string{"water", "water", "water", "water", "water", "water", "water", "water"},         //what       []string
					[]bool{false, false, false, false, false, false, false, false},                           //llf        []bool
				},
			},
			[]string{ //errors
				"(warn) Aspirate: While aspirating 100 ul of water to head 0 channels 0-7 - well A1@trough1 only contains 400 ul working volume, reducing aspirated volume by 50 ul",
			},
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{
					{0, "water", 50},
					{1, "water", 50},
					{2, "water", 50},
					{3, "water", 50},
					{4, "water", 50},
					{5, "water", 50},
					{6, "water", 50},
					{7, "water", 50},
				}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			"Fail - Aspirate with no tip",
			nil,
			[]*SetupFn{
				testLayout(),
				prefillWells("input_1", []string{"A1", "B1"}, "water", 200.),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "input_1", "", "", "", "", "", ""}, //deckposition
					[]string{"A1", "B1", "", "", "", "", "", ""},           //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                          //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},              //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},              //offsetY
					[]float64{1., 52.2, 1., 1., 1., 1., 1., 1.},            //offsetZ
					[]string{"plate", "plate", "", "", "", "", "", ""},     //plate_type
					0, //head
				},
				&Aspirate{
					[]float64{100., 100., 0., 0., 0., 0., 0., 0.},                  //volume     []float64
					[]bool{false, false, false, false, false, false, false, false}, //overstroke []bool
					0, //head       int
					2, //multi      int
					[]string{"plate", "plate", "", "", "", "", "", ""},             //platetype  []string
					[]string{"water", "water", "", "", "", "", "", ""},             //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //llf        []bool
				},
			},
			[]string{ //errors
				"(err) Aspirate: While aspirating 100 ul of water to head 0 channels 0-1 - missing tip on channel 1",
			},
			nil, //assertions
		},
		{
			"Fail - Underfull tip",
			nil,
			[]*SetupFn{
				testLayout(),
				prefillWells("input_1", []string{"A1"}, "water", 200.),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"A1", "", "", "", "", "", "", ""},      //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},       //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
				&Aspirate{
					[]float64{20., 0., 0., 0., 0., 0., 0., 0.},                     //volume     []float64
					[]bool{false, false, false, false, false, false, false, false}, //overstroke []bool
					0, //head       int
					1, //multi      int
					[]string{"plate", "", "", "", "", "", "", ""},                  //platetype  []string
					[]string{"water", "", "", "", "", "", "", ""},                  //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //llf        []bool
				},
			},
			[]string{ //errors
				"(warn) Aspirate: While aspirating 20 ul of water to head 0 channel 0 - minimum tip volume is 50 ul",
			},
			nil, //assertions
		},
		{
			"Fail - Overfull tip",
			nil,
			[]*SetupFn{
				testLayout(),
				prefillWells("input_1", []string{"A1", "B1", "C1", "D1", "E1", "F1"}, "water", 200.),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"A1", "", "", "", "", "", "", ""},      //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},       //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
				&Aspirate{
					[]float64{175., 0., 0., 0., 0., 0., 0., 0.},                    //volume     []float64
					[]bool{false, false, false, false, false, false, false, false}, //overstroke []bool
					0, //head       int
					1, //multi      int
					[]string{"plate", "", "", "", "", "", "", ""},                  //platetype  []string
					[]string{"water", "", "", "", "", "", "", ""},                  //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //llf        []bool
				},
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"B1", "", "", "", "", "", "", ""},      //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},       //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
				&Aspirate{
					[]float64{175., 0., 0., 0., 0., 0., 0., 0.},                    //volume     []float64
					[]bool{false, false, false, false, false, false, false, false}, //overstroke []bool
					0, //head       int
					1, //multi      int
					[]string{"plate", "", "", "", "", "", "", ""},                  //platetype  []string
					[]string{"water", "", "", "", "", "", "", ""},                  //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //llf        []bool
				},
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"C1", "", "", "", "", "", "", ""},      //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},       //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
				&Aspirate{
					[]float64{175., 0., 0., 0., 0., 0., 0., 0.},                    //volume     []float64
					[]bool{false, false, false, false, false, false, false, false}, //overstroke []bool
					0, //head       int
					1, //multi      int
					[]string{"plate", "", "", "", "", "", "", ""},                  //platetype  []string
					[]string{"water", "", "", "", "", "", "", ""},                  //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //llf        []bool
				},
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"D1", "", "", "", "", "", "", ""},      //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},       //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
				&Aspirate{
					[]float64{175., 0., 0., 0., 0., 0., 0., 0.},                    //volume     []float64
					[]bool{false, false, false, false, false, false, false, false}, //overstroke []bool
					0, //head       int
					1, //multi      int
					[]string{"plate", "", "", "", "", "", "", ""},                  //platetype  []string
					[]string{"water", "", "", "", "", "", "", ""},                  //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //llf        []bool
				},
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"E1", "", "", "", "", "", "", ""},      //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},       //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
				&Aspirate{
					[]float64{175., 0., 0., 0., 0., 0., 0., 0.},                    //volume     []float64
					[]bool{false, false, false, false, false, false, false, false}, //overstroke []bool
					0, //head       int
					1, //multi      int
					[]string{"plate", "", "", "", "", "", "", ""},                  //platetype  []string
					[]string{"water", "", "", "", "", "", "", ""},                  //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //llf        []bool
				},
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"F1", "", "", "", "", "", "", ""},      //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},       //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
				&Aspirate{
					[]float64{175., 0., 0., 0., 0., 0., 0., 0.},                    //volume     []float64
					[]bool{false, false, false, false, false, false, false, false}, //overstroke []bool
					0, //head       int
					1, //multi      int
					[]string{"plate", "", "", "", "", "", "", ""},                  //platetype  []string
					[]string{"water", "", "", "", "", "", "", ""},                  //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //llf        []bool
				},
			},
			[]string{ //errors
				"(err) Aspirate: While aspirating 175 ul of water to head 0 channel 0 - channel 0 contains 875 ul, command exceeds maximum volume 1000 ul",
			},
			nil, //assertions
		},
		{
			"Fail - non-independent head can only aspirate equal volumes",
			nil,
			[]*SetupFn{
				testLayout(),
				prefillWells("input_1", []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}, "water", 200.),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"}, //deckposition
					[]string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},                                         //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                                                                    //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},                                                        //offsetZ
					[]string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"},                 //plate_type
					0, //head
				},
				&Aspirate{
					[]float64{50., 60., 70., 80., 90., 100., 110., 120.},           //volume     []float64
					[]bool{false, false, false, false, false, false, false, false}, //overstroke []bool
					0, //head       int
					8, //multi      int
					[]string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"}, //platetype  []string
					[]string{"water", "water", "water", "water", "water", "water", "water", "water"}, //what       []string
					[]bool{false, false, false, false, false, false, false, false},                   //llf        []bool
				},
			},
			[]string{ //errors
				"(err) Aspirate: While aspirating {50,60,70,80,90,100,110,120} ul of water to head 0 channels 0-7 - channels cannot aspirate different volumes in non-independent head",
			},
			nil, //assertions
		},
		{
			"Fail - tip not in well",
			nil,
			[]*SetupFn{
				testLayout(),
				prefillWells("input_1", []string{"A1"}, "water", 200.),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"A1", "", "", "", "", "", "", ""},      //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{50., 1., 1., 1., 1., 1., 1., 1.},      //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
				&Aspirate{
					[]float64{100., 0., 0., 0., 0., 0., 0., 0.},                    //volume     []float64
					[]bool{false, false, false, false, false, false, false, false}, //overstroke []bool
					0, //head       int
					1, //multi      int
					[]string{"plate", "", "", "", "", "", "", ""},                  //platetype  []string
					[]string{"water", "", "", "", "", "", "", ""},                  //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //llf        []bool
				},
			},
			[]string{ //errors
				"(err) Aspirate: While aspirating 100 ul of water to head 0 channel 0 - tip on channel 0 not in a well",
			},
			nil, //assertions
		},
		{
			"Fail - Well doesn't contain enough",
			nil,
			[]*SetupFn{
				testLayout(),
				prefillWells("input_1", []string{"A1"}, "water", 200.),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"A1", "", "", "", "", "", "", ""},      //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},       //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
				&Aspirate{
					[]float64{535, 0., 0., 0., 0., 0., 0., 0.},                     //volume     []float64
					[]bool{false, false, false, false, false, false, false, false}, //overstroke []bool
					0, //head       int
					1, //multi      int
					[]string{"plate", "", "", "", "", "", "", ""},                  //platetype  []string
					[]string{"water", "", "", "", "", "", "", ""},                  //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //llf        []bool
				},
			},
			[]string{ //errors
				"(warn) Aspirate: While aspirating 535 ul of water to head 0 channel 0 - well A1@plate1 only contains 195 ul working volume, reducing aspirated volume by 340 ul",
			},
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				adaptorAssertion(0, []tipDesc{
					{0, "water", 195},
				}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		/*		{
				"Fail - wrong liquid type",
				nil,
				[]*SetupFn{
					testLayout(),
					prefillWells("input_1", []string{"A1"}, "water", 200.),
					preloadAdaptorTips(0, "tipbox_1", []int{0}),
				},
				[]TestRobotInstruction{
					&Move{
						[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
						[]string{"A1", "", "", "", "", "", "", ""},      //wellcoords
						[]int{0, 0, 0, 0, 0, 0, 0, 0},                   //reference
						[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
						[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
						[]float64{1., 1., 1., 1., 1., 1., 1., 1.},       //offsetZ
						[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
						0, //head
					},
					&Aspirate{
						[]float64{102.1, 0., 0., 0., 0., 0., 0., 0.},                   //volume     []float64
						[]bool{false, false, false, false, false, false, false, false}, //overstroke []bool
						0, //head       int
						1, //multi      int
						[]string{"plate", "", "", "", "", "", "", ""},                  //platetype  []string
						[]string{"ethanol", "", "", "", "", "", "", ""},                //what       []string
						[]bool{false, false, false, false, false, false, false, false}, //llf        []bool
					},
				},
				[]string{ //errors
					"(warn) Aspirate: While aspirating 102 ul of ethanol to head 0 channel 0 - well A1@plate1 contains water, not ethanol",
				},
				nil, //assertions
			},*/
		{
			"Fail - inadvertant aspiration",
			nil,
			[]*SetupFn{
				testLayout(),
				prefillWells("input_1", []string{"A1", "B1"}, "water", 200.),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1}),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"A1", "", "", "", "", "", "", ""},      //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},       //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
				&Aspirate{
					[]float64{98.6, 0., 0., 0., 0., 0., 0., 0.},                    //volume     []float64
					[]bool{false, false, false, false, false, false, false, false}, //overstroke []bool
					0, //head       int
					1, //multi      int
					[]string{"plate", "", "", "", "", "", "", ""},                  //platetype  []string
					[]string{"water", "", "", "", "", "", "", ""},                  //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //llf        []bool
				},
			},
			[]string{ //errors
				"(err) Aspirate: While aspirating 98.6 ul of water to head 0 channel 0 - channel 1 will inadvertantly aspirate water from well B1@plate1 as head is not independent",
			},
			nil, //assertions
		},
	}

	for _, test := range tests {
		test.run(t)
	}
}

func Test_Dispense(t *testing.T) {

	tests := []SimulatorTest{
		{
			"OK - single channel",
			nil,
			[]*SetupFn{
				testLayout(),
				preloadFilledTips(0, "tipbox_1", []int{0}, "water", 100.),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"A1", "", "", "", "", "", "", ""},      //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},       //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
				&Dispense{
					[]float64{50., 0., 0., 0., 0., 0., 0., 0.},                     //volume    []float64
					[]bool{false, false, false, false, false, false, false, false}, //blowout   []bool
					0, //head      int
					1, //multi     int
					[]string{"plate", "", "", "", "", "", "", ""},                  //platetype []string
					[]string{"water", "", "", "", "", "", "", ""},                  //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //llf        []bool
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				plateAssertion("input_1", []wellDesc{{"A1", "water", 50.}}),
				adaptorAssertion(0, []tipDesc{{0, "water", 50.}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			"OK - mixing",
			nil,
			[]*SetupFn{
				testLayout(),
				preloadFilledTips(0, "tipbox_1", []int{0}, "water", 100.),
				prefillWells("input_1", []string{"A1"}, "green", 50.),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"A1", "", "", "", "", "", "", ""},      //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},       //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
				&Dispense{
					[]float64{50., 0., 0., 0., 0., 0., 0., 0.},                     //volume    []float64
					[]bool{false, false, false, false, false, false, false, false}, //blowout   []bool
					0, //head      int
					1, //multi     int
					[]string{"plate", "", "", "", "", "", "", ""},                  //platetype []string
					[]string{"water", "", "", "", "", "", "", ""},                  //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //llf        []bool
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				plateAssertion("input_1", []wellDesc{{"A1", "green+water", 100.}}),
				adaptorAssertion(0, []tipDesc{{0, "water", 50.}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			"OK - single channel slightly above well",
			nil,
			[]*SetupFn{
				testLayout(),
				preloadFilledTips(0, "tipbox_1", []int{0}, "water", 100.),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"A1", "", "", "", "", "", "", ""},      //wellcoords
					[]int{1, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{3., 1., 1., 1., 1., 1., 1., 1.},       //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
				&Dispense{
					[]float64{50., 0., 0., 0., 0., 0., 0., 0.},                     //volume    []float64
					[]bool{false, false, false, false, false, false, false, false}, //blowout   []bool
					0, //head      int
					1, //multi     int
					[]string{"plate", "", "", "", "", "", "", ""},                  //platetype []string
					[]string{"water", "", "", "", "", "", "", ""},                  //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //llf        []bool
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				plateAssertion("input_1", []wellDesc{{"A1", "water", 50.}}),
				adaptorAssertion(0, []tipDesc{{0, "water", 50.}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			"OK - 8 channel",
			nil,
			[]*SetupFn{
				testLayout(),
				preloadFilledTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}, "water", 100.),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"}, //deckposition
					[]string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},                                         //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                                                                    //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},                                                        //offsetZ
					[]string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"},                 //plate_type
					0, //head
				},
				&Dispense{
					[]float64{50., 50., 50., 50., 50., 50., 50., 50.},              //volume    []float64
					[]bool{false, false, false, false, false, false, false, false}, //blowout   []bool
					0, //head      int
					8, //multi     int
					[]string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"}, //platetype []string
					[]string{"water", "water", "water", "water", "water", "water", "water", "water"}, //what       []string
					[]bool{false, false, false, false, false, false, false, false},                   //llf        []bool
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
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
			"Fail - no tips",
			nil,
			[]*SetupFn{
				testLayout(),
				//preloadFilledTips(0, "tipbox_1", []int{0}, "water", 100.),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"A1", "", "", "", "", "", "", ""},      //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},       //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
				&Dispense{
					[]float64{50., 0., 0., 0., 0., 0., 0., 0.},                     //volume    []float64
					[]bool{false, false, false, false, false, false, false, false}, //blowout   []bool
					0, //head      int
					1, //multi     int
					[]string{"plate", "", "", "", "", "", "", ""},                  //platetype []string
					[]string{"water", "", "", "", "", "", "", ""},                  //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //llf        []bool
				},
			},
			[]string{ //errors
				"(err) Dispense: 50 ul of water from head 0 channel 0 to A1@plate1 : no tip loaded on channel 0",
			},
			nil, //assertionsi
		},
		{
			"Fail - not enough in tip",
			nil,
			[]*SetupFn{
				testLayout(),
				preloadFilledTips(0, "tipbox_1", []int{0}, "water", 100.),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"A1", "", "", "", "", "", "", ""},      //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},       //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
				&Dispense{
					[]float64{150., 0., 0., 0., 0., 0., 0., 0.},                    //volume    []float64
					[]bool{false, false, false, false, false, false, false, false}, //blowout   []bool
					0, //head      int
					1, //multi     int
					[]string{"plate", "", "", "", "", "", "", ""},                  //platetype []string
					[]string{"water", "", "", "", "", "", "", ""},                  //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //llf        []bool
				},
			},
			[]string{ //errors
				"(warn) Dispense: 150 ul of water from head 0 channel 0 to A1@plate1 : tip on channel 0 contains only 100 ul, but blowout flag is false",
			},
			nil, //assertionsi
		},
		{
			"Fail - well over-full",
			nil,
			[]*SetupFn{
				testLayout(),
				preloadFilledTips(0, "tipbox_1", []int{0}, "water", 1000.),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"A1", "", "", "", "", "", "", ""},      //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},       //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
				&Dispense{
					[]float64{500., 0., 0., 0., 0., 0., 0., 0.},                    //volume    []float64
					[]bool{false, false, false, false, false, false, false, false}, //blowout   []bool
					0, //head      int
					1, //multi     int
					[]string{"plate", "", "", "", "", "", "", ""},                  //platetype []string
					[]string{"water", "", "", "", "", "", "", ""},                  //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //llf        []bool
				},
			},
			[]string{ //errors
				"(warn) Dispense: 500 ul of water from head 0 channel 0 to A1@plate1 : overfilling well A1@plate1 to 500 ul of 200 ul max volume",
			},
			nil, //assertionsi
		},
		{
			"Fail - not in a well",
			nil,
			[]*SetupFn{
				testLayout(),
				preloadFilledTips(0, "tipbox_1", []int{0}, "water", 100.),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"A1", "", "", "", "", "", "", ""},      //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{200., 1., 1., 1., 1., 1., 1., 1.},     //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
				&Dispense{
					[]float64{50., 0., 0., 0., 0., 0., 0., 0.},                     //volume    []float64
					[]bool{false, false, false, false, false, false, false, false}, //blowout   []bool
					0, //head      int
					1, //multi     int
					[]string{"plate", "", "", "", "", "", "", ""},                  //platetype []string
					[]string{"water", "", "", "", "", "", "", ""},                  //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //llf        []bool
				},
			},
			[]string{ //errors
				"(err) Dispense: 50 ul of water from head 0 channel 0 to @<unnamed> : no well within 5 mm below tip on channel 0",
			},
			nil, //assertionsi
		},
		{
			"Fail - dispensing to tipwaste",
			nil,
			[]*SetupFn{
				testLayout(),
				preloadFilledTips(0, "tipbox_1", []int{0}, "water", 100.),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"tipwaste", "", "", "", "", "", "", ""}, //deckposition
					[]string{"A1", "", "", "", "", "", "", ""},       //wellcoords
					[]int{1, 0, 0, 0, 0, 0, 0, 0},                    //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},        //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},        //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},        //offsetZ
					[]string{"tipwaste", "", "", "", "", "", "", ""}, //plate_type
					0, //head
				},
				&Dispense{
					[]float64{50., 0., 0., 0., 0., 0., 0., 0.},                     //volume    []float64
					[]bool{false, false, false, false, false, false, false, false}, //blowout   []bool
					0, //head      int
					1, //multi     int
					[]string{"tipwaste", "", "", "", "", "", "", ""},               //platetype []string
					[]string{"water", "", "", "", "", "", "", ""},                  //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //llf        []bool
				},
			},
			[]string{ //errors
				"(warn) Dispense: 50 ul of water from head 0 channel 0 to A1@tipwaste : dispensing to tipwaste",
			},
			nil, //assertionsi
		},
		{
			"fail - independence",
			nil,
			[]*SetupFn{
				testLayout(),
				preloadFilledTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}, "water", 100.),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"A1", "", "", "", "", "", "", ""},      //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},       //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
				&Dispense{
					[]float64{50., 0, 0, 0, 0, 0, 0, 0},                            //volume    []float64
					[]bool{false, false, false, false, false, false, false, false}, //blowout   []bool
					0, //head      int
					1, //multi     int
					[]string{"plate", "", "", "", "", "", "", ""},                  //platetype []string
					[]string{"water", "", "", "", "", "", "", ""},                  //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //llf        []bool
				},
			},
			[]string{ //errors
				"(err) Dispense: 50 ul of water from head 0 channel 0 to A1@plate1 : must also dispense 50 ul from channels 1-7 as head is not independent",
			},
			nil, //assertions
		},
		{
			"Fail - independence, different volumes",
			nil,
			[]*SetupFn{
				testLayout(),
				preloadFilledTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}, "water", 100.),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"}, //deckposition
					[]string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},                                         //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                                                                    //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},                                                        //offsetZ
					[]string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"},                 //plate_type
					0, //head
				},
				&Dispense{
					[]float64{50., 60., 50., 50., 50., 50., 50., 50.},              //volume    []float64
					[]bool{false, false, false, false, false, false, false, false}, //blowout   []bool
					0, //head      int
					8, //multi     int
					[]string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"}, //platetype []string
					[]string{"water", "water", "water", "water", "water", "water", "water", "water"}, //what       []string
					[]bool{false, false, false, false, false, false, false, false},                   //llf        []bool
				},
			},
			[]string{ //errors
				"(err) Dispense: {50,60,50,50,50,50,50,50} ul of water from head 0 channels 0-7 to A1-H1@plate1 : channels cannot dispense different volumes in non-independent head",
			},
			nil, //assertions
		},
	}

	for _, test := range tests {
		test.run(t)
	}
}

func Test_Mix(t *testing.T) {

	tests := []SimulatorTest{
		{
			"OK - single channel",
			nil,
			[]*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
				prefillWells("input_1", []string{"A1"}, "water", 200.),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"A1", "", "", "", "", "", "", ""},      //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},       //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
				&Mix{
					0, //head      int
					[]float64{50., 0., 0., 0., 0., 0., 0., 0.},    //volume    []float64
					[]string{"plate", "", "", "", "", "", "", ""}, //platetype []string
					[]int{5, 0, 0, 0, 0, 0, 0, 0},                 //cycles []int
					1, //multi     int
					[]string{"water", "", "", "", "", "", "", ""},                  //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //blowout   []bool
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				plateAssertion("input_1", []wellDesc{{"A1", "water", 200.}}),
				adaptorAssertion(0, []tipDesc{{0, "water", 0.}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			"OK - 8 channel",
			nil,
			[]*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}),
				prefillWells("input_1", []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}, "water", 200.),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"}, //deckposition
					[]string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},                                         //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                                                                    //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},                                                        //offsetZ
					[]string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"},                 //plate_type
					0, //head
				},
				&Mix{
					0, //head      int
					[]float64{50., 50., 50., 50., 50., 50., 50., 50.},                                //volume    []float64
					[]string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"}, //platetype []string
					[]int{5, 5, 5, 5, 5, 5, 5, 5},                                                    //cycles []int
					8, //multi     int
					[]string{"water", "water", "water", "water", "water", "water", "water", "water"}, //what       []string
					[]bool{false, false, false, false, false, false, false, false},                   //blowout   []bool
				},
			},
			nil, //errors
			[]*AssertionFn{ //assertions
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
					{0, "water", 0.},
					{1, "water", 0.},
					{2, "water", 0.},
					{3, "water", 0.},
					{4, "water", 0.},
					{5, "water", 0.},
					{6, "water", 0.},
					{7, "water", 0.},
				}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
		{
			"Fail - independece problems",
			nil,
			[]*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0, 1, 2, 3, 4, 5, 6, 7}),
				prefillWells("input_1", []string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"}, "water", 200.),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1", "input_1"}, //deckposition
					[]string{"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1"},                                         //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                                                                    //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},                                                        //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},                                                        //offsetZ
					[]string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"},                 //plate_type
					0, //head
				},
				&Mix{
					0, //head      int
					[]float64{50., 60., 50., 50., 50., 50., 50., 50.},                                //volume    []float64
					[]string{"plate", "plate", "plate", "plate", "plate", "plate", "plate", "plate"}, //platetype []string
					[]int{5, 5, 5, 5, 5, 2, 2, 2},                                                    //cycles []int
					8, //multi     int
					[]string{"water", "water", "water", "water", "water", "water", "water", "water"}, //what       []string
					[]bool{false, false, false, false, false, false, false, false},                   //blowout   []bool
				},
			},
			[]string{ //errors
				"(err) Mix: While mixing {50,60,50,50,50,50,50,50} ul {5,5,5,5,5,2,2,2} times in wells A1,B1,C1,D1,E1,F1,G1,H1 of plate \"plate1\" - cannot manipulate different volumes with non-independent head",
				"(err) Mix: While mixing {50,60,50,50,50,50,50,50} ul {5,5,5,5,5,2,2,2} times in wells A1,B1,C1,D1,E1,F1,G1,H1 of plate \"plate1\" - cannot vary number of mix cycles with non-independent head",
			},
			nil, //assertions
		},
		{
			"Fail - wrong platetype",
			nil,
			[]*SetupFn{
				testLayout(),
				preloadAdaptorTips(0, "tipbox_1", []int{0}),
				prefillWells("input_1", []string{"A1"}, "water", 200.),
			},
			[]TestRobotInstruction{
				&Move{
					[]string{"input_1", "", "", "", "", "", "", ""}, //deckposition
					[]string{"A1", "", "", "", "", "", "", ""},      //wellcoords
					[]int{0, 0, 0, 0, 0, 0, 0, 0},                   //reference
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetX
					[]float64{0., 0., 0., 0., 0., 0., 0., 0.},       //offsetY
					[]float64{1., 1., 1., 1., 1., 1., 1., 1.},       //offsetZ
					[]string{"plate", "", "", "", "", "", "", ""},   //plate_type
					0, //head
				},
				&Mix{
					0, //head      int
					[]float64{50., 0., 0., 0., 0., 0., 0., 0.},        //volume    []float64
					[]string{"notaplate", "", "", "", "", "", "", ""}, //platetype []string
					[]int{5, 0, 0, 0, 0, 0, 0, 0},                     //cycles []int
					1, //multi     int
					[]string{"water", "", "", "", "", "", "", ""},                  //what       []string
					[]bool{false, false, false, false, false, false, false, false}, //blowout   []bool
				},
			},
			[]string{ //errors
				"(warn) Mix: While mixing 50 ul 5 times in well A1 of plate \"plate1\" - plate \"plate1\" is of type \"plate\", not \"notaplate\"",
			},
			[]*AssertionFn{ //assertions
				tipboxAssertion("tipbox_1", []string{}),
				tipboxAssertion("tipbox_2", []string{}),
				plateAssertion("input_1", []wellDesc{{"A1", "water", 200.}}),
				adaptorAssertion(0, []tipDesc{{0, "water", 0.}}),
				tipwasteAssertion("tipwaste", 0),
			},
		},
	}

	for _, test := range tests {
		test.run(t)
	}
}

func component(name string) *wtype.LHComponent {
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
				[]bool{false, false, false, false, false, false, false, false},
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
	input_plate := default_lhplate("input")
	output_plate := default_lhplate("output")

	//tips - using small tipbox so I don't have to worry about using different tips
	tipbox := small_lhtipbox("tipbox")

	//tipwaste
	tipwaste := default_lhtipwaste("tipwaste")

	//setup the input plate
	wc := wtype.MakeWellCoords("A1")
	comp := []*wtype.LHComponent{
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

	st := SimulatorTest{
		"Run Workflow",
		nil,
		[]*SetupFn{},
		inst,
		nil, //errors
		[]*AssertionFn{ //assertions
			tipboxAssertion("tipbox_1", []string{
				"H12", "G12",
				"A11", "B11", "C11", "D11", "E11", "F11", "G11", "H11",
				"A10", "B10", "C10", "D10", "E10", "F10", "G10", "H10",
				"A9", "B9", "C9", "D9", "E9", "F9", "G9", "H9",
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
				{"A2", "red", 170.},
				{"B2", "red", 170.},
				{"C2", "red", 170.},
				{"D2", "red", 170.},
				{"E2", "red", 170.},
				{"F2", "red", 170.},
				{"G2", "red", 170.},
				{"H2", "red", 170.},
				{"A3", "green", 5.},
				{"B3", "green", 32.86},
				{"C3", "green", 60.71},
				{"D3", "green", 88.57},
				{"E3", "green", 116.43},
				{"F3", "green", 144.29},
				{"G3", "green", 172.14},
				{"H3", "green", 200.},
				{"A4", "water", 200.},
				{"B4", "water", 172.14},
				{"C4", "water", 144.29},
				{"D4", "water", 116.43},
				{"E4", "water", 88.57},
				{"F4", "water", 60.71},
				{"G4", "water", 32.86},
				{"H4", "water", 5.},
				{"A5", "green", 135.},
				{"B5", "green+water", 135.},
				{"C5", "green+water", 135.},
				{"D5", "green+water", 135.},
				{"E5", "green+water", 135.},
				{"F5", "green+water", 135.},
				{"G5", "green+water", 135.},
				{"H5", "water", 135.},
			}),
			plateAssertion("output_1", []wellDesc{
				{"A1", "green+red+water", 20.},
				{"B1", "green+red+water", 20.},
				{"C1", "green+red+water", 20.},
				{"D1", "green+red+water", 20.},
				{"E1", "green+red+water", 20.},
				{"F1", "green+red+water", 20.},
				{"G1", "green+red+water", 20.},
				{"H1", "red+water", 20.},
				{"A2", "green+red+water", 20.},
				{"B2", "green+red+water", 20.},
				{"C2", "green+red+water", 20.},
				{"D2", "green+red+water", 20.},
				{"E2", "green+red+water", 20.},
				{"F2", "green+red+water", 20.},
				{"G2", "green+red+water", 20.},
				{"H2", "red+water", 20.},
				{"A3", "green+red+water", 20.},
				{"B3", "green+red+water", 20.},
				{"C3", "green+red+water", 20.},
				{"D3", "green+red+water", 20.},
				{"E3", "green+red+water", 20.},
				{"F3", "green+red+water", 20.},
				{"G3", "green+red+water", 20.},
				{"H3", "red+water", 20.},
				{"A4", "green+red+water", 20.},
				{"B4", "green+red+water", 20.},
				{"C4", "green+red+water", 20.},
				{"D4", "green+red+water", 20.},
				{"E4", "green+red+water", 20.},
				{"F4", "green+red+water", 20.},
				{"G4", "green+red+water", 20.},
				{"H4", "red+water", 20.},
				{"A5", "green+red+water", 20.},
				{"B5", "green+red+water", 20.},
				{"C5", "green+red+water", 20.},
				{"D5", "green+red+water", 20.},
				{"E5", "green+red+water", 20.},
				{"F5", "green+red+water", 20.},
				{"G5", "green+red+water", 20.},
				{"H5", "red+water", 20.},
				{"A6", "green+red+water", 20.},
				{"B6", "green+red+water", 20.},
				{"C6", "green+red+water", 20.},
				{"D6", "green+red+water", 20.},
				{"E6", "green+red+water", 20.},
				{"F6", "green+red+water", 20.},
				{"G6", "green+red+water", 20.},
				{"H6", "red+water", 20.},
				{"A7", "green+red+water", 20.},
				{"B7", "green+red+water", 20.},
				{"C7", "green+red+water", 20.},
				{"D7", "green+red+water", 20.},
				{"E7", "green+red+water", 20.},
				{"F7", "green+red+water", 20.},
				{"G7", "green+red+water", 20.},
				{"H7", "red+water", 20.},
				{"A8", "green+red+water", 20.},
				{"B8", "green+red+water", 20.},
				{"C8", "green+red+water", 20.},
				{"D8", "green+red+water", 20.},
				{"E8", "green+red+water", 20.},
				{"F8", "green+red+water", 20.},
				{"G8", "green+red+water", 20.},
				{"H8", "red+water", 20.},
				{"A9", "green+red+water", 20.},
				{"B9", "green+red+water", 20.},
				{"C9", "green+red+water", 20.},
				{"D9", "green+red+water", 20.},
				{"E9", "green+red+water", 20.},
				{"F9", "green+red+water", 20.},
				{"G9", "green+red+water", 20.},
				{"H9", "red+water", 20.},
				{"A10", "green+red+water", 20.},
				{"B10", "green+red+water", 20.},
				{"C10", "green+red+water", 20.},
				{"D10", "green+red+water", 20.},
				{"E10", "green+red+water", 20.},
				{"F10", "green+red+water", 20.},
				{"G10", "green+red+water", 20.},
				{"H10", "red+water", 20.},
				{"A11", "green+red+water", 20.},
				{"B11", "green+red+water", 20.},
				{"C11", "green+red+water", 20.},
				{"D11", "green+red+water", 20.},
				{"E11", "green+red+water", 20.},
				{"F11", "green+red+water", 20.},
				{"G11", "green+red+water", 20.},
				{"H11", "red+water", 20.},
				{"A12", "green+water", 20.},
				{"B12", "green+water", 20.},
				{"C12", "green+water", 20.},
				{"D12", "green+water", 20.},
				{"E12", "green+water", 20.},
				{"F12", "green+water", 20.},
				{"G12", "green+water", 20.},
				{"H12", "water", 20.},
			}),
			adaptorAssertion(0, []tipDesc{}),
			tipwasteAssertion("tipwaste", 26),
		},
	}

	st.run(t)

}
