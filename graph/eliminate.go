package graph

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
		for i, inum := 0, graph.NumNodes(); i < inum; i += 1 {
			nodes = append(nodes, graph.Node(i))
		}
	}

	return
}

func inNodes(graph Graph) map[Node][]Node {
	ins := make(map[Node][]Node)
	for i, inum := 0, graph.NumNodes(); i < inum; i += 1 {
		n := graph.Node(i)
		for j, jnum := 0, graph.NumOuts(n); j < jnum; j += 1 {
			out := graph.Out(n, j)
			ins[out] = append(ins[out], n)
		}
	}
	return ins
}

// Return graph resulting from node elimination. Node elimination removes node
// n by adding edges (in(n), out(n)) for the product of incoming and outgoing
// neighbors.
func Eliminate(opt EliminateOpt) Graph {
	// Cache nodes to keep
	kmap := make(map[Node]bool)
	for i, inum := 0, opt.Graph.NumNodes(); i < inum; i += 1 {
		n := opt.Graph.Node(i)
		kmap[n] = opt.In(n)
	}

	// Retarget ins of eliminated nodes to outs of eliminated nodes
	nodes := eliminationOrder(opt.Graph)
	ins := inNodes(opt.Graph)
	outs := make(map[Node][]Node)
	eliminated := make(map[Node]bool)
	for idx, n := range nodes {
		// Remove processed ins as we go
		if idx > 0 {
			delete(ins, nodes[idx-1])
		}

		if kmap[n] {
			continue
		}

		// Make sure to process outs given to us by previously eliminated nodes
		nouts := outs[n]
		for j, jnum := 0, opt.Graph.NumOuts(n); j < jnum; j += 1 {
			nouts = append(nouts, opt.Graph.Out(n, j))
		}

		for _, out := range nouts {
			if kmap[out] {
				// Normal case: link ins of eliminated node to outs of
				// eliminated node
				for _, in := range ins[n] {
					outs[in] = append(outs[in], out)
				}
				continue
			}

			if eliminated[out] {
				continue
			}

			// If we aren't keeping out and we haven't yet processed it, update
			// its ins to include our ins
			for _, in := range ins[n] {
				ins[out] = append(ins[out], in)
			}
		}
		eliminated[n] = true
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
		for j, jnum := 0, opt.Graph.NumOuts(n); j < jnum; j += 1 {
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
