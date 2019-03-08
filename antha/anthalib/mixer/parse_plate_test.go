package mixer

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory/testinventory"
	"github.com/pkg/errors"
)

func nonEmpty(m map[string]*wtype.LHWell) map[string]*wtype.Liquid {
	r := make(map[string]*wtype.Liquid)
	for addr, c := range m {
		if c.WContents.IsZero() {
			continue
		}
		r[addr] = c.WContents
	}
	return r
}

func getComponentsFromPlate(plate *wtype.Plate) []*wtype.Liquid {

	var components []*wtype.Liquid
	allWellPositions := plate.AllWellPositions(false)

	for _, wellcontents := range allWellPositions {

		if !plate.WellMap()[wellcontents].IsEmpty() {

			component := plate.WellMap()[wellcontents].WContents
			components = append(components, component)

		}
	}
	return components
}

func allComponentsHaveWellLocation(plate *wtype.Plate) error {
	components := getComponentsFromPlate(plate)
	var errs []string
	for _, component := range components {
		if len(component.WellLocation()) == 0 {
			errs = append(errs, fmt.Errorf("no well location for %s after returning components from plate", component.Name()).Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, "\n"))
	}
	return nil
}

func samePlate(a, b *wtype.Plate) error {
	if a.Type != b.Type {
		return fmt.Errorf("different types %q != %q", a.Type, b.Type)
	}
	compsA := nonEmpty(a.Wellcoords)
	compsB := nonEmpty(b.Wellcoords)

	if numA, numB := len(compsA), len(compsB); numA != numB {
		return fmt.Errorf("different number of non-empty wells %d != %d", numA, numB)
	}

	for addr, compA := range compsA {

		compB, ok := compsB[addr]
		if !ok {
			return fmt.Errorf("missing component in well %q", addr)
		}

		volA, volB := compA.Vol, compB.Vol
		if volA != volB {
			return fmt.Errorf("different volume in well %q: %f != %f", addr, volA, volB)
		}
		vunitA, vunitB := compA.Vunit, compB.Vunit
		if vunitA != vunitB && volA != 0.0 {
			return fmt.Errorf("different volume unit in well %q: %s != %s", addr, vunitA, vunitB)
		}
		concA, concB := compA.Conc, compB.Conc
		if concA != concB {
			return fmt.Errorf("different concentration in well %q: expected: %f; found: %f", addr, concA, concB)
		}
		cunitA, cunitB := compA.Cunit, compB.Cunit
		if cunitA != cunitB && concA != 0.0 {
			return fmt.Errorf("different concetration unit in well %q: expected: %s; found: %s", addr, cunitA, cunitB)
		}

		if err := wtype.EqualLists(compA.SubComponents, compB.SubComponents); err != nil {
			return errors.Errorf("%s: %+v != %+v", err.Error(), compA.SubComponents, compB.SubComponents)
		}
	}

	return nil
}

func containsInvalidCharWarning(warnings []string) bool {
	for _, v := range warnings {
		if strings.Contains(v, "contains an invalid character \"+\"") {
			return true
		}
	}

	return false
}

func TestParsePlateWithValidation(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	file := []byte(
		`
pcrplate_with_cooler,
A1,water+soil,water,50.0,ul,0,g/l,
A4,tea,water,50.0,ul,0,g/l,
A5,milk,water,100.0,ul,0,g/l,
`)
	r, err := parsePlateCSVWithValidationConfig(ctx, bytes.NewBuffer(file), DefaultValidationConfig())

	if err != nil {
		t.Errorf("Failed to parse plate: %s ", err.Error())
	}
	if !containsInvalidCharWarning(r.Warnings) {
		t.Errorf("Default validation config must forbid + signs in component names")
	}
	r, err = parsePlateCSVWithValidationConfig(ctx, bytes.NewBuffer(file), PermissiveValidationConfig())

	if err != nil {
		t.Errorf("Failed to parse plate: %s ", err.Error())
	}

	if containsInvalidCharWarning(r.Warnings) {
		t.Errorf("Permissive validation config must allow + signs in component names")
	}
}

func TestParsePlate(t *testing.T) {
	type testCase struct {
		File              []byte
		Expected          *wtype.Plate
		NoWarnings        bool
		ReplacementConfig ValidationConfig
	}

	ctx := testinventory.NewContext(context.Background())

	// Read external file with carriage returns for that specific test.
	fileCarriage, err := ioutil.ReadFile("test_carriage.csv")
	if err != nil {
		t.Errorf("Failed to read test_carriage.csv: %s ", err.Error())
	}

	suite := []testCase{
		{
			File: []byte(
				`
pcrplate_with_cooler,
A1,water,water,50.0,ul,0,g/l,
A4,tea,water,50.0,ul,10.0,mMol/l,
A5,milk,water,100.0,ul,10.0,g/l,
A6,,,0,ul,0,g/l,
`),
			NoWarnings: false,
			Expected: &wtype.Plate{
				Type: "pcrplate_with_cooler",
				Wellcoords: map[string]*wtype.LHWell{
					"A1": {
						WContents: &wtype.Liquid{
							CName: "water",
							Type:  wtype.LTWater,
							Vol:   50.0,
							Vunit: "ul",
							Conc:  0.0,
							Cunit: "g/l",
						},
					},
					"A4": {
						WContents: &wtype.Liquid{
							CName: "tea",
							Type:  wtype.LTWater,
							Vol:   50.0,
							Vunit: "ul",
							Conc:  10.0,
							Cunit: "mMol/l",
						},
					},
					"A5": {
						WContents: &wtype.Liquid{
							CName: "milk",
							Type:  wtype.LTWater,
							Vol:   100.0,
							Vunit: "ul",
							Conc:  10.0,
							Cunit: "g/l",
						},
					},
				},
			},
		},
		{
			File: []byte(
				`
pcrplate_with_cooler, afternoon tea tray, LiquidType, Vol,Vol Unit,Conc,Conc Unit, SubComponents
A1,water,water,50.0,ul,0,g/l,
A4,tea,water,50.0,ul,10.0,mMol/l,tea leaves: ,5g/l,sugar:,1X,
A5,milk,water,100.0,ul,10.0,g/l,
A6,,,0,ul,0,g/l,
`),
			NoWarnings: false,
			Expected: &wtype.Plate{
				Type: "pcrplate_with_cooler",
				Wellcoords: map[string]*wtype.LHWell{
					"A1": {
						WContents: &wtype.Liquid{
							CName: "water",
							Type:  wtype.LTWater,
							Vol:   50.0,
							Vunit: "ul",
							Conc:  0.0,
							Cunit: "g/l",
						},
					},
					"A4": {
						WContents: &wtype.Liquid{
							CName: "tea",
							Type:  wtype.LTWater,
							Vol:   50.0,
							Vunit: "ul",
							Conc:  10.0,
							Cunit: "mMol/l",
							SubComponents: wtype.ComponentList{
								Components: map[string]wunit.Concentration{
									"tea leaves": wunit.NewConcentration(5.0, "g/L"),
									"sugar":      wunit.NewConcentration(1.0, "X"),
								},
							},
						},
					},
					"A5": {
						WContents: &wtype.Liquid{
							CName: "milk",
							Type:  wtype.LTWater,
							Vol:   100.0,
							Vunit: "ul",
							Conc:  10.0,
							Cunit: "g/l",
						},
					},
				},
			},
		},
		{
			File: []byte(
				`
pcrplate_with_cooler, afternoon tea tray, LiquidType, Vol,Vol Unit,Conc,Conc Unit, , ,SubComponents,
A1,water,water,50.0,ul,0,g/l,
A4,tea,water,50.0,ul,10.0,mMol/l, some random user text, ,tea leaves: ,5g/l,sugar:,1X,
A5,milk,water,100.0,ul,10.0,g/l,
A6,,,0,ul,0,g/l,
`),
			NoWarnings: false,
			Expected: &wtype.Plate{
				Type: "pcrplate_with_cooler",
				Wellcoords: map[string]*wtype.LHWell{
					"A1": {
						WContents: &wtype.Liquid{
							CName: "water",
							Type:  wtype.LTWater,
							Vol:   50.0,
							Vunit: "ul",
							Conc:  0.0,
							Cunit: "g/l",
						},
					},
					"A4": {
						WContents: &wtype.Liquid{
							CName: "tea",
							Type:  wtype.LTWater,
							Vol:   50.0,
							Vunit: "ul",
							Conc:  10.0,
							Cunit: "mMol/l",
							SubComponents: wtype.ComponentList{
								Components: map[string]wunit.Concentration{
									"tea leaves": wunit.NewConcentration(5.0, "g/L"),
									"sugar":      wunit.NewConcentration(1.0, "X"),
								},
							},
						},
					},
					"A5": {
						WContents: &wtype.Liquid{
							CName: "milk",
							Type:  wtype.LTWater,
							Vol:   100.0,
							Vunit: "ul",
							Conc:  10.0,
							Cunit: "g/l",
						},
					},
				},
			},
		},
		{
			File: []byte(
				`
pcrplate_with_cooler,
A1,water,water,50.0,ul,0,g/l,
A4,tea,water,50.0,ul,10.0,mMol/l,
A5,milk,water,100.0,ul,10.0,g/l,
A6,,,0,ul,0,g/l,
`),
			NoWarnings: false,
			ReplacementConfig: ValidationConfig{
				replaceField: map[string]string{
					plateTypeReplacementKey: "pcrplate_skirted",
				},
			},
			Expected: &wtype.Plate{
				Type: "pcrplate_skirted",
				Wellcoords: map[string]*wtype.LHWell{
					"A1": {
						WContents: &wtype.Liquid{
							CName: "water",
							Type:  wtype.LTWater,
							Vol:   50.0,
							Vunit: "ul",
							Conc:  0.0,
							Cunit: "g/l",
						},
					},
					"A4": {
						WContents: &wtype.Liquid{
							CName: "tea",
							Type:  wtype.LTWater,
							Vol:   50.0,
							Vunit: "ul",
							Conc:  10.0,
							Cunit: "mMol/l",
						},
					},
					"A5": {
						WContents: &wtype.Liquid{
							CName: "milk",
							Type:  wtype.LTWater,
							Vol:   100.0,
							Vunit: "ul",
							Conc:  10.0,
							Cunit: "g/l",
						},
					},
				},
			},
		},
		{
			File: []byte(
				`
pcrplate_skirted_riser40,Input_plate_1,LiquidType,Vol,Vol Unit,Conc,Conc Unit
A1,water,water,140.5,ul,0,mg/l
C1,neb5compcells,culture,20.5,ul,0,ng/ul
`),
			NoWarnings: true,
			Expected: &wtype.Plate{
				Type: "pcrplate_skirted_riser40",
				Wellcoords: map[string]*wtype.LHWell{
					"A1": {
						WContents: &wtype.Liquid{
							CName: "water",
							Type:  wtype.LTWater,
							Vol:   140.5,
							Vunit: "ul",
							Conc:  0,
							Cunit: "mg/l",
						},
					},
					"C1": {
						WContents: &wtype.Liquid{
							CName: "neb5compcells",
							Type:  wtype.LTCulture,
							Vol:   20.5,
							Vunit: "ul",
							Conc:  0,
							Cunit: "mg/l",
						},
					},
				},
			},
		},
		{
			// This is to test carriage returns.
			File:       fileCarriage,
			NoWarnings: false,
			Expected: &wtype.Plate{
				Type: "pcrplate_with_cooler",
				Wellcoords: map[string]*wtype.LHWell{
					"A1": {
						WContents: &wtype.Liquid{
							CName: "water",
							Type:  wtype.LTWater,
							Vol:   50.0,
							Vunit: "ul",
							Conc:  0.0,
							Cunit: "g/l",
						},
					},
					"A4": {
						WContents: &wtype.Liquid{
							CName: "tea",
							Type:  wtype.LTWater,
							Vol:   50.0,
							Vunit: "ul",
							Conc:  10.0,
							Cunit: "mMol/l",
						},
					},
					"A5": {
						WContents: &wtype.Liquid{
							CName: "milk",
							Type:  wtype.LTWater,
							Vol:   100.0,
							Vunit: "ul",
							Conc:  10.0,
							Cunit: "g/l",
						},
					},
				},
			},
		},
		{
			File: []byte(
				`
pcrplate_skirted_riser40,Input_plate_1,LiquidType,Vol,Vol Unit,Conc,Conc Unit
A1,water,randomType,140.5,ul,0,mg/l
C1,neb5compcells,culture,20.5,ul,0,ng/ul
`),
			NoWarnings: false,
			Expected: &wtype.Plate{
				Type: "pcrplate_skirted_riser40",
				Wellcoords: map[string]*wtype.LHWell{
					"A1": {
						WContents: &wtype.Liquid{
							CName: "water",
							Type:  wtype.LiquidType("randomType"),
							Vol:   140.5,
							Vunit: "ul",
							Conc:  0,
							Cunit: "mg/l",
						},
					},
					"C1": {
						WContents: &wtype.Liquid{
							CName: "neb5compcells",
							Type:  wtype.LTCulture,
							Vol:   20.5,
							Vunit: "ul",
							Conc:  0,
							Cunit: "mg/l",
						},
					},
				},
			},
		},
		{
			File: []byte(
				`
pcrplate_skirted_riser40,Input_plate_1,LiquidType,Vol,Vol Unit,Conc,Conc Unit
A1,water,randomType_modified_1,140.5,ul,0,mg/l
C1,neb5compcells,culture,20.5,ul,0,ng/ul
`),
			NoWarnings: false,
			Expected: &wtype.Plate{
				Type: "pcrplate_skirted_riser40",
				Wellcoords: map[string]*wtype.LHWell{
					"A1": {
						WContents: &wtype.Liquid{
							CName: "water",
							Type:  wtype.LiquidType("randomType"),
							Vol:   140.5,
							Vunit: "ul",
							Conc:  0,
							Cunit: "mg/l",
						},
					},
					"C1": {
						WContents: &wtype.Liquid{
							CName: "neb5compcells",
							Type:  wtype.LTCulture,
							Vol:   20.5,
							Vunit: "ul",
							Conc:  0,
							Cunit: "mg/l",
						},
					},
				},
			},
		},
	}

	for i, tc := range suite {
		p, err := ParsePlateCSV(ctx, bytes.NewBuffer(tc.File), tc.ReplacementConfig)
		if err != nil {
			t.Error(err)
		}
		if err := samePlate(tc.Expected, p.Plate); err != nil {
			t.Error(fmt.Sprintf("error in test %d: %s", i, err))
		}
		if tc.NoWarnings && len(p.Warnings) != 0 {
			t.Errorf("found warnings: %s", p.Warnings)
		}

		if err := allComponentsHaveWellLocation(p.Plate); err != nil {
			t.Error(err.Error())
		}
	}
}

func TestParsePlateOverfilled(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	file := []byte(
		`
pcrplate_with_cooler,
A1,water,water,50.0,ul,0,g/l,
A4,tea,water,50.0,ul,0,g/l,
A5,milk,water,500.0,ul,0,g/l,
`)
	_, err := parsePlateCSVWithValidationConfig(ctx, bytes.NewBuffer(file), DefaultValidationConfig())

	if err == nil {
		t.Error("Overfull well A5 failed to generate error")
	}
}

func TestUnModifyTypeName(t *testing.T) {

	type modifiedPolicyTest struct {
		Original, ExpectedProduct string
	}

	var tests = []modifiedPolicyTest{
		{
			Original:        "water_modified_1",
			ExpectedProduct: "water",
		},
		{
			Original:        "water_modified_2_modified_1",
			ExpectedProduct: "water",
		},
		{
			Original:        "water",
			ExpectedProduct: "water",
		},
		{
			Original:        "_modified_2",
			ExpectedProduct: "",
		},
	}

	for _, test := range tests {
		result := unModifyTypeName(test.Original)
		if result != test.ExpectedProduct {
			t.Error(
				"Error removing modified suffix from policy: \n",
				"For: ", test.Original, "\n",
				"Expected Product: ", test.ExpectedProduct, "\n",
				"Got: ", result, "\n",
			)
		}
	}

}
