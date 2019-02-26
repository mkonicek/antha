package mixer

import (
	"bytes"
	"context"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"
)

func makeTestPlate(ctx context.Context, in *wtype.Plate) *wtype.Plate {
	out, err := inventory.NewPlate(ctx, in.Type)
	if err != nil {
		panic(err)
	}

	out.PlateName = in.PlateName
	for coord, well := range in.Wellcoords {
		if w, ok := out.WellAt(wtype.MakeWellCoordsA1(coord)); ok {
			err := w.AddComponent(well.WContents)
			if err != nil {
				panic(err)
			}
		} else {
			panic("Couldn't find well")
		}
	}
	return out
}

func TestMarshalPlateCSV(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	type testCase struct {
		Plate    *wtype.Plate
		Expected []byte
	}

	suite := []testCase{
		{
			Expected: []byte(
				`
pcrplate_with_cooler,Input_plate_1,LiquidType,Vol,Vol Unit,Conc,Conc Unit
A1,water,water,50,ul,0,g/l
A4,tea,water,50,ul,10,mMol/l
A5,milk,water,100,ul,10,g/l
`),
			Plate: makeTestPlate(ctx, &wtype.Plate{
				PlateName: "Input_plate_1",
				Type:      "pcrplate_with_cooler",
				Wellcoords: map[string]*wtype.LHWell{
					"A1": {
						WContents: &wtype.Liquid{
							CName: "water",
							Type:  wtype.LTWater,
							Vol:   50.0,
							Vunit: "ul",
							Conc:  0,
							Cunit: "g/l",
						},
					},
					"A4": {
						WContents: &wtype.Liquid{
							CName: "tea",
							Type:  wtype.LTWater,
							Vol:   50.0,
							Vunit: "ul",
							Conc:  10,
							Cunit: "mMol/l",
						},
					},
					"A5": {
						WContents: &wtype.Liquid{
							CName: "milk",
							Type:  wtype.LTWater,
							Vol:   100.0,
							Vunit: "ul",
							Conc:  10,
							Cunit: "g/l",
						},
					},
				},
			}),
		},
		{
			Expected: []byte(
				`
pcrplate_skirted_riser40,Input_plate_1,LiquidType,Vol,Vol Unit,Conc,Conc Unit
A1,water,water,140.5,ul,0,g/l
C1,neb5compcells,culture,20.5,ul,0,g/l
`),
			Plate: makeTestPlate(ctx, &wtype.Plate{
				PlateName: "Input_plate_1",
				Type:      "pcrplate_skirted_riser40",
				Wellcoords: map[string]*wtype.LHWell{
					"A1": {
						WContents: &wtype.Liquid{
							CName: "water",
							Type:  wtype.LTWater,
							Vol:   140.5,
							Vunit: "ul",
							Conc:  0,
							Cunit: "g/l",
						},
					},
					"C1": {
						WContents: &wtype.Liquid{
							CName: "neb5compcells",
							Type:  wtype.LTCulture,
							Vol:   20.5,
							Vunit: "ul",
							Conc:  0,
							Cunit: "g/l",
						},
					},
				},
			}),
		},
	}

	for _, tc := range suite {
		bs, err := MarshalPlateCSV(tc.Plate)
		if err != nil {
			t.Error(err)
		}
		if e, f := bytes.TrimSpace(tc.Expected), bytes.TrimSpace(bs); !bytes.Equal(e, f) {
			t.Errorf("expected:\n%s\nfound:\n%s\n", string(e), string(f))
		}
	}
}
