package tests

import (
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/laboratory"
	"github.com/antha-lang/antha/laboratory/testlab"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

func TestParallelSetGeneration(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			// lazy way to get pre-prepared instructions
			tb, dstp := getTransferBlock2Component(lab) // defined in transferblock_test

			ins := make([]*wtype.LHInstruction, 0, len(tb.Inss)-1)

			for i := 0; i < len(tb.Inss); i++ {
				// make one hole
				if i == 4 {
					continue
				}
				ins = append(ins, tb.Inss[i])
			}

			tb.Inss = ins

			rbt := getTestRobot(lab, dstp, "pcrplate_skirted_riser40")

			// allow independent multichannel activity
			headsLoaded := rbt.GetLoadedHeads()
			headsLoaded[0].Params.Independent = true

			pol, err := wtype.GetLHPolicyForTest()
			if err != nil {
				return err
			}

			// allow multi

			pol.Policies["water"]["CAN_MULTI"] = true

			insIds, err := liquidhandling.GetParallelSetsHead(lab.LaboratoryEffects, headsLoaded[0], ins)

			if err != nil {
				return fmt.Errorf("Parallel set generation error: %s\n", err)
			}

			if len(insIds[0]) != 8 {
				return fmt.Errorf("Should have 8 insIDs, instead have %d", len(insIds[0]))
			}

			// insIds[0][4] should be ""

			if insIds[0][4] != "" {
				return fmt.Errorf("InsIds[0][4] should be \"\", it isn't, it's: %s", insIds[0][4])
			}
			return nil
		},
	})
}
