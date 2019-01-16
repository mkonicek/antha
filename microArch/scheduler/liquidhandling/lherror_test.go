package liquidhandling

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/inventory/testinventory"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"testing"
)

func GetLiquidHandlerForTest(tips []string) (*Liquidhandler, error) {
	if lh, err := liquidhandling.MakeLHForTest(tips); err != nil {
		return nil, err
	} else {
		return Init(lh), nil
	}
}

func GetIndependentLiquidHandlerForTest(tips []string) (*Liquidhandler, error) {
	if gilson, err := liquidhandling.MakeLHForTest(tips); err != nil {
		return nil, err
	} else {
		for _, head := range gilson.Heads {
			head.Params.Independent = true
		}

		for _, adaptor := range gilson.Adaptors {
			adaptor.Params.Independent = true
		}

		return Init(gilson), nil
	}
}

func GetLHRequestForTest() *LHRequest {
	req := NewLHRequest()
	return req
}

func TestNoInstructions(t *testing.T) {
	ctx := testinventory.NewContextForTest(context.Background())
	req := GetLHRequestForTest()
	if lh, err := GetLiquidHandlerForTest(nil); err != nil {
		t.Fatal(err)
	} else if err := lh.MakeSolutions(ctx, req); err.Error() != "9 (LH_ERR_OTHER) :  : Nil plan requested: no Mix Instructions present" {
		t.Fatal(fmt.Sprint("Unexpected error: ", err.Error()))
	}
}

func TestDeckSpace1(t *testing.T) {
	tipboxType := "Gilson200"

	lh, err := GetLiquidHandlerForTest([]string{tipboxType})
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < len(lh.Properties.Preferences.Tipboxes); i++ {
		if tb, err := lh.Properties.TipFactory.NewTipboxByTipType(tipboxType); err != nil {
			t.Fatal(err)
		} else if err := lh.Properties.AddTipBox(tb); err != nil {
			t.Fatalf("Should be able to fill deck with tipboxes, only managed %d", i+1)
		}
	}

	tb, err := lh.Properties.TipFactory.NewTipboxByTipType(tipboxType)
	if err != nil {
		t.Fatal(err)
	}
	err = lh.Properties.AddTipBox(tb)
	if e, f := "1 (LH_ERR_NO_DECK_SPACE) : insufficient deck space to fit all required items; this may be due to constraints : Trying to add tip box", err.Error(); e != f {
		t.Fatalf("Expected error %q found %q", e, f)
	}
}

func TestDeckSpace2(t *testing.T) {
	inv := testinventory.GetInventoryForTest()
	lh, err := GetLiquidHandlerForTest(nil)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < len(lh.Properties.Preferences.Inputs); i++ {
		if plate, err := inv.NewPlate("pcrplate_skirted"); err != nil {
			t.Fatal(err)
		} else if err := lh.Properties.AddPlateTo(lh.Properties.Preferences.Inputs[i], plate); err != nil {
			t.Fatalf("position %s is full, should be empty", lh.Properties.Preferences.Inputs[i])
		}
	}

	plate, err := inv.NewPlate("pcrplate_skirted")
	if err != nil {
		t.Fatal(err)
	}

	err = lh.Properties.AddPlateTo(lh.Properties.Preferences.Inputs[0], plate)
	if e, f := "1 (LH_ERR_NO_DECK_SPACE) : insufficient deck space to fit all required items; this may be due to constraints : Trying to add plate to full position position_4", err.Error(); e != f {
		t.Fatalf("Expected error %q found %q", e, f)
	}
}
