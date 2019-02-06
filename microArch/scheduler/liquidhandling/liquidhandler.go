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
	"github.com/antha-lang/antha/antha/anthalib/wutil"
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
	SetupAgent      func(context.Context, *LHRequest, *liquidhandling.LHProperties) (*LHRequest, error)
	LayoutAgent     func(context.Context, *LHRequest, *liquidhandling.LHProperties) (*LHRequest, error)
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

	var simulator *simulator_lh.VirtualLiquidHandler
	var err error

	err, simulator = this.Simulate(request)

	if err != nil && !request.Options.IgnorePhysicalSimulation {
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

	if !request.Options.IgnoreLogicalSimulation {
		// compare the declared output to what the simulator has found
		errs := simulator.CompareStateToDeclaredOutput(this.FinalProperties)

		if len(errs) != 0 {
			line := ""
			for _, err := range errs {
				line += err.Error() + "\n"
			}

			return fmt.Errorf("%s", line)
		}
	}

	return nil
}

//AddSetupInstructions add instructions to the instruction stream to setup
//the plate layout of the machine
func (this *Liquidhandler) AddSetupInstructions(request *LHRequest) error {
	if request.Instructions == nil {
		return wtype.LHError(wtype.LH_ERR_OTHER, "Cannot execute request: no instructions")
	}

	setup_insts := this.get_setup_instructions(request)
	if request.Instructions[0].Type() == liquidhandling.INI {
		request.Instructions = append(request.Instructions[:1], append(setup_insts, request.Instructions[1:]...)...)
	} else {
		request.Instructions = append(setup_insts, request.Instructions...)
	}
	return nil
}

// run the request via the physical simulator
func (this *Liquidhandler) Simulate(request *LHRequest) (error, *simulator_lh.VirtualLiquidHandler) {

	instructions := (*request).Instructions
	if instructions == nil {
		return wtype.LHError(wtype.LH_ERR_OTHER, "cannot simulate request: no instructions"), nil
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
		return err, vlh
	}

	triS := make([]liquidhandling.TerminalRobotInstruction, 0, len(instructions))
	for i, ins := range instructions {
		tri, ok := ins.(liquidhandling.TerminalRobotInstruction)
		if !ok {
			return fmt.Errorf("instruction %d not terminal", i), vlh
		}
		triS = append(triS, tri)

	}

	if request.Options.PrintInstructions {
		fmt.Printf("Simulating %d instructions\n", len(instructions))
		for i, ins := range instructions {
			if (*request).Options.PrintInstructions {
				fmt.Printf("%d: %s\n", i, liquidhandling.InsToString(ins))
			}
		}
	}

	if err := vlh.Simulate(triS); err != nil {
		return err, vlh
	}

	//if there were no errors or warnings
	numErrors := vlh.CountErrors()
	if numErrors == 0 {
		return nil, vlh
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
			errMsg), vlh
	}

	return nil, vlh
}

// run the request via the driver
func (this *Liquidhandler) Execute(request *LHRequest) error {
	//robot setup now included in instructions

	instructions := (*request).Instructions

	// some timing info for the log (only) for now

	timer := this.Properties.GetTimer()
	var d time.Duration

	err := this.update_metadata(request)
	if err != nil {
		return err
	}

	for _, ins := range instructions {

		if (*request).Options.PrintInstructions {
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

func (this *Liquidhandler) reviseVolumes(rq *LHRequest) error {
	// XXX -- HARD CODE 8 here
	lastPlate := make([]string, 8)
	lastWell := make([]string, 8)

	vols := make(map[string]map[string]wunit.Volume)

	rawvols := make(map[string]map[string]wunit.Volume)

	for _, ins := range rq.Instructions {
		ins.Visit(liquidhandling.RobotInstructionBaseVisitor{
			HandleMove: func(ins *liquidhandling.MoveInstruction) {
				lastPlate = make([]string, 8)

				for i, p := range ins.Pos {
					lastPlate[i] = this.Properties.PosLookup[p]
				}

				lastWell = ins.Well
			},
			HandleAspirate: func(ins *liquidhandling.AspirateInstruction) {
				for i := range lastPlate {
					if i >= len(lastWell) {
						break
					}
					lp := lastPlate[i]
					lw := lastWell[i]

					if lp == "" {
						continue
					}

					ppp := this.Properties.PlateLookup[lp].(*wtype.Plate)

					lwl := ppp.Wellcoords[lw]

					if !lwl.IsAutoallocated() {
						continue
					}

					_, ok := vols[lp]

					if !ok {
						vols[lp] = make(map[string]wunit.Volume)
						rawvols[lp] = make(map[string]wunit.Volume)
					}

					v, ok := vols[lp][lw]

					if !ok {
						v = wunit.NewVolume(0.0, "ul")
						vols[lp][lw] = v
						rawvols[lp][lw] = v.Dup()
					}
					//v.Add(ins.Volume[i])

					insvols := ins.Volume
					v.Add(insvols[i])
					v.Add(rq.CarryVolume)

					rawvols[lp][lw].Add(insvols[i])
				}
			},
			HandleTransfer: func(ins *liquidhandling.TransferInstruction) {
				for _, mtf := range ins.Transfers {
					for _, tf := range mtf.Transfers {
						lpos, lw := tf.PltFrom, tf.WellFrom

						lp := this.Properties.PosLookup[lpos]
						ppp := this.Properties.PlateLookup[lp].(*wtype.Plate)
						lwl := ppp.Wellcoords[lw]

						if !lwl.IsAutoallocated() {
							continue
						}

						_, ok := vols[lp]

						if !ok {
							vols[lp] = make(map[string]wunit.Volume)
						}

						v, ok := vols[lp][lw]

						if !ok {
							v = wunit.NewVolume(0.0, "ul")
							vols[lp][lw] = v
						}
						//v.Add(ins.Volume[i])

						v.Add(tf.Volume)
					}
				}
			},
		})
	}

	// apply evaporation
	for _, vc := range rq.Evaps {
		loctox := strings.Split(vc.Location, ":")

		// ignore anything where the location isn't properly set

		if len(loctox) < 2 {
			continue
		}

		plateID := loctox[0]
		wellcrds := loctox[1]

		wellmap, ok := vols[plateID]

		if !ok {
			continue
		}

		vol := wellmap[wellcrds]
		vol.Add(vc.Volume)
	}

	// now go through and set the plates up appropriately

	for plateID, wellmap := range vols {
		plate, ok := this.FinalProperties.Plates[this.Properties.PlateIDLookup[plateID]]
		plate2 := this.Properties.Plates[this.Properties.PlateIDLookup[plateID]]

		if !ok {
			err := wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprint("NO SUCH PLATE: ", plateID))
			return err
		}

		// what's it like here?

		for crd, unroundedvol := range wellmap {
			rv, _ := wutil.Roundto(unroundedvol.RawValue(), 1)
			vol := wunit.NewVolume(rv, unroundedvol.Unit().PrefixedSymbol())
			well := plate.Wellcoords[crd]
			well2 := plate2.Wellcoords[crd]

			// this logic is a bit complicated and questionable, however we need to
			// improve on how carry volumes are handled before it can be removed
			//
			// the idea is that if the total volume requested is greater than the
			// maximum the well can hold only because of carry volumes we round down
			// to the maximum the well can hold and hope for the best. This is a rare
			// case but still we should ideally be stricter
			//
			if well.IsAutoallocated() {
				vol.Add(well.ResidualVolume())

				if vol.GreaterThan(well.MaxVolume()) {
					rv := rawvols[plateID][crd]

					if rv.LessThan(well.MaxVolume()) || rv.EqualTo(well.MaxVolume()) {
						// don't exceed the well maximum by a trivial amount
						vol = well.MaxVolume()
					} else {
						return fmt.Errorf("Error autogenerating stock %s at plate %s (type %s) well %s: Volume requested (%s) over well capacity (%s)", well2.Contents().CName, plate2.Name(), plate2.Type, crd, vol.ToString(), well.MaxVolume().ToString())
					}
				}

				well2Contents := well2.Contents().Dup()
				well2Contents.SetVolume(vol)
				err := well2.SetContents(well2Contents)
				if err != nil {
					return err
				}

				wellContents := well.Contents().Dup()
				wellContents.SetVolume(well.ResidualVolume())
				wellContents.ID = wtype.GetUUID()
				err = well.SetContents(wellContents)
				if err != nil {
					return err
				}

				well.DeclareNotTemporary()
				well2.DeclareNotTemporary()
			}
		}
	}

	// finally get rid of any temporary stuff

	this.Properties.RemoveUnusedAutoallocatedComponents()
	this.FinalProperties.RemoveUnusedAutoallocatedComponents()

	pidm := make(map[string]string, len(this.Properties.Plates))
	for pos := range this.Properties.Plates {
		p1, ok1 := this.Properties.Plates[pos]
		p2, ok2 := this.FinalProperties.Plates[pos]

		if (!ok1 && ok2) || (ok1 && !ok2) {

			if ok1 {
				fmt.Println("BEFORE HAS: ", p1)
			}

			if ok2 {
				fmt.Println("AFTER  HAS: ", p2)
			}

			return (wtype.LHError(8, fmt.Sprintf("Plate disappeared from position %s", pos)))
		}

		if !(ok1 && ok2) {
			continue
		}

		this.plateIDMap[p1.ID] = p2.ID
		pidm[p2.ID] = p1.ID
	}

	// this is many shades of wrong but likely to save us a lot of time
	for _, pos := range this.Properties.InputSearchPreferences() {
		p1, ok1 := this.Properties.Plates[pos]
		p2, ok2 := this.FinalProperties.Plates[pos]

		if ok1 && ok2 {
			for _, wa := range p1.Cols {
				for _, w := range wa {
					// copy the outputs to the correct side
					// and remove the outputs from the initial state
					if !w.IsEmpty() {
						w2, ok := p2.Wellcoords[w.Crds.FormatA1()]
						if ok {
							// there's no strict separation between outputs and
							// inputs here
							if w.IsAutoallocated() {
								continue
							} else if w.IsUserAllocated() {
								// swap old and new
								c := w.WContents.Dup()
								w.Clear()
								c2 := w2.WContents.Dup()
								w2.Clear()
								err := w.AddComponent(c2)
								if err != nil {
									return wtype.LHError(wtype.LH_ERR_VOL, fmt.Sprintf("Scheduler : %s", err.Error()))
								}
								err = w2.AddComponent(c)
								if err != nil {
									return wtype.LHError(wtype.LH_ERR_VOL, fmt.Sprintf("Scheduler : %s", err.Error()))
								}
							} else {
								// replace
								w2.Clear()
								err := w2.AddComponent(w.Contents())
								if err != nil {
									return wtype.LHError(wtype.LH_ERR_VOL, fmt.Sprintf("Scheduler : %s", err.Error()))
								}
								w.Clear()
							}
						}
					}
				}

			}

			//fmt.Println(p2, " ", p1)
			//fmt.Println("Plate ID Map: ", p2.ID, " --> ", p1.ID)

			//	this.plateIDMap[p2.ID] = p1.ID
		}
	}

	// all done

	return nil
}

func (this *Liquidhandler) get_setup_instructions(rq *LHRequest) []liquidhandling.TerminalRobotInstruction {
	instructions := make([]liquidhandling.TerminalRobotInstruction, 0, 1+len(this.Properties.PosLookup))

	//first instruction is always to remove all plates
	instructions = append(instructions, liquidhandling.NewRemoveAllPlatesInstruction())

	for position, plateid := range this.Properties.PosLookup {
		if plateid == "" {
			continue
		}
		plate := this.Properties.PlateLookup[plateid]
		name := plate.(wtype.Named).GetName()

		ins := liquidhandling.NewAddPlateToInstruction(position, name, plate)

		instructions = append(instructions, ins)
	}
	return instructions
}

func (this *Liquidhandler) update_metadata(rq *LHRequest) error {
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
	if rq, err := this.Layout(ctx, request); err != nil {
		return err
	} else {
		request = rq
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
		if rq, err := FixVolumes(request); err != nil {
			return err
		} else {
			request = rq
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
	if inputSolutions, err := request.getInputs(); err != nil {
		return err
	} else {
		request.InputSolutions = inputSolutions
	}

	// define the input plates
	if err := request.inputPlateSetup(ctx); err != nil {
		return errors.WithMessage(err, "while setting up input plates")
	}

	// next we need to determine the liquid handler setup
	if rq, err := this.Setup(ctx, request); err != nil {
		return err
	} else {
		request = rq
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
	// nb. there is significant potential for confusion here:
	//    root.Generate(..., props LHProperties) is *destructive of state*, and leaves it's argument in the final state
	//    therefore from here until reviseVolumes is called,
	//      > this.Properties contains the final properties
	//      > this.FinalProperties contains the initial properties
	//    which cannot be changed until reviseVolumes is refactored
	this.FinalProperties = this.Properties.Dup()
	if root, err := liquidhandling.NewITreeRoot(request.InstructionChain); err != nil {
		return err
	} else if err := root.Generate(ctx, request.Policies(), this.Properties); err != nil {
		return err
	} else if tri, err := root.Leaves(); err != nil {
		return err
	} else {
		request.InstructionTree = root
		request.Instructions = tri
	}

	// counts tips used in this run -- reads instructions generated above so must happen
	// after execution planning
	if estimate, err := this.countTipsUsed(request.Instructions); err != nil {
		return err
	} else {
		request.TipsUsed = estimate
	}

	// Ensures tip boxes and wastes are correct for initial and final robot states
	this.refreshTipboxesTipwastes()

	// revise the volumes - this makes sure the volumes requested are correct
	if err := this.reviseVolumes(request); err != nil {
		return err
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
func (this *Liquidhandler) Setup(ctx context.Context, request *LHRequest) (*LHRequest, error) {
	// assign the plates to positions
	// this needs to be parameterizable
	return this.SetupAgent(ctx, request, this.Properties)
}

// generate the output layout
func (this *Liquidhandler) Layout(ctx context.Context, request *LHRequest) (*LHRequest, error) {
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
	channelPositions := make([]wtype.Coordinates, 0, adaptor.Params.Multi)
	for i := 0; i < adaptor.Params.Multi; i++ {
		channelPositions = append(channelPositions, wtype.Coordinates{Y: float64(i) * adaptorSpacing})
	}

	//ystart and count should come from some geometric calculation between channelPositions and well size
	ystart, count := getWellTargetYStart(plate.NRows())

	targets := make([]wtype.Coordinates, count)
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
	targets := make([]wtype.Coordinates, 0, adaptor.Params.Multi)
	for i := 0; i < adaptor.Params.Multi; i++ {
		targets = append(targets, wtype.Coordinates{Y: ystart + float64(i)*yinc})
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
