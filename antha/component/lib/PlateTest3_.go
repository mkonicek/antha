// Protocol to allow for rapid combinatorial testing of plate, liquid class combinations.
// Allows testing of effect of liquid handling changes such as offsets and liquid class changes
// Intended to be run prior to any liquid handling change before accepting pull requests.
// The element creates an output csv file which can be filled in by the user to log observed offsets
// for each condition
package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/cmd/antharun/cmd"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/factory"

	"encoding/csv"
	"os"
	//"path/filepath"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/component"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"golang.org/x/net/context"
	"time"
)

// Input parameters for this protocol (data)

// name of test e.g. branch name, date, name of project; csv file will be named after this
// can be whatever you want to call it
// List of volumes to test
// corresponding to valid antha liquid types
// list of out plate types to test
// optional slice of ints which should match the length and order of the OutPlates slice

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _PlateTest3Requirements() {

}

// Conditions to run on startup
func _PlateTest3Setup(_ctx context.Context, _input *PlateTest3Input) {

}

// The core process for this protocol, with the steps to be performed
// for every input
func _PlateTest3Steps(_ctx context.Context, _input *PlateTest3Input, _output *PlateTest3Output) {

	// prepare header to add data and export csv

	if _input.TestName == "" && len(_input.OutPlates) > 1 {
		_input.TestName = "PlateHeightTest" + fmt.Sprint(time.Now().Format("20060102150405"))
	} else if _input.TestName == "" && len(_input.OutPlates) == 1 {
		_input.TestName = "PlateHeightTest" + _input.OutPlates[0] + fmt.Sprint(time.Now().Format("20060102150405"))
	} else {
		_input.TestName = _input.TestName + fmt.Sprint(time.Now().Format("20060102150405"))
	}
	outputfilename := _input.TestName + ".csv"

	csvfile, err := os.Create(outputfilename)
	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	defer csvfile.Close()

	records := make([][]string, 0)

	title := []string{"Plate Height Test:", _input.TestName}
	time := []string{"Time:", fmt.Sprint(time.Now())}

	// find git commit id
	anthacommit, err := cmd.GitCommit()

	if err != nil {
		execute.Errorf(_ctx, err.Error())
	}

	gitcommit := []string{"antha-lang/antha commitID:", anthacommit}

	spacer := []string{}

	header := []string{"TestName", "Plate", "Liquid name", "Liquid type ", "Liquid Volume", "Well ", "Aspirate mm from Bottom of well? ", "Dispense mm from Bottom of well? ", "Acceptable? ", " Comments", "plate Z start", "Plate Height", "lhpolicy aspirate z offset", "lhpolicy dispense z offset"}
	records = append(records, title, time, gitcommit, spacer, header)

	// Make slices to fill up later before exporting as outputs
	_output.FinalSolutions = make([]*wtype.LHComponent, 0)
	_output.WellsUsedPostRunPerPlate = make([]int, 0)
	_output.PlatesUsedPostRunPerPlateType = make([]int, 0)

	// Get list of plates to check validity of plate names specified in parameters
	platelist := factory.GetPlateList()

	// This if statement ensures that default behaviour should be to assume that
	// all plates have no wells used if no WellsUsedperPlateTypeInorder []int is specified
	// in input parameters
	if _input.WellsUsedperOutPlateInorder == nil || len(_input.WellsUsedperOutPlateInorder) == 0 {
		_input.WellsUsedperOutPlateInorder = make([]int, len(_input.OutPlates))
		for l := range _input.OutPlates {
			_input.WellsUsedperOutPlateInorder[l] = 0
		}
	}

	// Range through all plates first
	for k := range _input.OutPlates {

		// set plate number to 1 to start with
		var platenumber int = 1

		// get all well positions from the plate

		lhplate := factory.GetPlateByType(_input.OutPlates[k])

		wellpositionsarray := lhplate.AllWellPositions(wtype.BYCOLUMN)

		// Initialise a counter to be equal to the number of wells used for that plate
		// The counter will be used to select the correct well position
		// if no well position is specified the scheduler will by default select the next well position
		// however using the counter gives flexibility to resume from a given well position if
		// a plate is already partially filled
		counter := _input.WellsUsedperOutPlateInorder[k]

		// range through different volumes to ensure correct behaviour with different pipette heads
		// recommended defaults would be "5ul" and "100"
		for j := range _input.LiquidVolumes {

			// range through liquid types
			//	for i := range LiquidTypes {

			for i, Startingsolution := range _input.Startingsolutions {

				startingsolutionname := Startingsolution.CName

				liquidtypestring, err := wtype.LiquidTypeFromString(_input.LiquidTypes[i])

				// check liquid type is valid
				if err != nil {
					execute.Errorf(_ctx, "Liquid type issue with ", _input.LiquidTypes[i], err.Error())
				}

				Startingsolution.CName = startingsolutionname + "_" + _input.LiquidTypes[i]
				// change liquid type to that specified in loop
				Startingsolution.Type = liquidtypestring

				// sample
				sample := mixer.Sample(Startingsolution, _input.LiquidVolumes[j])

				// check validity of plate name; is it in the plate factory?
				if !search.InSlice(_input.OutPlates[k], platelist) {
					execute.Errorf(_ctx, "No plate ", _input.OutPlates[k], " found in library ", platelist)
				}

				// Mix into a plate at next well position, plate name is given as the type of plate + platenumber

				platename := fmt.Sprint(_input.OutPlates[k], "_Platenumber_", platenumber)

				finalSolution := execute.MixNamed(_ctx, _input.OutPlates[k], wellpositionsarray[counter], platename, sample)
				_output.FinalSolutions = append(_output.FinalSolutions, finalSolution)

				// Append status
				_output.Status = _output.Status + fmt.Sprintln(_input.LiquidVolumes[j].ToString(), " of ", _input.Liquidname, "Liquid type ", _input.LiquidTypes[i], "was mixed into "+_input.OutPlates[k])

				// get specific plate info

				plateheight := lhplate.Height
				zstart := lhplate.WellZStart
				/*
					Height float64
					WellXOffset float64            // distance (mm) between well centres in X direction
					WellYOffset float64            // distance (mm) between well centres in Y direction
					WellXStart  float64            // offset (mm) to first well in X direction
					WellYStart  float64            // offset (mm) to first well in Y direction
					WellZStart  float64            // offset (mm) to bottom of well in Z direction
				*/
				// get lhpolicyinfo

				// print out LHPolicy info
				policy, ok := liquidhandling.GetPolicyByName(_input.LiquidTypes[i])

				if !ok {
					execute.Errorf(_ctx, fmt.Sprint("Liquid type, ", _input.LiquidTypes[i], "not found"))
				}

				aspz := policy["ASPZOFFSET"]
				dspz := policy["DSPZOFFSET"]

				record := []string{_input.TestName, platename, _input.Liquidname, _input.LiquidTypes[i], _input.LiquidVolumes[j].ToString(), wellpositionsarray[counter], "  ", "  ", " ", " ", fmt.Sprint(zstart), fmt.Sprint(plateheight), fmt.Sprint(aspz), fmt.Sprint(dspz)}
				records = append(records, record)

				// evaluate whether plate is full and if so add new plate
				if counter+1 == len(wellpositionsarray) {
					platenumber++
					counter = 0
					// else increase counter ready for next instance of loop
				} else {
					counter++
				}

			}
		}

		// export wells used once all aspirate and dispenses for a particular plate type
		// sticking to plate order specified in input parameters
		_output.WellsUsedPostRunPerPlate = append(_output.WellsUsedPostRunPerPlate, counter)

		if counter > 0 {
			_output.PlatesUsedPostRunPerPlateType = append(_output.PlatesUsedPostRunPerPlateType, platenumber)
		} else {
			_output.PlatesUsedPostRunPerPlateType = append(_output.PlatesUsedPostRunPerPlateType, platenumber-1)

		}
	}

	csvwriter := csv.NewWriter(csvfile)

	for _, record := range records {

		err = csvwriter.Write(record)

		if err != nil {
			execute.Errorf(_ctx, err.Error())
		}
	}
	csvwriter.Flush()

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _PlateTest3Analysis(_ctx context.Context, _input *PlateTest3Input, _output *PlateTest3Output) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _PlateTest3Validation(_ctx context.Context, _input *PlateTest3Input, _output *PlateTest3Output) {

}
func _PlateTest3Run(_ctx context.Context, input *PlateTest3Input) *PlateTest3Output {
	output := &PlateTest3Output{}
	_PlateTest3Setup(_ctx, input)
	_PlateTest3Steps(_ctx, input, output)
	_PlateTest3Analysis(_ctx, input, output)
	_PlateTest3Validation(_ctx, input, output)
	return output
}

func PlateTest3RunSteps(_ctx context.Context, input *PlateTest3Input) *PlateTest3SOutput {
	soutput := &PlateTest3SOutput{}
	output := _PlateTest3Run(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func PlateTest3New() interface{} {
	return &PlateTest3Element{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &PlateTest3Input{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _PlateTest3Run(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &PlateTest3Input{},
			Out: &PlateTest3Output{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type PlateTest3Element struct {
	inject.CheckedRunner
}

type PlateTest3Input struct {
	LiquidTypes                 []string
	LiquidVolumes               []wunit.Volume
	Liquidname                  string
	OutPlates                   []string
	Startingsolutions           []*wtype.LHComponent
	TestName                    string
	WellsUsedperOutPlateInorder []int
}

type PlateTest3Output struct {
	FinalSolutions                []*wtype.LHComponent
	PlatesUsedPostRunPerPlateType []int
	Status                        string
	WellsUsedPostRunPerPlate      []int
}

type PlateTest3SOutput struct {
	Data struct {
		PlatesUsedPostRunPerPlateType []int
		Status                        string
		WellsUsedPostRunPerPlate      []int
	}
	Outputs struct {
		FinalSolutions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(component.Component{Name: "PlateTest3",
		Constructor: PlateTest3New,
		Desc: component.ComponentDesc{
			Desc: "Protocol to allow for rapid combinatorial testing of plate, liquid class combinations.\nAllows testing of effect of liquid handling changes such as offsets and liquid class changes\nIntended to be run prior to any liquid handling change before accepting pull requests.\nThe element creates an output csv file which can be filled in by the user to log observed offsets\nfor each condition\n",
			Path: "src/github.com/antha-lang/antha/antha/component/an/Utility/PlateHeightTest3.an",
			Params: []component.ParamDesc{
				{Name: "LiquidTypes", Desc: "corresponding to valid antha liquid types\n", Kind: "Parameters"},
				{Name: "LiquidVolumes", Desc: "List of volumes to test\n", Kind: "Parameters"},
				{Name: "Liquidname", Desc: "can be whatever you want to call it\n", Kind: "Parameters"},
				{Name: "OutPlates", Desc: "list of out plate types to test\n", Kind: "Parameters"},
				{Name: "Startingsolutions", Desc: "", Kind: "Inputs"},
				{Name: "TestName", Desc: "name of test e.g. branch name, date, name of project; csv file will be named after this\n", Kind: "Parameters"},
				{Name: "WellsUsedperOutPlateInorder", Desc: "optional slice of ints which should match the length and order of the OutPlates slice\n", Kind: "Parameters"},
				{Name: "FinalSolutions", Desc: "", Kind: "Outputs"},
				{Name: "PlatesUsedPostRunPerPlateType", Desc: "", Kind: "Data"},
				{Name: "Status", Desc: "", Kind: "Data"},
				{Name: "WellsUsedPostRunPerPlate", Desc: "", Kind: "Data"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
