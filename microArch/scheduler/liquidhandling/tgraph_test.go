package liquidhandling

import (
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
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

	tgraph, err := MakeTGraph(tIns)
	if err != nil {
		t.Error(err)
	}

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
		t.Errorf("NumNodes should report 10, instead reports %d", tgraph.NumNodes())
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

// TestTGraphSplit checks whether split instructions are working correctly
// these are sorted so that they occur after use of their first return and
// before the use of their second - this is because they update the ID of their
// input component
func TestTGraphSplit(t *testing.T) {
	tIns := make([]*wtype.LHInstruction, 0, 3)

	cmpIn := wtype.NewLHComponent()
	moving, remaining := mixer.SplitSample(cmpIn, wunit.NewVolume(100.0, "ul"))

	cmpOut := wtype.NewLHComponent()

	// mix
	ins := wtype.NewLHMixInstruction()

	ins.AddComponent(moving)
	ins.AddProduct(cmpOut)
	tIns = append(tIns, ins)

	// split
	ins = wtype.NewLHSplitInstruction()
	ins.AddComponent(cmpOut)

	ins.AddProduct(moving)
	ins.AddProduct(remaining)

	tIns = append(tIns, ins)

	// mix again

	ins = wtype.NewLHMixInstruction()

	ins.AddComponent(remaining)
	ins.AddProduct(wtype.NewLHComponent())
	tIns = append(tIns, ins)

	tgraph, err := MakeTGraph(tIns)
	if err != nil {
		t.Error(err)
	}

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

	if tgraph.NumNodes() != 3 {
		t.Errorf("NumNodes should report 3, instead reports %d", tgraph.NumNodes())
	}

	// edge check... we should have 3->2->1

	for k := 2; k >= 1; k-- {
		expect := tgraph.Node(k - 1)
		got := tgraph.Out(tgraph.Node(k), 0)

		if expect != got {
			t.Errorf("this graph should be a chain - failed between nodes %d %d", k, k-1)
		}
	}
}
