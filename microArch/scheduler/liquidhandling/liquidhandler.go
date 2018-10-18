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
	"github.com/pkg/errors"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/cache"
	"github.com/antha-lang/antha/inventory/cache/plateCache"
	"github.com/antha-lang/antha/microArch/driver"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/logger"
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
// - execution (ExecutionPlanner): generates instructions to implement the required plan
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
	Properties       *liquidhandling.LHProperties
	FinalProperties  *liquidhandling.LHProperties
	SetupAgent       func(context.Context, *LHRequest, *liquidhandling.LHProperties) (*LHRequest, error)
	LayoutAgent      func(context.Context, *LHRequest, *liquidhandling.LHProperties) (*LHRequest, error)
	ExecutionPlanner func(context.Context, *LHRequest, *liquidhandling.LHProperties) (*LHRequest, error)
	PolicyManager    *LHPolicyManager
	plateIDMap       map[string]string // which plates are before / after versions
}

// initialize the liquid handling structure
func Init(properties *liquidhandling.LHProperties) *Liquidhandler {
	lh := Liquidhandler{}
	lh.SetupAgent = BasicSetupAgent
	lh.LayoutAgent = ImprovedLayoutAgent
	//lh.ExecutionPlanner = ImprovedExecutionPlanner
	lh.ExecutionPlanner = ExecutionPlanner3
	lh.Properties = properties
	lh.FinalProperties = properties
	lh.plateIDMap = make(map[string]string)
	return &lh
}

func (this *Liquidhandler) PlateIDMap() map[string]string {
	ret := make(map[string]string, len(this.plateIDMap))

	for k, v := range this.plateIDMap {
		ret[k] = v
	}

	return ret
}

// catch errors early
func ValidateRequest(request *LHRequest) error {
	if len(request.LHInstructions) == 0 {
		return wtype.LHError(wtype.LH_ERR_OTHER, "Nil plan requested: no Mix Instructions present")
	}

	// no component can have all three of Conc, Vol and TVol set to 0:

	for _, ins := range request.LHInstructions {
		// the check below makes sense only for mixes
		if ins.Type != wtype.LHIMIX {
			continue
		}
		for i, cmp := range ins.Inputs {
			if cmp.Vol == 0.0 && cmp.Conc == 0.0 && cmp.Tvol == 0.0 {
				errstr := fmt.Sprintf("Nil mix (no volume, concentration or total volume) requested: %d : ", i)

				for j := 0; j < len(ins.Inputs); j++ {
					ss := ins.Inputs[i].CName
					if j == i {
						ss = strings.ToUpper(ss)
					}

					if j != len(ins.Inputs)-1 {
						ss += ", "
					}

					errstr += ss
				}
				return wtype.LHError(wtype.LH_ERR_OTHER, errstr)
			}
		}
	}
	return nil
}

// high-level function which requests planning and execution for an incoming set of
// solutions
func (this *Liquidhandler) MakeSolutions(ctx context.Context, request *LHRequest) error {
	err := ValidateRequest(request)
	if err != nil {
		return err
	}

	err = this.Plan(ctx, request)
	if err != nil {
		return err
	}

	err = this.AddSetupInstructions(request)
	if err != nil {
		return err
	}

	fmt.Println("Tip Usage Summary:")
	for _, tipEstimate := range request.TipsUsed {
		fmt.Printf("  %v\n", tipEstimate)
	}

	err = this.Simulate(request)
	if err != nil && !request.Options.IgnorePhysicalSimulation {
		return errors.WithMessage(err, "during physical simulation")
	}

	err = this.Execute(request)
	if err != nil {
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

	setup_insts := this.get_setup_instructions(request)
	if request.Instructions[0].Type() == liquidhandling.INI {
		request.Instructions = append(request.Instructions[:1], append(setup_insts, request.Instructions[1:]...)...)
	} else {
		request.Instructions = append(setup_insts, request.Instructions...)
	}
	return nil
}

// run the request via the physical simulator
func (this *Liquidhandler) Simulate(request *LHRequest) error {

	instructions := (*request).Instructions
	if instructions == nil {
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
			fmt.Printf("%d: %s\n", i, liquidhandling.InsToString(ins))
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
	logger.Info(strings.Join(logLines, "\n"))

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

	logger.Debug(fmt.Sprintf("Total time estimate: %s", d.String()))
	request.TimeEstimate = d.Round(time.Second).Seconds()

	return nil
}

// shrinkVolumes reduce the autoallocated input volumes to match what was actually used
func (this *Liquidhandler) shrinkVolumes(rq *LHRequest) error {

	// first, iterate through the generated instructions and count up how much
	// of each autoallocated liquid was actually used
	var lastWells []*wtype.LHWell
	vols := make(map[*wtype.LHWell]wunit.Volume)
	usedPlates := make(map[*wtype.LHPlate]bool)
	for _, ins := range rq.Instructions {
		ins.Visit(liquidhandling.RobotInstructionBaseVisitor{
			HandleMove: func(ins *liquidhandling.MoveInstruction) {
				if len(ins.Plt) != len(lastWells) {
					lastWells = make([]*wtype.LHWell, len(ins.Plt))
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
						v, ok := vols[lastWell]
						if !ok {
							v = wunit.NewVolume(0.0, "ul")
							vols[lastWell] = v
						}
						v.Add(ins.Volume[i])
						v.Add(rq.CarryVolume)
					}
				}
			},
			HandleTransfer: func(ins *liquidhandling.TransferInstruction) {
				for _, mtf := range ins.Transfers {
					for _, tf := range mtf.Transfers {
						if plate, ok := this.Properties.PlateLookup[tf.PltFrom].(*wtype.LHPlate); ok {
							usedPlates[plate] = true
							if well := plate.Wellcoords[tf.WellFrom]; well.IsAutoallocated() {
								v, ok := vols[well]
								if !ok {
									v = wunit.NewVolume(0.0, "ul")
									vols[well] = v
								}
								v.Add(tf.Volume)
							}
						}
					}
				}
			},
		})
	}

	// second, apply pre-calculated evaporation volumes to the count
	for _, vc := range rq.Evaps {
		// ignore anything where the location isn't properly set ("<plateID>:<WellCoords>")
		if loctox := strings.Split(vc.Location, ":"); len(loctox) == 2 {
			plateID, wellCoords := loctox[0], loctox[1]

			if plate, ok := this.Properties.PlateLookup[plateID].(*wtype.LHPlate); ok {
				well := plate.Wellcoords[wellCoords]
				vols[well].IncrBy(vc.Volume) // nolint - (nill).IncrBy is noop
			}
		}
	}

	// third, set volumes for each autoallocated input as calculated
	for initialWell, volUsed := range vols {
		if initialWell.IsAutoallocated() {
			volUsed.IncrBy(initialWell.ResidualVolume()) // nolint - volumes are always compatible
			volUsed.DecrBy(rq.CarryVolume)               // nolint - final carry volume from well comes from residual

			initialContents := initialWell.Contents().Dup()
			initialContents.SetVolume(volUsed)
			if err := initialWell.SetContents(initialContents); err != nil {
				return err
			}

			initialWell.DeclareNotTemporary()
		}
	}

	// finally, remove anything which was autoallocated but not used at all by instructions
	toRemove := make([]string, 0, len(this.Properties.Plates))
	for _, plate := range this.Properties.Plates {
		if !usedPlates[plate] {
			toRemove = append(toRemove, plate.ID)
		} else {
			for _, well := range plate.Wellcoords {
				if _, used := vols[well]; !used && well.IsAutoallocated() {
					well.Clear()
				}
			}
		}
	}
	for _, id := range toRemove {
		this.Properties.RemovePlateWithID(id)
	}

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

	if _, ok := this.Properties.Driver.(liquidhandling.LowLevelLiquidhandlingDriver); ok {
		stat := this.Properties.Driver.(liquidhandling.LowLevelLiquidhandlingDriver).UpdateMetaData(this.Properties)
		if stat.Errorcode == driver.ERR {
			return wtype.LHError(wtype.LH_ERR_DRIV, stat.Msg)
		}
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

//assertVolumesNonNegative tests that the volumes within the LHRequest are zero or positive
func assertVolumesNonNegative(request *LHRequest) error {
	for _, ins := range request.LHInstructions {
		if ins.Type != wtype.LHIMIX {
			continue
		}

		for _, cmp := range ins.Inputs {
			if cmp.Volume().LessThan(wunit.ZeroVolume()) {
				return wtype.LHErrorf(wtype.LH_ERR_VOL, "negative volume for component \"%s\" in instruction:\n%s", cmp.CName, ins.Summarize(1))
			}
		}
	}
	return nil
}

//assertTotalVolumesMatch checks that component total volumes are all the same in mix instructions
func assertTotalVolumesMatch(request *LHRequest) error {
	for _, ins := range request.LHInstructions {
		if ins.Type != wtype.LHIMIX {
			continue
		}

		totalVolume := wunit.ZeroVolume()

		for _, cmp := range ins.Inputs {
			if tV := cmp.TotalVolume(); !tV.IsZero() {
				if !totalVolume.IsZero() && !tV.EqualTo(totalVolume) {
					return wtype.LHErrorf(wtype.LH_ERR_VOL, "multiple distinct total volumes specified in instruction:\n%s", ins.Summarize(1))
				}
				totalVolume = tV
			}
		}
	}
	return nil
}

//assertMixResultsCorrect checks that volumes of the mix result matches either the sum of the input, or the total volume if specified
func assertMixResultsCorrect(request *LHRequest) error {
	for _, ins := range request.LHInstructions {
		if ins.Type != wtype.LHIMIX {
			continue
		}

		totalVolume := wunit.ZeroVolume()
		volumeSum := wunit.ZeroVolume()

		for _, cmp := range ins.Inputs {
			if tV := cmp.TotalVolume(); !tV.IsZero() {
				totalVolume = tV
			} else if v := cmp.Volume(); !v.IsZero() {
				volumeSum.Add(v)
			}
		}

		if len(ins.Outputs) != 1 {
			return wtype.LHErrorf(wtype.LH_ERR_DIRE, "mix instruction has %d results specified, expecting one at instruction:\n%s",
				len(ins.Outputs), ins.Summarize(1))
		}

		resultVolume := ins.Outputs[0].Volume()

		if !totalVolume.IsZero() && !totalVolume.EqualTo(resultVolume) {
			return wtype.LHErrorf(wtype.LH_ERR_VOL, "total volume (%v) does not match resulting volume (%v) for instruction:\n%s",
				totalVolume, resultVolume, ins.Summarize(1))
		} else if totalVolume.IsZero() && !volumeSum.EqualTo(resultVolume) {
			return wtype.LHErrorf(wtype.LH_ERR_VOL, "sum of requested volumes (%v) does not match result volume (%v) for instruction:\n%s",
				volumeSum, resultVolume, ins.Summarize(1))
		}
	}
	return nil
}

//assertWellNotOverfilled checks that mix instructions aren't going to overfill the wells when a plate is specified
//assumes assertMixResultsCorrect returns nil
func assertWellNotOverfilled(ctx context.Context, request *LHRequest) error {
	for _, ins := range request.LHInstructions {
		if ins.Type != wtype.LHIMIX {
			continue
		}

		resultVolume := ins.Outputs[0].Volume()

		var plate *wtype.Plate
		if ins.OutPlate != nil {
			plate = ins.OutPlate
		} else if ins.PlateID != "" {
			if p, ok := request.GetPlate(ins.PlateID); !ok {
				continue
			} else {
				plate = p
			}
		} else if ins.Platetype != "" {
			if p, err := inventory.NewPlate(ctx, ins.Platetype); err != nil {
				continue
			} else {
				plate = p
			}
		} else {
			//couldn't find an appropriate plate
			continue
		}

		if maxVol := plate.Welltype.MaxVolume(); maxVol.LessThan(resultVolume) {
			//ignore if this is just numerical precision (#campainforintegervolume)
			delta := wunit.SubtractVolumes(resultVolume, maxVol)
			if delta.IsZero() {
				continue
			}
			return wtype.LHErrorf(wtype.LH_ERR_VOL, "volume of resulting mix (%v) exceeds the well maximum (%v) for instruction:\n%s",
				resultVolume, maxVol, ins.Summarize(1))
		}
	}
	return nil
}

func checkDestinationSanity(request *LHRequest) {
	for _, ins := range request.LHInstructions {
		// non-mix instructions are fine
		if ins.Type != wtype.LHIMIX {
			continue
		}

		if ins.PlateID == "" || ins.Platetype == "" || ins.Welladdress == "" {
			found := fmt.Sprintln("INS ", ins, " NOT WELL FORMED: HAS PlateID ", ins.PlateID != "", " HAS platetype ", ins.Platetype != "", " HAS WELLADDRESS ", ins.Welladdress != "")
			panic(fmt.Errorf("After layout all mix instructions must have plate IDs, plate types and well addresses, Found: \n %s", found))
		}
	}
}

func anotherSanityCheck(request *LHRequest) {
	p := map[*wtype.Liquid]*wtype.LHInstruction{}

	for _, ins := range request.LHInstructions {
		// we must not share pointers

		for _, c := range ins.Inputs {
			ins2, ok := p[c]
			if ok {
				panic(fmt.Sprintf("POINTER REUSE: Instructions %s %s for component %s %s", ins.ID, ins2.ID, c.ID, c.CName))
			}

			p[c] = ins
		}

		ins2, ok := p[ins.Outputs[0]]

		if ok {
			panic(fmt.Sprintf("POINTER REUSE: Instructions %s %s for component %s %s", ins.ID, ins2.ID, ins.Outputs[0].ID, ins.Outputs[0].CName))
		}

		p[ins.Outputs[0]] = ins
	}
}

func forceSanity(request *LHRequest) {
	for _, ins := range request.LHInstructions {
		for i := 0; i < len(ins.Inputs); i++ {
			ins.Inputs[i] = ins.Inputs[i].Dup()
		}

		ins.Outputs[0] = ins.Outputs[0].Dup()
	}
}

//check that none of the plates we're returning came from the cache
func assertNoTemporaryPlates(ctx context.Context, request *LHRequest) error {

	for id, plate := range request.Plates {
		if cache.IsFromCache(ctx, plate) {
			return wtype.LHErrorf(wtype.LH_ERR_DIRE, "found a temporary plate (id=%s) being returned in the request", id)
		}
	}

	return nil
}

func (this *Liquidhandler) Plan(ctx context.Context, request *LHRequest) error {

	//add in a plateCache for instruction generation
	ctx = plateCache.NewContext(ctx)

	// figure out the ordering for the high level instructions
	if ichain, err := buildInstructionChain(request.LHInstructions); err != nil {
		return err
	} else {
		//sort the instructions within each link of the instruction chain
		ichain.sortInstructions(request.Options.OutputSort)
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
	if err := request.InstructionChain.assertInstructionsSeparate(); err != nil {
		return err
	}

	forceSanity(request)
	// convert requests to volumes and determine required stock concentrations

	if err := assertVolumesNonNegative(request); err != nil {
		return err
	}
	if err := assertTotalVolumesMatch(request); err != nil {
		return err
	}
	if err := assertMixResultsCorrect(request); err != nil {
		return err
	}
	if err := assertWellNotOverfilled(ctx, request); err != nil {
		return err
	}

	instructions, stockconcs, err := solution_setup(request, this.Properties)
	if err != nil {
		return err
	}

	if err := assertVolumesNonNegative(request); err != nil {
		return err
	}
	if err := assertTotalVolumesMatch(request); err != nil {
		return err
	}
	if err := assertMixResultsCorrect(request); err != nil {
		return err
	}

	request.LHInstructions = instructions
	request.Stockconcs = stockconcs

	// set up the mapping of the outputs
	// tried moving here to see if we can use results in fixVolumes
	request, err = this.Layout(ctx, request)

	if err != nil {
		return err
	}
	forceSanity(request)
	anotherSanityCheck(request)

	// assert: all instructions should now be assigned specific plate IDs, types and wells
	checkDestinationSanity(request)

	if request.Options.FixVolumes {
		// see if volumes can be corrected
		request, err = FixVolumes(request)

		if err != nil {
			return err
		}
		if request.Options.PrintInstructions {
			fmt.Println("")
			fmt.Println("Instructions Post Volume Fix")
			for _, insID := range request.OutputOrder {
				fmt.Println(request.LHInstructions[insID])
			}
		}
	}

	if err := assertVolumesNonNegative(request); err != nil {
		return err
	}
	if err := assertTotalVolumesMatch(request); err != nil {
		return err
	}
	if err := assertMixResultsCorrect(request); err != nil {
		return err
	}
	if err := assertWellNotOverfilled(ctx, request); err != nil {
		return err
	}

	orderedInstructions, err := request.GetOrderedLHInstructions()
	if err != nil {
		return err
	}

	//find what liquids are explicitely provided by the user
	solutionsFromPlates, err := request.GetSolutionsFromInputPlates()
	if err != nil {
		return err
	}

	// looks at liquids provided, calculates liquids required
	if inputSolutions, err := GetInputs(orderedInstructions, solutionsFromPlates, request.CarryVolume); err != nil {
		return err
	} else {
		request.InputSolutions = inputSolutions
		if request.Options.PrintInstructions {
			fmt.Println(inputSolutions)
		}
	}

	// define the input plates
	// should be merged with the above
	request, err = input_plate_setup(ctx, request)

	if err != nil {
		return err
	}

	// next we need to determine the liquid handler setup
	request, err = this.Setup(ctx, request)
	if err != nil {
		return err
	}

	// final insurance that plate names will be safe

	request = fixDuplicatePlateNames(request)

	// remove dummy mix-in-place instructions

	request = removeDummyInstructions(request)

	//set the well targets
	err = this.addWellTargets()
	if err != nil {
		return err
	}

	// now make instructions
	if rq, finalProps, err := this.GenerateInstructions(ctx, request); err != nil {
		return err
	} else {
		request = rq
		this.FinalProperties = finalProps
	}

	// revise the volumes - this makes sure the autoallocated volumes are correct
	if err := this.shrinkVolumes(request); err != nil {
		return err
	}

	// now make instructions with the updated volumes
	if rq, finalProps, err := this.GenerateInstructions(ctx, request); err != nil {
		return errors.WithMessage(err, "in second round of execution planning")
	} else {
		request = rq
		// duplicate this time so that final IDs are different
		this.FinalProperties = finalProps.Dup()
	}

	// Ensures tip boxes and wastes are correct for initial and final robot states
	this.AddInitialTipboxes()

	// counts tips used in this run -- reads instructions generated above so must happen
	// after execution planning
	request, err = this.countTipsUsed(request)
	if err != nil {
		return err
	}

	if err := this.makePlateIDMap(); err != nil {
		return err
	}

	// ensure the after state is correct
	this.fixPostIDs()
	if err := this.fixPostNames(request); err != nil {
		return err
	}

	err = assertNoTemporaryPlates(ctx, request)

	return err
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

// GenerateInstructions generate the low level liquidhandling instructions (LHRequest.Instructions)
// from the high level liquidhandling instructions (LHRequest.LHInstructions) with the initial
// robot state given by this.Properties
// returns a new request object containing the TerminalRobotInstructions and the final robot
// state after all instructions are executed
func (this *Liquidhandler) GenerateInstructions(ctx context.Context, request *LHRequest) (*LHRequest, *liquidhandling.LHProperties, error) {
	robot := this.Properties.DupKeepIDs()
	if rq, err := this.ExecutionPlanner(ctx, request, robot); err != nil {
		return nil, nil, err
	} else {
		return rq, robot, err
	}
}

func OutputSetup(robot *liquidhandling.LHProperties) {
	logger.Debug("DECK SETUP INFO")
	logger.Debug("Tipboxes: ")

	for k, v := range robot.Tipboxes {
		logger.Debug(fmt.Sprintf("%s %s: %s", k, robot.PlateIDLookup[k], v.Type))
	}

	logger.Debug("Plates:")

	for k, v := range robot.Plates {

		logger.Debug(fmt.Sprintf("%s %s: %s %s", k, robot.PlateIDLookup[k], v.PlateName, v.Type))

		//TODO Deprecate
		if strings.Contains(v.GetName(), "Input") {
			_, err := wtype.AutoExportPlateCSV(v.GetName()+".csv", v)
			if err != nil {
				logger.Debug(fmt.Sprintf("export plate csv (deprecated): %s", err.Error()))
			}
		}

		v.OutputLayout()
	}

	logger.Debug("Tipwastes: ")

	for k, v := range robot.Tipwastes {
		logger.Debug(fmt.Sprintf("%s %s: %s capacity %d", k, robot.PlateIDLookup[k], v.Type, v.Capacity))
	}

}

// fixPostIDs change the final component IDs such that all component IDs are changed by the liquidhandling
func (lh *Liquidhandler) fixPostIDs() {
	for _, p := range lh.FinalProperties.Plates {
		for _, w := range p.Wellcoords {
			w.WContents.ID = wtype.GetUUID()
		}
	}
}

// fixPostNames change the final component names to the value set by the last
// LHIMIX command which refers to that component
func (lh *Liquidhandler) fixPostNames(rq *LHRequest) error {
	// Instructions updating a well
	assignment := make(map[*wtype.LHWell]*wtype.LHInstruction)
	for _, inst := range rq.LHInstructions {
		// ignore non -mix instructions
		if inst.Type != wtype.LHIMIX {
			continue
		}

		tx := strings.Split(inst.Outputs[0].Loc, ":")
		newid, ok := lh.plateIDMap[tx[0]]
		if !ok {
			return wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprintf("No output plate mapped to %s", tx[0]))
		}

		ip, ok := lh.FinalProperties.PlateLookup[newid]
		if !ok {
			return wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprintf("No output plate %s", newid))
		}

		p, ok := ip.(*wtype.Plate)
		if !ok {
			return wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprintf("Got %s, should have *wtype.LHPlate", reflect.TypeOf(ip)))
		}

		well, ok := p.Wellcoords[tx[1]]
		if !ok {
			return wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprintf("No well %s on plate %s", tx[1], tx[0]))
		}

		oldInst := assignment[well]
		if oldInst == nil {
			assignment[well] = inst
		} else if prev, cur := oldInst.Outputs[0].Generation(), inst.Outputs[0].Generation(); prev < cur {
			assignment[well] = inst
		}
	}

	for well, inst := range assignment {
		well.WContents.CName = inst.Outputs[0].CName
	}

	return nil
}

func removeDummyInstructions(rq *LHRequest) *LHRequest {
	toRemove := make(map[string]bool, len(rq.LHInstructions))
	for _, ins := range rq.LHInstructions {
		if ins.IsDummy() {
			toRemove[ins.ID] = true
		}
	}

	if len(toRemove) == 0 {
		//no dummies
		return rq
	}

	oo := make([]string, 0, len(rq.OutputOrder)-len(toRemove))

	for _, ins := range rq.OutputOrder {
		if toRemove[ins] {
			continue
		} else {
			oo = append(oo, ins)
		}
	}

	if len(oo) != len(rq.OutputOrder)-len(toRemove) {
		panic(fmt.Sprintf("Dummy instruction prune failed: before %d dummies %d after %d", len(rq.OutputOrder), len(toRemove), len(oo)))
	}

	rq.OutputOrder = oo

	// prune instructionChain

	rq.InstructionChain.PruneOut(toRemove)

	return rq
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
		logger.Info("Not adding well target data for non-gilson adaptor")
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

func (req *LHRequest) MergedInputOutputPlates() map[string]*wtype.Plate {
	m := make(map[string]*wtype.Plate, len(req.InputPlates)+len(req.OutputPlates))
	addToMap(m, req.InputPlates)
	addToMap(m, req.OutputPlates)
	return m
}

func addToMap(m, a map[string]*wtype.Plate) {
	for k, v := range a {
		m[k] = v
	}
}

func fixDuplicatePlateNames(rq *LHRequest) *LHRequest {
	seen := make(map[string]int, 1)
	fixNames := func(sa []string, pm map[string]*wtype.Plate) {
		for _, id := range sa {
			p, foundPlate := pm[id]

			if !foundPlate {
				panic(fmt.Sprintf("Inconsistency in plate order / map for plate ID %s ", id))
			}

			n, ok := seen[p.PlateName]

			if ok {
				newName := fmt.Sprintf("%s_%d", p.PlateName, n)
				seen[p.PlateName] += 1
				p.PlateName = newName
			} else {
				seen[p.PlateName] = 1
			}
		}
	}

	fixNames(rq.InputPlateOrder, rq.InputPlates)
	fixNames(rq.OutputPlateOrder, rq.OutputPlates)

	return rq
}

// makePlateIDMap build the map between initial plate ID and final plate ID
func (this *Liquidhandler) makePlateIDMap() error {
	// map plates, tipboxes and tipwastes from initial to final assuming positions don't change
	this.plateIDMap = make(map[string]string, len(this.Properties.Positions))
	for position := range this.Properties.Positions {
		if initialID, finalID := this.Properties.PosLookup[position], this.FinalProperties.PosLookup[position]; initialID == "" || finalID == "" {
			if initialID != finalID {
				return wtype.LHErrorf(wtype.LH_ERR_DIRE, "layout changed for position %s: initially object with ID %q, then %q", position, initialID, finalID)
			}
		} else if initial, ok := this.Properties.PlateLookup[initialID]; !ok {
			return wtype.LHErrorf(wtype.LH_ERR_DIRE, "initial state inconsistent: object with id %q at %s not present in PlateLookup", initialID, position)
		} else if final, ok := this.FinalProperties.PlateLookup[finalID]; !ok {
			return wtype.LHErrorf(wtype.LH_ERR_DIRE, "final state inconsistent: object with id %q at %s not present in PlateLookup", finalID, position)
		} else if initialClass, finalClass := wtype.ClassOf(initial), wtype.ClassOf(final); initialClass != finalClass {
			return wtype.LHErrorf(wtype.LH_ERR_DIRE, "cannot map object of class %q in position %s to object of class %q", initialClass, position, finalClass)
		} else {
			this.plateIDMap[initialID] = finalID
		}
	}

	return nil
}

// AddInitialTipboxes adds full tipboxes to the initial Properties objects
// based on where tipboxes are found in FinalProperties
func (lh *Liquidhandler) AddInitialTipboxes() {
	for pos, final := range lh.FinalProperties.Tipboxes {
		initial := final.Dup()
		initial.Refresh()
		lh.Properties.AddTipBoxTo(pos, initial)
	}
}
