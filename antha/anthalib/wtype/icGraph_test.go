package wtype

import (
	"github.com/antha-lang/antha/graph"
	"testing"
)

func TestCullDeadNodes(t *testing.T) {
	icA := []*IChain{{Values: []*LHInstruction{NewLHMixInstruction()}}, {Values: []*LHInstruction{}}}

	icB := cullDeadNodes(icA)

	if len(icB) != 1 {
		t.Errorf("Expected 1 node after culling, instead got %d", len(icB))
	}
}

func TestNodesMixedOK(t *testing.T) {
	tests := [][]*LHInstruction{{{Type: LHIMIX}, {Type: LHISPL}},
		{{Type: LHIMIX}, {Type: LHIMIX}},
		{{Type: LHIPRM}, {Type: LHISPL}},
		{{Type: LHIMIX}, {Type: LHIPRM}}}
	wants := []bool{true, false, false, false}
	names := []string{"MixSplit", "MixMix", "PromptSplit", "MixPrompt"}

	for i := range tests {
		doTheTest := func(t *testing.T) {
			ic := &IChain{Values: tests[i]}
			if got := ic.hasMixAndSplitOnly(); got != wants[i] {
				t.Errorf("Expected %t got %t", wants[i], got)
			}
		}

		t.Run(names[i], doTheTest)
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

func cullDeadNodes(in []*IChain) (out []*IChain) {
	for _, v := range in {
		if len(v.Values) != 0 {
			out = append(out, v)
		}
	}
	return out
}
