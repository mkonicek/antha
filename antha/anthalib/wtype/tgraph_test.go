package wtype

import (
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"testing"
)

func splitSample(l *Liquid, v wunit.Volume) (moving, remaining *Liquid) {
	remaining = l.Dup()

	moving = sample(remaining, v)

	remaining.Vol -= v.ConvertToString(remaining.Vunit)
	remaining.ID = GetUUID()

	return
}

// sample takes a sample of volume v from this liquid
func sample(l *Liquid, v wunit.Volume) *Liquid {
	ret := NewLHComponent()
	//      ret.ID = l.ID
	l.AddDaughterComponent(ret)
	ret.ParentID = l.ID
	ret.CName = l.Name()
	ret.Type = l.Type
	ret.Vol = v.RawValue()
	ret.Vunit = v.Unit().PrefixedSymbol()
	ret.Extra = l.GetExtra()
	ret.SubComponents = l.SubComponents
	ret.Smax = l.GetSmax()
	ret.Visc = l.GetVisc()
	if l.Conc > 0 && len(l.Cunit) > 0 {
		ret.SetConcentration(wunit.NewConcentration(l.Conc, l.Cunit))
	}

	ret.SetSample(true)

	return ret
}

func TestTGraph(t *testing.T) {
	tIns := make([]*LHInstruction, 0, 10)

	cmpIn := NewLHComponent()

	for k := 0; k < 10; k++ {
		ins := NewLHMixInstruction()
		cmpOut := NewLHComponent()
		ins.AddInput(cmpIn)
		ins.AddOutput(cmpOut)
		tIns = append(tIns, ins)
		cmpIn = cmpOut
	}

	tgraph, err := MakeTGraph(tIns)
	if err != nil {
		t.Error(err)
	}

	arrEq := func(ar1 []*LHInstruction, ar2 []*LHInstruction) bool {
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
	tIns := make([]*LHInstruction, 0, 3)

	cmpIn := NewLHComponent()
	moving, remaining := splitSample(cmpIn, wunit.NewVolume(100.0, "ul"))

	cmpOut := NewLHComponent()

	// mix
	ins := NewLHMixInstruction()

	ins.AddInput(moving)
	ins.AddOutput(cmpOut)
	tIns = append(tIns, ins)

	// split
	ins = NewLHSplitInstruction()
	ins.AddInput(cmpOut)

	ins.AddOutput(moving)
	ins.AddOutput(remaining)

	tIns = append(tIns, ins)

	// mix again

	ins = NewLHMixInstruction()

	ins.AddInput(remaining)
	ins.AddOutput(NewLHComponent())
	tIns = append(tIns, ins)

	tgraph, err := MakeTGraph(tIns)
	if err != nil {
		t.Error(err)
	}

	arrEq := func(ar1 []*LHInstruction, ar2 []*LHInstruction) bool {
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
