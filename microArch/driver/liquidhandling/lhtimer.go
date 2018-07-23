package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"time"
)

// LHTimer provides timing for instructions
type LHTimer interface {
	TimeFor(r RobotInstruction) time.Duration
}

// deprecate this mess
type OldLHTimer struct {
	Times map[*InstructionType]time.Duration
}

func NewTimer() *OldLHTimer {
	return &OldLHTimer{
		Times: make(map[*InstructionType]time.Duration),
	}
}

func (t *OldLHTimer) TimeFor(r RobotInstruction) time.Duration {
	d := t.Times[r.Type()]
	if r.Type() == MIX {
		// get cycles

		prm := r.GetParameter(CYCLES)

		cyc, ok := prm.([]int)

		if ok {
			max := func(ds []int) int {
				res := 0
				for _, elem := range ds {
					if elem > res {
						res = elem
					}
				}
				return res
			}
			d = time.Duration(int64(max(cyc)) * int64(d))
		}
	}

	return d
}

type highLeveltimer struct {
	name     string
	model    string
	flowRate float64 // nl/s
	moveRate float64 // secs/well
	scanRate float64 // secs/well
}

func (hlt highLeveltimer) TimeFor(ins RobotInstruction) time.Duration {
	var totaltime float64

	if ins.Type() == TFR {
		tfr := ins.(*TransferInstruction)
		lastFrom := wtype.WellCoords{}
		lastTo := wtype.WellCoords{}
		for _, mt := range tfr.Transfers {
			for _, t := range mt.Transfers {
				wcF := wtype.MakeWellCoords(t.WellFrom)
				wcT := wtype.MakeWellCoords(t.WellTo)
				totaltime += (manhattan(wcF, lastFrom) + manhattan(wcT, lastTo)) * hlt.moveRate // time to move plates
				totaltime += t.Volume.ConvertToString("nl") / hlt.flowRate                      // time to do fluid transfer
				totaltime += hlt.scanRate                                                       // time to scan src well

			}
		}
	}

	return time.Duration(int64(wutil.RoundInt(totaltime)) * 1000000000)
}

func manhattan(a, b wtype.WellCoords) float64 {
	return float64(absI(a.X-b.X) + absI(a.Y-b.Y))
}

func absI(i int) int {
	if i < 0 {
		return -i
	}

	return i
}
