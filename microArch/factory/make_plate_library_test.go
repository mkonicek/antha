// make_plate_library_test.go
package factory

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"strings"
	"testing"
)

type platetest struct {
	TestPlateName string
}

var tests = []platetest{platetest{TestPlateName: "reservoir"}}

func TestAddRiser(t *testing.T) {

	for _, test := range tests {
		for _, device := range Devices {

			testplatename := test.TestPlateName
			testplate := GetPlateByType(testplatename)
			testname := testplatename + "_" + device.GetName()

			defaultPlateInventory.AddRiser(testplate, device)
			if _, found := defaultPlateInventory.inv[testname]; !found {
				t.Error(
					"for", device, "\n",
					"testname", testname, "\n",
					"not in platelist", GetPlateList(), "\n",
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
	testdevice{name: "incubator", constraintdevice: "Pipetmax", constraintposition1: "position_1", height: 55.92},
}

type deviceExceptions map[string][]string // key is device name, exceptions are the plates which will give a result which differs from norm

var exceptions deviceExceptions = map[string][]string{
	"incubator": []string{"EGEL96_1", "EGEL96_2", "EPAGE48", "Nuncon12wellAgarD_incubator"},
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

var testPlateInventory *PlateInventory

func init() {
	testPlateInventory = &PlateInventory{
		inv: makePlateLibrary(),
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

					testplate, ok = testPlateInventory.inv[testname]

					if !ok {
						t.Error(
							"for", device.GetName(), "\n",
							"testname", testname, "not added correctly", "\n",
						)
					} else {

						positionsinterface, found := testplate.Welltype.Extra[platform]
						positions, ok := positionsinterface.([]string)
						fmt.Println("testplate: ", testname, " Constraints: ", positions)
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

					testplate, ok = defaultPlateInventory.inv[testname]

					if !ok {
						t.Error(
							"for", device.GetName(), "\n",
							"testname", testname, "not added correctly", "\n",
						)
					} else {

						positionsinterface, found := testplate.Welltype.Extra[platform]
						positions, ok := positionsinterface.([]string)
						fmt.Println("testplate: ", testname, " Constraints: ", positions)
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
