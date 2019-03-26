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
	"context"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

//
//	at this point (i.e. in a TransferBlock) the instructions have potentially been grouped into sets
//	with simultaneously servicable destinations - row or column-wise depending on the head
//	orientation chosen
//
//	The main goal here is to find sources in appropriate structure (rows or columns)
//	to allow them to be done simultaneously. This essentially follows a greedy strategy
//	in which the required components are aligned to the available sources to see how many
//	can be taken at once. There are some tricks involved to make this work with troughs
//	etc.
//

func ConvertInstructions(ctx context.Context, inssIn LHIVector, robot *LHProperties, channelprms *wtype.LHChannelParameter, multi int, legacyVolume bool, policy *wtype.LHPolicyRuleSet) ([]*TransferInstruction, error) {
	// we call convertInstructions twice because
	// 1) calling convertInstructions with multi = 8 when there are no actual multichannel instructions causes
	//    undesirable source volume selection, see tests "TestExecutionPlanning/single_channel_well_use", and
	//    "TestExecutionPlanning/single_channel_auto_allocation"
	// 2) convertInstructions makes changes to robot, meaning that it must be called exactly once with the the copy of robot passed to the function
	if transfers, err := convertInstructions(inssIn, robot.DupKeepIDs(), channelprms, multi, legacyVolume); err != nil {
		return nil, err
	} else if hasMCB, err := hasMultiChannelBlock(ctx, transfers, robot, policy); err != nil {
		return nil, err
	} else if hasMCB {
		return convertInstructions(inssIn, robot, channelprms, multi, legacyVolume)
	} else {
		return convertInstructions(inssIn, robot, channelprms, 1, legacyVolume)
	}
}

func hasMultiChannelBlock(ctx context.Context, tfrs []*TransferInstruction, rbt *LHProperties, policy *wtype.LHPolicyRuleSet) (bool, error) {
	for _, tfr := range tfrs {
		instrx, err := tfr.Dup().Generate(ctx, policy, rbt)

		if err != nil {
			return false, err
		}

		for _, ins := range instrx {
			if mcb, ok := ins.(*ChannelBlockInstruction); ok {
				if mcb.MaxMulti() > 1 {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func singlemakeComponentSets(inssIn LHIVector) ([][]*wtype.Liquid, []LHIVector) {
	var componentsToMove [][]*wtype.Liquid
	var instructionsToUse []LHIVector
	//make lists of components to attempt to transfer simultaneously
	for i := 0; i < len(inssIn); i++ {
		if inssIn[i] == nil {
			continue
		}

		cmps := inssIn[i].ComponentsMoving()
		for _, cmp := range cmps {
			if cmp != nil {
				if cmp.CName == "" {
					panic("COMPONENTS MUST HAVE NAMES")
				}
				componentsToMove = append(componentsToMove, []*wtype.Liquid{cmp})
				instructionsToUse = append(instructionsToUse, LHIVector{inssIn[i]})
			}
		}
	}

	return componentsToMove, instructionsToUse
}

func multimakeComponentSets(inssIn LHIVector) ([][]*wtype.Liquid, []LHIVector) {
	var componentsToMove [][]*wtype.Liquid
	var instructionsToUse []LHIVector
	//make lists of components to attempt to transfer simultaneously
	for i := 0; i < inssIn.MaxLen(); i++ {
		var inssToUse LHIVector
		cmps := inssIn.CompsAt(i)
		inssToUse = inssIn
		lenToMake := 0

		for _, c := range cmps {
			if c != nil {
				if c.CName == "" {
					panic("COMPONENTS MUST HAVE NAMES")
				}
				lenToMake += 1
			}
		}

		if lenToMake == 0 {
			// don't make empty transfers
			continue
		}

		componentsToMove = append(componentsToMove, cmps)
		instructionsToUse = append(instructionsToUse, inssToUse)
	}

	return componentsToMove, instructionsToUse

}

// this function takes a set of instructions and, depending on whether parallelisation
// is chosen, either groups components into simultaneously movable chunks or splits
// the components up, to be found one at a time. The dichotomy here is necessary to
// avoid this part overriding mix ordering requested by the user since the source matching
// routine favours higher-volume components over lower ones, so that if all components
// for single-channeling are fed in as a block they will actually be moved from highest
// to lowest volume rather than preserving order.
func makeComponentSets(inssIn LHIVector, multi int) ([][]*wtype.Liquid, []LHIVector) {
	if multi == 1 {
		return singlemakeComponentSets(inssIn)
	} else {
		return multimakeComponentSets(inssIn)
	}

}

func convertInstructions(inssIn LHIVector, robot *LHProperties, channelprms *wtype.LHChannelParameter, multi int, legacyVolume bool) ([]*TransferInstruction, error) {

	componentsToMove, instructionsToUse := makeComponentSets(inssIn, multi)

	orientation := wtype.LHVChannel
	independent := false
	if channelprms != nil {
		orientation = channelprms.Orientation
		independent = channelprms.Independent
	}

	ret := make([]*TransferInstruction, 0)
	for i := 0; i < len(componentsToMove); i++ {
		parallelTransfers, err := robot.GetComponents(componentsToMove[i], orientation, multi, independent, legacyVolume)
		if err != nil {
			return nil, err
		}

		for _, t := range parallelTransfers {
			if transfers, err := makeTransfers(t, componentsToMove[i], robot, instructionsToUse[i]); err != nil {
				return nil, err
			} else {
				ret = append(ret, transfers...)
			}
		}

	}

	return ret, nil
}

func makeTransfers(parallelTransfer ParallelTransfer, cmps []*wtype.Liquid, robot *LHProperties, inssIn []*wtype.LHInstruction) ([]*TransferInstruction, error) {
	fromPlateIDs := parallelTransfer.PlateIDs
	fromWells := parallelTransfer.WellCoords
	vols := parallelTransfer.Vols

	insOut := make([]*TransferInstruction, 0, 1)
	// mt counts up the arrays got by GetComponents
	// each array refers to a transfer
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
	cnames := make([]string, len(cmps))   //        component names
	policies := make([]wtype.LHPolicy, len(cmps))
	// ci counts up cmps

	for ci := 0; ci < len(cmps); ci++ {
		if len(fromPlateIDs) <= ci || fromPlateIDs[ci] == "" || fromPlateIDs[ci] == "<No-ID>" {
			continue
		}

		wh[ci] = cmps[ci].TypeName()

		// source plate position

		ppf, ok := robot.PlateIDLookup[fromPlateIDs[ci]]

		if !ok {
			return insOut, wtype.LHError(wtype.LH_ERR_DIRE, "Planning inconsistency: input plate ID not found on robot - please report this error to the authors")
		}

		pf[ci] = ppf

		// destination plate position

		ppt, ok := robot.PlateIDLookup[inssIn[ci].PlateID]

		if !ok {
			return insOut, wtype.LHError(wtype.LH_ERR_DIRE, "Planning inconsistency: destination plate ID not found on robot - please report this error to the authors")
		}

		pt[ci] = ppt

		// source well

		wf[ci] = fromWells[ci]

		wt[ci] = inssIn[ci].Welladdress

		// source plate type

		srcPlate, ok := robot.Plates[ppf]

		if !ok {
			return insOut, wtype.LHError(wtype.LH_ERR_DIRE, "Planning inconsistency: input plate ID not found on robot (#2) - please report this error to the authors")
		}

		ptf[ci] = srcPlate.Type

		// destination plate type

		dstPlate, ok := robot.Plates[ppt]

		if !ok {
			return insOut, wtype.LHError(wtype.LH_ERR_DIRE, "Planning inconsistency: destination plate ID not found on robot - please report this error to the authors")
		}

		ptt[ci] = dstPlate.Type

		// volume being moved

		va[ci] = vols[ci]

		// source well volume

		wellFrom, ok := srcPlate.WellAtString(wf[ci])

		if !ok {
			return insOut, wtype.LHError(wtype.LH_ERR_DIRE, "Planning inconsistency: source well not found on source plate - plate report this error to the authors")
		}

		vf[ci] = wellFrom.CurrentVolume()

		// dest well volume

		wellTo, ok := dstPlate.Wellcoords[wt[ci]]

		if !ok {
			return insOut, wtype.LHError(wtype.LH_ERR_DIRE, "Planning inconsistency: dest well not found on dest plate - please report this error to the authors")
		}

		vt[ci] = wellTo.CurrentVolume()

		// source plate dimensions

		pfwx[ci] = srcPlate.WellsX()
		pfwy[ci] = srcPlate.WellsY()

		// dest plate dimensions

		ptwx[ci] = dstPlate.WellsX()
		ptwy[ci] = dstPlate.WellsY()

		cnames[ci] = wellFrom.WContents.CName

		cmpFrom, err := wellFrom.RemoveVolume(va[ci])
		if err != nil {
			return insOut, wtype.LHErrorf(wtype.LH_ERR_DIRE, "Planning inconsistency: %s - please report this error to the authors", err.Error())
		}

		// silently remove the carry
		wellFrom.RemoveCarry(robot.CarryVolume())

		err = wellTo.AddComponent(cmpFrom)
		if err != nil {
			return insOut, wtype.LHErrorf(wtype.LH_ERR_VOL, "Planning inconsistency : %s", err.Error())
		}

		// make sure the cmp loc is set

		wellTo.WContents.Loc = wtype.IDOf(wellTo.Plate) + ":" + wellTo.Crds.FormatA1()

		// make sure the wellTo gets the right ID (ultimately)
		cmpFrom.ReplaceDaughterID(wellTo.WContents.ID, inssIn[ci].Outputs[0].ID)
		wellTo.WContents.ID = inssIn[ci].Outputs[0].ID
		wellTo.WContents.DeclareInstance()
		//fmt.Println("ADDED :", cmpFrom.CName, " ", cmpFrom.Vol, " TO ", dstPlate.ID, " ", wt[ci])
	}

	//}

	tfr := NewTransferInstruction(wh, pf, pt, wf, wt, ptf, ptt, va, vf, vt, pfwx, pfwy, ptwx, ptwy, cnames, policies)
	insOut = append(insOut, tfr)

	return insOut, nil
}
