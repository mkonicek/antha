// anthalib//liquidhandling/newexecutionplanner.go: Part of the Antha language
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
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"time"
)

// robot here should be a copy... this routine will be destructive of state
func ImprovedExecutionPlanner(request *LHRequest, robot *liquidhandling.LHProperties) (*LHRequest, error) {
	rbtcpy := robot.Dup()

	// get timer to assess evaporation etc.

	timer := robot.GetTimer()
	// 1 -- generate high level instructions

	// aggregation now works by lumping together stuff that makes the same components
	// until it finds something different

	agg := make([][]int, 0, 1)
	curragg := make([]int, 0, 1)

	transfers := make([]liquidhandling.RobotInstruction, 0, len(request.LHInstructions))
	evaps := make([]wtype.VolumeCorrection, 0, 10)
	for ix, insID := range request.Output_order {
		//	request.InstructionSet.Add(ConvertInstruction(request.LHInstructions[insID], robot))

		ris := liquidhandling.NewRobotInstructionSet(nil)

		transIns, err := ConvertInstruction(request.LHInstructions[insID], robot, request.CarryVolume)

		if err != nil {
			return request, err
		}

		ris.Add(transIns)

		transfers = append(transfers, transIns)
		cmp := fmt.Sprintf("%s_%s", request.LHInstructions[insID].ComponentsMoving(), request.LHInstructions[insID].Generation())

		/*
			ar, ok := agg[cmp]
			if !ok {
				ar = make([]int, 0, 1)
			}

			ar = append(ar, ix)
			agg[cmp] = ar

		*/

		if canaggregate(curragg, cmp, request.Output_order, request.LHInstructions) {
			// true if either curragg empty or cmp is same
			curragg = append(curragg, ix)
		} else {
			agg = append(agg, curragg)
			curragg = make([]int, 0, 1)
			curragg = append(curragg, ix)

		}

		if request.Options.ModelEvaporation {
			// we should be able to model evaporation here

			instrx, _ := ris.Generate(request.Policies, rbtcpy)

			if timer != nil {
				var totaltime time.Duration
				for _, instr := range instrx {
					totaltime += timer.TimeFor(instr)
				}

				// evaporate stuff

				myevap := robot.Evaporate(totaltime)
				evaps = append(evaps, myevap...)
			}
		}
	}
	agg = append(agg, curragg)

	// 2 -- see if any of the above can be aggregated, if so we merge them

	transfers = merge_transfers(transfers, agg)

	// 3 -- add them to the instruction set

	for _, tfr := range transfers {
		request.InstructionSet.Add(tfr)
	}

	// 4 -- make the low-level instructions

	inx, err := request.InstructionSet.Generate(request.Policies, robot)

	if err != nil {
		return nil, err
	}

	instrx := make([]liquidhandling.TerminalRobotInstruction, len(inx))
	for i := 0; i < len(inx); i++ {
		//fmt.Println(liquidhandling.InsToString(inx[i]))
		instrx[i] = inx[i].(liquidhandling.TerminalRobotInstruction)
	}
	request.Instructions = instrx

	request.Evaps = evaps

	return request, nil
}

func canaggregate(agg []int, cmp string, outorder []string, cmps map[string]*wtype.LHInstruction) bool {

	if len(agg) == 0 {
		return true
	}

	return cmps[outorder[agg[0]]].ComponentsMoving() == cmp
}
