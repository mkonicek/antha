// make_plate_library_test.go
package factory

import (
	//"fmt"
	"strings"
	"testing"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
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

	for _, test := range tests {
		for _, device := range Devices {

			testplatename := test.TestPlateName
			testplate := GetPlateByType(testplatename)

			testname := testplatename + "_" + device.GetName()

			testPlateInventory2.AddRiser(testplate, device)

			offset, _ := platespecificoffset[testplatename]

			// check if new plate with device is in inventory
			if _, found := testPlateInventory2.lib[testname]; !found {
				t.Error(
					"for", device, "\n",
					"testname", testname, "\n",
					"not in platelist", GetPlateList(), "\n",
				)
			}
			// check that the height is as expected using default inventory
			if testplate.Height != test.ExpectedHeight {
				t.Error(
					"for", testplatename, "\n",
					"Expected plate height:", test.ExpectedHeight, "\n",
					"got:", testplate.Height, "\n",
				)
			}
			// check that the height is as expected using replicated default inventory following AddRiser()
			if testPlateInventory2.lib[test.TestPlateName].Height != test.ExpectedHeight {
				t.Error(
					"for", "no device", "\n",
					"testname", test.TestPlateName, "\n",
					"Expected plate height:", test.ExpectedHeight, "\n",
					"got:", testPlateInventory2.lib[test.TestPlateName].Height, "\n",
				)
			}
			// check that the height is as expected with riser added
			if testPlateInventory2.lib[testname].Height != test.ExpectedHeight {
				t.Error(
					"for", device, "\n",
					"testname", testname, "\n",
					"Expected plate height:", test.ExpectedHeight, "\n",
					"got:", testPlateInventory2.lib[testname].Height, "\n",
				)
			}
			// now test z offsets
			if testplate.WellZStart != test.ExpectedZStart {
				t.Error(
					"for", testplatename, "\n",
					"Expected plate height:", test.ExpectedZStart, "\n",
					"got:", testplate.WellZStart, "\n",
				)
			}

			if testPlateInventory2.lib[testname].WellZStart != test.ExpectedZStart+device.GetHeightInmm()-offset {
				t.Error(
					"for", device, "\n",
					"testname", testname, "\n",
					"Expected plate height:", test.ExpectedZStart, "+", "device:", device.GetHeightInmm(), "=", test.ExpectedZStart+device.GetHeightInmm()-offset, "\n",
					"got:", testPlateInventory2.lib[testname].WellZStart, "\n",
				)
			}
			if testPlateInventory2.lib[test.TestPlateName].WellZStart != test.ExpectedZStart {
				t.Error(
					"for", "no device", "\n",
					"testname", test.TestPlateName, "\n",
					"Expected plate height:", test.ExpectedZStart, "\n",
					"got:", testPlateInventory2.lib[test.TestPlateName].WellZStart, "\n",
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

		_, ok := Devices[device.name]

		//n := Devices[device.name].GetName()

		if !ok {
			t.Error(
				"for", device.name, "\n",
				"not found in devices", Devices, "\n",
			)
		} else {
			c := Devices[device.name].GetConstraints()
			h := Devices[device.name].GetHeightInmm()
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

var testPlateInventory *plateLibrary

func init() {
	testPlateInventory = &plateLibrary{
		lib: makePlateLibrary(),
	}

}

var testPlateInventory2 *plateLibrary

func init() {
	testPlateInventory2 = &plateLibrary{
		lib: makePlateLibrary(),
	}

	//defaultPlateInventory.AddAllDevices()
	//defaultPlateInventory.AddAllRisers()
}

func TestSetConstraints(t *testing.T) {

	var ok bool
	allplates := GetPlateList()
	platform := "Pipetmax"
	expectedpositions := []string{"position_1"}
	var testname string
	for _, testplatename := range allplates {
		for _, device := range Devices {

			if device.GetConstraints() != nil {

				testplate := GetPlateByType(testplatename)

				if !search.InSlice(testplatename, exceptions[device.GetName()]) {

					if strings.Contains(testplatename, device.GetName()) {
						testname = testplatename
					} else if !ContainsRiser(testplate) {
						testname = testplatename + "_" + device.GetName()
					} else {
						continue
					}

					testPlateInventory.AddRiser(testplate, device)

					testplate, ok = testPlateInventory.lib[testname]

					if !ok {
						t.Error(
							"for", device.GetName(), "\n",
							"testname", testname, "not added correctly", "\n",
						)
					} else {

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
		}
	}
}

func TestGetConstraints(t *testing.T) {

	var ok bool
	allplates := GetPlateList()
	platform := "Pipetmax"
	expectedpositions := []string{"position_1"}
	var testname string
	for _, testplatename := range allplates {
		for _, device := range Devices {

			if device.GetConstraints() != nil {

				testplate := GetPlateByType(testplatename)

				if !search.InSlice(testplatename, exceptions[device.GetName()]) {

					if strings.Contains(testplatename, device.GetName()) {
						testname = testplatename
					} else if !ContainsRiser(testplate) {
						testname = testplatename + "_" + device.GetName()
					} else {
						continue
					}

					testplate, ok = defaultPlateLibrary.lib[testname]

					if !ok {
						t.Error(
							"for", device.GetName(), "\n",
							"testname", testname, "not added correctly", "\n",
						)
					} else {

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
		}
	}
}

func TestPlateZs(t *testing.T) {
	for _, test := range tests {

		testplate := GetPlateByType(test.TestPlateName)

		if testplate.WellZStart != test.ExpectedZStart {

			t.Error(
				"for", test.TestPlateName, "\n",
				"expected height: ", test.ExpectedZStart, "\n",
				"got height :", testplate.WellZStart, "\n",
			)
		}
	}
}
