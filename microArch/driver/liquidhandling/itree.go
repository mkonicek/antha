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
	"context"
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
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

	for {
		if ch == nil {
			break
		}

		if ch.Values[0].Type == wtype.LHIPRM {
			root.AddChild(NewMessageInstruction(ch.Values[0]))
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

			root.AddChild(NewSplitBlockInstruction(ch.Values))
		} else {
			// otherwise...
			// make a transfer block instruction out of the incoming instructions
			// -- essentially each node of the topological graph is passed wholesale
			// into the instruction generator to be teased apart as appropriate

			root.AddChild(NewTransferBlockInstruction(ch.Values))
		}
		ch = ch.Child
	}

	// last thing is always to finalize the liquidhandler
	root.AddChild(NewFinalizeInstruction())

	return root, nil
}

// AddChild add a child to this node of the ITree
func (tree *ITree) AddChild(ins RobotInstruction) {
	tree.children = append(tree.children, NewITree(ins))
}

// Generate generate the tree of instructions, returning an error on failure
func (tree *ITree) Generate(ctx context.Context, lhpr *wtype.LHPolicyRuleSet, lhpm *LHProperties) error {

	// the root node (instruction == nil) already has children set
	if tree.instruction != nil {
		if children, err := tree.instruction.Generate(ctx, lhpr, lhpm); err != nil {
			return err
		} else {
			for _, child := range children {
				tree.AddChild(child)
			}
		}
	}

	// call generate on the children recursively
	for _, child := range tree.children {
		if err := child.Generate(ctx, lhpr, lhpm); err != nil {
			return err
		}
	}

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
func (tree *ITree) Leaves() ([]TerminalRobotInstruction, error) {
	return tree.addLeaves(make([]TerminalRobotInstruction, 0, tree.NumLeaves()))
}

// addLeaves add the leaves to the accumulator
func (tree *ITree) addLeaves(acc []TerminalRobotInstruction) ([]TerminalRobotInstruction, error) {
	if len(tree.children) == 0 {
		// i am leaf (on the wind)
		// ignore instructions which aren't terminal instructions, these are probably things like split which don't
		// actually generate terminal instructions at all
		if tri, ok := tree.instruction.(TerminalRobotInstruction); ok {
			return append(acc, tri), nil
		}
	} else {
		for _, child := range tree.children {
			if nac, err := child.addLeaves(acc); err != nil {
				return nac, err
			} else {
				acc = nac
			}
		}
	}
	return acc, nil
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
