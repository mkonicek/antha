package workflowtest

import (
	"errors"
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/target"
)

// A TestOpt is an option for running a test
type TestOpt struct {
	ComparisonOptions   string
	CompareInstructions bool
	CompareOutputs      bool
	Results             TestResults
}

// TestResults are the results of running a set of tests
type TestResults struct {
	MixTaskResults []MixTaskResult
}

// A MixTaskResult is the result if running a mix task
type MixTaskResult struct {
	Instructions liquidhandling.SetOfRobotInstructions
	Outputs      map[string]*wtype.Plate
}

func generaliseInstructions(insIn []liquidhandling.TerminalRobotInstruction) []liquidhandling.RobotInstruction {
	insOut := make([]liquidhandling.RobotInstruction, len(insIn))

	for i := 0; i < len(insIn); i++ {
		// this must succeed since TerminalRobotInstructions ARE RobotInstructions, as the name suggests...
		// kinda makes you wonder why Go forces you to do this, really, doesn't it?
		insOut[i] = insIn[i].(liquidhandling.RobotInstruction)
	}
	return insOut
}

func compareOutputs(outputs1, outputs2 map[string]*wtype.Plate, opt TestOpt) string {
	return joinErrors(CompareMixOutputs(outputs1, outputs2, unpackOutputComparisonOptions(opt.ComparisonOptions)).Errors)
}

func compareInstructions(genIns1, genIns2 []liquidhandling.RobotInstruction, opt TestOpt) string {
	return joinErrors(
		liquidhandling.CompareInstructionSets(
			genIns1,
			genIns2,
			liquidhandling.ComparisonOpt{
				InstructionParameters: unpackInstructionComparisonOptions(opt.ComparisonOptions),
			}).Errors)
}

// CompareTestResults compares an execution with an expected output
func CompareTestResults(runResult *execute.Result, opt TestOpt) error {
	// pull out mix tasks from the Result

	mixTasks := getMixTasks(runResult)

	if len(mixTasks) != len(opt.Results.MixTaskResults) {
		return fmt.Errorf("Number of mix tasks differs: expected %d got %d", len(opt.Results.MixTaskResults), len(mixTasks))
	}

	errstr := ""

	for i := 0; i < len(mixTasks); i++ {
		if opt.CompareInstructions {
			genIns1 := opt.Results.MixTaskResults[i].Instructions.RobotInstructions // already RobotInstructions
			genIns2 := generaliseInstructions(mixTasks[i].Request.Instructions)
			ssss := compareInstructions(genIns1, genIns2, opt)
			if ssss != "" {
				errstr += ssss + "\n"
			}
		} else if opt.CompareOutputs {
			ssss := compareOutputs(opt.Results.MixTaskResults[i].Outputs, getMixTaskOutputs(mixTasks[i]), opt)
			if ssss != "" {
				errstr += ssss + "\n"
			}
		}
	}

	if errstr != "" {
		return errors.New(errstr)
	}
	return nil
}

func getMixTaskOutputs(mix *target.Mix) map[string]*wtype.Plate {
	outputs := make(map[string]*wtype.Plate)

	// get output plates (ONLY)

	for _, pos := range mix.FinalProperties.Output_preferences {
		plate, ok := mix.FinalProperties.Plates[pos]

		if ok {
			outputs[pos] = plate
		}
	}

	return outputs
}

// SaveTestOutputs extracts a TestOpt from an execution result
func SaveTestOutputs(runResult *execute.Result, comparisonOptions string) TestOpt {
	// get mix tasks
	mixTasks := getMixTasks(runResult)

	mixTaskResults := make([]MixTaskResult, len(mixTasks))

	for i := 0; i < len(mixTasks); i++ {
		outputs := getMixTaskOutputs(mixTasks[i])
		mixTaskResults[i] = MixTaskResult{
			Instructions: liquidhandling.SetOfRobotInstructions{
				RobotInstructions: generaliseInstructions(mixTasks[i].Request.Instructions),
			}, Outputs: outputs,
		}
	}

	results := TestResults{MixTaskResults: mixTaskResults}
	return TestOpt{
		Results:           results,
		ComparisonOptions: comparisonOptions,
		CompareOutputs:    true,
	}
}

func unpackOutputComparisonOptions(optIn string) ComparisonMode {
	// v0 do the sensible thing
	return ComparePlateTypesVolumes
}

func unpackInstructionComparisonOptions(optIn string) map[string][]string {
	// v0 --> just compare everything
	return liquidhandling.CompareAllParameters()
}

func joinErrors(errors []error) string {
	r := make([]string, len(errors))
	for i, err := range errors {
		r[i] = err.Error()
	}

	return strings.Join(r, "\n")
}

func getMixTasks(runResult *execute.Result) []*target.Mix {
	ret := make([]*target.Mix, 0, len(runResult.Insts))

	for _, inst := range runResult.Insts {
		switch inst.(type) {
		case *target.Mix:
			ret = append(ret, inst.(*target.Mix))
		}
	}

	return ret
}
