package laboratory

import (
	"fmt"
	"math"

	runner "github.com/Synthace/antha-runner/export"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/workflow"
)

func (lb *LaboratoryBuilder) Compare() error {
	// No testing data available.
	if lb.workflow.Testing == nil {
		lb.Logger.Log("msg", "No comparison test data supplied.")
		return nil
	}

	tasks, err := runner.InstrsToTasks(lb.effects.IDGenerator, lb.instrs)
	mixIdx := 0

	for i, instr := range lb.instrs {
		switch t := instr.(type) {
		case *target.Mix:
			lb.Logger.Log("msg", fmt.Sprintf("[%d] Checking mix instruction, %f seconds", i, t.GetTimeEstimate()))
			lb.compareTimings(t, mixIdx)
			mixIdx++
		}
	}

	if mixIdx != len(lb.workflow.Testing.MixTaskChecks) {
		return fmt.Errorf("Expected %d mix tasks, found %d", len(lb.workflow.Testing.MixTaskChecks), mixIdx)
	}

	if err != nil {
		lb.Logger.Log("msg", "Error in task generation.", "err", err)
		return err
	}

	lb.Logger.Log("msg", fmt.Sprint("Generated ", len(tasks), " tasks."))

	for i, t := range tasks {
		lb.Logger.Log("msg", fmt.Sprintf("Task [%d] : %v", i, t))
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

func compareInstructions(w *workflow.Workflow, m *target.Mix, idx int) error {
	return nil
}

func compareOutputs(w *workflow.Workflow, m *target.Mix, idx int) error {
	return nil
}
