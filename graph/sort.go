package graph

import (
	"sort"
)

type nodeSliceSorter struct {
	Nodes    []Node
	LessThan LessThan
}

func (a *nodeSliceSorter) Less(i, j int) bool {
	return a.LessThan(a.Nodes[i], a.Nodes[j])
}

func (a *nodeSliceSorter) Swap(i, j int) {
	a.Nodes[i], a.Nodes[j] = a.Nodes[j], a.Nodes[i]
}

func (a *nodeSliceSorter) Len() int {
	return len(a.Nodes)
}

// Create a graph that has a deterministic iteration order. Package-private
// because Graph interface should not imply a particular iteration order.
func makeSortedGraph(g Graph, less LessThan) *qgraph {
	ret := &qgraph{
		Outs: make(map[Node][]Node),
	}

	for i, inum := 0, g.NumNodes(); i < inum; i++ {
		n := g.Node(i)
		ret.Nodes = append(ret.Nodes, n)
		for j, jnum := 0, g.NumOuts(n); j < jnum; j++ {
			ret.Outs[n] = append(ret.Outs[n], g.Out(n, j))
		}
		sort.Sort(&nodeSliceSorter{
			Nodes:    ret.Outs[n],
			LessThan: less,
		})
	}

	sort.Sort(&nodeSliceSorter{
		Nodes:    ret.Nodes,
		LessThan: less,
	})

	return ret
}
