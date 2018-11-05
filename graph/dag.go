package graph

import (
	"fmt"
)

type topoOrder struct {
	Graph Graph
	Order []Node  // Order in topographical sort
	Cycle []Node  // If not DAG, which nodes participate in a cycle
	black nodeSet // Fully processed nodes
	gray  nodeSet // Nodes currently being processed
}

// TODO(ddn): For performance, consider avoiding recursion

// Visit leaves, then nodes next to leaves, etc.
func (a *topoOrder) visit(n Node) {
	if a.gray[n] {
		if !a.black[n] {
			a.black[n] = true
			a.Cycle = append(a.Cycle, n)
		}
		return
	}
	if a.black[n] {
		return
	}
	a.gray[n] = true
	for i, num := 0, a.Graph.NumOuts(n); i < num; i++ {
		a.visit(a.Graph.Out(n, i))
	}
	delete(a.gray, n)
	if !a.black[n] {
		a.black[n] = true
		a.Order = append(a.Order, n)
	}
}

func (a *topoOrder) cycleError() error {
	if len(a.Cycle) == 0 {
		return nil
	}
	return fmt.Errorf("cycle containing %p", a.Cycle[0])
}

// Run topographic sort
func topoSort(opt TopoSortOpt) *topoOrder {
	g := opt.Graph
	if opt.NodeOrder != nil {
		g = makeSortedGraph(opt.Graph, opt.NodeOrder)
	}
	to := &topoOrder{
		Graph: g,
		black: make(nodeSet),
		gray:  make(nodeSet),
	}
	for i, num := 0, g.NumNodes(); i < num; i++ {
		n := g.Node(i)
		if _, seen := to.black[n]; !seen {
			to.visit(n)
		}
	}
	return to
}

// IsDag returns nil if graph is acyclic. If graph contains a cycle, return
// error.
func IsDag(g Graph) error {
	return topoSort(TopoSortOpt{Graph: g}).cycleError()
}

// A TopoSortOpt are options to TopoSort
type TopoSortOpt struct {
	Graph     Graph
	NodeOrder LessThan // Optional argument to ensure deterministic output
}

// TopoSort returns topological sort of graph. If edge (a, b) is in g, then b <
// a in the resulting order.  Returns an error if graph contains a cycle.
func TopoSort(opt TopoSortOpt) ([]Node, error) {
	to := topoSort(opt)
	if err := to.cycleError(); err != nil {
		return nil, err
	}
	return to.Order, nil
}

// TransitiveReduction computes transitive reduction of a graph. Relatively
// expensive operation: O(nm).
func TransitiveReduction(graph Graph) (Graph, error) {
	ret := &qgraph{
		Outs: make(map[Node][]Node),
	}

	if err := IsDag(graph); err != nil {
		// TODO(ddn): transitive reductions exist for cyclic graphs but we just
		// can't use SSSP to find them
		return nil, fmt.Errorf("not yet implemented: %s", err)
	}

	dag := Schedule(graph)
	for len(dag.Roots) > 0 {
		for _, root := range dag.Roots {
			// In DAG, solving shortest path with -w() is the solution to the
			// longest path problem
			dist := ShortestPath(ShortestPathOpt{
				Graph:   graph,
				Sources: []Node{root},
				Weight: func(x, y Node) int64 {
					return -1
				},
			})

			ret.Nodes = append(ret.Nodes, root)
			for i, inum := 0, graph.NumOuts(root); i < inum; i++ {
				dst := graph.Out(root, i)
				if dist[dst] == -1 {
					ret.Outs[root] = append(ret.Outs[root], dst)
				}
			}
		}

		var next []Node
		for _, root := range dag.Roots {
			next = append(next, dag.Visit(root)...)
		}
		dag.Roots = next
	}

	return ret, nil
}
