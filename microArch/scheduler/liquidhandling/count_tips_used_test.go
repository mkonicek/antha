package liquidhandling

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory/testinventory"
	"reflect"
	"testing"
)

func TestTipCounting(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	lh := GetLiquidHandlerForTest(ctx)
	lh.ExecutionPlanner = ExecutionPlanner3
	rq := GetLHRequestForTest()
	configure_request_simple(ctx, rq)
	rq.InputPlatetypes = append(rq.InputPlatetypes, GetPlateForTest())
	rq.Output_platetypes = append(rq.Output_platetypes, GetPlateForTest())

	err := lh.Plan(ctx, rq)

	if err != nil {
		t.Fatal(fmt.Sprint("Got planning error: ", err))
	}

	// [{DFL10 Tip Rack (PIPETMAX 8x20) 27 1}]

	expected := []wtype.TipEstimate{{TipType: "DFL10 Tip Rack (PIPETMAX 8x20)", NTips: 27, NTipBoxes: 1}}

	if !reflect.DeepEqual(expected, rq.TipsUsed) {
		t.Errorf("Expected %v Got %v", expected, rq.TipsUsed)
	}
}
