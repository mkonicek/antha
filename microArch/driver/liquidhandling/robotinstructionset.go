// /anthalib/driver/liquidhandling/robotinstructionset.go: Part of the Antha language
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
	"github.com/antha-lang/antha/laboratory"
)

type RobotInstructionSet struct {
	parent       RobotInstruction
	instructions []*RobotInstructionSet
}

func NewRobotInstructionSet(p RobotInstruction) *RobotInstructionSet {
	var ret RobotInstructionSet
	ret.instructions = make([]*RobotInstructionSet, 0)
	ret.parent = p
	return &ret
}

func (ri *RobotInstructionSet) Add(ins RobotInstruction) {
	ris := NewRobotInstructionSet(ins)
	ri.instructions = append(ri.instructions, ris)
}

func (ri *RobotInstructionSet) Generate(lab *laboratory.Laboratory, lhpr *wtype.LHPolicyRuleSet, lhpm *LHProperties) ([]RobotInstruction, error) {
	ret := make([]RobotInstruction, 0, 1)

	if ri.parent != nil {
		arr, err := ri.parent.Generate(lab, lhpr, lhpm)

		if err != nil {
			return ret, err
		}

		// if the parent doesn't generate anything then it is our return - bottom out here
		// assuming it's a Terminal
		if len(arr) == 0 {
			_, ok := ri.parent.(TerminalRobotInstruction)
			if ok {
				ret = append(ret, ri.parent)
				return ret, nil
			}
		} else {
			for _, ins := range arr {
				ri.Add(ins)
			}
		}
	}

	for _, ins := range ri.instructions {
		arr, err := ins.Generate(lab, lhpr, lhpm)

		if err != nil {
			return arr, err
		}
		ret = append(ret, arr...)
	}

	if ri.parent == nil {
		// add the initialize and finalize instructions
		ini := NewInitializeInstruction()
		newret := make([]RobotInstruction, 0, len(ret)+2)
		newret = append(newret, ini)
		newret = append(newret, ret...)
		fin := NewFinalizeInstruction()
		newret = append(newret, fin)
		ret = newret
	}

	// might need to do this instead of current version
	/*
		else if ri.parent.Type == TFR {
			// update the vols
			prms.Evaporate()
		}
	*/

	return ret, nil
}

func (ri *RobotInstructionSet) ToString(level int) string {

	name := ""

	if ri.parent != nil {
		name = ri.parent.Type().Name
	}
	s := ""
	for i := 0; i < level-1; i++ {
		s += fmt.Sprintf("\t")
	}
	s += fmt.Sprintf("%s\n", name)
	for i := 0; i < level; i++ {
		s += fmt.Sprintf("\t")
	}
	s += fmt.Sprintf("{\n")
	for _, ins := range ri.instructions {
		s += ins.ToString(level + 1)
	}
	for i := 0; i < level; i++ {
		s += fmt.Sprintf("\t")
	}
	s += "}\n"
	return s
}
