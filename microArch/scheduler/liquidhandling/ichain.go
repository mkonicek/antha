package liquidhandling

import (
	"fmt"

	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/graph"
)

type IChain struct {
	Parent *IChain
	Child  *IChain
	Values []*wtype.LHInstruction
	Depth  int
}

func NewIChain(parent *IChain) *IChain {
	var it IChain
	it.Parent = parent
	it.Values = make([]*wtype.LHInstruction, 0, 1)
	if parent != nil {
		it.Depth = parent.Depth + 1
	}
	return &it
}

// depth from here
func (it *IChain) Height() int {
	if it == nil {
		return 0
	}

	return it.Child.Height() + 1
}

func (it *IChain) PruneOut(Remove map[string]bool) *IChain {
	if it == nil || len(Remove) == 0 || len(it.Values) == 0 {
		return it
	}

	it.Child = it.Child.PruneOut(Remove)

	newValues := make([]*wtype.LHInstruction, 0, len(it.Values))

	for _, v := range it.Values {
		if Remove[v.ID] {
			continue
		}
		newValues = append(newValues, v)
		delete(Remove, v.ID)
	}

	// if we've removed a whole layer, get rid of it

	if len(newValues) == 0 {

		if it.Child != nil {
			it.Child.Parent = it.Parent
		}

		if it.Parent != nil {
			it.Parent.Child = it.Child
		}

		return it.Child

	} else {
		it.Values = newValues
		return it
	}

}

func (it *IChain) AsList(ica []*IChain) []*IChain {
	if it == nil {
		return ica
	}

	ica = append(ica, it)

	return it.Child.AsList(ica)
}

func (it *IChain) Reverse() {
	if it.Child != nil {
		it.Child.Reverse()
	}
	// swap parent and child
	p := it.Parent
	it.Parent = it.Child
	it.Child = p
}

func (it *IChain) ValueIDs() []string {
	r := make([]string, 0, 1)

	for _, v := range it.Values {
		r = append(r, v.ID)
	}
	return r
}

func (it *IChain) Add(ins *wtype.LHInstruction) {
	if it.Depth < ins.Generation() {
		it.GetChild().Add(ins)
	} else {
		it.Values = append(it.Values, ins)
	}
}

func (it *IChain) GetChild() *IChain {
	if it.Child == nil {
		it.Child = NewIChain(it)
	}
	return it.Child
}

func (it *IChain) Print() {
	fmt.Println("****")
	fmt.Println("\tPARENT NIL: ", it.Parent == nil)
	if len(it.Values) > 0 {
		for j := 0; j < len(it.Values); j++ {
			if it.Values[j].Type == wtype.LHIMIX {
				fmt.Printf("MIX    %2d: %s \n", j, it.Values[j].ID)
				for i := 0; i < len(it.Values[j].Components); i++ {
					fmt.Print(" ", it.Values[j].Components[i].ID, ":", it.Values[j].Components[i].FullyQualifiedName(), "@", it.Values[j].Components[i].Volume().ToString(), " \n")
				}
				fmt.Println(":", it.Values[j].Results[0].ID, ":", it.Values[j].Platetype, " ", it.Values[j].PlateName, " ", it.Values[j].Welladdress)
				fmt.Printf("-- ")
			} else if it.Values[j].Type == wtype.LHIPRM {
				fmt.Println("PROMPT ", it.Values[j].Message, "-- ")
				for in, out := range it.Values[j].PassThrough {
					fmt.Println(in, ":::", out.ID, " --")
				}
			} else if it.Values[j].Type == wtype.LHISPL {
				fmt.Printf("SPLIT %2d: %s ", j, it.Values[j].ID)
				fmt.Println(" ", it.Values[j].Components[0].ID, ":", it.Values[j].Components[0].FullyQualifiedName(), " : ", it.Values[j].PlateName, " ", it.Values[j].Welladdress)
				fmt.Println(" MOVE:", it.Values[j].Results[0].ID, ":", it.Values[j].Results[0].FullyQualifiedName(), "@", it.Values[j].Results[0].Volume().ToString())
				fmt.Println(" STAY:", it.Values[j].Results[1].ID, ":", it.Values[j].Results[1].FullyQualifiedName(), "@", it.Values[j].Results[1].Volume().ToString())
				fmt.Printf("-- \n")
			} else {
				fmt.Println("WTF?   ", wtype.InsType(it.Values[j].Type), "-- ")
			}
		}
		fmt.Println()
	}
	if it.Child != nil {
		it.Child.Print()
	}
}

func (it *IChain) InputIDs() string {
	s := ""

	for _, ins := range it.Values {
		for _, c := range ins.Components {
			s += c.ID + "   "
		}
		s += ","
	}

	return s
}

func (it *IChain) ProductIDs() string {
	s := ""

	for _, ins := range it.Values {
		s += strings.Join(ins.ProductIDs(), " ") + "   "
	}
	return s
}

//Flatten get the instruction IDs in order
func (it *IChain) Flatten() []string {
	var ret []string

	if it == nil {
		return ret
	}

	for _, v := range it.Values {
		ret = append(ret, v.ID)
	}

	ret = append(ret, it.Child.Flatten()...)

	return ret
}

//GetOrderedLHInstructions get the instructions in order
func (it *IChain) GetOrderedLHInstructions() []*wtype.LHInstruction {
	var ret []*wtype.LHInstruction

	if it == nil {
		return ret
	}

	ret = append(ret, it.Values...)

	ret = append(ret, it.Child.GetOrderedLHInstructions()...)

	return ret
}

func (it *IChain) SplitMixedNodes() {
	if it.hasMixAndSplitOnly() {
		it.splitMixedNode()
	}

	// stop if we reach the end
	if it.Child == nil {
		return
	}

	// carry on
	it.Child.SplitMixedNodes()
}

func (it *IChain) splitMixedNode() {
	// put mixes first, then splits

	mixValues := make([]*wtype.LHInstruction, 0, len(it.Values))
	splitValues := make([]*wtype.LHInstruction, 0, len(it.Values))

	for _, v := range it.Values {
		if v.Type == wtype.LHIMIX {
			mixValues = append(mixValues, v)
		} else if v.Type == wtype.LHISPL {
			splitValues = append(splitValues, v)
		} else {
			panic("Wrong instruction type passed through to instruction chain split")
		}
	}

	// it == Mix level
	it.Values = mixValues
	// ch == Split level
	ch := NewIChain(it)
	ch.Values = splitValues
	ch.Child = it.Child

	if ch.Child != nil {
		ch.Child.Parent = ch
	}
	it.Child = ch
}

type icGraph struct {
	edges map[graph.Node][]graph.Node
	nodes []graph.Node
}

func (g icGraph) NumNodes() int {
	return len(g.nodes)
}

func (g icGraph) Node(i int) graph.Node {
	return g.nodes[i]
}

func (g icGraph) NumOuts(n graph.Node) int {
	a, ok := g.edges[n]

	if ok {
		return len(a)
	} else {
		return 0
	}
}

func (g icGraph) Out(n graph.Node, i int) graph.Node {
	return g.edges[n][i]
}

//AsGraph returns the chain in graph form, unidirectional only
func (ic *IChain) AsGraph() graph.Graph {
	edges := make(map[graph.Node][]graph.Node)
	nodes := make([]graph.Node, 0, 1)
	n := 1
	for curr := ic; curr != nil; curr = curr.Child {
		nodes = append(nodes, curr)
		if curr.Child != nil {
			edges[graph.Node(curr)] = []graph.Node{graph.Node(curr.Child)}
		}
		n += 1
	}

	return icGraph{
		edges: edges,
		nodes: nodes,
	}
}

func (ic *IChain) hasMixAndSplitOnly() bool {
	/// true iff we have exactly two types of node: split and mix
	insTypes := ic.getInstructionTypes()

	return len(insTypes) == 2 && insTypes[wtype.InsNames[wtype.LHIMIX]] && insTypes[wtype.InsNames[wtype.LHISPL]]
}

func (self *IChain) getInstructionTypes() map[string]bool {
	types := make(map[string]bool, len(self.Values))
	for _, v := range self.Values {
		types[v.InsType()] = true
	}

	return types
}

//assertInstructionsSeparate check that there's only one type of instruction
//in each link of the chain
func (self *IChain) assertInstructionsSeparate() error {
	if self == nil {
		return nil
	}

	types := self.getInstructionTypes()

	if len(types) != 1 {
		return fmt.Errorf("Only one instruction type per stage is allowed, found %v at stage %d", len(types), self.Depth)
	}

	return self.Child.assertInstructionsSeparate()
}
