package liquidhandling

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"testing"
)

// choose_plate_assignments(component_volumes map[string]wunit.Volume, plate_types []*wtype.LHPlate, weight_constraint map[string]float64) (map[string]map[*wtype.LHPlate]int, error)
func TestIPL1(t *testing.T) {
	type testCase struct {
		Name              string
		Component_volumes map[string]wunit.Volume
		Plate_types       []*wtype.LHPlate
		Weight_constraint map[string]float64
		Expected          map[string]map[string]int
	}

	testCases := []testCase{
		{
			Name:              "",
			Component_volumes: nil,
			Plate_types:       nil,
			Weight_constraint: nil,
			Expected:          nil,
		},
	}

	for _, testCase := range testCases {
		doTheTest := func(t *testing.T) {
			results, err := choose_plate_assignments(testCase.Component_volumes, testCase.Plate_types, testCase.Weight_constraint)

			if err != nil {
				t.Errorf("Expected %s got error %v", fmtIt(testCase.Expected), err.Error())
			}

			if !compare(results, testCase.Expected) {
				t.Errorf("Expected: %s got %s", fmtIt(testCase.Expected), fmtIt(results))
			}
		}

		t.Run(testCase.Name, doTheTest)
	}

}

// true iff results == expected
func compare(results map[string]map[*wtype.LHPlate]int, expected map[string]map[string]int) bool {
	if len(results) != len(expected) {
		return false
	}

	for k, v := range expected {
		v2, ok := results[k]

		if !ok {
			return false
		}

		if len(v) != len(v2) {
			return false
		}

		for k, i := range v2 {
			i2, ok := v[k.Type]

			if !ok || i != i2 {
				return false
			}
		}
	}

	return true
}

func fmtIt(rest interface{}) string {
	var ret string

	switch rest.(type) {
	case map[string]map[string]int:

		t := rest.(map[string]map[string]int)

		for k, v := range t {
			ret += fmt.Sprintf("%s: ", k)
			for k2, v2 := range v {
				ret += fmt.Sprintf("%s(%d) ", k2, v2)
			}
		}

	case map[string]map[*wtype.LHPlate]int:
		t := rest.(map[string]map[*wtype.LHPlate]int)

		for k, v := range t {
			ret += fmt.Sprintf("%s: ", k)
			for k2, v2 := range v {
				ret += fmt.Sprintf("%s(%d) ", k2.Type, v2)
			}
		}

	}

	return ret
}
