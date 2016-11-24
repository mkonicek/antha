// make_plate_library_test.go
package factory

import (
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
