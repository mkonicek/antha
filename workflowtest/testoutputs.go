package workflowtest

import (
	"errors"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/target"
	"strings"
)

type TestOpt struct {
	ComparisonOptions string
	Results           TestResults
}

type TestResults struct {
	MixTaskResults []MixTaskResult
}

type MixTaskResult struct {
	Instructions liquidhandling.SetOfRobotInstructions
	Outputs      []*wtype.LHPlate
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

func CompareTestResults(runResult *execute.Result, opt TestOpt) error {
	// pull out mix tasks from the Result

	mixTasks := getMixTasks(runResult)

	if len(mixTasks) != len(opt.Results.MixTaskResults) {
		return fmt.Errorf("Number of mix tasks differs: expected %d got %d", len(opt.Results.MixTaskResults), len(mixTasks))
	}

	errstr := ""

	for i := 0; i < len(mixTasks); i++ {
		genIns1 := opt.Results.MixTaskResults[i].Instructions.RobotInstructions // already RobotInstructions
		genIns2 := generaliseInstructions(mixTasks[i].Request.Instructions)

		ssss := joinErrors(liquidhandling.CompareInstructionSets(genIns1, genIns2, liquidhandling.ComparisonOpt{unpackOpt(opt.ComparisonOptions)}).Errors)
		if ssss != "" {
			errstr += ssss + "\n"
		}
	}

	return errors.New(errstr)
}

func SaveTestOutputs(runResult *execute.Result, comparisonOptions string) TestOpt {
	// get mix tasks
	mixTasks := getMixTasks(runResult)

	mixTaskResults := make([]MixTaskResult, len(mixTasks))

	for i := 0; i < len(mixTasks); i++ {
		mixTaskResults[i] = MixTaskResult{Instructions: liquidhandling.SetOfRobotInstructions{RobotInstructions: generaliseInstructions(mixTasks[i].Request.Instructions)}}
	}

	results := TestResults{MixTaskResults: mixTaskResults}
	return TestOpt{Results: results, ComparisonOptions: comparisonOptions}
}

func unpackOpt(optIn string) map[string][]string {
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
