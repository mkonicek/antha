// /anthalib/driver/liquidhandling/itree.go: Part of the Antha language
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
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/logger/levlog"
)

// ITree the Instruction Tree - takes a high level liquid handling instruction
// and calls Generate() recursively to produce a tree of instructions with the
// lowest level "Terminal" robot instructions at the bottom
type ITree struct {
	instruction RobotInstruction
	children    []*ITree
}

// NewITree initialize a new tree with the given parent instruction
func NewITree(p RobotInstruction) *ITree {
	return &ITree{instruction: p}
}

// Children get the children of this node of the tree
func (tree *ITree) Children() []*ITree {
	return tree.children
}

// Instruction get the instruction associated with this node of the tree
func (tree *ITree) Instruction() RobotInstruction {
	return tree.instruction
}

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

// NewITreeRoot create the root layer of the instruction tree. The root layer
// has a nil instruction with children which correspond to each layer of the
// IChain
func NewITreeRoot(ch *wtype.IChain) (*ITree, error) {

	root := NewITree(nil)

	// first thing is always to initialize the liquidhandler
	root.AddChild(NewInitializeInstruction())
	counter := 0
	for {
		if ch == nil {
			levlog.Debug(" ---- Breaking out after ", counter, " iterations")
			break
		}

		if ch.Values[0].Type == wtype.LHIPRM {
			levlog.Debug(" --> Adding PROMPT ", ch.Values[0].Message, " Counter : ", counter)
			root.AddChild(NewMessageInstruction(ch.Values[0]))
		} else if hasSplit(ch.Values) {
			levlog.Debug(" --> Adding Split ", ch.Values[0].Message, "  Counter : ", counter)
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

			root.AddChild(NewSplitBlockInstruction(ch.Values))
		} else {
			// otherwise...
			// make a transfer block instruction out of the incoming instructions
			// -- essentially each node of the topological graph is passed wholesale
			// into the instruction generator to be teased apart as appropriate
			levlog.Debug(" Adding TransferBlockInstruction -- len(ch.Values) = ", len(ch.Values), " counter= ", counter)
			root.AddChild(NewTransferBlockInstruction(ch.Values))
		}
		ch = ch.Child
		counter += 1
	}

	// last thing is always to finalize the liquidhandler
	root.AddChild(NewFinalizeInstruction())

	return root, nil
}

// AddChild add a child to this node of the ITree
func (tree *ITree) AddChild(ins RobotInstruction) {
	levlog.Debug("  --  -- appending NewITree ", ins, " to tree ")
	tree.children = append(tree.children, NewITree(ins))
}

// Build the tree of instructions, starting from the root, discarding any existing children
// returns the final state and does not alter the initial state
func (tree *ITree) Build(labEffects *effects.LaboratoryEffects, lhpr *wtype.LHPolicyRuleSet, initial *LHProperties) (*LHProperties, error) {
	final := initial.DupKeepIDs(labEffects.IDGenerator)
	err := tree.addChildren(labEffects, lhpr, final, 0)
	return final, err
}

// addChildren recursively Generate() the children of this instruction, adding them to the tree.
// the effects of the generated instructions are applied to lhpm
func (tree *ITree) addChildren(labEffects *effects.LaboratoryEffects, lhpr *wtype.LHPolicyRuleSet, lhpm *LHProperties, depth int) error {

	// the root node (instruction == nil) already has children set
	if tree.instruction != nil {
		levlog.Debug("===== Generating instruction ====== ")
		if children, err := tree.instruction.Generate(labEffects, lhpr, lhpm); err != nil {
			return err
		} else {
			levlog.Debug(" ---- adding ", len(children), " children to tree ---- ", len(children))
			for _, child := range children {
				tree.AddChild(child)
			}
		}
	}

	// call generate on the children recursively
	levlog.Debug(" # Children: ", len(tree.children))
	for ii, child := range tree.children {
		levlog.Debug(" Depth ", depth, " -- Loop ", ii)
		if err := child.addChildren(labEffects, lhpr, lhpm, depth+1); err != nil {
			return err
		}
	}

	levlog.Debug(" ")
	return nil
}

// NumLeaves returns the number of leaves at the bottom of the tree
func (tree *ITree) NumLeaves() int {
	if len(tree.children) == 0 {
		return 1
	} else {
		ret := 0
		for _, child := range tree.children {
			ret += child.NumLeaves()
		}
		return ret
	}
}

// Leaves returns the leaves of the tree - i.e. the TerminalRobotInstructions
func (tree *ITree) Leaves() []TerminalRobotInstruction {
	leaves := tree.Refine()
	ret := make([]TerminalRobotInstruction, 0, len(leaves))
	for _, leaf := range leaves {
		if tri, ok := leaf.instruction.(TerminalRobotInstruction); ok {
			ret = append(ret, tri)
		}
	}
	return ret
}

// Refine returns a slice containing more precise instructions. Returns the lowest possible descendents of this node,
// stopping either at leaves or if the instruction type is given in the argument
// e.g. for a tree with instruction types A,B,C,D,E of the form
//       A
//     /   \
//    B     C
//    |    / \
//    D   D'  E
// then:
//   - A.Refine(C) returns [D, C]
//   - A.Refine(B) returns [B, D', E]
//   - A.Refine() returns [D, D', E], the leaves of the tree
func (tree *ITree) Refine(types ...*InstructionType) []*ITree {
	tmap := make(map[*InstructionType]bool, len(types))
	for _, t := range types {
		tmap[t] = true
	}

	ret := make([]*ITree, 0, tree.NumLeaves())
	return tree.refine(tmap, ret)
}

func (tree *ITree) refine(tmap map[*InstructionType]bool, acc []*ITree) []*ITree {
	if tree == nil {
		return acc
	} else if (tree.instruction != nil && tmap[tree.instruction.Type()]) || len(tree.children) == 0 {
		// the type is in tmap or we've reached the bottom of the tree
		return append(acc, tree)
	} else {
		for _, child := range tree.children {
			acc = child.refine(tmap, acc)
		}
		return acc
	}
}

// String return a multi-line representation of the ITree showing instructions types
func (tree *ITree) String() string {
	return tree.toString(0)
}

// toString return a multi-line representation of the tree below this node, starting
// at the given indentation level
func (tree *ITree) toString(level int) string {

	name := ""

	if tree.instruction != nil {
		name = tree.instruction.Type().Name
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
	for _, ins := range tree.children {
		s += ins.toString(level + 1)
	}
	for i := 0; i < level; i++ {
		s += fmt.Sprintf("\t")
	}
	s += "}\n"
	return s
}
