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
	Times []time.Duration
}

func NewTimer() *OldLHTimer {
	var t OldLHTimer
	t.Times = make([]time.Duration, 50)
	return &t
}

func (t *OldLHTimer) TimeFor(r RobotInstruction) time.Duration {
	var d time.Duration

	if r.InstructionType() > 0 && r.InstructionType() < len(t.Times) {
		d = t.Times[r.InstructionType()]
		max := func(a []int) int {
			m := a[0]
			for i := 1; i < len(a); i++ {
				if m < a[i] {
					m = a[i]
				}
			}

			return m
		}
		if r.InstructionType() == 34 { // MIX
			// get cycles

			prm := r.GetParameter("CYCLES")

			cyc, ok := prm.([]int)

			if ok {
				d = time.Duration(int64(max(cyc)) * int64(d))
			}
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

	if InstructionTypeName(ins) == "TFR" {
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
