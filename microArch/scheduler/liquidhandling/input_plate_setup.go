// anthalib//liquidhandling/input_plate_setup.go: Part of the Antha language
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
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/dustinkirkland/golang-petname"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/microArch/sampletracker"
)

type InputSorter struct {
	Ordered []string
	Values  map[string]wunit.Volume
}

// @implement sort.Interface
func (is InputSorter) Len() int {
	return len(is.Ordered)
}

func (is InputSorter) Swap(i, j int) {
	s := is.Ordered[i]
	is.Ordered[i] = is.Ordered[j]
	is.Ordered[j] = s
}

func (is InputSorter) Less(i, j int) bool {
	vv1 := is.Values[is.Ordered[i]]
	vv2 := is.Values[is.Ordered[j]]

	v1 := vv1.SIValue()
	v2 := vv2.SIValue()

	// we want ascending sort here
	if v1 < v2 {
		return false
	} else if v1 > v2 {
		return true
	}

	// volumes are equal

	ss := sort.StringSlice(is.Ordered)

	return ss.Less(i, j)
}

// inputPlateSetup map input liquids to input plates
// INPUT: 	"input_platetype", "inputs"
//OUTPUT: 	"input_plates"      -- these each have components in wells
//		"input_assignments" -- map with arrays of assignment strings, i.e. {tea: [plate1:A:1, plate1:A:2...] }etc.
func (rq *LHRequest) inputPlateSetup(ctx context.Context, carryVolume wunit.Volume) error {
	st := sampletracker.FromContext(ctx)

	input_platetypes := rq.InputPlatetypes

	// we assume that input_plates is set if any locs are set
	input_plates := rq.InputPlates

	if len(input_plates) == 0 {
		input_plates = make(map[string]*wtype.Plate, 3)
	}

	// need to fill each plate type

	var curr_plate *wtype.Plate

	inputs := rq.InputSolutions.Solutions

	input_order := make([]string, len(rq.InputSolutions.Order))
	copy(input_order, rq.InputSolutions.Order)

	// this needs to be passed in via the request... must specify how much of inputs cannot
	// be satisfied by what's already passed in

	input_volumes := rq.InputSolutions.VolumesWanting

	// sort to make deterministic
	// we sort by a) volume (descending) b) name (alphabetically)

	isrt := InputSorter{input_order, input_volumes}

	sort.Sort(isrt)

	input_order = isrt.Ordered

	weights_constraints := rq.InputSetupWeights

	// get the assignment

	var well_count_assignments map[string]map[*wtype.Plate]int

	if len(input_volumes) != 0 {
		// If any input solutions need to be set up then we now check if there any input plate types set.
		if len(input_platetypes) == 0 {
			return fmt.Errorf("no input plate set: \n  - Please upload plate file or select at least one input plate type in Configuration > Preferences > inputPlateTypes. \n - Important: Please add a riser to the plate choice for low profile plates such as PCR plates, 96 and 384 well plates. ")
		}
		var err error
		well_count_assignments, err = choosePlateAssignments(input_volumes, input_platetypes, weights_constraints)

		if err != nil {
			return err
		}
	}

	input_assignments := make(map[string][]string, len(well_count_assignments))

	plates_in_play := make(map[string]*wtype.Plate)

	curplaten := 1
	for _, cname := range input_order {
		volume, ok := input_volumes[cname]

		if !ok {
			continue
		}

		// this needs to get the right thing:
		// -- anonymous components are fine but
		//    identified ones need to come out correctly
		component := inputs[cname][0]

		well_assignments, ok := well_count_assignments[cname]

		// is this really OK?!
		if !ok {
			continue
		}

		// check here
		if isInstance(cname) && len(well_assignments) != 1 {
			return fmt.Errorf("Error: Autoallocated mix-in-place components cannot be spread across multiple wells")
		}

		var curr_well *wtype.LHWell
		var assignments []string

		// best hack so far: add an extra well of everything
		// in case we run out
		for platetype, nwells := range well_assignments {
			WellTot := nwells + 1

			// unless it's an instance
			if isInstance(cname) {
				WellTot = nwells
			}

			for i := 0; i < WellTot; i++ {
				curr_plate = plates_in_play[platetype.Type]

				if curr_plate == nil {
					p, err := inventory.NewPlate(ctx, platetype.Type)
					if err != nil {
						return err
					}

					plates_in_play[platetype.Type] = p
					curr_plate = plates_in_play[platetype.Type]
					platename := rq.getSafeInputPlateName(curplaten)
					curr_plate.PlateName = platename
					curplaten += 1
					//curr_plate.DeclareAutoallocated()
				}

				// find somewhere to put it
				curr_well, ok = wtype.Get_Next_Well(curr_plate, component, curr_well)

				if !ok {
					// if no space, reset
					plates_in_play[platetype.Type] = nil
					curr_plate = nil
					curr_well = nil
					i -= 1
					continue
				}

				// now put it there

				location := curr_plate.ID + ":" + curr_well.Crds.FormatA1()
				assignments = append(assignments, location)

				var newcomponent *wtype.Liquid

				if isInstance(cname) {
					newcomponent = component
					newcomponent.Loc = location
					// don't let these get deleted...
					curr_well.SetUserAllocated()
				} else {
					newcomponent = component.Dup()
					newcomponent.Vol = curr_well.MaxVolume().RawValue()
					newcomponent.Vunit = curr_well.MaxVolume().Unit().PrefixedSymbol()
					newcomponent.Loc = location

					//usefulVolume is the most we can get from the well assuming one transfer
					usefulVolume := curr_well.CurrentWorkingVolume()
					usefulVolume.Subtract(carryVolume)
					volume.Subtract(usefulVolume)
				}

				st.SetLocationOf(component.ID, location)

				err := curr_well.AddComponent(newcomponent)
				if err != nil {
					return wtype.LHError(wtype.LH_ERR_VOL, fmt.Sprintf("Input plate setup : %s", err.Error()))
				}
				curr_well.DeclareAutoallocated()
				input_plates[curr_plate.ID] = curr_plate
			}
		}

		input_assignments[cname] = assignments
	}

	// add any remaining assignments

	for _, v := range inputs {
		for _, vv := range v {
			// this now means input assignments is always set...
			// previously this was empty
			if vv.Loc != "" && !vv.Volume().IsZero() {
				// append it
				input_assignments[vv.CName] = append(input_assignments[vv.CName], vv.Loc)
			}
		}
	}

	rq.InputPlates = input_plates
	rq.InputAssignments = input_assignments

	return nil
}

func isInstance(s string) bool {
	// we need to forbid this prefix in component names
	if strings.HasPrefix(s, "CNID:") {
		return true
	} else {
		return false
	}
}

// returns a unique plate name
func (rq *LHRequest) getSafeInputPlateName(curplaten int) string {
	return rq.getSafePlateName("auto_input_plate", "_", curplaten)
}

func (rq *LHRequest) getSafePlateName(prefix, sep string, curplaten int) string {
	trialPlateName := randomPlateName(prefix, sep, curplaten)

	for {
		if rq.HasPlateNamed(trialPlateName) {
			trialPlateName = randomPlateName(prefix, sep, curplaten)

		} else {
			break
		}
	}

	return trialPlateName
}

func randomPlateName(prefix, sep string, order int) string {

	blackListed := []string{"crappie", "titmouse", "stinkbug"}

	randomName := petname.Generate(1, "")
	for search.InStrings(blackListed, randomName) {
		randomName = petname.Generate(1, "")
	}
	tox := []string{prefix, fmt.Sprintf("%d", order), randomName}
	return strings.Join(tox, sep)
}
