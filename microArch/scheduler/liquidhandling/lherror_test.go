package liquidhandling

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"
	"testing"
)

func GetLiquidHandlerForTest(ctx context.Context) *Liquidhandler {
	return Init(makeGilson(ctx))
}

func GetIndependentLiquidHandlerForTest(ctx context.Context) *Liquidhandler {
	gilson := makeGilson(ctx)
	for _, head := range gilson.Heads {
		head.Params.Independent = true
	}

	for _, adaptor := range gilson.Adaptors {
		adaptor.Params.Independent = true
	}

	return Init(gilson)
}

func GetLHRequestForTest() *LHRequest {
	req := NewLHRequest()
	return req
}

func TestNoInstructions(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())
	req := GetLHRequestForTest()
	lh := GetLiquidHandlerForTest(ctx)
	err := lh.MakeSolutions(ctx, req)

	if err.Error() != "9 (LH_ERR_OTHER) :  : Nil plan requested: no Mix Instructions present" {
		t.Fatal(fmt.Sprint("Unexpected error: ", err.Error()))
	}
}

func TestDeckSpace1(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	lh := GetLiquidHandlerForTest(ctx)

	for i := 0; i < len(lh.Properties.Preferences.Tipboxes); i++ {
		tb, err := inventory.NewTipbox(ctx, lh.Properties.Tips[0].Type)
		if err != nil {
			t.Fatal(err)
		}

		if err := lh.Properties.AddTipBox(tb); err != nil {
			t.Fatalf("Should be able to fill deck with tipboxes, only managed %d", i+1)
		}
	}

	tb, err := inventory.NewTipbox(ctx, lh.Properties.Tips[0].Type)
	if err != nil {
		t.Fatal(err)
	}
	err = lh.Properties.AddTipBox(tb)
	if e, f := "1 (LH_ERR_NO_DECK_SPACE) : insufficient deck space to fit all required items; this may be due to constraints : Trying to add tip box", err.Error(); e != f {
		t.Fatalf("Expected error %q found %q", e, f)
	}
}

func TestDeckSpace2(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())
	lh := GetLiquidHandlerForTest(ctx)

	for i := 0; i < len(lh.Properties.Preferences.Inputs); i++ {
		plate, err := inventory.NewPlate(ctx, "pcrplate_skirted")
		if err != nil {
			t.Fatal(err)
		}

		if err := lh.Properties.AddPlateTo(lh.Properties.Preferences.Inputs[i], plate); err != nil {
			t.Fatalf("position %s is full, should be empty", lh.Properties.Preferences.Inputs[i])
		}
	}

	plate, err := inventory.NewPlate(ctx, "pcrplate_skirted")
	if err != nil {
		t.Fatal(err)
	}

	err = lh.Properties.AddPlateTo(lh.Properties.Preferences.Inputs[0], plate)
	if e, f := "1 (LH_ERR_NO_DECK_SPACE) : insufficient deck space to fit all required items; this may be due to constraints : Trying to add plate to full position position_4", err.Error(); e != f {
		t.Fatalf("Expected error %q found %q", e, f)
	}
}
