package graph

// Merge n and src into dst, modifying dst. Assumes src != dst.
func mergeInto(n Node, src, dst map[Node]bool) {
	dst[n] = true
	for k := range src {
		dst[k] = true
	}
}

// Compute reachability over graph. A reaches B if there is a path from A to B.
func Reaches(g Graph) Graph {
	reaches := make(map[Node]map[Node]bool)

	// Get int set for node i
	get := func(n Node) map[Node]bool {
		m, ok := reaches[n]
		if ok {
			return m
		}
		m = make(map[Node]bool)
		reaches[n] = m
		return m
	}

	// Initialize
	changed := make(map[Node]bool)
	for i, inum := 0, g.NumNodes(); i < inum; i++ {
		node := g.Node(i)
		changed[node] = true
	}

	// Fixpoint
	for len(changed) != 0 {
		next := make(map[Node]bool)

		for node := range changed {
			src := get(node)
			for j, jnum := 0, g.NumOuts(node); j < jnum; j++ {
				neigh := g.Out(node, j)
				dst := get(neigh)
				prev := len(dst)
				if node != neigh {
					mergeInto(node, src, dst)
				} else {
					dst[node] = true
				}
				if len(dst) != prev {
					next[neigh] = true
				}
			}
		}

		changed = next
	}

	// Reaches is all the nodes that reach n
	ret := &qgraph{
		Outs: make(map[Node][]Node),
	}

	for i, inum := 0, g.NumNodes(); i < inum; i++ {
		node := g.Node(i)
		ret.Nodes = append(ret.Nodes, node)
		for r := range reaches[node] {
			ret.Outs[node] = append(ret.Outs[node], r)
		}
	}

	// Reverse to get nodes that n can reach
	return Reverse(ret)
}
