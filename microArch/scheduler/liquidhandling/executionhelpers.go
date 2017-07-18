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
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/graph"
	driver "github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"sort"
	"strings"
)

const (
	COLWISE = iota
	ROWWISE
	RANDOM
)

func roundup(f float64) float64 {
	return float64(int(f) + 1)
}

func get_aggregate_component(sol *wtype.LHSolution, name string) *wtype.LHComponent {
	components := sol.Components

	ret := wtype.NewLHComponent()

	ret.CName = name

	vol := 0.0
	found := false

	for _, component := range components {
		nm := component.CName

		if nm == name {
			ret.Type = component.Type
			vol += component.Vol
			ret.Vunit = component.Vunit
			ret.Order = component.Order
			found = true
		}
	}
	if !found {
		return nil
	}
	ret.Vol = vol
	return ret
}

func get_assignment(assignments []string, plates *map[string]*wtype.LHPlate, vol wunit.Volume) (string, wunit.Volume, bool) {
	assignment := ""
	ok := false
	prevol := wunit.NewVolume(0.0, "ul")

	for _, assignment = range assignments {
		asstx := strings.Split(assignment, ":")
		plate := (*plates)[asstx[0]]

		crds := asstx[1] + ":" + asstx[2]
		wellidlkp := plate.Wellcoords
		well := wellidlkp[crds]

		currvol := well.CurrVolume()
		currvol.Subtract(well.ResidualVolume())
		if currvol.GreaterThan(vol) || currvol.EqualTo(vol) {
			prevol = well.CurrVolume()
			well.Remove(vol)
			plate.HWells[well.ID] = well
			(*plates)[asstx[0]] = plate
			ok = true
			break
		}
	}

	return assignment, prevol, ok
}

func copyplates(plts map[string]*wtype.LHPlate) map[string]*wtype.LHPlate {
	ret := make(map[string]*wtype.LHPlate, len(plts))

	for k, v := range plts {
		ret[k] = v.Dup()
	}

	return ret
}

func insSliceFromMap(m map[string]*wtype.LHInstruction) []*wtype.LHInstruction {
	ret := make([]*wtype.LHInstruction, 0, len(m))

	for _, v := range m {
		ret = append(ret, v)
	}

	return ret
}

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
		c := strings.Compare(bg[i].Result.CName, bg[j].Result.CName)

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
	// compare the plate names (which must exist now)
	//	 -- oops, I think this has ben violated by moving the sort
	// 	 TODO check and fix

	c := strings.Compare(bg[i].PlateName, bg[j].PlateName)

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
	// compare the names of the resultant components
	c := strings.Compare(bg[i].Result.CName, bg[j].Result.CName)

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

func aggregateAppropriateInstructions(inss []*wtype.LHInstruction) []*wtype.LHInstruction {
	agg := make([]map[string]*wtype.LHInstruction, len(wtype.InsNames))
	for i := 0; i < len(wtype.InsNames); i++ {
		agg[i] = make(map[string]*wtype.LHInstruction, 10)
	}

	for _, ins := range inss {
		// just prompts
		if ins.Type == wtype.LHIPRM {
			cur := agg[ins.Type][ins.Message]
			if cur == nil || cur.Generation() < ins.Generation() {
				agg[ins.Type][ins.Message] = ins
			}
		}
	}

	// now filter
	insout := make([]*wtype.LHInstruction, 0, len(inss))
	for _, ins := range inss {
		if ins.Type == wtype.LHIPRM {
			if agg[ins.Type][ins.Message].ID != ins.ID {
				continue
			}
		}
		insout = append(insout, ins)
	}

	return insout
}

func convertToInstructionChain(sortedNodes []graph.Node, tg graph.Graph, sort bool) *IChain {
	ic := NewIChain(nil)

	// the nodes are now ordered according to dependency relations
	// this routine defines equivalence classes of nodes

	for _, n := range sortedNodes {
		addToIChain(ic, n, tg)
	}

	// finally we sort outputs
	sortOutputs(ic, sort)

	return ic
}

func sortOutputs(ic *IChain, byComponent bool) {
	// recursively progress through the chain, sorting values as we go

	if ic == nil {
		return
	}

	if byComponent {
		sort.Sort(ByComponent(ic.Values))
	} else {
		sort.Sort(ByColumn(ic.Values))
	}

	sortOutputs(ic.Child, byComponent)
}

func addToIChain(ic *IChain, n graph.Node, tg graph.Graph) {
	cur := findNode(ic, n, tg)
	cur.Values = append(cur.Values, n.(*wtype.LHInstruction))
}

func findNode(ic *IChain, n graph.Node, tg graph.Graph) *IChain {
	if thisNode(ic, n, tg) {
		return ic
	} else if ic.Child != nil {
		return findNode(ic.Child, n, tg)
	} else {
		return NewIChain(ic)
	}
}

func thisNode(ic *IChain, n graph.Node, tg graph.Graph) bool {
	inEdges := nodeIns(n, tg)

	// if this node has no in edges we return immediately
	if len(inEdges) == 0 {
		return true
	}

	// look one up the chain to see if this instruction is in the output
	outEdges := chainOuts(ic, tg)

	if len(outEdges == 0) {
		// this node has ins but the parent in the chain has no outs, so not this one
		return false
	}

	// compare the sets

	return len(meetInsSets(inEdges, outEdges)) != 0
}

func nodeIns(n graph.Node, tg graph.Graph) []*wtype.LHInstruction {
	ins := n.(*wtype.LHInstruction)

	r := make([]*wtype.LHInstruction, 0, 10)

	for insFrom, v := range tg.Edges {
		for _, e := range v {
			ins2 := v.(*wtype.LHInstruction)
			if ins2 == ins {
				r = append(r, insFrom.(*wtype.LHInstruction))
				break
			}
		}
	}

	return r
}

func chainOuts(ic *IChain, tg graph.Graph) []*wtype.LHInstruction {
	// we don't really care about duplicates in this list so just take unions
	ret := make([]*wtype.LHInstruction, 0, 2*len(ic.Values))

	for _, v := range ret.Values {
		edge := tg.Edges[graph.Node(v)]

		for _, i := range edge {
			ret = append(ret, i.(*wtype.LHInstruction))
		}
	}

	return ret
}

func meetInsSets(ins, outs []*wtype.LHInstruction) []*wtype.LHInstruction {
	m := make(map[*wtype.LHInstruction]bool, len(ins))

	for _, i := range ins {
		m[i] = true
	}

	r := make([]string, 0, len(outs))

	for _, o := range outs {
		if m[o] {
			r = append(r, o)
		}
	}

	return r
}

func set_output_order(rq *LHRequest) error {
	// use topoSort
	tg := makeTGraph(rq.LHInstructions)

	// TODO --> add NodeOrder definition to allow output sorting
	sorted, err := graph.TopoSort(TopoSortOpt{Graph: tg})

	if err != nil {
		return err
	}

	it := convertToInstructionChain(sorted, rq.Options.OutputSort)
	rq.InstructionChain = it
	return nil
}

func set_output_order_orig(rq *LHRequest) error {
	// sort into equivalence classes by generation

	sorted := insSliceFromMap(rq.LHInstructions)

	sorted = aggregateAppropriateInstructions(sorted)

	if rq.Options.OutputSort {
		sort.Sort(ByGenerationOpt(sorted))
	} else {
		sort.Sort(ByGeneration(sorted))
	}

	it := NewIChain(nil)

	// aggregation of instructions effectively happens here. This entire level is
	// passed as a block to the instruction generator as a TransferBlock (TFB)
	// to be picked apart sequentially into sets which can be serviced simultaneously
	// etc.

	for _, v := range sorted {
		// fmt.Println("V: ", v.Result.CName, " ID: ", v.Result.ID, " PARENTS: ", v.ParentString(), " GENERATION: ", v.Generation())

		it.Add(v)
	}

	it.Print()

	rq.Output_order = it.Flatten()

	rq.InstructionChain = it

	//rq.InstructionSets = make_instruction_sets(it)

	return nil
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

		newtfr := insIn[ar[0]].(*driver.TransferInstruction)

		for k := 1; k < len(ar); k++ {
			newtfr.MergeWith(insIn[ar[k]].(*driver.TransferInstruction))
		}

		ret = append(ret, newtfr)
	}

	return ret
}

// TODO -- refactor this to pass robot through
func ConvertInstruction(insIn *wtype.LHInstruction, robot *driver.LHProperties, carryvol wunit.Volume) (insOut *driver.TransferInstruction, err error) {
	cmps := insIn.Components

	lenToMake := len(insIn.Components)

	/*	TODO -- remove
		fmt.Println("MIX (IN PLACE: ", insIn.IsMixInPlace(), ") CMPS ", len(cmps), " RES: ", insIn.Result.ID, " NAME: ", insIn.Result.CName, " ADDRESS: ", insIn.Welladdress)
		fmt.Println("FIRST CMPID: ", cmps[0].ID, " AND NAME ", cmps[0].CName)
		fmt.Println("---")
	*/
	if insIn.IsMixInPlace() {
		lenToMake = lenToMake - 1
		cmps = cmps[1:len(cmps)]
	}

	wh := make([]string, 0, lenToMake)       // component types
	va := make([]wunit.Volume, 0, lenToMake) // volumes

	fromPlateIDs, fromWellss, volss, err := robot.GetComponents(cmps, carryvol, wtype.LHVChannel, 1, true)

	if err != nil {
		return nil, err
	}

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

	for i, v := range cmps {
		for xx, _ := range fromPlateIDs[i] {
			// get dem big ole plates out
			// TODO -- pass them in instead of all this nonsense

			var flhp, tlhp *wtype.LHPlate

			flhif := robot.PlateLookup[fromPlateIDs[i][xx]]

			if flhif != nil {
				flhp = flhif.(*wtype.LHPlate)
			} else {
				s := fmt.Sprint("NO SRC PLATE FOUND : ", i, " ", xx, " ", fromPlateIDs[i][xx])
				err := wtype.LHError(wtype.LH_ERR_DIRE, s)

				return nil, err
			}

			tlhif := robot.PlateLookup[insIn.PlateID()]

			if tlhif != nil {
				tlhp = tlhif.(*wtype.LHPlate)
			} else {
				s := fmt.Sprint("NO DST PLATE FOUND : ", i, " ", xx, " ", insIn.PlateID())
				err := wtype.LHError(wtype.LH_ERR_DIRE, s)

				return nil, err
			}

			wlt, ok := tlhp.WellAtString(insIn.Welladdress)

			if !ok {
				return nil, wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprintf("Well %s not found on dest plate %s", insIn.Welladdress, insIn.PlateID()))
			}

			//v2 := wunit.NewVolume(v.Vol, v.Vunit)
			v2 := volss[i][xx]
			vt = append(vt, wlt.CurrVolume())
			wh = append(wh, v.TypeName())
			va = append(va, v2)
			pt = append(pt, robot.PlateIDLookup[insIn.PlateID()])
			wt = append(wt, insIn.Welladdress)
			ptwx = append(ptwx, tlhp.WellsX())
			ptwy = append(ptwy, tlhp.WellsY())
			ptt = append(ptt, tlhp.Type)

			wlf, ok := flhp.WellAtString(fromWellss[i][xx])

			if !ok {
				//logger.Fatal(fmt.Sprint("Well ", fromWells[ix], " not found on source plate ", fromPlateID[ix]))
				err = wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprint("Well ", fromWellss[i][xx], " not found on source plate ", fromPlateIDs[i][xx]))
				return nil, err
			}

			vf = append(vf, wlf.CurrVolume())
			vrm := v2.Dup()
			vrm.Add(carryvol)
			wlf.Remove(vrm)

			pf = append(pf, robot.PlateIDLookup[fromPlateIDs[i][xx]])
			wf = append(wf, fromWellss[i][xx])
			pfwx = append(pfwx, flhp.WellsX())
			pfwy = append(pfwy, flhp.WellsY())
			ptf = append(ptf, flhp.Type)

			if v.Loc == "" {
				v.Loc = fromPlateIDs[i][xx] + ":" + fromWellss[i][xx]
			}
			// add component to destination

			// ensure we keep results straight
			vd := v.Dup()
			vd.ID = wlf.WContents.ID
			vd.ParentID = wlf.WContents.ParentID
			wlt.Add(vd)

			// TODO -- danger here, is result definitely set?
			wlt.WContents.ID = insIn.Result.ID
			wlf.WContents.AddDaughterComponent(wlt.WContents)

			//fmt.Println("HERE GOES: ", i, wh[i], vf[i].ToString(), vt[i].ToString(), va[i].ToString(), pt[i], wt[i], pf[i], wf[i], pfwx[i], pfwy[i], ptwx[i], ptwy[i])

		}
	}

	// what, pltfrom, pltto, wellfrom, wellto, fplatetype, tplatetype []string, volume, fvolume, tvolume []wunit.Volume, FPlateWX, FPlateWY, TPlateWX, TPlateWY []int
	ti := driver.NewTransferInstruction(wh, pf, pt, wf, wt, ptf, ptt, va, vf, vt, pfwx, pfwy, ptwx, ptwy)

	return ti, nil
}
