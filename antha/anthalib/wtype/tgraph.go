package wtype

import (
	"github.com/antha-lang/antha/graph"
	"github.com/pkg/errors"
)

func MakeTGraph(inss []*LHInstruction) (tGraph, error) {
	edges := make(map[*LHInstruction][]*LHInstruction, len(inss))
	rcm := resultCmpMap(inss)
	ccm := cmpInsMap(inss)

	var err error
	for _, ins := range inss {
		edges[ins], err = getEdges(ins, rcm, ccm)
		if err != nil {
			return tGraph{}, err
		}
	}

	return tGraph{Nodes: inss, Edges: edges}, nil
}

type tGraph struct {
	Nodes []*LHInstruction
	Edges map[*LHInstruction][]*LHInstruction
}

func (tg tGraph) NumNodes() int {
	return len(tg.Nodes)
}

func (tg tGraph) Node(i int) graph.Node {
	return graph.Node(tg.Nodes[i])
}

func (tg tGraph) NumOuts(n graph.Node) int {
	return len(tg.Edges[n.(*LHInstruction)])
}

func (tg tGraph) Out(n graph.Node, i int) graph.Node {
	return graph.Node(tg.Edges[n.(*LHInstruction)][i])
}

func (tg *tGraph) Add(n *LHInstruction, edges []*LHInstruction) {
	// enforce uniqueness

	if !uniqueIn(n, tg.Nodes) {
		return
	} else {
		tg.Nodes = append(tg.Nodes, n)
		tg.Edges[n] = edges
	}
}

// for mixes this is 1:1
// but prompts may be aggregated first
// splits are more complex
func resultCmpMap(inss []*LHInstruction) map[string]*LHInstruction {
	res := make(map[string]*LHInstruction, len(inss))
	for _, ins := range inss {
		if ins.Type == LHIMIX {
			res[ins.Outputs[0].ID] = ins
		} else if ins.Type == LHIPRM {
			for _, cmp := range ins.Outputs {
				res[cmp.ID] = ins
			}
		} else if ins.Type == LHISPL {
			// Splits need to go after the use of result 0
			// and before the use of result 1
			res[ins.Outputs[1].ID] = ins
		}
	}

	return res
}

func cmpInsMap(inss []*LHInstruction) map[string]*LHInstruction {
	res := make(map[string]*LHInstruction, len(inss))
	for _, ins := range inss {
		if ins.Type == LHIMIX {
			for _, c := range ins.Inputs {
				res[c.ID] = ins
			}
		}
	}

	return res

}

// inss maps result (i.e. component) IDs to instructions
func getEdges(n *LHInstruction, resultMap, cmpMap map[string]*LHInstruction) ([]*LHInstruction, error) {
	ret := make([]*LHInstruction, 0, 1)

	// don't make cycles containing split instructions

	if n.Type == LHISPL {
		// cmpMap answers the question "which instruction *uses* this?"
		cmp := n.Outputs[0]

		var lhi *LHInstruction
		var ok bool

		lhi, ok = cmpMap[cmp.ID]

		if ok {
			ret = append(ret, lhi)
		} else {
			return nil, errors.Errorf("SplitSample called without use of component. Moving components must be moved using Mix. Component name %s ID %s", cmp.CName, cmp.ID)
		}
		return ret, nil
	}

	// we make this backwards since it's easier to say where something's coming from than where
	// it's going to

	for _, cmp := range n.Inputs {
		// resultMap answers the question "which instruction *makes* this?"
		// for samples we need to ask for the parent component
		var lhi *LHInstruction
		var ok bool
		if cmp.IsSample() {
			lhi, ok = resultMap[cmp.ParentID]
		} else {
			lhi, ok = resultMap[cmp.ID]
		}
		if ok {
			ret = append(ret, lhi)
		}
	}

	return ret, nil
}

func uniqueIn(n *LHInstruction, s []*LHInstruction) bool {
	for _, n2 := range s {
		if n == n2 {
			return false
		}
	}

	return true
}
