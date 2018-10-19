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
)

const (
	COLWISE = iota
	ROWWISE
	RANDOM
)

func convertToInstructionChain(sortedNodes []graph.Node, tg graph.Graph) *IChain {
	ic := NewIChain(nil)

	// the nodes are now ordered according to dependency relations
	// *IN REVERSE ORDER*

	// this routine defines equivalence classes of nodes

	for _, n := range sortedNodes {
		addToIChain(ic, n, tg)
	}

	// we need to ensure that splits, prompts and mixes are kept separate by fissioning nodes

	ic = ic.SplitMixedNodes()

	return ic
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
			ins.AddOutput(wtype.NewLHComponent())
			for _, ins2 := range ar {
				for _, cmp := range ins2.Inputs {
					ins.Inputs = append(ins.Inputs, cmp)
					ins.PassThrough[cmp.ID] = ins2.Outputs[0]
				}
			}
			insOut = append(insOut, graph.Node(ins))
		}

	}

	return insOut
}

//buildInstructionChain guarantee all nodes are dependency-ordered
//in order to aggregate without introducing cycles
func buildInstructionChain(unsorted map[string]*wtype.LHInstruction) (*IChain, error) {
	unsortedSlice := make([]*wtype.LHInstruction, 0, len(unsorted))
	for _, instruction := range unsorted {
		unsortedSlice = append(unsortedSlice, instruction)
	}

	tg, err := MakeTGraph(unsortedSlice)
	if err != nil {
		return nil, err
	}

	sorted, err := graph.TopoSort(graph.TopoSortOpt{Graph: tg})
	if err != nil {
		return nil, err
	}

	sortedAsIns := make([]*wtype.LHInstruction, len(sorted))
	for i := 0; i < len(sorted); i++ {
		sortedAsIns[i] = sorted[i].(*wtype.LHInstruction)
	}

	sorted = aggregatePromptsWithSameMessage(sortedAsIns, tg)

	// aggregate sorted again
	sortedAsIns = make([]*wtype.LHInstruction, len(sorted))
	for i, nIns := range sorted {
		sortedAsIns[i] = nIns.(*wtype.LHInstruction)
	}

	// sort again post aggregation
	tg, err = MakeTGraph(sortedAsIns)
	if err != nil {
		return nil, err
	}

	sorted, err = graph.TopoSort(graph.TopoSortOpt{Graph: tg})
	if err != nil {
		return nil, err
	}

	// make into equivalence classes and sort according to defined order
	ic := convertToInstructionChain(sorted, tg)

	return ic, nil
}

type ByOrdinal [][]int

func (bo ByOrdinal) Len() int      { return len(bo) }
func (bo ByOrdinal) Swap(i, j int) { bo[i], bo[j] = bo[j], bo[i] }
func (bo ByOrdinal) Less(i, j int) bool {
	// just compare the first one

	return bo[i][0] < bo[j][0]
}

// TODO -- refactor this to pass robot through
func ConvertInstruction(insIn *wtype.LHInstruction, robot *driver.LHProperties, carryvol wunit.Volume, legacyVolume bool) (insOut *driver.TransferInstruction, err error) {
	cmps := insIn.Inputs

	lenToMake := len(insIn.Inputs)

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

			var flhp, tlhp *wtype.Plate

			flhif := robot.PlateLookup[tfrs[i].PlateIDs[xx]] //[fromPlateIDs[i][xx]]

			if flhif != nil {
				flhp = flhif.(*wtype.Plate)
			} else {
				s := fmt.Sprint("NO SRC PLATE FOUND : ", i, " ", xx, " ", tfrs[i].PlateIDs[xx]) //fromPlateIDs[i][xx])
				err := wtype.LHError(wtype.LH_ERR_DIRE, s)

				return nil, err
			}

			tlhif := robot.PlateLookup[insIn.PlateID]

			if tlhif != nil {
				tlhp = tlhif.(*wtype.Plate)
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
			wlt.WContents.ID = insIn.Outputs[0].ID
			wlf.WContents.AddDaughterComponent(wlt.WContents)

			//fmt.Println("HERE GOES: ", i, wh[i], vf[i].ToString(), vt[i].ToString(), va[i].ToString(), pt[i], wt[i], pf[i], wf[i], pfwx[i], pfwy[i], ptwx[i], ptwy[i])

		}
	}

	// what, pltfrom, pltto, wellfrom, wellto, fplatetype, tplatetype []string, volume, fvolume, tvolume []wunit.Volume, FPlateWX, FPlateWY, TPlateWX, TPlateWY []int
	ti := driver.NewTransferInstruction(wh, pf, pt, wf, wt, ptf, ptt, va, vf, vt, pfwx, pfwy, ptwx, ptwy, cnames, policies)

	return ti, nil
}
