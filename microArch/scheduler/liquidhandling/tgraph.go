package liquidhandling

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/graph"
)

func MakeTGraph(inss []*wtype.LHInstruction) tGraph {
	edges := make(map[*wtype.LHInstruction][]*wtype.LHInstruction, len(inss))
	rcm := resultCmpMap(inss)

	for _, ins := range inss {
		edges[ins] = getEdges(ins, rcm)
	}

	return tGraph{Nodes: inss, Edges: edges}
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
func resultCmpMap(inss []*wtype.LHInstruction) map[string]*wtype.LHInstruction {
	res := make(map[string]*wtype.LHInstruction, len(inss))
	for _, ins := range inss {
		if ins.Type == wtype.LHIMIX {
			res[ins.Result.ID] = ins
		} else if ins.Type == wtype.LHIPRM {
			// we use passthrough instead
			for _, cmp := range ins.PassThrough {
				res[cmp.ID] = ins
			}
		}
	}

	return res
}

// inss maps result (i.e. component) IDs to instructions
func getEdges(n *wtype.LHInstruction, inss map[string]*wtype.LHInstruction) []*wtype.LHInstruction {
	ret := make([]*wtype.LHInstruction, 0, 1)

	// we make this backwards since it's easier to say where something's coming from than where
	// it's going to

	for _, cmp := range n.Components {
		// inss answers the question "which instruction made this?"
		// for samples we need to ask for the parent component
		var lhi *wtype.LHInstruction
		var ok bool
		if cmp.IsSample() {
			lhi, ok = inss[cmp.ParentID]
		} else {
			lhi, ok = inss[cmp.ID]
		}
		if ok {
			ret = append(ret, lhi)
		}
	}

	return ret
}

func uniqueIn(n *wtype.LHInstruction, s []*wtype.LHInstruction) bool {
	for _, n2 := range s {
		if n == n2 {
			return false
		}
	}

	return true
}
