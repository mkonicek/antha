package testinventory

import (
	"context"
	"strings"
	"testing"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/inventory"
)

type platetest struct {
	TestPlateName  string
	ExpectedHeight float64
	ExpectedZStart float64
}

var tests = []platetest{
	platetest{TestPlateName: "reservoir", ExpectedZStart: 10.0, ExpectedHeight: 45.0},
	platetest{TestPlateName: "pcrplate_skirted", ExpectedZStart: 0.636, ExpectedHeight: 15.5},
	platetest{TestPlateName: "greiner384", ExpectedZStart: 2.5, ExpectedHeight: 14.0},
	platetest{TestPlateName: "Nuncon12well", ExpectedZStart: 4.0, ExpectedHeight: 19.0},
	platetest{TestPlateName: "Nuncon12wellAgar", ExpectedZStart: 9.0, ExpectedHeight: 19.0},
}

func TestAddRiser(t *testing.T) {
	ctx := NewContext(context.Background())

	for _, test := range tests {
		for _, device := range defaultDevices {
			testPlate, err := inventory.NewPlate(ctx, test.TestPlateName)
			if err != nil {
				t.Error(err)
				continue
			}

			testname := test.TestPlateName + "_" + device.GetName()

			newPlates := addRiser(testPlate, device)
			if e, f := 0, len(newPlates); e == f {
				t.Errorf("expected some plates but none found")
				continue
			}

			newPlate := newPlates[0]
			if e, f := testname, newPlate.Type; e != f {
				t.Errorf("expected %s but found %s", e, f)
			}

			offset, _ := platespecificoffset[test.TestPlateName]

			// check that the height is as expected using default inventory
			if testPlate.Height != test.ExpectedHeight {
				t.Error(
					"for", test.TestPlateName, "\n",
					"Expected plate height:", test.ExpectedHeight, "\n",
					"got:", testPlate.Height, "\n",
				)
			}

			// check that the height is as expected with riser added
			if f, e := newPlate.Height, test.ExpectedHeight; e != f {
				t.Error(
					"for", device, "\n",
					"testname", testname, "\n",
					"Expected plate height:", e, "\n",
					"got:", f, "\n",
				)
			}

			// now test z offsets
			if testPlate.WellZStart != test.ExpectedZStart {
				t.Error(
					"for", test.TestPlateName, "\n",
					"Expected plate height:", test.ExpectedZStart, "\n",
					"got:", testPlate.WellZStart, "\n",
				)
			}

			if f, e := newPlate.WellZStart, test.ExpectedZStart+device.GetHeightInmm()-offset; e != f {
				t.Error(
					"for", device, "\n",
					"testname", testname, "\n",
					"Expected plate height:", test.ExpectedZStart, "+",
					"device:", device.GetHeightInmm(), "=", e, "\n",
					"got:", f, "\n",
				)
			}

			if f, e := testPlate.WellZStart, test.ExpectedZStart; e != f {
				t.Error(
					"for", "no device", "\n",
					"testname", test.TestPlateName, "\n",
					"Expected plate height:", e, "\n",
					"got:", f, "\n",
				)
			}
		}
	}
}

type testdevice struct {
	name                string
	constraintdevice    string
	constraintposition1 string
	height              float64
}

var testdevices = []testdevice{
	testdevice{name: "bioshake", constraintdevice: "Pipetmax", constraintposition1: "position_1", height: 55.92},
}

type deviceExceptions map[string][]string // key is device name, exceptions are the plates which will give a result which differs from norm

var exceptions deviceExceptions = map[string][]string{
	"bioshake":                  []string{"EGEL96_1", "EGEL96_2", "EPAGE48", "EGEL48", "Nuncon12wellAgarD_incubator"},
	"bioshake_96well_adaptor":   []string{"EGEL96_1", "EGEL96_2", "EPAGE48", "EGEL48", "Nuncon12wellAgarD_incubator"},
	"bioshake_standard_adaptor": []string{"EGEL96_1", "EGEL96_2", "EPAGE48", "EGEL48", "Nuncon12wellAgarD_incubator"},
}

func TestDeviceMethods(t *testing.T) {

	for _, device := range testdevices {

		_, ok := defaultDevices[device.name]

		if !ok {
			t.Error(
				"for", device.name, "\n",
				"not found in devices", defaultDevices, "\n",
			)
		} else {
			c := defaultDevices[device.name].GetConstraints()
			h := defaultDevices[device.name].GetHeightInmm()
			//r := Devices[device].GetRiser()

			if constraints, ok := c[device.constraintdevice]; !ok || constraints[0] != device.constraintposition1 {
				t.Error(
					"for", device.name, "\n",
					"Constraints", c, "\n",
					"expected key", device.constraintdevice, "\n",
					"expected 1st position", device.constraintposition1, "\n",
				)
			}

			if h != device.height {
				t.Error(
					"for", device.name, "\n",
					"expectd height", device.height, "\n",
					"got", h, "\n",
				)
			}
		}

	}

}

func TestSetConstraints(t *testing.T) {
	ctx := NewContext(context.Background())

	platform := "Pipetmax"
	expectedpositions := []string{"position_1"}

	for _, testplate := range GetPlates(ctx) {
		for _, device := range defaultDevices {

			if device.GetConstraints() == nil {
				continue
			}

			if search.InSlice(testplate.Type, exceptions[device.GetName()]) {
				continue
			}

			newPlates := addRiser(testplate, device)

			if strings.Contains(testplate.Type, device.GetName()) {
				if e, f := 0, len(newPlates); e != f {
					t.Errorf("expecting %d plates found %d", e, f)
					continue
				}
			} else if !containsRiser(testplate) {
				if e, f := 1, len(newPlates); e != f {
					t.Errorf("expecting %d plates found %d", e, f)
					continue
				} else if e, f := testplate.Type+"_"+device.GetName(), newPlates[0].Type; e != f {
					t.Errorf("expecting type %s found %s", e, f)
					continue
				}
			} else {
				continue
			}

			for _, testplate := range newPlates {
				positionsinterface, found := testplate.Welltype.Extra[platform]
				positions, ok := positionsinterface.([]string)
				if !ok || !found || positions == nil || len(positions) != len(expectedpositions) || positions[0] != expectedpositions[0] {
					t.Error(
						"for", device, "\n",
						"testname", testplate.Type, "\n",
						"Extra found", found, "\n",
						"[]string?: ", ok, "\n",
						"Positions: ", positions, "\n",
						"expected positions: ", expectedpositions, "\n",
						"Constraints expected :", device.GetConstraints()[platform], "\n",
						"Constraints got :", testplate.Welltype.Extra[platform], "\n",
					)
				}
			}
		}
	}
}

func TestGetConstraints(t *testing.T) {
	ctx := NewContext(context.Background())

	platform := "Pipetmax"
	expectedpositions := []string{"position_1"}
	for _, testplate := range GetPlates(ctx) {
		for _, device := range defaultDevices {

			if device.GetConstraints() == nil {
				continue
			}

			if search.InSlice(testplate.Type, exceptions[device.GetName()]) {
				continue
			}

			var testname string
			if strings.Contains(testplate.Type, device.GetName()) {
				testname = testplate.Type
			} else if !containsRiser(testplate) {
				testname = testplate.Type + "_" + device.GetName()
			} else {
				continue
			}

			testplate, err := inventory.NewPlate(ctx, testname)
			if err != nil {
				t.Error(err)
				continue
			}

			positionsinterface, found := testplate.Welltype.Extra[platform]
			positions, ok := positionsinterface.([]string)
			if !ok || !found || positions == nil || len(positions) != len(expectedpositions) || positions[0] != expectedpositions[0] {
				t.Error(
					"for", device, "\n",
					"testname", testname, "\n",

					"Extra found", found, "\n",
					"[]string?: ", ok, "\n",
					"Positions: ", positions, "\n",
					"expected positions: ", expectedpositions, "\n",
					"Constraints expected :", device.GetConstraints()[platform], "\n",
					"Constraints got :", testplate.Welltype.Extra[platform], "\n",
				)
			}
		}
	}
}

func TestPlateZs(t *testing.T) {
	ctx := NewContext(context.Background())

	for _, test := range tests {

		testplate, err := inventory.NewPlate(ctx, test.TestPlateName)
		if err != nil {
			t.Error(err)
			continue
		}

		if testplate.WellZStart != test.ExpectedZStart {
			t.Error(
				"for", test.TestPlateName, "\n",
				"expected height: ", test.ExpectedZStart, "\n",
				"got height :", testplate.WellZStart, "\n",
			)
		}
	}
}
