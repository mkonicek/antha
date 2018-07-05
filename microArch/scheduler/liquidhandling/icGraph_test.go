package liquidhandling

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/graph"
	"reflect"
	"testing"
)

func TestCullDeadNodes(t *testing.T) {
	icA := []*IChain{{Values: []*wtype.LHInstruction{wtype.NewLHMixInstruction()}}, {Values: []*wtype.LHInstruction{}}}

	icB := cullDeadNodes(icA)

	if len(icB) != 1 {
		t.Errorf("Expected 1 node after culling, instead got %d", len(icB))
	}
}

func TestGetStartEnd(t *testing.T) {
	m := map[string]string{"": "a", "a": "b", "b": "c", "c": ""}

	s, e := getStartEnd(m)

	if s != "a" {
		t.Errorf("getStartEnd error: Expected start result \"a\", instead got \"%s\"", s)
	}

	if e != "c" {
		t.Errorf("getStartEnd error: Expected end result \"c\", instead got \"%s\"", e)
	}
}

func TestGetSplitUpdateMap(t *testing.T) {
	m := []map[string]string{{"": "a", "a": "b", "b": "c", "c": ""}, {"": "d", "d": "e", "e": "f", "f": ""}}
	expected := map[string]string{"a": "c", "d": "f"}

	ret := getSplitUpdateMap(m)

	if !reflect.DeepEqual(ret, expected) {
		t.Errorf("getSplitUpdateMap error: Expected %v got %v", expected, ret)
	}
}

func TestGetOnePassUpdateMap(t *testing.T) {
	m := []map[string]string{{"": "a", "a": "b", "b": "c", "c": ""}, {"": "d", "d": "e", "e": "f", "f": ""}}
	expected := map[string]string{"a": "a", "b": "a", "d": "d", "e": "d"}

	ret := getOnePassUpdateMap(m)

	if !reflect.DeepEqual(ret, expected) {
		t.Errorf("getOnePassUpdateMap error: Expected %v got %v", expected, ret)
	}
}

func TestConvertWithChain(t *testing.T) {
	m := []map[string]string{{"": "a", "a": "b", "b": "c", "c": ""}, {"": "d", "d": "e", "e": "f", "f": ""}}

	inss := []*wtype.LHInstruction{getAnInstruction(getComponentsWithIDsParents([]string{"v", "w"}, []string{"a", "d"}), getComponentsWithIDs([]string{"p"}), wtype.LHIMIX), getAnInstruction(getComponentsWithIDsParents([]string{"x", "y"}, []string{"b", "e"}), getComponentsWithIDs([]string{"q"}), wtype.LHIMIX)}
	nodes := getIChain([][]*wtype.LHInstruction{inss}).AsList([]*IChain{})

	inss2 := []*wtype.LHInstruction{getAnInstruction(getComponentsWithIDsParents([]string{"v", "w"}, []string{"a", "d"}), getComponentsWithIDs([]string{"p"}), wtype.LHIMIX), getAnInstruction(getComponentsWithIDsParents([]string{"x", "y"}, []string{"a", "d"}), getComponentsWithIDs([]string{"q"}), wtype.LHIMIX)}
	expected := getIChain([][]*wtype.LHInstruction{inss2}).AsList([]*IChain{})

	convertWithChain(nodes, m)

	for i := range nodes {
		n := nodes[i]
		e := expected[i]

		for j := range n.Values {
			if n.Values[j].Components[0].ID != e.Values[j].Components[0].ID {
				t.Errorf("convertWithChain error: Expected %v got %v", expected, nodes)
			}
		}
	}
}

func TestAddNewNodesTo(t *testing.T) {
	getEnd := func(ic *IChain) *IChain {
		cur := ic
		for ; cur.Child != nil; cur = cur.Child {
		}
		return cur
	}

	lastHead := &IChain{}
	head := &IChain{Child: &IChain{Child: &IChain{Child: lastHead}}}
	lastTail := &IChain{}
	tail := &IChain{Child: &IChain{Child: lastTail}}

	// test 1: empty head

	r := addNewNodesTo(nil, tail)

	if r != tail {
		t.Errorf("addNewNodesTo Failure: must return tail if head is nil")
	}

	// test 2: empty tail

	r = addNewNodesTo(head, nil)

	if r != head {
		t.Errorf("addNewNodesTo Failure: must return head if tail is nil")
	}

	// test 3: both supplied

	r = addNewNodesTo(head, tail)

	if r == nil {
		t.Errorf("addNewNodesTo Failure: must not return nil")
	}

	e := getEnd(r)

	if e != lastTail {
		t.Errorf("addNewNodesTo Failure: head not joined to tail")
	}
}

func TestNodesMixedOK(t *testing.T) {
	tests := [][]*wtype.LHInstruction{{{Type: wtype.LHIMIX}, {Type: wtype.LHISPL}},
		{{Type: wtype.LHIMIX}, {Type: wtype.LHIMIX}},
		{{Type: wtype.LHIPRM}, {Type: wtype.LHISPL}},
		{{Type: wtype.LHIMIX}, {Type: wtype.LHIPRM}}}
	wants := []bool{true, false, false, false}
	names := []string{"MixSplit", "MixMix", "PromptSplit", "MixPrompt"}

	for i := range tests {
		doTheTest := func(t *testing.T) {
			got := nodesMixedOK(tests[i])
			if got != wants[i] {
				t.Errorf("Expected %t got %t", wants[i], got)
			}
		}

		t.Run(names[i], doTheTest)
	}
}

func TestUpdateCmpMap(t *testing.T) {
	componentsLive := map[string]bool{"A": true, "B": true, "C": true, "D": true, "E": true}
	mix := wtype.NewLHMixInstruction()
	mix.Components = []*wtype.Liquid{{ID: "A"}}
	mix.Results = []*wtype.Liquid{{ID: "W"}}

	split := wtype.NewLHSplitInstruction()

	split.Components = []*wtype.Liquid{{ID: "B"}}
	split.Results = []*wtype.Liquid{{ID: "S"}, {ID: "X"}}

	prompt := wtype.NewLHPromptInstruction()
	prompt.PassThrough = map[string]*wtype.Liquid{"C": {ID: "Y"}, "D": {ID: "Z"}}

	expected := map[string]bool{"W": true, "X": true, "Y": true, "Z": true, "E": true}

	values := []*wtype.LHInstruction{mix, split, prompt}

	updateCmpMap(values, componentsLive)

	if !reflect.DeepEqual(expected, componentsLive) {
		t.Errorf("updateCmpMap ERROR: expected %v got %v", expected, componentsLive)
	}

}

func TestIcAsGraph(t *testing.T) {
	node3 := &IChain{}
	node2 := &IChain{Child: node3}
	node1 := &IChain{Child: node2}

	g := node1.AsGraph()

	if g.NumNodes() != 3 {
		t.Errorf("AsGraph error: Expected 3 nodes got %d", g.NumNodes())
	}

	n := 1
	for curr := node1; curr != nil; curr = curr.Child {
		expected := 1
		if curr == node3 {
			expected = 0
		}

		out := g.NumOuts(graph.Node(curr))

		if out != expected {
			t.Errorf("Expected %d outs for node, got %d", expected, out)
		}

		if expected != 0 {
			exOut := graph.Node(curr.Child)
			if g.Out(curr, 0) != exOut {
				t.Errorf("Expected to get out %v instead got %v", exOut, g.Out(curr, 0))
			}
		}
		n += 1
	}
}

func getComponentsWithIDs(IDs []string) []*wtype.Liquid {
	r := make([]*wtype.Liquid, 0, len(IDs))

	for _, id := range IDs {
		r = append(r, &wtype.Liquid{ID: id})
	}

	return r
}

func getComponentsWithIDsParents(IDs, parents []string) []*wtype.Liquid {
	if len(IDs) != len(parents) {
		panic(fmt.Sprintf("IDS and parents not same lengths: %d vs %d", len(IDs), len(parents)))
	}
	cmps := getComponentsWithIDs(IDs)

	for i, v := range parents {
		cmps[i].ParentID = v
	}

	return cmps
}

func getIChain(inss [][]*wtype.LHInstruction) *IChain {
	chain := &IChain{}
	cur := chain
	for i, a := range inss {
		cur.Values = a
		if i != len(inss)-1 {
			cur.Child = &IChain{}
			cur = cur.Child
		}
	}

	return chain
}

func TestGetNodeColourMapSimple(t *testing.T) {
	inputs := map[string][]*wtype.Liquid{"whyisthisamapanyway": getComponentsWithIDs([]string{"A1", "B1"})}

	mix1 := wtype.NewLHMixInstruction()
	mix1.Components = getComponentsWithIDs([]string{"S1", "T1"})
	mix1.Results = getComponentsWithIDs([]string{"P1"})
	split1 := wtype.NewLHSplitInstruction()
	split1.Components = getComponentsWithIDs([]string{"A1"})
	split1.Results = getComponentsWithIDs([]string{"S1", "A2"})
	split2 := wtype.NewLHSplitInstruction()
	split2.Components = getComponentsWithIDs([]string{"B1"})
	split2.Results = getComponentsWithIDs([]string{"T1", "B2"})
	mix2 := wtype.NewLHMixInstruction()
	mix2.Components = getComponentsWithIDs([]string{"A2", "B2"})
	mix2.Results = getComponentsWithIDs([]string{"P2"})

	ic := getIChain([][]*wtype.LHInstruction{{mix1}, {split1, split2}, {mix2}})

	colourMap, noColourMap := getNodeColourMap(ic, inputs)

	nodes := ic.AsList([]*IChain{})

	// all nodes should have a colour and be in the order 0,0,1

	expected := []int{0, 0, 1}

	for i := range expected {
		c, ok := colourMap[nodes[i]]

		if !ok {
			t.Errorf("Expected colour for node %d, got none", i)
		}

		if c != expected[i] {
			t.Errorf("Expected colour %d for node %d, got %d", expected[i], i, c)
		}

		hasColour := noColourMap[nodes[i]]

		if !hasColour {
			t.Errorf("Node %d must not report no colour yet it does", i)
		}
	}
}

// TODO--> add the below
/*
func TestgetNodeColourMap2(t *testing.T)
*/

func getAnInstruction(in, out []*wtype.Liquid, whatType int) *wtype.LHInstruction {
	var ins *wtype.LHInstruction
	if whatType == wtype.LHIMIX {
		ins = wtype.NewLHMixInstruction()
	} else if whatType == wtype.LHISPL {
		ins = wtype.NewLHSplitInstruction()
	}
	ins.Components = in
	ins.Results = out
	return ins
}

func getTestNodesForSplit() ([]*IChain, []map[string]string, []*wtype.LHInstruction) {
	nodes := []*IChain{}

	nodes = append(nodes, &IChain{Values: []*wtype.LHInstruction{getAnInstruction(getComponentsWithIDs([]string{"A"}), getComponentsWithIDs([]string{"S", "B"}), wtype.LHISPL)}})
	nodes = append(nodes, &IChain{Values: []*wtype.LHInstruction{getAnInstruction(getComponentsWithIDs([]string{"B"}), getComponentsWithIDs([]string{"T", "C"}), wtype.LHISPL)}})
	nodes = append(nodes, &IChain{Values: []*wtype.LHInstruction{getAnInstruction(getComponentsWithIDs([]string{"C"}), getComponentsWithIDs([]string{"U", "D"}), wtype.LHISPL)}})
	nodes = append(nodes, &IChain{Values: []*wtype.LHInstruction{getAnInstruction(getComponentsWithIDs([]string{"X"}), getComponentsWithIDs([]string{"V", "Y"}), wtype.LHISPL)}})
	nodes = append(nodes, &IChain{Values: []*wtype.LHInstruction{getAnInstruction(getComponentsWithIDs([]string{"Y"}), getComponentsWithIDs([]string{"W", "Z"}), wtype.LHISPL)}})
	chain := []map[string]string{{"": "A", "A": "B", "B": "C", "C": "D", "D": ""}, {"": "X", "X": "Y", "Y": "Z", "Z": ""}}

	postMerge := []*wtype.LHInstruction{getAnInstruction(getComponentsWithIDs([]string{"A"}), getComponentsWithIDs([]string{"S", "D"}), wtype.LHISPL), getAnInstruction(getComponentsWithIDs([]string{"X"}), getComponentsWithIDs([]string{"V", "Z"}), wtype.LHISPL)}

	return nodes, chain, postMerge
}

func TestGetUpdateChain(t *testing.T) {
	nodes, expected, _ := getTestNodesForSplit()
	updateChains := getUpdateChain(nodes)

	if !reflect.DeepEqual(updateChains, expected) {
		t.Errorf("getUpdateChain error: expected %v got %v", expected, updateChains)
	}
}

func TestGetUpdateChain2(t *testing.T) {
	splitNodes, expected, _ := getTestNodesForSplit()
	// merge some of the nodes
	splitNodes[0].Values = append(splitNodes[0].Values, splitNodes[3].Values...)
	splitNodes[1].Values = append(splitNodes[1].Values, splitNodes[4].Values...)
	splitNodes2 := []*IChain{splitNodes[0], splitNodes[1], splitNodes[2]}
	updateChains := getUpdateChain(splitNodes2)

	if !reflect.DeepEqual(updateChains, expected) {
		t.Errorf("getUpdateChain error: expected %v got %v", expected, updateChains)
	}
}

func TestPruneSplits(t *testing.T) {
	nodes, chain, _ := getTestNodesForSplit()
	pruned := pruneSplits(nodes, chain)

	expected := []*IChain{{Values: []*wtype.LHInstruction{getAnInstruction(getComponentsWithIDs([]string{"A"}), getComponentsWithIDs([]string{"S", "D"}), wtype.LHISPL)}}, {Values: []*wtype.LHInstruction{getAnInstruction(getComponentsWithIDs([]string{"X"}), getComponentsWithIDs([]string{"V", "Z"}), wtype.LHISPL)}}}

	if len(pruned) != len(expected) {
		t.Errorf("Unexpected number of nodes after pruning: expect %d got %d", len(expected), len(pruned))
	}

	for i := range expected {
		n1 := expected[i]
		n2 := pruned[i]

		if len(n1.Values) != len(n2.Values) {
			t.Errorf("Unexpected number of instructions in node %d after pruning: expect %d got %d", i, len(n1.Values), len(n2.Values))
		}

		for j := range n1.Values {
			ins1 := n1.Values[j]
			ins2 := n2.Values[j]

			getIDs := func(ins *wtype.LHInstruction) []string {
				r := []string{}

				for _, c := range ins.Components {
					r = append(r, c.ID)
				}

				for _, rslt := range ins.Results {
					r = append(r, rslt.ID)
				}

				return r
			}

			if !reflect.DeepEqual(getIDs(ins1), getIDs(ins2)) {
				t.Errorf("Splits not correct post pruning node %d ins %d. Expect %v got %v", i, j, getIDs(ins1), getIDs(ins2))
			}
		}
	}
}

func TestMergeSingleTypeNodes(t *testing.T) {
	nodes, _, _ := getTestNodesForSplit()
	merged := mergeSingleTypeNodes(nodes)

	if len(merged.Values) != len(nodes) {
		t.Errorf("mergeSingleTypeNodes error: expected %d values in result got %d", len(nodes), len(merged.Values))
	}

	for i := range nodes {
		if !reflect.DeepEqual(merged.Values[i], nodes[i].Values[0]) {
			t.Errorf("mergeSingleTypeNodes error: expected instruction %v at position %d, got %v", nodes[i].Values[0], i, merged.Values[i])
		}
	}
}

func getTestNodesForMix() ([]*IChain, []*wtype.LHInstruction) {
	ret := make([]*IChain, 0, 3)
	cids := [][]string{{"S", "V"}, {"T", "W"}, {"U"}}
	pids := [][]string{{"P1"}, {"P2"}, {"P3"}}
	parentids := [][]string{{"A", "X"}, {"B", "Y"}, {"C"}}
	postparentids := [][]string{{"A", "X"}, {"A", "X"}, {"A"}}

	postMerge := make([]*wtype.LHInstruction, 0, 3)

	for i := range cids {
		mix := getAnInstruction(getComponentsWithIDsParents(cids[i], parentids[i]), getComponentsWithIDs(pids[i]), wtype.LHIMIX)
		ret = append(ret, getIChain([][]*wtype.LHInstruction{{mix}}))

		postMerge = append(postMerge, getAnInstruction(getComponentsWithIDsParents(cids[i], postparentids[i]), getComponentsWithIDs(pids[i]), wtype.LHIMIX))
	}

	return ret, postMerge
}

func TestMergeMixedNodes(t *testing.T) {
	mixNodes, expectedPostMerge := getTestNodesForMix()

	splitNodes, _, splitPostMerge := getTestNodesForSplit()
	// merge some of the nodes
	splitNodes[0].Values = append(splitNodes[0].Values, splitNodes[3].Values...)
	splitNodes[1].Values = append(splitNodes[1].Values, splitNodes[4].Values...)
	splitNodes2 := []*IChain{splitNodes[0], splitNodes[1], splitNodes[2]}

	// interleave the nodes

	allNodes := []*IChain{}

	for i := 0; i < len(mixNodes); i++ {
		allNodes = append(allNodes, mixNodes[i])
		allNodes = append(allNodes, splitNodes2[i])
	}

	ch := mergeMixedNodes(allNodes)

	expected := &IChain{Values: expectedPostMerge, Child: &IChain{Values: splitPostMerge}}

	if ch.Height() != expected.Height() {
		t.Errorf("LENGTH MISMATCH WANT %d GOT %d", ch.Height(), expected.Height())
	}

	cur2 := expected
	for cur := ch; cur != nil; cur = cur.Child {
		if cur2 == nil {
			t.Errorf("Chain of excess length returned")
		}

		if len(cur.Values) != len(cur2.Values) {
			t.Errorf("Wrong length for node of iChain: expected %d got %d", len(cur2.Values), len(cur.Values))
		}

		for i := range cur.Values {
			if cur.Values[i].InsType() != cur2.Values[i].InsType() {
				t.Errorf("Instruction type mismatch in iChain: expected %s got %s", cur2.Values[i].InsType(), cur.Values[i].InsType())
			}

			gcmp := cur.Values[i].Components
			grsl := cur.Values[i].Results

			ecmp := cur2.Values[i].Components
			ersl := cur2.Values[i].Results

			if len(ecmp) != len(gcmp) {
				t.Errorf("Wrong number of components for instruction in iChain: expected %d got %d", len(ecmp), len(gcmp))
			}

			for j := range gcmp {
				if ecmp[j].ID != gcmp[j].ID || ecmp[j].ParentID != gcmp[j].ParentID {
					t.Errorf("ID Mismatch for instruction type %s in iChain: expected ID %s Parent %s got ID %s Parent %s", cur.Values[i].InsType(), ecmp[j].ID, ecmp[j].ParentID, gcmp[j].ID, gcmp[j].ParentID)
				}
			}

			for j := range grsl {
				if ersl[j].ID != grsl[j].ID || ersl[j].ParentID != grsl[j].ParentID {
					t.Errorf("ID Mismatch for instruction type %s in iChain: expected ID %s Parent %s got ID %s Parent %s", cur.Values[i].InsType(), ersl[j].ID, ersl[j].ParentID, grsl[j].ID, grsl[j].ParentID)
				}
			}
		}

		cur2 = cur2.Child
	}

}

func TestMakeNewNodes(t *testing.T) {
	splitNodes, _, splitPostMerge := getTestNodesForSplit()

	// merge some of the nodes
	splitNodes[0].Values = append(splitNodes[0].Values, splitNodes[3].Values...)
	splitNodes[1].Values = append(splitNodes[1].Values, splitNodes[4].Values...)
	splitNodes2 := []*IChain{splitNodes[0], splitNodes[1], splitNodes[2]}

	mixNodes, expectedPostMerge := getTestNodesForMix()

	// interleave the nodes

	allNodes := []*IChain{}

	for i := 0; i < len(mixNodes); i++ {
		allNodes = append(allNodes, mixNodes[i])
		allNodes = append(allNodes, splitNodes2[i])
	}

	ch := makeNewNodes(allNodes)

	expected := &IChain{Values: expectedPostMerge, Child: &IChain{Values: splitPostMerge}}

	cur2 := expected
	for cur := ch; cur != nil; cur = cur.Child {
		if cur2 == nil {
			t.Errorf("Chain of excess length returned")
		}

		if len(cur.Values) != len(cur2.Values) {
			t.Errorf("Wrong length for node of iChain: expected %d got %d", len(cur2.Values), len(cur.Values))
		}

		for i := range cur.Values {
			if cur.Values[i].InsType() != cur2.Values[i].InsType() {
				t.Errorf("Instruction type mismatch in iChain: expected %s got %s", cur2.Values[i].InsType(), cur.Values[i].InsType())
			}

			gcmp := cur.Values[i].Components
			grsl := cur.Values[i].Results

			ecmp := cur2.Values[i].Components
			ersl := cur2.Values[i].Results

			if len(ecmp) != len(gcmp) {
				t.Errorf("Wrong number of components for instruction in iChain: expected %d got %d", len(ecmp), len(gcmp))
			}

			for j := range gcmp {
				if ecmp[j].ID != gcmp[j].ID || ecmp[j].ParentID != gcmp[j].ParentID {
					t.Errorf("ID Mismatch for instruction type %s in iChain: expected ID %s Parent %s got ID %s Parent %s", cur.Values[i].InsType(), ecmp[j].ID, ecmp[j].ParentID, gcmp[j].ID, gcmp[j].ParentID)
				}
			}

			for j := range grsl {
				if ersl[j].ID != grsl[j].ID || ersl[j].ParentID != grsl[j].ParentID {
					t.Errorf("ID Mismatch for instruction type %s in iChain: expected ID %s Parent %s got ID %s Parent %s", cur.Values[i].InsType(), ersl[j].ID, ersl[j].ParentID, grsl[j].ID, grsl[j].ParentID)
				}
			}
		}

		cur2 = cur2.Child
	}

}

func cullDeadNodes(in []*IChain) (out []*IChain) {
	for _, v := range in {
		if len(v.Values) != 0 {
			out = append(out, v)
		}
	}
	return out
}

// getNodeColourMap(ic *IChain, inputs map[string][]*wtype.LHComponent) (map[graph.Node]interface{}, map[graph.Node]bool)
// getAnInstruction(in, out []*wtype.LHComponent, whatType int) *wtype.LHInstruction {
// getComponentsWithIDsParents(IDs, parents []string)
// getIChain([][]*wtype.LHInstruction)

// just a chain of single instructions, each using the output of the last
func getMixChainWithIDsParents(IDs, parents []string) *IChain {
	cmps := getComponentsWithIDsParents(IDs, parents)
	inss := make([][]*wtype.LHInstruction, len(IDs)-1)
	for i := 0; i < len(IDs)-1; i++ {
		inss[i] = []*wtype.LHInstruction{getAnInstruction([]*wtype.Liquid{cmps[i]}, []*wtype.Liquid{cmps[i+1]}, wtype.LHIMIX)}
	}

	return getIChain(inss)
}

func getMixChainWithIDsParentsProducts(IDs, parents, products []string) *IChain {
	cmps := getComponentsWithIDsParents(IDs, parents)
	prds := getComponentsWithIDs(products)
	inss := make([][]*wtype.LHInstruction, len(IDs))
	for i := 0; i < len(IDs); i++ {
		inss[i] = []*wtype.LHInstruction{getAnInstruction([]*wtype.Liquid{cmps[i]}, []*wtype.Liquid{prds[i]}, wtype.LHIMIX)}
	}

	return getIChain(inss)
}

func TestColourMapNoSplits(t *testing.T) {
	mixChain := getMixChainWithIDsParents([]string{"A", "B", "C", "D", "E"}, []string{"", "", "", "", ""})

	inputs := map[string][]*wtype.Liquid{"A": getComponentsWithIDsParents([]string{"A"}, []string{""})}

	colourMap, hasColour := getNodeColourMap(mixChain, inputs)

	// all should have different colours

	expected := []int{1, 2, 3, 4}

	i := 0
	for curr := mixChain; curr != nil; curr = curr.Child {

		if !hasColour[curr] {
			t.Errorf("Expected colour for node %d in chain", i)
		}

		if colourMap[curr] != expected[i] {
			t.Errorf("Expected colour %d for node %d, got %d", expected[i], i, colourMap[curr])
		}

		i += 1
	}

}

func TestColourMapNoSplits2(t *testing.T) {
	mixChain := getMixChainWithIDsParentsProducts([]string{"A", "B", "C", "D", "E"}, []string{"U", "V", "W", "X", "Y"}, []string{"V", "W", "X", "Y", "Z"})

	inputs := map[string][]*wtype.Liquid{"U": getComponentsWithIDsParents([]string{"U"}, []string{""})}

	colourMap, hasColour := getNodeColourMap(mixChain, inputs)

	// all should have different colours

	expected := []int{0, 1, 2, 3, 4}

	i := 0
	for curr := mixChain; curr != nil; curr = curr.Child {

		if !hasColour[curr] {
			t.Errorf("Expected colour for node %d in chain", i)
		}

		if colourMap[curr] != expected[i] {
			t.Errorf("Expected colour %d for node %d, got %d", expected[i], i, colourMap[curr])
		}

		i += 1
	}

}
