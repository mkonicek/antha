// liquidhandling/layoutagent2.go: Part of the Antha language
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
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/sampletracker"
)

//ImprovedLayoutAgent assigns destinations to mix instructions
//don't ask about how bad the original one (upon which the 'improvements' here were made) was...
func ImprovedLayoutAgent(ctx context.Context, request *LHRequest, params *liquidhandling.LHProperties) error {
	// do this multiply based on the order in the chain

	ch := request.InstructionChain
	pc := make([]PlateChoice, 0, 3)
	mp := make(map[string]string)
	var err error

	// stage zero: seed in user plates if destinations are required
	pc = map_in_user_plates(request, pc)

	k := 1

	for {
		if ch == nil {
			break
		}
		pc, mp, err = LayoutStage(ctx, request, params, ch, pc, mp)

		k += 1
		if err != nil {
			return err
		}
		ch = ch.Child
	}

	// prune out dead instructions from the assignments
	filtered := make(map[string][]string)

	for k, insAr := range request.OutputAssignments {
		ar := make([]string, 0, len(insAr))
		for _, v := range insAr {
			_, ok := request.LHInstructions[v]

			if ok {
				ar = append(ar, v)
			}
		}

		if len(ar) != 0 {
			filtered[k] = ar
		}
	}

	request.OutputAssignments = filtered

	return err
}

func map_in_user_plates(rq *LHRequest, pc []PlateChoice) []PlateChoice {
	for _, p := range rq.InputPlates {
		pc = map_in_user_plate(p, pc, rq)
	}

	for _, p := range rq.OutputPlates {
		pc = map_in_user_plate(p, pc, rq)
	}

	return pc
}

func findInPC(ass, w string, pc PlateChoice) int {

	i := -1

	for ix := 0; ix < len(pc.Assigned); ix++ {
		if pc.Assigned[ix] == ass && pc.Wells[ix] == w {
			i = ix
			break
		}
	}

	return i
}

func map_in_user_plate(p *wtype.Plate, pc []PlateChoice, rq *LHRequest) []PlateChoice {
	nm := p.PlateName

	it := wtype.NewAddressIterator(p, wtype.ColumnWise, wtype.TopToBottom, wtype.LeftToRight, false)

	for wc := it.Curr(); it.Valid(); wc = it.Next() {
		w := p.Wellcoords[wc.FormatA1()]

		if w.IsEmpty() {
			continue
		}

		i := defined(p.ID, pc)
		cnt := w.WContents

		if i == -1 {
			pc = append(pc, PlateChoice{Platetype: p.Type, Assigned: []string{cnt.ID}, ID: p.ID, Wells: []string{wc.FormatA1()}, Name: nm, Output: []bool{false}})
		} else {
			ass := findInPC(cnt.ID, wc.FormatA1(), pc[i])

			if ass == -1 {
				pc[i].Assigned = append(pc[i].Assigned, cnt.ID)
				pc[i].Wells = append(pc[i].Wells, wc.FormatA1())
				pc[i].Output = append(pc[i].Output, false)
			}
		}
	}
	return pc
}

func getNameForID(pc []PlateChoice, id string) string {
	for _, p := range pc {
		if p.ID == id {
			return p.Name
		}
	}

	return fmt.Sprintf("Output_plate_%s", id[0:6])
}

func LayoutStage(ctx context.Context, request *LHRequest, params *liquidhandling.LHProperties, chain *wtype.IChain, plate_choices []PlateChoice, mapchoices map[string]string) ([]PlateChoice, map[string]string, error) {
	// considering only plate assignments,
	// we have three kinds of solution
	// 1- ones going to a specific plate
	// 2- ones going to a specific plate type
	// 3- ones going to a plate of our choosing

	// find existing assignments and copy into the plate_choices structure
	// this may be because 1) the user has set the assignment 2) the assignment derives from a component
	st := sampletracker.FromContext(ctx)
	plate_choices, mapchoices, err := getAndCompleteAssignments(st, request, chain.ValueIDs(), plate_choices, mapchoices)

	// map choices maps layout groups to (temp)plate IDs

	if err != nil {
		return nil, nil, err
	}
	// now we know what remains unassigned, we assign it

	plate_choices, err = choose_plates(ctx, request, plate_choices, chain.ValueIDs())
	if err != nil {
		return nil, nil, err
	}

	// now we have solutions of type 1 & 2

	// make specific plates... this may mean splitting stuff out into multiple plates

	remap, err := make_plates(ctx, request, chain.ValueIDs())
	if err != nil {
		return nil, nil, err
	}

	// I fix da map

	for k, v := range remap {
		for kk, vv := range mapchoices {
			if vv == k {
				mapchoices[kk] = v
			}
		}
	}

	// give them names

	for _, v := range request.OutputPlates {
		if wtype.NameOf(v) == "" {
			v.PlateName = getNameForID(plate_choices, v.ID)
		}
	}

	// now we have solutions of type 1 only -- we just need to
	// say where on each plate they will go
	// this needs to set OutputAssignments
	if err := make_layouts(ctx, request, plate_choices); err != nil {
		return nil, nil, err
	}

	lkp := make(map[string][]*wtype.Liquid)
	lk2 := make(map[string]string)
	// fix the output locations correctly

	//for _, v := range request.LHInstructions {
	order := chain.ValueIDs()
	for _, id := range order {
		v := request.LHInstructions[id]
		// pass ID through chain if not a mix
		if v.Type == wtype.LHIPRM {
			// the current contract on prompt instructions is to pass through a set of components
			// on which basis we need only make sure that each result has the same location
			// as its corresponding input

			for i := range v.Inputs {
				v.Outputs[i].Loc = v.Inputs[i].Loc
			}

			continue
		} else if v.Type == wtype.LHISPL {
			// similar to the above, just ensure the results both have the right location set
			v.Outputs[0].Loc = v.Inputs[0].Loc
			v.Outputs[1].Loc = v.Inputs[0].Loc
		}

		lkp[v.ID] = make([]*wtype.Liquid, 0, 1) //v.Output
		lk2[v.Outputs[0].ID] = v.ID
	}

	for _, id := range order {
		v := request.LHInstructions[id]
		for _, c := range v.Inputs {
			// if this component has the same ID
			// as the result of another instruction
			// we map it in
			iID, ok := lk2[c.ID]

			if ok {
				// iID is an instruction ID
				lkp[iID] = append(lkp[iID], c)
			}
		}

		// now we put the actual result in
		lkp[v.ID] = append(lkp[v.ID], v.Outputs[0])
	}

	// now map the output assignments in
	for k, v := range request.OutputAssignments {
		for _, id := range v {
			l := lkp[id]
			for _, x := range l {
				// x.Loc = k
				// also need to remap the plate id
				tx := strings.Split(k, ":")
				_, ok := remap[tx[0]]

				if ok {
					x.Loc = remap[tx[0]] + ":" + tx[1]
					st.SetLocationOf(x.ID, x.Loc)
				} else {
					x.Loc = tx[0] + ":" + tx[1]
					st.SetLocationOf(x.ID, x.Loc)
				}
			}
		}
	}

	// make sure plate choices is remapped
	for i, v := range plate_choices {
		_, ok := remap[v.ID]

		if ok {
			plate_choices[i].ID = remap[v.ID]
		}
	}

	return plate_choices, mapchoices, nil
}

type PlateChoice struct {
	Platetype string
	Assigned  []string
	ID        string
	Wells     []string
	Name      string
	Output    []bool
}

func getAndCompleteAssignments(st *sampletracker.SampleTracker, request *LHRequest, order []string, s []PlateChoice, m map[string]string) ([]PlateChoice, map[string]string, error) {
	//s := make([]PlateChoice, 0, 3)
	//m := make(map[int]string)

	// inconsistent plate types will be assigned randomly!
	x := 0
	for _, k := range order {
		x += 1
		v := request.LHInstructions[k]

		// ignore non-mixes
		if v.Type != wtype.LHIMIX {
			continue
		}

		// if plate ID set
		if v.PlateID != "" {
			//MixInto
			i := defined(v.PlateID, s)

			nm := v.PlateName

			if nm == "" {
				nm = fmt.Sprintf("Output_plate_%s", v.PlateID[0:6])
			}

			if i == -1 {
				s = append(s, PlateChoice{Platetype: v.Platetype, Assigned: []string{v.ID}, ID: v.PlateID, Wells: []string{v.Welladdress}, Name: nm, Output: []bool{true}})
			} else {

				s[i].Assigned = append(s[i].Assigned, v.ID)
				s[i].Wells = append(s[i].Wells, v.Welladdress)
				s[i].Output = append(s[i].Output, true)
			}

		} else if v.Majorlayoutgroup != -1 || v.PlateName != "" {
			//MixTo / MixNamed
			nm := "Output_plate"
			mlg := fmt.Sprintf("%d", v.Majorlayoutgroup)

			if mlg == "-1" {
				mlg = v.PlateName
				nm = v.PlateName
			}

			id, ok := m[mlg]
			// if no plate assigned so far, assign a temp ID for grouping
			if !ok {
				id = wtype.NewUUID()
				m[mlg] = id
				if nm == "Output_plate" {
					nm += "_" + id[0:6]
				}
			}

			//  fix the plate id to this temporary one
			request.LHInstructions[k].SetPlateID(id)

			i := defined(id, s)

			if i == -1 {
				s = append(s, PlateChoice{Platetype: v.Platetype, Assigned: []string{v.ID}, ID: id, Wells: []string{v.Welladdress}, Name: nm, Output: []bool{true}})
			} else {
				// check if this well is used... if so, we need another plate

				if v.Welladdress != "" && wutil.StrInStrArray(v.Welladdress, s[i].Wells) {

					// see if we can find a plate

					i = findPlateWithWellFree(s, v.Platetype, v.Welladdress, v.PlateName)

					if i == -1 {
						// a '-1' means we didn't find one
						id := wtype.NewUUID()
						request.LHInstructions[k].SetPlateID(id)
						s = append(s, PlateChoice{Platetype: v.Platetype, Assigned: []string{v.ID}, ID: v.PlateID, Wells: []string{v.Welladdress}, Name: nm, Output: []bool{true}})
						i = len(s) - 1
					}
				}

				s[i].Assigned = append(s[i].Assigned, v.ID)
				s[i].Wells = append(s[i].Wells, v.Welladdress)
				s[i].Output = append(s[i].Output, true)
			}
		} else if v.IsMixInPlace() {
			// the first component sets the destination
			// and now it should indeed be set

			// really?
			if len(v.Inputs) == 0 {
				continue
			}

			if v.Inputs[0].PlateLocation().ID == "" {
				addr, ok := st.GetLocationOf(v.Inputs[0].ID)

				if !ok {
					err := wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprintf("MIX IN PLACE WITH NO LOCATION SET FOR %s", v.Inputs[0].Name()))
					return s, m, err
				}

				v.Inputs[0].Loc = addr
			}

			addr := v.Inputs[0].Loc
			tx := strings.Split(addr, ":")

			// do we know about the plate?

			lookUp, ok := request.GetPlate(tx[0])

			if !ok {
				panic(fmt.Sprint("No information on plate ", tx[0], "  available for layout "))
			}

			request.LHInstructions[k].Welladdress = tx[1]
			request.LHInstructions[k].SetPlateID(tx[0])
			request.LHInstructions[k].Platetype = lookUp.Type
			request.LHInstructions[k].OutPlate = lookUp
			request.LHInstructions[k].Outputs[0].Loc = addr

			// same as condition 1 except we get the plate id somewhere else
			i := defined(tx[0], s)

			// we should check for it in OutputPlates as well
			// this could be a mix in place which has been split

			if i == -1 {
				fmt.Println("CONTRADICTORY PLATE ID SITUATION ", v)
			}
			for i2, v2 := range s[i].Wells {
				if v2 == tx[1] {
					/*
						if s[i].Output[i2] {
							s[i].Assigned[i2] = v.ID
						} else {
							s[i].Assigned[i2] = v.ProductIDs()[0]
						}
					*/
					s[i].Assigned[i2] = v.ID
					//		found = true
					break
				}
			}

		}

		//else {
		// bare mix
		// this is handled later
		//}
	}

	// make sure the plate choices all have defined types

	for i := range s {
		if s[i].Platetype == "" {
			s[i].Platetype = request.OutputPlatetypes[0].Type
		}
	}

	return s, m, nil
}

func defined(s string, pc []PlateChoice) int {
	r := -1

	for i, v := range pc {
		if v.ID == s {
			r = i
			break
		}
	}
	return r
}

func choose_plates(ctx context.Context, request *LHRequest, pc []PlateChoice, order []string) ([]PlateChoice, error) {
	for _, k := range order {
		v := request.LHInstructions[k]

		// ignore non-mix instructions

		if v.Type != wtype.LHIMIX {
			continue
		}

		// this id may be temporary, only things without it still are not assigned to a
		// plate, even a virtual one
		if v.PlateID == "" {
			pt := v.Platetype
			// find a plate choice to put it in or return -1 for a new one
			ass := -1

			if pt != "" {
				ass = assignmentWithType(pt, pc)
			} else if len(pc) != 0 {
				// just stick it in the first one
				ass = 0
			}

			if ass == -1 {
				// make a new plate
				if len(request.OutputPlatetypes) == 0 {
					return nil, fmt.Errorf("no output plate types specified. \n If not specifying output plate type in a Mix Command, at least one output plate type must be specified in config > outputPlateTypes.")
				}
				pc = append(pc, PlateChoice{Platetype: chooseAPlate(request, v), Assigned: []string{v.ID}, ID: wtype.GetUUID(), Wells: []string{""}, Name: "Output_plate_" + v.ID[0:6], Output: []bool{true}})
				continue
			}

			pc[ass].Assigned = append(pc[ass].Assigned, v.ID)
			pc[ass].Wells = append(pc[ass].Wells, "")
			pc[ass].Output = append(pc[ass].Output, true)
		}
	}

	// now we have everything assigned to virtual plates
	// make sure the plates aren't too full

	pc2 := make([]PlateChoice, 0, len(pc))

	for _, v := range pc {
		plate, err := inventory.NewPlate(ctx, v.Platetype)
		if err != nil {
			return nil, err
		}

		// chop the assignments up

		pc2 = append(pc2, modpc(v, plate.Nwells)...)
	}

	// copy the choices in

	for _, c := range pc2 {
		for _, i := range c.Assigned {
			_, ok := request.LHInstructions[i]

			if !ok {
				continue
			}
			request.LHInstructions[i].SetPlateID(c.ID)
			request.LHInstructions[i].Platetype = c.Platetype
			request.LHInstructions[i].PlateName = c.Name
		}
	}

	return pc2, nil
}

// chop the assignments up modulo plate size
func modpc(choice PlateChoice, nwell int) []PlateChoice {
	r := make([]PlateChoice, 0, 1)

	seen := make(map[string]bool)

	for s := 0; s < len(choice.Assigned); s += nwell {
		e := s + nwell
		if e > len(choice.Assigned) {
			e = len(choice.Assigned)
		}
		ID := choice.ID
		if s != 0 {
			// new ID
			ID = wtype.GetUUID()
		}

		nm := uniquePlateName(choice.Name, seen, 100)

		r = append(r, PlateChoice{Platetype: choice.Platetype, Assigned: choice.Assigned[s:e], ID: ID, Wells: choice.Wells[s:e], Name: nm, Output: choice.Output[s:e]})
	}
	return r
}

func uniquePlateName(namein string, seen map[string]bool, maxtries int) string {
	nm := namein

	_, ok := seen[nm]

	if ok {
		for k := 0; k < maxtries; k++ {
			nm2 := fmt.Sprintf("%s%d", nm, k+2)
			_, ok = seen[nm2]
			if !ok {
				nm = nm2
				break
			}
		}

		if ok {
			fmt.Printf("Tried to assign more than %d output plates\n", maxtries)
		}

	}

	seen[nm] = true

	return nm
}

func assignmentWithType(pt string, pc []PlateChoice) int {
	r := -1

	if pt == "" {
		if len(pc) != 0 {
			//r = 0
			// assume previous plates are all full
			// TODO -- much more sensible choice method
			r = len(pc) - 1
		}
		return r
	}

	for i, v := range pc {
		if pt == v.Platetype {
			r = i
			//			break
		}
	}

	return r
}

func chooseAPlate(request *LHRequest, ins *wtype.LHInstruction) string {
	// for now we ignore ins and just choose the First Output Platetype
	return request.OutputPlatetypes[0].Type
}

// we have potentially added extra theoretical plates above
// now we make real plates and swap them in

func make_plates(ctx context.Context, request *LHRequest, order []string) (map[string]string, error) {
	remap := make(map[string]string)
	//for k, v := range request.LHInstructions {
	for _, k := range order {
		v := request.LHInstructions[k]

		// ignore non-mix instructions

		if v.Type != wtype.LHIMIX {
			continue
		}

		_, skip := remap[v.PlateID]

		if skip {
			request.LHInstructions[k].SetPlateID(remap[v.PlateID])
			continue
		}
		_, ok := request.OutputPlates[v.PlateID]

		// we don't remap input plates
		_, ok2 := request.InputPlates[v.PlateID]

		// need to assign a new plate
		if !(ok || ok2) {
			plate, err := inventory.NewPlate(ctx, v.Platetype)
			if err != nil {
				return nil, fmt.Errorf("cannot make plate %s: %s", v.Platetype, err)
			}
			plate.PlateName = request.LHInstructions[k].PlateName
			request.OutputPlates[plate.ID] = plate
			remap[v.PlateID] = plate.ID
			request.LHInstructions[k].SetPlateID(remap[v.PlateID])
		}

	}

	return remap, nil
}

func make_layouts(ctx context.Context, request *LHRequest, pc []PlateChoice) error {
	// we need to fill in the platechoice structure then
	// transfer the info across to the solutions
	//opa := request.OutputAssignments
	opa := make(map[string][]string)

	for _, c := range pc {
		// make a temporary plate to hold info

		plat, err := inventory.NewPlate(ctx, c.Platetype)
		if err != nil {
			return err
		}

		// make an iterator for it

		it := request.OutputIteratorFactory(plat)

		//put a dummy component in the assigned wells to mark them as used

		for _, w := range c.Wells {
			if w != "" {
				wc := wtype.MakeWellCoords(w)

				well, ok := plat.WellAt(wc)
				if !ok {
					return wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprintf("well (%s) specified is out of range of available wells for plate type %s", w, plat.Type))
				}
				err := markWellUsed(well)
				if err != nil {
					return err
				}
			}
		}

		for i := range c.Assigned {
			sID := c.Assigned[i]
			well := ""
			if i < len(c.Wells) {
				well = c.Wells[i]
			}

			var assignment string

			if well == "" {
				wc := plat.NextEmptyWell(it)
				well, ok := plat.WellAt(wc)
				if !ok {
					return wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprintf("too many assignments made to output plate \"%s\"", c.Platetype))
				}

				err := markWellUsed(well)
				if err != nil {
					return err
				}

				request.LHInstructions[sID].Welladdress = wc.FormatA1()
				assignment = c.ID + ":" + wc.FormatA1()
				c.Wells[i] = wc.FormatA1()
			} else {
				assignment = c.ID + ":" + well
			}

			opa[assignment] = append(opa[assignment], sID)
		}
	}

	request.OutputAssignments = opa
	return nil
}

//markWellUsed add a dummy component to the well so that it's marked as having been used
func markWellUsed(well *wtype.LHWell) error {
	//avoid adding a dummy component if one's already been added
	if well.IsEmpty() {
		dummycmp := wtype.NewLHComponent()
		dummycmp.SetVolume(well.MaxVolume())
		err := well.AddComponent(dummycmp)
		if err != nil {
			return wtype.LHError(wtype.LH_ERR_VOL, fmt.Sprintf("Layout Agent : %s", err.Error()))
		}
	}
	return nil
}

//findPlateWithWellFree(s, v.Platetype, v.Welladdress, v.PlateName)

//findPlateWithWellFree looks in our array of plate choices to see if there already exists a plate of this type with this well free
// optionally we can specify a name
func findPlateWithWellFree(plateChoices []PlateChoice, plateType, wellAddress, plateName string) int {
	// -1 indicates not found
	ret := -1

	for i := 0; i < len(plateChoices); i++ {
		pc := plateChoices[i]
		nm := pc.Name

		// ensure that if name is empty it does not act as a constraint
		if plateName == "" {
			nm = ""
		}

		if pc.Platetype == plateType && nm == plateName && !wutil.StrInStrArray(wellAddress, pc.Wells) {
			ret = i
			break
		}
	}

	return ret
}
