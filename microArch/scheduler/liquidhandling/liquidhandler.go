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
	"fmt"
	"time"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/driver"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/factory"
	"github.com/antha-lang/antha/microArch/logger"
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
	SetupAgent       func(*LHRequest, *liquidhandling.LHProperties) (*LHRequest, error)
	LayoutAgent      func(*LHRequest, *liquidhandling.LHProperties) (*LHRequest, error)
	ExecutionPlanner func(*LHRequest, *liquidhandling.LHProperties) (*LHRequest, error)
	PolicyManager    *LHPolicyManager
	plateIDMap       map[string]string // which plates are before / after versions
}

// initialize the liquid handling structure
func Init(properties *liquidhandling.LHProperties) *Liquidhandler {
	lh := Liquidhandler{}
	lh.SetupAgent = BasicSetupAgent
	lh.LayoutAgent = ImprovedLayoutAgent
	lh.ExecutionPlanner = ImprovedExecutionPlanner
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

// high-level function which requests planning and execution for an incoming set of
// solutions
func (this *Liquidhandler) MakeSolutions(request *LHRequest) error {
	// the minimal request which is possible defines what solutions are to be made
	if len(request.LHInstructions) == 0 {
		return wtype.LHError(wtype.LH_ERR_OTHER, "Nil plan requested: no Mix Instructions present")
	}

	//f := func() {
	err := this.Plan(request)

	if err != nil {
		return err
	}
	err = this.Execute(request)

	if err != nil {
		return err
	}

	// output some info on the final setup

	OutputSetup(this.Properties)

	//}

	// this is protective, should not be needed
	//err := this.Once.Do(f)

	return nil
}

// run the request via the driver
func (this *Liquidhandler) Execute(request *LHRequest) error {
	// set up the robot

	err := this.do_setup(request)

	if err != nil {
		return err
	}

	instructions := (*request).Instructions

	if instructions == nil {
		return wtype.LHError(wtype.LH_ERR_OTHER, "Cannot execute request: no instructions")
	}

	// some timing info for the log (only) for now

	timer := this.Properties.GetTimer()
	var d time.Duration

	for _, ins := range instructions {
		//logger.Debug(fmt.Sprintln(liquidhandling.InsToString(ins)))
		//fmt.Println(liquidhandling.InsToString(ins))
		ins.(liquidhandling.TerminalRobotInstruction).OutputTo(this.Properties.Driver)

		if timer != nil {
			d += timer.TimeFor(ins)
		}
	}

	logger.Debug(fmt.Sprintf("Total time estimate: %s", d.String()))
	request.TimeEstimate = d.Seconds()

	return nil
}

func (this *Liquidhandler) revise_volumes(rq *LHRequest) error {
	// XXX -- HARD CODE 8 here
	lastPlate := make([]string, 8)
	lastWell := make([]string, 8)

	vols := make(map[string]map[string]wunit.Volume)

	for _, ins := range rq.Instructions {
		if ins.InstructionType() == liquidhandling.MOV {
			lastPlate = make([]string, 8)
			lastPos := ins.GetParameter("POSTO").([]string)

			for i, p := range lastPos {
				lastPlate[i] = this.FinalProperties.PosLookup[p]
			}

			lastWell = ins.GetParameter("WELLTO").([]string)
		} else if ins.InstructionType() == liquidhandling.ASP {
			for i, _ := range lastPlate {
				if i >= len(lastWell) {
					break
				}
				lp := lastPlate[i]
				lw := lastWell[i]

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

				insvols := ins.GetParameter("VOLUME").([]wunit.Volume)
				v.Add(insvols[i])
				v.Add(rq.CarryVolume)
			}
		}
	}

	// now go through and set the plates up appropriately

	for plateID, wellmap := range vols {
		plate, ok := this.FinalProperties.Plates[this.FinalProperties.PlateIDLookup[plateID]]
		plate2, _ := this.Properties.Plates[this.FinalProperties.PlateIDLookup[plateID]]

		if !ok {
			err := wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprint("NO SUCH PLATE: ", plateID))
			return err
		}

		for crd, vol := range wellmap {
			well := plate.Wellcoords[crd]
			well2 := plate2.Wellcoords[crd]
			if well.IsAutoallocated() {
				vol.Add(well.ResidualVolume())
				well2.WContents.SetVolume(vol)
				well.WContents.SetVolume(well.ResidualVolume())
				well.DeclareNotTemporary()
				well2.DeclareNotTemporary()
			}
		}
	}

	// finally get rid of any temporary stuff

	this.Properties.RemoveTemporaryComponents()
	this.FinalProperties.RemoveTemporaryComponents()

	// now ensure the mapping, excluding any temp plates added above, is recorded

	for plateID, _ := range vols {
		plate, ok := this.FinalProperties.Plates[this.FinalProperties.PlateIDLookup[plateID]]

		if !ok {
			// plate deleted
			continue
		}

		pos := this.FinalProperties.PlateIDLookup[plateID]

		plate2 := this.Properties.Plates[pos]

		this.plateIDMap[plate2.ID] = plate.ID
	}

	// all done

	return nil
}

func (this *Liquidhandler) do_setup(rq *LHRequest) error {

	stat := this.Properties.Driver.RemoveAllPlates()

	if stat.Errorcode == driver.ERR {
		return wtype.LHError(wtype.LH_ERR_DRIV, stat.Msg)
	}

	for position, plateid := range this.Properties.PosLookup {
		if plateid == "" {
			continue
		}
		plate := this.Properties.PlateLookup[plateid]
		name := plate.(wtype.Named).GetName()
		stat = this.Properties.Driver.AddPlateTo(position, plate, name)

		if stat.Errorcode == driver.ERR {
			return wtype.LHError(wtype.LH_ERR_DRIV, stat.Msg)
		}
	}

	stat = this.Properties.Driver.(liquidhandling.ExtendedLiquidhandlingDriver).UpdateMetaData(this.Properties)
	if stat.Errorcode == driver.ERR {
		return wtype.LHError(wtype.LH_ERR_DRIV, stat.Msg)
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

func (this *Liquidhandler) Plan(request *LHRequest) error {
	// figure out the output order

	err := set_output_order(request)

	if err != nil {
		return err
	}

	request.InstructionChain.Print()

	// convert requests to volumes and determine required stock concentrations
	instructions, stockconcs, err := solution_setup(request, this.Properties)

	if err != nil {
		return err
	}

	request.LHInstructions = instructions
	request.Stockconcs = stockconcs

	// looks at components, determines what inputs are required
	request, err = this.GetInputs(request)

	if err != nil {
		return err
	}
	// define the input plates
	// should be merged with the above

	request, err = input_plate_setup(request)

	if err != nil {
		return err
	}
	// set up the mapping of the outputs
	request, err = this.Layout(request)

	if err != nil {
		return err
	}

	// next we need to determine the liquid handler setup
	request, err = this.Setup(request)
	if err != nil {
		return err
	}

	// now make instructions
	request, err = this.ExecutionPlan(request)

	if err != nil {
		return err
	}
	// fix the deck setup
	request, err = this.Tip_box_setup(request)
	if err != nil {
		return err
	}

	// revise the volumes
	err = this.revise_volumes(request)

	if err != nil {
		return err
	}

	return nil
}

// sort out inputs
func (this *Liquidhandler) GetInputs(request *LHRequest) (*LHRequest, error) {
	instructions := (*request).LHInstructions
	inputs := make(map[string][]*wtype.LHComponent, 3)
	order := make(map[string]map[string]int, 3)
	vmap := make(map[string]wunit.Volume)

	allinputs := make([]string, 0, 10)

	for _, instruction := range instructions {
		components := instruction.Components

		for _, component := range components {
			// ignore anything which is made in another mix

			if component.HasAnyParent() {
				continue
			}

			cmps, ok := inputs[component.CName]
			if !ok {
				cmps = make([]*wtype.LHComponent, 0, 3)
				allinputs = append(allinputs, component.CName)
			}

			cmps = append(cmps, component)
			inputs[component.CName] = cmps

			// similarly add the volumes up

			vol := vmap[component.CName]

			if vol.IsNil() {
				vol = wunit.NewVolume(0.0, "ul")
			}

			v2a := wunit.NewVolume(component.Vol, component.Vunit)

			// we have to add the carry volume here
			// this is roughly per transfer so should be OK
			v2a.Add(request.CarryVolume)
			vol.Add(v2a)

			vmap[component.CName] = vol

			for j := 0; j < len(components); j++ {
				// again exempt those parented components
				if components[j].HasAnyParent() {
					continue
				}
				if component.Order < components[j].Order {
					m, ok := order[component.CName]
					if !ok {
						m = make(map[string]int, len(components))
						order[component.CName] = m
					}

					m[components[j].CName] += 1
				} else {
					m, ok := order[components[j].CName]
					if !ok {
						m = make(map[string]int, len(components))
						order[components[j].CName] = m
					}
					m[component.CName] += 1
				}
			}

		}
	}

	// define component ordering

	component_order, err := DefineOrderOrFail(order)

	if err != nil {
		return request, err
	}

	(*request).Input_order = component_order

	// work out how much we have and how much we need
	// need to consider what to do with IDs

	var requestinputs map[string][]*wtype.LHComponent
	requestinputs = request.Input_solutions

	if len(requestinputs) == 0 {
		requestinputs = make(map[string][]*wtype.LHComponent, 5)
	}

	vmap2 := make(map[string]wunit.Volume, len(vmap))
	vmap3 := make(map[string]wunit.Volume, len(vmap))

	for _, k := range allinputs {
		// vola: how much comes in
		ar := requestinputs[k]
		vola := wunit.NewVolume(0.00, "ul")
		for _, cmp := range ar {
			vold := wunit.NewVolume(cmp.Vol, cmp.Vunit)
			vola.Add(vold)
		}
		// volb: how much we asked for
		volb := vmap[k].Dup()
		volb.Subtract(vola)
		vmap2[k] = vola

		if volb.GreaterThanFloat(0.0001) {
			vmap3[k] = volb
		}
		//	volc := vmap[k]
		//logger.Debug(fmt.Sprint("COMPONENT ", k, " HAVE : ", vola.ToString(), " WANT: ", volc.ToString(), " DIFF: ", volb.ToString()))
	}

	(*request).Input_vols_required = vmap
	(*request).Input_vols_supplied = vmap2
	(*request).Input_vols_wanting = vmap3

	// add any new inputs

	for k, v := range inputs {
		if requestinputs[k] == nil {
			requestinputs[k] = v
		}
	}

	(*request).Input_solutions = requestinputs

	return request, nil
}

func DefineOrderOrFail(mapin map[string]map[string]int) ([]string, error) {
	cmps := make([]string, 0, 1)

	for name, _ := range mapin {
		cmps = append(cmps, name)
	}

	ord := make([][]string, len(cmps))

	mx := 0
	for i := 0; i < len(cmps); i++ {
		cnt := 0
		for j := 0; j < len(cmps); j++ {
			if i == j {
				continue
			}

			// PREVIOUSLY
			// only one side can be > 0
			// NOW we don't care

			c1 := mapin[cmps[i]][cmps[j]]
			//c2 := mapin[cmps[j]][cmps[i]]

			/*
				if c1 > 0 && c2 > 0 {
					log.Fatalf(fmt.Sprint("CANNOT DEAL WITH INCONSISTENT COMPONENT ORDERING ", cmps[i], " ", cmps[j], " ", c1, " ", c2))
				}

			*/
			// if c1 > 0 we add to the count

			if c1 > 0 {
				cnt += 1
			}
		}

		a := ord[cnt]

		if a == nil {
			a = make([]string, 0, 3)
		}

		a = append(a, cmps[i])
		if cnt > mx {
			mx = cnt
		}
		ord[cnt] = a
	}

	ret := make([]string, 0, len(cmps))

	// take in reverse order
	if len(cmps) > 0 {
		for j := mx; j >= 0; j-- {
			a := ord[j]
			if a == nil {
				continue
			}

			for _, name := range a {
				ret = append(ret, name)
			}
		}
	}

	return ret, nil
}

// define which labware to use
func (this *Liquidhandler) GetPlates(plates map[string]*wtype.LHPlate, major_layouts map[int][]string, ptype *wtype.LHPlate) map[string]*wtype.LHPlate {
	if plates == nil {
		plates = make(map[string]*wtype.LHPlate, len(major_layouts))

		// assign new plates
		for i := 0; i < len(major_layouts); i++ {
			//newplate := wtype.New_Plate(ptype)
			newplate := factory.GetPlateByType(ptype.Type)
			plates[newplate.ID] = newplate
		}
	}

	// we should know how many plates we need
	for k, plate := range plates {
		if plate.Inst == "" {
			//stockrequest := execution.GetContext().StockMgr.RequestStock(makePlateStockRequest(plate))
			//plate.Inst = stockrequest["inst"].(string)
		}

		plates[k] = plate
	}

	return plates
}

// generate setup for the robot
func (this *Liquidhandler) Setup(request *LHRequest) (*LHRequest, error) {
	// assign the plates to positions
	// this needs to be parameterizable
	return this.SetupAgent(request, this.Properties)
}

// generate the output layout
func (this *Liquidhandler) Layout(request *LHRequest) (*LHRequest, error) {
	// assign the results to destinations
	// again needs to be parameterized

	return this.LayoutAgent(request, this.Properties)
}

// make the instructions for executing this request
func (this *Liquidhandler) ExecutionPlan(request *LHRequest) (*LHRequest, error) {
	//rbtcpy := this.Properties.Dup()
	this.FinalProperties = this.Properties.Dup()

	rq, err := this.ExecutionPlanner(request, this.FinalProperties)

	// ensure any relevant state changes are noted

	/*
		for _, v := range rbtcpy.Plates {
			fmt.Println("PLATE NAME: ", v.Name())
			for _, w := range v.Wellcoords {
				if !w.Empty() {
					fmt.Println("WELL: ", w.Crds, " HAS: ", w.WContents.Volume().ToString(), " of ", w.WContents.CName)
				}
			}
		}
	*/

	return rq, err
}

func OutputSetup(robot *liquidhandling.LHProperties) {
	logger.Debug("DECK SETUP INFO")
	logger.Debug("Tipboxes: ")

	for k, v := range robot.Tipboxes {
		logger.Debug(fmt.Sprintf("%s %s: %s", k, robot.PlateIDLookup[k], v.Type))
	}

	logger.Debug("Plates:")

	for k, v := range robot.Plates {
		logger.Debug(fmt.Sprintf("%s %s: %s", k, robot.PlateIDLookup[k], v.PlateName))
	}

	logger.Debug("Tipwastes: ")

	for k, v := range robot.Tipwastes {
		logger.Debug(fmt.Sprintf("%s %s: %s capacity %d", k, robot.PlateIDLookup[k], v.Type, v.Capacity))
	}

}
