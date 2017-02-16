package mixer

import (
	"bytes"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/factory"
)

func makeTestPlate(in *wtype.LHPlate) *wtype.LHPlate {
	out := factory.GetPlateByType(in.Type)
	out.PlateName = in.PlateName
	for coord, well := range in.Wellcoords {
		out.WellAt(wtype.MakeWellCoordsA1(coord)).Add(well.WContents)
	}
	return out
}

func TestMarshalPlateCSV(t *testing.T) {
	type testCase struct {
		Plate    *wtype.LHPlate
		Expected []byte
	}

	suite := []testCase{
		testCase{
			Expected: []byte(
				`
pcrplate_with_cooler,Input_plate_1,LiquidType,Vol,Vol Unit,Conc,Conc Unit
A1,water,water,50,ul,0,
A4,tea,water,50,ul,0,
A5,milk,water,100,ul,0,
`),
			Plate: makeTestPlate(&wtype.LHPlate{
				PlateName: "Input_plate_1",
				Type:      "pcrplate_with_cooler",
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
			}),
		},
		testCase{
			Expected: []byte(
				`
pcrplate_skirted_riser40,Input_plate_1,LiquidType,Vol,Vol Unit,Conc,Conc Unit
A1,water,water,140.5,ul,0,
C1,neb5compcells,culture,20.5,ul,0,
`),
			Plate: makeTestPlate(&wtype.LHPlate{
				PlateName: "Input_plate_1",
				Type:      "pcrplate_skirted_riser40",
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
