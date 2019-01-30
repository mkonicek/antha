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
)

func allSplits(inss []*wtype.LHInstruction) bool {
	for _, ins := range inss {
		if ins.Type != wtype.LHISPL {
			return false
		}
	}
	return true
}

func hasSplit(inss []*wtype.LHInstruction) bool {
	for _, ins := range inss {
		if ins.Type == wtype.LHISPL {
			return true
		}
	}
	return false
}

// ExecutionPlanner create the instruction tree
func ExecutionPlanner3(ch *wtype.IChain) (*liquidhandling.RobotInstructionSet, error) {

	root := liquidhandling.NewRobotInstructionSet(nil)

	// first thing is always to initialize the liquidhandler
	root.AddChild(liquidhandling.NewInitializeInstruction())

	for {
		if ch == nil {
			break
		}

		if ch.Values[0].Type == wtype.LHIPRM {
			prm := liquidhandling.NewMessageInstruction(ch.Values[0])
			root.AddChild(prm)
		} else if hasSplit(ch.Values) {
			if !allSplits(ch.Values) {
				insTypes := func(inss []*wtype.LHInstruction) string {
					s := ""
					for _, ins := range inss {
						s += ins.InsType() + " "
					}

					return s
				}
				return nil, fmt.Errorf("Internal error: Failure in instruction sorting - got types %s in layer starting with split", insTypes(ch.Values))
			}

			splitBlock := liquidhandling.NewSplitBlockInstruction(ch.Values)
			root.AddChild(splitBlock)
		} else {
			// otherwise...
			// make a transfer block instruction out of the incoming instructions
			// -- essentially each node of the topological graph is passed wholesale
			// into the instruction generator to be teased apart as appropriate

			tfb := liquidhandling.NewTransferBlockInstruction(ch.Values)

			root.AddChild(tfb)
		}
		ch = ch.Child
	}

	// last thing is always to finalize the liquidhandler
	root.AddChild(liquidhandling.NewFinalizeInstruction())

	return root, nil
}
