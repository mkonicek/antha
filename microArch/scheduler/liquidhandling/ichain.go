package liquidhandling

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/graph"
	"strings"
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
				fmt.Printf("MIX    %2d: %s ", j, it.Values[j].ID)
				for i := 0; i < len(it.Values[j].Components); i++ {
					fmt.Print(" ", it.Values[j].Components[i].ID, ":", it.Values[j].Components[i].FullyQualifiedName(), "@", it.Values[j].Components[i].Volume().Summary(), " ")
				}
				fmt.Print(":", it.Values[j].Results[0].ID, ":", it.Values[j].Platetype, " ", it.Values[j].PlateName, " ", it.Values[j].Welladdress)
				fmt.Printf("-- ")
			} else if it.Values[j].Type == wtype.LHIPRM {
				fmt.Print("PROMPT ", it.Values[j].Message, "-- ")
				for in, out := range it.Values[j].PassThrough {
					fmt.Print(in, ":::", out.ID, " --")
				}
			} else if it.Values[j].Type == wtype.LHISPL {
				fmt.Printf("SPLIT %2d: %s ", j, it.Values[j].ID)
				fmt.Print(" ", it.Values[j].Components[0].ID, ":", it.Values[j].Components[0].FullyQualifiedName(), " : ", it.Values[j].PlateName, " ", it.Values[j].Welladdress, " ")
				fmt.Print(" MOVE:", it.Values[j].Results[0].ID, ":", it.Values[j].Results[0].FullyQualifiedName(), "@", it.Values[j].Results[0].Volume().Summary())
				fmt.Print(" STAY:", it.Values[j].Results[1].ID, ":", it.Values[j].Results[1].FullyQualifiedName(), "@", it.Values[j].Results[1].Volume().Summary())
				fmt.Printf("-- ")
			} else {
				fmt.Print("WTF?   ", wtype.InsType(it.Values[j].Type), "-- ")
			}
		}

		fmt.Println("End of Instruction")
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

func (it *IChain) SplitMixedNodes() {
	if nodesMixedOK(it.Values) {
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
	ch.Child.Parent = ch
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

func nodesMixedOK(values []*wtype.LHInstruction) bool {
	/// true iff we have exactly two types of node: split and mix
	insTypes := countInstructionTypes(values)

	return len(insTypes) == 2 && insTypes[wtype.InsNames[wtype.LHIMIX]] && insTypes[wtype.InsNames[wtype.LHISPL]]
}

func hasAnySplitNodes(ic *IChain) bool {
	if ic == nil {
		return false
	}

	if ic.Values[0].Type == wtype.LHISPL {
		return true
	}

	return hasAnySplitNodes(ic.Child)
}
func simplifyIChain(ic *IChain, inputs map[string][]*wtype.LHComponent) *IChain {

	if !hasAnySplitNodes(ic) {
		return ic
	}

	// define a graph

	icg := ic.AsGraph()

	// define a graph colouring function

	colourMap, hasColour := getNodeColourMap(ic, inputs)

	colorer := func(n graph.Node) interface{} {
		return colourMap[n]
	}

	hascolour := func(n graph.Node) bool {
		return hasColour[n]
	}

	// derive the quotient graph

	qg := graph.MakeQuotient(graph.MakeQuotientOpt{Graph: icg, Colorer: colorer, HasColor: hascolour, KeepSelfEdges: false})

	// merge to make the new IChain

	return qGraphToIChain(qg, ic)
}

func maxGen(inss []*wtype.LHInstruction, componentGen map[string]int) int {
	max := 0
	for _, ins := range inss {
		for _, c := range ins.Components {
			g, ok := componentGen[c.ID]

			if ok && g > max {
				max = g
			}

			g, ok = componentGen[c.ParentID]

			if ok && g > max {
				max = g
			}
		}
	}

	return max
}

func getNodeColourMap(ic *IChain, inputs map[string][]*wtype.LHComponent) (map[graph.Node]interface{}, map[graph.Node]bool) {
	fmt.Println("GET NODE COLOUR MAP")
	ret := make(map[graph.Node]interface{})
	hc := make(map[graph.Node]bool)

	colour := 0

	// components count as 'live' until something is physically added to them
	// i.e. they are mixed with something else wholesale
	componentsLive := make(map[string]bool)
	componentGen := make(map[string]int)

	// seed live components with inputs

	for _, a := range inputs {
		for _, c := range a {
			componentsLive[c.ID] = true
			componentGen[c.ID] = 0
		}
	}

	front := 0
	for cur := ic; cur != nil; cur = cur.Child {
		if cur.Values[0].Type == wtype.LHIMIX {
			// mix nodes have a different colour if they use any live component

			useOfLiveComponent := func(inss []*wtype.LHInstruction, componentsLive map[string]bool) bool {
				for _, ins := range inss {
					for _, c := range ins.Components {
						if _, ok := componentsLive[c.ID]; ok {
							return true
						}
					}
				}
				return false
			}

			g := maxGen(cur.Values, componentGen)

			if useOfLiveComponent(cur.Values, componentsLive) || g > front {
				colour += 1
			}

			if g > front {
				front = g
			}
			hc[graph.Node(cur)] = true
		} else if cur.Values[0].Type == wtype.LHIPRM {
			// prompts are always kept separate, hence have no colour
			colour += 1
			hc[graph.Node(cur)] = false
		} else if cur.Values[0].Type == wtype.LHISPL {
			hc[graph.Node(cur)] = true
		}

		// delete used components, add products, update IDs
		updateCmpMap(cur.Values, componentsLive)
		updateGenMap(cur.Values, componentGen)
		ret[graph.Node(cur)] = colour
	}

	return ret, hc
}

func updateCmpMap(values []*wtype.LHInstruction, componentsLive map[string]bool) {
	updateMap := func(id1, id2 string, m map[string]bool) {
		_, ok := m[id1]

		if ok {
			delete(m, id1)
			m[id2] = true
		}

	}

	for _, v := range values {
		switch v.Type {
		case wtype.LHISPL:
			updateMap(v.Components[0].ID, v.Results[1].ID, componentsLive)
		case wtype.LHIPRM:
			for in, out := range v.PassThrough {
				updateMap(in, out.ID, componentsLive)
			}
		case wtype.LHIMIX:
			// use inputs
			for _, c := range v.Components {
				delete(componentsLive, c.ID)
			}

			// add output
			componentsLive[v.Results[0].ID] = true
		default:
			panic(fmt.Sprintf("Unknown or irrelevant instruction of type %s passed to instruction sorting", v.InsType()))
		}
	}
}
func updateGenMap(values []*wtype.LHInstruction, componentGen map[string]int) {
	updateMap := func(id1, id2 string, m map[string]int) {
		i, ok := m[id1]

		if ok {
			delete(m, id1)
			m[id2] = i
		}

	}

	for _, v := range values {
		switch v.Type {
		case wtype.LHISPL:
			updateMap(v.Components[0].ID, v.Results[1].ID, componentGen)
		case wtype.LHIPRM:
			for in, out := range v.PassThrough {
				updateMap(in, out.ID, componentGen)
			}
		case wtype.LHIMIX:
			// use inputs
			for _, c := range v.Components {
				delete(componentGen, c.ID)
			}

			// add output
			componentGen[v.Results[0].ID] = maxGen([]*wtype.LHInstruction{v}, componentGen) + 1
		default:
			panic(fmt.Sprintf("Unknown or irrelevant instruction of type %s passed to instruction sorting", v.InsType()))
		}
	}
}

func qGraphToIChain(qg graph.QGraph, orig *IChain) *IChain {
	findRootNode := func(qg graph.QGraph, c *IChain) graph.Node {
		for i := 0; i < qg.NumNodes(); i++ {
			n := qg.Node(i)
			for j := 0; j < qg.NumOrigs(n); j++ {
				if qg.Orig(n, j) == graph.Node(orig) {
					return n
				}
			}
		}

		return nil
	}

	// find the root
	qRootNode := findRootNode(qg, orig)

	if qRootNode == nil {
		fmt.Println(graph.Print(graph.PrintOpt{Graph: qg}))
		panic("No root node found for quotient graph!")
	}

	// now navigate the graph from the root, gathering and updating nodes as we go

	cur := qRootNode
	var ic *IChain

	for {
		// get original nodes
		origNodes := getOrigsAsSlice(qg, cur)

		// convert original nodes to allow component updating to occur correctly

		newNodes := makeNewNodes(origNodes)

		// finally add these to the chain

		ic = addNewNodesTo(ic, newNodes)

		if qg.NumOuts(cur) == 0 {
			break
		}

		// qg here must be a chain
		cur = qg.Out(cur, 0)
	}

	return ic
}

// make slice of original nodes
func getOrigsAsSlice(qg graph.QGraph, cur graph.Node) []*IChain {
	ret := make([]*IChain, 0, qg.NumOrigs(cur))

	for i := 0; i < qg.NumOrigs(cur); i++ {
		ret = append(ret, qg.Orig(cur, i).(*IChain))
	}

	return ret
}

// generate new nodes from old... this involves
// - allowing prompts through as-is
// - finding chains of splits and converting all intermediates into
//   the initial ID then adding a new split to convert the initial to the
//   final ID in the chain
func makeNewNodes(oldNodes []*IChain) *IChain {
	if oldNodes[0].Values[0].Type == wtype.LHIPRM {
		// must be solo

		if len(oldNodes) != 1 {
			panic(fmt.Sprintf("Error: Prompt nodes must appear singly, instead got %d\n", len(oldNodes)))
		}

		oldNodes[0].Parent = nil
		oldNodes[0].Child = nil
		// no need for any replacements
		return oldNodes[0]
	}
	// find out what's in this set

	getNodeTypes := func(nodes []*IChain) map[string]int {
		ret := make(map[string]int)
		for _, n := range nodes {
			_, ok := ret[n.Values[0].InsType()]

			if !ok {
				ret[n.Values[0].InsType()] = 1
			} else {
				ret[n.Values[0].InsType()] += 1
			}
		}
		return ret
	}

	// if just mix or just split, we just merge

	nodeTypes := getNodeTypes(oldNodes)

	if len(nodeTypes) == 1 {
		return mergeSingleTypeNodes(oldNodes)
	}

	return mergeMixedNodes(oldNodes)
}

func mergeSingleTypeNodes(oldNodes []*IChain) *IChain {
	values := make([]*wtype.LHInstruction, 0, len(oldNodes[0].Values))

	for _, ic := range oldNodes {
		values = append(values, ic.Values...)
	}

	return &IChain{Values: values}
}

func mergeMixedNodes(oldNodes []*IChain) *IChain {
	findNodes := func(nodes []*IChain, nodeType int) []*IChain {
		ret := make([]*IChain, 0, len(nodes))

		for _, v := range nodes {
			if v.Values[0].Type == nodeType {
				ret = append(ret, v)
			}
		}

		return ret
	}

	// generate two nodes: one of splits, one of mixes

	splitNodes := findNodes(oldNodes, wtype.LHISPL)
	mixNodes := findNodes(oldNodes, wtype.LHIMIX)

	if len(splitNodes)+len(mixNodes) != len(oldNodes) {
		fmt.Println("SPLIT: ", splitNodes)
		fmt.Println("MIX  : ", mixNodes)
		fmt.Println("OLD  : ", oldNodes)
		panic(fmt.Sprintf("oldNodes (%d) does not partition into splitNodes (%d) + mixNodes (%d)", len(oldNodes), len(splitNodes), len(mixNodes)))
	}

	// get chain of updates
	updateChain := getUpdateChain(splitNodes)

	// convert the mixes to all use the same initial component
	convertWithChain(mixNodes, updateChain)

	// convert the splits
	splitNodes = pruneSplits(splitNodes, updateChain)

	// pare down
	ret := mergeSingleTypeNodes(mixNodes)
	ret.Child = mergeSingleTypeNodes(splitNodes)

	return ret
}

// getUpdateChain reads the ordered list of split instructions passed in and finds
// chains of ID updates from a->b->c->... Each map returned contains one such chain
func getUpdateChain(justSplitInstructionNodes []*IChain) []map[string]string {
	// find all chains of split style updates from a->b->c->... and return
	// instructions come in ordered

	ret := make([]map[string]string, 0, len(justSplitInstructionNodes))

	currMap := make(map[string]string)

	ret = append(ret, currMap)

	for _, node := range justSplitInstructionNodes {

		for _, ins := range node.Values {

			if ins.Type != wtype.LHISPL {
				panic(fmt.Sprintf("Error: passed non-split (type %s) instruction to getUpdateChain", ins.InsType()))
			}

			// update to find if we have any previous entry

			ok := false

			for _, m := range ret {
				_, ok = m[ins.Components[0].ID]

				if ok {
					currMap = m
					break
				}
			}

			if !ok {
				if len(currMap) != 0 {
					// start a new map
					currMap = make(map[string]string)
					ret = append(ret, currMap)
				}
				currMap[""] = ins.Components[0].ID
			}

			currMap[ins.Components[0].ID] = ins.Results[1].ID
			currMap[ins.Results[1].ID] = ""
		}
	}

	return ret
}

//convertWithChain takes a list of mix instruction nodes
// and converts them all to use the first version of components mentioned in any splits we've seen
func convertWithChain(justMixNodes []*IChain, updateChains []map[string]string) {
	// find what everything maps to
	updateMap := getOnePassUpdateMap(updateChains)

	// now apply to the mix instructions

	for _, node := range justMixNodes {
		for _, ins := range node.Values {
			for _, c := range ins.Components {
				newID, ok := updateMap[c.ParentID]

				if ok {
					c.ParentID = newID
				}

			}
		}
	}

}

// getOnePassUpdateMap gives us a map to look up what the new IDs for components are
// so that all mixes can be simply assigned the first version of the component's ID
// since chains are distinct we can merge them here
// e.g. this function takes a set of maps map which goes [{"":"a", "a":"b", "b":"c", "c":""}, {"":"d", "d":"e","e":"}]
// and returns one which goes {"a":"a", "b","a", "d":"d"}
func getOnePassUpdateMap(chains []map[string]string) map[string]string {
	// all chains are distinct here, so we can merge
	r := make(map[string]string, len(chains[0]))

	for _, chain := range chains {
		first := chain[""]

		for k, v := range chain {
			if k == "" || v == "" {
				continue
			} else {
				r[k] = first
			}
		}
	}

	return r
}

// takes a set of maps like this [{"":"a", "a":"b", "b":""}, {"":"c", "c":"d", "d","e", "e":""}]
// and returns a single map like this {"a":"b", "c","e"}
func getSplitUpdateMap(chains []map[string]string) map[string]string {
	ret := make(map[string]string, len(chains))
	for _, chain := range chains {
		start, end := getStartEnd(chain)
		ret[start] = end
	}

	return ret
}

// takes a map like this {"":"a", "a":"b", "b":""}
// and returns a,b
func getStartEnd(chain map[string]string) (start, end string) {
	start = chain[""]

	for k, v := range chain {
		if v == "" {
			end = k
		}
	}

	return
}

func pruneSplits(justSplitNodes []*IChain, chains []map[string]string) []*IChain {
	// collapse split instructions which are covered by chains found above into
	// single splits (i.e. a->...->z ---> a->z)

	updateMap := getSplitUpdateMap(chains)

	ret := make([]*IChain, 0, len(justSplitNodes))

	// now remove all splits which concern component versions we are invalidating

	for _, node := range justSplitNodes {
		inss := make([]*wtype.LHInstruction, 0, len(justSplitNodes))
		for _, split := range node.Values {
			finalID, ok := updateMap[split.Components[0].ID]

			if ok {
				// we only keep the split that goes from component v1
				split.Results[1].ID = finalID
				inss = append(inss, split)
			}
		}
		node.Values = inss
		if len(node.Values) != 0 {
			ret = append(ret, node)
		}
	}

	return ret
}

// append nodes to chain
// ic, newNodes must already be linked, any existing parent link from newNodes[0] is overwritten
// any child link from ic[len(ic)-1] is also overwritten
func addNewNodesTo(ic *IChain, newNodes *IChain) *IChain {
	// newNodes must be a chain

	if ic == nil {
		return newNodes
	}

	if newNodes == nil {
		return ic
	}

	var cur *IChain

	for cur = ic; cur.Child != nil; cur = cur.Child {
	}

	cur.Child = newNodes
	newNodes.Parent = cur

	last := cur

	for cur := newNodes; cur != nil; cur = cur.Child {

		if cur.Parent == nil {
			cur.Parent = last
		}

		if cur.Parent != nil {
			cur.Depth = cur.Parent.Depth + 1
		}

		last = cur
	}

	return ic
}
