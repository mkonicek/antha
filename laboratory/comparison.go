package laboratory

import (
	"fmt"

	runner "github.com/Synthace/antha-runner/export"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/workflow"
)

func (labBuild *LaboratoryBuilder) Compare() error {
	// No testing data available.
	labBuild.Logger.Log("msg", "Checking comparisons.")
	if labBuild.workflow.Testing == nil {
		labBuild.Logger.Log("msg", "Checking comparisons skipped.")
		return nil
	}
	labBuild.Logger.Log("msg", "Checking comparisons trying...")

	tasks, err := runner.InstrsToTasks(labBuild.effects.IDGenerator, labBuild.instrs)

	for i, instr := range labBuild.instrs {
		switch t := instr.(type) {
		case *target.Mix:
			labBuild.Logger.Log("msg", fmt.Sprintf("[%d] Mix instruction, %f seconds", i, t.GetTimeEstimate()))
			//labBuild.Logger.Log("msg", fmt.Sprintf("[%d] Mix instruction : %v", i, t))
			for j, rb := range t.Request.Instructions {
				labBuild.Logger.Log("msg", fmt.Sprintf("[%d][%d] - Instruction [%s]", i, j, rb.Type().Name))
			}
		default:
			labBuild.Logger.Log("msg", fmt.Sprintf("Other instruction %v", t))
		}
	}

	if err != nil {
		labBuild.Logger.Log("msg", "Error in task generation.", "err", err)
		return err
	}

	labBuild.Logger.Log("msg", fmt.Sprint("Generated ", len(tasks), " tasks."))

	for i, t := range tasks {
		labBuild.Logger.Log("msg", fmt.Sprintf("Task [%d] : %v", i, t))
	}
	return nil
}

func expectedMix(w *workflow.Workflow, idx mixIndex) (*workflow.MixTaskCheck, error) {
	if mtc, ok := w.Testing.MixTaskChecks[idx]; ok != nil {
		return nil, fmt.Errorf("mix comparison %d not found, only %d mixes are expected", idx, len(w.Testing.MixTaskChecks))
	}

	return mtc, nil
} 

func compareTimings(w *Workflow, m *targetMix, idx mixIndex) {
	expected := w.Testing.MixTaskChecks[]
}

func compareInstructions(w *Workflow, m *target.Mix, idx mixIndex) {

}
