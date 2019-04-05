package laboratory

import (
	"fmt"
	"math"

	"github.com/antha-lang/antha/laboratory/compare"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/workflow"
)

// Compare compares output generated with any supplied test data in the workflow
func (lb *LaboratoryBuilder) Compare() error {
	// No testing data available.
	if lb.workflow.Testing == nil {
		lb.Logger.Log("msg", "No comparison test data supplied.")
		return nil
	}

	mixIdx := 0

	for i, instr := range lb.instrs {
		switch t := instr.(type) {
		case *target.Mix:
			lb.Logger.Log("msg", fmt.Sprintf("[%d] Checking mix instruction, %f seconds", i, t.GetTimeEstimate()))
			lb.compareTimings(t, mixIdx)
			lb.compareOutputs(t, mixIdx)
			mixIdx++
		}
	}

	lb.Logger.Log("msg", "Comparison test data passed.")
	return nil
}

func expectedMix(w *workflow.Workflow, idx int) (*workflow.MixTaskCheck, error) {
	if idx >= len(w.Testing.MixTaskChecks) {
		return nil, fmt.Errorf("mix comparison %d not found, only %d mixes are expected", idx, len(w.Testing.MixTaskChecks))
	}

	return &w.Testing.MixTaskChecks[idx], nil
}

func (lb *LaboratoryBuilder) compareTimings(m *target.Mix, idx int) error {
	em, err := expectedMix(lb.workflow, idx)

	if err != nil {
		return fmt.Errorf("timing check failed, %v", err)
	} else if err = compareToPercent(em.TimeEstimate.Seconds(), m.GetTimeEstimate(), 10); err != nil {
		return fmt.Errorf("timing check failed, %v", err)
	}

	lb.Logger.Log("msg", fmt.Sprintf("Passed timing check. Expected %.3gs, found %.3gs.", em.TimeEstimate.Seconds(), m.GetTimeEstimate()))
	return nil
}

func compareToPercent(expected float64, actual float64, percent float64) error {
	if math.Abs(expected-actual) > math.Abs(0.01*percent*expected) {
		return fmt.Errorf("expected %.2g but found %.2g (checked to %.2g%%)", expected, actual, percent)
	}

	return nil
}

func (lb *LaboratoryBuilder) compareOutputs(m *target.Mix, idx int) error {
	em, err := expectedMix(lb.workflow, idx)
	if err != nil {
		return err
	}

	if em.Outputs == nil || len(em.Outputs) == 0 {
		lb.Logger.Log("msg", fmt.Sprintf("No output comparison data supplied for mix task %d", idx))
		return nil
	}

	return compare.Plates(em.Outputs, m.FinalProperties.Plates, lb.effects.IDGenerator)
}
