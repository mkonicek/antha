package liquidhandling

import (
	"github.com/pkg/errors"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/graph"
)

func MakeTGraph(inss []*wtype.LHInstruction) (tGraph, error) {
	edges := make(map[*wtype.LHInstruction][]*wtype.LHInstruction, len(inss))
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
	Nodes []*wtype.LHInstruction
	Edges map[*wtype.LHInstruction][]*wtype.LHInstruction
}

func (tg tGraph) NumNodes() int {
	return len(tg.Nodes)
}

func (tg tGraph) Node(i int) graph.Node {
	return graph.Node(tg.Nodes[i])
}

func (tg tGraph) NumOuts(n graph.Node) int {
	return len(tg.Edges[n.(*wtype.LHInstruction)])
}

func (tg tGraph) Out(n graph.Node, i int) graph.Node {
	return graph.Node(tg.Edges[n.(*wtype.LHInstruction)][i])
}

func (tg *tGraph) Add(n *wtype.LHInstruction, edges []*wtype.LHInstruction) {
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
func resultCmpMap(inss []*wtype.LHInstruction) map[string]*wtype.LHInstruction {
	res := make(map[string]*wtype.LHInstruction, len(inss))
	for _, ins := range inss {
		if ins.Type == wtype.LHIMIX {
			res[ins.Results[0].ID] = ins
		} else if ins.Type == wtype.LHIPRM {
			// we use passthrough instead
			for _, cmp := range ins.PassThrough {
				res[cmp.ID] = ins
			}
		} else if ins.Type == wtype.LHISPL {
			// Splits need to go after the use of result 0
			// and before the use of result 1
			res[ins.Results[1].ID] = ins
		}
	}

	return res
}

func cmpInsMap(inss []*wtype.LHInstruction) map[string]*wtype.LHInstruction {
	res := make(map[string]*wtype.LHInstruction, len(inss))
	for _, ins := range inss {
		if ins.Type == wtype.LHIMIX {
			for _, c := range ins.Components {
				res[c.ID] = ins
			}
		} else if ins.Type == wtype.LHIPRM {
			// we use passthrough instead
			for ID := range ins.PassThrough {
				res[ID] = ins
			}
		}
	}

	return res

}

// inss maps result (i.e. component) IDs to instructions
func getEdges(n *wtype.LHInstruction, resultMap, cmpMap map[string]*wtype.LHInstruction) ([]*wtype.LHInstruction, error) {
	ret := make([]*wtype.LHInstruction, 0, 1)

	// don't make cycles containing split instructions

	if n.Type == wtype.LHISPL {
		// cmpMap answers the question "which instruction *uses* this?"
		cmp := n.Results[0]

		var lhi *wtype.LHInstruction
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

	for _, cmp := range n.Components {
		// resultMap answers the question "which instruction *makes* this?"
		// for samples we need to ask for the parent component
		var lhi *wtype.LHInstruction
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

func uniqueIn(n *wtype.LHInstruction, s []*wtype.LHInstruction) bool {
	for _, n2 := range s {
		if n == n2 {
			return false
		}
	}

	return true
}
