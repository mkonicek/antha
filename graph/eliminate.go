package graph

type EliminateOpt struct {
	Graph          Graph
	In             func(Node) bool // Should node be included
	KeepMultiEdges bool
}

// Return graph resulting from node elimination. Node elimination removes node
// n by adding edges (in(n), out(n)) for the product of incoming and outgoing
// neighbors.
func Eliminate(opt EliminateOpt) Graph {
	kmap := make(map[Node]bool)
	keep := func(n Node) bool {
		r, ok := kmap[n]
		if !ok {
			r = opt.In(n)
			kmap[n] = r
		}
		return r
	}

	ret := &qgraph{
		Outs: make(map[Node][]Node),
	}

	// Elimination can be quadratic. Reduce to linear on special case of trees
	// with reverse topological order.
	var nodes []Node
	if order, err := TopoSort(TopoSortOpt{
		Graph: opt.Graph,
	}); err != nil {
		for i, inum := 0, opt.Graph.NumNodes(); i < inum; i += 1 {
			nodes = append(nodes, opt.Graph.Node(i))
		}
	} else {
		for i := len(order) - 1; i >= 0; i -= 1 {
			nodes = append(nodes, order[i])
		}
	}

	ins := make(map[Node][]Node)
	for _, n := range nodes {
		for j, jnum := 0, opt.Graph.NumOuts(n); j < jnum; j += 1 {
			out := opt.Graph.Out(n, j)
			ins[out] = append(ins[out], n)
		}
	}

	outs := make(map[Node][]Node)
	for _, n := range nodes {
		if keep(n) {
			continue
		}
		for j, jnum := 0, opt.Graph.NumOuts(n); j < jnum; j += 1 {
			out := opt.Graph.Out(n, j)
			for _, in := range ins[n] {
				outs[in] = append(outs[in], out)
				ins[out] = append(ins[out], in)
			}
		}
	}

	// Filter out nodes
	for _, n := range nodes {
		if !keep(n) {
			continue
		}

		ret.Nodes = append(ret.Nodes, n)

		seen := make(map[Node]bool)
		for j, jnum := 0, opt.Graph.NumOuts(n); j < jnum; j += 1 {
			dst := opt.Graph.Out(n, j)
			if !keep(dst) {
				continue
			}
			seen[dst] = true
			ret.Outs[n] = append(ret.Outs[n], dst)
		}
		for _, dst := range outs[n] {
			if !keep(dst) {
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
