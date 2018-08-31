package liquidhandling

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	driver "github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

func (lh Liquidhandler) countTipsUsed(rq *LHRequest) (*LHRequest, error) {
	teHash := make(map[string]wtype.TipEstimate)

	var err error
	for _, ins := range rq.Instructions {
		ins.Visit(driver.RobotInstructionBaseVisitor{
			HandleLoadTips: func(ins *driver.LoadTipsInstruction) {
				for i := 0; i < len(ins.Pos); i++ {
					if ins.Pos[i] == "" {
						continue
					}
					bx, ok := lh.Properties.Tipboxes[ins.Pos[i]]

					if !ok {
						err = fmt.Errorf("Instruction %s requests tips from an empty position", driver.InsToString(ins))
						return
					}

					tt := bx.Type
					te, ok := teHash[tt]

					if !ok {
						te = wtype.TipEstimate{TipType: tt, NTipBoxes: bx.NTips}
					}

					te.NTips += 1

					teHash[te.TipType] = te
				}
			},
		})
	}
	if err != nil {
		return nil, err
	}

	// output to the request

	for _, te := range teHash {
		if te.NTips == 0 {
			// should not be possible but I don't want to generate confusion
			continue
		}
		// above we have recorded the total number of tips in a box of lh type
		// in NTipBoxes, here we use it to determine how many boxes are needed
		dv := te.NTips / te.NTipBoxes
		mod := te.NTips % te.NTipBoxes

		te.NTipBoxes = dv

		if mod != 0 {
			te.NTipBoxes += 1
		}

		rq.TipsUsed = append(rq.TipsUsed, te)
	}

	return rq, nil
}

//func countWellMulti(sa []string) int {
//	r := 0
//	for _, s := range sa {
//		if s != "" {
//			r += 1
//		}
//	}

// 	return r
// }
