package liquidhandling

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/cache/plateCache"
	"github.com/antha-lang/antha/inventory/testinventory"
	"testing"
)

func summarize(m map[string]map[*wtype.LHPlate]int) map[string]map[string]int {
	r := make(map[string]map[string]int, len(m))

	for k, v := range m {
		mm := make(map[string]int, len(v))

		for kk, vv := range v {
			mm[kk.Type] = vv
		}

		r[k] = mm
	}

	return r
}

// choose_plate_assignments(component_volumes map[string]wunit.Volume, plate_types []*wtype.LHPlate, weight_constraint map[string]float64) (map[string]map[*wtype.LHPlate]int, error)
func TestIPL1(t *testing.T) {
	// a few plates

	ctx := testinventory.NewContext(context.Background())
	ctx = plateCache.NewContext(ctx)

	pcrplateSkirted, err := inventory.NewPlate(ctx, "pcrplate_skirted")

	if err != nil {
		t.Errorf("Can't get pcrplate_skirted")
	}

	DWST12, err := inventory.NewPlate(ctx, "DWST12")
	if err != nil {
		t.Errorf("Can't get DWST12")
	}
	DSW96, err := inventory.NewPlate(ctx, "DSW96")
	if err != nil {
		t.Errorf("Can't get DSW96")
	}

	/*

		reservoir, err := inventory.NewPlate(ctx, "reservoir")
		if err != nil {
			t.Errorf("Can't get reservoir")
		}
	*/

	weightConstraint := map[string]float64{"MAX_N_WELLS": 98, "RESIDUAL_VOLUME_WEIGHT": 1.0}

	type testCase struct {
		Name              string
		Component_volumes map[string]wunit.Volume
		Plate_types       []*wtype.LHPlate
		Weight_constraint map[string]float64
		Expected          map[string]map[string]int
	}

	testCases := []testCase{
		{
			Name:              "1thing1well",
			Component_volumes: map[string]wunit.Volume{"water": wunit.NewVolume(50, "ul")},
			Plate_types:       []*wtype.LHPlate{pcrplateSkirted},
			Weight_constraint: weightConstraint,
			Expected:          map[string]map[string]int{"water": {"pcrplate_skirted": 1}},
		},
		{
			Name:              "1thing3wells",
			Component_volumes: map[string]wunit.Volume{"water": wunit.NewVolume(500, "ul")},
			Plate_types:       []*wtype.LHPlate{pcrplateSkirted},
			Weight_constraint: weightConstraint,
			Expected:          map[string]map[string]int{"water": {"pcrplate_skirted": 3}},
		},
		{
			Name:              "1thing2plates",
			Component_volumes: map[string]wunit.Volume{"water": wunit.NewVolume(18915, "ul")},
			Plate_types:       []*wtype.LHPlate{pcrplateSkirted},
			Weight_constraint: weightConstraint,
			Expected:          map[string]map[string]int{"water": {"pcrplate_skirted": 97}},
		},
		{
			Name:              "failWellConstraint129wells",
			Component_volumes: map[string]wunit.Volume{"water": wunit.NewVolume(25155, "ul")},
			Plate_types:       []*wtype.LHPlate{pcrplateSkirted},
			Weight_constraint: weightConstraint,
			Expected:          map[string]map[string]int{"FAIL": {}},
		},
		{
			Name:              "simplechoiceforabigwell",
			Component_volumes: map[string]wunit.Volume{"water": wunit.NewVolume(7000, "ul")},
			Plate_types:       []*wtype.LHPlate{DSW96, DWST12},
			Weight_constraint: weightConstraint,
			Expected:          map[string]map[string]int{"water": {"DWST12": 1}},
		},
		{
			Name:              "2things2plates",
			Component_volumes: map[string]wunit.Volume{"water": wunit.NewVolume(25155, "ul"), "milk": wunit.NewVolume(300, "ul")},
			Plate_types:       []*wtype.LHPlate{pcrplateSkirted, DWST12},
			Weight_constraint: weightConstraint,
			Expected:          map[string]map[string]int{"water": {"DWST12": 3}, "milk": {"pcrplate_skirted": 2}},
		},
	}

	for _, testCase := range testCases {
		doTheTest := func(t *testing.T) {
			fmt.Println("AHMA DOIN MA TEST ", testCase.Name)
			results, err := choose_plate_assignments(testCase.Component_volumes, testCase.Plate_types, testCase.Weight_constraint)

			_, fail := testCase.Expected["FAIL"]

			if fail {
				// routine should fail

				if err == nil {
					t.Errorf("Expected failure, got %v", summarize(results))
				}

			} else {
				// routine should pass

				if err != nil {
					t.Errorf("Expected %s got error %v", fmtIt(testCase.Expected), err.Error())
				}

				if !compare(results, testCase.Expected) {
					t.Errorf("Expected: %s got %s", fmtIt(testCase.Expected), fmtIt(results))
				}
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
