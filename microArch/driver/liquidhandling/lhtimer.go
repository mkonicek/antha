package liquidhandling

import "time"

// records timing info
// preliminary implementation assumes all instructions of a given
// type have the same timing, TimeFor is expressed in terms of the instruction
// however so it will be possible to modify this behaviour in future

type LHTimer struct {
	Times []time.Duration
}

func NewTimer() *LHTimer {
	var t LHTimer
	t.Times = make([]time.Duration, 50)
	return &t
}

func (t *LHTimer) TimeFor(r RobotInstruction) time.Duration {
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

	} else {
	}
	return d
}
