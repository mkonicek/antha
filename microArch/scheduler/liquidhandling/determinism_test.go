package liquidhandling

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"testing"
)

func configure_request_quitebig(rq *LHRequest) {
	water := GetComponentForTest("water", wunit.NewVolume(5000.0, "ul"))
	mmx := GetComponentForTest("mastermix_sapI", wunit.NewVolume(5000.0, "ul"))
	part := GetComponentForTest("dna", wunit.NewVolume(5000.0, "ul"))

	for k := 0; k < 130; k++ {
		ins := wtype.NewLHInstruction()
		ws := mixer.Sample(water, wunit.NewVolume(21.0, "ul"))
		mmxs := mixer.Sample(mmx, wunit.NewVolume(21.0, "ul"))
		ps := mixer.Sample(part, wunit.NewVolume(1.0, "ul"))

		ins.AddComponent(ws)
		ins.AddComponent(mmxs)
		ins.AddComponent(ps)
		ins.AddProduct(GetComponentForTest("water", wunit.NewVolume(43.0, "ul")))
		ins.Result.CName = fmt.Sprintf("DANGER_MIX_%d", k)
		ins.SetGeneration(k + 1)
		rq.Add_instruction(ins)
	}
}

func GetItHere(t *testing.T) (*Liquidhandler, *LHRequest) {
	lh := GetLiquidHandlerForTest()
	rq := GetLHRequestForTest()
	configure_request_quitebig(rq)
	rq.Input_platetypes = append(rq.Input_platetypes, GetPlateForTest())
	rq.Output_platetypes = append(rq.Output_platetypes, GetPlateForTest())

	rq.ConfigureYourself()

	err := lh.Plan(rq)

	if err != nil {
		t.Fatal(fmt.Sprint("Got an error planning with no inputs: ", err))
	}
	return lh, rq
}

func whereISthatplate(name string, robot *liquidhandling.LHProperties) string {
	for pos, plt := range robot.Plates {
		if itshere(name, plt) {
			return pos
		}
	}

	return "notheremate"
}

func itshere(name string, plate *wtype.LHPlate) bool {
	for _, w := range plate.Wellcoords {
		if w.Empty() {
			continue
		}
		if w.WContents.CName == name {
			return true
		}
	}

	return false
}

func TestLayoutDeterminism(t *testing.T) {
	t.Skip() // pending final changes
	lastLH, _ := GetItHere(t)

	for i := 0; i < 10; i++ {
		lh, _ := GetItHere(t)

		was := whereISthatplate("DANGER_MIX_0", lastLH.FinalProperties)

		if was == "notheremate" {
			t.Fatal("BIG, WEIRD ERROR! No plate found in before time")
		}

		is := whereISthatplate("DANGER_MIX_0", lh.FinalProperties)

		if is == "notheremate" {
			t.Fatal("BIG, WEIRD ERROR! No plate found in after time")
		}

		if was != is {
			t.Fatal(fmt.Sprintf("Think again, boyo - your layout ain't deterministic nohow %s =/= %s", was, is))
		}
	}
}
