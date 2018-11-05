package liquidhandling

import (
	"context"
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

func configure_request_quitebig(ctx context.Context, rq *LHRequest) {
	water := GetComponentForTest(ctx, "water", wunit.NewVolume(5000.0, "ul"))
	mmx := GetComponentForTest(ctx, "mastermix_sapI", wunit.NewVolume(5000.0, "ul"))
	part := GetComponentForTest(ctx, "dna", wunit.NewVolume(5000.0, "ul"))

	for k := 0; k < 130; k++ {
		ins := wtype.NewLHMixInstruction()
		ws := mixer.Sample(water, wunit.NewVolume(21.0, "ul"))
		mmxs := mixer.Sample(mmx, wunit.NewVolume(21.0, "ul"))
		ps := mixer.Sample(part, wunit.NewVolume(1.0, "ul"))

		ins.AddInput(ws)
		ins.AddInput(mmxs)
		ins.AddInput(ps)
		ins.AddOutput(GetComponentForTest(ctx, "water", wunit.NewVolume(43.0, "ul")))
		ins.Outputs[0].CName = fmt.Sprintf("DANGER_MIX_%d", k)
		ins.SetGeneration(k + 1)
		rq.Add_instruction(ins)
	}
}

func GetComponentForTest(ctx context.Context, name string, vol wunit.Volume) *wtype.Liquid {
	c, err := inventory.NewComponent(ctx, name)
	if err != nil {
		panic(err)
	}
	c.ID = wtype.GetUUID()
	c.SetVolume(vol)
	return c
}

func GetItHere(ctx context.Context) (*Liquidhandler, *LHRequest, error) {
	lh := GetLiquidHandlerForTest(ctx)
	rq := GetLHRequestForTest()
	configure_request_quitebig(ctx, rq)
	rq.InputPlatetypes = append(rq.InputPlatetypes, GetPlateForTest())
	rq.OutputPlatetypes = append(rq.OutputPlatetypes, GetPlateForTest())

	err := lh.Plan(ctx, rq)
	if err != nil {
		return nil, nil, err
	}
	return lh, rq, nil
}

func whereISthatplate(name string, robot *liquidhandling.LHProperties) string {
	for pos, plt := range robot.Plates {
		if itshere(name, plt) {
			return pos
		}
	}

	return "notheremate"
}

func itshere(name string, plate *wtype.Plate) bool {
	for _, w := range plate.Wellcoords {
		if w.IsEmpty() {
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

	ctx := testinventory.NewContext(context.Background())
	lastLH, _, err := GetItHere(ctx)
	if err != nil {
		t.Fatalf("Got an error planning with no inputs: %s", err)
	}

	for i := 0; i < 10; i++ {
		lh, _, err := GetItHere(ctx)
		if err != nil {
			t.Fatalf("Got an error planning with no inputs: %s", err)
		}

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
