package liquidhandling

import (
	"context"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory/testinventory"
)

func TestParallelSetGeneration(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	// lazy way to get pre-prepared instructions
	tb, dstp := getTransferBlock2Component(ctx) // defined in transferblock_test

	ins := make([]*wtype.LHInstruction, 0, len(tb.Inss)-1)

	for i := 0; i < len(tb.Inss); i++ {
		// make one hole
		if i == 4 {
			continue
		}
		ins = append(ins, tb.Inss[i])
	}

	tb.Inss = ins

	rbt := getTestRobot(ctx, dstp, "pcrplate_skirted_riser40")

	// allow independent multichannel activity
	headsLoaded := rbt.GetLoadedHeads()
	headsLoaded[0].Params.Independent = true

	pol, err := wtype.GetLHPolicyForTest()
	if err != nil {
		t.Error(err)
	}

	// allow multi

	pol.Policies["water"]["CAN_MULTI"] = true

	//get_parallel_sets_head(ctx context.Context, head *wtype.LHHead, ins []*wtype.LHInstruction)

	insIds, err := get_parallel_sets_head(ctx, headsLoaded[0], ins)

	if err != nil {
		t.Errorf("Parallel set generation error: %s\n", err)
	}

	if len(insIds[0]) != 8 {
		t.Errorf("Should have 8 insIDs, instead have %d", len(insIds[0]))
	}

	// insIds[0][4] should be ""

	if insIds[0][4] != "" {
		t.Errorf("InsIds[0][4] should be \"\", it isn't, it's: %s", insIds[0][4])
	}
}
