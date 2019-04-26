package tests

import (
	"fmt"
	"testing"

	"github.com/antha-lang/antha/laboratory"
	"github.com/antha-lang/antha/laboratory/effects/id"
	"github.com/antha-lang/antha/laboratory/testlab"
	lh "github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
)

func GetLiquidHandlerForTest(lab *laboratory.Laboratory) *lh.Liquidhandler {
	return lh.Init(makeGilson(lab))
}

func GetIndependentLiquidHandlerForTest(lab *laboratory.Laboratory) *lh.Liquidhandler {
	gilson := makeGilson(lab)
	for _, head := range gilson.Heads {
		head.Params.Independent = true
	}

	for _, adaptor := range gilson.Adaptors {
		adaptor.Params.Independent = true
	}

	return lh.Init(gilson)
}

func GetLHRequestForTest(idGen *id.IDGenerator) *lh.LHRequest {
	return lh.NewLHRequest(idGen)
}

func TestNoInstructions(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			req := GetLHRequestForTest(lab.IDGenerator)
			lh := GetLiquidHandlerForTest(lab)
			err := lh.MakeSolutions(lab.LaboratoryEffects, req)

			if err.Error() != "9 (LH_ERR_OTHER) :  : Nil plan requested: no Mix Instructions present" {
				return fmt.Errorf("Unexpected error: %v", err)
			}
			return nil
		},
	})
}

func TestDeckSpace1(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			lh := GetLiquidHandlerForTest(lab)

			for i := 0; i < len(lh.Properties.Preferences.Tipboxes); i++ {
				tb, err := lab.Inventory.TipBoxes.NewTipbox(lh.Properties.Tips[0].Type)
				if err != nil {
					return err
				}

				if err := lh.Properties.AddTipBox(tb); err != nil {
					return fmt.Errorf("Should be able to fill deck with tipboxes, only managed %d (%v)", i+1, err)
				}
			}

			tb, err := lab.Inventory.TipBoxes.NewTipbox(lh.Properties.Tips[0].Type)
			if err != nil {
				return err
			}
			err = lh.Properties.AddTipBox(tb)
			if e, f := "1 (LH_ERR_NO_DECK_SPACE) : insufficient deck space to fit all required items; this may be due to constraints : Trying to add tip box", err.Error(); e != f {
				return fmt.Errorf("Expected error %q found %q", e, f)
			}
			return nil
		},
	})
}

func TestDeckSpace2(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			lh := GetLiquidHandlerForTest(lab)

			for i := 0; i < len(lh.Properties.Preferences.Inputs); i++ {
				plate, err := lab.Inventory.Plates.NewPlate("pcrplate_skirted")
				if err != nil {
					return err
				}

				if err := lh.Properties.AddPlateTo(lh.Properties.Preferences.Inputs[i], plate); err != nil {
					return fmt.Errorf("position %s is full, should be empty (%v)", lh.Properties.Preferences.Inputs[i], err)
				}
			}

			plate, err := lab.Inventory.Plates.NewPlate("pcrplate_skirted")
			if err != nil {
				return err
			}

			err = lh.Properties.AddPlateTo(lh.Properties.Preferences.Inputs[0], plate)
			if e, f := "1 (LH_ERR_NO_DECK_SPACE) : insufficient deck space to fit all required items; this may be due to constraints : Trying to add plate to full position position_4", err.Error(); e != f {
				return fmt.Errorf("Expected error %q found %q", e, f)
			}
			return nil
		},
	})
}
