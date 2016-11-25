// make_plate_library_test.go
package factory

import (
	"fmt"
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

func TestDeviceMethods(t *testing.T) {

	for _, device := range Devices {

		fmt.Println(
			device.GetName(),
			device.GetConstraints(),
			device.GetHeightInmm(),
			device.GetRiser(),
		)
		/*
			if _, found := defaultPlateInventory.inv[testname]; !found {
				t.Error(
					"for", device, "\n",
					"testname", testname, "\n",
					"not in platelist", GetPlateList(), "\n",
				)
			}
		*/
	}

}

func TestSetConstraints(t *testing.T) {

	allplates := GetPlateList()
	platform := "PipetMax"
	expectedpositions := []string{"position_1"}
	var testname string
	for _, testplatename := range allplates {
		for _, device := range Devices {

			if device.GetConstraints() != nil {

				testplate := GetPlateByType(testplatename)

				if strings.Contains(testplatename, device.GetName()) {
					testname = testplatename
				} else if !ContainsRiser(testplate) {
					testname = testplatename + "_" + device.GetName()
				} else {
					continue
				}

				defaultPlateInventory.AddRiser(testplate, device)
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
