package liquidhandling

import (
	"context"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"
	"reflect"
	"testing"
)

type initFinalCmp struct {
	CNameI string
	CNameF string
	VolI   float64
	VolF   float64
}

func (ifc initFinalCmp) IsZero() bool {
	v := initFinalCmp{}
	return reflect.DeepEqual(v, ifc)
}

func TestStateBeforeAfter(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())
	rq := makeRequest()
	lh := makeLiquidhandler(ctx)
	cmp1, err := inventory.NewComponent(ctx, inventory.WaterType)
	if err != nil {
		t.Fatal(err)
	}
	cmp2, err := inventory.NewComponent(ctx, "dna_part")
	if err != nil {
		t.Fatal(err)
	}

	s1 := mixer.Sample(cmp1, wunit.NewVolume(50.0, "ul"))
	s2 := mixer.Sample(cmp2, wunit.NewVolume(25.0, "ul"))

	mo := mixer.MixOptions{
		Components: []*wtype.LHComponent{s1, s2},
		PlateType:  "pcrplate_skirted_riser20",
		Address:    "A1",
		PlateNum:   1,
	}

	ins := mixer.GenericMix(mo)

	rq.LHInstructions[ins.ID] = ins

	pl, err := inventory.NewPlate(ctx, "pcrplate_skirted_riser20")
	if err != nil {
		t.Fatal(err)
	}

	rq.Input_platetypes = append(rq.Input_platetypes, pl)

	rq.ConfigureYourself()

	lh.Plan(ctx, rq)

	/*
		expectedInitial := make(map[string]float64)
		expectedInitial["dna_part"] = 30.5
		expectedInitial["water"] = 55.5

		expectedFinal := make(map[string]float64)
		expectedFinal["water+dna_part"] = 75.0
		expectedFinal["water"] = 5.0
		expectedFinal["dna_part"] = 5.0
	*/

	expected := make(map[string][]initFinalCmp)

	expected["dna_part"] = []initFinalCmp{initFinalCmp{CNameI: "dna_part", CNameF: "dna_part", VolI: 30.5, VolF: 5.0}}
	expected["water"] = []initFinalCmp{initFinalCmp{CNameI: "water", CNameF: "water", VolI: 55.5, VolF: 5.0}}

	expected["water+dna_part"] = []initFinalCmp{initFinalCmp{CNameI: "", CNameF: "water+dna_part", VolI: 0.0, VolF: 75.0}}

	compareInitFinalStates(t, lh, expected)
}

// all means just one or the whole lot
func del(a initFinalCmp, ar []initFinalCmp, all bool) []initFinalCmp {
	ar2 := make([]initFinalCmp, 0, len(ar)-1)
	d := false
	for _, b := range ar {
		if !reflect.DeepEqual(a, b) || (d && !all) {
			ar2 = append(ar2, b)
			d = true
		}
	}

	return ar2
}

/*

type initFinalCmp struct {
	CNameI string
	CNameF string
	VolI   float64
	VolF   float64
}

*/
func findWells(wi, wf *wtype.LHWell, ar []initFinalCmp) initFinalCmp {
	ifc := initFinalCmp{CNameI: wi.WContents.CName, CNameF: wf.WContents.CName, VolI: wi.WContents.Vol, VolF: wf.WContents.Vol}

	return findIFC(ifc, ar)
}

func findIFC(ifc initFinalCmp, ar []initFinalCmp) initFinalCmp {
	r := initFinalCmp{}

	for _, ifc2 := range ar {
		if reflect.DeepEqual(ifc, ifc2) {
			r = ifc2
			break
		}
	}

	return r
}

func compareInitFinalStates(t *testing.T, lh *Liquidhandler, expected map[string][]initFinalCmp) {
	xor := func(a, b bool) bool {
		return (a && !b) || (b && !a)
	}
	for _, pos := range lh.Properties.InputSearchPreferences() {
		p, ok := lh.Properties.Plates[pos]
		p2, ok2 := lh.FinalProperties.Plates[pos]

		if xor(ok, ok2) {
			t.Errorf("Plates moving in simple liquid handling plan")
		}

		if ok {
			// find each component and make sure it has stayed in the same place
			for _, crd := range p.AllWellPositions(false) {
				w := p.Wellcoords[crd]
				w2 := p2.Wellcoords[crd]

				v, ok3 := expected[w2.WContents.CName]

				if ok3 {
					ifc := findWells(w, w2, v)

					if ifc.IsZero() {
						t.Errorf("Extraneous components in before / after: %s %f %s %f", w.WContents.CName, w.WContents.Vol, w2.WContents.CName, w2.WContents.Vol)
					}

					// good, delete this now

					expected[w2.WContents.CName] = del(ifc, v, false)
				}
			}
		}
	}

	// is anything left in the expected pile?

	for k, v := range expected {
		if len(v) != 0 {
			t.Errorf("Unmatched components of type %s : %d total", k, len(v))
			for _, vv := range v {
				t.Errorf("%+v", vv)
			}
		}
	}
}

func makeLiquidhandler(ctx context.Context) *Liquidhandler {
	rbt := makeGilson(ctx)
	lh := Init(rbt)
	return lh
}

func makeRequest() *LHRequest {
	rq := NewLHRequest()
	return rq
}
