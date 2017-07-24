package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"testing"
)

func TestTGraph(t *testing.T) {
	tIns := make([]*wtype.LHInstruction, 0, 10)

	cmpIn := wtype.NewLHComponent()

	for k := 0; k < 10; k++ {
		ins := wtype.NewLHMixInstruction()
		cmpOut := wtype.NewLHComponent()
		ins.AddComponent(cmpIn)
		ins.AddProduct(cmpOut)
		tIns = append(tIns, ins)
		cmpIn = cmpOut
	}

	tgraph := MakeTGraph(tIns)

	arrEq := func(ar1 []*wtype.LHInstruction, ar2 []*wtype.LHInstruction) bool {
		if len(ar1) != len(ar2) {
			return false
		}

		for i := 0; i < len(ar1); i++ {
			if ar1[i].ID != ar2[i].ID {
				return false
			}
		}

		return true
	}

	if !arrEq(tIns, tgraph.Nodes) {
		t.Errorf("Nodes in tGraph should be identical to nodes inputted")
	}

	if tgraph.NumNodes() != 10 {
		t.Errorf("NumNodes should report 10, instead reports %d", tgraph.NumNodes)
	}

	// edge check... we should have 9->8->7->6->5

	for k := 9; k >= 1; k-- {
		expect := tgraph.Node(k - 1)
		got := tgraph.Out(tgraph.Node(k), 0)

		if expect != got {
			t.Errorf("this graph should be a chain - failed between nodes %d %d", k, k-1)
		}
	}
}
