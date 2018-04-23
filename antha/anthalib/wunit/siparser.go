// /anthalib/wunit/siparser.go: Part of the Antha language
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

package wunit

type NodeType uint8

const (
	TopNode = iota
	UnitPlusPrefixNode
	LeafNode
)

type PNode struct {
	Name     string
	Type     NodeType
	Up       *PNode
	Children []*PNode
	Value    interface{}
}

type SIPrefixedUnit struct {
	Stack   []*PNode
	Text    string
	TreeTop *PNode
	CurNode *PNode
}

func (p *SIPrefixedUnit) Init(text []byte) {
	p.Text = string(text)
	node := NewNode("Top", UnitPlusPrefixNode, 2)
	p.TreeTop = node
	p.CurNode = node
}

// couple of utility functions

func (p *SIPrefixedUnit) AddNodeToStack(node *PNode) {
	p.Stack = append(p.Stack, node)

	/*
		for _, n := range p.Stack {
			fmt.Print(n.Name, ":")
			logger.Debug(fmt.Sprintln())
		}
	*/
}

func (p *SIPrefixedUnit) AddStackToNode(node *PNode) {
	node.Children = p.Stack

	for _, n := range node.Children {
		n.Up = node
	}
}

func (p *SIPrefixedUnit) PopStack() (node *PNode) {
	lastItem := len(p.Stack) - 1

	if lastItem == -1 {
		return nil
	}
	node = p.Stack[lastItem]

	if lastItem == 0 {
		p.NewStack()
	} else {
		p.Stack = p.Stack[:lastItem]
	}

	return node
}

func (p *SIPrefixedUnit) PopStackAndAddTo(node *PNode) {
	n := p.PopStack()
	if n != nil {
		node.AddChild(n)
	}
}

func (pn *PNode) AddChild(child *PNode) {
	if child != nil {
		child.Up = pn
		pn.Children = append(pn.Children, child)
	}
}

func (p *SIPrefixedUnit) NewStack() {
	p.Stack = make([]*PNode, 0, 4)
}

func (p *SIPrefixedUnit) AddNodeToCurrent(node *PNode) {
	p.CurNode.Children = append(p.CurNode.Children, node)
}

func NewNode(name string, typ NodeType, cap uint8) *PNode {
	var children []*PNode

	if cap != 0 {
		children = make([]*PNode, 0, cap)
	}
	node := PNode{name, typ, nil, children, nil}
	return &node
}

// Functions for building the tree

func (p *SIPrefixedUnit) AddUnit(s string) {
	//	logger.Debug(fmt.Sprintln("Adding Unit", s))
	node := NewNode("Unit", LeafNode, 0)
	node.Value = s
	p.AddNodeToStack(node)
}

func (p *SIPrefixedUnit) AddUnitPrefix(s string) {
	//	logger.Debug(fmt.Sprintln("Adding unit prefix", s))
	node := NewNode("UnitPrefix", LeafNode, 0)
	node.Value = s
	p.AddNodeToStack(node)
}

func (p *SIPrefixedUnit) AddUnitPlusPrefixNode() {
	node := NewNode("UnitPlusPrefix", UnitPlusPrefixNode, 2)
	pref := p.PopStack()
	p.PopStackAndAddTo(node)
	node.AddChild(pref)
	p.TreeTop = node
	p.CurNode = node
}
