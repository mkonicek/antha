// liquidhandling/Liquidhandler.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package liquidhandling

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/cache/plateCache"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/simulator"
	simulator_lh "github.com/antha-lang/antha/microArch/simulator/liquidhandling"
)

// the liquid handler structure defines the interface to a particular liquid handling
// platform. The structure holds the following items:
// - an LHRequest structure defining the characteristics of the platform
// - a channel for communicating with the liquid handler
// additionally three functions are defined to implement platform-specific
// implementation requirements
// in each case the LHRequest structure passed in has some additional information
// added and is then passed out. Features which are already defined (e.g. by the
// scheduler or the user) are respected as constraints and will be left unchanged.
// The three functions define
// - setup (SetupAgent): How sources are assigned to plates and plates to positions
// - layout (LayoutAgent): how experiments are assigned to outputs
//
// The general mechanism by which requests which refer to specific items as opposed to
// those which only state that an item of a particular kind is required is by the definition
// of an 'inst' tag in the request structure with a guid. If this is defined and valid
// it indicates that this item in the request (e.g. a plate, stock etc.) is a specific
// instance. If this is absent then the GUID will either be created or requested
//

// NB for flexibility we should not make the properties object part of this but rather
// send it in as an argument

type Liquidhandler struct {
	Properties      *liquidhandling.LHProperties
	FinalProperties *liquidhandling.LHProperties
	SetupAgent      func(context.Context, *LHRequest, *liquidhandling.LHProperties) error
	LayoutAgent     func(context.Context, *LHRequest, *liquidhandling.LHProperties) error
	plateIDMap      map[string]string // which plates are before / after versions
}

// initialize the liquid handling structure
func Init(properties *liquidhandling.LHProperties) *Liquidhandler {
	return &Liquidhandler{
		Properties:      properties,
		FinalProperties: properties,
		SetupAgent:      BasicSetupAgent,
		LayoutAgent:     ImprovedLayoutAgent,
		plateIDMap:      make(map[string]string),
	}
}

func (this *Liquidhandler) PlateIDMap() map[string]string {
	ret := make(map[string]string, len(this.plateIDMap))

	for k, v := range this.plateIDMap {
		ret[k] = v
	}

	return ret
}

// high-level function which requests planning and execution for an incoming set of
// solutions
func (this *Liquidhandler) MakeSolutions(ctx context.Context, request *LHRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	if err := this.Plan(ctx, request); err != nil {
		return err
	}

	if err := this.AddSetupInstructions(request); err != nil {
		return err
	}

	fmt.Println("Tip Usage Summary:")
	for _, tipEstimate := range request.TipsUsed {
		fmt.Printf("  %v\n", tipEstimate)
	}

	if err := this.Simulate(request); err != nil && !request.Options.IgnorePhysicalSimulation {
		return errors.WithMessage(err, "during physical simulation")
	}

	if err := this.Execute(request); err != nil {
		return err
	}

	// output some info on the final setup
	OutputSetup(this.Properties)

	// and after
	fmt.Println("SETUP AFTER: ")
	OutputSetup(this.FinalProperties)

	return nil
}

//AddSetupInstructions add instructions to the instruction stream to setup
//the plate layout of the machine
func (this *Liquidhandler) AddSetupInstructions(request *LHRequest) error {
	if request.Instructions == nil {
		return wtype.LHError(wtype.LH_ERR_OTHER, "Cannot execute request: no instructions")
	}

	setup_insts := this.Properties.GetSetupInstructions()
	if request.Instructions[0].Type() == liquidhandling.INI {
		request.Instructions = append(request.Instructions[:1], append(setup_insts, request.Instructions[1:]...)...)
	} else {
		request.Instructions = append(setup_insts, request.Instructions...)
	}
	return nil
}

// run the request via the physical simulator
func (this *Liquidhandler) Simulate(request *LHRequest) error {

	instructions := request.Instructions
	if len(instructions) == 0 {
		return wtype.LHError(wtype.LH_ERR_OTHER, "cannot simulate request: no instructions")
	}

	// set up the simulator with default settings
	props := this.Properties.DupKeepIDs()

	settings := simulator_lh.DefaultSimulatorSettings()

	//Make this warning less noisy since it's not really important
	settings.EnablePipetteSpeedWarning(simulator_lh.WarnOnce)
	//again, something we should fix, but not important to users to quieten
	settings.EnableAutoChannelWarning(simulator_lh.WarnOnce)
	//this is probably not even an error as liquid types are more about LHPolicies than what's actually in the well
	settings.EnableLiquidTypeWarning(simulator_lh.WarnNever)
	//disable tipbox collision. Tipboxes are narrower at the top than the bottom, so bounding box collision falsely predicts
	//collisions when when tips are picked up sequentially
	settings.EnableTipboxCollision(false)

	vlh, err := simulator_lh.NewVirtualLiquidHandler(props, settings)
	if err != nil {
		return err
	}

	triS := make([]liquidhandling.TerminalRobotInstruction, 0, len(instructions))
	for i, ins := range instructions {
		tri, ok := ins.(liquidhandling.TerminalRobotInstruction)
		if !ok {
			return fmt.Errorf("instruction %d not terminal", i)
		}
		triS = append(triS, tri)

	}

	if request.Options.PrintInstructions {
		fmt.Printf("Simulating %d instructions\n", len(instructions))
		for i, ins := range instructions {
			if request.Options.PrintInstructions {
				fmt.Printf("%d: %s\n", i, liquidhandling.InsToString(ins))
			}
		}
	}

	if err := vlh.Simulate(triS); err != nil {
		return err
	}

	//if there were no errors or warnings
	numErrors := vlh.CountErrors()
	if numErrors == 0 {
		return nil
	}

	//Output all the messages from the simulator in one logger call
	pMessage := func(n int) string {
		if n == 1 {
			return "message"
		}
		return "messages"
	}
	logLines := make([]string, 0, numErrors+1)
	logLines = append(logLines, fmt.Sprintf("showing %d %s from physical simulation:", numErrors, pMessage(numErrors)))
	//Format numbers at consistent width so messages line up
	fmtString := fmt.Sprintf("  %%%dd: simulator: %%s", 1+int(math.Floor(math.Log10(float64(numErrors)))))
	for i, err := range vlh.GetErrors() {
		logLines = append(logLines, fmt.Sprintf(fmtString, i+1, err.Error()))
	}
	fmt.Println(strings.Join(logLines, "\n"))

	//return the worst error if it's actually an error
	if simErr := vlh.GetFirstError(simulator.SeverityError); simErr != nil {
		errMsg := simErr.Error()
		if dErr, ok := simErr.(simulator_lh.DetailedLHError); ok {
			//include physical 'stack'
			errMsg += "\n\t" + strings.Replace(dErr.GetStateAtError(), "\n", "\n\t", -1)
		}
		return errors.Errorf("%s\n\tPhysical simulation can be overridden using the \"IgnorePhysicalSimulation\" configuration option.",
			errMsg)
	}

	return nil
}

// run the request via the driver
func (this *Liquidhandler) Execute(request *LHRequest) error {
	//robot setup now included in instructions

	instructions := request.Instructions

	// some timing info for the log (only) for now

	timer := this.Properties.GetTimer()
	var d time.Duration

	err := this.update_metadata()
	if err != nil {
		return err
	}

	for _, ins := range instructions {

		if request.Options.PrintInstructions {
			fmt.Println(liquidhandling.InsToString(ins))

		}
		_, ok := ins.(liquidhandling.TerminalRobotInstruction)

		if !ok {
			fmt.Printf("ERROR: Got instruction \"%s\" which is wrong type", liquidhandling.InsToString(ins))
			continue
		}

		err := ins.(liquidhandling.TerminalRobotInstruction).OutputTo(this.Properties.Driver)

		if err != nil {
			return wtype.LHError(wtype.LH_ERR_DRIV, err.Error())
		}

		// The graph view depends on the string generated in this step
		str := ""
		if ins.Type() == liquidhandling.TFR {
			mocks := liquidhandling.MockAspDsp(ins)
			for _, ii := range mocks {
				str += liquidhandling.InsToString(ii) + "\n"
			}
		} else {
			str = liquidhandling.InsToString(ins) + "\n"
		}

		request.InstructionText += str

		//fmt.Println(liquidhandling.InsToString(ins))

		if timer != nil {
			d += timer.TimeFor(ins)
		}
	}

	fmt.Printf("Total time estimate: %s\n", d.String())
	request.TimeEstimate = d.Round(time.Second).Seconds()

	return nil
}

// shrinkVolumes reduce autoallocated volumes to the amount we actually need, removing
// any unused wells or plates
func (this *Liquidhandler) shrinkVolumes(rq *LHRequest) error {

	// first, iterate through the generated instructions and count up how much
	// of each autoallocated liquid was actually used
	var lastWells []*wtype.LHWell
	vols := make(map[*wtype.LHWell]wunit.Volume)
	rawVols := make(map[*wtype.LHWell]wunit.Volume)
	useVol := func(well *wtype.LHWell, vol wunit.Volume, carry wunit.Volume) {
		v, ok := vols[well]
		r := rawVols[well]
		if !ok {
			v = wunit.NewVolume(0.0, "ul")
			r = wunit.NewVolume(0.0, "ul")
			vols[well] = v
			rawVols[well] = r
		}
		v.Add(vol)
		v.Add(carry)
		r.Add(vol)
	}
	usedPlates := make(map[*wtype.LHPlate]bool)

	for _, ins := range rq.Instructions {
		ins.Visit(liquidhandling.RobotInstructionBaseVisitor{
			HandleMove: func(ins *liquidhandling.MoveInstruction) {
				if len(ins.Pos) != len(lastWells) {
					lastWells = make([]*wtype.LHWell, len(ins.Pos))
				}
				for i, position := range ins.Pos {
					plateID := this.Properties.PosLookup[position]
					if plate, ok := this.Properties.PlateLookup[plateID].(*wtype.LHPlate); ok {
						usedPlates[plate] = true
						lastWells[i] = plate.Wellcoords[ins.Well[i]]
					}
				}
			},
			HandleAspirate: func(ins *liquidhandling.AspirateInstruction) {
				for i, lastWell := range lastWells {
					if lastWell.IsAutoallocated() && i < len(ins.Volume) {
						useVol(lastWell, ins.Volume[i], this.Properties.CarryVolume())
					}
				}
			},
			HandleTransfer: func(ins *liquidhandling.TransferInstruction) {
				for _, mtf := range ins.Transfers {
					for _, tf := range mtf.Transfers {
						if plate, ok := this.Properties.PlateLookup[tf.PltFrom].(*wtype.LHPlate); ok {
							usedPlates[plate] = true
							if well := plate.Wellcoords[tf.WellFrom]; well.IsAutoallocated() {
								useVol(well, tf.Volume, wunit.ZeroVolume())
							}
						}
					}
				}
			},
		})
	}

	getFinalWell := func(initialWell *wtype.LHWell) (*wtype.LHWell, error) {
		// assumption: plate locations don't change
		plateID := wtype.IDOf(initialWell.GetParent())
		if platePos, ok := this.Properties.PlateIDLookup[plateID]; !ok {
			return nil, errors.Errorf("couldn't find position of plate %s", plateID)
		} else if finalPlate, ok := this.FinalProperties.Plates[platePos]; !ok {
			return nil, errors.Errorf("couldn't find final plate for initial plate %s at %s", plateID, platePos)
		} else if finalWell, ok := finalPlate.WellAt(initialWell.Crds); !ok {
			return nil, errors.Errorf("couldn't find well %s in final plate at %s", initialWell.Crds.FormatA1(), platePos)
		} else {
			return finalWell, nil
		}

	}

	// second, set volumes for each autoallocated input as calculated
	for initialWell, volUsed := range vols {
		if initialWell.IsAutoallocated() {
			remainingVolume := initialWell.ResidualVolume()
			initialVolume := wunit.AddVolumes(volUsed, remainingVolume)

			if initialVolume.GreaterThan(initialWell.MaxVolume()) {
				// this logic is a bit complicated and questionable, however we need to
				// improve on how carry volumes are handled before it can be removed
				//
				// the idea is that if the total volume requested is greater than the
				// maximum the well can hold only because of carry volumes we round down
				// to the maximum the well can hold and hope for the best. This is a rare
				// case but still we should ideally be stricter
				//
				rv := rawVols[initialWell]
				if rv.LessThan(initialWell.MaxVolume()) || rv.EqualTo(initialWell.MaxVolume()) {
					// don't exceed the well maximum by a trivial amount
					initialVolume = initialWell.MaxVolume()
				} else {
					return fmt.Errorf("Error autogenerating stock %s at plate %s (type %s) well %s: Volume requested (%s) over well capacity (%s)",
						initialWell.Contents().CName, wtype.NameOf(initialWell.Plate), wtype.TypeOf(initialWell.Plate), initialWell.Crds.FormatA1(),
						initialVolume.ToString(), initialWell.MaxVolume().ToString())
				}
			}

			initialContents := initialWell.Contents().Dup()
			initialContents.SetVolume(initialVolume)
			if err := initialWell.SetContents(initialContents); err != nil {
				return err
			}

			// since we aren't yet re-generating the instructions, we need to update the final volume as well
			finalContents := initialContents.Dup()
			finalContents.SetVolume(remainingVolume)
			if finalWell, err := getFinalWell(initialWell); err != nil {
				return err
			} else if err := finalWell.SetContents(finalContents); err != nil {
				return err
			}
		}
	}

	// third, remove anything which was autoallocated but not used at all by instructions
	toRemove := make([]string, 0, len(this.Properties.Plates))
	for _, plate := range this.Properties.Plates {
		if !usedPlates[plate] && plate.AllAutoallocated() {
			toRemove = append(toRemove, plate.ID)
		} else {
			for _, initialWell := range plate.Wellcoords {
				if _, used := vols[initialWell]; !used && initialWell.IsAutoallocated() {
					initialWell.Clear()

					// as in step 2, we need to update the final well volume as well
					if finalWell, err := getFinalWell(initialWell); err != nil {
						return err
					} else {
						finalWell.Clear()
					}
				}
			}
		}
	}
	for _, removeID := range toRemove {
		if platePos, ok := this.Properties.PlateIDLookup[removeID]; !ok {
			return errors.Errorf("cannot remove plate that doesn't exist")
		} else {
			this.Properties.RemovePlateWithID(removeID)
			this.FinalProperties.RemovePlateAtPosition(platePos)
			ipo := make([]string, 0, len(rq.InputPlateOrder))
			for _, id := range rq.InputPlateOrder {
				if id != removeID {
					ipo = append(ipo, id)
				}
			}
			rq.InputPlateOrder = ipo
		}
	}

	return nil
}

// updateIDs
func (this *Liquidhandler) updateIDs() error {

	// keep track of object ID changes
	this.plateIDMap = make(map[string]string, len(this.Properties.PlateLookup))
	for id := range this.Properties.PlateIDLookup {
		if newID, err := this.FinalProperties.UpdateID(id); err != nil {
			return err
		} else {
			this.plateIDMap[id] = newID
		}
	}

	// update component IDs
	for _, plate := range this.FinalProperties.Plates {
		for _, row := range plate.Wells() {
			for _, well := range row {
				if !well.IsEmpty() {
					well.Contents().ID = wtype.GetUUID()
				}
			}
		}
	}

	return nil
}

func (this *Liquidhandler) update_metadata() error {
	if drv, ok := this.Properties.Driver.(liquidhandling.LowLevelLiquidhandlingDriver); ok {
		return drv.UpdateMetaData(this.Properties).GetError()
	}
	return nil
}

// This runs the following steps in order:
// - determine required inputs
// - request inputs	--- should be moved out
// - define robot setup
// - define output layout
// - generate the robot instructions
// - request consumables and other device setups e.g. heater setting
//
// as described above, steps only have an effect if the required inputs are
// not defined beforehand
//
// so essentially the idea is to parameterise all requests to liquid handlers
// using a Command structure called LHRequest
//
// Depending on its state of completeness, the request structure may be executable
// immediately or may need some additional definition. The purpose of the liquid
// handling service is to provide methods to invoke when parts of the request need
// further definition.
//
// when running a request we should be able to provide mechanisms for pushing requests
// back into the queue to allow them to be cached
//
// this should be OK since the LHRequest parameterises all state including instructions
// for asynchronous drivers we have to determine how far the program got before it was
// paused, which should be tricky but possible.
//

func (this *Liquidhandler) Plan(ctx context.Context, request *LHRequest) error {

	//add in a plateCache for instruction generation
	ctx = plateCache.NewContext(ctx)

	// figure out the ordering for the high level instructions
	if ichain, err := buildInstructionChain(request.LHInstructions); err != nil {
		return err
	} else {
		//sort the instructions within each link of the instruction chain
		ichain.SortInstructions(request.Options.OutputSort)
		request.InstructionChain = ichain
		request.updateWithNewLHInstructions(ichain.GetOrderedLHInstructions())
		request.OutputOrder = ichain.FlattenInstructionIDs()
	}

	if request.Options.PrintInstructions {
		fmt.Println("")
		fmt.Printf("Ordered Instructions:")
		for _, insID := range request.OutputOrder {
			fmt.Println(request.LHInstructions[insID])
		}
		request.InstructionChain.Print()
	}

	// assert we should have some instruction ordering
	if len(request.OutputOrder) == 0 {
		return fmt.Errorf("Error with instruction sorting: Have %d want %d instructions", len(request.OutputOrder), len(request.LHInstructions))
	}

	// assert that we must keep prompts and splits separate from mixes
	if err := request.InstructionChain.AssertInstructionsSeparate(); err != nil {
		return err
	}

	// make sure instructions don't share pointers to inputs or outputs
	request.LHInstructions.DupLiquids()

	if err := request.LHInstructions.AssertVolumesNonNegative(); err != nil {
		return err
	} else if err := request.LHInstructions.AssertTotalVolumesMatch(); err != nil {
		return err
	} else if err := request.LHInstructions.AssertMixResultsCorrect(); err != nil {
		return err
	} else if err := request.assertWellNotOverfilled(ctx); err != nil {
		return err
	}

	if instructions, stockconcs, err := request.solutionSetup(); err != nil {
		return errors.WithMessage(err, "during solution setup")
	} else if err := instructions.AssertVolumesNonNegative(); err != nil {
		return errors.WithMessage(err, "after solution setup")
	} else if err := instructions.AssertTotalVolumesMatch(); err != nil {
		return errors.WithMessage(err, "after solution setup")
	} else if err := instructions.AssertMixResultsCorrect(); err != nil {
		return errors.WithMessage(err, "after solution setup")
	} else {
		request.LHInstructions = instructions
		request.Stockconcs = stockconcs
	}

	// set up the mapping of the outputs
	// tried moving here to see if we can use results in fixVolumes
	if err := this.Layout(ctx, request); err != nil {
		return err
	}

	request.LHInstructions.DupLiquids()
	if err := request.LHInstructions.AssertNoPointerReuse(); err != nil {
		return errors.WithMessage(err, "failed to prevent pointer re-use")
	}

	// assert: all instructions should now be assigned specific plate IDs, types and wells
	if err := request.LHInstructions.AssertDestinationsSet(); err != nil {
		return errors.WithMessage(err, "some mix instructions missing destinations after layout")
	}

	if request.Options.FixVolumes {
		// see if volumes can be corrected
		if err := FixVolumes(request, this.Properties.CarryVolume()); err != nil {
			return err
		} else {
			if request.Options.PrintInstructions {
				fmt.Println("\nInstructions Post Volume Fix")
				for _, insID := range request.OutputOrder {
					fmt.Println(request.LHInstructions[insID])
				}
			}
		}
	}

	if err := request.LHInstructions.AssertVolumesNonNegative(); err != nil {
		return errors.WithMessage(err, "after fixing volumes")
	} else if err := request.LHInstructions.AssertTotalVolumesMatch(); err != nil {
		return errors.WithMessage(err, "after fixing volumes")
	} else if err := request.LHInstructions.AssertMixResultsCorrect(); err != nil {
		return errors.WithMessage(err, "after fixing volumes")
	} else if err := request.assertWellNotOverfilled(ctx); err != nil {
		return errors.WithMessage(err, "after fixing volumes")
	}

	// looks at liquids provided, calculates liquids required
	if inputSolutions, err := request.getInputs(this.Properties.CarryVolume()); err != nil {
		return err
	} else {
		request.InputSolutions = inputSolutions
		if request.Options.PrintInstructions {
			// print out a summary of the input solutions
			s := strings.Join(strings.Split(inputSolutions.String(), "\n"), "\n  ")
			fmt.Printf("===================\n  Input Solutions\n-------------------\n  %s\n===================\n", s)
		}
	}

	// define the input plates
	if err := request.inputPlateSetup(ctx, this.Properties.CarryVolume()); err != nil {
		return errors.WithMessage(err, "while setting up input plates")
	}

	// next we need to determine the liquid handler setup
	if err := this.Setup(ctx, request); err != nil {
		return err
	}

	// final insurance that plate names will be safe
	request.fixDuplicatePlateNames()

	// remove dummy mix-in-place instructions
	request.removeDummyInstructions()

	//set the well targets
	if err := this.addWellTargets(); err != nil {
		return err
	}

	// make the instructions for executing this request by first building the ITree root, then generating the lower level instructions
	if root, err := liquidhandling.NewITreeRoot(request.InstructionChain); err != nil {
		return err
	} else if final, err := root.Build(ctx, request.Policies(), this.Properties); err != nil {
		return err
	} else {
		request.InstructionTree = root
		request.Instructions = root.Leaves()
		this.FinalProperties = final
	}

	// tipboxes are added during the tree building, so only exist in the final state
	// copy them accross to the initial properties
	for pos, tb := range this.FinalProperties.Tipboxes {
		initialTb := tb.DupKeepIDs()
		initialTb.Refresh()
		this.Properties.AddTipBoxTo(pos, initialTb)
	}

	// revise the volumes - this makes sure the volumes requested are correct
	if err := this.shrinkVolumes(request); err != nil {
		return err
	}

	// make certain the IDs have been changed by the liquidhandling step
	if err := this.updateIDs(); err != nil {
		return err
	}

	// counts tips used in this run -- reads instructions generated above so must happen
	// after execution planning
	if estimate, err := this.countTipsUsed(request.Instructions); err != nil {
		return err
	} else {
		request.TipsUsed = estimate
	}

	// change component IDs in final state to be certain that all compoenent IDs were changed by the liquidhandling
	for _, p := range this.FinalProperties.Plates {
		for _, w := range p.Wellcoords {
			if w.IsUserAllocated() {
				w.WContents.ID = wtype.GetUUID()
			}
		}
	}

	if err := this.updateComponentNames(request); err != nil {
		return err
	}

	return request.assertNoTemporaryPlates(ctx)
}

// define which labware to use
func (this *Liquidhandler) GetPlates(ctx context.Context, plates map[string]*wtype.Plate, major_layouts map[int][]string, ptype *wtype.Plate) (map[string]*wtype.Plate, error) {
	if plates == nil {
		plates = make(map[string]*wtype.Plate, len(major_layouts))

		// assign new plates
		for i := 0; i < len(major_layouts); i++ {
			//newplate := wtype.New_Plate(ptype)
			newplate, err := inventory.NewPlate(ctx, ptype.Type)
			if err != nil {
				return nil, err
			}
			plates[newplate.ID] = newplate
		}
	}

	// we should know how many plates we need
	for k, plate := range plates {
		//if plate.Inst == "" {
		//stockrequest := execution.GetContext().StockMgr.RequestStock(makePlateStockRequest(plate))
		//plate.Inst = stockrequest["inst"].(string)
		//}

		plates[k] = plate
	}

	return plates, nil
}

// generate setup for the robot
func (this *Liquidhandler) Setup(ctx context.Context, request *LHRequest) error {
	// assign the plates to positions
	// this needs to be parameterizable
	return this.SetupAgent(ctx, request, this.Properties)
}

// generate the output layout
func (this *Liquidhandler) Layout(ctx context.Context, request *LHRequest) error {
	// assign the results to destinations
	// again needs to be parameterized

	return this.LayoutAgent(ctx, request, this.Properties)
}

func OutputSetup(robot *liquidhandling.LHProperties) {
	fmt.Println("DECK SETUP INFO")
	fmt.Println("Tipboxes: ")

	for k, v := range robot.Tipboxes {
		fmt.Printf("%s %s: %s\n", k, robot.PlateIDLookup[k], v.Type)
	}

	fmt.Println("Plates:")

	for k, v := range robot.Plates {

		fmt.Printf("%s %s: %s %s\n", k, robot.PlateIDLookup[k], v.PlateName, v.Type)

		//TODO Deprecate
		if strings.Contains(v.GetName(), "Input") {
			_, err := wtype.AutoExportPlateCSV(v.GetName()+".csv", v)
			if err != nil {
				fmt.Printf("export plate csv (deprecated): %s\n", err.Error())
			}
		}

		v.OutputLayout()
	}

	fmt.Println("Tipwastes: ")

	for k, v := range robot.Tipwastes {
		fmt.Printf("%s %s: %s capacity %d\n", k, robot.PlateIDLookup[k], v.Type, v.Capacity)
	}

}

// updateComponentNames the name (CName) of an output liquid can be overridden by the
// LHInstruction which generated it, so this function updates the name of each liquid
// to be equal to that which was set in the Mix instruction which created it
func (lh *Liquidhandler) updateComponentNames(rq *LHRequest) error {

	// get the instructions in the ordere that they happen
	orderedInstructions, err := rq.GetOrderedLHInstructions()
	if err != nil {
		return err
	}

	for _, inst := range orderedInstructions {
		// only mixes update component names
		if inst.Type != wtype.LHIMIX {
			continue
		}

		// find the well from the output's platelocation
		if pl := inst.Outputs[0].PlateLocation(); pl.IsZero() {
			return wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprintf("malformed output component location \"%s\" for instruction %s", inst.Outputs[0].Loc, inst.ID))
		} else if newid, ok := lh.plateIDMap[pl.ID]; !ok {
			return wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprintf("No output plate mapped to %s", pl.ID))
		} else if ip, ok := lh.FinalProperties.PlateLookup[newid]; !ok {
			return wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprintf("No output plate %s", newid))
		} else if p, ok := ip.(*wtype.Plate); !ok {
			return wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprintf("Got %s, should have *wtype.LHPlate", reflect.TypeOf(ip)))
		} else if well, ok := p.WellAt(pl.Coords); !ok {
			return wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprintf("No well %s on plate %s", pl.Coords.FormatA1(), pl.ID))
		} else {
			well.WContents.CName = inst.Outputs[0].CName
		}
	}

	return nil
}

//addWellTargets for all the adaptors and plates available
func (lh *Liquidhandler) addWellTargets() error {
	for _, head := range lh.Properties.Heads {
		for _, plate := range lh.Properties.Plates {
			addWellTargetsPlate(head.Adaptor, plate)
		}
		for _, plate := range lh.Properties.Wastes {
			addWellTargetsPlate(head.Adaptor, plate)
		}
		for _, plate := range lh.Properties.Washes {
			addWellTargetsPlate(head.Adaptor, plate)
		}
		for _, tipWaste := range lh.Properties.Tipwastes {
			addWellTargetsTipWaste(head.Adaptor, tipWaste)
		}
	}
	return nil
}

const adaptorSpacing = 9.0

func addWellTargetsPlate(adaptor *wtype.LHAdaptor, plate *wtype.Plate) {
	if adaptor.Manufacturer != "Gilson" {
		fmt.Println("Not adding well target data for non-gilson adaptor")
		return
	}

	if !plate.AreWellTargetsEnabled(adaptor.Params.Multi, adaptorSpacing) {
		if plate.NRows() < 8 {
			//declare special so that the driver knows not to expect well targets
			plate.DeclareSpecial()
		}
		return
	}

	//channelPositions should come from the adaptor
	channelPositions := make([]wtype.Coordinates3D, 0, adaptor.Params.Multi)
	for i := 0; i < adaptor.Params.Multi; i++ {
		channelPositions = append(channelPositions, wtype.Coordinates3D{Y: float64(i) * adaptorSpacing})
	}

	//ystart and count should come from some geometric calculation between channelPositions and well size
	ystart, count := getWellTargetYStart(plate.NRows())

	targets := make([]wtype.Coordinates3D, count)
	copy(targets, channelPositions)

	for i := 0; i < count; i++ {
		targets[i].Y += ystart
	}

	plate.Welltype.SetWellTargets(adaptor.Name, targets)
}

func addWellTargetsTipWaste(adaptor *wtype.LHAdaptor, waste *wtype.LHTipwaste) {
	// this may vary in future but for now we just need to add the following eight entries:
	ystart := -31.5
	yinc := 9.0
	targets := make([]wtype.Coordinates3D, 0, adaptor.Params.Multi)
	for i := 0; i < adaptor.Params.Multi; i++ {
		targets = append(targets, wtype.Coordinates3D{Y: ystart + float64(i)*yinc})
	}

	waste.AsWell.SetWellTargets(adaptor.Name, targets)
}

//getWellTargetYStart pmdriver mapping from number of y wells to number of tips and y start
func getWellTargetYStart(wy int) (float64, int) {
	// this is pretty simple to start with
	// OK but there are a few issues with special plate types

	switch wy {
	case 1:
		return -31.5, 8
	case 2:
		return -13.5, 4
	case 3:
		return -9.0, 3 // check
	case 4:
		return -4.5, 2
	case 5:
		return 0.0, 1
	case 6:
		return 0.0, 1
	case 7:
		return 0.0, 1
	}

	return 0.0, 0
}
