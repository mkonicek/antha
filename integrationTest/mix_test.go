package main

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"
	dlh "github.com/antha-lang/antha/microArch/driver/liquidhandling"
	slh "github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
	"strings"
	"testing"
)

type TestCase struct {
	Context    context.Context
	S1         bool   // is S1 a sample
	Command    string // command name
	PlateName  string // plate name if specified
	PlateType  string // plate type if needed
	PlateID    string // plate ID if specified
	Well       string // well location if specified
	Consistent bool   // if no, generate a case which errors
	Error      string // expect this error
	InPlace    bool   // expect result to be a mix-in-place
}

func (tc TestCase) Name() string {
	return fmt.Sprintf("%s_sample1:%t__platename:%s_platetype:%s_plateid:%s_well:%s_Consistent:%t", tc.Command, tc.S1, tc.PlateName, tc.PlateType, tc.PlateID, tc.Well, tc.Consistent)
}

func (tc TestCase) MixOptions() mixer.MixOptions {
	components := standardComponents()
	components[0].SetSample(tc.S1)
	var destPlate *wtype.LHPlate
	var err error

	// PlateID and PlateName are mutually exclusive
	if tc.PlateID != "" {
		// need to create a plate
		destPlate, err = inventory.NewPlate(tc.Context, tc.PlateType)

		if err != nil {
			panic(err.Error())
		}

		destPlate.ResetID(tc.PlateID)

		if tc.Consistent {
			// the first component must be where we expect
			well := tc.Well
			if well == "" {
				well = "A1"
			}

			destPlate.Wellcoords[well].AddComponent(components[0])
			components[0].Loc = destPlate.Wellcoords[well].WContents.Loc
		}
	} else if tc.PlateName != "" {
		// just here to keep this simple
	} else {
		if tc.Consistent {
			if !tc.S1 {
				// c1 needs a location
				destPlate, err = inventory.NewPlate(tc.Context, tc.PlateType)
				if err != nil {
					panic(err.Error())
				}
				destPlate.Wellcoords["A1"].AddComponent(components[0])
				components[0].Loc = destPlate.Wellcoords["A1"].WContents.Loc
			}
		}
	}

	return mixer.MixOptions{
		Components:  components,
		PlateType:   tc.PlateType,
		Address:     tc.Well,
		PlateName:   tc.PlateName,
		Destination: destPlate,
	}
}

func makeTestCases() []TestCase {
	commands := []string{"Mix", "MixInto", "MixNamed"}
	bools := []bool{false, true}

	testcases := make([]TestCase, 0, 1)

	for _, cmd := range commands {
		plateid, platetype, platename := getParamsForCommand(cmd)

		for _, s1 := range bools {
			for _, c := range bools {
				if cmd == "Mix" && !s1 && !c {
					// inconsistency not relevant here
					continue
				}
				errString := getErrorForTest(cmd, s1, c)
				inplace := isInPlace(cmd, s1, c)
				testcases = append(testcases, TestCase{
					Context:    getContextForTest(),
					Command:    cmd,
					S1:         s1,
					PlateType:  platetype,
					PlateID:    plateid,
					PlateName:  platename,
					Consistent: c,
					Error:      errString,
					InPlace:    inplace,
				})
			}
		}
	}

	return testcases
}

func Test(t *testing.T) {
	testCases := makeTestCases()

	for _, testCase := range testCases {
		test := func(t *testing.T) {
			robot := dlh.MakeGilsonWithTipboxesForTest()
			lh := slh.Init(robot)
			rq := slh.NewLHRequest()
			pt, _ := inventory.NewPlate(testCase.Context, "pcrplate_skirted_riser18")
			rq.Input_platetypes = append(rq.Input_platetypes, pt.Dup())
			rq.Output_platetypes = append(rq.Output_platetypes, pt.Dup())
			mixOptions := testCase.MixOptions()
			ins := mixer.GenericMix(mixOptions)
			rq.LHInstructions[ins.ID] = ins
			if mixOptions.Destination != nil {
				rq.Input_plates[mixOptions.Destination.ID] = mixOptions.Destination
			}
			err := rq.ConfigureYourself()

			if err != nil {
				t.Errorf(err.Error())
			}

			err = lh.Plan(testCase.Context, rq)

			if testCase.Error != "" {
				if err == nil {
					t.Errorf("Expected error %s, got nil", testCase.Error)
				} else if err.Error() != testCase.Error {
					t.Errorf("Wrong error returned: got %s expected %s", err.Error(), testCase.Error)
				}
			} else if err != nil {
				t.Errorf(err.Error())
			}

			if testCase.Consistent {
				// test 1 - did we get a mix in place when we expected one?
				testInPlace(t, rq, testCase.InPlace)

				if testCase.PlateID != "" {
					// test 2a - if we specified a plateID did we get it?
					testPlateID(t, rq, testCase.PlateID)
				} else if testCase.PlateName != "" {
					// test 2b - if we specified a plate name did both components move?
					testPlateName(t, rq)
				}
			}
		}

		t.Run(testCase.Name(), test)
	}
}

func standardComponents() []*wtype.Liquid {
	return []*wtype.Liquid{
		wtype.NewLiquid("water", 50, "ul", wtype.LTWater),
		wtype.NewLiquid("coffee", 30, "ul", wtype.LTWater),
	}
}

func getContextForTest() context.Context {
	ctx := context.Background()
	invCtx := testinventory.NewContext(ctx)
	return invCtx
}

func standardPlate() *wtype.Plate {
	invCtx := getContextForTest()
	plate, _ := inventory.NewPlate(invCtx, "pcrplate_skirted_riser18")
	plate.PlateName = "destplate"
	return plate
}

func getParamsForCommand(cmd string) (plateid, platetype, platename string) {
	switch cmd {
	case "Mix":
		platetype = "pcrplate_skirted_riser18"
	case "MixInto":
		plateid = wtype.GetUUID()
		platetype = "pcrplate_skirted_riser18"
	case "MixNamed":
		platename = "destplate"
		platetype = "pcrplate_skirted"
	}

	return
}

func getErrorForTest(cmd string, s1 bool, consistent bool) string {
	if (cmd == "Mix" || cmd == "MixInto") && !s1 {
		if !consistent {
			return "8 (LH_ERR_DIRE) : an internal error : MIX IN PLACE WITH NO LOCATION SET"
		}
	}

	return ""
}

func isInPlace(cmd string, s1, consistent bool) bool {
	if cmd == "MixNamed" || !consistent {
		return false
	}

	return !s1
}

// in this case we are mixing 'coffee' onto 'water'
// so we check whether the input assignment for water
// is also used as an output assignment
func testInPlace(t *testing.T, rq *slh.LHRequest, expected bool) {
	// does not apply here
	if len(rq.LHInstructions) != 1 {
		t.Errorf("testInPlace called with more than one LHInstruction in request: %d", len(rq.LHInstructions))
		return
	}

	assWater, ok := rq.Input_assignments["water"]

	if !ok {
		t.Errorf("No input assignment for water")
		return
	}

	if expected && len(assWater) != 1 {
		t.Errorf("Wrong number of input assignments for water: %d (need 1)", len(assWater))
		return
	}

	_, found := rq.Output_assignments[assWater[0]]

	if expected && !found {
		t.Errorf("Mix in place expected but not found")
	} else if !expected && found {
		t.Errorf("Mix in place found but not expected")
	}
}

func testPlateID(t *testing.T, rq *slh.LHRequest, plateID string) {
	// does not apply here
	if len(rq.LHInstructions) != 1 {
		t.Errorf("testPlateID called with more than one LHInstruction in request: %d", len(rq.LHInstructions))
		return
	}

	if len(rq.Output_assignments) != 1 {
		t.Errorf("testPlateID - wrong number of output assignments: %d (need 1)", len(rq.Output_assignments))
		return
	}

	for k := range rq.Output_assignments {
		if !strings.HasPrefix(k, plateID) {
			t.Errorf("testPlateID: plate ID %s expected but not used in assignment %s", plateID, k)
			return
		}
	}
}

func testPlateName(t *testing.T, rq *slh.LHRequest) {
	// does not apply here
	if len(rq.LHInstructions) != 1 {
		t.Errorf("testPlateName called with more than one LHInstruction in request: %d", len(rq.LHInstructions))
		return
	}

	if len(rq.Output_assignments) != 1 {
		t.Errorf("testPlateName - wrong number of output assignments: %d (need 1)", len(rq.Output_assignments))
		return
	}

	// ensure no assignments used for inputs are also used as outputs

	for _, ass := range rq.Input_assignments {
		for _, str := range ass {
			_, used := rq.Output_assignments[str]

			if used {
				t.Errorf("Input assignment reused as output assignment: %s ", str)
				return
			}
		}
	}
}
