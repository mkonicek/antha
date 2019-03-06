package workflowtest

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

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
	TimeEstimate time.Duration
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

func compareInstructions(genIns1, genIns2 []liquidhandling.RobotInstruction, opt TestOpt) error {
	return liquidhandling.CompareInstructionSets(
		genIns1,
		genIns2,
		unpackInstructionComparisonOptions(opt.ComparisonOptions)...,
	)
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
			err := compareInstructions(genIns1, genIns2, opt)
			if err != nil {
				errstr += err.Error() + "\n"
			}
			// step to compare time estimate
			if err := compareTimeEstimates(
				opt.Results.MixTaskResults[i].TimeEstimate.Seconds(),
				mixTasks[i].Request.TimeEstimate,
				timeEstimatePrecisionFactor); err != nil {
				errstr += err.Error() + "\n"
			}
		} else if opt.CompareOutputs {
			ssss := compareOutputs(opt.Results.MixTaskResults[i].Outputs, mixTasks[i].FinalProperties.Plates, opt)
			if ssss != "" {
				errstr += ssss + "\n"
			}
			// step to compare time estimate
			if err := compareTimeEstimates(
				opt.Results.MixTaskResults[i].TimeEstimate.Seconds(),
				mixTasks[i].Request.TimeEstimate,
				timeEstimatePrecisionFactor); err != nil {
				errstr += err.Error() + "\n"
			}
		}
	}

	if errstr != "" {
		return errors.New(errstr)
	}
	return nil
}

// permitted proportional difference between test result time estimate and returned result.
const timeEstimatePrecisionFactor = 0.1

// compareTimeEstimates returns an error if the testTimeInSecs deviates from the expectedTimeInSecs
// by greater than the expectedTimeInSecs * precisionFactor
func compareTimeEstimates(expectedTimeInSecs, testTimeInSecs, precisionFactor float64) error {
	if math.Abs(expectedTimeInSecs-testTimeInSecs) > (timeEstimatePrecisionFactor * expectedTimeInSecs) {
		return fmt.Errorf(
			"Expected time estimate %f seconds but got %f seconds; \n"+
				"Time estimates must be equal within %f %% to be permitted",
			expectedTimeInSecs,
			testTimeInSecs,
			precisionFactor*100.0,
		)
	}
	return nil
}

// SaveTestOutputs extracts a TestOpt from an execution result
func SaveTestOutputs(runResult *execute.Result, comparisonOptions string) TestOpt {
	// get mix tasks
	mixTasks := getMixTasks(runResult)

	mixTaskResults := make([]MixTaskResult, len(mixTasks))

	for i := 0; i < len(mixTasks); i++ {
		outputs := mixTasks[i].FinalProperties.Plates
		mixTaskResults[i] = MixTaskResult{
			Instructions: liquidhandling.SetOfRobotInstructions{
				RobotInstructions: generaliseInstructions(mixTasks[i].Request.Instructions),
			},
			Outputs: outputs,
			// We're ok with truncating to the nearest second by casting into an int64;
			// this is much more precise than the required precision.
			TimeEstimate: time.Duration(time.Duration(mixTasks[i].Request.TimeEstimate) * time.Second),
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

func unpackInstructionComparisonOptions(optIn string) []liquidhandling.RobotInstructionComparatorFunc {
	// v0 --> just compare everything
	return liquidhandling.CompareAllParameters
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
		if mix, ok := inst.(*target.Mix); ok {
			ret = append(ret, mix)
		}
	}

	return ret
}
