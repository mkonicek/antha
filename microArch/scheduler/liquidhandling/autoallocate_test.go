package liquidhandling

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/factory"
	"testing"
)

func TestInputSampleAutoAllocate(t *testing.T) {
	rbt := makeGilson()
	rq := NewLHRequest()

	cmp1 := factory.GetComponentByType("water")
	cmp2 := factory.GetComponentByType("dna_part")

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

	pl := factory.GetPlateByType("pcrplate_skirted_riser20")

	rq.Input_platetypes = append(rq.Input_platetypes, pl)

	rq.ConfigureYourself()

	lh := Init(rbt)

	lh.Plan(rq)

	expected := make(map[string]float64)

	expected["dna_part"] = 30.5
	expected["water"] = 55.5

	testSetup(rbt, expected, t)
}

func testSetup(rbt *liquidhandling.LHProperties, expected map[string]float64, t *testing.T) {
	for _, p := range rbt.Plates {
		for _, w := range p.Wellcoords {
			if !w.Empty() {
				v, ok := expected[w.WContents.CName]

				if !ok {
					t.Fatal(fmt.Sprint("ERROR: unexpected component in plating area: ", w.WContents.CName))
				}

				if v != w.WContents.Vol {
					t.Fatal(fmt.Sprint("ERROR: Volume of component ", w.WContents.CName, " was ", w.WContents.Vol, " should be ", v))
				}

				delete(expected, w.WContents.CName)
			}
		}
	}

	if len(expected) != 0 {
		t.Fatal(fmt.Sprint("ERROR: Expected components remaining: ", expected))
	}

}
func TestInPlaceAutoAllocate(t *testing.T) {
	rbt := makeGilson()
	rq := NewLHRequest()

	cmp1 := factory.GetComponentByType("water")
	cmp2 := factory.GetComponentByType("dna_part")

	cmp1.Vol = 100.0
	cmp2.Vol = 50.0

	mo := mixer.MixOptions{
		Components: []*wtype.LHComponent{cmp1, cmp2},
		PlateType:  "pcrplate_skirted_riser20",
		Address:    "A1",
		PlateNum:   1,
	}

	ins := mixer.GenericMix(mo)

	rq.LHInstructions[ins.ID] = ins

	pl := factory.GetPlateByType("pcrplate_skirted_riser20")

	rq.Input_platetypes = append(rq.Input_platetypes, pl)

	rq.ConfigureYourself()

	lh := Init(rbt)

	lh.Plan(rq)

	expected := make(map[string]float64)

	expected["dna_part"] = 55.5
	expected["water"] = 100.0

	testSetup(rbt, expected, t)

}
