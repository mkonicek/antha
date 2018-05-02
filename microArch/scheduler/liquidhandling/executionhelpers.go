// anthalib//liquidhandling/executionplanner.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
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
	"sort"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/graph"
	driver "github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

const (
	COLWISE = iota
	ROWWISE
	RANDOM
)

type ByGeneration []*wtype.LHInstruction

func (bg ByGeneration) Len() int      { return len(bg) }
func (bg ByGeneration) Swap(i, j int) { bg[i], bg[j] = bg[j], bg[i] }
func (bg ByGeneration) Less(i, j int) bool {
	if bg[i].Generation() == bg[j].Generation() {

		// compare the plate names (which must exist now)
		//	 -- oops, I think this has ben violated by moving the sort
		// 	 TODO check and fix

		c := strings.Compare(bg[i].PlateName, bg[j].PlateName)

		if c != 0 {
			return c < 0
		}

		// finally go down columns (nb need to add option)

		return wtype.CompareStringWellCoordsCol(bg[i].Welladdress, bg[j].Welladdress) < 0
	}

	return bg[i].Generation() < bg[j].Generation()
}

// Optimally - order by component.
type ByGenerationOpt []*wtype.LHInstruction

func (bg ByGenerationOpt) Len() int      { return len(bg) }
func (bg ByGenerationOpt) Swap(i, j int) { bg[i], bg[j] = bg[j], bg[i] }
func (bg ByGenerationOpt) Less(i, j int) bool {
	if bg[i].Generation() == bg[j].Generation() {

		// compare the names of the resultant components
		c := strings.Compare(bg[i].Results[0].CName, bg[j].Results[0].CName)

		if c != 0 {
			return c < 0
		}

		// if two components names are equal, then compare the plates
		c = strings.Compare(bg[i].PlateName, bg[j].PlateName)

		if c != 0 {
			return c < 0
		}

		// finally go down columns (nb need to add option)

		return wtype.CompareStringWellCoordsCol(bg[i].Welladdress, bg[j].Welladdress) < 0
	}

	return bg[i].Generation() < bg[j].Generation()
}

type ByColumn []*wtype.LHInstruction

func (bg ByColumn) Len() int      { return len(bg) }
func (bg ByColumn) Swap(i, j int) { bg[i], bg[j] = bg[j], bg[i] }
func (bg ByColumn) Less(i, j int) bool {
	// compare any messages present (only really applies to prompts)
	c := strings.Compare(bg[i].Message, bg[j].Message)

	if c != 0 {
		return c < 0
	}
	// compare the plate names (which must exist now)
	//	 -- oops, I think this has ben violated by moving the sort
	// 	 TODO check and fix

	c = strings.Compare(bg[i].PlateName, bg[j].PlateName)

	if c != 0 {
		return c < 0
	}

	// Go Down Columns

	return wtype.CompareStringWellCoordsCol(bg[i].Welladdress, bg[j].Welladdress) < 0
}

// Optimally - order by component.
type ByResultComponent []*wtype.LHInstruction

func (bg ByResultComponent) Len() int      { return len(bg) }
func (bg ByResultComponent) Swap(i, j int) { bg[i], bg[j] = bg[j], bg[i] }
func (bg ByResultComponent) Less(i, j int) bool {
	// compare any messages present

	c := strings.Compare(bg[i].Message, bg[j].Message)

	if c != 0 {
		return c < 0
	}

	// compare the names of the resultant components
	c = strings.Compare(bg[i].Results[0].CName, bg[j].Results[0].CName)

	if c != 0 {
		return c < 0
	}

	// if two components names are equal, then compare the plates
	c = strings.Compare(bg[i].PlateName, bg[j].PlateName)

	if c != 0 {
		return c < 0
	}

	// finally go down columns (nb need to add option)

	return wtype.CompareStringWellCoordsCol(bg[i].Welladdress, bg[j].Welladdress) < 0
}

func convertToInstructionChain(sortedNodes []graph.Node, tg graph.Graph, sort bool, inputs map[string][]*wtype.LHComponent) *IChain {
	ic := NewIChain(nil)

	// the nodes are now ordered according to dependency relations
	// *IN REVERSE ORDER*

	// this routine defines equivalence classes of nodes

	for _, n := range sortedNodes {
		addToIChain(ic, n, tg)
	}

	// we need to ensure that splits and mixes are kept separate by fissioning nodes

	ic.SplitMixedNodes()

	// this routine ensures that instructions can be executed in parallel

	if ic == nil {
		ic = simplifyIChain(ic, inputs)
	}

	sortOutputs(ic, sort)

	return ic
}

func sortOutputs(ic *IChain, byComponent bool) {
	// recursively progress through the chain, sorting values as we go

	if ic == nil {
		return
	}

	if byComponent {
		sort.Sort(ByResultComponent(ic.Values))
	} else {
		sort.Sort(ByColumn(ic.Values))
	}

	sortOutputs(ic.Child, byComponent)
}

func addToIChain(ic *IChain, n graph.Node, tg graph.Graph) {
	deps := make(map[graph.Node]bool)

	for i := 0; i < tg.NumOuts(n); i++ {
		deps[tg.Out(n, i)] = true
	}

	cur := findNode(ic, n, tg, deps)
	cur.Values = append(cur.Values, n.(*wtype.LHInstruction))
}

func findNode(ic *IChain, n graph.Node, tg graph.Graph, deps map[graph.Node]bool) *IChain {
	if thisNode(ic, n, tg, deps) {
		return ic
	} else if ic.Child != nil {
		return findNode(ic.Child, n, tg, deps)
	} else {
		newNode := NewIChain(ic)
		ic.Child = newNode
		return newNode
	}
}

func thisNode(ic *IChain, n graph.Node, tg graph.Graph, deps map[graph.Node]bool) bool {
	// if this looks weird it's because "output" below really means "input"
	// since we have reversed dependency order

	// delete any deps satisfied by this node

	if ic.Parent != nil {
		for _, v := range ic.Parent.Values {
			delete(deps, graph.Node(v))
		}
	}

	// have we seen all of the outputs? If so, stop here

	if len(deps) == 0 {
		return true
	}

	// if not

	return false
}

func getInstructionSet(rq *LHRequest) []*wtype.LHInstruction {
	ret := make([]*wtype.LHInstruction, 0, len(rq.LHInstructions))
	for _, v := range rq.LHInstructions {
		ret = append(ret, v)
	}

	return ret
}

// is n1 an ancestor of n2?
func ancestor(n1, n2 graph.Node, topolGraph graph.Graph) bool {
	if n1 == n2 {
		return true
	}

	for i := 0; i < topolGraph.NumOuts(n2); i++ {
		if ancestor(n1, topolGraph.Out(n2, i), topolGraph) {
			return true
		}
	}

	return false
}

// track back and see if one depends on the other
func related(i1, i2 *wtype.LHInstruction, topolGraph graph.Graph) bool {
	if ancestor(graph.Node(i1), graph.Node(i2), topolGraph) || ancestor(graph.Node(i2), graph.Node(i1), topolGraph) {
		return true
	}

	return false
}

func canAggHere(ar []*wtype.LHInstruction, ins *wtype.LHInstruction, topolGraph graph.Graph) bool {
	for _, i2 := range ar {
		if related(ins, i2, topolGraph) {
			return false
		}
	}

	return true
}

// we can only append if we don't create cycles
// this function makes sure this is OK
func appendSensitively(iar [][]*wtype.LHInstruction, ins *wtype.LHInstruction, topolGraph graph.Graph) [][]*wtype.LHInstruction {
	done := false
	for i := 0; i < len(iar); i++ {
		ar := iar[i]
		// just add to the first available one
		if canAggHere(ar, ins, topolGraph) {
			ar = append(ar, ins)
			iar[i] = ar
			done = true
			break
		}
	}

	if !done {
		ar := make([]*wtype.LHInstruction, 0, 1)
		ar = append(ar, ins)
		iar = append(iar, ar)
	}

	return iar
}

func aggregatePromptsWithSameMessage(inss []*wtype.LHInstruction, topolGraph graph.Graph) []graph.Node {
	// merge dependencies of any prompts which have a message in common
	prMessage := make(map[string][][]*wtype.LHInstruction, len(inss))
	insOut := make([]graph.Node, 0, len(inss))

	for _, ins := range inss {
		if ins.Type == wtype.LHIPRM {
			iar, ok := prMessage[ins.Message]

			if !ok {
				iar = make([][]*wtype.LHInstruction, 0, len(inss)/2)
			}

			//iar = append(iar, ins)

			iar = appendSensitively(iar, ins, topolGraph)

			prMessage[ins.Message] = iar
		} else {
			insOut = append(insOut, graph.Node(ins))
		}
	}

	// aggregate instructions
	// TODO --> user control of scope of this aggregation
	//          i.e. break every plate, some other subset

	for msg, iar := range prMessage {
		// single message may appear multiply in the chain
		for _, ar := range iar {
			ins := wtype.NewLHPromptInstruction()
			ins.Message = msg
			ins.AddResult(wtype.NewLHComponent())
			for _, ins2 := range ar {
				for _, cmp := range ins2.Components {
					ins.Components = append(ins.Components, cmp)
					ins.PassThrough[cmp.ID] = ins2.Results[0]
				}
			}
			insOut = append(insOut, graph.Node(ins))
		}

	}

	return insOut
}

func set_output_order(rq *LHRequest) error {
	// guarantee all nodes are dependency-ordered
	// in order to aggregate without introducing cycles

	unsorted := getInstructionSet(rq)

	tg := MakeTGraph(unsorted)

	sorted, err := graph.TopoSort(graph.TopoSortOpt{Graph: tg})

	if err != nil {
		return err
	}

	sortedAsIns := make([]*wtype.LHInstruction, len(sorted))

	for i := 0; i < len(sorted); i++ {
		sortedAsIns[i] = sorted[i].(*wtype.LHInstruction)
	}

	sorted = aggregatePromptsWithSameMessage(sortedAsIns, tg)

	// aggregate sorted again
	sortedAsIns = make([]*wtype.LHInstruction, len(sorted))
	for i, nIns := range sorted {
		ins := nIns.(*wtype.LHInstruction)
		sortedAsIns[i] = ins
	}

	// update request to be consistent with new instructions
	rq = updateRequestWithNewInstructions(rq, sortedAsIns)

	// sort again post aggregation
	tg = MakeTGraph(sortedAsIns)

	sorted, err = graph.TopoSort(graph.TopoSortOpt{Graph: tg})

	if err != nil {
		return err
	}

	// make into equivalence classes and sort according to defined order
	it := convertToInstructionChain(sorted, tg, rq.Options.OutputSort, rq.Input_solutions)

	// populate the request
	rq.InstructionChain = it
	rq.Output_order = it.Flatten()

	return nil
}

func updateRequestWithNewInstructions(rq *LHRequest, sorted []*wtype.LHInstruction) *LHRequest {
	// make sure the request contains the new instructions if aggregation has occurred here
	for _, ins := range sorted {
		_, ok := rq.LHInstructions[ins.ID]
		if !ok {
			rq.LHInstructions[ins.ID] = ins
		}
	}
	return rq
}

type ByOrdinal [][]int

func (bo ByOrdinal) Len() int      { return len(bo) }
func (bo ByOrdinal) Swap(i, j int) { bo[i], bo[j] = bo[j], bo[i] }
func (bo ByOrdinal) Less(i, j int) bool {
	// just compare the first one

	return bo[i][0] < bo[j][0]
}

func merge_instructions(insIn []driver.RobotInstruction, aggregates [][]int) []driver.RobotInstruction {
	ret := make([]driver.RobotInstruction, 0, len(insIn))

	for _, ar := range aggregates {
		if len(ar) == 1 {
			// just push it in and move on
			ret = append(ret, insIn[ar[0]])
			continue
		}

		// otherwise more than one here

		newtfr, ok := insIn[ar[0]].(*driver.TransferInstruction)

		if ok {
			for k := 1; k < len(ar); k++ {
				newtfr.MergeWith(insIn[ar[k]].(*driver.TransferInstruction))
			}

			ret = append(ret, newtfr)
		} else {
			// must be a message
			ins1 := insIn[ar[0]]
			ret = append(ret, ins1)

			// put in any distinct instructions

			for i := 1; i < len(ar); i++ {
				if insIn[ar[i]].(*driver.MessageInstruction).Message != ins1.(*driver.MessageInstruction).Message {
					ret = append(ret, insIn[ar[i]])
					ins1 = insIn[ar[i]]
				}
			}

		}
	}

	return ret
}

// TODO -- refactor this to pass robot through
func ConvertInstruction(insIn *wtype.LHInstruction, robot *driver.LHProperties, carryvol wunit.Volume, legacyVolume bool) (insOut *driver.TransferInstruction, err error) {
	cmps := insIn.Components

	lenToMake := len(insIn.Components)

	if insIn.IsMixInPlace() {
		lenToMake = lenToMake - 1
		cmps = cmps[1:]
	}

	wh := make([]string, 0, lenToMake)       // component types
	va := make([]wunit.Volume, 0, lenToMake) // volumes

	//fromPlateIDs, fromWellss, volss, err := robot.GetComponents(driver.GetComponentsOptions{Cmps: cmps, Carryvol: carryvol, Ori: wtype.LHVChannel, Multi: 1, Independent: true, LegacyVolume: legacyVolume})
	getComponentsReply, err := robot.GetComponents(driver.GetComponentsOptions{Cmps: cmps, Carryvol: carryvol, Ori: wtype.LHVChannel, Multi: 1, Independent: true, LegacyVolume: legacyVolume})

	if err != nil {
		return nil, err
	}

	tfrs := getComponentsReply.Transfers

	pf := make([]string, 0, lenToMake)
	wf := make([]string, 0, lenToMake)
	pfwx := make([]int, 0, lenToMake)
	pfwy := make([]int, 0, lenToMake)
	vf := make([]wunit.Volume, 0, lenToMake)
	ptt := make([]string, 0, lenToMake)

	// six parameters applying to the destination

	pt := make([]string, 0, lenToMake)       // dest plate positions
	wt := make([]string, 0, lenToMake)       // dest wells
	ptwx := make([]int, 0, lenToMake)        // dimensions of plate pipetting to (X)
	ptwy := make([]int, 0, lenToMake)        // dimensions of plate pipetting to (Y)
	vt := make([]wunit.Volume, 0, lenToMake) // volume in well to
	ptf := make([]string, 0, lenToMake)      // plate types
	cnames := make([]string, 0, lenToMake)   // actual Component names
	policies := make([]wtype.LHPolicy, 0, lenToMake)

	for i, v := range cmps {
		for xx := range tfrs[i].PlateIDs { //fromPlateIDs[i] {
			// get dem big ole plates out
			// TODO -- pass them in instead of all this nonsense

			var flhp, tlhp *wtype.LHPlate

			flhif := robot.PlateLookup[tfrs[i].PlateIDs[xx]] //[fromPlateIDs[i][xx]]

			if flhif != nil {
				flhp = flhif.(*wtype.LHPlate)
			} else {
				s := fmt.Sprint("NO SRC PLATE FOUND : ", i, " ", xx, " ", tfrs[i].PlateIDs[xx]) //fromPlateIDs[i][xx])
				err := wtype.LHError(wtype.LH_ERR_DIRE, s)

				return nil, err
			}

			tlhif := robot.PlateLookup[insIn.PlateID]

			if tlhif != nil {
				tlhp = tlhif.(*wtype.LHPlate)
			} else {
				s := fmt.Sprint("NO DST PLATE FOUND : ", i, " ", xx, " ", insIn.PlateID)
				err := wtype.LHError(wtype.LH_ERR_DIRE, s)

				return nil, err
			}

			wlt, ok := tlhp.WellAtString(insIn.Welladdress)

			if !ok {
				return nil, wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprintf("Well %s not found on dest plate %s", insIn.Welladdress, insIn.PlateID))
			}

			//v2 := wunit.NewVolume(v.Vol, v.Vunit)
			v2 := tfrs[i].Vols[xx] // volss[i][xx]
			vt = append(vt, wlt.CurrentVolume())
			wh = append(wh, v.TypeName())
			va = append(va, v2)
			pt = append(pt, robot.PlateIDLookup[insIn.PlateID])
			wt = append(wt, insIn.Welladdress)
			ptwx = append(ptwx, tlhp.WellsX())
			ptwy = append(ptwy, tlhp.WellsY())
			ptt = append(ptt, tlhp.Type)

			wlf, ok := flhp.WellAtString(tfrs[i].WellCoords[xx]) //fromWellss[i][xx])

			if !ok {
				//logger.Fatal(fmt.Sprint("Well ", fromWells[ix], " not found on source plate ", fromPlateID[ix]))
				//err = wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprint("Well ", fromWellss[i][xx], " not found on source plate ", fromPlateIDs[i][xx]))
				err = wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprint("Well ", tfrs[i].WellCoords[xx], " not found on source plate ", tfrs[i].PlateIDs[xx]))
				return nil, err
			}

			vf = append(vf, wlf.CurrentVolume())
			vrm := v2.Dup()
			vrm.Add(carryvol)
			cnames = append(cnames, wlf.WContents.CName)
			policies = append(policies, wlf.WContents.Policy)
			if _, err := wlf.RemoveVolume(vrm); err != nil {
				return nil, err
			}

			pf = append(pf, robot.PlateIDLookup[tfrs[i].PlateIDs[xx]])
			wf = append(wf, tfrs[i].WellCoords[xx])
			pfwx = append(pfwx, flhp.WellsX())
			pfwy = append(pfwy, flhp.WellsY())
			ptf = append(ptf, flhp.Type)

			if v.Loc == "" {
				v.Loc = tfrs[i].PlateIDs[xx] + ":" + tfrs[i].WellCoords[xx]
			}
			// add component to destination

			// ensure we keep results straight
			vd := v.Dup()
			// volumes need to come from volss
			vd.Vol = v2.ConvertToString(vd.Vunit)
			vd.ID = wlf.WContents.ID
			vd.ParentID = wlf.WContents.ParentID
			err := wlt.AddComponent(vd)
			if err != nil {
				return nil, wtype.LHError(wtype.LH_ERR_VOL, fmt.Sprintf("Scheduler couldn't add volume to well : %s", err.Error()))
			}

			// TODO -- danger here, is result definitely set?
			wlt.WContents.ID = insIn.Results[0].ID
			wlf.WContents.AddDaughterComponent(wlt.WContents)

			//fmt.Println("HERE GOES: ", i, wh[i], vf[i].ToString(), vt[i].ToString(), va[i].ToString(), pt[i], wt[i], pf[i], wf[i], pfwx[i], pfwy[i], ptwx[i], ptwy[i])

		}
	}

	// what, pltfrom, pltto, wellfrom, wellto, fplatetype, tplatetype []string, volume, fvolume, tvolume []wunit.Volume, FPlateWX, FPlateWY, TPlateWX, TPlateWY []int
	ti := driver.NewTransferInstruction(wh, pf, pt, wf, wt, ptf, ptt, va, vf, vt, pfwx, pfwy, ptwx, ptwy, cnames, policies)

	return ti, nil
}
