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
)

func ConvertInstructions(inssIn LHIVector, robot *LHProperties, carryvol wunit.Volume, channelprms *wtype.LHChannelParameter, multi int) (insOut []*TransferInstruction, err error) {
	for i := 0; i < inssIn.MaxLen(); i++ {
		comps := inssIn.CompsAt(i)

		lenCmps := 0

		for _, ccc := range comps {
			if ccc != nil {
				lenCmps += 1
			}
		}
		if lenCmps == 0 {
			continue
		}

		fromPlateID, fromWells, err := robot.GetComponents(comps, carryvol, channelprms.Orientation, multi)

	}

}

/*
func yeahyeahyeah(){
	cmps := insIn.Components

	lenToMake := len(insIn.Components)

	if insIn.IsMixInPlace() {
		lenToMake = lenToMake - 1
		cmps = cmps[1:len(cmps)]
	}

	wh := make([]string, lenToMake)       // component types
	va := make([]wunit.Volume, lenToMake) // volumes

	// six parameters applying to the source

	fromPlateID, fromWells, err := robot.GetComponents(cmps, carryvol, channelprms.Orientation, multi)

	if err != nil {
		return nil, err
	}

	pf := make([]string, lenToMake)
	wf := make([]string, lenToMake)
	pfwx := make([]int, lenToMake)
	pfwy := make([]int, lenToMake)
	vf := make([]wunit.Volume, lenToMake)
	ptt := make([]string, lenToMake)

	// six parameters applying to the destination

	pt := make([]string, lenToMake)       // dest plate positions
	wt := make([]string, lenToMake)       // dest wells
	ptwx := make([]int, lenToMake)        // dimensions of plate pipetting to (X)
	ptwy := make([]int, lenToMake)        // dimensions of plate pipetting to (Y)
	vt := make([]wunit.Volume, lenToMake) // volume in well to
	ptf := make([]string, lenToMake)      // plate types

	ix := 0

	for i, v := range insIn.Components {
		if insIn.IsMixInPlace() && i == 0 {
			continue
		}

		// get dem big ole plates out
		// TODO -- pass them in instead of all this nonsense

		var flhp, tlhp *wtype.LHPlate

		flhif := robot.PlateLookup[fromPlateID[ix]]

		if flhif != nil {
			flhp = flhif.(*wtype.LHPlate)
		} else {
			s := fmt.Sprint("NO SRC PLATE FOUND : ", ix, " ", fromPlateID[ix])
			err := wtype.LHError(wtype.LH_ERR_DIRE, s)

			return nil, err
		}

		tlhif := robot.PlateLookup[insIn.PlateID()]

		if tlhif != nil {
			tlhp = tlhif.(*wtype.LHPlate)
		} else {
			s := fmt.Sprint("NO DST PLATE FOUND : ", ix, " ", insIn.PlateID())
			err := wtype.LHError(wtype.LH_ERR_DIRE, s)

			return nil, err
		}

		wlt, ok := tlhp.WellAtString(insIn.Welladdress)

		if !ok {
			err = wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprint("Well ", insIn.Welladdress, " not found on dest plate ", insIn.PlateID))
			return nil, err
		}

		v2 := wunit.NewVolume(v.Vol, v.Vunit)
		vt[ix] = wlt.CurrVolume()
		wh[ix] = v.TypeName()
		va[ix] = v2
		pt[ix] = robot.PlateIDLookup[insIn.PlateID()]
		wt[ix] = insIn.Welladdress
		ptwx[ix] = tlhp.WellsX()
		ptwy[ix] = tlhp.WellsY()
		ptt[ix] = tlhp.Type

		wlf, ok := flhp.WellAtString(fromWells[ix])

		if !ok {
			//logger.Fatal(fmt.Sprint("Well ", fromWells[ix], " not found on source plate ", fromPlateID[ix]))
			err = wtype.LHError(wtype.LH_ERR_DIRE, fmt.Sprint("Well ", fromWells[ix], " not found on source plate ", fromPlateID[ix]))
			return nil, err
		}

		vf[ix] = wlf.CurrVolume()
		//wlf.Remove(va[ix])

		pf[ix] = robot.PlateIDLookup[fromPlateID[ix]]
		wf[ix] = fromWells[ix]
		pfwx[ix] = flhp.WellsX()
		pfwy[ix] = flhp.WellsY()
		ptf[ix] = flhp.Type

		if v.Loc == "" {
			v.Loc = fromPlateID[ix] + ":" + fromWells[ix]
		}
		// add component to destination
		// need to ensure data are consistent
		vd := v.Dup()
		vd.ID = wlf.WContents.ID
		vd.ParentID = wlf.WContents.ParentID
		wlt.Add(vd)

		// add daughter ID to component in

		wlf.WContents.AddDaughterComponent(wlt.WContents)

		//fmt.Println("HERE GOES: ", i, wh[i], vf[i].ToString(), vt[i].ToString(), va[i].ToString(), pt[i], wt[i], pf[i], wf[i], pfwx[i], pfwy[i], ptwx[i], ptwy[i])

		ix += 1
	}

	ti := NewTransferInstruction(wh, pf, pt, wf, wt, ptf, ptt, va, vf, vt, pfwx, pfwy, ptwx, ptwy)

	// what, pltfrom, pltto, wellfrom, wellto, fplatetype, tplatetype []string, volume, fvolume, tvolume []wunit.Volume, FPlateWX, FPlateWY, TPlateWX, TPlateWY []int
	return ti, nil
}
*/
