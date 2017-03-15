// liquidhandling/convertinstructions.go Part of the Antha language
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
)

//	this section aggregates instructions with the following constraints:
//
//	1) obey any requirement to do one sample at a time
//		-- bullet bitten: we cannot permit transfer to split up any multichannel instructions
//		   into singles here... this is a bit tricky but we must make it so
//		   some revision to how pragmas work may be needed: extend only to component type etc.
//
//	here is what a single sample assembled one component at a time looks like
//	|
//	|		here is one sample assembled all at once looks like
//	|		|
//	i1(A)		i2(A B C)	--> the LHIVector contains these two, maxlen = 3, CmpAt (0) = [A A]
//	--										  CmpAt (1) = [  B]
//	i3(B) <------									  CmpAt (2) = [  C]
//	--          |-- these two are done separately (so they're boring)
//	i4(C) <------
//
// 	this should produce the output:
//	TFR(A A d1 d2), TFR(B d2), TFR(C d2), TFR(B d1), TFR(C d1)
//	iow it does i1 + first part of i2 in parallel, then the rest of i2 then i3 then i4

// 	issue is we cannot tolerate this situation
//
//	i1(A)		i2(A B)		i3(A C)
//	so we have to ensure the components line up

func readableComponentArray(arr []*wtype.LHComponent) string {
	ret := ""

	for i, v := range arr {
		if v != nil {
			ret += fmt.Sprintf("%s:%-6.2f%s", v.CName, v.Vol, v.Vunit)
		} else {
			ret += "_nil_"
		}
		if i < len(arr)-1 {

			ret += ", "
		}
	}

	return ret
}

//
func ConvertInstructions(inssIn LHIVector, robot *LHProperties, carryvol wunit.Volume, channelprms *wtype.LHChannelParameter, multi int) (insOut []*TransferInstruction, err error) {
	insOut = make([]*TransferInstruction, 0, 1)

	// TODO TODO TODO
	// -- we have to choose channels somewhere... probably during the *generate* method
	//    of transferinstruction

	for i := 0; i < inssIn.MaxLen(); i++ {
		cmps := inssIn.CompsAt(i)
		lenToMake := 0
		for _, c := range cmps {
			if c != nil {
				lenToMake += 1
			}
		}

		if lenToMake == 0 {
			// don't make empty transfers
			continue
		}

		orientation := wtype.LHVChannel
		independent := false

		if channelprms != nil {
			orientation = channelprms.Orientation
			independent = channelprms.Independent
		}

		// the alignment here just says component i comes from fromWells[i]
		// it says nothing about which channel should be used
		fromPlateIDs, fromWells, vols, err := robot.GetComponents(cmps, carryvol, orientation, multi, independent)

		if err != nil {
			return nil, err
		}

		// mt counts up the arrays got by GetComponents
		for mt := 0; mt < len(fromPlateIDs); mt++ {
			wh := make([]string, len(cmps))       //	what
			pf := make([]string, len(cmps))       //	position from
			pt := make([]string, len(cmps))       //	position to
			wf := make([]string, len(cmps))       //	well from
			wt := make([]string, len(cmps))       //	well to
			ptf := make([]string, len(cmps))      //	plate type from
			ptt := make([]string, len(cmps))      //	plate type to
			va := make([]wunit.Volume, len(cmps)) //	volume
			vf := make([]wunit.Volume, len(cmps)) //	volume in well from
			vt := make([]wunit.Volume, len(cmps)) //	volume in well to
			pfwx := make([]int, len(cmps))        //	plate from wells x
			pfwy := make([]int, len(cmps))        //	  "     "    "   y
			ptwx := make([]int, len(cmps))        //	  "    to    "   x
			ptwy := make([]int, len(cmps))        //	  "     "    "   y

			for ci := 0; ci < len(cmps); ci++ {
				if fromPlateIDs[mt][ci] == "" {
					continue
				}

				// what type is this component?

				wh[ci] = cmps[ci].TypeName()

				// source plate position

				ppf, ok := robot.PlateIDLookup[fromPlateIDs[mt][ci]]

				if !ok {
					return insOut, wtype.LHError(wtype.LH_ERR_DIRE, "Planning inconsistency: input plate ID not found on robot - please report this error")
				}

				pf[ci] = ppf

				// destination plate position

				ppt, ok := robot.PlateIDLookup[inssIn[ci].PlateID()]

				if !ok {
					return insOut, wtype.LHError(wtype.LH_ERR_DIRE, "Planning inconsistency: destination plate ID not found on robot - please report this error")
				}

				pt[ci] = ppt

				// source well

				wf[ci] = fromWells[mt][ci]

				// destination well

				wt[ci] = inssIn[ci].Welladdress

				// source plate type

				srcPlate, ok := robot.Plates[ppf]

				if !ok {
					return insOut, wtype.LHError(wtype.LH_ERR_DIRE, "Planning inconsistency: input plate ID not found on robot (#2) - please report this error")
				}

				ptf[ci] = srcPlate.Type

				// destination plate type

				dstPlate, ok := robot.Plates[ppt]

				if !ok {
					return insOut, wtype.LHError(wtype.LH_ERR_DIRE, "Planning inconsistency: destination plate ID not found on robot - please report this error")
				}

				ptt[ci] = dstPlate.Type

				// volume being moved

				va[ci] = vols[mt][ci]

				// source well volume

				wellFrom, ok := srcPlate.Wellcoords[wf[ci]]

				if !ok {
					return insOut, wtype.LHError(wtype.LH_ERR_DIRE, "Planning inconsistency: source well not found on source plate - plate report this error")
				}

				vf[ci] = wellFrom.CurrVolume()

				// dest well volume

				wellTo, ok := dstPlate.Wellcoords[wt[ci]]

				if !ok {
					return insOut, wtype.LHError(wtype.LH_ERR_DIRE, "Planning inconsistency: dest well not found on dest plate - please report this error")
				}

				vt[ci] = wellTo.CurrVolume()

				// source plate dimensions

				pfwx[ci] = srcPlate.WellsX()
				pfwy[ci] = srcPlate.WellsY()

				// dest plate dimensions

				ptwx[ci] = dstPlate.WellsX()
				ptwy[ci] = dstPlate.WellsY()

				// do the bookkeeping

				cmpFrom := wellFrom.Remove(va[ci])

				if cmpFrom == nil {
					return insOut, wtype.LHError(wtype.LH_ERR_DIRE, "Planning inconsistency: src well does not contain sufficient volume - please report this error")
				}

				wellTo.Add(cmpFrom)

			}

			tfr := NewTransferInstruction(wh, pf, pt, wf, wt, ptf, ptt, va, vf, vt, pfwx, pfwy, ptwx, ptwy)
			insOut = append(insOut, tfr)
		}
	}

	return insOut, nil
}
