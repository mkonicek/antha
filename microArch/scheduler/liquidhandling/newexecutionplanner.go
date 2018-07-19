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
	"context"
	"time"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

// robot here should be a copy... this routine will be destructive of state
func ImprovedExecutionPlanner(ctx context.Context, request *LHRequest, robot *liquidhandling.LHProperties) (*LHRequest, error) {
	rbtcpy := robot.Dup()

	// get timer to assess evaporation etc.

	timer := robot.GetTimer()
	// 1 -- generate high level instructions

	// aggregation now works by lumping together stuff that makes the same components
	// until it finds something different

	agg := make([][]int, 0, 1)
	curragg := make([]int, 0, 1)

	instrx := make([]liquidhandling.RobotInstruction, 0, len(request.LHInstructions))
	evaps := make([]wtype.VolumeCorrection, 0, 10)

	for ix, insID := range request.OutputOrder {
		//	request.InstructionSet.Add(ConvertInstruction(request.LHInstructions[insID], robot))

		ins := request.LHInstructions[insID]

		ris := liquidhandling.NewRobotInstructionSet(nil)

		if ins.Type == wtype.LHIPRM {
			// prompt
			prm := liquidhandling.NewMessageInstruction(ins)
			ris.Add(prm)
			instrx = append(instrx, prm)
			robot.UpdateComponentIDs(ins.PassThrough) // prompting changes IDs
		} else {
			transIns, err := ConvertInstruction(ins, robot, request.CarryVolume, request.UseLegacyVolume())

			if err != nil {
				return request, err
			}
			ris.Add(transIns)
			instrx = append(instrx, transIns)
		}

		var cmp string // partition key

		if request.LHInstructions[insID].Type == wtype.LHIMIX {
			cmp = request.LHInstructions[insID].NamesOfComponentsMoving()
		} else if request.LHInstructions[insID].Type == wtype.LHIPRM {
			cmp = request.LHInstructions[insID].Message
		}

		if canaggregate(curragg, cmp, request.OutputOrder, request.LHInstructions) {
			// true if either curragg empty or cmp is same
			curragg = append(curragg, ix)
		} else {
			agg = append(agg, curragg)
			curragg = make([]int, 0, 1)
			curragg = append(curragg, ix)
		}

		if request.Options.ModelEvaporation {
			// we should be able to model evaporation here
			instrx, _ := ris.Generate(ctx, request.Policies(), rbtcpy)

			if timer != nil {
				var totaltime time.Duration
				for _, instr := range instrx {
					totaltime += timer.TimeFor(instr)
					// PROMPTS really screw this up... this should generate
					// a warning. We will assume zero time
				}

				// evaporate stuff

				myevap := robot.Evaporate(totaltime)
				evaps = append(evaps, myevap...)
			}
		}
	}
	agg = append(agg, curragg)

	// 2 -- see if any of the above can be aggregated, if so we merge them

	mergedInstructions := merge_instructions(instrx, agg)

	// 3 -- add them to the instruction set

	for _, ins := range mergedInstructions {
		request.InstructionSet.Add(ins)
	}

	// 4 -- make the low-level instructions

	inx, err := request.InstructionSet.Generate(ctx, request.Policies(), robot)

	if err != nil {
		return nil, err
	}

	finalInstrx := make([]liquidhandling.TerminalRobotInstruction, len(inx))
	for i := 0; i < len(inx); i++ {
		finalInstrx[i] = inx[i].(liquidhandling.TerminalRobotInstruction)
	}
	request.Instructions = finalInstrx

	request.Evaps = evaps

	return request, nil
}

func canaggregate(agg []int, cmp string, outorder []string, cmps map[string]*wtype.LHInstruction) bool {

	if len(agg) == 0 {
		return true
	}

	if !singleInstructionType(agg, outorder, cmps) {
		return false
	}

	if cmps[outorder[agg[0]]].Type == wtype.LHIPRM {
		return cmps[outorder[agg[0]]].Message == cmp
	} else {
		return cmps[outorder[agg[0]]].NamesOfComponentsMoving() == cmp
	}
}

func singleInstructionType(agg []int, outorder []string, cmps map[string]*wtype.LHInstruction) bool {
	instype := cmps[outorder[agg[0]]].Type
	for i := 1; i < len(agg); i++ {
		if cmps[outorder[agg[i]]].Type != instype {
			return false
		}
	}

	return true
}
