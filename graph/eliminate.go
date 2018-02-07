package graph

// An EliminateOpt are a set of options to Eliminate
type EliminateOpt struct {
	Graph          Graph
	In             func(Node) bool // Should node be included
	KeepMultiEdges bool
}

// Elimination can be quadratic. Reduce to average cost by processing nodes in
// topological order.
func eliminationOrder(graph Graph) (nodes []Node) {
	order, err := TopoSort(TopoSortOpt{
		Graph: graph,
	})
	if err == nil {
		nodes = order
	} else {
		for i, inum := 0, graph.NumNodes(); i < inum; i++ {
			nodes = append(nodes, graph.Node(i))
		}
	}

	return
}

func inNodes(graph Graph) map[Node][]Node {
	ins := make(map[Node][]Node)
	for i, inum := 0, graph.NumNodes(); i < inum; i++ {
		n := graph.Node(i)
		for j, jnum := 0, graph.NumOuts(n); j < jnum; j++ {
			out := graph.Out(n, j)
			ins[out] = append(ins[out], n)
		}
	}
	return ins
}

func outNodes(graph Graph) map[Node][]Node {
	outs := make(map[Node][]Node)
	for i, inum := 0, graph.NumNodes(); i < inum; i++ {
		n := graph.Node(i)
		for j, jnum := 0, graph.NumOuts(n); j < jnum; j++ {
			out := graph.Out(n, j)
			outs[n] = append(outs[n], out)
		}
	}
	return outs
}

// Eliminate returns the graph resulting from node elimination. Node
// elimination removes node n by adding edges (in(n), out(n)) for the product
// of incoming and outgoing neighbors.
func Eliminate(opt EliminateOpt) Graph {
	// Cache nodes to keep
	kmap := make(map[Node]bool)
	for i, inum := 0, opt.Graph.NumNodes(); i < inum; i++ {
		n := opt.Graph.Node(i)
		kmap[n] = opt.In(n)
	}

	// Retarget ins of eliminated nodes to outs of eliminated nodes
	nodes := eliminationOrder(opt.Graph)
	ins := inNodes(opt.Graph)
	outs := outNodes(opt.Graph)

	for _, n := range nodes {
		if kmap[n] {
			continue
		}

		for _, out := range outs[n] {
			for _, in := range ins[n] {
				if n == in {
					continue
				}

				outs[in] = append(outs[in], out)
			}

			for _, in := range ins[n] {
				if n == out {
					continue
				}
				ins[out] = append(ins[out], in)
			}
		}
	}

	// Create eliminated graph
	ret := &qgraph{
		Outs: make(map[Node][]Node),
	}

	// Filter out nodes
	for _, n := range nodes {
		if !kmap[n] {
			continue
		}

		ret.Nodes = append(ret.Nodes, n)

		seen := make(map[Node]bool)
		for j, jnum := 0, opt.Graph.NumOuts(n); j < jnum; j++ {
			dst := opt.Graph.Out(n, j)
			if !kmap[dst] {
				continue
			}
			seen[dst] = true
			ret.Outs[n] = append(ret.Outs[n], dst)
		}
		for _, dst := range outs[n] {
			if !kmap[dst] {
				continue
			}
			if !opt.KeepMultiEdges && seen[dst] {
				continue
			}
			ret.Outs[n] = append(ret.Outs[n], dst)
			seen[dst] = true
		}
	}

	return ret
}
