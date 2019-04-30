package laboratory

import (
	"fmt"
	"math"
	"path/filepath"

	"github.com/antha-lang/antha/laboratory/compare"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/utils"
	"github.com/antha-lang/antha/workflow"
)

// Compare compares output generated with any supplied test data in the workflow
func (labBuild *LaboratoryBuilder) Compare() {

	if len(labBuild.Workflow.Testing.MixTaskChecks) == 0 {
		return
	}

	mixIdx := 0
	errs := make(utils.ErrorSlice, 0, len(labBuild.instrs))

	for i, instr := range labBuild.instrs {
		if t, ok := instr.(*target.Mix); ok {
			labBuild.Logger.Log("msg", "checking mix instruction", "index", i)

			if expected, err := expectedMix(labBuild.Workflow, mixIdx); err != nil {
				errs = append(errs, err)
			} else {
				errs = append(errs, labBuild.compareTimings(t, expected))
				errs = append(errs, labBuild.compareOutputs(t, expected)...)
			}
			mixIdx++
		}
	}

	filename := filepath.Join(labBuild.outDir, "comparisons.json")
	if err := errs.WriteToFile(filename); err != nil {
		labBuild.RecordError(fmt.Errorf("errors writing comparison results to file %s: %v", filename, err), true)
	} else if errs.Pack() != nil {
		labBuild.RecordError(fmt.Errorf("errors in comparison tests, details in %s", filename), true)
	} else {
		labBuild.Logger.Log("msg", "Comparison test data passed.")
	}
}

func expectedMix(w *workflow.Workflow, idx int) (*workflow.MixTaskCheck, error) {
	if idx >= len(w.Testing.MixTaskChecks) {
		return nil, fmt.Errorf("mix comparison %d not found, only %d mixes are expected", idx, len(w.Testing.MixTaskChecks))
	}

	return &w.Testing.MixTaskChecks[idx], nil
}

func (labBuild *LaboratoryBuilder) compareTimings(m *target.Mix, expectedMix *workflow.MixTaskCheck) error {
	const timeAccuracyPercent = 10
	if err := compareToPercent(expectedMix.TimeEstimate.Seconds(), m.GetTimeEstimate(), timeAccuracyPercent); err != nil {
		return fmt.Errorf("timing check failed, %v", err)
	}

	labBuild.Logger.Log("msg", fmt.Sprintf("Passed timing check. Expected %.3gs, found %.3gs.", expectedMix.TimeEstimate.Seconds(), m.GetTimeEstimate()))
	return nil
}

func compareToPercent(expected float64, actual float64, percent float64) error {
	const onePercent = 0.01
	if math.Abs(expected-actual) > math.Abs(onePercent*percent*expected) {
		return fmt.Errorf("expected %.2g but found %.2g (checked to %.2g%%)", expected, actual, percent)
	}

	return nil
}

func (labBuild *LaboratoryBuilder) compareOutputs(m *target.Mix, expectedMix *workflow.MixTaskCheck) utils.ErrorSlice {
	if expectedMix.Outputs == nil || len(expectedMix.Outputs) == 0 {
		labBuild.Logger.Log("msg", "No output comparison data supplied for mix task.")
		return nil
	}

	return compare.Plates(labBuild.effects.IDGenerator, expectedMix.Outputs, m.FinalProperties.Plates)
}
