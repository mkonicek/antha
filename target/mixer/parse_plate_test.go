package mixer

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func nonEmpty(m map[string]*wtype.LHWell) map[string]*wtype.LHComponent {
	r := make(map[string]*wtype.LHComponent)
	for addr, c := range m {
		if c.WContents.IsZero() {
			continue
		}
		r[addr] = c.WContents
	}
	return r
}

func samePlate(a, b *wtype.LHPlate) error {
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
			return fmt.Errorf("different volume in well %q: %d != %d", addr, volA, volB)
		}
		vunitA, vunitB := compA.Vunit, compB.Vunit
		if vunitA != vunitB && volA != 0.0 {
			return fmt.Errorf("different volume unit in well %q: %d != %d", addr, vunitA, vunitB)
		}
		concA, concB := compA.Conc, compB.Conc
		if concA != concB {
			return fmt.Errorf("different concentration in well %q: %d != %d", addr, concA, concB)
		}
		cunitA, cunitB := compA.Cunit, compB.Cunit
		if cunitA != cunitB && concA != 0.0 {
			return fmt.Errorf("different concetration unit in well %q: %d != %d", addr, cunitA, cunitB)
		}
	}

	return nil
}

func TestParsePlate(t *testing.T) {
	type testCase struct {
		File     []byte
		Expected *wtype.LHPlate
	}

	suite := []testCase{
		testCase{
			File: []byte(
				`
pcrplate_with_cooler,
A1,water,water,50.0,ul,
A4,tea,water,50.0,ul,
A5,milk,water,100.0,ul,
`),
			Expected: &wtype.LHPlate{
				Type: "pcrplate_with_cooler",
				Wellcoords: map[string]*wtype.LHWell{
					"A1": &wtype.LHWell{
						WContents: &wtype.LHComponent{
							CName: "water",
							Type:  wtype.LTWater,
							Vol:   50.0,
							Vunit: "ul",
						},
					},
					"A4": &wtype.LHWell{
						WContents: &wtype.LHComponent{
							CName: "tea",
							Type:  wtype.LTWater,
							Vol:   50.0,
							Vunit: "ul",
						},
					},
					"A5": &wtype.LHWell{
						WContents: &wtype.LHComponent{
							CName: "milk",
							Type:  wtype.LTWater,
							Vol:   100.0,
							Vunit: "ul",
						},
					},
				},
			},
		},
		testCase{
			File: []byte(
				`
pcrplate_skirted_riser40,Input_plate_1,LiquidType,Vol,Vol Unit,Conc,Conc Unit
A1,water,water,140.5,ul,0,mg/l
C1,neb5compcells,culture,20.5,ul,0,mg/l
`),
			Expected: &wtype.LHPlate{
				Type: "pcrplate_skirted_riser40",
				Wellcoords: map[string]*wtype.LHWell{
					"A1": &wtype.LHWell{
						WContents: &wtype.LHComponent{
							CName: "water",
							Type:  wtype.LTWater,
							Vol:   140.5,
							Vunit: "ul",
						},
					},
					"C1": &wtype.LHWell{
						WContents: &wtype.LHComponent{
							CName: "neb5compcells",
							Type:  wtype.LTCulture,
							Vol:   20.5,
							Vunit: "ul",
						},
					},
				},
			},
		},
	}

	for _, tc := range suite {
		p, err := ParsePlateCSV(bytes.NewBuffer(tc.File))
		if err != nil {
			t.Error(err)
		}
		if err := samePlate(tc.Expected, p.Plate); err != nil {
			t.Error(err)
		}
	}
}
